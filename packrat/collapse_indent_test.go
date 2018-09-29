package packrat

import (
	"bytes"
	"testing"
)

func TestIndentationCollapse(t *testing.T) {
	type test struct {
		input    []byte
		expected []byte
	}
	data := []test{
		{
			[]byte(`Package: PKPDmisc
Source: CRAN
Version: 2.1.1
Hash: 5cb470e20683af7fe4ee68681194d3dd
Requires: BH, Rcpp, data.table, dplyr, ggplot2, lazyeval, magrittr,
	purrr, readr, rlang, stringr, tibble
`),
			[]byte(`Package: PKPDmisc
Source: CRAN
Version: 2.1.1
Hash: 5cb470e20683af7fe4ee68681194d3dd
Requires: BH, Rcpp, data.table, dplyr, ggplot2, lazyeval, magrittr,purrr, readr, rlang, stringr, tibble
`),
		},
	}

	for i, d := range data {
		res := CollapseIndentation(d.input)
		if !bytes.Equal(d.expected, res) {
			t.Errorf("Test %d failed. Expected %s got %s", i, d.expected, res)
		}
	}
}
