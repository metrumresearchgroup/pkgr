package rollback

import (
	"github.com/metrumresearchgroup/pkgr/cran"
	"github.com/metrumresearchgroup/pkgr/desc"
	"github.com/metrumresearchgroup/pkgr/gpsr"
	"github.com/spf13/afero"
	"github.com/thoas/go-funk"
)

// RollbackPlan maintains information about what's being changed on an install, which it can use to "undo" the application of an InstallPlan.
type RollbackPlan struct {
	AllPackages            []string
	NewPackages            []string
	UpdateRollbacks        []UpdateAttempt
	AdditionalPkgRollbacks []UpdateAttempt
	PreinstalledPackages   map[string]desc.Desc
	InstallPlan            gpsr.InstallPlan
	Library                string
}

// CreateRollbackPlan creates a RollbackPlan to track changes to the package environment and undo those changes if necessary.
func CreateRollbackPlan (library string, installPlan gpsr.InstallPlan, preinstalledPackages map[string]desc.Desc) RollbackPlan {
	ap := installPlan.GetAllPackages()

	// We need to determine which packages are both "new" AND part of the installPlan, otherwise we end up
	// changing miscellaneous packages.
	np := discernNewPackages(ap, preinstalledPackages)
	return RollbackPlan {
			AllPackages: ap,
			NewPackages: np,
			PreinstalledPackages: preinstalledPackages,
			InstallPlan: installPlan,
			Library: library,
		}
}

// createUpdateBackupFolders backs up outdated packages in the library by renaming them, thus making space for the updated versions to install.
func (rp *RollbackPlan) PreparePackagesForUpdate(fs afero.Fs, library string) {

	//InstallPlan is aware of _all_ preinstalled packages, even if they're not in pkgr.yml. We don't want to touch
	//preinstalled packages that weren't in pkgr.yml, so we just want to take the active packages from the install plan.
	outdatedPackages := rp.InstallPlan.OutdatedPackages
	var opFiltered []cran.OutdatedPackage
	for _, op := range outdatedPackages {
		if funk.Contains(rp.AllPackages, op.Package) {
			opFiltered = append(opFiltered, op)
		}
	}

	// Wrap this function for now until we're ready to move it over.
	updateAttempts := createUpdateBackupFolders(fs, library, opFiltered)
	rp.UpdateRollbacks = updateAttempts
}

func (rp *RollbackPlan) PrepareAdditionalPackagesForOverwrite(fs afero.Fs, library string) {
	var overwriteAttempts []UpdateAttempt
	for pkg := range rp.InstallPlan.AdditionalPackageSources {
		preinstPkg, isPreinstalled := rp.PreinstalledPackages[pkg]
		if isPreinstalled {
			overwriteAttempt := tagOldInstallation(fs, library, cran.OutdatedPackage{
				Package: pkg,
				OldVersion: preinstPkg.Version,
				NewVersion: "[from tarball]",
			})
			overwriteAttempts = append(overwriteAttempts, overwriteAttempt)
		}
	}
	rp.AdditionalPkgRollbacks = overwriteAttempts
}

// Helper function to determine which packages out of a list are not already installed. Used to determine which packages pkgr specifically will be installing fresh.
func discernNewPackages(toInstallPackageNames []string, preinstalledPackages map[string]desc.Desc) []string {
	var newPackages []string

	for _, name := range toInstallPackageNames {
		_, found := preinstalledPackages[name]
		if !found {
			newPackages = append(newPackages, name)
		}
	}

	return newPackages
}



