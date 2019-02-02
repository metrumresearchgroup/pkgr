package rcmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigureArgs(t *testing.T) {
	assert := assert.New(t)

	var installArgsTests = []struct {
		in       string
		expected []string
	}{
		{
			"",
			[]string{"R_LIBS_SITE=path/to/install/lib", "R_LIBS_USER=path/to/install/lib"},
		},
		{
			"dplyr",
			[]string{"R_LIBS_SITE=path/to/install/lib", "R_LIBS_USER=path/to/install/lib"},
		},
	}
	defaultRS := RSettings{}
	defaultRS.LibPaths = []string{"path/to/install/lib"}
	for i, tt := range installArgsTests {
		actual := configureEnv([]string{}, defaultRS, tt.in)
		assert.Equal(actual, tt.expected, fmt.Sprintf("test num: %v", i+1))

	}
}
