package rollback

import (
	"github.com/metrumresearchgroup/pkgr/pacman"
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
			}).Warn("failed to remove created library during rollback", err0)
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
		err2 := pacman.RollbackUpdatePackages(fileSystem, rbp.UpdateRollbacks)
		if err2 != nil {
			return err2
		}
	}

	return nil
}