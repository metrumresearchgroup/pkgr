package desc

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDescParsing(t *testing.T) {
	assert := assert.New(t)
	var data = []struct {
		in       string
		expected Desc
	}{
		{
			"testdata/D1",
			Desc{Package: "desc", Source: "",
				Version: "1.0.0", Maintainer: "Gábor Csárdi <csardi.gabor@gmail.com>",
				Description: "Tools to read, write, create, and manipulate DESCRIPTION\n   files. It is intented for packages that create or manipulate other\n   packages.\n", MD5sum: "",
				Imports:   map[string]Dep{"R6": Dep{Name: "R6", Version: Version{Major: 0, Minor: 0, Patch: 0, Dev: 0, Other: 0}, Constraint: 0}},
				Suggests:  map[string]Dep{"testthat": Dep{Name: "testthat", Version: Version{Major: 0, Minor: 0, Patch: 0, Dev: 0, Other: 0}, Constraint: 0}},
				Depends:   map[string]Dep{},
				LinkingTo: map[string]Dep{},
			},
		},
		{
			"testdata/D2",
			Desc{Package: "roxygen2", Source: "", Version: "4.1.1.9000", Maintainer: "",
				Description: "A 'Doxygen'-like in-source documentation system\n   for Rd, collation, and 'NAMESPACE' files.\n",
				MD5sum:      "",
				Remotes: []string{
					"foo/digest", "svn::https://github.com/hadley/stringr",
					"local::/pkgs/testthat"},
				Imports: map[string]Dep{
					"digest":  Dep{Name: "digest", Version: Version{Major: 0, Minor: 0, Patch: 0, Dev: 0, Other: 0}, Constraint: 0},
					"methods": Dep{Name: "methods", Version: Version{Major: 0, Minor: 0, Patch: 0, Dev: 0, Other: 0}, Constraint: 0},
					"Rcpp":    Dep{Name: "Rcpp", Version: Version{Major: 0, Minor: 11, Patch: 0, Dev: 0, Other: 0, String: "0.11.0"}, Constraint: 1},
					"stringr": Dep{Name: "stringr", Version: Version{Major: 0, Minor: 5, Patch: 0, Dev: 0, Other: 0, String: "0.5"}, Constraint: 1},
					"brew":    Dep{Name: "brew", Version: Version{Major: 0, Minor: 0, Patch: 0, Dev: 0, Other: 0}, Constraint: 0}},
				Suggests: map[string]Dep{"knitr": Dep{Name: "knitr", Version: Version{Major: 0, Minor: 0, Patch: 0, Dev: 0, Other: 0}, Constraint: 0},
					"testthat": Dep{Name: "testthat", Version: Version{Major: 0, Minor: 8, Patch: 0, Dev: 0, Other: 0, String: "0.8.0"}, Constraint: 1},
				},
				Depends:   map[string]Dep{"R": Dep{Name: "R", Version: Version{Major: 3, Minor: 0, Patch: 2, Dev: 0, Other: 0, String: "3.0.2"}, Constraint: 1}},
				LinkingTo: map[string]Dep{"Rcpp": Dep{Name: "Rcpp"}}},
		},
		{
			"testdata/D5",
			Desc{
				Package:     "desc",
				Source:      "",
				Version:     "1.0.0",
				Maintainer:  "Gábor Csárdi <csardi.gabor@gmail.com>",
				Description: "Tools to read, write, create, and manipulate DESCRIPTION\n   files. It is intented for packages that create or manipulate other\n   packages.\n",
				MD5sum:      "",
				Imports:     map[string]Dep{"httrmock": Dep{Name: "httrmock", Version: Version{Major: 0, Minor: 0, Patch: 0, Dev: 0, Other: 0}, Constraint: 0}, "lintr": Dep{Name: "lintr", Version: Version{Major: 0, Minor: 0, Patch: 0, Dev: 0, Other: 0}, Constraint: 0}},
				Suggests:    map[string]Dep{},
				Depends:     map[string]Dep{},
				LinkingTo:   map[string]Dep{},
				Remotes: []string{
					"jimhester/lintr",
				}},
		},
	}
	for i, tt := range data {
		actual, err := ReadDesc(tt.in)
		if err != nil {
			assert.FailNowf("reading failed", "err: %s", err)
		}
		assert.Equal(tt.expected, actual, fmt.Sprintf("test num: %v", i+1))
	}
}
