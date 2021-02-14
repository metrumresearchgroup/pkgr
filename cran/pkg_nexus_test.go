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
	var pkgName = "sankey"
	var urls = []RepoURL{
		RepoURL{
			Name: "CRAN_2018_11_11",
			URL:  "https://cran.microsoft.com/snapshot/2018-11-11",
		},
	}
	var installConfig = InstallConfig{
		Packages: map[string]PkgConfig{},
	}

	pkgNexus, _ := NewPkgDb(urls, Source, &installConfig, RVersion{}, false)
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
	var pkgName = "sankey"
	var urls = []RepoURL{
		RepoURL{
			Name: "CRAN_2018_11_11",
			URL:  "https://cran.microsoft.com/snapshot/2018-11-11",
		},
		RepoURL{
			Name: "CRAN_2018_11_12",
			URL:  "https://cran.microsoft.com/snapshot/2018-11-12",
		},
	}
	var installConfig = InstallConfig{
		Packages: map[string]PkgConfig{},
	}

	pkgNexus, _ := NewPkgDb(urls, Source, &installConfig, RVersion{}, false)

	_, pkgCfg, _ := pkgNexus.GetPackage(pkgName)
	assert.Equal(t, "CRAN_2018_11_11", pkgCfg.Repo.Name, "Error getting repo CRAN_2018_11_11")

	pkgNexus.SetPackageRepo(pkgName, "CRAN_2018_11_12")
	_, pkgCfg, _ = pkgNexus.GetPackage(pkgName)
	assert.Equal(t, "CRAN_2018_11_12", pkgCfg.Repo.Name, "Error setting repo CRAN_2018_11_12")
}
