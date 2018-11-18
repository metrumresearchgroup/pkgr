package cmd

import (
	"os"
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
	return cdir
}
