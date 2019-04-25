package gpsr

import (
	"errors"
	"fmt"
	"github.com/metrumresearchgroup/pkgr/cran"

	"github.com/deckarep/golang-set"
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

// InvertDependencies provides an inversion of the dependencies
// such that each element contains a slice of all packages that depend on it
// This can be used when a package is installed to identify which
// other packages may have all their dependencies satisfied.
func (ip *InstallPlan) InvertDependencies() map[string][]string {
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

func (ip *InstallPlan) Pack(pkgNexus *cran.PkgNexus) {
	var toDl []cran.PkgDl
	// starting packages
	for _, p := range ip.StartingPackages {
		pkg, cfg, _ := pkgNexus.GetPackage(p)
		toDl = append(toDl, cran.PkgDl{Package: pkg, Config: cfg})
	}
	// all other packages
	for p := range ip.DepDb {
		pkg, cfg, _ := pkgNexus.GetPackage(p)
		toDl = append(toDl, cran.PkgDl{Package: pkg, Config: cfg})
	}
	ip.PackageDownloads = toDl
	//return toDl
}