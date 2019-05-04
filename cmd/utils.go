package cmd

import (
	"os"
	"path/filepath"
	"runtime"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// returns the cache or sets to a cache dir
func userCache(pc string) string {
	// if actually set then use that cache dir
	if pc != "" {
		log.WithField("dir", pc).Trace("package cache directory set by user")
		return pc
	}
	cdir, err := os.UserCacheDir()
	if err != nil {
		log.Warn("could not use user cache dir, using temp dir")
		cdir = os.TempDir()
	}
	log.WithField("dir", cdir).Trace("default package cache directory")

	pkgrCacheDir := filepath.Join(cdir, "pkgr")

	return pkgrCacheDir
}

func getWorkerCount() int {
	var nworkers int
	if viper.GetInt("threads") < 1 {
		nworkers = runtime.NumCPU()
		if nworkers > 2 {
			nworkers = nworkers - 1
		}
	} else {
		nworkers = viper.GetInt("threads")
		if nworkers > runtime.NumCPU()+2 {
			log.Warn("number of workers exceeds the number of threads on machine by at least 2, this may result in degraded performance")
		}
	}
	return nworkers
}
func stringInSlice(s string, slice []string) bool {
	for _, entry := range slice {
		if s == entry {
			return true
		}
	}
	return false
}
