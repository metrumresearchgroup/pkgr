package rcmd

import (
	"fmt"
	"testing"

	"github.com/metrumresearchgroup/pkgr/cran"
	"github.com/stretchr/testify/assert"
)

func TestLibPathsEnv(t *testing.T) {
	var libPathTests = []struct {
		in       RSettings
		expected string
	}{
		{
			RSettings{
				LibPaths: []string{
					// TODO: check if paths need to be checked for trailing /
					"path/to/folder1/",
					"path/to/folder2/",
				},
			},
			"R_LIBS_SITE=path/to/folder1/:path/to/folder2/",
		},
		{
			RSettings{
				LibPaths: []string{},
			},
			"",
		},
	}
	for _, tt := range libPathTests {
		ok, actual := tt.in.LibPathsEnv()
		if actual != "" && !ok {
			t.Errorf("LibPaths present, should be ok")
		}
		if actual != tt.expected {
			t.Errorf("GOT: %s, WANT: %s", actual, tt.expected)
		}
	}
}

func TestGetRVersion(t *testing.T) {
	var rVersionTests = []struct {
		version    cran.RVersion
		platform   string
		message    string
		assertTrue bool
	}{
		{
			cran.RVersion{
				Major: 3,
				Minor: 5,
				Patch: 3,
			},
			"x86_64-apple-darwin15.6.0",
			"version info - accurate on mac",
			true,
		},
		{
			cran.RVersion{
				Major: 9,
				Minor: 9,
				Patch: 9,
			},
			"999_x86_64-apple-darwin15.6.0",
			"version info - wrong version",
			false,
		},
	}
	for _, tt := range rVersionTests {
		rSettings := NewRSettings("")
		version := GetRVersion(&rSettings)
		if tt.assertTrue {
			assert.Equal(t, tt.version, version, fmt.Sprintf("Not equal: %s", tt.message))
			assert.Equal(t, tt.version, rSettings.Version, fmt.Sprintf("Not equal: %s", tt.message))
			assert.Equal(t, tt.platform, rSettings.Platform, fmt.Sprintf("Not equal: %s", tt.message))
		} else {
			assert.NotEqual(t, tt.version, version, fmt.Sprintf("Equal: %s", tt.message))
			assert.NotEqual(t, tt.version, rSettings.Version, fmt.Sprintf("Equal: %s", tt.message))
			assert.NotEqual(t, tt.platform, rSettings.Platform, fmt.Sprintf("Equal: %s", tt.message))
		}
	}
}

func BenchmarkGetVersionInfo(b *testing.B) {
	rSettings := NewRSettings("")
	for n := 0; n < 10; n++ {
		GetRVersion(&rSettings)
	}
}
