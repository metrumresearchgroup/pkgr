package rcmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstallArgs(t *testing.T) {
	assert := assert.New(t)

	var installArgsTests = []struct {
		in       *InstallArgs
		expected []string
	}{
		{
			&InstallArgs{},
			[]string{},
		},
		{
			NewDefaultInstallArgs(),
			[]string{"--build", "--no-multiarch", "--with-keep.source"},
		},
	}
	for i, tt := range installArgsTests {
		actual := tt.in.CliArgs()
		assert.Equal(actual, tt.expected, fmt.Sprintf("test num: %v", i+1))

	}
}
