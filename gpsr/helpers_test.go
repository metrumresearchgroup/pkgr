package gpsr

import (
	"fmt"
	"testing"

	"github.com/metrumresearchgroup/pkgr/cran"
	"github.com/metrumresearchgroup/pkgr/desc"
	"github.com/stretchr/testify/assert"
)

func TestAppendGraph(t *testing.T) {
	workingGraph := NewGraph()
	dependencyConfigurations := NewDefaultInstallDeps()
	var pkgDesc = desc.Desc{
		Package: "roxygen2", Source: "", Version: "4.1.1.9000", Maintainer: "",
		Imports: map[string]desc.Dep{
			"brew": desc.Dep{Name: "brew", Version: desc.Version{Major: 0, Minor: 0, Patch: 0, Dev: 0, Other: 0}, Constraint: 0},
		},
	}
	var urls = []cran.RepoURL{
		cran.RepoURL{
			Name: "CRAN",
			URL:  "https://cran.microsoft.com/snapshot/2018-11-11",
		},
	}
	packages := map[string]cran.PkgConfig{
		"brew": cran.PkgConfig{
			Repo: urls[0],
			Type: cran.Source,
		},
	}
	var installConfig = cran.InstallConfig{
		Packages: packages,
	}
	pkgNexus, _ := cran.NewPkgDb(urls, cran.Source, &installConfig, cran.RVersion{})

	// call the function to test
	appendToGraph(workingGraph, pkgDesc, dependencyConfigurations, pkgNexus)

	// brew is a dep of roxygen2
	m := workingGraph["roxygen2"]
	assert.NotEqual(t, nil, m, fmt.Sprintf("Graph Error"))

	n := len(m.Deps)
	assert.GreaterOrEqual(t, 1, n, fmt.Sprintf("Length Error"))

	md := m.Deps[0]
	assert.Equal(t, "brew", md, fmt.Sprintf("Deps Error"))
}
