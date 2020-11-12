package rcmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

// sysEnvVars contains the default environment variables usually from
// os.Environ()
func configureEnv(sysEnvVars []string, rs RSettings, pkg string) []string {
	envList := NvpList{}
	envVars := []string{}

	for _, p := range rs.GlobalEnvVars.Pairs {
		_, exists := envList.Get(p.Name)
		if !exists {
			envList.Append(p.Name, p.Value)
		}
	}

	pkgEnv, hasCustomEnv := rs.PkgEnvVars[pkg]
	if hasCustomEnv {
		// not sure if this is needed when logging maps but for simple json want a single string
		// so will also collect in a separate set of envs and log as a single combined string
		for k, v := range pkgEnv {
			envList.Append(k, v)
		}
		log.WithFields(log.Fields{
			"envs":    pkgEnv,
			"package": pkg,
		}).Trace("Custom Environment Variables")
	}

	// system env vars generally
	for _, ev := range sysEnvVars {
		evs := strings.SplitN(ev, "=", 2)
		if len(evs) > 1 && evs[1] != "" {

			// we don't want to track the order of these anyway since they should take priority in the end
			// R_LIBS_USER takes precedence over R_LIBS_SITE
			// so will cause the loading characteristics to
			// not be representative of the hierarchy specified
			// in Library/Libpaths in the pkgr configuration
			// we only want R_LIBS_SITE set to control all relevant library paths for the user to
			if evs[0] == "R_LIBS_USER" {
				log.WithField("path", evs[1]).Debug("overriding system R_LIBS_USER")
				continue
			}
			if evs[0] == "R_LIBS_SITE" {
				log.WithField("path", evs[1]).Debug("overriding system R_LIBS_SITE")
				continue
			}
			if evs[0] == "PATH" {
				condaPath := "/opt/anaconda/bin"
				if strings.HasPrefix(evs[1], condaPath) && !rs.PreservePath {
					evs[1] = strings.TrimPrefix(evs[1], fmt.Sprintf("%s:", condaPath))
				}
				rDir := filepath.Dir(rs.Rpath)
				if rDir != "" && rDir != "." && !strings.HasPrefix(evs[1], rDir) {
					evs[1] = fmt.Sprintf("%s:%s", rDir, evs[1])
					log.WithField("path", evs[1]).Debug("adding Rpath to front of system PATH")
				}
			}
			// if exists would be custom to the package hence should not accept the system env
			_, exists := envList.Get(evs[0])
			if !exists {
				displayValue := evs[1]
				censoredVars := censoredEnvVars(nil)
				_, isCensored := censoredVars[strings.ToUpper(evs[0])]

				if isCensored {
					displayValue = "**HIDDEN**"
				}
				envList.Append(evs[0], displayValue)
			}
		}
	}

	// Force R_LIBS_USER to be an empty dir so that we can be sure it won't get overridden by default R paths.
	tmpdir := filepath.Join(
		os.TempDir(),
		randomString(12),
	)
	err := os.MkdirAll(tmpdir, 0777)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("error making temporary directory while overriding R_LIBS_USER for install.")
	}
	envList.Append("R_LIBS_USER", tmpdir)

	ok, lp := rs.LibPathsEnv()
	if ok {
		envList.AppendNvp(lp)
	}

	for _, p := range envList.Pairs {
		// the one and only place to append name=value strings to envVars
		envVars = append(envVars, p.GetString(p.Name))
	}

	return envVars
}

// Returns a constant set of env vars to be hidden in logs.
// Calling function may pass in additional values to include in the censor list.
func censoredEnvVars(add []string) map[string]string {
	censoredVarsMap := map[string]string{
		"GITHUB_TOKEN":      "GITHUB_TOKEN",
		"GITHUB_PAT":        "GITHUB_PAT",
		"GHE_TOKEN":         "GHE_TOKEN",
		"GHE_PAT":           "GHE_PAT",
		"AWS_ACCESS_KEY_ID": "AWS_ACCESS_KEY_ID",
		"AWS_SECRET_KEY":    "AWS_SECRET_KEY",
	}
	if add != nil {
		for _, v := range add {
			censoredVarsMap[strings.ToUpper(v)] = strings.ToUpper(v) // append(censoredVars, strings.ToUpper(v))
		}
	}
	return censoredVarsMap
}
