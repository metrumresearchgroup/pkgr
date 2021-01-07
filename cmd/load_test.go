package cmd

import (
	"fmt"
	"github.com/metrumresearchgroup/pkgr/logger"
	"github.com/metrumresearchgroup/pkgr/rcmd"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/thoas/go-funk"
	"path/filepath"
	"testing"
)

func TestLoadSucceeds(t *testing.T) {
	type testCase struct {
		expectedToLoad []string
		allOption      bool
	}

	testCases := map[string]testCase{
		"load user packages": testCase{
			expectedToLoad: []string{"R6", "pillar"},
			allOption:      false,
		},
		"load dependencies": testCase{
			expectedToLoad: []string{"R6", "pillar", "cli", "assertthat", "crayon", "fansi", "rlang", "utf8"},
			allOption:      true,
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {

			// Setup
			rDir := InitializeGoldenTestSiteWorking("load")
			libraryPath := InitializeTestLibrary("load", filepath.Join("testsite", "golden", "load", "test-library-4.0"))// set library to "test-library-3.0" if user R 3.X
			InitializeGlobalsForTest()

			localRepoPath := filepath.Join("..", "localrepos", "simple")

			userPackages := []string{
				"R6",
				"pillar",
			}

			InitGlobalConfig(libraryPath, localRepoPath, false, false, "source", userPackages)
			rSettings := rcmd.NewRSettings(cfg.RPath)

			/// Create a test logrus hook to check logs
			_, hook := test.NewNullLogger()
			logrus.AddHook(hook)
			logger.SetLogLevel("DEBUG")

			// Execution
			load(userPackages, rSettings, rDir, cfg.Threads, tc.allOption, false)

			// Validation
			var loadedPackages []string
			packageLoadSuccessful := false
			for _, entry := range hook.AllEntries() {
				if entry.Message == "all packages loaded successfully" {
					packageLoadSuccessful = true
				}
				if entry.Message == "Package loaded successfully" {
					pkg := entry.Data["pkg"]
					loadedPackages = append(loadedPackages, fmt.Sprintf("%v", pkg))
				}
			}
			assert.Equal(t, len(tc.expectedToLoad), len(loadedPackages), "length of expected packages and loaded packages not equal")
			for _, p := range loadedPackages {
				assert.Contains(t, tc.expectedToLoad, p, "expected package not loaded")
			}
			assert.True(t, packageLoadSuccessful, "packages were not successfully loaded")
		})
	}
}

func TestLoadFails(t *testing.T) {
	type testCase struct {
		expectedToLoad []string
		expectedToFail []string
		allOption      bool
	}

	testCases := map[string]testCase{
		"load user packages": testCase{
			expectedToLoad: []string{"pillar"},
			expectedToFail: []string{"R6"},
			allOption:      false,
		},
		"load dependencies": testCase{
			expectedToLoad: []string{"pillar", "cli", "assertthat", "crayon", "rlang"},
			expectedToFail: []string{"R6", "fansi", "utf8"},
			allOption:      true,
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {

			// Setup
			rDir := InitializeGoldenTestSiteWorking("load-fail")
			InitializeGlobalsForTest()
			libraryPath := InitializeTestLibrary("load-fail", filepath.Join("testsite", "golden", "load-fail", "test-library-4.0"))// set library to "test-library-3.0" if user R 3.X
			localRepoPath := filepath.Join("..", "localrepos", "simple")

			userPackages := []string{
				"R6",
				"pillar",
			}

			InitGlobalConfig(libraryPath, localRepoPath, false, false, "source", userPackages)
			rSettings := rcmd.NewRSettings(cfg.RPath)

			/// Create a test logrus hook to check logs
			_, hook := test.NewNullLogger()
			logrus.AddHook(hook)
			logger.SetLogLevel("DEBUG")

			// Execution
			load(userPackages, rSettings, rDir, cfg.Threads, tc.allOption, false)

			// Validation
			var loadedPackages []string
			var failedPackages []string

			packageLoadFailed := false
			for _, entry := range hook.AllEntries() {
				if entry.Message == "some packages failed to load." {
					packageLoadFailed = true
				} else if entry.Message == "Package loaded successfully" {
					pkg := entry.Data["pkg"]
					loadedPackages = append(loadedPackages, fmt.Sprintf("%v", pkg))
				} else if funk.Contains(entry.Message, "error loading package via") {
					pkg := entry.Data["pkg"]
					failedPackages = append(failedPackages, fmt.Sprintf("%v", pkg))
				}
			}
			assert.Equal(t, len(tc.expectedToLoad), len(loadedPackages), "length of expected packages and loaded packages not equal")
			for _, p := range loadedPackages {
				assert.Contains(t, tc.expectedToLoad, p, "package did not load")
			}
			assert.Equal(t, len(tc.expectedToFail), len(failedPackages), "length of expected failure packages and actual failed packages not equal")
			for _, p := range failedPackages {
				assert.Contains(t, tc.expectedToFail, p, "package did not fail")
			}

			assert.True(t, packageLoadFailed, "packages succeeded in loading where they should have failed")
		})
	}
}
