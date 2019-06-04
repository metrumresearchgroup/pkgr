package rcmd

import (
	"fmt"
	"testing"

	"github.com/spf13/afero"

	"github.com/metrumresearchgroup/pkgr/desc"
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

func TestUpdateDcfFile(t *testing.T) {
	var tests = []struct {
		filename    string
		version     string
		installType string
		repoURL     string
		repo        string
		message     string
	}{
		{
			filename:    "../integration_tests/simple/test-library/R6/Description",
			version:     "version",
			installType: "binary",
			repoURL:     "myURL",
			repo:        "CRAN",
			message:     "R6 test",
		},
	}

	osfs := afero.NewOsFs()
	dfcname := "/tmp/Description"
	for _, tt := range tests {

		dcf, err := updateDcfFile(tt.filename, tt.version, tt.installType, tt.repoURL, tt.repo)
		assert.Equal(t, nil, err, fmt.Sprintf("Failed: %s", tt.message))

		err = osfs.Remove(dfcname)
		assert.Equal(t, nil, err, fmt.Sprintf("Failed: %s", tt.message))

		descriptionFile, err := osfs.Create(dfcname)
		assert.Equal(t, nil, err, fmt.Sprintf("Failed: %s", tt.message))

		n, err := descriptionFile.WriteString(string(dcf))
		assert.Equal(t, nil, err, fmt.Sprintf("Failed: %s", tt.message))
		assert.NotEqual(t, 0, n, fmt.Sprintf("Failed: %s", tt.message))

		//installedPackage, err := desc.ParseDesc(descriptionFile)
		installedPackage, err := desc.ReadDesc(dfcname)
		assert.Equal(t, nil, err, fmt.Sprintf("Failed: %s", tt.message))

		assert.Equal(t, tt.repo, installedPackage.Repository, fmt.Sprintf("Failed: %s", tt.message))
		assert.Equal(t, tt.version, installedPackage.PkgrVersion, fmt.Sprintf("Failed: %s", tt.message))
		assert.Equal(t, tt.repoURL, installedPackage.PkgrRepositoryURL, fmt.Sprintf("Failed: %s", tt.message))
		assert.Equal(t, tt.installType, installedPackage.PkgrInstallType, fmt.Sprintf("Failed: %s", tt.message))
	}
}
