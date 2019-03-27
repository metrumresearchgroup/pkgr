package cmd

import (
	"os"
	"path/filepath"
	"runtime"
	. "github.com/metrumresearchgroup/pkgr/logger"
	"github.com/spf13/viper"
)

// returns the cache or sets to a cache dir
func userCache(pc string) string {
	// if actually set then use that cache dir
	if pc != "" {
		Log.WithField("dir", pc).Trace("package cache directory set by user")
		return pc
	}
	cdir, err := os.UserCacheDir()
	if err != nil {
		Log.Warn("could not use user cache dir, using temp dir")
		cdir = os.TempDir()
	}
	Log.WithField("dir", cdir).Trace("default package cache directory")

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
			Log.Warn("number of workers exceeds the number of threads on machine by at least 2, this may result in degraded performance")
		}
	}
	return nworkers
}
