package main

import (
	"encoding/json"
	"fmt"

	"github.com/metrumresearchgroup/pkgr/configlib"
)

func main() {
	// appFS := afero.NewOsFs()
	// afero.ReadFile(appFS, "../integration_tests/")
	yml := `
Version: 1
Packages:
- test1 # some inline comment
# some other comment
- test2

Library: some/lib/path
Customizations:
  Packages:
    - R6:
        Type: binary

`
	fmtedyml, err := configlib.Format([]byte(yml))
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(fmtedyml))
	}

}
func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}


