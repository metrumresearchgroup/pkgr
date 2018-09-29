package main

import (
	"fmt"

	"github.com/dpastoor/rpackagemanager/installer"
	"github.com/dpastoor/rpackagemanager/packrat"
	"github.com/spf13/afero"
)

func main() {

	appFS := afero.NewOsFs()
	lf, _ := afero.ReadFile(appFS, "installer/testdata/01-mixed_gh_cran_packrat.lock")
	pm := packrat.ChunkLockfile(lf)

	var workingGraph installer.Graph
	for _, p := range pm.CRANlike {
		workingGraph = append(workingGraph, installer.NewNode(p.Package, p.Requires))
	}
	for _, p := range pm.Github {
		workingGraph = append(workingGraph, installer.NewNode(p.Reqs.Package, p.Reqs.Requires))
	}

	installer.DisplayGraph(workingGraph)

	resolved, err := installer.ResolveGraph(workingGraph)
	if err != nil {
		fmt.Printf("Failed to resolve dependency graph: %s\n", err)
	} else {
		fmt.Println("The dependency graph resolved successfully")
	}

	for _, pkglayer := range resolved {
		fmt.Println(pkglayer)
	}
}
