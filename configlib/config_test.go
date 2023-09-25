package configlib

import (
	"bytes"
	"fmt"
	"github.com/metrumresearchgroup/pkgr/testhelper"
	"github.com/thoas/go-funk"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/metrumresearchgroup/pkgr/cran"
	"github.com/metrumresearchgroup/pkgr/gpsr"
	"github.com/metrumresearchgroup/pkgr/rcmd"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// Test names for tests that are relevant to validation.
const (
	configUnitTest1 = "CFG-UNIT-001"
	configUnitTest2 = "CFG-UNIT-002"
	configUnitTest3 = "CFG-UNIT-003"
	configUnitTest4 = "CFG-UNIT-004"
	configUnitTest5 = "CFG-UNIT-005"
	configUnitTest6 = "CFG-UNIT-006"
	configUnitTest7 = "CFG-UNIT-007"
	configUnitTest8 = "CFG-UNIT-008"
)

func TestExpandTilde(t *testing.T) {
	t.Run(testhelper.MakeTestName(configUnitTest1, "Test expand tildes 1"), func(t *testing.T) {
		homeDirectory, err := os.UserHomeDir()
		if err != nil {
			t.Fatal("error determining home directory to use for testing: ", err)
		}

		type testCase struct {
			path           string
			expectedResult string
			testSubId 	   string
		}

		tests := map[string]testCase{
			"expands tilde": {
				path:           filepath.Join("~/Desktop/folderA"),
				expectedResult: filepath.Join(homeDirectory, "Desktop", "folderA"),
				testSubId: "1",
			},
			"does not modify regular path": {
				path:           filepath.Join("A/B/C"),
				expectedResult: filepath.Join("A/B/C"),
				testSubId: "2",
			},
			"does not modify local path": {
				path:           filepath.Join("../A/B/C"),
				expectedResult: filepath.Join("../A/B/C"),
				testSubId: "3",
			},
			"tilde must be prefix": {
				path:           filepath.Join("A/B/~/C"),
				expectedResult: filepath.Join("A/B/~/C"),
				testSubId: "4",
			},
			"works with empty path": {
				path:           "",
				expectedResult: "",
				testSubId: "5",
			},
		}

		for testName, tc := range tests {
			t.Run(testhelper.MakeSubtestName(configUnitTest1, tc.testSubId, testName), func(t *testing.T) {
				actualResult := expandTilde(tc.path)
				assert.Equal(t, tc.expectedResult, actualResult)
			})
		}
	})

	t.Run(testhelper.MakeTestName(configUnitTest2, "Test expand tildes 2"), func(t *testing.T) {
		homeDirectory, err := os.UserHomeDir()
		if err != nil {
			t.Fatal("error determining home directory to use for testing: ", err)
		}

		type testCase struct {
			paths           []string
			expectedResults []string
			testSubId 	   string
		}

		tests := map[string]testCase{
			"expands tildes": {
				paths: []string{
					filepath.Join("~/Desktop/folderA"),
					filepath.Join("~/Documents/folderB"),
				},
				expectedResults: []string{
					filepath.Join(homeDirectory, "Desktop", "folderA"),
					filepath.Join(homeDirectory, "Documents", "folderB"),
				},
				testSubId : "1",
			},
			"expands tildes but not others": {
				paths: []string{
					filepath.Join("~/Desktop/folderA"),
					filepath.Join("/TopDir/Documents/folderB"),
				},
				expectedResults: []string{
					filepath.Join(homeDirectory, "Desktop", "folderA"),
					filepath.Join("/TopDir", "Documents", "folderB"),
				},
				testSubId : "2",
			},
			"does not modify non-tilde repos": {
				paths: []string{
					filepath.Join("A", "B", "C"),
					filepath.Join("D"),
					filepath.Join("/srv", "shiny-server", "log.txt"),
					"",
					filepath.Join("..", "E", "F"),
				},
				expectedResults: []string{
					filepath.Join("A", "B", "C"),
					filepath.Join("D"),
					filepath.Join("/srv", "shiny-server", "log.txt"),
					"",
					filepath.Join("..", "E", "F"),
				},
				testSubId : "3",
			},
		}

		for testName, tc := range tests {
			t.Run(testhelper.MakeSubtestName(configUnitTest1, string(tc.testSubId), testName), func(t *testing.T) {
				actualResults := expandTildes(tc.paths)
				assert.Equal(t, tc.expectedResults, actualResults)
			})
		}
	})

	t.Run(testhelper.MakeTestName(configUnitTest3, "Test expand tildes 3"), func(t *testing.T) {
		homeDirectory, err := os.UserHomeDir()
		if err != nil {
			t.Fatal("error determining home directory to use for testing: ", err)
		}

		type testCase struct {
			repos           []map[string]string
			expectedResults []map[string]string
			testSubId 		string
		}

		tests := map[string]testCase{
			"expands tildes": {
				repos: []map[string]string{
					{"A": filepath.Join("~/Desktop/folderA")},
					{"B": filepath.Join("~/Documents/folderB")},
				},
				expectedResults: []map[string]string{
					{"A": filepath.Join(homeDirectory, "Desktop/folderA")},
					{"B": filepath.Join(homeDirectory, "Documents/folderB")},
				},
				testSubId : "1",
			},
			"expands tildes but not others": {
				repos: []map[string]string{
					{"A": filepath.Join("~/Desktop/folderA")},
					{"B": filepath.Join("/TopDir/Documents/folderB")},
				},
				expectedResults: []map[string]string{
					{"A": filepath.Join(homeDirectory, "Desktop", "folderA")},
					{"B": filepath.Join("/TopDir", "Documents", "folderB")},
				},
				testSubId : "2",
			},
			"does not modify non-tilde repos": {
				repos: []map[string]string{
					{"1": filepath.Join("A", "B", "C")},
					{"2": filepath.Join("D")},
					{"3": filepath.Join("/srv", "shiny-server", "log.txt")},
					{"4": ""},
					{"5": filepath.Join("..", "E", "F")},
				},
				expectedResults: []map[string]string{
					{"1": filepath.Join("A", "B", "C")},
					{"2": filepath.Join("D")},
					{"3": filepath.Join("/srv", "shiny-server", "log.txt")},
					{"4": ""},
					{"5": filepath.Join("..", "E", "F")},
				},
				testSubId : "3",
			},
		}

		for testName, tc := range tests {
			t.Run(testhelper.MakeSubtestName(configUnitTest3, tc.testSubId, testName), func(t *testing.T) {
				actualResults := expandTildesRepos(tc.repos)
				assert.Equal(t, tc.expectedResults, actualResults)
			})
		}
	})

}

func TestAddRemovePackage(t *testing.T) {
	tests := []struct {
		fileName    string
		packageName string
	}{
		{
			fileName:    "simple",
			packageName: "packageTestName",
		},
		{
			fileName:    "simple-suggests",
			packageName: "packageTestName",
		},
		{
			fileName:    "mixed-source",
			packageName: "packageTestName",
		},
		{
			fileName:    "outdated-pkgs",
			packageName: "packageTestName",
		},
		{
			fileName:    "outdated-pkgs-no-update",
			packageName: "packageTestName",
		},
		{
			fileName:    "repo-order",
			packageName: "packageTestName",
		},
	}

	t.Run(testhelper.MakeTestName(configUnitTest4, "packages can be added and removed"), func(t *testing.T) {
		appFS := afero.NewOsFs()
		for ttIndex, tt := range tests {
			t.Run(testhelper.MakeSubtestName(configUnitTest4, fmt.Sprint(ttIndex), fmt.Sprintf("test filename: %s", tt.fileName)), func(t *testing.T) {
				t.Logf("Test filename: %s", tt.fileName)
				fileName := filepath.Join(getIntegrationTestFolder(t, tt.fileName), "pkgr.yml")
				viper.SetConfigFile(fileName)

				b, _ := afero.Exists(appFS, fileName)
				assert.Equal(t, true, b, fmt.Sprintf("yml file not found:%s", fileName))

				ymlStart, _ := afero.ReadFile(appFS, fileName)

				AddPackages([]string{tt.packageName})
				b, _ = afero.FileContainsBytes(appFS, fileName, []byte(tt.packageName))
				assert.Equal(t, true, b, fmt.Sprintf("Package not added:%s", fileName))

				remove(fileName, tt.packageName)
				b, _ = afero.FileContainsBytes(appFS, fileName, []byte(tt.packageName))
				assert.Equal(t, false, b, fmt.Sprintf("Package not removed:%s", fileName))

				ymlEnd, _ := afero.ReadFile(appFS, fileName)
				b = equal(ymlStart, ymlEnd, false)
				assert.Equal(t, true, b, fmt.Sprintf("Start and End yml files differ:%s", fileName))

				// put file back for Git
				fi, _ := os.Stat(fileName)
				err := afero.WriteFile(appFS, fileName, ymlStart, fi.Mode())
				assert.Equal(t, nil, err, fmt.Sprintf("Error writing file back to original state:%s", fileName))
			})
		}
	})


}

func TestAddPackageWithDuplicate(t *testing.T) {
	type test struct {
		testFolder   string
		packageToAdd string
	}

	tests := []test{
		{
			"simple-modify",
			"R6",
		},
	}

	fs := afero.NewOsFs()

	t.Run(testhelper.MakeTestName(configUnitTest5, "add package with duplicates"), func(t *testing.T) {
		for _, testCase := range tests {
			pkgrYamlContent := []byte(`
Version: 1
# top level packages
Packages:
  - R6
  - pillar

# any repositories, order matters
Repos:
  - MPN: "https://mpn.metworx.com/snapshots/stable/2023-06-29"


Library: "test-library"

`)
			configFilePath := filepath.Join("testsite", testCase.testFolder, "pkgr.yml")
			_ = fs.Remove(configFilePath)
			err := afero.WriteFile(fs, configFilePath, pkgrYamlContent, 0755)
			if err != nil {
				t.Error("Could not write test pkgr.yml file in " + testCase.testFolder)
				t.Fail()
			}
			viper.SetConfigFile(configFilePath)
			resultErr := AddPackages([]string{testCase.packageToAdd})
			assert.Nil(t, resultErr, "failed to add package")

			var actualResult PkgrConfig
			postChangeConfig, err := afero.ReadFile(fs, configFilePath)
			assert.Nil(t, err, "Could not read in updated yml file for folder "+testCase.testFolder)
			err = yaml.Unmarshal(postChangeConfig, &actualResult)
			assert.Nil(t, err, "Could not unmarshal updated yml file for folder "+testCase.testFolder)

			pkgCount := 0
			for _, p := range actualResult.Packages {
				if p == testCase.packageToAdd {
					pkgCount++
				}
			}
			assert.Equal(t, 1, pkgCount, fmt.Sprintf("expected to find exactly one occurence of package %s in %s, found %d", testCase.packageToAdd, configFilePath, pkgCount))
		}
	})



}

func TestAddPackageLockfileConfig(t *testing.T) {
	type test struct {
		testFolder   string
		lockfileType string
		packageToAdd string
	}

	tests := []test{
		{
			"renv-library-modify",
			"renv",
			"renv",
		},
		{
			"packrat-library-modify",
			"packrat",
			"packrat",
		},
	}

	fs := afero.NewOsFs()

	t.Run(testhelper.MakeTestName(configUnitTest6, "can add packages in config with lockfile"), func(t *testing.T) {
		for _, testCase := range tests {
			pkgrYamlContent := []byte(fmt.Sprintf(`
Version: 1

Packages:
  - fansi
Repos:
  - MPN: "https://mpn.metworx.com/snapshots/stable/2019-12-02"

Lockfile:
  Type: %s
`, testCase.lockfileType))
			configFilePath := filepath.Join("testsite", testCase.testFolder, "pkgr.yml")
			_ = fs.Remove(configFilePath)
			err := afero.WriteFile(fs, configFilePath, pkgrYamlContent, 0755)
			if err != nil {
				t.Error("Could not write test pkgr.yml file in " + testCase.testFolder)
				t.Fail()
			}
			viper.SetConfigFile(configFilePath)
			resultErr := AddPackages([]string{testCase.packageToAdd})
			assert.Nil(t, resultErr, "failed to add package")

			// Find packageToRemove in yml file under Packages:
			var actualResult PkgrConfig
			postChangeConfig, err := afero.ReadFile(fs, configFilePath)
			assert.Nil(t, err, "Could not read in updated yml file for folder "+testCase.testFolder)
			err = yaml.Unmarshal(postChangeConfig, &actualResult)
			assert.Nil(t, err, "Could not unmarshal updated yml file for folder "+testCase.testFolder)

			pkgWasAdded := funk.Contains(actualResult.Packages, testCase.packageToAdd)
			assert.True(t, pkgWasAdded, fmt.Sprintf("package not found after add command: %s", testCase.packageToAdd))
		}
	})



}

func TestRemovePackageLockfileConfig(t *testing.T) {
	type test struct {
		testFolder      string
		lockfileType    string
		packageToRemove string
	}

	tests := []test{
		{
			"renv-library-modify",
			"renv",
			"renv",
		},
		{
			"packrat-library-modify",
			"packrat",
			"packrat",
		},
	}

	fs := afero.NewOsFs()

	t.Run(testhelper.MakeTestName(configUnitTest7, "can remove packages in config with lockfile"), func(t *testing.T) {
		for _, testCase := range tests {
			pkgrYamlContent := []byte(fmt.Sprintf(`
Version: 1

Packages:
  - %s
  - fansi
Repos:
  - MPN: "https://mpn.metworx.com/snapshots/stable/2019-12-02"

Lockfile:
  Type: %s
`,
				testCase.packageToRemove, testCase.lockfileType))

			configFilePath := filepath.Join("testsite", testCase.testFolder, "pkgr.yml")
			_ = fs.Remove(configFilePath)
			err := afero.WriteFile(fs, configFilePath, pkgrYamlContent, 0755)
			if err != nil {
				t.Error("Could not write test pkgr.yml file in " + testCase.testFolder)
				t.Fail()
			}
			resultErr := remove(configFilePath, testCase.packageToRemove)
			assert.Nil(t, resultErr, "failed to add package")

			// Verify packageToRemove is not in yml file.
			var actualResult PkgrConfig
			postChangeConfig, err := afero.ReadFile(fs, configFilePath)
			assert.Nil(t, err, "Could not read in updated yml file for folder "+testCase.testFolder)
			err = yaml.Unmarshal(postChangeConfig, &actualResult)
			assert.Nil(t, err, "Could not unmarshal updated yml file for folder "+testCase.testFolder)

			pkgWasRemoved := !funk.Contains(actualResult.Packages, testCase.packageToRemove)

			assert.True(t, pkgWasRemoved, fmt.Sprintf("package not excluded after add command: %s", testCase.packageToRemove))
		}
	})

}

func TestRemoveWhitespace(t *testing.T) {

	tests := []struct {
		in       string
		expected string
		message  string
	}{
		{
			in:       "hello world\n",
			expected: "helloworld",
			message:  "newline",
		},
		{
			in:       "hello world\t",
			expected: "helloworld",
			message:  "h tab",
		},
		{
			in:       "hello world\v",
			expected: "helloworld",
			message:  "v tab",
		},
		{
			in:       "hello world\f",
			expected: "helloworld",
			message:  "feed",
		},
		{
			in:       "hello world\r",
			expected: "helloworld",
			message:  "return",
		},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.expected, string(removeWhitespace([]byte(tt.in))), fmt.Sprintf("fail to remove:%s", tt.message))
	}
}

