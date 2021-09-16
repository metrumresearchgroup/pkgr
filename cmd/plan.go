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
	"path/filepath"
	"runtime"

	"github.com/spf13/afero"
	"github.com/thoas/go-funk"

	"github.com/metrumresearchgroup/pkgr/rollback"

	"github.com/metrumresearchgroup/pkgr/desc"
	"github.com/metrumresearchgroup/pkgr/pacman"

	"os"

	"strings"
	"time"

	"github.com/metrumresearchgroup/pkgr/rcmd"

	"github.com/sajari/fuzzy"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/metrumresearchgroup/pkgr/configlib"
	"github.com/metrumresearchgroup/pkgr/cran"
	"github.com/metrumresearchgroup/pkgr/gpsr"
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
	log.Infof("Installation would launch %v workers\n", getWorkerCount(cfg.Threads, runtime.GOMAXPROCS(0)))
	rs := rcmd.NewRSettings(cfg.RPath)
	rVersion := rcmd.GetRVersion(&rs)
	log.Infoln("R Version " + rVersion.ToFullString())
	log.Infoln("OS Platform " + rs.Platform)
	_, ip, _ := planInstall(rVersion, true)
	if viper.GetBool("show-deps") {
		for pkg, deps := range ip.DepDb {
			fmt.Println("-----------  ", pkg, "   ------------")
			fmt.Println(deps)
		}
	}
	return nil
}

