package cran

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetType(t *testing.T) {
	var pkgConfig PkgConfig

	setType(&pkgConfig, "source")
	assert.Equal(t, Source, pkgConfig.Type, "Error setting type source")

	setType(&pkgConfig, "binary")
	assert.Equal(t, Binary, pkgConfig.Type, "Error setting type binary")

	err := setType(&pkgConfig, "invalid")
	assert.Equal(t, "invalid source type: invalid", err.Error(), "Error setting type invalid value")
}

func TestSetType2(t *testing.T) {
	var pkgName = "mrgsolve"
	var urls = []RepoURL{
		RepoURL{
			Name: "MPN_2023_03_13",
			URL:  "https://mpn.metworx.com/snapshots/stable/2023-03-13",
		},
	}
	var installConfig = InstallConfig{
		Packages: map[string]PkgConfig{},
	}

	pkgNexus, _ := NewPkgDb(urls, Source, &installConfig, RVersion{4, 1, 3}, false)
	_, pkgCfg, _ := pkgNexus.GetPackage(pkgName)
	assert.Equal(t, Source, pkgCfg.Type, "Error getting type source")

	pkgNexus.SetPackageType(pkgName, "binary")
	assert.Equal(t, Binary, pkgNexus.Config.Packages[pkgName].Type, "Error setting type binary")

	pkgNexus.SetPackageType(pkgName, "source")
	assert.Equal(t, Source, pkgNexus.Config.Packages[pkgName].Type, "Error setting type source")

	err := pkgNexus.SetPackageType(pkgName, "invalid")
	assert.Equal(t, "invalid source type: invalid", err.Error(), "Error setting type invalid value")
}

func TestSetRepo(t *testing.T) {
	var pkgName = "mrgsolve"
	var urls = []RepoURL{
		RepoURL{
			Name: "MPN_2023_03_13",
			URL:  "https://mpn.metworx.com/snapshots/stable/2023-03-13",
		},
		RepoURL{
			Name: "MPN_2022_06_15",
			URL:  "https://mpn.metworx.com/snapshots/stable/2022-06-15",
		},
	}
	var installConfig = InstallConfig{
		Packages: map[string]PkgConfig{},
	}

	pkgNexus, _ := NewPkgDb(urls, Source, &installConfig, RVersion{4, 1, 3}, false)

	_, pkgCfg, _ := pkgNexus.GetPackage(pkgName)
	assert.Equal(t, "MPN_2023_03_13", pkgCfg.Repo.Name, "Error getting repo MPN_2023_03_13")

	pkgNexus.SetPackageRepo(pkgName, "MPN_2022_06_15")
	_, pkgCfg, _ = pkgNexus.GetPackage(pkgName)
	assert.Equal(t, "MPN_2022_06_15", pkgCfg.Repo.Name, "Error getting repo MPN_2022_06_15")
}
