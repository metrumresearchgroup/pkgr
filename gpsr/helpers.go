package gpsr

import (
	"github.com/dpastoor/rpackagemanager/cran"
	"github.com/dpastoor/rpackagemanager/desc"
)

func isDefaultPackage(pkg string) bool {
	_, exists := DefaultPackages[pkg]
	return exists
}

func appendToGraph(m Graph, d desc.Desc, pkgdb *cran.PkgDb) {
	var reqs []string
	for r := range d.Imports {
		_, _, ok := pkgdb.GetPackage(r)
		if ok && !isDefaultPackage(r) {
			reqs = append(reqs, r)
		}
	}
	for r := range d.Depends {
		_, _, ok := pkgdb.GetPackage(r)
		if ok && !isDefaultPackage(r) {
			reqs = append(reqs, r)
		}
	}
	for r := range d.LinkingTo {
		_, _, ok := pkgdb.GetPackage(r)
		if ok && !isDefaultPackage(r) {
			reqs = append(reqs, r)
		}
	}
	m[d.Package] = NewNode(d.Package, reqs)
	if len(reqs) > 0 {
		for _, pn := range reqs {
			_, ok := m[pn]
			if pn != "R" && !ok {
				pkg, _, _ := pkgdb.GetPackage(pn)
				appendToGraph(m, pkg, pkgdb)
			}
		}
	}
}