func planInstall(rv cran.RVersion, exitOnMissing bool) (*cran.PkgNexus, gpsr.InstallPlan, rollback.RollbackPlan) {
	startTime := time.Now()

	//Check library existence
	libraryExists, err := afero.DirExists(fs, cfg.Library)

	if err != nil {
		log.WithFields(log.Fields{
			"library": cfg.Library,
			"error":   err,
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
			repo, _ := configlib.GetRepoCustomizationByName(nm, cfg.Customizations)
			// for now no need to check if customization exists as the repo will have a default empty string
			// regardless so no additional logic needed
			repos = append(repos, cran.RepoURL{Name: nm, URL: url, Suffix: repo.RepoSuffix})
		}
	}
	st := cran.DefaultType()
	cic := cran.NewInstallConfig()
	for _, repoSlice := range cfg.Customizations.Repos {
		for rn, val := range repoSlice {
			rc := cran.RepoConfig{}
			if strings.EqualFold(val.RepoType, "MPN") {
				rc.RepoType = cran.MPN
				rc.DefaultSourceType = cran.Binary
			}
			if strings.EqualFold(val.RepoType, "RSPM") {
				rc.RepoType = cran.RSPM
			}
			if strings.EqualFold(val.Type, "binary") {
				rc.DefaultSourceType = cran.Binary
			}
			if strings.EqualFold(val.Type, "source") {
				rc.DefaultSourceType = cran.Source
			}
			if val.RepoSuffix != "" {
				rc.RepoSuffix = val.RepoSuffix
			}
			cic.Repos[rn] = rc
		}
	}
	pkgNexus, err := cran.NewPkgDb(repos, st, cic, rv, cfg.NoSecure)
	if err != nil {
		log.Panicln("error getting pkgdb ", err)
	}
	log.Infoln("Default package installation type: ", st.String())
	for _, db := range pkgNexus.Db {
		log.Infoln(fmt.Sprintf("%v:%v (binary:source) packages available in for %s from %s", len(db.DescriptionsBySourceType[cran.Binary]), len(db.DescriptionsBySourceType[cran.Source]), db.Repo.Name, db.Repo.URL))
		for _, pkg := range cfg.IgnorePackages {
			// to "skip" packages, we'll just completely nuke them from the pkgdb so they'll never even come up in the plan
			// this is probably overly hacky
			log.Debugln("ignoring by deleting pkg: ", pkg)
			delete(db.DescriptionsBySourceType[cran.Binary], pkg)
			delete(db.DescriptionsBySourceType[cran.Source], pkg)
		}
	}
	log.Infoln("Package installation cache directory: ", userCache(cfg.Cache))
	log.Infoln("Database cache directory: ", filepath.Dir(pkgNexus.Db[0].GetRepoDbCacheFilePath(rv.ToFullString())))

	dependencyConfigurations := gpsr.NewDefaultInstallDeps()
	dependencyConfigurations.Default.NoRecommended = cfg.NoRecommended
	configlib.SetPlanCustomizations(cfg, dependencyConfigurations, pkgNexus)

	// Set tarball dependencies as user-packages, for convenience.
	var tarballDescriptions []desc.Desc
	var unpackedTarballPkgs map[string]gpsr.AdditionalPkg

	if len(cfg.Tarballs) > 0 {
		tarballDescriptions, unpackedTarballPkgs = unpackTarballs(fs, cfg.Tarballs, cfg.Cache)
		for _, tarballDesc := range tarballDescriptions {
			tarballDeps := tarballDesc.GetCombinedDependencies(false)
			for _, d := range tarballDeps {
				if !funk.Contains(cfg.Packages, d.Name) {
					cfg.Packages = append(cfg.Packages, d.Name)
				}
			}
		}
	}
	// end tarball deps

	// Set dependencies from Descriptions files as user-packages, for convenience.
	var descDescriptions []desc.Desc

	if len(cfg.Descriptions) > 0 {
		descDescriptions = unpackDescriptions(fs, cfg.Descriptions)
		for _, desc := range descDescriptions {
			descDeps := desc.GetCombinedDependencies(true)
			for _, d := range descDeps {
				if !funk.Contains(cfg.Packages, d.Name) {
					cfg.Packages = append(cfg.Packages, d.Name)
				}
			}
		}
	}
	// end Descriptions deps

	cfg.Packages = removeBasePackages(cfg.Packages)

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
		//tarballDescriptions,
		dependencyConfigurations,
		pkgNexus,
		!cfg.NoUpdate,
		libraryExists,
		cfg.NoRecommended,
	)

	installPlan.AdditionalPackageSources = unpackedTarballPkgs

	logAdditionalPackageOrigins(installPlan.AdditionalPackageSources)

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
		if cfg.NoUpdate {
			log.WithFields(updateLogFields).Warn("outdated package found")
		} else {
			log.WithFields(updateLogFields).Info("package will be updated")
			pkgsToUpdateCount = len(installPlan.OutdatedPackages)
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
	installSources["tarballs"] = len(installPlan.AdditionalPackageSources)
	fields := make(log.Fields)
	for k, v := range installSources {
		fields[k] = v
	}

	log.WithFields(fields).Info("package installation sources")

	log.WithFields(log.Fields{
		"to_install": toInstall,
		"to_update":  pkgsToUpdateCount,
	}).Info("package installation plan")
	log.Infof("Library path to install packages: %s\n", cfg.Library)

	if toInstall > 0 && toInstall != totalPackagesRequired {
		// log which packages to install, but only if doing an incremental install
		for _, pn := range pkgs {
			//_, isAdditionalPkg := installPlan.AdditionalPackageSources[pn]
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

// Removes any "base" packages from the given list.
func removeBasePackages(pkgList []string) []string {
	var nonbasePkgList []string
	for _, p := range pkgList {
		pType, found := gpsr.DefaultPackages[p]
		if !found || pType != "base" {
			nonbasePkgList = append(nonbasePkgList, p)
		} else {
			log.WithFields(log.Fields{
				"pkg": p,
			}).Warn("removing base package from user-defined package list")
		}
	}
	return nonbasePkgList
}

func logAdditionalPackageOrigins(additionalPackages map[string]gpsr.AdditionalPkg) {
	for pkg, details := range additionalPackages {
		log.WithFields(log.Fields{
			"pkg":          pkg,
			"origin":       details.OriginPath,
			"method":       details.Type,
			"install_from": details.InstallPath,
		}).Debug("additional installation set")
	}
}

func logUserPackageRepos(packageDownloads []cran.PkgDl) {
	for _, pkg := range packageDownloads {
		log.WithFields(log.Fields{
			"pkg":          pkg.Package.Package,
			"repo":         pkg.Config.Repo.Name,
			"type":         pkg.Config.Type,
			"version":      pkg.Package.Version,
			"relationship": "user_defined",
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
