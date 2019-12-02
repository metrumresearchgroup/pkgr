package configlib

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// This test passes when executed from this file or executed alone using:
// /usr/local/go/bin/go test -timeout 30s github.com/metrumresearchgroup/pkgr/configlib -run TestViper

// This test INTERMITTENTLY fails using "go test ./..." with the below error.
// Debugging shows the file "integration_tests/simple/pkgr.yml" contains two extra packages "abc" and
// "shiny", however the file on the file system does not contain those packages.

// --- FAIL: TestViper (0.00s)
//     viper_test.go:51:
//         	Error Trace:	viper_test.go:51
//         	Error:      	Not equal:
//         	            	expected: false
//         	            	actual  : true
//         	Test:       	TestViper
//         	Messages:   	Package not equal.
//         	            	Expected:[R6 pillar]
//         	            	Actual:[abc shiny R6 pillar]

// ATTENTION:
// This test is misconfigured, it shouldn't be pointing to the integration test folders as those folders are only valid
// after make test-install has been run from the integration_tests folder.
// This test might fail falsey depending on the current state of pkgr/integration_tests/simple/test-library
func TestViper(t *testing.T) {
	tests := []struct {
		ymlfolder string
		message   string
		expected  []string
	}{
		{
			ymlfolder: "simple",
			expected: []string{
				"R6",
				"pillar",
			},
		},
	}

	t.Skip("test is dependent on the state of an integration_test folder, meaning it cannot be trusted to produce consistent results")

	fs := afero.NewOsFs()
	for _, tt := range tests {

		// get the yml file
		ymlFile := filepath.Join(getTestFolder(t, tt.ymlfolder), "pkgr.yml")
		b, _ := afero.Exists(fs, ymlFile)
		assert.Equal(t, true, b, fmt.Sprintf("yml file not found:%s", ymlFile))
		// read text for error output
		text, err := ioutil.ReadFile(ymlFile)
		assert.Equal(t, nil, err, fmt.Sprintf("error reading yml file:%s", ymlFile))
		fmt.Println(string(text))

		// config viper
		viper.Reset()
		viper.SetConfigFile(ymlFile)
		viper.ReadInConfig()

		// get package data from viper
		actual := viper.Get("Packages")
		expected := []interface{}{}
		for _, p := range tt.expected {
			expected = append(expected, p)
		}

		help := "Run this to view file:\n /bin/cat " + ymlFile
		b = reflect.DeepEqual(actual, expected)
		assert.Equal(t, b, true, "Packages not equal\nExpected:%s\nActual:%s\nStart_file<\n%s\n>End_file\n\n%s\n",
			expected, actual, text, help)
	}
}
