package cmd

import (
	"fmt"
	"github.com/metrumresearchgroup/pkgr/configlib"
	"github.com/metrumresearchgroup/pkgr/testhelper"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func InitializeEmptyTestSiteWorking() {
	fileSystem := afero.NewOsFs()
	testWorkDir := filepath.Join("testsite", "working")
	fileSystem.RemoveAll(testWorkDir)
	fileSystem.MkdirAll(testWorkDir, 0755)
}

// Returns path to working test directory
func InitializeGoldenTestSiteWorking(goldenSet string) string {
	fileSystem := afero.NewOsFs()
	goldenSetPath := filepath.Join("testsite", "golden", goldenSet)
	testWorkDir := filepath.Join("testsite", "working")
	fileSystem.RemoveAll(testWorkDir)
	fileSystem.MkdirAll(testWorkDir, 0755)

	fileSystem.MkdirAll(testWorkDir, 0755)

	err := testhelper.CopyDir(fileSystem, goldenSetPath, testWorkDir)

	if err != nil {
		panic(err)
	}

	return filepath.Join(testWorkDir)
}

// Initializes a symlink to "test-library" in the testsite for the goldenSet. Links to the provided LibraryToUse and returns the symlink.
func InitializeTestLibrary(goldenSet, libraryToUse string) string {
	//fs := &afero.OsFs{}

	testLibrary := filepath.Join("testsite", "working", "test-library")
	err := os.Symlink(libraryToUse, testLibrary)
	//err := fs.SymlinkIfPossible(testLibrary, libraryToUse)
	if err != nil {
		fmt.Sprintf("error creating symlink for test library: %s", err)
		panic("could not symlink test-library")
	}
	return testLibrary

}

func InitGlobalConfig(libraryPath, localRepo string, update, suggests bool, installType string, packages []string) {



	cfg = configlib.PkgrConfig{
		Threads:  5,
		NoUpdate:   !update,
		NoRollback: false,
		Strict:   false,
		Packages: packages,
		Library:  libraryPath,
		Version:  1,
		//Logging: configlib.LogConfig{Level: "debug"}, // Controlled before cfg unmarshalling, I think
		Cache: "./testsite/working/localcache",
		Customizations: configlib.Customizations{
			Repos: []map[string]configlib.RepoConfig{{
				"testRepo": configlib.RepoConfig{
					Type: installType,
				},
			},
			},
		},
		//LibPaths: nil,
		//Lockfile: nil,
		Repos: []map[string]string{
			{
				"testRepo": localRepo,
			},
		},
		//RPath: nil,
		Suggests: suggests,
	}
}

func InitializeGlobalsForTest() {
	// Overwrite the global root cmd to "fake" the parts we need for cobra.
	RootCmd = &cobra.Command{
		Use:   "pkgr",
		Short: "package manager",
	}

	// Run the "set globals" function to init the "fs" object.
	setGlobals()
	//logger.SetLogLevel("trace") // Overwriting what's done in setGlobals() if we need to for testing.
}

func TestPackagesInstalled(t *testing.T) {

	type TestCase struct {
		localRepoName     string
		installUpdates    bool
		installSuggests   bool
		toInstall         []string // Equivalent to  "Packages" in pkgr.yml
		expectedInstalled []string
	}

	testCases := map[string]TestCase{
		"Basic Check": TestCase{
			localRepoName:   "simple",
			installUpdates:  false,
			installSuggests: false,
			toInstall: []string{
				"R6",
				"pillar",
			},
			expectedInstalled: []string{
				"assertthat",
				"cli",
				"crayon",
				"fansi",
				"pillar",
				"R6",
				"rlang",
				"utf8",
			},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {

			// Setup
			InitializeEmptyTestSiteWorking()
			InitializeGlobalsForTest()

			libraryPath := filepath.Join("testsite", "working", "libs")
			localRepoPath, err := filepath.Abs(filepath.Join("..", "localrepos", tc.localRepoName))
			checkError(t, err)
			InitGlobalConfig(libraryPath, localRepoPath, tc.installUpdates, tc.installSuggests, "source", tc.toInstall)

			// Execution
			_ = rInstall(nil, []string{})

			//Validation
			libSubDirectories, err := afero.ReadDir(fs, libraryPath)
			checkError(t, err)
			numInstalled := len(libSubDirectories)
			assert.Equalf(t, len(tc.expectedInstalled), numInstalled, "Expected %d packages to be installed but found %d", len(tc.expectedInstalled), numInstalled)

			for _, p := range tc.expectedInstalled {
				assert.DirExists(t, filepath.Join(libraryPath, p), "Package missing from final results: "+p)
			}
		})
	}

}

func TestTarballInstall(t *testing.T) {

	libraryPath := filepath.Join("testsite", "working", "libs")
	localReposDir, err := filepath.Abs(filepath.Join("..", "localrepos"))
	checkError(t, err)

	type TestCase struct {
		localRepoName     string
		installUpdates    bool
		installSuggests   bool
		toInstall         []string // Equivalent to  "Packages" in pkgr.yml
		toInstallTarballs []string
		expectedInstalled []string
	}

	testCases := map[string]TestCase{
		"Tarball no dependencies": TestCase{
			localRepoName:   "simple-no-R6",
			installUpdates:  false,
			installSuggests: false,
			toInstall: []string{
				"pillar",
			},
			toInstallTarballs: []string{
				filepath.Join(localReposDir, "tarballs", "R6_2.4.0.tar.gz"),
			},
			expectedInstalled: []string{
				"assertthat",
				"cli",
				"crayon",
				"fansi",
				"pillar",
				"R6", //Should be installed through tarball
				"rlang",
				"utf8",
			},
		},
		"Tarball with dependencies": TestCase{
			localRepoName:   "tarball-deps",
			installUpdates:  false,
			installSuggests: false,
			toInstall: []string{
				"openssl",
			},
			toInstallTarballs: []string{
				filepath.Join(localReposDir, "tarballs", "pillar_1.3.1.tar.gz"),
			},
			expectedInstalled: []string{
				"askpass", // Pulled in by openssl.
				"assertthat",
				"cli",
				"crayon",
				"fansi",
				"openssl",
				"pillar",
				"rlang",
				"sys", // Pulled in by openssl.
				"utf8",
			},
		},
		// Right after release 1.0.0, we found a bug on a specific tarball.
		// For this tarball, the root archive folder was not the first thing read by our tarball reader. This caused problems,
		// as our tarball-unpacker function would try to extract files (which it found first) into folders that we hadn't created
		// yet.
		// This automated test is meant to regression-test the aforementioned bug. We created the test file for this by using
		// "devtools::build()" on MacOS Mojave for the R6 package, pulled from GitHub (commit 8e0b3182cdcc5047343e9c590816578472ec9dfa)
		"Unordered tarball": TestCase{
			localRepoName:   "simple-no-R6",
			installUpdates:  false,
			installSuggests: false,
			toInstall: []string{
				"utf8",
			},
			toInstallTarballs: []string{
				filepath.Join(localReposDir, "tarballs", "R6_unordered_tarball_read.tar.gz"),
			},
			expectedInstalled: []string{
				"utf8",
				"R6",
			},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {

			// Setup
			InitializeEmptyTestSiteWorking()
			InitializeGlobalsForTest()

			localRepoPath := filepath.Join(localReposDir, tc.localRepoName)

			InitGlobalConfig(libraryPath, localRepoPath, tc.installUpdates, tc.installSuggests, "source", tc.toInstall)
			cfg.Tarballs = tc.toInstallTarballs

			// Execution
			_ = rInstall(nil, []string{})

			//Validation
			libSubDirectories, err := afero.ReadDir(fs, libraryPath)
			checkError(t, err)
			numInstalled := len(libSubDirectories)
			assert.Equalf(t, len(tc.expectedInstalled), numInstalled, "Expected %d packages to be installed but found %d", len(tc.expectedInstalled), numInstalled)

			for _, p := range tc.expectedInstalled {
				assert.DirExists(t, filepath.Join(libraryPath, p), "Package missing from final results: "+p)
			}
		})
	}

}

func TestInstallWithoutRollback(t *testing.T) {
	// Setup
	_ = InitializeGoldenTestSiteWorking("rollback-disabled")
	testLibrary := filepath.Join("testsite", "working", "libs")

	// Overwrite the global root cmd to "fake" the parts we need for cobra.
	RootCmd = &cobra.Command{
		Use:   "pkgr",
		Short: "package manager",
	}

	// Run the "set globals" function to init the "fs" object.
	setGlobals()

	// Create a fake config (will work for commands that don't use viper.Get[...])
	cfg = configlib.PkgrConfig{
		Threads:  5,
		NoUpdate:   false,
		NoRollback: true,
		Strict:   false,
		Packages: []string{"xml2", "crayon", "R6", "Rcpp", "crayon", "fansi", "flatxml"},
		Library:  testLibrary,
		Version:  1,
		//Logging: nil,
		//Cache: nil,
		Customizations: configlib.Customizations{
			Repos: []map[string]configlib.RepoConfig{
				{
					"local58": configlib.RepoConfig{
						Type: "source",
					},
				},
			},
		},
		//LibPaths: nil,
		//Lockfile: nil,
		Repos: []map[string]string{
			{
				"local58": "../localrepos/bad-xml2",
			},
		},
		//RPath: nil,
		Suggests: false,
	}

	// Run the actual test
	// Are we supposed to pass in RootCmd?
	_ = rInstall(nil, []string{})

	//Verify things look as we expect

	// Regular packages (either were installed during run or were preinstalled and up to date)
	assert.DirExists(t, filepath.Join(testLibrary, "bitops"), "Package missing from final results")
	assert.DirExists(t, filepath.Join(testLibrary, "crayon"), "Package missing from final results")
	assert.DirExists(t, filepath.Join(testLibrary, "RCurl"), "Package missing from final results")
	assert.DirExists(t, filepath.Join(testLibrary, "fansi"), "Package missing from final results")

	// Preinstalled packages not managed by pkgr
	assert.DirExists(t, filepath.Join(testLibrary, "utf8"), "Preinstalled, non-pkgr package missing from final results")

	// Outdated packages are still updated
	assert.DirExists(t, filepath.Join(testLibrary, "R6"), "Package missing from final results")
	fileExistsCheck, _ := afero.Exists(fs, filepath.Join(testLibrary, "R6", "THIS_PACKAGE_IS_OUTDATED"))
	assert.False(t, fileExistsCheck)

	assert.DirExists(t, filepath.Join(testLibrary, "Rcpp"), "Package missing from final results")
	fileExistsCheck, _ = afero.Exists(fs, filepath.Join(testLibrary, "Rcpp", "THIS_PACKAGE_IS_OUTDATED"))
	assert.False(t, fileExistsCheck)

	//Fail to install
	dirExistsCheck, _ := afero.DirExists(fs, filepath.Join(testLibrary, "xml2"))
	assert.False(t, dirExistsCheck, "Package was not properly removed or was installed when it shouldn't have been")
	dirExistsCheck, _ = afero.DirExists(fs, filepath.Join(testLibrary, "flatxml"))
	assert.False(t, dirExistsCheck, "Package was not properly removed or was installed when it shouldn't have been")
}

// Utility

func checkError(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}
