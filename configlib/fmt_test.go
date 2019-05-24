package configlib

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestFormat(t *testing.T) {

	tests := []struct {
		ymlfolder string
		expected  []string
		message   string
	}{
		{
			ymlfolder: "../integration_tests/simple",
			expected: []string{
				"R6",
				"pillar",
			},
			message: "add whitespace",
		},
		{
			ymlfolder: "../integration_tests/simple",
			expected: []string{
				"R6",
				"pillar",
			},
			message: "remove whitespace",
		},

		{
			ymlfolder: "../integration_tests/mixed-source",
			expected: []string{
				"rmarkdown",
				"shiny",
				"devtools",
				"mrgsolve",
			},
			message: "add whitespace",
		},

		// fails because of the nested structure ...
		// {
		// 	ymlfolder: "../integration_tests/mixed-source",
		// 	expected: []string{
		// 		"rmarkdown",
		// 		"shiny",
		// 		"devtools",
		// 		"mrgsolve",
		// 	},
		// 	message: "remove whitespace",
		// },
	}
	viper.Reset()
	fs := afero.NewOsFs()
	for _, tt := range tests {

		ymlFile := filepath.Join(tt.ymlfolder, "pkgr.yml")
		b, _ := afero.Exists(fs, ymlFile)
		assert.Equal(t, true, b, fmt.Sprintf("yml file not found:%s", ymlFile))

		// capture initial state to restore
		ymlSave, _ := afero.ReadFile(fs, ymlFile)
		fi, _ := os.Stat(ymlFile)

		// get a formatted version of the initial file
		ymlStart, err := Format(ymlSave)
		assert.Equal(t, nil, err, "Error formatting starting yml file")

		// change the file
		var ymlTest []byte
		if tt.message == "add whitespace" {
			for _, line := range bytes.Split(ymlStart, []byte("\n")) {
				ymlTest = append(ymlTest, " "...)
				ymlTest = append(ymlTest, line...)
				ymlTest = append(ymlTest, "\n"...)
			}
		} else if tt.message == "remove whitespace" {

			lines := bytes.Split(ymlStart, []byte("\n"))
			for _, line := range lines {
				ymlTest = append(ymlTest, bytes.Trim(line, " ")...)
				ymlTest = append(ymlTest, "\n"...)
			}
		}

		// write the formatted changes to the original yml file
		ymlTest, err = Format(ymlTest)
		assert.Equal(t, nil, err, "Error formatting yml test file")

		err = afero.WriteFile(fs, ymlFile, ymlTest, fi.Mode())
		assert.Equal(t, nil, err, "Error writing yml test file")

		fmt.Println(string(ymlTest))
		fmt.Print(string(ymlStart))

		// viper-load the yml
		viper.Reset()
		viper.SetConfigName("pkgr")
		viper.AddConfigPath(tt.ymlfolder)
		err = viper.ReadInConfig()
		assert.Equal(t, nil, err, "Error reading yml file")

		// get package data from viper
		actual := viper.Get("Packages") //.([]interface{})
		expected := []interface{}{}
		for _, p := range tt.expected {
			expected = append(expected, p)
		}

		b = reflect.DeepEqual(actual, expected)
		assert.Equal(t, b, true, "Error reading yml cfg")

		// restore the original yml
		err = afero.WriteFile(fs, ymlFile, ymlSave, fi.Mode())
		assert.Equal(t, nil, err, "Error restoring yml file")
	}
}
