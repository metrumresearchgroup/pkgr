package gpsr

import (
	"fmt"
	"github.com/metrumresearchgroup/pkgr/desc"

	"github.com/metrumresearchgroup/pkgr/cran"
)

// ResolveInstallationReqs resolves all the installation requirements
func ResolveInstallationReqs(pkgs []string, dependencyConfigs InstallDeps, pkgNexus *cran.PkgNexus, installedPackages map[string]desc.Desc) (InstallPlan, error) {

	var excludedPackages []desc.Desc
	for _, d := range installedPackages {
		excludedPackages = append(excludedPackages, d)
	}

	workingGraph := NewGraph()
	defaultDependencyConfigs := NewDefaultInstallDeps()
	depDb := make(map[string][]string)

	for _, p := range pkgs {
		pkg, _, _ := pkgNexus.GetPackage(p)
		appendToGraph(workingGraph, pkg, dependencyConfigs, pkgNexus)
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
			pkg, _, _ := pkgNexus.GetPackage(p)
			// for dependencies don't want to propogate custom config such as suggests TRUE/FALSE
			// as this should just be about what packages that need to be present
			// to kick off installation - aka Dep/Import/LinkingTo thus we can use
			// the default
			appendToGraph(workingGraph, pkg, defaultDependencyConfigs, pkgNexus)
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

	var startingPackagesPruned  []string
	depDbPruned := make(map[string][]string)

	//TODO: Rewrite these for loops using the installedPackages map instead of iterating
	for _, sp := range resolved[0] {
		isExcluded := false
		for _, excl := range excludedPackages {
			isExcluded = isExcluded || sp == excl.Package
		}

		if !isExcluded {
			startingPackagesPruned = append(startingPackagesPruned, sp)
		}
	}

	//TODO: Rewrite these for loops using the installedPackages map instead of iterating
	// Prune the insides of depDb
	for key, _ := range depDb {
		group := depDb[key]
		var prunedGroup []string
		for _, dep := range group {
			isExcluded := false
			for _, excl := range excludedPackages {
				isExcluded = isExcluded || dep == excl.Package
			}
			if !isExcluded {
				prunedGroup = append(prunedGroup, dep)
			}
		}
		depDb[key] = prunedGroup
	}

	//Prune the top level of depDb
	for key, _ := range depDb {
		if _, found := installedPackages[key]; !found {
			depDbPruned[key] = depDb[key]
		}
	}


	installPlan := InstallPlan{StartingPackages: startingPackagesPruned,//resolved[0],
		DepDb: depDbPruned}//depDb}

	installPlan.Pack(pkgNexus, excludedPackages)
	return installPlan, nil
}
