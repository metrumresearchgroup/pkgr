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
	"github.com/metrumresearchgroup/pkgr/gpsr"
	"path/filepath"
	"time"

	"github.com/metrumresearchgroup/pkgr/cran"
	"github.com/metrumresearchgroup/pkgr/logger"
	"github.com/metrumresearchgroup/pkgr/rcmd"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)


var updateArgument bool

// installCmd represents the R CMD install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install a package",
	Long: `
	install a package
 `,
	RunE: rInstall,
}

func init() {
	installCmd.Flags().BoolVar(&updateArgument, "update", false, "Update outdated packages during installation.")
	RootCmd.AddCommand(installCmd)
}

func rInstall(cmd *cobra.Command, args []string) error {

	// Initialize log and start time.
	initInstallLog()
	startTime := time.Now()

	// Initialize objects to hold R settings and metadata.
	rSettings := rcmd.NewRSettings()
	rVersion := rcmd.GetRVersion(&rSettings)
	log.Infoln("R Version " + rVersion.ToFullString())

	// Get master object containing the packages available in each repository (pkgNexus),
	//  as well as a master install plan to guide our process.
	pkgNexus, installPlan := planInstall(rVersion)

	//Prepare our environment to update outdated packages if the "--update" flag is set.
	var packageUpdateAttempts []UpdateAttempt
	if updateArgument {
		log.Info("update argument passed. staging packages for update...")
		packageUpdateAttempts = preparePackagesForUpdate(fs, cfg.Library, installPlan.OutdatedPackages)
	}

	// Create a list of package download objects using our install plan and our "nexus" object.
	pkgsToDownload := getPackagesToDownload(installPlan, pkgNexus)

	// Retrieve a cache to store any packages we need to download for the install.
	packageCache := rcmd.NewPackageCache(userCache(cfg.Cache), false)

	//Create a pkgMap object, which helps us with parallel downloads (?)
	pkgMap, err := cran.DownloadPackages(fs, pkgsToDownload, packageCache.BaseDir, rVersion)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	//Set the arguments to be passed in to the R Package Installer
	pkgInstallArgs := rcmd.NewDefaultInstallArgs()
	pkgInstallArgs.Library, _ = filepath.Abs(cfg.Library)

	// Get the number of workers.
	// leave at least 1 thread open for coordination, given more than 2 threads available.
	// if only 2 available, will let the OS hypervisor coordinate some else would drop the
	// install time too much for the little bit of additional coordination going on.
	nworkers := getWorkerCount()

	// Process any customizations set in the yaml config file for individual packages.
	// TODO: Refactor this into its own method.
	pkgCustomizations := cfg.Customizations.Packages
	for n, v := range pkgCustomizations {
		if v.Env != nil {
			rSettings.PkgEnvVars[n] = v.Env
		}
	}

	//
	// Install the packages
	//
	err = rcmd.InstallPackagePlan(fs, installPlan, pkgMap, packageCache, pkgInstallArgs, rSettings, rcmd.ExecSettings{}, nworkers)
	if err != nil {
		fmt.Println("failed package install")
		fmt.Println(err)
	}

	// After package installation, fix any problems that occurred during reinstallation of
	//  packages that were to be updated.
	restoreUnupdatedPackages(fs, packageUpdateAttempts)

	fmt.Println("duration:", time.Since(startTime))
	return nil
}

func initInstallLog() {
	//Init install-specific log, if one has been set. This overwrites the default log.
	if cfg.Logging.Install != "" {
		logger.AddLogFile(cfg.Logging.Install, cfg.Logging.Overwrite)
	} else {
		logger.AddLogFile(cfg.Logging.All, cfg.Logging.Overwrite)
	}
}

func getPackagesToDownload(installPlan gpsr.InstallPlan, pkgNexus *cran.PkgNexus) []cran.PkgDl {
	var toDl []cran.PkgDl
	// starting packages
	for _, p := range installPlan.StartingPackages {
		pkg, cfg, _ := pkgNexus.GetPackage(p)
		toDl = append(toDl, cran.PkgDl{Package: pkg, Config: cfg})
	}
	// all other packages
	for p := range installPlan.DepDb {
		pkg, cfg, _ := pkgNexus.GetPackage(p)
		toDl = append(toDl, cran.PkgDl{Package: pkg, Config: cfg})
	}
	return toDl
}
