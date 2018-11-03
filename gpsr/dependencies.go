package gpsr

import (
	"errors"
	"fmt"

	"github.com/deckarep/golang-set"
	"github.com/dpastoor/rpackagemanager/cran"
)

type Graph map[string]*Node

// NewGraph initializes a new graph
func NewGraph() Graph {
	return make(map[string]*Node)
}

// Node represents a single node in the graph with it's dependencies
type Node struct {
	// Name of the node
	Name string

	// Dependencies of the node
	Deps []string
}

// NewNode creates a new node
func NewNode(name string, deps []string) *Node {
	n := &Node{
		Name: name,
		Deps: deps,
	}

	return n
}

// DisplayGraph shows the dependency graph
func DisplayGraph(graph Graph) {
	for _, node := range graph {
		for _, dep := range node.Deps {
			fmt.Printf("%s -> %s\n", node.Name, dep)
		}
	}
}

// ResolveLayers Resolves the dependency graph
// providing the layers of dependencies such that
// each layer only dependends on elements from above
// allowing confident parallel installation at any
// layer
func ResolveLayers(graph Graph) ([][]string, error) {
	// A map containing the nodes and their dependencies
	nodeDependencies := make(map[string]mapset.Set)

	// Populate the maps
	for nm, node := range graph {

		dependencySet := mapset.NewSet()
		for _, dep := range node.Deps {
			if !isDefaultPackage(dep) {
				dependencySet.Add(dep)
			}
		}
		nodeDependencies[nm] = dependencySet
	}

	// Iteratively find and remove nodes from the graph which have no dependencies.
	// If at some point there are still nodes in the graph and we cannot find
	// nodes without dependencies, that means we have a circular dependency
	var resolved [][]string
	for len(nodeDependencies) != 0 {
		// Get all nodes from the graph which have no dependencies
		readySet := mapset.NewSet()
		for name, deps := range nodeDependencies {
			if deps.Cardinality() == 0 {
				readySet.Add(name)
			}
		}

		// If there aren't any ready nodes, then we have a cicular dependency
		if readySet.Cardinality() == 0 {
			return resolved, errors.New("Circular dependency found")
		}

		// Remove the ready nodes and add them to the resolved graph
		var dl []string
		for name := range readySet.Iter() {
			delete(nodeDependencies, name.(string))
			dl = append(dl, name.(string))
		}
		resolved = append(resolved, dl)

		// Also make sure to remove the ready nodes from the
		// remaining node dependencies as well
		for name, deps := range nodeDependencies {
			diff := deps.Difference(readySet)
			nodeDependencies[name] = diff
		}
	}

	return resolved, nil
}

// ResolveInstallationReqs resolves all the installation requirements
func ResolveInstallationReqs(pkgs []string, pkgdb *cran.PkgDb) (InstallPlan, error) {
	workingGraph := NewGraph()
	depDb := make(map[string][]string)
	for _, p := range pkgs {
		pkg, _, _ := pkgdb.GetPackage(p)
		appendToGraph(workingGraph, pkg, pkgdb)
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
			appendToGraph(workingGraph, pkg, pkgdb)
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

// InvertDependencies provides an inversion of the dependencies
// such that each element contains a slice of all packages that depend on it
// This can be used when a package is installed to identify which
// other packages may have all their dependencies satisfied.
func (ip InstallPlan) InvertDependencies() map[string][]string {
	ddb := ip.DepDb
	idb := make(map[string][]string)
	for pkg, deps := range ddb {
		for _, p := range deps {
			d := idb[p]
			d = append(d, pkg)
			idb[p] = d
		}
	}
	return idb
}
