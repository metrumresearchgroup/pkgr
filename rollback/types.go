package rollback

import (
	"github.com/metrumresearchgroup/pkgr/cran"
	"github.com/metrumresearchgroup/pkgr/desc"
	"github.com/metrumresearchgroup/pkgr/gpsr"
	"github.com/metrumresearchgroup/pkgr/pacman"
	"github.com/spf13/afero"
	"github.com/thoas/go-funk"
)

type RollbackPlan struct {
	AllPackages []string
	NewPackages []string
	UpdateRollbacks []pacman.UpdateAttempt
	InstallPlan gpsr.InstallPlan
	Library string
}

//Constructor
func CreateRollbackPlan (library string, installPlan gpsr.InstallPlan, preinstalledPackages map[string]desc.Desc) RollbackPlan {
	ap := installPlan.GetAllPackages()
	np := DiscernNewPackages(installPlan.GetAllPackages(), preinstalledPackages)
	return RollbackPlan {
			AllPackages: ap,
			NewPackages: np,
			InstallPlan: installPlan,
			Library: library,
		}
}


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
	updateAttempts := pacman.PreparePackagesForUpdate(fs, library, opFiltered)
	rp.UpdateRollbacks = updateAttempts
}

func DiscernNewPackages(toInstallPackageNames []string, preinstalledPackages map[string]desc.Desc) []string {
	newPackages := make([]string, 0)

	for _, name := range toInstallPackageNames {
		_, found := preinstalledPackages[name]
		if !found {
			newPackages = append(newPackages, name)
		}
	}

	return newPackages
}



