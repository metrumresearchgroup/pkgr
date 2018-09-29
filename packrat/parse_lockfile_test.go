package packrat

import (
	"reflect"
	"testing"
)

func TestParsePackageReqs(t *testing.T) {
	type test struct {
		input    []byte
		expected PackageReqs
	}
	data := []test{
		{
			[]byte(`Package: fakepkg
Source: CRAN
Version: 0.0.1
Hash: 5cb470e20683af7fe4ee68681194d3dd
Requires: BH, Rcpp, data.table,
	purrr, readr
`),
			PackageReqs{
				Package: "fakepkg",
				Source:  "CRAN",
				Version: "0.0.1",
				Hash:    "5cb470e20683af7fe4ee68681194d3dd",
				Requires: []string{
					"BH",
					"Rcpp",
					"data.table",
					"purrr",
					"readr",
				},
			},
		},
	}

	for i, d := range data {
		res := ParsePackageReqs(CollapseIndentation(d.input))
		if !reflect.DeepEqual(d.expected, res) {
			t.Errorf("Test %d failed. Expected %s got %s", i, d.expected, res)
		}
	}
}
