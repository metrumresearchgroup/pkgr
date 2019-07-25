package pacman

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/dpastoor/goutils"
	"github.com/metrumresearchgroup/pkgr/cran"
	"github.com/metrumresearchgroup/pkgr/desc"
	"github.com/metrumresearchgroup/pkgr/gpsr"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// GetPriorInstalledPackages ...
func GetPriorInstalledPackages(fileSystem afero.Fs, libraryPath string) map[string]desc.Desc {
	installed := make(map[string]desc.Desc)
	installedLibrary, err := fileSystem.Open(libraryPath)

	if err != nil {
		log.WithField("libraryPath", libraryPath).Fatal("package library not found at given library path")
		return installed
	}
	defer installedLibrary.Close()

	fileInfo, _ := installedLibrary.Readdir(0)
	installedPackageFolders := goutils.ListDirNames(fileInfo)

	for _, pkgFolder := range installedPackageFolders {
		descriptionFilePath := filepath.Join(libraryPath, pkgFolder, "DESCRIPTION")
		installedPackage, err := scanInstalledPackage(descriptionFilePath, fileSystem)

		if err != nil {
			log.Error(err)
			err = nil
		} else {
			log.WithFields(log.Fields{
				"package":     installedPackage.Package,
				"version":     installedPackage.Version,
				"source repo": installedPackage.Repository,
			}).Debug("found installed package")

			installed[installedPackage.Package] = installedPackage
		}
	}

	return installed
}

// GetInstallers returns the installers for the installed packages
func GetInstallers(ip map[string]desc.Desc) InstalledFromPkgs {
	var pkgr, packrat, unknown []string
	for k, v := range ip {
		if v.PkgrVersion == "" {
			packrat = append(packrat, k)
		} else {
			pkgr = append(pkgr, k)
		}
	}
	return InstalledFromPkgs{
		Pkgr:    pkgr,
		Packrat: packrat,
		Unknown: unknown,
	}

}

// GetPackagesByInstalledFrom returns InstalledFromPkgs structure
// single location where business rule of "not pkgr" is applied
func GetPackagesByInstalledFrom(fileSystem afero.Fs, libraryPath string) InstalledFromPkgs {
	var pkgr, packrat, unknown []string
	ip := GetPriorInstalledPackages(fileSystem, libraryPath)
	for k, v := range ip {
		if v.PkgrVersion == "" {
			packrat = append(packrat, k)
		} else {
			pkgr = append(pkgr, k)
		}
	}
	return InstalledFromPkgs{
		Pkgr:    pkgr,
		Packrat: packrat,
		Unknown: unknown,
	}
}

func scanInstalledPackage(
	descriptionFilePath string, fileSystem afero.Fs) (desc.Desc, error) {

	descriptionFile, err := fileSystem.Open(descriptionFilePath)

	if err != nil {
		log.WithField("file", descriptionFilePath).Warn("DESCRIPTION missing from installed package.")
		return desc.Desc{}, err
	}
	defer descriptionFile.Close()

	log.WithField("description_file", descriptionFilePath).Trace("scanning DESCRIPTION file")

	installedPackage, err := desc.ParseDesc(descriptionFile)

	if installedPackage.Package == "" {
		err = errors.New(fmt.Sprintf("encountered a description file without package name: %s", descriptionFile))
		log.WithField("description file", descriptionFilePath).Error(err)
		return desc.Desc{}, err
	}

	return installedPackage, nil
}

// GetOutdatedPackages returns a list of outdated packages
func GetOutdatedPackages(installed map[string]desc.Desc, availablePackages []cran.PkgDl) []gpsr.OutdatedPackage {
	var outdatedPackages []gpsr.OutdatedPackage

	for _, pkgDl := range availablePackages {

		pkgName := pkgDl.Package.Package
		availableVersion := pkgDl.Package.Version

		if installedPkg, found := installed[pkgName]; found {

			// If available version is later than currently installed version
			if desc.CompareVersionStrings(availableVersion, installedPkg.Version) > 0 {
				outdatedPackages = append(outdatedPackages, gpsr.OutdatedPackage{
					Package:    pkgName,
					OldVersion: installed[pkgName].Version,
					NewVersion: pkgDl.Package.Version,
				})
			}
		}
	}
	return outdatedPackages
}

// PreparePackagesForUpdate ...
func PreparePackagesForUpdate(fileSystem afero.Fs, libraryPath string, outdatedPackages []gpsr.OutdatedPackage) []UpdateAttempt {
	var updateAttempts []UpdateAttempt

	//Tag each outdated pkg directory in library with "__OLD__"
	for _, pkg := range outdatedPackages {
		updateAttempts = append(updateAttempts, tagOldInstallation(fileSystem, libraryPath, pkg))
	}

	return updateAttempts
}

