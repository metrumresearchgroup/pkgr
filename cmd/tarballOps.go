package cmd

import (
	"archive/tar"
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"github.com/metrumresearchgroup/pkgr/desc"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"io"
	"os"
	"path/filepath"
)

// Tarball manipulation code taken from https://gist.github.com/indraniel/1a91458984179ab4cf80 -- is there a built-in function that does this?
func unpackTarballs(fs afero.Fs, tarballs []string, cache string) ([]desc.Desc, map[string]string) {
	cacheDir := userCache(cache)

	untarredMap := make(map[string]string)

	var untarredPaths []string
	for _, path := range tarballs {
		untarredFolder := untar(fs, path, cacheDir)
		untarredPaths = append(untarredPaths, untarredFolder)
	}

	var descriptions []desc.Desc
	for _, path := range untarredPaths {
		reader, err := fs.Open(filepath.Join(path, "DESCRIPTION"))
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"file":  path,
				"error": err,
			}).Fatal("error opening DESCRIPTION file for tarball package")
		}

		desc, err := desc.ParseDesc(reader)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"file":  path,
				"error": err,
			}).Fatal("error parsing DESCRIPTION file for tarball package")
		}
		descriptions = append(descriptions, desc)
		untarredMap[desc.Package] = path
	}

	return descriptions, untarredMap
}

// Returns path to top-level package folder of untarred files
func untar(fs afero.Fs, path string, cacheDir string) string {

	// Part 1
	tgzFile, err := fs.Open(path)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"path": path,
		}).Fatal("error processing specified tarball")
	}
	defer tgzFile.Close()
	tgzFileForHash, err := fs.Open(path) // Shouldn't fail if the first one passed, but I'll check anyway.
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"path": path,
		}).Fatal("error opening second copy of specified tarball for hashing")
	}
	defer tgzFileForHash.Close()

	// Part 1.5
	//Use a hash of the file so that we always regenerate when the tarball is updated.
	tarballDirectoryName, err := getHashedTarballName(tgzFileForHash)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"path": path,
		}).Fatalf("error while creating hash for tarball in cache: %s", err)
	}
	tarballDirectoryPath := filepath.Join(cacheDir, tarballDirectoryName)

	//Part 2
	gzipStream, err := gzip.NewReader(tgzFile)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"path": path,
		}).Fatal("error creating gzip stream for specified tarball")
	}
	defer gzipStream.Close()
	tarStream := tar.NewReader(gzipStream)
	for {
		header, err := tarStream.Next()

		if err == io.EOF {
			break
		} else if err != nil {
			logrus.Error("could not process file in tar stream. Error was: ", err)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			fs.MkdirAll(filepath.Join(tarballDirectoryPath, header.Name), 0755)
			break
		case tar.TypeReg:
			dstFile := filepath.Join(tarballDirectoryPath, header.Name)
			extractFile(dstFile, tarStream)
			break
		default:
			logrus.WithFields(logrus.Fields{
				"type": header.Typeflag,
				"path": path,
			}).Error("unknown file type found while processing tarball")
			break
		}
	}

	dirEntries, err := afero.ReadDir(fs, tarballDirectoryPath)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"directory": tarballDirectoryPath,
			"tarball":   path,
			"error":     err,
		}).Error("error encountered while reading untarred directory")
	}

	if len(dirEntries) == 0 {
		logrus.WithFields(logrus.Fields{
			"directory": tarballDirectoryPath,
			"tarball":   path,
		}).Fatal("untarred directory is empty -- cannot install tarball package")
	} else if len(dirEntries) > 1 {
		logrus.WithFields(logrus.Fields{
			"directory": tarballDirectoryPath,
			"tarball":   path,
		}).Warn("found more than one item at top level in unarchived tarball -- assuming first alphabetical entry is package directory")
	}

	return filepath.Join(tarballDirectoryPath, dirEntries[0].Name())
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

func extractFile(dstFile string, tarStream *tar.Reader) {
	outFile, err := os.Create(dstFile)
	if err != nil {
		logrus.Fatalf("ExtractTarGz: Create() failed: %s", err.Error())
	}
	defer outFile.Close()
	if _, err := io.Copy(outFile, tarStream); err != nil {
		logrus.Fatalf("ExtractTarGz: Copy() failed: %s", err.Error())
	}
}

