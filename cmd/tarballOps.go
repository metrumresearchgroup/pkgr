package cmd

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/metrumresearchgroup/pkgr/desc"
	"github.com/metrumresearchgroup/pkgr/gpsr"
	"github.com/mholt/archiver/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"io"
	"os"
	"path/filepath"
)

func unpackTarballs(fs afero.Fs, tarballs []string, cache string) ([]desc.Desc, map[string]gpsr.AdditionalPkg) {
	cacheDir := userCache(cache)

	var descriptions []desc.Desc
	untarredMap := make(map[string]gpsr.AdditionalPkg)

	for _, tarballPath := range tarballs {

		// replace
		untarredFolder := untar(fs, tarballPath, cacheDir)

		reader, err := fs.Open(filepath.Join(untarredFolder, "DESCRIPTION"))
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"file":  untarredFolder,
				"error": err,
			}).Fatal("error opening DESCRIPTION file for tarball package")
		}

		desc, err := desc.ParseDesc(reader)
		reader.Close()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"file":  untarredFolder,
				"error": err,
			}).Fatal("error parsing DESCRIPTION file for tarball package")
		}
		descriptions = append(descriptions, desc)
		untarredMap[desc.Package] = gpsr.AdditionalPkg{InstallPath: untarredFolder, OriginPath: tarballPath, Type: "tarball"}
	}

	return descriptions, untarredMap
}

// Returns path to top-level package folder of untarred files
func untar(fs afero.Fs, path string, cacheDir string) string {

	// Part 1
	// Create hash of Tarball to use for a folder name in the cache.
	tgzFileForHash, err := fs.Open(path)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"path": path,
		}).Fatal("error opening second copy of specified tarball for hashing")
	}
	defer tgzFileForHash.Close()

	// Use a hash of the file so that we always regenerate when the tarball is updated.
	// Probably unnecessary, but prevents problems in imaginary cases such as "two instances of pkgr are sharing a cache
	// and installing two different versions of the 'same' tarball."
	tarballDirectoryName, err := getHashedTarballName(tgzFileForHash)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"path": path,
		}).Fatalf("error while creating hash for tarball in cache: %s", err)
	}
	tarballDirectoryPath := filepath.Join(cacheDir, tarballDirectoryName)
	tarballInCache, err := afero.DirExists(fs, tarballDirectoryPath)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"cached_location": tarballDirectoryPath,
			"source_tarball":  path,
			"error":           err,
		}).Warn("encountered problem checking cache for existing tarball. Extracting tarball to cache_location anyway.")
	}
	if tarballInCache {
		logrus.WithFields(logrus.Fields{
			"cached_tarball": tarballDirectoryPath,
			"source_tarball": path,
		}).Debug("using cached tarball.")
	} else {
		err = archiver.Unarchive(path, tarballDirectoryPath)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"path": path,
			}).Fatalf("error while unarchiving tarball: %s", err)
		}
	}

	//Part 2
	dirEntries, err := afero.ReadDir(fs, tarballDirectoryPath)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"directory": tarballDirectoryPath,
			"tarball":   path,
			"error":     err,
		}).Error("error encountered while reading untarred directory")
	}

	var extractedDirs []os.FileInfo
	for _, entry := range dirEntries {
		if entry.IsDir() {
			extractedDirs = append(extractedDirs, entry)
		} else {
			logrus.WithFields(logrus.Fields{
				"tarball":           path,
				"untarredDirectory": tarballDirectoryPath,
				"item":              entry.Name(),
			}).Trace("extraneous item found in untarred directory. If you have previously installed this tarball via pkgr, it may be a build artifact from the last installation")
		}
	}

	if len(extractedDirs) == 0 {
		logrus.WithFields(logrus.Fields{
			"directory": tarballDirectoryPath,
			"tarball":   path,
		}).Fatal("untarred directory is empty -- cannot install tarball package")
	} else if len(extractedDirs) > 1 {
		var entries []string
		for _, entry := range extractedDirs {
			entries = append(entries, entry.Name())
		}
		logrus.WithFields(logrus.Fields{
			"directory": tarballDirectoryPath,
			"tarball":   path,
			"items":     entries,
			"using":     extractedDirs[0].Name(),
		}).Warn("found more than one directory at top level in unarchived tarball -- assuming first alphabetical entry is package directory")
	}

	return filepath.Join(tarballDirectoryPath, extractedDirs[0].Name())
}

func getHashedTarballName(tgzFile afero.File) (string, error) {
	hash := md5.New()
	_, err := io.Copy(hash, tgzFile)
	hashInBytes := hash.Sum(nil)[:8]
	//Convert the bytes to a string, used as a directory name for the package.
	tarballDirectoryName := hex.EncodeToString(hashInBytes)
	// Hashing code adapted from https://mrwaggel.be/post/generate-md5-hash-of-a-file-in-golang/
	return tarballDirectoryName, err
}
