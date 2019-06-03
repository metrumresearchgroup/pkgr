package rcmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstallArgs(t *testing.T) {
	assert := assert.New(t)

	var installArgsTests = []struct {
		in       InstallArgs
		expected []string
	}{
		{
			InstallArgs{},
			[]string{},
		},
		{
			InstallArgs{Clean: true},
			[]string{"--clean"},
		},
		{
			InstallArgs{Library: "path/to/lib"},
			[]string{"--library=path/to/lib"},
		},
		{
			NewDefaultInstallArgs(),
			[]string{"--build", "--install-tests", "--no-multiarch", "--with-keep.source"},
		},
	}
	for i, tt := range installArgsTests {
		actual := tt.in.CliArgs()
		assert.Equal(tt.expected, actual, fmt.Sprintf("test num: %v", i+1))
	}
}

func TestParseDescriptionFile(t *testing.T) {
	var tests = []struct {
		filename string
		fields   []string
		expected map[string]string
		message  string
	}{
		{
			filename: "../integration_tests/simple/test-library/R6/Description",
			fields: []string{
				"Repository",
				"URL",
				"NeedsCompilation",
			},
			expected: map[string]string{
				"Repository":       "CRAN",
				"URL":              "https://r6.r-lib.org, https://github.com/r-lib/R6/",
				"NeedsCompilation": "no",
			},
			message: "simple test",
		},
		// {
		// 	filename: "../integration_tests/simple/test-library/rlang/Description",
		// 	fields: []string{
		// 		"Repository",
		// 		"URL",
		// 		"NeedsCompilation",
		// 	},
		// 	expected: map[string]string{
		// 		"Repository":       "CRAN",
		// 		"URL":              "http://rlang.r-lib.org, https://github.com/r-lib/rlang",
		// 		"NeedsCompilation": "no",
		// 	},
		// 	message: "simple test",
		// },
	}
	for _, tt := range tests {
		actual, err := parseDescriptionFile(tt.filename, tt.fields)
		fail := false

		if err != nil {
			fail = true
		} else {

			for _, field := range tt.fields {
				if tt.expected[field] != actual[field] {
					fail = true
				}
			}
		}

		assert.Equal(t, false, fail, fmt.Sprintf("Failed: %s", tt.message))
	}
}
