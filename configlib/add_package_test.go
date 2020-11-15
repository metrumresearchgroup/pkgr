package configlib


import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddPackage(t *testing.T) {
	type TT struct {
		Input []string
		Output []string
		Msg string
	}
	tt := []TT{
		{
			Input: []string{"Packages: ", "- dplyr"},
			Output: []string{"Packages: ", "- newpkg", "- dplyr"},
			Msg: "simple",
	    },
		{
			Input: []string{"Version: 1", "Packages: ", "- dplyr"},
			Output: []string{"Version: 1", "Packages: ", "- newpkg", "- dplyr"},
			Msg: "simple with version",
		},
		{
			Input: []string{"Packages: ", "  - dplyr"},
			Output: []string{"Packages: ", "  - newpkg", "  - dplyr"},
			Msg: "simple indented",
		},
		{
			Input: []string{"Version: 1", "Packages: ", "  - dplyr"},
			Output: []string{"Version: 1", "Packages: ", "  - newpkg", "  - dplyr"},
			Msg: "simple indented with version",
		},
		{
			Input: []string{"Packages: ", "# some comment", "- dplyr"},
			Output: []string{"Packages: ", "- newpkg", "# some comment", "- dplyr"},
			Msg: "has comment before first package",
		},
		{
			Input: []string{"Packages: ", "# some comment", "  - dplyr"},
			Output: []string{"Packages: ", "  - newpkg", "# some comment", "  - dplyr"},
			Msg: "has comment before first package with indentation",
		},
		{
			Input: []string{"Packages: ", "- dplyr", "# a comment", "- ggplot2"},
			Output: []string{"Packages: ", "- newpkg", "- dplyr", "# a comment", "- ggplot2"},
			Msg: "has comment in the middle of packages",
		},
	}
	for _, tst := range tt {
		assert.Equal(t, tst.Output, insertPackages(tst.Input, "newpkg"), tst.Msg)
	}
}
