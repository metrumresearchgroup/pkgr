package rcmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBinaryName(t *testing.T) {
	tests := []struct {
		os       string
		platform string
		expected string
	}{
		{
			os:       "linux",
			platform: "x86_64-pc-linux-gnu",
			expected: "pkg_version_R_x86_64-pc-linux-gnu.tar.gz",
		},
		{
			os:       "darwin",
			expected: "pkg_version.tgz",
		},
		{
			os:       "windows",
			expected: "pkg_version.zip",
		},
	}
	for _, tt := range tests {
		name := binaryNameOs(tt.os, "pkg", "version", tt.platform)
		assert.Equal(t, tt.expected, name, fmt.Sprintf("Not equal: %s", tt.os))
	}
}

func TestBinaryExt(t *testing.T) {
	tests := []struct {
		os       string
		path     string
		expected string
	}{
		{
			os:       "linux",
			path:     "/var/tmp/gz",
			expected: "gz",
		},
		{
			os:       "darwin",
			path:     "/var/tmp/tgz",
			expected: "tgz",
		},
		{
			os:       "windows",
			path:     "/var/tmp/zip",
			expected: "zip",
		},
	}
	for _, tt := range tests {
		name := binaryExtOs(tt.os, tt.path, tt.os)
		assert.Equal(t, tt.expected, name, fmt.Sprintf("Not equal: %s", tt.os))
	}
}
