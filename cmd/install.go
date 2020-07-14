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
	"errors"
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
				"library": cfg.Library,
				"error":   err,
			}).Fatal("could not create library directory")
		}
	}

	if cfg.Update { //} && cfg.Rollback { We actually need this to run either way, as the "prepare packages for update" operation moves the current installations to __OLD__ folders, thus allowing updated versions to be installed.
		log.Info("update argument passed. staging packages for update...")
		rollbackPlan.PreparePackagesForUpdate(fs, cfg.Library)
	}
	rollbackPlan.PrepareAdditionalPackagesForOverwrite(fs, cfg.Library)

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
	errInstallAdditional := installAdditionalPackages(installPlan, rSettings, cfg.Library, cfg.Cache)

	log.WithField("duration", time.Since(startTime)).Info("total package install time")

	if cfg.Rollback {
		//If anything went wrong during the installation, rollback the environment.
		if err != nil || errInstallAdditional != nil {
			errRollback := rollback.RollbackPackageEnvironment(fs, rollbackPlan)
			if errRollback != nil {
				log.WithFields(log.Fields{}).Error("failed to reset package environment after bad installation. Your package Library will be in a corrupt state. It is recommended you delete your Library and reinstall all packages.")
			}
		}
	}
	// If any packages were being updated, we need to remove any leftover backup folders that were created.
	// Errors are handled in the lower functions
	_ = rollbackPlan.DeleteBackupPackageFolders(fs)

	log.Info("duration:", time.Since(startTime))

	if err != nil {
		log.Errorf("failed package install with err, %s", err)
	}

	return nil
}

func installAdditionalPackages(installPlan gpsr.InstallPlan, rSettings rcmd.RSettings, library, cache string) error {

	toInstallCount := len(installPlan.AdditionalPackageSources) - 1

	// Set up installArgs object
	iargs := rcmd.NewDefaultInstallArgs()

	// Need to use absolute path to library to accomodate our "extracurricular usage" of Install method.
	libraryAbs, err := filepath.Abs(library)
	if err != nil {
		log.WithFields(log.Fields{
			"lib_folder": library,
			"error":      err,
		}).Error("error installing tarball -- could not find absolute path for library folder")
		return err
	}
	iargs.Library = libraryAbs

	log.Info("starting individual tarball install")

	var errorAggregator []error

	for pkgName, additionalPkg := range installPlan.AdditionalPackageSources {

		log.WithFields(log.Fields{
			"package":   pkgName,
			"pkgSource": additionalPkg.InstallPath,
		}).Debug("installing tarball")

		// Need to use absolute path or else we encounter a weird bug from filepath.Clean in the Install function.
		// 	(Instead of cleaning the local path, it was basically duplicating the path onto itself: A/B became A/B/A/B)
		pkgSourcePathAbs, err := filepath.Abs(additionalPkg.InstallPath)
		if err != nil {
			log.WithFields(log.Fields{
				"pkg":       pkgName,
				"pkgSource": additionalPkg.InstallPath,
				"error":     err,
			}).Error("error installing tarball -- could not find absolute path for tarball")
			errorAggregator = append(errorAggregator, err)
		}

		res, err := rcmd.Install(
			fs,
			pkgName,
			pkgSourcePathAbs,
			iargs,
			rSettings,
			rcmd.ExecSettings{
				PkgrVersion: VERSION,
				WorkDir:     filepath.Dir(additionalPkg.InstallPath),
			},
			rcmd.InstallRequest{
				Package: pkgName,
				Cache: rcmd.PackageCache{
					BaseDir: userCache(cache),
				},
				InstallArgs: iargs,
				ExecSettings: rcmd.ExecSettings{ // Needed for updating description file
					PkgrVersion: VERSION, // Needed for updating description file
				},
				Metadata: cran.Download{ // Needed for updating description file
					Metadata: cran.PkgDl{ // Needed for updating description file
						Config: cran.PkgConfig{ // Needed for updating description file
							Type: cran.Source, // Needed for updating description file
							Repo: cran.RepoURL{ // Needed for updating description file
								URL:  pkgSourcePathAbs,    // Needed for updating description file
								Name: "IndividualPackage", // Needed for updating description file
							},
						},
					},
				},
			},
		)

		if err != nil {
			log.WithFields(log.Fields{
				"pkg":           pkgName,
				"source":        additionalPkg.OriginPath,
				"installedFrom": additionalPkg.InstallPath,
				"installType":   additionalPkg.Type,
				"error":         err,
			}).Error("error installing package")
			log.WithFields(log.Fields{
				"pkg":           pkgName,
				"source":        additionalPkg.OriginPath,
				"installedFrom": additionalPkg.InstallPath,
				"installType":   additionalPkg.Type,
				"remaining":     toInstallCount,
				"stdout":        res.Stdout,
				"stderr":        res.Stderr,
			}).Debug("error installing package")
			errorAggregator = append(errorAggregator, err)
		} else {
			log.WithFields(log.Fields{
				"pkg":         pkgName,
				"source":      additionalPkg.OriginPath,
				"installType": additionalPkg.Type,
				"remaining":   toInstallCount,
			}).Info("Successfully Installed Package.")
			log.WithFields(log.Fields{
				"pkg":           pkgName,
				"source":        additionalPkg.OriginPath,
				"installedFrom": additionalPkg.InstallPath,
				"installType":   additionalPkg.Type,
				"remaining":     toInstallCount,
				"stdout":        res.Stdout,
			}).Trace("Successfully Installed Package.")
		}

		toInstallCount--
	}
	if len(errorAggregator) > 0 {
		return errors.New("errorAggregator occured while installing additional packages. see logs for more detail")
	}
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
