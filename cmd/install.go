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
	"github.com/metrumresearchgroup/pkgr/gpsr"
	"github.com/metrumresearchgroup/pkgr/rollback"
	"path/filepath"
	"runtime"
	"time"

	"github.com/metrumresearchgroup/pkgr/configlib"
	"github.com/metrumresearchgroup/pkgr/cran"
	"github.com/metrumresearchgroup/pkgr/logger"
	"github.com/metrumresearchgroup/pkgr/rcmd"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

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
	RootCmd.AddCommand(installCmd)
}

func rInstall(cmd *cobra.Command, args []string) error {

	// Initialize log and start time.
	initInstallLog()
	startTime := time.Now()
	rSettings := rcmd.NewRSettings(cfg.RPath)
	rVersion := rcmd.GetRVersion(&rSettings)
	log.Infoln("R Version " + rVersion.ToFullString())
	// most people should know what platform they are on
	log.Debugln("OS Platform " + rSettings.Platform)

	// Get master object containing the packages available in each repository (pkgNexus),
	//  as well as a master install plan to guide our process.
	_, installPlan, rollbackPlan := planInstall(rVersion, true)

	if installPlan.CreateLibrary {
		if cfg.Strict {
			log.WithFields(log.Fields{
				"library": cfg.Library,
			}).Fatal("library directory must exist before running pkgr in strict mode -- halting execution")
		}

		err := fs.MkdirAll(cfg.Library, 0755)
		if err != nil {
			log.WithFields(log.Fields{
				"library" : cfg.Library,
				"error" : err,
			}).Fatal("could not create library directory")
		}
	}

	if cfg.Update { //} && cfg.Rollback { We actually need this to run either way, as the "prepare packages for update" operation moves the current installations to __OLD__ folders, thus allowing updated versions to be installed.
		log.Info("update argument passed. staging packages for update...")
		rollbackPlan.PreparePackagesForUpdate(fs, cfg.Library)
	}

	// Create a list of package download objects using our install plan and our "nexus" object.
	//pkgsToDownload := getPackagesToDownload(installPlan, pkgNexus)

	// Retrieve a cache to store any packages we need to download for the install.
	packageCache := rcmd.NewPackageCache(userCache(cfg.Cache), false)

	//Create a pkgMap object, which helps us with parallel downloads (?)
	pkgMap, err := cran.DownloadPackages(fs, installPlan.PackageDownloads, packageCache.BaseDir, rVersion)
	if err != nil {
		log.Fatalf("error downloading packages: %s", err)
	}

	//Set the arguments to be passed in to the R Package Installer
	pkgInstallArgs := rcmd.NewDefaultInstallArgs()
	pkgInstallArgs.Library, _ = filepath.Abs(cfg.Library)

	// Get the number of workers.
	// Use number of user defined threads if set. Otherwise, use the number of CPUs (up to 8).
	nworkers := getWorkerCount(cfg.Threads, runtime.NumCPU())

	// Process any customizations set in the yaml config file for individual packages.
	// Set ENV values in rSettings
	rSettings = configlib.SetCustomizations(rSettings, cfg)

	//
	// Install the packages
	//
	err = rcmd.InstallPackagePlan(fs, installPlan, pkgMap, packageCache, pkgInstallArgs, rSettings, rcmd.ExecSettings{PkgrVersion: VERSION}, nworkers)

	//
	// Install the tarballs, if applicable.
	//
	installTarballs(installPlan, rSettings)

	if cfg.Rollback {
		//If anything went wrong during the installation, rollback the environment.
		if err != nil {
			errRollback := rollback.RollbackPackageEnvironment(fs, rollbackPlan)
			if errRollback != nil {
				log.WithFields(log.Fields{

				}).Error("failed to reset package environment after bad installation. Your package Library will be in a corrupt state. It is recommended you delete your Library and reinstall all packages.")
			}
		}
	}
	// If any packages were being updated, we need to remove any leftover backup folders that were created.
	rollback.DeleteBackupPackageFolders(fs, rollbackPlan.UpdateRollbacks)

	log.Info("duration:", time.Since(startTime))

	if err != nil {
		log.Errorf("failed package install with err, %s", err)
	}

	return nil
}

func installTarballs(installPlan gpsr.InstallPlan, rSettings rcmd.RSettings) {
	// Install the tarballs
	iargs := rcmd.NewDefaultInstallArgs()

	libraryAbs, err := filepath.Abs(cfg.Library) // Need to use absolute path to library to accomodate our "extracurricular usage" of Install method.
	if err != nil {
		log.WithFields(log.Fields{
			"lib_folder": cfg.Library,
			"error":   err,
		}).Error("error installing tarball -- could not find absolute path for library folder")
	}

	iargs.Library = libraryAbs
	for tarballPkg, tarballPath := range installPlan.Tarballs {
		log.WithFields(log.Fields{
			"package":     tarballPkg,
			"tarball": tarballPath,
		}).Info("installing tarball")

		//installReqest2 := rcmd.InstallRequest{
		//	Package: tarballPkg,
		//	Cache: rcmd.PackageCache {BaseDir: userCache(cfg.Cache)},
		//	Metadata : cran.Download {Path : tarballPath},
		//	InstallArgs: iargs,
		//	RSettings: rSettings,
		//}
		//
		//res, _, err := rcmd.InstallThroughBinary(fs, installReqest2, rcmd.PackageCache{ BaseDir: filepath.Dir(tarballPath)}) // We pass in the tarball path for the pkg cache here to assure that pkgr doens't find a prexisting binary and thus alwasy rebuilds.
		//if err != nil || res.ExitCode != 0 {
		//	log.WithFields(log.Fields{
		//		"err":       err,
		//		"exit_code": res.ExitCode,
		//		"tarball":   tarballPath,
		//	}).Error("could not install package, problem in function InstallThroughBinary")
		//}
		tarballPathAbs, err := filepath.Abs(tarballPath) // Need to do this or else we get a weird bug from filepath.Clean in the Install function.
		if err != nil {
			log.WithFields(log.Fields{
				"pkg":     tarballPkg,
				"tarball": tarballPath,
				"error":   err,
			}).Error("error installing tarball -- could not find absolute path for tarball")
		}

		res, err := rcmd.Install(
			fs,
			tarballPkg,
			tarballPathAbs,
			iargs,
			rSettings,
			rcmd.ExecSettings{
				PkgrVersion: VERSION,
				WorkDir: filepath.Dir(tarballPath),
			},
			rcmd.InstallRequest{
				Package: tarballPkg,
				Cache: rcmd.PackageCache{
					BaseDir: userCache(cfg.Cache),
				},
				InstallArgs: iargs,
			},
		)

		if err != nil {
			log.WithFields(log.Fields{
				"pkg":     tarballPkg,
				"tarball": tarballPath,
				"error":   err,
			}).Error("error installing tarball")
		} else {
			log.WithFields(log.Fields{
				"pkg":     tarballPkg,
				"tarball": tarballPath,
				"output":  res.Stdout,
			}).Info("tarball installed successfully")
		}
	}
}


func initInstallLog() {
	//Init install-specific log, if one has been set. This overwrites the default log.
	if cfg.Logging.Install != "" {
		logger.AddLogFile(cfg.Logging.Install, cfg.Logging.Overwrite)
	} else {
		logger.AddLogFile(cfg.Logging.All, cfg.Logging.Overwrite)
	}
}
