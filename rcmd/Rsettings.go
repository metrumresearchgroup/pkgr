package rcmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/metrumresearchgroup/pkgr/cran"
	"github.com/metrumresearchgroup/pkgr/rcmd/rp"
	. "github.com/metrumresearchgroup/pkgr/logger"
	"github.com/spf13/afero"
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

// GetRVersion returns the R version as set in RSettings
// unlike the other methods, this one is a pointer, as RVersion mutates the known R Version,
// as if it is not defined, it will shell out to R to determine the version, and mutate itself
// to set that value, while also returning the RVersion.
// This will keep any program using rs from needing to shell out multiple times
func GetRVersion(rs *RSettings) cran.RVersion {
	if rs.Version.ToString() == "0.0" {
		res, err := RunR(afero.NewOsFs(), "", *rs, "paste0(R.Version()$major,'.',R.Version()$minor)", "")
		if err != nil {
			Log.Fatal("error getting R version")
			return cran.RVersion{}
		}
		rVersionString := rp.ScanLines(res)[0]
		rsp := strings.Split(strings.Replace(rVersionString, "\"", "", -1), ".")
		if len(rsp) == 3 {
			maj, _ := strconv.Atoi(rsp[0])
			min, _ := strconv.Atoi(rsp[1])
			pat, _ := strconv.Atoi(rsp[2])
			// this should now make it so in the future it will be set so should only need to shell out to R once
			rs.Version = cran.RVersion{
				Major: maj,
				Minor: min,
				Patch: pat,
			}
		} else {
			Log.Fatal("error getting R version")
		}
	}
	return rs.Version
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
