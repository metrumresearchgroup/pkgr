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
