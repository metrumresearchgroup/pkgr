package gpsr

import (
	"github.com/metrumresearchgroup/pkgr/cran"
	"github.com/metrumresearchgroup/pkgr/desc"
	log "github.com/sirupsen/logrus"
)

func isExcludedPackage(pkg string, noRecommended bool) bool {
	pkgType, exists := DefaultPackages[pkg]
	if exists && !noRecommended && pkgType == "recommended" {
		return false
	}
	return exists
}

func appendToGraph(m Graph, d desc.Desc, dependencyConfigs InstallDeps, pkgNexus *cran.PkgNexus) {
	var reqs []string
	dependencyConfig, exists := dependencyConfigs.Deps[d.Package]
	if !exists {
		dependencyConfig = dependencyConfigs.Default
	}
	log.WithField("pkg", d.Package).WithField("config", dependencyConfig).Trace("dep config")
	if dependencyConfig.Depends {
		for r := range d.Depends {
			_, _, ok := pkgNexus.GetPackage(r)
			if ok && !isExcludedPackage(r, dependencyConfig.NoRecommended) {
				reqs = append(reqs, r)
			} else {
				log.WithField("pkg", d.Package).WithField("dep", r).Trace("skipping Depends dep")
			}
		}
	}
	if dependencyConfig.Imports {
		for r := range d.Imports {
			_, _, ok := pkgNexus.GetPackage(r)
			if ok && !isExcludedPackage(r, dependencyConfig.NoRecommended) {
				reqs = append(reqs, r)
			} else {
				log.WithField("pkg", d.Package).WithField("dep", r).Trace("skipping Imports dep")
			}
		}
	}
	if dependencyConfig.LinkingTo {
		for r := range d.LinkingTo {
			_, _, ok := pkgNexus.GetPackage(r)
			if ok && !isExcludedPackage(r, dependencyConfig.NoRecommended) {
				reqs = append(reqs, r)
			} else {
				log.WithField("pkg", d.Package).WithField("dep", r).Trace("skipping LinkingTo dep")
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


// Left off here -- we can easily gather a list of Dep objects and try to append them,
// but if we add the full-blown Desc objects, then we'll be putting the tarball package in the install plan.
// We're probably just going to do a hacky workaround for this.
func appendDependencyToGraph(m Graph, d desc.Dep, pkgNexus *cran.PkgNexus) {

}
