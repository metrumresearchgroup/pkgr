package gpsr

import (
	"fmt"

	"github.com/dpastoor/rpackagemanager/cran"
)

// ResolveInstallationReqs resolves all the installation requirements
func ResolveInstallationReqs(pkgs []string, ids InstallDeps, pkgdb *cran.PkgDb) (InstallPlan, error) {
	workingGraph := NewGraph()
	depDb := make(map[string][]string)
	for _, p := range pkgs {
		pkg, _, _ := pkgdb.GetPackage(p)
		appendToGraph(workingGraph, pkg, ids, pkgdb)
	}
	resolved, err := ResolveLayers(workingGraph)
	if err != nil {
		fmt.Println("error resolving graph")
		return InstallPlan{}, err
	}
	for i, l := range resolved {
		if i == 0 {
			// don't need to know dep tree for first layer as shouldn't have any deps
			continue
		}
		for _, p := range l {
			workingGraph := NewGraph()
			pkg, _, _ := pkgdb.GetPackage(p)
			appendToGraph(workingGraph, pkg, ids, pkgdb)
			resolved, err := ResolveLayers(workingGraph)
			if err != nil {
				fmt.Println("error resolving graph")
				return InstallPlan{}, err
			}
			allDeps := resolved[0]
			for j, rl := range resolved {
				if j == 0 {
					continue
				}
				if j+1 == len(resolved) {
					// for last layer don't add the package itself
					for _, pkg := range rl {
						if pkg != p {
							allDeps = append(allDeps, pkg)
						}
					}
				} else {
					allDeps = append(allDeps, rl...)
				}
			}
			depDb[p] = allDeps
		}
	}
	return InstallPlan{StartingPackages: resolved[0],
		DepDb: depDb}, nil
}
