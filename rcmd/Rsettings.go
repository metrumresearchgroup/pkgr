package rcmd

import (
	"fmt"
	"strings"
)

// R provides a cleaned path to the R executable
func (rs RSettings) R() string {
	// TODO: check if this could have problems with trailing slash on windows
	// TODO: better to use something like filepath.clean? would that sanitize better?

	// Need to trim trailing slash as will form the R CMD syntax
	// eg /path/to/R CMD, so can't have /path/to/R/ CMD
	return strings.TrimSuffix(rs.Rpath, "/")
}

// LibPathsEnv returns the libpaths formatted in the style to be set as an environment variable
func (rs RSettings) LibPathsEnv() (bool, string) {
	if len(rs.LibPaths) == 0 {
		return false, ""
	}
	if len(rs.LibPaths) == 1 && rs.LibPaths[0] == "" {
		return false, ""
	}
	return true, fmt.Sprintf("R_LIBS_SITE=%s", strings.Join(rs.LibPaths, ":"))
}
