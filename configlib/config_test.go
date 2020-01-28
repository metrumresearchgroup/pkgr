package configlib

import (
	"bytes"
	"fmt"
	"github.com/thoas/go-funk"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
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

	appFS := afero.NewOsFs()
	for _, tt := range tests {
		fileName := filepath.Join(getIntegrationTestFolder(t, tt.fileName), "pkgr.yml")

		b, _ := afero.Exists(appFS, fileName)
		assert.Equal(t, true, b, fmt.Sprintf("yml file not found:%s", fileName))

		ymlStart, _ := afero.ReadFile(appFS, fileName)

		add(fileName, tt.packageName)
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
	}
}

func TestAddPackageWithDuplicate(t *testing.T) {
	type test struct {
		testFolder string
		packageToAdd string
	}

	tests := []test {
		{
			"simple-modify",
			"R6",
		},
	}

	fs := afero.NewOsFs()

	for _, testCase := range tests {
		pkgrYamlContent := []byte(`
Version: 1
# top level packages
Packages:
  - R6
  - pillar

# any repositories, order matters
Repos:
  - CRAN: "https://cran.microsoft.com/snapshot/2019-05-01"


Library: "test-library"

`)
		configFilePath := filepath.Join("testsite", testCase.testFolder, "pkgr.yml")
		_ = fs.Remove(configFilePath)
		err := afero.WriteFile(fs, configFilePath, pkgrYamlContent,  0755)
		if err != nil {
			t.Error("Could not write test pkgr.yml file in " + testCase.testFolder)
			t.Fail()
		}
		resultErr := add(configFilePath, testCase.packageToAdd)
		assert.Nil(t, resultErr, "failed to add package")

		var actualResult PkgrConfig
		postChangeConfig, err := afero.ReadFile(fs, configFilePath)
		assert.Nil(t, err, "Could not read in updated yml file for folder " + testCase.testFolder)
		err = yaml.Unmarshal(postChangeConfig, &actualResult)
		assert.Nil(t, err, "Could not unmarshal updated yml file for folder " + testCase.testFolder)

		pkgCount := 0
		for _, p := range actualResult.Packages {
			if p == testCase.packageToAdd {
				pkgCount++
			}
		}
		assert.Equal(t, 1, pkgCount, fmt.Sprintf("expected to find exactly one occurence of package %s in %s, found %d", testCase.packageToAdd, configFilePath, pkgCount))
	}

}

func TestAddPackageLockfileConfig(t *testing.T) {
	type test struct {
		testFolder string
		lockfileType string
		packageToAdd string
	}

	tests := []test {
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
		err := afero.WriteFile(fs, configFilePath, pkgrYamlContent,  0755)
		if err != nil {
			t.Error("Could not write test pkgr.yml file in " + testCase.testFolder)
			t.Fail()
		}
		resultErr := add(configFilePath, testCase.packageToAdd)
		assert.Nil(t, resultErr, "failed to add package")

		// Find packageToRemove in yml file under Packages:
		var actualResult PkgrConfig
		postChangeConfig, err := afero.ReadFile(fs, configFilePath)
		assert.Nil(t, err, "Could not read in updated yml file for folder " + testCase.testFolder)
		err = yaml.Unmarshal(postChangeConfig, &actualResult)
		assert.Nil(t, err, "Could not unmarshal updated yml file for folder " + testCase.testFolder)

		pkgWasAdded := funk.Contains(actualResult.Packages, testCase.packageToAdd)
		assert.True(t, pkgWasAdded, fmt.Sprintf("package not found after add command: %s", testCase.packageToAdd))
	}

}

func TestRemovePackageLockfileConfig(t *testing.T) {
	type test struct {
		testFolder      string
		lockfileType    string
		packageToRemove string
	}

	tests := []test {
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
		err := afero.WriteFile(fs, configFilePath, pkgrYamlContent,  0755)
		if err != nil {
			t.Error("Could not write test pkgr.yml file in " + testCase.testFolder)
			t.Fail()
		}
		resultErr := remove(configFilePath, testCase.packageToRemove)
		assert.Nil(t, resultErr, "failed to add package")

		// Verify packageToRemove is not in yml file.
		var actualResult PkgrConfig
		postChangeConfig, err := afero.ReadFile(fs, configFilePath)
		assert.Nil(t, err, "Could not read in updated yml file for folder " + testCase.testFolder)
		err = yaml.Unmarshal(postChangeConfig, &actualResult)
		assert.Nil(t, err, "Could not unmarshal updated yml file for folder " + testCase.testFolder)

		pkgWasRemoved := !funk.Contains(actualResult.Packages, testCase.packageToRemove)

		assert.True(t, pkgWasRemoved, fmt.Sprintf("package not excluded after add command: %s", testCase.packageToRemove))
	}

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
	tests := []struct {
		folder   string
		expected string
		message  string
	}{
		{
			folder:   "packrat-library",
			expected: "packrat",
			message:  "packrat exists",
		},
		{
			folder: "renv-library",
			expected: "renv",
			message: "renv exists",
		},
	}
	for _, tt := range tests {
		var cfg PkgrConfig
		_ = os.Chdir(getTestFolder(t, tt.folder))
		_ = LoadConfigFromPath(viper.GetString("config"))
		NewConfig(&cfg)
		assert.Equal(t, tt.expected, cfg.Lockfile.Type, fmt.Sprintf("Fail:%s", tt.message))
	}
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
		_ = LoadConfigFromPath(viper.GetString("config"))
		NewConfig(&cfg)
		assert.Equal(t, tt.expected, cfg.Lockfile.Type, fmt.Sprintf("Fail:%s", tt.message))
	}
}

