package gpsr

import (
	"github.com/dpastoor/rpackagemanager/cran"
	"github.com/dpastoor/rpackagemanager/desc"
)

func isDefaultPackage(pkg string) bool {
	_, exists := DefaultPackages[pkg]
	return exists
}

func appendToGraph(m Graph, d desc.Desc, ids map[string]InstallDeps, pkgdb *cran.PkgDb) {
	var reqs []string
	id, exists := ids[d.Package]
	if !exists {
		id = NewDefaultInstallDeps()
	}
	if id.Depends {
		for r := range d.Depends {
			_, _, ok := pkgdb.GetPackage(r)
			if ok && !isDefaultPackage(r) {
				reqs = append(reqs, r)
			}
		}
	}
	if id.Imports {
		for r := range d.Imports {
			_, _, ok := pkgdb.GetPackage(r)
			if ok && !isDefaultPackage(r) {
				reqs = append(reqs, r)
			}
		}
	}
	if id.LinkingTo {
		for r := range d.LinkingTo {
			_, _, ok := pkgdb.GetPackage(r)
			if ok && !isDefaultPackage(r) {
				reqs = append(reqs, r)
			}
		}
	}
	if id.Suggests {
		for r := range d.Suggests {
			_, _, ok := pkgdb.GetPackage(r)
			if ok && !isDefaultPackage(r) {
				reqs = append(reqs, r)
			}
		}
	}
	m[d.Package] = NewNode(d.Package, reqs)
	if len(reqs) > 0 {
		for _, pn := range reqs {
			_, ok := m[pn]
			if pn != "R" && !ok {
				pkg, _, _ := pkgdb.GetPackage(pn)
				appendToGraph(m, pkg, ids, pkgdb)
			}
		}
	}
}
