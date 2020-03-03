// Copyright © 2018 Devin Pastoor <devin.pastoor@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"github.com/spf13/afero"
	"runtime"

	"github.com/metrumresearchgroup/pkgr/rollback"

	"github.com/metrumresearchgroup/pkgr/desc"
	"github.com/metrumresearchgroup/pkgr/pacman"

	"os"

	"strings"
	"time"

	"github.com/metrumresearchgroup/pkgr/rcmd"

	"github.com/metrumresearchgroup/pkgr/configlib"
	"github.com/metrumresearchgroup/pkgr/cran"
	"github.com/metrumresearchgroup/pkgr/gpsr"
	"github.com/sajari/fuzzy"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	funk "github.com/thoas/go-funk"
)

// planCmd shows the install plan
var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "plan a full installation",
	Long: `
	see the plan for an install
 `,
	RunE: plan,
}

func init() {
	planCmd.PersistentFlags().Bool("show-deps", false, "show the (required) dependencies for each package")
	viper.BindPFlag("show-deps", planCmd.PersistentFlags().Lookup("show-deps"))
	RootCmd.AddCommand(planCmd)
}

func plan(cmd *cobra.Command, args []string) error {
	log.Infof("Installation would launch %v workers\n", getWorkerCount(viper.GetInt("threads"), runtime.NumCPU()))
	rs := rcmd.NewRSettings(cfg.RPath)
	rVersion := rcmd.GetRVersion(&rs)
	log.Infoln("R Version " + rVersion.ToFullString())
	log.Debugln("OS Platform " + rs.Platform)
	_, ip, _ := planInstall(rVersion, cran.DefaultType(), true)
	if viper.GetBool("show-deps") {
		for pkg, deps := range ip.DepDb {
			fmt.Println("-----------  ", pkg, "   ------------")
			fmt.Println(deps)
		}
	}
	return nil
}