func removeWhitespace(b []byte) []byte {
	var ws = []byte{'\t', '\n', '\v', '\f', '\r', ' '}
	for _, r := range ws {
		b = bytes.ReplaceAll(b, []byte(string(r)), []byte(""))
	}
	return b
}

func equal(a []byte, b []byte, compareWs bool) bool {
	if compareWs == false {
		a = removeWhitespace(a)
		b = removeWhitespace(b)
	}
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

func TestNewConfigPackrat(t *testing.T) {
	var cfg PkgrConfig
	_ = os.Chdir(getTestFolder(t, "packrat-library"))
	NewConfig(viper.GetString("config"), &cfg)
	assert.Equal(t, "packrat", cfg.Lockfile.Type, fmt.Sprintf("Fail:%s", "packrat exists"))
}

func TestNewConfigRenv(t *testing.T) {
	// The "no renv" error case is covered by the integration
	// tests when PKGR_TESTS_SYS_RENV is an empty string.
	renv := os.Getenv("PKGR_TESTS_SYS_RENV")
	if renv == "" {
		t.Skip("Skipping: empty PKGR_TESTS_SYS_RENV indicates renv not available")
	}

	var cfg PkgrConfig
	_ = os.Chdir(getTestFolder(t, "renv-library"))
	NewConfig(viper.GetString("config"), &cfg)
	assert.Equal(t, "renv", cfg.Lockfile.Type, fmt.Sprintf("Fail:%s", "renv exists"))
}

func TestNewConfigNoLockfile(t *testing.T) {
	tests := []struct {
		folder   string
		expected string
		message  string
	}{
		{
			folder:   "simple",
			expected: "",
			message:  "packrat does not exist",
		},
	}
	for _, tt := range tests {
		var cfg PkgrConfig
		_ = os.Chdir(getTestFolder(t, tt.folder))
		NewConfig(viper.GetString("config"), &cfg)
		assert.Equal(t, tt.expected, cfg.Lockfile.Type, fmt.Sprintf("Fail:%s", tt.message))
	}
}

func TestGetLibraryPath(t *testing.T) {
	tests := []struct {
		lftype   string
		expected string
	}{
		{
			lftype:   "packrat",
			expected: "packrat/lib/apple/1.2.3",
		},
		{
			lftype:   "pkgr",
			expected: "original",
		},
	}
	for _, tt := range tests {
		var rv = cran.RVersion{
			Major: 1,
			Minor: 2,
			Patch: 3,
		}
		library := getLibraryPath(tt.lftype, "myRpath", rv, "apple", "original")
		assert.Equal(t, tt.expected, library, fmt.Sprintf("Fail:%s", tt.expected))
	}
}

func TestGetLibraryPathRenv(t *testing.T) {
	// The "no renv" error case is covered by the integration
	// tests when PKGR_TESTS_SYS_RENV is an empty string.
	renv := os.Getenv("PKGR_TESTS_SYS_RENV")
	if renv == "" {
		t.Skip("Skipping: empty PKGR_TESTS_SYS_RENV indicates renv not available")
	}

	var rs rcmd.RSettings
	rv := rcmd.GetRVersion(&rs)

	library := getLibraryPath("renv", "R", rv,
		"platform (ignored)", "original (ignored)")

	assert.Contains(t, library,
		fmt.Sprintf("renv/library/R-%d.%d/", rv.Major, rv.Minor))
}

func TestSetCustomizations(t *testing.T) {
	tests := []struct {
		pkg   string
		name  string
		value string
	}{
		{
			pkg:   "data.table",
			name:  "R_MAKEVARS_USER",
			value: "~/.R/Makevars_data.table",
		},
		{
			pkg:   "boo",
			name:  "foo",
			value: "soo",
		},
	}
	for _, tt := range tests {
		var cfg PkgrConfig
		_ = os.Chdir(getTestFolder(t, "simple"))
		NewConfig(viper.GetString("config"), &cfg) // Just need to slurp in misc stuff to keep tests working.
		cfg.Packages = []string{tt.pkg}            // Overwrites whatever packages were in "simple"
		cfg.Customizations.Packages = []map[string]PkgConfig{
			{
				tt.pkg: PkgConfig{
					Env: map[string]string{
						tt.name: tt.value,
					},
				},
			},
		}
		var rs rcmd.RSettings
		rs.PkgEnvVars = make(map[string]map[string]string)
		rs2 := SetCustomizations(rs, cfg)
		assert.Equal(t, tt.value, rs2.PkgEnvVars[tt.pkg][tt.name], fmt.Sprintf("Fail to get: %s", tt.value))
	}
}

func TestSetCfgCustomizations(t *testing.T) {
	tests := []struct {
		pkg string
	}{
		{
			pkg: "data.table",
		},
		{
			pkg: "boo",
		},
	}
	for _, tt := range tests {
		dependencyConfigurations := gpsr.NewDefaultInstallDeps()
		var cfg PkgrConfig
		//NewConfig(viper.GetString("config"), &cfg) // Don't need to load full config for this test
		cfg.Suggests = true
		cfg.Packages = []string{ // Overwrites
			tt.pkg,
		}
		setCfgCustomizations(cfg, &dependencyConfigurations)
		_, found := dependencyConfigurations.Deps[tt.pkg]
		assert.Equal(t, true, found, fmt.Sprintf("Fail to get: %s", tt.pkg))
		assert.Equal(t, cfg.Suggests, dependencyConfigurations.Deps[tt.pkg].Suggests, fmt.Sprintf("Suggest not correct: %s", tt.pkg))
	}
}
func TestSetViperCustomizations(t *testing.T) {
	tests := []struct {
		pkg      string
		repo     string
		suggests bool
		stype    string
		source   cran.SourceType
	}{
		{
			pkg:      "pillar",
			repo:     "CRAN2",
			suggests: false,
			stype:    "source",
			source:   cran.Source,
		},
		{
			pkg:      "R6",
			repo:     "CRAN",
			suggests: true,
			stype:    "binary",
			source:   cran.Binary,
		},
	}
	var installConfig = cran.InstallConfig{
		Packages: map[string]cran.PkgConfig{},
	}

	for _, tt := range tests {
		var cfg PkgrConfig
		//NewConfig(viper.GetString("config"), &cfg) // Not needed for test to run.
		cfg.Customizations.Packages = []map[string]PkgConfig{ {
			tt.pkg: PkgConfig{
				Env: map[string]string{
					"": "",
				},
				Repo:     tt.repo,
				Type:     tt.stype,
				Suggests: tt.suggests,
			},
		},
		}
		var urls = []cran.RepoURL{
			cran.RepoURL{
				Name: tt.repo,
				URL:  "https://mpn.metworx.com/snapshots/stable/2019-12-02",
			},
		}
		pkgNexus, _ := cran.NewPkgDb(urls, cran.Source, &installConfig, cran.RVersion{}, false)
		dependencyConfigurations := gpsr.NewDefaultInstallDeps()

		// build pkgSettings, normally read from the yml file, via viper
		var pkgSettings = []interface{}{
			map[interface{}]interface{}{
				tt.pkg: map[interface{}]interface{}{
					"Suggests": tt.suggests,
					"Repo":     tt.repo,
					"Type":     tt.stype,
				},
			},
		}

		// call the function to test
		setViperCustomizations(cfg, pkgSettings, dependencyConfigurations, pkgNexus)

		val := getCustomizationValue("Type", pkgSettings, tt.pkg)
		assert.Equal(t, tt.stype, val, fmt.Sprintf("Error setting type %s for pkg %s", tt.stype, tt.pkg))
		assert.Equal(t, tt.source, pkgNexus.Config.Packages[tt.pkg].Type, fmt.Sprintf("Error setting type %s for pkg %s", tt.stype, tt.pkg))

		// Repo
		val = getCustomizationValue("Repo", pkgSettings, tt.pkg)
		assert.Equal(t, tt.repo, val, fmt.Sprintf("Error setting repo %s for pkg %s", tt.repo, tt.pkg))
		assert.Equal(t, tt.repo, pkgNexus.Config.Packages[tt.pkg].Repo.Name, fmt.Sprintf("Error setting repo %s for pkg %s", tt.repo, tt.pkg))

		// Suggests
		pkgDepTypes, found := dependencyConfigurations.Deps[tt.pkg]
		assert.Equal(t, true, found, fmt.Sprintf("Deps not found:%s", tt.pkg))
		assert.Equal(t, tt.suggests, pkgDepTypes.Suggests, fmt.Sprintf("Suggests error for pkg %s", tt.pkg))
		val = getCustomizationValue("Suggests", pkgSettings, tt.pkg)
		assert.Equal(t, tt.suggests, val, fmt.Sprintf("Suggests error for pkg %s", tt.pkg))
	}
}

func getCustomizationValue(key string, elems []interface{}, elem string) interface{} {
	for _, v := range elems {
		for k, iv := range v.(map[interface{}]interface{}) {
			if k == elem {
				for val, k2 := range iv.(map[interface{}]interface{}) {
					if val == key {
						return k2
					}
				}
			}
		}
	}
	return ""
}

// compare pkgSettings map to viper interface
// pkgNexus should be the same
func TestSetViperCustomizations2(t *testing.T) {
	tests := []struct {
		pkg      string
		repo     string
		suggests bool
		stype    string
		source   cran.SourceType
	}{
		{
			pkg:      "pillar",
			repo:     "CRAN2",
			suggests: false,
			stype:    "source",
			source:   cran.Source,
		},
		{
			pkg:      "R6",
			repo:     "CRAN",
			suggests: true,
			stype:    "binary",
			source:   cran.Binary,
		},
	}
	var installConfig = cran.InstallConfig{
		Packages: map[string]cran.PkgConfig{},
	}

	for _, tt := range tests {
		var cfg PkgrConfig
		//NewConfig(viper.GetString("config"), &cfg) // Not needed for test to run
		cfg.Customizations.Packages = []map[string]PkgConfig{
			{
				tt.pkg: PkgConfig{
					Env: map[string]string{
						"": "",
					},
					Repo:     tt.repo,
					Type:     tt.stype,
					Suggests: tt.suggests,
				},
			},
		}
		var urls = []cran.RepoURL{
			cran.RepoURL{
				Name: tt.repo,
				URL:  "https://mpn.metworx.com/snapshots/stable/2019-12-02",
			},
		}

		pkgNexus, _ := cran.NewPkgDb(urls, cran.Source, &installConfig, cran.RVersion{}, false)
		pkgNexus2, _ := cran.NewPkgDb(urls, cran.Source, &installConfig, cran.RVersion{}, false)
		dependencyConfigurations := gpsr.NewDefaultInstallDeps()
		dependencyConfigurations2 := gpsr.NewDefaultInstallDeps()

		viper.Reset()
		viper.SetConfigName("pkgr")
		ymlfile := getIntegrationTestFolder(t, "customization")
		viper.AddConfigPath(ymlfile)

		err := viper.ReadInConfig()
		assert.Equal(t, nil, err, "Error reading yml file")
		pkgSettings := viper.Sub("Customizations").AllSettings()["packages"].([]interface{})
		setViperCustomizations(cfg, pkgSettings, dependencyConfigurations, pkgNexus)

		pkgSettings2 := PkgSettingsMap{
			"pillar": PkgConfig{
				Suggests: false,
				Repo:     "CRAN2",
				Type:     "source",
			},
			"R6": PkgConfig{
				Suggests: true,
				Repo:     "CRAN",
				Type:     "binary",
			},
		}
		setViperCustomizations2(cfg, pkgSettings2, dependencyConfigurations2, pkgNexus2)

		// Type
		assert.Equal(t, pkgNexus2.Config.Packages[tt.pkg].Type, pkgNexus.Config.Packages[tt.pkg].Type, fmt.Sprintf("Error setting type %s for pkg %s", tt.stype, tt.pkg))

		// Repo
		assert.Equal(t, pkgNexus2.Config.Packages[tt.pkg].Repo.Name, pkgNexus.Config.Packages[tt.pkg].Repo.Name, fmt.Sprintf("Error setting repo %s for pkg %s", tt.repo, tt.pkg))

		pkgDepTypes, found := dependencyConfigurations.Deps[tt.pkg]
		assert.Equal(t, true, found, fmt.Sprintf("Deps not found:%s", tt.pkg))
		pkgDepTypes2, found2 := dependencyConfigurations.Deps[tt.pkg]
		assert.Equal(t, true, found2, fmt.Sprintf("Deps2 not found:%s", tt.pkg))
		assert.Equal(t, pkgDepTypes2.Suggests, pkgDepTypes.Suggests, fmt.Sprintf("Suggests error for pkg %s", tt.pkg))

		val := getCustomizationValue("Suggests", pkgSettings, tt.pkg)
		val2 := pkgSettings2[tt.pkg].Suggests
		assert.Equal(t, val2, val, fmt.Sprintf("Suggests error for pkg %s", tt.pkg))
	}
}

func TestSetPkgConfig(t *testing.T) {
	tests := []struct {
		pkg      string
		repo     string
		suggests bool
		stype    string
		source   cran.SourceType
	}{
		{
			pkg:      "pillar",
			repo:     "CRAN2",
			suggests: false,
			stype:    "source",
			source:   cran.Source,
		},
		{
			pkg:      "R6",
			repo:     "CRAN",
			suggests: true,
			stype:    "binary",
			source:   cran.Binary,
		},
	}
	var installConfig = cran.InstallConfig{
		Packages: map[string]cran.PkgConfig{},
	}

	for _, tt := range tests {
		var cfg PkgrConfig
		//NewConfig(viper.GetString("config"), &cfg) // not needed for test to run
		cfg.Customizations.Packages = []map[string]PkgConfig{
			{
				tt.pkg: PkgConfig{
					Env: map[string]string{
						"": "",
					},
					Repo:     tt.repo,
					Type:     tt.stype,
					Suggests: tt.suggests,
				},
			},
		}
		var urls = []cran.RepoURL{
			cran.RepoURL{
				Name: tt.repo,
				URL:  "https://mpn.metworx.com/snapshots/stable/2019-12-02",
			},
		}

		pkgNexus, _ := cran.NewPkgDb(urls, cran.Source, &installConfig, cran.RVersion{}, false)
		pkgNexus2, _ := cran.NewPkgDb(urls, cran.Source, &installConfig, cran.RVersion{}, false)
		dependencyConfigurations := gpsr.NewDefaultInstallDeps()
		dependencyConfigurations2 := gpsr.NewDefaultInstallDeps()

		viper.Reset()
		viper.SetConfigName("pkgr")

		ymlfile := getIntegrationTestFolder(t, "customization")

		viper.AddConfigPath(ymlfile)
		err := viper.ReadInConfig()
		assert.Equal(t, nil, err, "Error reading yml file")
		pkgSettings := viper.Sub("Customizations").AllSettings()["packages"].([]interface{})
		setViperCustomizations(cfg, pkgSettings, dependencyConfigurations, pkgNexus)

		// get pksettings as map
		pkgSettings2 := getPkgSettingsMap(pkgSettings)

		setViperCustomizations2(cfg, pkgSettings2, dependencyConfigurations2, pkgNexus2)

		// Type
		assert.Equal(t, pkgNexus2.Config.Packages[tt.pkg].Type, pkgNexus.Config.Packages[tt.pkg].Type, fmt.Sprintf("Error setting type %s for pkg %s", tt.stype, tt.pkg))

		// Repo
		assert.Equal(t, pkgNexus2.Config.Packages[tt.pkg].Repo.Name, pkgNexus.Config.Packages[tt.pkg].Repo.Name, fmt.Sprintf("Error setting repo %s for pkg %s", tt.repo, tt.pkg))

		pkgDepTypes, found := dependencyConfigurations.Deps[tt.pkg]
		assert.Equal(t, true, found, fmt.Sprintf("Deps not found:%s", tt.pkg))
		pkgDepTypes2, found2 := dependencyConfigurations.Deps[tt.pkg]
		assert.Equal(t, true, found2, fmt.Sprintf("Deps2 not found:%s", tt.pkg))
		assert.Equal(t, pkgDepTypes2.Suggests, pkgDepTypes.Suggests, fmt.Sprintf("Suggests error for pkg %s", tt.pkg))

		val := getCustomizationValue("Suggests", pkgSettings, tt.pkg)
		val2 := pkgSettings2[tt.pkg].Suggests
		assert.Equal(t, val2, val, fmt.Sprintf("Suggests error for pkg %s", tt.pkg))
	}
}

func TestNewConfigSimple(t *testing.T) {
	//type testCase struct {
	//	testSet string
	//
	//}
	testYamlFile := filepath.Join(getTestFolder(t, "simple"), "pkgr.yml")

	var cfg PkgrConfig
	NewConfig(testYamlFile, &cfg)

	//cfg.Library -
	//cfg.Packages -
	//cfg.Version
	//cfg.Cache
	//cfg.Tarballs
	//cfg.RPath
	//cfg.Logging
	//cfg.Update
	//cfg.Suggests
	//cfg.Customizations
	//cfg.Repos
	//cfg.Strict
	//cfg.Lockfile
	//cfg.Rollback
	//cfg.Threads
	//cfg.NoRecommended

	assert.Contains(t, cfg.Packages, "R6")
	assert.Contains(t, cfg.Packages, "pillar")
	assert.Equal(t, cfg.Library, "test-library", cfg.Library)
	assert.Equal(t, 1, cfg.Version)
	assert.Equal(t, "", cfg.Cache)
	assert.Empty(t, cfg.Tarballs)
	assert.Equal(t, "R", filepath.Base(cfg.RPath)) // Just make sure the path ends in R executable. May not work on Windows.
	assert.Equal(t, LogConfig{}, cfg.Logging)
	assert.Equal(t, false, cfg.NoUpdate)
	assert.Equal(t, false, cfg.Suggests)
	assert.Empty(t, cfg.Customizations)
	assert.Equal(t, []map[string]string{{"MPN": "https://mpn.metworx.com/snapshots/stable/2023-06-29"}}, cfg.Repos)
	assert.Equal(t, false, cfg.Strict)
	assert.Equal(t, Lockfile{}, cfg.Lockfile)
	assert.Equal(t, false, cfg.NoRollback)
	assert.True(t, reflect.TypeOf(cfg.Threads).String() == "int")
	assert.Equal(t, false, cfg.NoRecommended)

}

func TestNewConfigNonDefaults(t *testing.T) {
	renv := os.Getenv("PKGR_TESTS_SYS_RENV")
	if renv == "" {
		t.Skip("Skipping: empty PKGR_TESTS_SYS_RENV indicates renv not available")
	}
	//type testCase struct {
	//	testSet string
	//
	//}
	testYamlFile := filepath.Join(getTestFolder(t, "many-settings"), "pkgr.yml")

	var cfg PkgrConfig
	NewConfig(testYamlFile, &cfg)

	assert.Contains(t, cfg.Packages, "R6")
	assert.Contains(t, cfg.Packages, "pillar")

	assert.True(t, strings.Contains(cfg.Library, "renv/")) // Should be set because of Lockfile setting.

	assert.Equal(t, 1, cfg.Version)
	assert.Equal(t, "./localcache", cfg.Cache)
	assert.True(t, strings.Contains(cfg.Tarballs[0], "folder/tarball.tar.gz")) // Should be set somewhere in the homedir.
	// assert.Equal(t, "../R", cfg.RPath) // Disabling this to make the test easier.
	assert.Equal(t, LogConfig{
		All:       "log/log.txt",
		Install:   "log/install.txt",
		Level:     "debug",
		Overwrite: true,
	}, cfg.Logging)
	assert.Equal(t, true, cfg.NoUpdate)
	assert.Equal(t, true, cfg.Suggests)
	assert.Empty(t, cfg.Customizations) // Customizations are tested elsewhere.
	assert.Equal(t, []map[string]string{{"MPN": "https://mpn.metworx.com/snapshots/stable/2023-06-29"}}, cfg.Repos)
	assert.Equal(t, true, cfg.Strict)
	assert.Equal(t, Lockfile{
		Type: "renv",
	}, cfg.Lockfile)
	assert.Equal(t, true, cfg.NoRollback)
	assert.Equal(t, 9, cfg.Threads)
	assert.Equal(t, true, cfg.NoRecommended)

}

func getPkgSettingsMap(elems []interface{}) PkgSettingsMap {
	pkgSettingsMap := make(map[string]PkgConfig)
	for _, e := range elems {
		for pkg, iv := range e.(map[interface{}]interface{}) {
			var pkgConfig PkgConfig
			for name, value := range iv.(map[interface{}]interface{}) {
				switch name {
				case "Suggests":
					pkgConfig.Suggests = value.(bool)
				case "Type":
					pkgConfig.Type = value.(string)
				case "Repo":
					pkgConfig.Repo = value.(string)
					// TODO: Env map
				}
			}
			pkgSettingsMap[pkg.(string)] = pkgConfig
		}
	}
	return pkgSettingsMap
}

func setViperCustomizations2(cfg PkgrConfig, pkgSettings PkgSettingsMap, dependencyConfigurations gpsr.InstallDeps, pkgNexus *cran.PkgNexus) {
	for _, pkgArray := range cfg.Customizations.Packages {
		for pkg, v := range pkgArray {
			if pkgSettings[pkg].Suggests {
				pkgDepTypes := dependencyConfigurations.Default
				pkgDepTypes.Suggests = v.Suggests
				dependencyConfigurations.Deps[pkg] = pkgDepTypes
			}
			if len(pkgSettings[pkg].Repo) > 0 {
				err := pkgNexus.SetPackageRepo(pkg, v.Repo)
				if err != nil {
					log.WithFields(log.Fields{
						"pkg":  pkg,
						"repo": v.Repo,
					}).Fatal("error finding custom repo to set")
				}
			}
			if len(pkgSettings[pkg].Type) > 0 {
				err := pkgNexus.SetPackageType(pkg, v.Type)
				if err != nil {
					log.WithFields(log.Fields{
						"pkg":  pkg,
						"repo": v.Repo,
					}).Fatal("error finding custom repo to set")
				}
			}
		}
	}
}

func getTestFolder(t *testing.T, folder string) string {
	_, filename, _, _ := runtime.Caller(0)
	top := filepath.Join(filepath.Dir(filename), "..")

	testFolder := filepath.Join(top, "configlib", "testsite", folder)
	return testFolder
}

func getIntegrationTestFolder(t *testing.T, folder string) string {
	_, filename, _, _ := runtime.Caller(0)
	top := filepath.Join(filepath.Dir(filename), "..")
	return filepath.Join(top, "configlib", "testsite", "integration_test_archive", folder)
}