func TestGetLibraryPath(t *testing.T) {
	tests := []struct {
		lftype   string
		expected string
		message  string
	}{
		{
			lftype:   "renv",
			expected: "renv/library/R-1.2/apple",
		},
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
		NewConfig(&cfg)
		cfg.Customizations.Packages = map[string]PkgConfig{
			tt.pkg: PkgConfig{
				Env: map[string]string{
					tt.name: tt.value,
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
		NewConfig(&cfg)
		cfg.Suggests = true
		cfg.Packages = []string{
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
		NewConfig(&cfg)
		cfg.Customizations.Packages = map[string]PkgConfig{
			tt.pkg: PkgConfig{
				Env: map[string]string{
					"": "",
				},
				Repo:     tt.repo,
				Type:     tt.stype,
				Suggests: tt.suggests,
			},
		}
		var urls = []cran.RepoURL{
			cran.RepoURL{
				Name: tt.repo,
				URL:  "https://cran.microsoft.com/snapshot/2018-11-11",
			},
		}
		pkgNexus, _ := cran.NewPkgDb(urls, cran.Source, &installConfig, cran.RVersion{})
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
		assert.Equal(t, tt.source, pkgNexus.Config.Packages[tt.pkg].GetSourceType2(), fmt.Sprintf("Error setting type %s for pkg %s", tt.stype, tt.pkg))

		// Repo
		val = getCustomizationValue("Repo", pkgSettings, tt.pkg)
		assert.Equal(t, tt.repo, val, fmt.Sprintf("Error setting repo %s for pkg %s", tt.repo, tt.pkg))
		assert.Equal(t, tt.repo, pkgNexus.Config.Packages[tt.pkg].GetOrigin().Name, fmt.Sprintf("Error setting repo %s for pkg %s", tt.repo, tt.pkg))

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
		NewConfig(&cfg)
		cfg.Customizations.Packages = map[string]PkgConfig{
			tt.pkg: PkgConfig{
				Env: map[string]string{
					"": "",
				},
				Repo:     tt.repo,
				Type:     tt.stype,
				Suggests: tt.suggests,
			},
		}
		var urls = []cran.RepoURL{
			cran.RepoURL{
				Name: tt.repo,
				URL:  "https://cran.microsoft.com/snapshot/2018-11-11",
			},
		}

		pkgNexus, _ := cran.NewPkgDb(urls, cran.Source, &installConfig, cran.RVersion{})
		pkgNexus2, _ := cran.NewPkgDb(urls, cran.Source, &installConfig, cran.RVersion{})
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
		assert.Equal(t, pkgNexus2.Config.Packages[tt.pkg].GetSourceType2(), pkgNexus.Config.Packages[tt.pkg].GetSourceType2(), fmt.Sprintf("Error setting type %s for pkg %s", tt.stype, tt.pkg))

		// Repo
		assert.Equal(t, pkgNexus2.Config.Packages[tt.pkg].GetOrigin().Name, pkgNexus.Config.Packages[tt.pkg].GetOrigin().Name, fmt.Sprintf("Error setting repo %s for pkg %s", tt.repo, tt.pkg))

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
		NewConfig(&cfg)
		cfg.Customizations.Packages = map[string]PkgConfig{
			tt.pkg: PkgConfig{
				Env: map[string]string{
					"": "",
				},
				Repo:     tt.repo,
				Type:     tt.stype,
				Suggests: tt.suggests,
			},
		}
		var urls = []cran.RepoURL{
			cran.RepoURL{
				Name: tt.repo,
				URL:  "https://cran.microsoft.com/snapshot/2018-11-11",
			},
		}

		pkgNexus, _ := cran.NewPkgDb(urls, cran.Source, &installConfig, cran.RVersion{})
		pkgNexus2, _ := cran.NewPkgDb(urls, cran.Source, &installConfig, cran.RVersion{})
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
		assert.Equal(t, pkgNexus2.Config.Packages[tt.pkg].GetSourceType2(), pkgNexus.Config.Packages[tt.pkg].GetSourceType2(), fmt.Sprintf("Error setting type %s for pkg %s", tt.stype, tt.pkg))

		// Repo
		assert.Equal(t, pkgNexus2.Config.Packages[tt.pkg].GetOrigin().Name, pkgNexus.Config.Packages[tt.pkg].GetOrigin().Name, fmt.Sprintf("Error setting repo %s for pkg %s", tt.repo, tt.pkg))

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
	for pkg, v := range cfg.Customizations.Packages {
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

func getTestFolder(t *testing.T, folder string) string {
	_, filename, _, _ := runtime.Caller(0)
	sa := strings.SplitAfter(filename, "/pkgr/")

	testFolder := filepath.Join(filepath.Dir(sa[0]), "configlib", "testsite", folder)
	return testFolder
}

// NOTE: This should NOT be used, but I'm creating this function as a patch while we decide the best way to test these things.
// It is only acceptable to use this function for tests that access ONLY the pkgr.yml files in the integration tests in
// a read-only capacity.
func getIntegrationTestFolder(t *testing.T, folder string) string {
	_, filename, _, _ := runtime.Caller(0)
	sa := strings.SplitAfter(filename, "/pkgr/")
	return filepath.Join(filepath.Dir(sa[0]), "integration_tests", folder)
}
