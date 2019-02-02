package rcmd

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// envVars contains the default environment variables usually from
// os.Environ()
func configureEnv(envVars []string, rs RSettings, pkg string) []string {
	envMap := make(map[string]string)
	for _, ev := range envVars {
		evs := strings.SplitN(ev, "=", 2)
		if len(evs) > 1 && evs[1] != "" {
			envMap[evs[0]] = evs[1]
		}
	}
	rlu, exists := envMap["R_LIBS_USER"]
	if exists {
		// R_LIBS_USER takes precidence over R_LIBS_SITE
		// so will cause the loading characteristics to
		// not be representative of the hierarchy specified
		// in Library/Libpaths in the pkgr configuration
		delete(envMap, "R_LIBS_USER")
		log.WithField("path", rlu).Debug("deleting set R_LIBS_USER")
	}
	envVars = []string{}
	for k, v := range rs.GlobalEnvVars {
		envMap[k] = v
	}

	ok, lp := rs.LibPathsEnv()
	if ok {
		// if LibPaths set, lets drop R_LIBS_SITE set as an ENV and instead
		// add the generated R_LIBS_SITE from LibPathsEnv
		delete(envMap, "R_LIBS_SITE")
		envVars = append(envVars, lp)
	}
	for k, v := range envMap {
		envVars = append(envVars, fmt.Sprintf("%s=%s", k, v))
	}
	pkgEnv, hasCustomEnv := rs.PkgEnvVars[pkg]
	if hasCustomEnv {
		// not sure if this is needed when logging maps but for simple json want a single string
		// so will also collect in a separate set of envs and log as a single combined string
		var pkgEnvForLog []string
		for k, v := range pkgEnv {
			envVars = append(envVars, fmt.Sprintf("%s=%s", k, v))
			pkgEnvForLog = append(pkgEnvForLog, fmt.Sprintf("%s=%s", k, v))
		}
		log.WithFields(logrus.Fields{
			"envs":    pkgEnvForLog,
			"package": pkg,
		}).Trace("Custom Environment Variables")
	}
	// double down on overwriting any specification of user customization
	// and set R_LIBS_SITE to the same as the user
	envVars = append(envVars, strings.Replace(lp, "R_LIBS_SITE", "R_LIBS_USER", 1))

	return envVars
}
