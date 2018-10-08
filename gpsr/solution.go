package gpsr

import (
	"github.com/dpastoor/rpackagemanager/packrat"
)

// SolveLockfile provides a solution give a lockfile spec
func SolveLockfile(lf []byte) ([][]string, error) {
	pm := packrat.ChunkLockfile(lf)
	var workingGraph Graph
	for _, p := range pm.CRANlike {
		workingGraph = append(workingGraph, NewNode(p.Package, p.Requires))
	}
	for _, p := range pm.Github {
		workingGraph = append(workingGraph, NewNode(p.Reqs.Package, p.Reqs.Requires))
	}

	resolved, err := ResolveGraph(workingGraph)

	return resolved, err
}
