package main

import (
	"fmt"

	"github.com/metrumresearchgroup/pkgr/gpsr"
	"github.com/metrumresearchgroup/pkgr/packrat"
	"github.com/spf13/afero"
)

func main() {

	appFS := afero.NewOsFs()
	lf, _ := afero.ReadFile(appFS, "gpsr/testdata/01-mixed_gh_cran_packrat.lock")
	pm := packrat.ChunkLockfile(lf)

	var workingGraph gpsr.Graph

	for _, p := range pm.CRANlike {
		workingGraph[p.Package] = gpsr.NewNode(p.Package, p.Requires) // append(workingGraph, gpsr.NewNode(p.Package, p.Requires))
	}
	for _, p := range pm.Github {
		workingGraph[p.Reqs.Package] = gpsr.NewNode(p.Reqs.Package, p.Reqs.Requires) //append(workingGraph, gpsr.NewNode(p.Reqs.Package, p.Reqs.Requires))
	}

	gpsr.DisplayGraph(workingGraph)

	resolved, err := gpsr.ResolveLayers(workingGraph)
	if err != nil {
		fmt.Printf("Failed to resolve dependency graph: %s\n", err)
	} else {
		fmt.Println("The dependency graph resolved successfully")
	}

	for _, pkglayer := range resolved {
		fmt.Println(pkglayer)
	}
}
