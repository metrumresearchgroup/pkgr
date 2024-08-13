package library

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/metrumresearchgroup/command"
	"github.com/stretchr/testify/assert"

	. "github.com/metrumresearchgroup/pkgr/testhelper"
)

const (
	librariesE2ETest1 = "LIB-E2E-001"
	librariesE2ETest2 = "LIB-E2E-002"
	librariesE2ETest3 = "LIB-E2E-003"
	librariesE2ETest4 = "LIB-E2E-004"
	librariesE2ETest6 = "LIB-E2E-006"
	librariesE2ETest7 = "LIB-E2E-007"
)

func TestLibrary(t *testing.T) {
	t.Run(MakeTestName(librariesE2ETest1, "strict mode stops pkgr when library doesn't exist"), func(t *testing.T) {
		DeleteTestFolder(t, "test-library")
		DeleteTestFolder(t, "test-cache")

		installCmd := command.New("pkgr", "install", "--config=pkgr-strict.yml", "--logjson")

		capture, err := installCmd.CombinedOutput()
		assert.Error(t, err, "install succeeded, but it should have failed due to strict mode")
		assert.NoDirExists(t, "test-library", "test-library was created when it should not have been")
		//t.Log(string(capture.Output))
		logs := CollectGenericLogs(t, capture, "library directory must exist before running pkgr in strict mode")
		assert.Len(t, logs, 2, "expected exactly one error and one fatal message")
		for _, log := range logs {
			assert.Equal(t, "test-library", log.Library)
		}
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

func TestLibraryRenv(t *testing.T) {
	renv := os.Getenv("PKGR_TESTS_SYS_RENV")
	expectError := renv == ""

	t.Run(MakeTestName(librariesE2ETest2, "lockfile type renv installs to renv/library"), func(t *testing.T) {
		DeleteTestFolder(t, "test-cache")
		cmdout, cmderr := SetupEndToEndWithInstallFull(t, "pkgr-renv.yml", "renv",
			nil, "", expectError)

		if renv == "" {
			assert.Error(t, cmderr, "expected 'pkgr install' error")
			assert.Contains(t, cmdout, "calling renv to find library path failed")
			t.Skip("Skipping: empty PKGR_TESTS_SYS_RENV indicates renv not available")
		} else {
			r6FolderFound := false
			err := filepath.Walk("renv/library",
				func(path string, info os.FileInfo, err error) error {
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
		}
	})

	t.Run("lockfile type renv errors renv if is unavailable", func(t *testing.T) {
		if renv != "" {
			t.Skip("Skipping: non-empty PKGR_TESTS_SYS_RENV indicates renv is available")
		}

		DeleteTestFolder(t, "test-cache")
		cmdout, _ := SetupEndToEndWithInstallFull(t, "pkgr-renv.yml", "renv",
			nil, "", true)

		assert.Contains(t, cmdout, "calling renv to find library path failed")
	})

	t.Run(MakeTestName(librariesE2ETest6, "lockfile type renv respects RENV_PATHS_LIBRARY_ROOT"), func(t *testing.T) {
		DeleteTestFolder(t, "test-cache")

		lib := "renv-custom"
		env := append(os.Environ(), "RENV_PATHS_LIBRARY_ROOT="+lib)
		cmdout, _ := SetupEndToEndWithInstallFull(t, "pkgr-renv.yml", lib,
			env, "", expectError)

		if expectError {
			assert.Contains(t, cmdout, "calling renv to find library path failed")
			t.Skip("Skipping: empty PKGR_TESTS_SYS_RENV indicates renv not available")
		} else {
			r6FolderFound := false
			err := filepath.Walk(filepath.Join(lib),
				func(path string, info os.FileInfo, err error) error {
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
		}
	})

	t.Run(MakeTestName(librariesE2ETest7, "lockfile type renv detects package project library"), func(t *testing.T) {
		DeleteTestFolder(t, "test-cache")
		DeleteTestFolder(t, "pkg")

		err := os.MkdirAll(filepath.Join("pkg", "renv"), 0777)
		if err != nil {
			t.Fatalf("os.MkdirAll() error: %s", err)
		}

		err = os.WriteFile(filepath.Join("pkg", ".Rprofile"),
			[]byte("cat('i do not interfere\n')\nsource('renv/activate.R')\n"),
			0666)
		if err != nil {
			t.Fatalf("os.WriteFile() error: %s", err)
		}

		err = os.WriteFile(filepath.Join("pkg", "DESCRIPTION"),
			[]byte("Package: pkg\n"),
			0666)
		if err != nil {
			t.Fatalf("os.WriteFile() error: %s", err)
		}

		err = os.Link(filepath.Join("testdata", "activate.R"),
			filepath.Join("pkg", "renv", "activate.R"))
		if err != nil {
			t.Fatalf("os.Link() error: %s", err)
		}

		cwd, err := os.Getwd()
		if err != nil {
			t.Fatalf("os.Getwd() error: %s", err)
		}

		libcache := filepath.Join(cwd, "r-user-cache")
		env := append(os.Environ(), "R_USER_CACHE_DIR="+libcache)
		dir := "pkg"

		SetupEndToEndWithInstallFull(t,
			filepath.Join(cwd, "pkgr-renv.yml"),
			libcache, env, dir, false)

		r6FolderFound := false
		err = filepath.Walk(filepath.Join(libcache),
			func(path string, info os.FileInfo, err error) error {
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

		loadcmd := command.New("Rscript", "-e", "library(R6)")
		loadcmd.Env = env
		loadcmd.Dir = dir
		err = loadcmd.Run()
		if err != nil {
			t.Fatalf("loading R6 failed: %s", err)
		}
	})
}
