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
			&InstallArgs{Clean: true},
			[]string{"--clean"},
		},
		{
			NewDefaultInstallArgs(),
			[]string{"--build", "--install-tests", "--no-multiarch", "--with-keep.source"},
		},
	}
	for i, tt := range installArgsTests {
		actual := tt.in.CliArgs()
		assert.Equal(actual, tt.expected, fmt.Sprintf("test num: %v", i+1))

	}
}
