package configlib

import (
	"fmt"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestAddRemovePackage(t *testing.T) {
	tests := []struct {
		fileName    string
		packageName string
	}{
		{
			fileName:    "../integration_tests/simple/pkgr.yml",
			packageName: "packageTestName",
		},
		{
			fileName:    "../integration_tests/simple-suggests/pkgr.yml",
			packageName: "packageTestName",
		},
		{
			fileName:    "../integration_tests/mixed-source/pkgr.yml",
			packageName: "packageTestName",
		},
		{
			fileName:    "../integration_tests/outdated-pkgs/pkgr.yml",
			packageName: "packageTestName",
		},
		{
			fileName:    "../integration_tests/outdated-pkgs-no-update/pkgr.yml",
			packageName: "packageTestName",
		},
		{
			fileName:    "../integration_tests/repo-order/pkgr.yml",
			packageName: "packageTestName",
		},
	}

	appFS := afero.NewOsFs()
	for _, tt := range tests {
		b, _ := afero.Exists(appFS, tt.fileName)
		assert.Equal(t, true, b, fmt.Sprintf("yml file not found:%s", tt.fileName))

		ymlStart, _ := afero.ReadFile(appFS, tt.fileName)

		add(tt.fileName, tt.packageName)
		b, _ = afero.FileContainsBytes(appFS, tt.fileName, []byte(tt.packageName))
		assert.Equal(t, true, b, fmt.Sprintf("Package not added:%s", tt.fileName))

		remove(tt.fileName, tt.packageName)
		b, _ = afero.FileContainsBytes(appFS, tt.fileName, []byte(tt.packageName))
		assert.Equal(t, false, b, fmt.Sprintf("Package not removed:%s", tt.fileName))

		ymlEnd, _ := afero.ReadFile(appFS, tt.fileName)
		b = equal(ymlStart, ymlEnd)
		assert.Equal(t, true, b, fmt.Sprintf("Start and End yml files differ:%s", tt.fileName))
	}
}

func equal(a []byte, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