func tagOldInstallation(fileSystem afero.Fs, libraryPath string, outdatedPackage gpsr.OutdatedPackage) UpdateAttempt {
	packageDir := filepath.Join(libraryPath, outdatedPackage.Package)
	taggedPackageDir := filepath.Join(libraryPath, "__OLD__"+outdatedPackage.Package)

	err := RenameDirRecursive(fileSystem, packageDir, taggedPackageDir)

	if err != nil {
		log.WithField("package dir", packageDir).Warn("error when backing up outdated package")
		log.Error(err)
	}

	return UpdateAttempt{
		Package:                outdatedPackage.Package,
		ActivePackageDirectory: packageDir,
		BackupPackageDirectory: taggedPackageDir,
		OldVersion:             outdatedPackage.OldVersion,
		NewVersion:             outdatedPackage.NewVersion,
	}
}

// RestoreUnupdatedPackages ...
func RestoreUnupdatedPackages(fileSystem afero.Fs, packageBackupInfo []UpdateAttempt) {

	if len(packageBackupInfo) == 0 {
		return
	}

	//libraryDirectoryFsObject, _ := fs.Open(libraryPath)
	//packageFolderObjects, _ := libraryDirectoryFsObject.Readdir(0)

	for _, info := range packageBackupInfo {
		_, err := fileSystem.Stat(info.ActivePackageDirectory) //Checking existence
		if err == nil {

			fileSystem.RemoveAll(info.BackupPackageDirectory)

			log.WithFields(log.Fields{
				"pkg":         info.Package,
				"old_version": info.OldVersion,
				"new_version": info.NewVersion,
			}).Info("successfully updated package")

		} else {
			log.WithFields(log.Fields{
				"pkg":         info.Package,
				"old_version": info.OldVersion,
				"new_version": info.NewVersion,
			}).Warn("could not update package, restoring last-installed version")
			err := RenameDirRecursive(fileSystem, info.BackupPackageDirectory, info.ActivePackageDirectory)
			if err != nil {
				log.WithField("pkg", info.Package).Error(err)
			}
		}
	}
}

// RenameDirRecursive ...
func RenameDirRecursive(fileSystem afero.Fs, oldPath string, newPath string) error {
	err := CopyDir(fileSystem, oldPath, newPath)

	if err != nil {
		return err
	}

	err = fileSystem.RemoveAll(oldPath)
	if err != nil {
		return err
	}

	return nil
}

// CopyDir ...
//TODO: Move into goutils.
func CopyDir(fs afero.Fs, src string, dst string) error {

	err := fs.MkdirAll(dst, 0755)
	if err != nil {
		return err
	}

	openedDir, err := fs.Open(src)
	if err != nil {
		return err
	}

	directoryContents, err := openedDir.Readdir(0)

	if err != nil {
		return err
	}

	for _, item := range directoryContents {
		srcSubPath := filepath.Join(src, item.Name())
		dstSubPath := filepath.Join(dst, item.Name())
		if item.IsDir() {
			fs.Mkdir(dstSubPath, item.Mode())
			err := CopyDir(fs, srcSubPath, dstSubPath)
			if err != nil {
				return err
			}
		} else {
			_, err := goutils.CopyFS(fs, srcSubPath, dstSubPath)
			if err != nil {
				fmt.Print("Received error: ")
				fmt.Println(err)
				return err
			}
		}
	}
	return nil
}

func stringInSlice(s string, slice []string) bool {
	for _, entry := range slice {
		if s == entry {
			return true
		}
	}
	return false
}

// UpdateAttempt ...
type UpdateAttempt struct {
	Package                string
	ActivePackageDirectory string
	BackupPackageDirectory string
	OldVersion             string
	NewVersion             string
}

// InstalledFromPkgs ...
type InstalledFromPkgs struct {
	Pkgr    []string `json:"pkgr"`
	Packrat []string `json:"packrat"`
	Unknown []string `json:"unknown"`
}

// NotPkgr returns a list of packages not installed by Pkgr
func (ip *InstalledFromPkgs) NotFromPkgr() []string {
	var list []string
	for _, p := range ip.Packrat {
		list = append(list, p)
	}
	for _, p := range ip.Unknown {
		list = append(list, p)
	}
	return list
}

// IsPkgr returns a list of packages installed by Pkgr
func (ip *InstalledFromPkgs) FromPkgr() []string {
	return ip.Pkgr
}
