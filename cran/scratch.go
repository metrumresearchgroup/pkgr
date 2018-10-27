package cran

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/dpastoor/rpackagemanager/desc"
)

// ParsePACKAGES parses a PACKAGES file from a cran-like repo
func ParsePACKAGES() {
	packages := []byte(`Package: logrrr
Version: 0.1.0.9000
Imports: R6, crayon, glue (>= 1.3.0), rlang (>=
	0.2.1)
Suggests: testthat, jsonlite, covr, sessioninfo
License: MIT + file LICENSE
MD5sum: 2ac5e74d5161c40fb26dd78b8c19cc8d
NeedsCompilation: no

Package: rrsq
Version: 0.0.2.9000
Imports: R6, purrr, magrittr, httr, jsonlite,
     logrrr (>= 0.0.1)
Suggests: testthat
License: MIT + file LICENSE
MD5sum: ec57bc69465f4ae31c2d59c673f3113d
NeedsCompilation: no
`)
	cb := bytes.Split(packages, []byte("\n\n"))
	for _, p := range cb {
		reader := bytes.NewReader(p)
		d, err := desc.ParseDesc(reader)
		if err != nil {
			panic(err)
		}
		prettyPrint(d)
	}
}

func prettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}
