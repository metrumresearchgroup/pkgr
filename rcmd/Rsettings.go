package rcmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/metrumresearchgroup/pkgr/cran"
	"github.com/metrumresearchgroup/pkgr/rcmd/rp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// NewRSettings initializes RSettings
func NewRSettings(rPath string) RSettings {
	return RSettings{
		GlobalEnvVars: make(map[string]string),
		PkgEnvVars:    make(map[string]map[string]string),
		Rpath:         rPath,
	}
}

// R provides a cleaned path to the R executable
func (rs RSettings) R() string {
	if rs.Rpath == "" {
		return ("R")
	}
	// TODO: check if this could have problems with trailing slash on windows
	// TODO: better to use something like filepath.clean? would that sanitize better?

	// Need to trim trailing slash as will form the R CMD syntax
	// eg /path/to/R CMD, so can't have /path/to/R/ CMD
	return strings.TrimSuffix(rs.Rpath, "/")
}

// GetRVersion returns the R version, and sets R Version and R platform in RSettings
// unlike the other methods, this one is a pointer, as RVersion mutates the known R Version,
// as if it is not defined, it will shell out to R to determine the version, and mutate itself
// to set that value, while also returning the RVersion.
// This will keep any program using rs from needing to shell out multiple times
func GetRVersion(rs *RSettings) cran.RVersion {
	if rs.Version.ToString() == "0.0" {
		res, err := RunR(afero.NewOsFs(), "", *rs, "version", "")
		if err != nil {
			log.Fatal("error getting R version info")
			return cran.RVersion{}
		}

		rs.Version, rs.Platform = parseVersionData(res)

		// lines := rp.ScanLines(res)
		// for _, line := range lines {
		// 	s := strings.Fields(line)
		// 	switch s[0] {
		// 	case "version.string":
		// 		rsp := strings.Split(strings.Replace(s[3], "\"", "", -1), ".")
		// 		if len(rsp) == 3 {
		// 			maj, _ := strconv.Atoi(rsp[0])
		// 			min, _ := strconv.Atoi(rsp[1])
		// 			pat, _ := strconv.Atoi(rsp[2])
		// 			// this should now make it so in the future it will be set so should only need to shell out to R once
		// 			rs.Version = cran.RVersion{
		// 				Major: maj,
		// 				Minor: min,
		// 				Patch: pat,
		// 			}
		// 		} else {
		// 			log.Fatal("error getting R version")
		// 		}
		// 	case "platform":
		// 		// set platform in Rsettings for future access
		// 		rs.Platform = s[1]
		// 	default:
		// 	}
		// }
	}
	return rs.Version
}

func parseVersionData(data []byte) (version cran.RVersion, platform string) {
	lines := rp.ScanLines(data)
	for _, line := range lines {
		s := strings.Fields(line)
		switch s[0] {
		case "version.string":
			rsp := strings.Split(strings.Replace(s[3], "\"", "", -1), ".")
			if len(rsp) == 3 {
				maj, _ := strconv.Atoi(rsp[0])
				min, _ := strconv.Atoi(rsp[1])
				pat, _ := strconv.Atoi(rsp[2])
				// this should now make it so in the future it will be set so should only need to shell out to R once
				version = cran.RVersion{
					Major: maj,
					Minor: min,
					Patch: pat,
				}
			} else {
				log.Fatal("error getting R version")
			}
		case "platform":
			// set platform in Rsettings for future access
			platform = s[1]
		default:
		}
	}
	return version, platform
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
