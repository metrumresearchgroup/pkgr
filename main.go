package main

import (
	"fmt"

	"github.com/dpastoor/rpackagemanager/installer"
	"github.com/dpastoor/rpackagemanager/packrat"
)

func main() {
	// 	pkg1 := []byte(`Package: PKPDmisc
	// Source: CRAN
	// Version: 2.1.1
	// Hash: 5cb470e20683af7fe4ee68681194d3dd
	// Requires: BH, Rcpp, data.table, dplyr, ggplot2, lazyeval, magrittr,
	// 	purrr, readr, rlang, stringr, tibble
	// `)
	pkg1 := []byte(`Package: PKPDmisc
Source: CRAN
Version: 2.1.1
Hash: 5cb470e20683af7fe4ee68681194d3dd
Requires: BH, Rcpp
`)
	pkg2 := []byte(`Package: BH
Source: CRAN
Version: 1.66.0-1
Hash: 4cc8883584b955ed01f38f68bc03af6d`)
	pkg3 := []byte(`Package: Rcpp
Source: CRAN
Version: 1.66.0-1
Hash: 4cc8883584b955ed01f38f68bc03af6d`)

	var pkgs [][]byte
	var workingGraph installer.Graph
	pkgs = append(pkgs, pkg1, pkg2, pkg3)
	for _, p := range pkgs {
		pc := packrat.CollapseIndentation(p)
		pv := packrat.ParsePackageReqs(pc)
		fmt.Println(pv)
		workingGraph = append(workingGraph, installer.NewNode(pv.Package, pv.Requires))
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
