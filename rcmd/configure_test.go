package rcmd

import (
	"fmt"
	"github.com/metrumresearchgroup/pkgr/testhelper"
	"io/ioutil"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type configureArgsTestCase struct {
	context string
	in      string
	// mocked system environment variables per os.Environ()
	sysEnv   []string
	expected []string
}

// Utility functions
func checkEnvVarsValid(t *testing.T, testCase configureArgsTestCase, actualResults []string) {
	rLibsUserFound := false
	for _, envVar := range actualResults {
		if strings.Contains(envVar, "R_LIBS_USER") {
			rLibsUserFound = true
			tmpDir := strings.Split(envVar, "=")[1]
			checkIsTempDir(t, tmpDir)
			assert.DirExists(t, tmpDir)
			dirEntries, err := ioutil.ReadDir(tmpDir)
			assert.Nil(t, err)
			assert.Empty(t, dirEntries, "failure: R_LIBS_USER was not set to an EMPTY temp directory")
		} else {
			assert.Contains(t, testCase.expected, envVar, "excess environment vars found")
			// assert.Equal(testCase.expected[index], envVar) // We are no longer claiming that order matters.
		}
	}
	assert.True(t, rLibsUserFound, "R_LIBS_USER was not set -- we expect it to always be set")
	// Make sure we're not missing any expected vars. A little redundant, but the easiest way to do this.
	for _, envVar := range testCase.expected {
		if strings.Contains(envVar, "R_LIBS_USER") {
			continue
		} else {
			assert.Contains(t, actualResults, envVar, "missing expected environment var")
		}
	}
}

func checkIsTempDir(t *testing.T, tmpDir string) {
	switch runtime.GOOS {
	case "darwin":
		assert.True(t, strings.Contains(tmpDir, "var/folders"), "R_LIBS_USER not set to temp directory: Dir found: %s", tmpDir)
		break
	case "linux":
		t.Skip("tmp dir check not implemented for linux")
		break
	case "windows":
		t.Skip("tmp dir check not implemented for linux")
		break
	default:
		t.Skip("tmp dir check not implemented for detected os")
	}
}

// end Utility functions

//These tests are going to be important for validation and must be named.
const (
	configureEnvVarsTest = "ENV-UNIT-001"
)

func TestConfigureArgs(t *testing.T) {
	t.Run(testhelper.MakeTestName(configureEnvVarsTest, "ensure environment variables are set correctly in R session"), func(t *testing.T) {
		defaultRS := NewRSettings("")
		// there should always be at least one libpath
		defaultRS.LibPaths = []string{"path/to/install/lib"}
		defaultRS.PkgEnvVars["dplyr"] = map[string]string{"DPLYR_ENV": "true"}
		var installArgsTests = []configureArgsTestCase{
			{
				"minimal",
				"",
				[]string{},
				[]string{"R_LIBS_SITE=path/to/install/lib", "R_LIBS_USER=SHOULD_BE_TMP_DIR"},
			},
			{
				"non-impactful system env set",
				"",
				[]string{"MISC_ENV=foo", "MISC2=bar"},
				[]string{"MISC_ENV=foo", "MISC2=bar", "R_LIBS_SITE=path/to/install/lib"},
			},
			{
				"non-impactful system env set with known package",
				"dplyr",
				[]string{"MISC_ENV=foo", "MISC2=bar"},
				[]string{"DPLYR_ENV=true", "MISC_ENV=foo", "MISC2=bar", "R_LIBS_SITE=path/to/install/lib", "R_LIBS_USER=SHOULD_BE_TMP_DIR"},
			},
			{
				"impactful system env set on separate package",
				"",
				[]string{"MISC_ENV=foo", "MISC2=bar", "DPLYR_ENV=false"},
				[]string{"MISC_ENV=foo", "MISC2=bar", "DPLYR_ENV=false", "R_LIBS_SITE=path/to/install/lib", "R_LIBS_USER=SHOULD_BE_TMP_DIR"},
			},
			{
				"impactful system env set with known package",
				"dplyr",
				[]string{"MISC_ENV=foo", "MISC2=bar", "DPLYR_ENV=false"},
				[]string{"DPLYR_ENV=true", "MISC_ENV=foo", "MISC2=bar", "R_LIBS_SITE=path/to/install/lib", "R_LIBS_USER=SHOULD_BE_TMP_DIR"},
			},
			{
				"R_LIBS_SITE env set",
				"",
				[]string{"R_LIBS_SITE=original/path", "MISC2=bar"},
				[]string{"MISC2=bar", "R_LIBS_SITE=path/to/install/lib", "R_LIBS_USER=SHOULD_BE_TMP_DIR"},
			},
			{
				"R_LIBS_SITE env set with known package",
				"dplyr",
				[]string{"R_LIBS_SITE=original/path", "MISC2=bar"},
				[]string{"DPLYR_ENV=true", "MISC2=bar", "R_LIBS_SITE=path/to/install/lib", "R_LIBS_USER=SHOULD_BE_TMP_DIR"},
			},
			{
				"R_LIBS_USER env set",
				"",
				[]string{"R_LIBS_USER=original/path", "MISC2=bar"},
				[]string{"MISC2=bar", "R_LIBS_SITE=path/to/install/lib", "R_LIBS_USER=SHOULD_BE_TMP_DIR"},
			},
			{
				"R_LIBS_USER env set with known package",
				"dplyr",
				[]string{"R_LIBS_USER=original/path", "MISC2=bar"},
				[]string{"DPLYR_ENV=true", "MISC2=bar", "R_LIBS_SITE=path/to/install/lib", "R_LIBS_USER=SHOULD_BE_TMP_DIR"},
			},
			{
				"R_LIBS_SITE and R_LIBS_USER env set",
				"",
				[]string{"R_LIBS_USER=original/path", "R_LIBS_SITE=original/site/path", "MISC2=bar"},
				[]string{"MISC2=bar", "R_LIBS_SITE=path/to/install/lib", "R_LIBS_USER=SHOULD_BE_TMP_DIR"},
			},
			{
				"R_LIBS_SITE and R_LIBS_USER env set",
				"dplyr",
				[]string{"R_LIBS_USER=original/path", "R_LIBS_SITE=original/site/path", "MISC2=bar"},
				[]string{"DPLYR_ENV=true", "MISC2=bar", "R_LIBS_SITE=path/to/install/lib", "R_LIBS_USER=SHOULD_BE_TMP_DIR"},
			},
			{
				"System contains sensitive information",
				"",
				[]string{"R_LIBS_USER=original/path", "GITHUB_PAT=should_get_hidden1", "ghe_token=should_get_hidden2", "ghe_PAT=should_get_hidden3", "github_token=should_get_hidden4"},
				[]string{"GITHUB_PAT=**HIDDEN**", "ghe_token=**HIDDEN**", "ghe_PAT=**HIDDEN**", "github_token=**HIDDEN**", "R_LIBS_SITE=path/to/install/lib", "R_LIBS_USER=SHOULD_BE_TMP_DIR"},
			},
		}
		for testNum, tt := range installArgsTests {
			t.Run(testhelper.MakeSubtestName(configureEnvVarsTest, fmt.Sprint(testNum), tt.context), func(t *testing.T) {
				actual := configureEnv(tt.sysEnv, defaultRS, tt.in)

				// Make sure that all environment variables are present
				// Also make sure that R_LIBS_USER is set.
				checkEnvVarsValid(t, tt, actual)

				//assert.Equal(tt.expected, actual, fmt.Sprintf("%s, test num: %v", tt.context, i+1))
			})
		}
	})
}

// Since the other tests are failing and we haven't addressed them yet, I'm just breaking this one out into its own group.
func TestConfigureArgs2(t *testing.T) {
	defaultRS := NewRSettings("")
	// there should always be at least one libpath
	defaultRS.LibPaths = []string{"path/to/install/lib"}
	defaultRS.PkgEnvVars["dplyr"] = map[string]string{"DPLYR_ENV": "true"}
	var installArgsTests = []configureArgsTestCase{
		{
			"System contains sensitive information",
			"",
			[]string{
				"R_LIBS_USER=original/path",
				"GITHUB_PAT=should_get_hidden1",
				"ghe_token=should_get_hidden2",
				"ghe_PAT=should_get_hidden3",
				"github_token=should_get_hidden4",
				"AWS_ACCESS_KEY_ID=should_get_hidden5",
				"AWS_SECRET_KEY=should_get_hidden6",
			},
			[]string{
				"GITHUB_PAT=**HIDDEN**",
				"ghe_token=**HIDDEN**",
				"ghe_PAT=**HIDDEN**",
				"github_token=**HIDDEN**",
				"AWS_ACCESS_KEY_ID=**HIDDEN**",
				"AWS_SECRET_KEY=**HIDDEN**",
				"R_LIBS_SITE=path/to/install/lib",
				"R_LIBS_USER=SHOULD_BE_TMP_DIR",
			},
		},
	}
	for _, tt := range installArgsTests {
		actual := configureEnv(tt.sysEnv, defaultRS, tt.in)
		//assert.Equal(tt.expected, actual, fmt.Sprintf("%s, test num: %v", tt.context, i+1))

		checkEnvVarsValid(t, tt, actual)
	}
}

func TestCensoredEnvVars(t *testing.T) {
	tests := map[string]struct {
		additionalVars []string
		expected       map[string]string
	}{
		"Default": {
			additionalVars: nil,
			expected: map[string]string{
				"GITHUB_TOKEN":      "GITHUB_TOKEN",
				"GITHUB_PAT":        "GITHUB_PAT",
				"GHE_TOKEN":         "GHE_TOKEN",
				"GHE_PAT":           "GHE_PAT",
				"AWS_ACCESS_KEY_ID": "AWS_ACCESS_KEY_ID",
				"AWS_SECRET_KEY":    "AWS_SECRET_KEY",
			},
		},
		"Empty arg": {
			additionalVars: []string{},
			expected: map[string]string{
				"GITHUB_TOKEN":      "GITHUB_TOKEN",
				"GITHUB_PAT":        "GITHUB_PAT",
				"GHE_TOKEN":         "GHE_TOKEN",
				"GHE_PAT":           "GHE_PAT",
				"AWS_ACCESS_KEY_ID": "AWS_ACCESS_KEY_ID",
				"AWS_SECRET_KEY":    "AWS_SECRET_KEY",
			},
		},
		"Add one": {
			additionalVars: []string{"cats"},
			expected: map[string]string{
				"GITHUB_TOKEN":      "GITHUB_TOKEN",
				"GITHUB_PAT":        "GITHUB_PAT",
				"GHE_TOKEN":         "GHE_TOKEN",
				"GHE_PAT":           "GHE_PAT",
				"AWS_ACCESS_KEY_ID": "AWS_ACCESS_KEY_ID",
				"AWS_SECRET_KEY":    "AWS_SECRET_KEY",
				"CATS":              "CATS",
			},
		},
		"Add two": {
			additionalVars: []string{"CATS", "and_oranges"},
			expected: map[string]string{
				"GITHUB_TOKEN":      "GITHUB_TOKEN",
				"GITHUB_PAT":        "GITHUB_PAT",
				"GHE_TOKEN":         "GHE_TOKEN",
				"GHE_PAT":           "GHE_PAT",
				"AWS_ACCESS_KEY_ID": "AWS_ACCESS_KEY_ID",
				"AWS_SECRET_KEY":    "AWS_SECRET_KEY",
				"CATS":              "CATS",
				"AND_ORANGES":       "AND_ORANGES",
			},
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			results := censoredEnvVars(test.additionalVars)
			assert.Equal(t, test.expected, results, fmt.Sprintf("failure in test %s: expected not equal to actual", testName))
		})

	}
}
