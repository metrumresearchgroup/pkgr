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
			data: []byte(`
			platform       x86_64-apple-darwin15.6.0
			arch           x86_64
			os             darwin15.6.0
			system         x86_64, darwin15.6.0
			status
			major          3
			minor          5.3
			year           2019
			month          03
			day            11
			svn rev        76217
			language       R
			version.string R version 3.5.3 (2019-03-11)
			nickname       Great Truth
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
			data: []byte(`
			platform       i386-w64-mingw32            
			arch           i386                        
			os             mingw32                     
			system         i386, mingw32               
			status                                     
			major          3                           
			minor          1.2                         
			year           2014                        
			month          10                          
			day            31                          
			svn rev        66913                       
			language       R                           
			version.string R version 3.1.2 (2014-10-31)
			nickname       Pumpkin Helmet              
			`),
			version: cran.RVersion{
				Major: 3,
				Minor: 1,
				Patch: 2,
			},
			platform: "i386-w64-mingw32",
			message:  "windows test",
		},
		{
			data: []byte(`
			platform       x86_64-pc-linux-gnu (64-bit)            
			version.string R version 3.4.4 (2018-03-15)
			`),
			version: cran.RVersion{
				Major: 3,
				Minor: 4,
				Patch: 4,
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
