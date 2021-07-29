package library

import (
	"context"
	"github.com/metrumresearchgroup/command"
	. "github.com/metrumresearchgroup/pkgr/testhelper"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

const (
	librariesE2ETest1 = "LIB-E2E-001"
	librariesE2ETest2 = "LIB-E2E-002"
	librariesE2ETest3 = "LIB-E2E-003"
	librariesE2ETest4 = "LIB-E2E-004"
)

func TestLibrary(t *testing.T) {
	t.Run(MakeTestName(librariesE2ETest1, "strict mode stops pkgr when library doesn't exist"), func(t *testing.T) {
		DeleteTestFolder(t, "test-library")
		DeleteTestFolder(t, "test-cache")
		ctx := context.TODO()
		installCmd := command.New()

		capture, err := installCmd.Run(ctx, "pkgr", "install", "--config=pkgr-strict.yml", "--logjson")
		assert.Error(t, err, "install succeeded, but it should have failed due to strict mode")
		assert.NoDirExists(t, "test-library", "test-library was created when it should not have been")
		//t.Log(string(capture.Output))
		logs := CollectGenericLogs(t, capture, "library directory must exist before running pkgr in strict mode")
		assert.Len(t, logs, 2, "expected exactly one error and one fatal message")
		for _, log := range logs {
			assert.Equal(t, "test-library", log.Library)
		}
	})

	t.Run(MakeTestName(librariesE2ETest2, "lockfile type renv installs to renv/library"), func(t *testing.T) {
		DeleteTestFolder(t, "test-cache")
		SetupEndToEndWithInstall(t, "pkgr-renv.yml", "renv")

		r6FolderFound := false
		err := filepath.Walk("renv/library", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() && info.Name() == "R6" {
				r6FolderFound = true
			}
			return nil
		})
		if err != nil {
			t.Fatalf("error when walking renv folder to find installed pkgs: %s", err)
		}
		assert.True(t, r6FolderFound, "failed to find installation of R6")
	})

	t.Run(MakeTestName(librariesE2ETest3, "lockfile type packrat installs to packrat/library"), func(t *testing.T) {
		DeleteTestFolder(t, "test-cache")
		SetupEndToEndWithInstall(t, "pkgr-packrat.yml", "packrat")

		r6FolderFound := false
		err := filepath.Walk("packrat/lib", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() && info.Name() == "R6" {
				r6FolderFound = true
			}
			return nil
		})
		if err != nil {
			t.Fatalf("error when walking renv folder to find installed pkgs: %s", err)
		}
		assert.True(t, r6FolderFound, "failed to find installation of R6")
	})
}