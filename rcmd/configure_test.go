package rcmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigureArgs(t *testing.T) {
	assert := assert.New(t)

	defaultRS := NewRSettings("")
	// there should always be at least one libpath
	defaultRS.LibPaths = []string{"path/to/install/lib"}
	defaultRS.PkgEnvVars["dplyr"] = map[string]string{"DPLYR_ENV": "true"}
	var installArgsTests = []struct {
		context string
		in      string
		// mocked system environment variables per os.Environ()
		sysEnv   []string
		expected []string
	}{
		{
			"minimal",
			"",
			[]string{},
			[]string{"R_LIBS_SITE=path/to/install/lib"},
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
			[]string{"DPLYR_ENV=true", "MISC_ENV=foo", "MISC2=bar", "R_LIBS_SITE=path/to/install/lib"},
		},
		{
			"impactful system env set on separate package",
			"",
			[]string{"MISC_ENV=foo", "MISC2=bar", "DPLYR_ENV=false"},
			[]string{"MISC_ENV=foo", "MISC2=bar", "DPLYR_ENV=false", "R_LIBS_SITE=path/to/install/lib"},
		},
		{
			"impactful system env set with known package",
			"dplyr",
			[]string{"MISC_ENV=foo", "MISC2=bar", "DPLYR_ENV=false"},
			[]string{"DPLYR_ENV=true", "MISC_ENV=foo", "MISC2=bar", "R_LIBS_SITE=path/to/install/lib"},
		},
		{
			"R_LIBS_SITE env set",
			"",
			[]string{"R_LIBS_SITE=original/path", "MISC2=bar"},
			[]string{"MISC2=bar", "R_LIBS_SITE=path/to/install/lib",},
		},
		{
			"R_LIBS_SITE env set with known package",
			"dplyr",
			[]string{"R_LIBS_SITE=original/path", "MISC2=bar"},
			[]string{"DPLYR_ENV=true", "MISC2=bar", "R_LIBS_SITE=path/to/install/lib"},
		},
		{
			"R_LIBS_USER env set",
			"",
			[]string{"R_LIBS_USER=original/path", "MISC2=bar"},
			[]string{"MISC2=bar", "R_LIBS_SITE=path/to/install/lib"},
		},
		{
			"R_LIBS_USER env set with known package",
			"dplyr",
			[]string{"R_LIBS_USER=original/path", "MISC2=bar"},
			[]string{"DPLYR_ENV=true", "MISC2=bar", "R_LIBS_SITE=path/to/install/lib"},
		},
		{
			"R_LIBS_SITE and R_LIBS_USER env set",
			"",
			[]string{"R_LIBS_USER=original/path", "R_LIBS_SITE=original/site/path", "MISC2=bar"},
			[]string{"MISC2=bar", "R_LIBS_SITE=path/to/install/lib"},
		},
		{
			"R_LIBS_SITE and R_LIBS_USER env set",
			"dplyr",
			[]string{"R_LIBS_USER=original/path", "R_LIBS_SITE=original/site/path", "MISC2=bar"},
			[]string{"DPLYR_ENV=true", "MISC2=bar", "R_LIBS_SITE=path/to/install/lib",},
		},
	}
	for i, tt := range installArgsTests {
		actual := configureEnv(tt.sysEnv, defaultRS, tt.in)
		assert.Equal(tt.expected, actual, fmt.Sprintf("%s, test num: %v", tt.context, i+1))
	}

}
