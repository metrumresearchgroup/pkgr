package gpsr

import (
	"fmt"
	"github.com/metrumresearchgroup/pkgr/desc"
	"github.com/metrumresearchgroup/pkgr/pacman"

	"github.com/metrumresearchgroup/pkgr/cran"
)

// ResolveInstallationReqs resolves all the installation requirements
func ResolveInstallationReqs(
		pkgs []string,
		preinstalledPkgs map[string]desc.Desc,
		dependencyConfigs InstallDeps,
		pkgNexus *cran.PkgNexus,
		update bool,
		libraryExists bool,
	) (InstallPlan, error) {

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


	outdatedPackages := pacman.GetOutdatedPackages(preinstalledPkgs, pkgNexus.GetPackages(extractNamesFromDesc(preinstalledPkgs)).Packages)

	installPlan := InstallPlan{
		StartingPackages:  resolved[0],
		DepDb:             depDb,
		InstalledPackages: preinstalledPkgs,
		OutdatedPackages:  outdatedPackages,
		CreateLibrary:     !libraryExists,
		Update:            update,
	}
	installPlan.Pack(pkgNexus)
	return installPlan, nil
}

func extractNamesFromDesc(installedPackages map[string]desc.Desc) []string {
	var installedPackageNames []string
	for key := range installedPackages {
		installedPackageNames = append(installedPackageNames, key)
	}
	return installedPackageNames
}