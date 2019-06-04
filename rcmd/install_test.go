package rcmd

import (
	"fmt"
	"testing"

	"github.com/metrumresearchgroup/pkgr/desc"
	"github.com/spf13/afero"
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
			repoURL:     "?",
			repo:        "CRAN",
			message:     "R6 test",
		},
	}
	fs := afero.NewMemMapFs()
	//var err error
	for _, tt := range tests {

		dcf, _ := updateDcfFile(tt.filename, tt.version, tt.installType, tt.repoURL, tt.repo)
		descriptionFile, _ := fs.Create("descriptionFilePath")
		descriptionFile.WriteString(string(dcf))
		installedPackage, _ := desc.ParseDesc(descriptionFile)
		assert.Equal(t, tt.repo, installedPackage.Repository, fmt.Sprintf("Failed: %s", tt.message))
	}
}
