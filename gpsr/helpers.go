package gpsr

import (
	"github.com/dpastoor/rpackagemanager/cran"
	"github.com/dpastoor/rpackagemanager/desc"
)

func isDefaultPackage(pkg string) bool {
	_, exists := DefaultPackages[pkg]
	return exists
}

func appendToGraph(m Graph, d desc.Desc, ids InstallDeps, pkgdb *cran.PkgDb) {
	var reqs []string
	id, exists := ids.Deps[d.Package]
	if !exists {
		id = ids.Default
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
	m[d.Package] = NewNode(d.Package, reqs)
	if id.Suggests {
		// suggests can't be requirements, as otherwise will end up getting
		// many circular dependencies, hence instead, we just
		// want to add these to the dependency graph without tying them
		// to the package specifically as requirements
		for r := range d.Suggests {
			_, ok := m[r]
			if r != "R" && !ok {
				//fmt.Println(d.Package, "-->", r)
				pkg, _, exists := pkgdb.GetPackage(r)
				if exists {
					appendToGraph(m, pkg, ids, pkgdb)
				}
			}
		}
	}
	if len(reqs) > 0 {
		for _, pn := range reqs {
			_, ok := m[pn]
			if pn != "R" && !ok {
				pkg, _, exists := pkgdb.GetPackage(pn)
				if exists {
					appendToGraph(m, pkg, ids, pkgdb)
				}
			}
		}
	}
}
