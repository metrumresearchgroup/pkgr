package installer

import (
	"errors"
	"fmt"

	"github.com/deckarep/golang-set"
)

type Graph []*Node

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

// ResolveGraph Resolves the dependency graph
func ResolveGraph(graph Graph) (Graph, error) {
	// A map containing the node names and the actual node object
	nodeNames := make(map[string]*Node)

	// A map containing the nodes and their dependencies
	nodeDependencies := make(map[string]mapset.Set)

	// Populate the maps
	for _, node := range graph {
		nodeNames[node.Name] = node

		dependencySet := mapset.NewSet()
		for _, dep := range node.Deps {
			dependencySet.Add(dep)
		}
		nodeDependencies[node.Name] = dependencySet
	}

	// Iteratively find and remove nodes from the graph which have no dependencies.
	// If at some point there are still nodes in the graph and we cannot find
	// nodes without dependencies, that means we have a circular dependency
	var resolved Graph
	for len(nodeDependencies) != 0 {
		// Get all nodes from the graph which have no dependencies
		readySet := mapset.NewSet()
		for name, deps := range nodeDependencies {
			fmt.Printf("identifying deps for: %s\n", name)
			if deps.Cardinality() == 0 {
				fmt.Printf("adding to readyset: %s\n", name)
				readySet.Add(name)
			}
		}

		// If there aren't any ready nodes, then we have a cicular dependency
		if readySet.Cardinality() == 0 {
			var g Graph
			for name := range nodeDependencies {
				g = append(g, nodeNames[name])
			}

			return g, errors.New("Circular dependency found")
		}

		// Remove the ready nodes and add them to the resolved graph
		for name := range readySet.Iter() {
			delete(nodeDependencies, name.(string))
			resolved = append(resolved, nodeNames[name.(string)])
		}

		// Also make sure to remove the ready nodes from the
		// remaining node dependencies as well
		for name, deps := range nodeDependencies {
			diff := deps.Difference(readySet)
			nodeDependencies[name] = diff
		}
	}

	return resolved, nil
}
