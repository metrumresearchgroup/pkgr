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

func TestParseVersionData(t *testing.T) {
	var rVersionTests = []struct {
		data     []byte
		version  cran.RVersion
		platform string
		message  string
	}{
		{
			data: []byte(`R version 3.5.3 (2019-03-11) -- "Great Truth"
Copyright (C) 2019 The R Foundation for Statistical Computing
Platform: x86_64-apple-darwin15.6.0 (64-bit)

R is free software and comes with ABSOLUTELY NO WARRANTY.
You are welcome to redistribute it under the terms of the
GNU General Public License versions 2 or 3.
For more information about these matters see
http://www.gnu.org/licenses/.

`),
			version: cran.RVersion{
				Major: 3,
				Minor: 5,
				Patch: 3,
			},
			platform: "x86_64-apple-darwin15.6.0",
			message:  "darwin test",
		},
		{
			data: []byte(`R version 3.5.2 (2018-12-20) -- "Eggshell Igloo"
Copyright (C) 2018 The R Foundation for Statistical Computing
Platform: x86_64-w64-mingw32/x64 (64-bit)
			
R is free software and comes with ABSOLUTELY NO WARRANTY.
You are welcome to redistribute it under the terms of the
GNU General Public License versions 2 or 3.
For more information about these matters see
http://www.gnu.org/licenses/.

`),
			version: cran.RVersion{
				Major: 3,
				Minor: 5,
				Patch: 2,
			},
			platform: "x86_64-w64-mingw32/x64",
			message:  "windows test",
		},
		{
			data: []byte(`
			R version 1.2.3 (2018-12-20) -- "name for Ubuntu"            
			Platform: x86_64-pc-linux-gnu (64-bit)
			`),
			version: cran.RVersion{
				Major: 1,
				Minor: 2,
				Patch: 3,
			},
			platform: "x86_64-pc-linux-gnu",
			message:  "Manually built Ubuntu test",
		},
	}
	for _, tt := range rVersionTests {
		version, platform := parseVersionData(tt.data)
		assert.Equal(t, tt.version, version, fmt.Sprintf("Version not equal: %s", tt.message))
		assert.Equal(t, tt.platform, platform, fmt.Sprintf("Platform not equal: %s", tt.message))
	}
}
