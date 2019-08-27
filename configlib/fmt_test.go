package configlib

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestFormat(t *testing.T) {

	var simple = []byte(`
Packages:
  - R6
  - pillar
# any repositories, order matters
Repos:
- CRAN: "https://cran.rstudio.com"
Library: "test-library"


`)

	var mixed = []byte(`
Version: 1
# top level packages
Packages:
  - rmarkdown
  - shiny
  - devtools
  - pmplots
# any repositories, order matters
Repos:
- CRAN: "https://cran.rstudio.com"
Library: "test-library"


`)

	tests := []struct {
		yml      []byte
		expected []string
		message  string
	}{
		{
			yml: simple,
			expected: []string{
				"R6",
				"pillar",
			},
			message: "add whitespace",
		},
		{
			yml: simple,
			expected: []string{
				"R6",
				"pillar",
			},
			message: "remove whitespace",
		},
		{
			yml: mixed,
			expected: []string{
				"rmarkdown",
				"shiny",
				"devtools",
				"pmplots",
			},
			message: "add whitespace",
		},
		{
			yml: mixed,
			expected: []string{
				"rmarkdown",
				"shiny",
				"devtools",
				"pmplots",
			},
			message: "remove whitespace",
		},
	}
	for _, tt := range tests {
		viper.Reset()
		viper.SetConfigType("yaml")

		if tt.message == "add whitespace" {
			tt.yml = append(tt.yml, []byte("\n\n")...)
		} else if tt.message == "remove whitespace" {
			tt.yml = []byte(strings.Replace(string(tt.yml), "\n\t", "", -1))
		}

		err := viper.ReadConfig(bytes.NewBuffer(tt.yml))
		assert.Equal(t, nil, err, fmt.Sprintf("%s", err))

		// get package data from viper
		actual := viper.Get("Packages")
		expected := []interface{}{}
		for _, p := range tt.expected {
			expected = append(expected, p)
		}
		b := reflect.DeepEqual(actual, expected)
		assert.Equal(t, b, true, "Error reading yml cfg")
	}
}
