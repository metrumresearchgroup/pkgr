package cmd

import (
	"bytes"

	"encoding/json"

	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"github.com/metrumresearchgroup/pkgr/desc"
)

// returns the cache or sets to a cache dir
// TODO: verify that tests use localized caches instead of relying on this universal cache
//  as it may cause test interference.
// TODO: instead of userCache, let's call this something closer to CacheDir() at least. Especially if we have an option
//  to collect errors.
func userCache(pc string) string {
	// if actually set then use that cache dir
	if pc != "" {
		log.WithField("dir", pc).Trace("package cache directory set by user")
		return pc
	}
	// TODO: consider UserCacheDir/errors for dumping of error reports, if generated, and a tool to report the error
	//  from this cache of errors, including possibly values passed into the function.
	cdir, err := os.UserCacheDir()
	if err != nil {
		log.Warn("could not use user cache dir, using temp dir")
		cdir = os.TempDir()
	}
	log.WithField("dir", cdir).Trace("default package cache directory")

	pkgrCacheDir := filepath.Join(cdir, "pkgr")

	return pkgrCacheDir
}

// TODO: Are we expecting high load or high io-waiting state to constrain us? Go has pretty free forking internally, but
//  exec is constrained by shell overhead running the commands to fetch the packages.
// If user has not specified a thread count themselves, will limit the user to 8 threads max to avoid issues.
func getWorkerCount(threadCount, numCpus int) int {
	// TODO: I'd probably put the logic to call runtime.NumCPU() here if numCpus == 0 so the zero value has meaning.
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
		// TODO: this can be re-arranged into a multi-return instead of single-return.
		// TODO: the negative case can be a short  > 0 check with an early-return, reducing indent of 0 state calculation.
		// TODO: The warning doesn't have a place in the function that retrieves that data. This is a dual purpose.
		//  as it has an effect on the outside world in this log line.
		// TODO: where it is implemented, it could be avoided to set it to numCpus instead, overriding the user setting for safety.
		if nworkers > numCpus+2 {
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

// TODO: check if this is needed or is something used and then disused.
func libraryExists(fileSystem afero.Fs, libraryPath string) bool {
	result, _ := afero.Exists(fileSystem, libraryPath)
	return result
}

// Adapted from https://stackoverflow.com/questions/28595664/how-to-stop-json-marshal-from-escaping-and
func JsonMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "    ")
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

// TODO: name does not describe what this function does. Intent is muddied by this.
//  Did this originally take a []desc.Desc?
func extractNamesFromDesc(installedPackages map[string]desc.Desc) []string {
	// TODO: make from len(installedPackages) as capacity
	var installedPackageNames []string
	for key := range installedPackages {
		installedPackageNames = append(installedPackageNames, key)
	}
	return installedPackageNames
}
