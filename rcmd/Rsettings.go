package rcmd

import (
	"fmt"
	"strings"
)

// NewRSettings initializes RSettings
func NewRSettings() RSettings {
	return RSettings{
		GlobalEnvVars: make(map[string]string),
		PkgEnvVars:    make(map[string]map[string]string),
	}
}

// R provides a cleaned path to the R executable
func (rs RSettings) R() string {
	rpath := rs.Rpath
	if rpath == "" {
		return ("R")
	}
	// TODO: check if this could have problems with trailing slash on windows
	// TODO: better to use something like filepath.clean? would that sanitize better?

	// Need to trim trailing slash as will form the R CMD syntax
	// eg /path/to/R CMD, so can't have /path/to/R/ CMD
	return strings.TrimSuffix(rpath, "/")
}

// LibPathsEnv returns the library formatted in the style to be set as an environment variable
func (rs RSettings) LibPathsEnv() (bool, string) {
	if len(rs.LibPaths) == 0 {
		return false, ""
	}
	if len(rs.LibPaths) == 1 && rs.LibPaths[0] == "" {
		return false, ""
	}
	return true, fmt.Sprintf("R_LIBS_SITE=%s", strings.Join(rs.LibPaths, ":"))
}
