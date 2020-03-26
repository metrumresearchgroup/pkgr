package rollback

import (
	"github.com/metrumresearchgroup/pkgr/cran"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"path/filepath"
)

// RollbackPackageEnvironment "executes" the given RollbackPlan, thereby resetting the environment to the state it was in before rbp.InstallPlan was applied.
func RollbackPackageEnvironment(fileSystem afero.Fs, rbp RollbackPlan) error {

	//reset packages
	logrus.Trace("Resetting package environment")

	if rbp.InstallPlan.CreateLibrary {
		err0 := fileSystem.RemoveAll(rbp.Library)

		if err0 != nil {
			logrus.WithFields(logrus.Fields{
				"library": rbp.Library,
				"error": err0,
			}).Warn("failed to remove created library during rollback")
			return err0
		}

		return nil
	}

	for _, pkg := range rbp.NewPackages {
		err1 := fileSystem.RemoveAll(filepath.Join(rbp.Library, pkg))
		if err1 != nil {
			logrus.WithFields(logrus.Fields{
				"library": rbp.Library,
				"package": pkg,
			}).Warn("failed to remove package during rollback", err1)
			return err1
		}
	}


	//Rollback updated packages -- we have to do this differently than the rest, because updated packages need to be
	//restored from backups.
	if len(rbp.UpdateRollbacks) > 0 {
		err2 := rollbackChangedPackages(fileSystem, rbp.UpdateRollbacks)
		if err2 != nil {
			return err2
		}
	}

	//Rollback additional packages -- same process as updated packages, but we hold them in a separate object.
	if len(rbp.AdditionalPkgRollbacks) > 0 {
		err3 := rollbackChangedPackages(fileSystem, rbp.AdditionalPkgRollbacks)
		if err3 != nil {
			return err3
		}
	}

	return nil
}

func DeleteBackupPackageFolders(fileSystem afero.Fs, packageBackupInfo []UpdateAttempt) error {
	if len(packageBackupInfo) == 0 {
		logrus.Debug("No update-packages to restore.")
		return nil
	}

	for _, info := range packageBackupInfo {

		backupExists, _ := afero.Exists(fileSystem, info.BackupPackageDirectory)
		//_, err1 := fileSystem.Stat(info.BackupPackageDirectory) // Checking existence
		if backupExists {
			err1 := fileSystem.RemoveAll(info.BackupPackageDirectory)
			if err1 != nil {
				logrus.WithFields(logrus.Fields{
					"package":           info.Package,
					"problem_directory": info.BackupPackageDirectory,
				}).Warn("could not delete directory during cleanup")
				return err1
			}
		}
	}

	return nil
}

func rollbackChangedPackages(fileSystem afero.Fs, packageBackupInfo []UpdateAttempt) error {

	if len(packageBackupInfo) == 0 {
		logrus.Debug("Not update-packages to restore.")
		return nil
	}

	for _, info := range packageBackupInfo {

		logrus.WithFields(logrus.Fields{
			"pkg":                 info.Package,
			"rolling back to":     info.OldVersion,
			"failed to update to": info.NewVersion,
		}).Warn("did not update package, restoring last-installed version")

		_, err1 := fileSystem.Stat(info.ActivePackageDirectory) // Checking existence
		if err1 == nil {
			err1 = fileSystem.RemoveAll(info.ActivePackageDirectory)
			if err1 != nil {
				logrus.WithFields(logrus.Fields{
					"package":           info.Package,
					"problem_directory": info.ActivePackageDirectory,
				}).Warn("could not delete directory during package rollback")
				return err1
			}
		}

		err2 := fileSystem.Rename(info.BackupPackageDirectory, info.ActivePackageDirectory)

		if err2 != nil {
			logrus.WithFields(logrus.Fields{
				"pkg": info.Package,
			}).Warn("could not rollback package -- package will need reinstallation.")
			return err2
		}

	}

	return nil
}

// createUpdateBackupFolders ...
func createUpdateBackupFolders(fileSystem afero.Fs, libraryPath string, outdatedPackages []cran.OutdatedPackage) []UpdateAttempt {
	var updateAttempts []UpdateAttempt

	//Tag each outdated pkg directory in library with "__OLD__"
	for _, pkg := range outdatedPackages {
		updateAttempts = append(updateAttempts, tagOldInstallation(fileSystem, libraryPath, pkg))
	}

	return updateAttempts
}

func tagOldInstallation(fileSystem afero.Fs, libraryPath string, outdatedPackage cran.OutdatedPackage) UpdateAttempt {
	packageDir := filepath.Join(libraryPath, outdatedPackage.Package)
	taggedPackageDir := filepath.Join(libraryPath, "__OLD__"+outdatedPackage.Package)

	err := fileSystem.Rename(packageDir, taggedPackageDir)
	//err := RenameDirRecursive(fileSystem, packageDir, taggedPackageDir)

	if err != nil {
		logrus.WithField("package dir", packageDir).Warn("error when backing up outdated package")
		logrus.Error(err)
	}

	return UpdateAttempt{
		Package:                outdatedPackage.Package,
		ActivePackageDirectory: packageDir,
		BackupPackageDirectory: taggedPackageDir,
		OldVersion:             outdatedPackage.OldVersion,
		NewVersion:             outdatedPackage.NewVersion,
	}
}