func planInstall(rv cran.RVersion, st cran.SourceType, exitOnMissing bool) (*cran.PkgNexus, gpsr.InstallPlan, rollback.RollbackPlan) {
	startTime := time.Now()

	//Check library existence
	libraryExists, err := afero.DirExists(fs, cfg.Library)

	if err != nil {
		log.WithFields(log.Fields{
			"library": cfg.Library,
			"error" : err,
		}).Error("unexpected error when checking existence of library")
	}

	if !libraryExists && cfg.Strict {
		log.WithFields(log.Fields{
			"library": cfg.Library,
		}).Error("library directory must exist before running pkgr in strict mode")
	}

	var installedPackageNames []string
	var installedPackages map[string]desc.Desc
	var whereInstalledFrom pacman.InstalledFromPkgs

	if libraryExists {
		installedPackages = pacman.GetPriorInstalledPackages(fs, cfg.Library)
		installedPackageNames = extractNamesFromDesc(installedPackages)
		log.WithField("count", len(installedPackages)).Info("found installed packages")
		whereInstalledFrom = pacman.GetInstallers(installedPackages)
		notPkgr := whereInstalledFrom.NotFromPkgr()
		if len(notPkgr) > 0 {
			// TODO: should this say "prior installed packages" not ...
			log.WithFields(log.Fields{
				"packages": notPkgr,
			}).Warn("Packages not installed by pkgr")
		}
	} else {
		log.WithFields(log.Fields{
			"path": cfg.Library,
		}).Info("Package Library will be created")
		//fs.Create(cfg.Library)
		//fs.Chmod(cfg.Library, 0755)
	}

	var repos []cran.RepoURL
	for _, r := range cfg.Repos {
		for nm, url := range r {
			repos = append(repos, cran.RepoURL{Name: nm, URL: url})
		}
	}
	cic := cran.NewInstallConfig()
	for rn, val := range cfg.Customizations.Repos {
		if strings.EqualFold(val.Type, "binary") {
			cic.Repos[rn] = cran.RepoConfig{DefaultSourceType: cran.Binary}
		}
		if strings.EqualFold(val.Type, "source") {
			cic.Repos[rn] = cran.RepoConfig{DefaultSourceType: cran.Source}
		}
	}
	pkgNexus, err := cran.NewPkgDb(repos, st, cic, rv)
	if err != nil {
		log.Panicln("error getting pkgdb ", err)
	}
	log.Infoln("Default package installation type: ", st.String())
	for _, db := range pkgNexus.Db {
		log.Infoln(fmt.Sprintf("%v:%v (binary:source) packages available in for %s from %s", len(db.DescriptionsBySourceType[st]), len(db.DescriptionsBySourceType[cran.Source]), db.Repo.Name, db.Repo.URL))
	}

	dependencyConfigurations := gpsr.NewDefaultInstallDeps()
	configlib.SetPlanCustomizations(cfg, dependencyConfigurations, pkgNexus)

	availableUserPackages := pkgNexus.GetPackages(cfg.Packages)
	if len(availableUserPackages.Missing) > 0 {
		log.Errorln("missing packages: ", availableUserPackages.Missing)
		model := fuzzy.NewModel()

		// For testing only, this is not advisable on production
		model.SetThreshold(1)

		// This expands the distance searched, but costs more resources (memory and time).
		// For spell checking, "2" is typically enough, for query suggestions this can be higher
		model.SetDepth(1)
		pkgs := pkgNexus.GetAllPkgsByName()
		model.Train(pkgs)
		for _, mp := range availableUserPackages.Missing {
			log.Warnln("did you mean one of: ", model.Suggestions(mp, false))
		}
		if exitOnMissing {
			os.Exit(1)
		} else {
			return pkgNexus, gpsr.InstallPlan{}, rollback.RollbackPlan{}
		}
	}
	logUserPackageRepos(availableUserPackages.Packages)
	installPlan, err := gpsr.ResolveInstallationReqs(
		cfg.Packages,
		installedPackages,
		dependencyConfigurations,
		pkgNexus,
		cfg.Update,
		libraryExists)
	rollbackPlan := rollback.CreateRollbackPlan(cfg.Library, installPlan, installedPackages)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	logDependencyRepos(installPlan.PackageDownloads)

	pkgs := installPlan.GetAllPackages()

	pkgsToUpdateCount := 0
	for _, p := range installPlan.OutdatedPackages {
		updateLogFields := log.Fields{
			"pkg":               p.Package,
			"installed_version": p.OldVersion,
			"update_version":    p.NewVersion,
		}
		if viper.GetBool("update") {
			log.WithFields(updateLogFields).Info("package will be updated")
			pkgsToUpdateCount = len(installPlan.OutdatedPackages)
		} else {
			log.WithFields(updateLogFields).Warn("outdated package found")
		}
	}

	totalPackagesRequired := len(pkgs)
	toInstall := installPlan.GetNumPackagesToInstall()
	log.WithFields(log.Fields{
		"total_packages_required": totalPackagesRequired,
		"installed":               len(installedPackages),
		"outdated":                len(installPlan.OutdatedPackages),
		"not_from_pkgr":           len(whereInstalledFrom.NotFromPkgr()),
	}).Info("package installation status")

	installSources := make(map[string]int)
	for _, pkgdl := range installPlan.PackageDownloads {
		_, rn := pkgdl.PkgAndRepoNames()
		installSources[rn]++
	}
	fields := make(log.Fields)
	for k, v := range installSources {
		fields[k] = v
	}
	log.WithFields(fields).Info("package installation sources")

	log.WithFields(log.Fields{
		"to_install": toInstall,
		"to_update":  pkgsToUpdateCount,
	}).Info("package installation plan")

	if toInstall > 0 && toInstall != totalPackagesRequired {
		// log which packages to install, but only if doing an incremental install
		for _, pn := range pkgs {
			if !funk.ContainsString(installedPackageNames, pn) {
				pkgDesc, cfg, _ := pkgNexus.GetPackage(pn)
				log.WithFields(log.Fields{
					"package": pkgDesc.Package,
					"version": pkgDesc.Version,
					"repo":    cfg.Repo.Name,
					"type":    cfg.Type,
				}).Info("to install")
			}
		}
	}

	log.Infoln("resolution time", time.Since(startTime))
	return pkgNexus, installPlan, rollbackPlan
}

func logUserPackageRepos(packageDownloads []cran.PkgDl) {
	for _, pkg := range packageDownloads {
		log.WithFields(log.Fields{
			"pkg":          pkg.Package.Package,
			"repo":         pkg.Config.Repo.Name,
			"type":         pkg.Config.Type,
			"version":      pkg.Package.Version,
			"relationship": "user package",
		}).Debug("package repository set")
	}
}

func logDependencyRepos(dependencyDownloads []cran.PkgDl) {
	for _, pkgToDownload := range dependencyDownloads {
		pkg := pkgToDownload.Package.Package

		if !stringInSlice(pkg, cfg.Packages) {
			log.WithFields(log.Fields{
				"pkg":          pkgToDownload.Package.Package,
				"repo":         pkgToDownload.Config.Repo.Name,
				"type":         pkgToDownload.Config.Type,
				"version":      pkgToDownload.Package.Version,
				"relationship": "dependency",
			}).Debug("package repository set")
		}
	}
}

func extractNamesFromDesc(installedPackages map[string]desc.Desc) []string {
	var installedPackageNames []string
	for key := range installedPackages {
		installedPackageNames = append(installedPackageNames, key)
	}
	return installedPackageNames
}
