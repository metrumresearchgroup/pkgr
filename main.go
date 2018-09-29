package main

import (
	"fmt"

	"github.com/dpastoor/rpackagemanager/packrat"
)

func main() {
	str := []byte(`Package: PKPDmisc
Source: CRAN
Version: 2.1.1
Hash: 5cb470e20683af7fe4ee68681194d3dd
Requires: BH, Rcpp, data.table, dplyr, ggplot2, lazyeval, magrittr,
	purrr, readr, rlang, stringr, tibble
`)
	strc := packrat.CollapseIndentation(str)
	pv := packrat.ParsePackageReqs(strc)
	fmt.Println(pv)
}
