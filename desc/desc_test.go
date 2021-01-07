package desc

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDescParsing(t *testing.T) {
	assert := assert.New(t)
	// TODO: add testdata element from recommended package that has Priority: recommended
	// TODO: add testdata element that has a Path for PACKAGES file
	// TODO: update new fields such as License
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
				License: "MIT + file LICENSE",
				NeedsCompilation: false,
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
				LinkingTo: map[string]Dep{"Rcpp": Dep{Name: "Rcpp"}},
				License: "GPL (>= 2)",
				NeedsCompilation: false,
			},

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
				},
				License: "MIT + file LICENSE",
				NeedsCompilation: false,
			},
		},
		{
			"testdata/dplyr.macbinary",
			Desc{Package: "dplyr",
				Source:           "",
				Version:          "0.8.0",
				Maintainer:       "",
				Description:      "",
				MD5sum:           "",
				Remotes:          []string(nil),
				Repository:       "",
				Imports:          map[string]Dep{"R6": Dep{Name: "R6", Version: Version{Major: 2, Minor: 2, Patch: 2, Dev: 0, Other: 0, String: "2.2.2"}, Constraint: 1}, "rlang": Dep{Name: "rlang", Version: Version{Major: 0, Minor: 3, Patch: 0, Dev: 0, Other: 0, String: "0.3.0"}, Constraint: 1}, "tibble": Dep{Name: "tibble", Version: Version{Major: 2, Minor: 0, Patch: 0, Dev: 0, Other: 0, String: "2.0.0"}, Constraint: 1}, "methods": Dep{Name: "methods", Version: Version{Major: 0, Minor: 0, Patch: 0, Dev: 0, Other: 0, String: ""}, Constraint: 0}, "glue": Dep{Name: "glue", Version: Version{Major: 1, Minor: 1, Patch: 1, Dev: 0, Other: 0, String: "1.1.1"}, Constraint: 1}, "magrittr": Dep{Name: "magrittr", Version: Version{Major: 1, Minor: 5, Patch: 0, Dev: 0, Other: 0, String: "1.5"}, Constraint: 1}, "pkgconfig": Dep{Name: "pkgconfig", Version: Version{Major: 2, Minor: 0, Patch: 1, Dev: 0, Other: 0, String: "2.0.1"}, Constraint: 1}, "Rcpp": Dep{Name: "Rcpp", Version: Version{Major: 1, Minor: 0, Patch: 0, Dev: 0, Other: 0, String: "1.0.0"}, Constraint: 1}, "tidyselect": Dep{Name: "tidyselect", Version: Version{Major: 0, Minor: 2, Patch: 5, Dev: 0, Other: 0, String: "0.2.5"}, Constraint: 1}, "utils": Dep{Name: "utils", Version: Version{Major: 0, Minor: 0, Patch: 0, Dev: 0, Other: 0, String: ""}, Constraint: 0}, "assertthat": Dep{Name: "assertthat", Version: Version{Major: 0, Minor: 2, Patch: 0, Dev: 0, Other: 0, String: "0.2.0"}, Constraint: 1}},
				Suggests:         map[string]Dep{"callr": Dep{Name: "callr", Version: Version{Major: 3, Minor: 1, Patch: 1, Dev: 0, Other: 0, String: "3.1.1"}, Constraint: 1}, "lubridate": Dep{Name: "lubridate", Version: Version{Major: 1, Minor: 7, Patch: 4, Dev: 0, Other: 0, String: "1.7.4"}, Constraint: 1}, "mgcv": Dep{Name: "mgcv", Version: Version{Major: 1, Minor: 8, Patch: 23, Dev: 0, Other: 0, String: "1.8.23"}, Constraint: 1}, "rmarkdown": Dep{Name: "rmarkdown", Version: Version{Major: 1, Minor: 8, Patch: 0, Dev: 0, Other: 0, String: "1.8"}, Constraint: 1}, "RPostgreSQL": Dep{Name: "RPostgreSQL", Version: Version{Major: 0, Minor: 6, Patch: 2, Dev: 0, Other: 0, String: "0.6.2"}, Constraint: 1}, "RSQLite": Dep{Name: "RSQLite", Version: Version{Major: 2, Minor: 0, Patch: 0, Dev: 0, Other: 0, String: "2.0"}, Constraint: 1}, "testthat": Dep{Name: "testthat", Version: Version{Major: 2, Minor: 0, Patch: 0, Dev: 0, Other: 0, String: "2.0.0"}, Constraint: 1}, "dtplyr": Dep{Name: "dtplyr", Version: Version{Major: 0, Minor: 0, Patch: 2, Dev: 0, Other: 0, String: "0.0.2"}, Constraint: 1}, "microbenchmark": Dep{Name: "microbenchmark", Version: Version{Major: 1, Minor: 4, Patch: 4, Dev: 0, Other: 0, String: "1.4.4"}, Constraint: 1}, "withr": Dep{Name: "withr", Version: Version{Major: 2, Minor: 1, Patch: 1, Dev: 0, Other: 0, String: "2.1.1"}, Constraint: 1}, "broom": Dep{Name: "broom", Version: Version{Major: 0, Minor: 5, Patch: 1, Dev: 0, Other: 0, String: "0.5.1"}, Constraint: 1}, "bit64": Dep{Name: "bit64", Version: Version{Major: 0, Minor: 9, Patch: 7, Dev: 0, Other: 0, String: "0.9.7"}, Constraint: 1}, "ggplot2": Dep{Name: "ggplot2", Version: Version{Major: 2, Minor: 2, Patch: 1, Dev: 0, Other: 0, String: "2.2.1"}, Constraint: 1}, "hms": Dep{Name: "hms", Version: Version{Major: 0, Minor: 4, Patch: 1, Dev: 0, Other: 0, String: "0.4.1"}, Constraint: 1}, "nycflights13": Dep{Name: "nycflights13", Version: Version{Major: 0, Minor: 2, Patch: 2, Dev: 0, Other: 0, String: "0.2.2"}, Constraint: 1}, "RMySQL": Dep{Name: "RMySQL", Version: Version{Major: 0, Minor: 10, Patch: 13, Dev: 0, Other: 0, String: "0.10.13"}, Constraint: 1}, "purrr": Dep{Name: "purrr", Version: Version{Major: 0, Minor: 3, Patch: 0, Dev: 0, Other: 0, String: "0.3.0"}, Constraint: 1}, "crayon": Dep{Name: "crayon", Version: Version{Major: 1, Minor: 3, Patch: 4, Dev: 0, Other: 0, String: "1.3.4"}, Constraint: 1}, "covr": Dep{Name: "covr", Version: Version{Major: 3, Minor: 0, Patch: 1, Dev: 0, Other: 0, String: "3.0.1"}, Constraint: 1}, "DBI": Dep{Name: "DBI", Version: Version{Major: 0, Minor: 7, Patch: 14, Dev: 0, Other: 0, String: "0.7.14"}, Constraint: 1}, "dbplyr": Dep{Name: "dbplyr", Version: Version{Major: 1, Minor: 2, Patch: 0, Dev: 0, Other: 0, String: "1.2.0"}, Constraint: 1}, "knitr": Dep{Name: "knitr", Version: Version{Major: 1, Minor: 19, Patch: 0, Dev: 0, Other: 0, String: "1.19"}, Constraint: 1}, "Lahman": Dep{Name: "Lahman", Version: Version{Major: 3, Minor: 0, Patch: 1, Dev: 0, Other: 0, String: "3.0-1"}, Constraint: 1}, "MASS": Dep{Name: "MASS", Version: Version{Major: 0, Minor: 0, Patch: 0, Dev: 0, Other: 0, String: ""}, Constraint: 0}, "readr": Dep{Name: "readr", Version: Version{Major: 1, Minor: 3, Patch: 1, Dev: 0, Other: 0, String: "1.3.1"}, Constraint: 1}},
				Depends:          map[string]Dep{"R": Dep{Name: "R", Version: Version{Major: 3, Minor: 1, Patch: 2, Dev: 0, Other: 0, String: "3.1.2"}, Constraint: 1}},
				LinkingTo:        map[string]Dep{},
				License: 		  "",
				NeedsCompilation: false,
			},
		},
		{
			"testdata/dplyr.source",
			Desc{Package: "dplyr",
				Source:      "",
				Version:     "0.8.0.1",
				Maintainer:  "",
				Description: "",
				MD5sum:      "",
				Remotes:     []string(nil),
				Repository:  "",
				Imports:     map[string]Dep{"pkgconfig": Dep{Name: "pkgconfig", Version: Version{Major: 2, Minor: 0, Patch: 1, Dev: 0, Other: 0, String: "2.0.1"}, Constraint: 1}, "Rcpp": Dep{Name: "Rcpp", Version: Version{Major: 1, Minor: 0, Patch: 0, Dev: 0, Other: 0, String: "1.0.0"}, Constraint: 1}, "rlang": Dep{Name: "rlang", Version: Version{Major: 0, Minor: 3, Patch: 0, Dev: 0, Other: 0, String: "0.3.0"}, Constraint: 1}, "tibble": Dep{Name: "tibble", Version: Version{Major: 2, Minor: 0, Patch: 0, Dev: 0, Other: 0, String: "2.0.0"}, Constraint: 1}, "utils": Dep{Name: "utils", Version: Version{Major: 0, Minor: 0, Patch: 0, Dev: 0, Other: 0, String: ""}, Constraint: 0}, "assertthat": Dep{Name: "assertthat", Version: Version{Major: 0, Minor: 2, Patch: 0, Dev: 0, Other: 0, String: "0.2.0"}, Constraint: 1}, "glue": Dep{Name: "glue", Version: Version{Major: 1, Minor: 1, Patch: 1, Dev: 0, Other: 0, String: "1.1.1"}, Constraint: 1}, "magrittr": Dep{Name: "magrittr", Version: Version{Major: 1, Minor: 5, Patch: 0, Dev: 0, Other: 0, String: "1.5"}, Constraint: 1}, "methods": Dep{Name: "methods", Version: Version{Major: 0, Minor: 0, Patch: 0, Dev: 0, Other: 0, String: ""}, Constraint: 0}, "R6": Dep{Name: "R6", Version: Version{Major: 2, Minor: 2, Patch: 2, Dev: 0, Other: 0, String: "2.2.2"}, Constraint: 1}, "tidyselect": Dep{Name: "tidyselect", Version: Version{Major: 0, Minor: 2, Patch: 5, Dev: 0, Other: 0, String: "0.2.5"}, Constraint: 1}},
				Suggests:    map[string]Dep{"covr": Dep{Name: "covr", Version: Version{Major: 3, Minor: 0, Patch: 1, Dev: 0, Other: 0, String: "3.0.1"}, Constraint: 1}, "dbplyr": Dep{Name: "dbplyr", Version: Version{Major: 1, Minor: 2, Patch: 0, Dev: 0, Other: 0, String: "1.2.0"}, Constraint: 1}, "dtplyr": Dep{Name: "dtplyr", Version: Version{Major: 0, Minor: 0, Patch: 2, Dev: 0, Other: 0, String: "0.0.2"}, Constraint: 1}, "knitr": Dep{Name: "knitr", Version: Version{Major: 1, Minor: 19, Patch: 0, Dev: 0, Other: 0, String: "1.19"}, Constraint: 1}, "Lahman": Dep{Name: "Lahman", Version: Version{Major: 3, Minor: 0, Patch: 1, Dev: 0, Other: 0, String: "3.0-1"}, Constraint: 1}, "lubridate": Dep{Name: "lubridate", Version: Version{Major: 1, Minor: 7, Patch: 4, Dev: 0, Other: 0, String: "1.7.4"}, Constraint: 1}, "MASS": Dep{Name: "MASS", Version: Version{Major: 0, Minor: 0, Patch: 0, Dev: 0, Other: 0, String: ""}, Constraint: 0}, "microbenchmark": Dep{Name: "microbenchmark", Version: Version{Major: 1, Minor: 4, Patch: 4, Dev: 0, Other: 0, String: "1.4.4"}, Constraint: 1}, "RMySQL": Dep{Name: "RMySQL", Version: Version{Major: 0, Minor: 10, Patch: 13, Dev: 0, Other: 0, String: "0.10.13"}, Constraint: 1}, "bit64": Dep{Name: "bit64", Version: Version{Major: 0, Minor: 9, Patch: 7, Dev: 0, Other: 0, String: "0.9.7"}, Constraint: 1}, "DBI": Dep{Name: "DBI", Version: Version{Major: 0, Minor: 7, Patch: 14, Dev: 0, Other: 0, String: "0.7.14"}, Constraint: 1}, "ggplot2": Dep{Name: "ggplot2", Version: Version{Major: 2, Minor: 2, Patch: 1, Dev: 0, Other: 0, String: "2.2.1"}, Constraint: 1}, "testthat": Dep{Name: "testthat", Version: Version{Major: 2, Minor: 0, Patch: 0, Dev: 0, Other: 0, String: "2.0.0"}, Constraint: 1}, "crayon": Dep{Name: "crayon", Version: Version{Major: 1, Minor: 3, Patch: 4, Dev: 0, Other: 0, String: "1.3.4"}, Constraint: 1}, "mgcv": Dep{Name: "mgcv", Version: Version{Major: 1, Minor: 8, Patch: 23, Dev: 0, Other: 0, String: "1.8.23"}, Constraint: 1}, "purrr": Dep{Name: "purrr", Version: Version{Major: 0, Minor: 3, Patch: 0, Dev: 0, Other: 0, String: "0.3.0"}, Constraint: 1}, "callr": Dep{Name: "callr", Version: Version{Major: 3, Minor: 1, Patch: 1, Dev: 0, Other: 0, String: "3.1.1"}, Constraint: 1}, "hms": Dep{Name: "hms", Version: Version{Major: 0, Minor: 4, Patch: 1, Dev: 0, Other: 0, String: "0.4.1"}, Constraint: 1}, "nycflights13": Dep{Name: "nycflights13", Version: Version{Major: 0, Minor: 2, Patch: 2, Dev: 0, Other: 0, String: "0.2.2"}, Constraint: 1}, "rmarkdown": Dep{Name: "rmarkdown", Version: Version{Major: 1, Minor: 8, Patch: 0, Dev: 0, Other: 0, String: "1.8"}, Constraint: 1}, "RPostgreSQL": Dep{Name: "RPostgreSQL", Version: Version{Major: 0, Minor: 6, Patch: 2, Dev: 0, Other: 0, String: "0.6.2"}, Constraint: 1}, "RSQLite": Dep{Name: "RSQLite", Version: Version{Major: 2, Minor: 0, Patch: 0, Dev: 0, Other: 0, String: "2.0"}, Constraint: 1}, "withr": Dep{Name: "withr", Version: Version{Major: 2, Minor: 1, Patch: 1, Dev: 0, Other: 0, String: "2.1.1"}, Constraint: 1}, "broom": Dep{Name: "broom", Version: Version{Major: 0, Minor: 5, Patch: 1, Dev: 0, Other: 0, String: "0.5.1"}, Constraint: 1}, "readr": Dep{Name: "readr", Version: Version{Major: 1, Minor: 3, Patch: 1, Dev: 0, Other: 0, String: "1.3.1"}, Constraint: 1}},
				Depends:     map[string]Dep{"R": Dep{Name: "R", Version: Version{Major: 3, Minor: 1, Patch: 2, Dev: 0, Other: 0, String: "3.1.2"}, Constraint: 1}},
				LinkingTo:   map[string]Dep{"BH": Dep{Name: "BH", Version: Version{Major: 1, Minor: 58, Patch: 0, Dev: 1, Other: 0, String: "1.58.0-1"}, Constraint: 1}, "plogr": Dep{Name: "plogr", Version: Version{Major: 0, Minor: 1, Patch: 10, Dev: 0, Other: 0, String: "0.1.10"}, Constraint: 1}, "Rcpp": Dep{Name: "Rcpp", Version: Version{Major: 1, Minor: 0, Patch: 0, Dev: 0, Other: 0, String: "1.0.0"}, Constraint: 1}},
				License: "MIT + file LICENSE",
				NeedsCompilation: true,
			},
		}}
	for i, tt := range data {
		actual, err := ReadDesc(tt.in)
		if err != nil {
			assert.FailNowf("reading failed", "err: %s", err)
		}
		assert.Equal(tt.expected, actual, fmt.Sprintf("test num: %v", i+1))
	}
}
