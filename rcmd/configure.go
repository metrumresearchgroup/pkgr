package rcmd

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus"
)

// sysEnvVars contains the default environment variables usually from
// os.Environ()
func configureEnv(sysEnvVars []string, rs RSettings, pkg string) []string {
	envMap := make(map[string]string)
	envVars := []string{}
	envOrder := []string{}

	pkgEnv, hasCustomEnv := rs.PkgEnvVars[pkg]
	if hasCustomEnv {
		// not sure if this is needed when logging maps but for simple json want a single string
		// so will also collect in a separate set of envs and log as a single combined string
		for k, v := range pkgEnv {
			envVars = append(envVars, fmt.Sprintf("%s=%s", k, v))
			envMap[k] = v
		}
		log.WithFields(logrus.Fields{
			"envs":    envVars,
			"package": pkg,
		}).Trace("Custom Environment Variables")
	}
	for k, v := range rs.GlobalEnvVars {
		_, exists := envMap[k]
		if !exists {
			envMap[k] = v
			envOrder = append(envOrder, k)
		}
	}
	// system env vars generally
	for _, ev := range sysEnvVars {
		evs := strings.SplitN(ev, "=", 2)
		if len(evs) > 1 && evs[1] != "" {

			// we don't want to track the order of these anyway since they should take priority in the end
			// R_LIBS_USER takes precidence over R_LIBS_SITE
			// so will cause the loading characteristics to
			// not be representative of the hierarchy specified
			// in Library/Libpaths in the pkgr configuration
			// we only want R_LIBS_SITE set to control all relevant library paths for the user to
			if evs[0] == "R_LIBS_USER" {
				log.WithField("path", evs[1]).Debug("overriding system R_LIBS_USER")
				continue
			}
			if evs[0] == "R_LIBS_SITE" {
				log.WithField("path", evs[1]).Debug("overriding system R_LIBS_USER")
				continue
			}

			// if exists would be custom to the package hence should not accept the system env
			_, exists := envMap[evs[0]]
			if !exists {
				envMap[evs[0]] = evs[1]
				envOrder = append(envOrder, evs[0])
			}
		}
	}

	// TODO: determine if using globalenvvars as a map could cause subtle bug given ordering precedence

	ok, lp := rs.LibPathsEnv()
	if ok {
		envVars = append(envVars, lp)
	}

	for _, ev := range envOrder {
		val := envMap[ev]
		envVars = append(envVars, fmt.Sprintf("%s=%s", ev, val))
	}

	return envVars
}
