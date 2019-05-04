package gpsr

import (
	"github.com/metrumresearchgroup/pkgr/cran"
	"github.com/metrumresearchgroup/pkgr/desc"
)

func isDefaultPackage(pkg string) bool {
	_, exists := DefaultPackages[pkg]
	return exists
}

func appendToGraph(m Graph, d desc.Desc, dependencyConfigs InstallDeps, pkgNexus *cran.PkgNexus) {
	var reqs []string
	dependencyConfig, exists := dependencyConfigs.Deps[d.Package]
	if !exists {
		dependencyConfig = dependencyConfigs.Default
	}
	if dependencyConfig.Depends {
		for r := range d.Depends {
			_, _, ok := pkgNexus.GetPackage(r)
			if ok && !isDefaultPackage(r) {
				reqs = append(reqs, r)
			}
		}
	}
	if dependencyConfig.Imports {
		for r := range d.Imports {
			_, _, ok := pkgNexus.GetPackage(r)
			if ok && !isDefaultPackage(r) {
				reqs = append(reqs, r)
			}
		}
	}
	if dependencyConfig.LinkingTo {
		for r := range d.LinkingTo {
			_, _, ok := pkgNexus.GetPackage(r)
			if ok && !isDefaultPackage(r) {
				reqs = append(reqs, r)
			}
		}
	}
	m[d.Package] = NewNode(d.Package, reqs)
	if dependencyConfig.Suggests {
		// suggests can't be requirements, as otherwise will end up getting
		// many circular dependencies, hence instead, we just
		// want to add these to the dependencyConfig graph without tying them
		// to the package specifically as requirements
		for r := range d.Suggests {
			_, ok := m[r]
			if r != "R" && !ok {
				//fmt.Println(d.Package, "-->", r)
				pkg, _, exists := pkgNexus.GetPackage(r)
				if exists {
					appendToGraph(m, pkg, dependencyConfigs, pkgNexus)
				}
			}
		}
	}
	if len(reqs) > 0 {
		for _, pn := range reqs {
			_, ok := m[pn]
			if pn != "R" && !ok {
				pkg, _, exists := pkgNexus.GetPackage(pn)
				if exists {
					appendToGraph(m, pkg, dependencyConfigs, pkgNexus)
				}
			}
		}
	}
}
