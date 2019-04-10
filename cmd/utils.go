package cmd

import (
	"errors"
	"fmt"
	"github.com/metrumresearchgroup/pkgr/cran"
	"github.com/metrumresearchgroup/pkgr/desc"
	"github.com/metrumresearchgroup/pkgr/gpsr"
	"github.com/spf13/afero"
	"os"
	"path/filepath"
	"runtime"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// returns the cache or sets to a cache dir
func userCache(pc string) string {
	// if actually set then use that cache dir
	if pc != "" {
		log.WithField("dir", pc).Trace("package cache directory set by user")
		return pc
	}
	cdir, err := os.UserCacheDir()
	if err != nil {
		log.Warn("could not use user cache dir, using temp dir")
		cdir = os.TempDir()
	}
	log.WithField("dir", cdir).Trace("default package cache directory")

	pkgrCacheDir := filepath.Join(cdir, "pkgr")

	return pkgrCacheDir
}

func getWorkerCount() int {
	var nworkers int
	if viper.GetInt("threads") < 1 {
		nworkers = runtime.NumCPU()
		if nworkers > 2 {
			nworkers = nworkers - 1
		}
	} else {
		nworkers = viper.GetInt("threads")
		if nworkers > runtime.NumCPU()+2 {
			log.Warn("number of workers exceeds the number of threads on machine by at least 2, this may result in degraded performance")
		}
	}
	return nworkers
}

func GetPriorInstalledPackages(fileSystem afero.Fs, libraryPath string) map[string]desc.Desc {
	installed := make(map[string]desc.Desc)
	installedLibrary, err := fileSystem.Open(libraryPath)

	if err != nil {
		log.WithField("libraryPath", libraryPath).Fatal("package library not found at given library path")
		return installed
	}
	defer installedLibrary.Close()

	installedPackageFolders, _ := installedLibrary.Readdir(0)

	for _, pkgFolder := range installedPackageFolders {
		descriptionFilePath := filepath.Join(libraryPath, pkgFolder.Name(), "DESCRIPTION")
		installedPackage, err := scanInstalledPackage(descriptionFilePath, fileSystem)

		if err != nil {
			log.Error(err)
			err = nil
		} else {
			log.WithFields(log.Fields{
				"package":     installedPackage.Package,
				"version":     installedPackage.Version,
				"source repo": installedPackage.Repository,
			}).Info("found installed package")

			installed[installedPackage.Package] = installedPackage
		}
	}

	return installed
}


func scanInstalledPackage(
	descriptionFilePath string,	fileSystem afero.Fs) (desc.Desc, error) {

	descriptionFile, err := fileSystem.Open(descriptionFilePath)

	if err != nil {
		log.WithField("file", descriptionFilePath).Warn("DESCRIPTION missing from installed package.")
		return desc.Desc{}, err
	}
	defer descriptionFile.Close()

	log.WithField("description file", descriptionFilePath).Debug("scanning DESCRIPTION file")

	installedPackage, err := desc.ParseDesc(descriptionFile)

	if installedPackage.Package != "" {
		return installedPackage, nil
	} else {
		err = errors.New(fmt.Sprintf("encountered a description file without package name: %s", descriptionFile))
		log.WithField("description file", descriptionFilePath).Error(err)
		return desc.Desc{}, err
	}
}

func GetOutOfDatePackages(installed map[string]desc.Desc, availablePackages cran.AvailablePkgs) []gpsr.OutdatedPackage {
	var outdatedPackages []gpsr.OutdatedPackage

	for _, pkgDl := range availablePackages.Packages {

		pkgName := pkgDl.Package.Package

		if _, found := installed[pkgName]; found {
			outdatedPackages = append(outdatedPackages, gpsr.OutdatedPackage {
				Package:    pkgName,
				OldVersion: installed[pkgName].Version,
				NewVersion: pkgDl.Package.Version,
			})
		}
	}
	return outdatedPackages
}
