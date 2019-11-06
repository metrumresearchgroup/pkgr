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

func GetWorkerCount() int {
	return GetRestrictedWorkerCount(viper.GetInt("threads"), runtime.NumCPU())
}

// If user has not specified a thread count themselves, will limit the user to 8 threads max to avoid issues.
func GetRestrictedWorkerCount(threadCount, numCpus int) int {
	var nworkers int
	if threadCount < 1 { // This indicates that the user has not specified a thread count, i.e. we're using the default thread count of "0".

		workerCap := 8 // We have decided to cap this number at 8 unless the user requests otherwise. This is to prevent issues with cyclic Suggests installs.
		if numCpus <= workerCap {
			nworkers = numCpus
		} else {
			nworkers = workerCap
		}
	} else {
		nworkers = threadCount
		if nworkers > numCpus + 2 {
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
