package cmd

import (
	"archive/tar"
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"github.com/metrumresearchgroup/pkgr/desc"
	"github.com/metrumresearchgroup/pkgr/gpsr"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	log "github.com/sirupsen/logrus"
)

// Tarball manipulation code taken from https://gist.github.com/indraniel/1a91458984179ab4cf80 -- is there a built-in function that does this?
func unpackTarballs(fs afero.Fs, tarballs []string, cache string) ([]desc.Desc, map[string]gpsr.AdditionalPkg) {
	cacheDir := userCache(cache)
	tmpDownloadDir := filepath.Join(cacheDir, "additionalDownloads")
	err := fs.MkdirAll(tmpDownloadDir, 0777)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
			"dir": tmpDownloadDir,
		}).Error("could not make temp directory to hold downloaded tarballs")
	}

	var descriptions []desc.Desc
	untarredMap := make(map[string]gpsr.AdditionalPkg)

	for _, tarballPath := range tarballs {

		var untarredFolder string
		///Download code ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
		var downloadReader io.ReadCloser
		//var downloadReaderForHash io.ReadCloser
		if strings.HasPrefix(tarballPath, "http") {
			urlObj, err := url.Parse(tarballPath)
			if err != nil {
				log.WithFields(log.Fields{
					"err": err,
					"url": tarballPath,
				}).Error("could not parse tarball URL")
				continue
			}
			resp, err := http.Get(urlObj.String())
			if resp.StatusCode != 200 {
				log.WithFields(log.Fields{
					//"package":     d.Package.Package,
					"url":         tarballPath,
					"status":      resp.Status,
					"status_code": resp.StatusCode,
				}).Error("bad server response")
				respBody, _ := ioutil.ReadAll(resp.Body)
				log.WithField("tarball", tarballPath).Println("body: ", string(respBody))
				//return cran.Download{}, err
			}
			// TODO: the response will often be valid, but return like a server 404 or other error
			if err != nil {
				log.WithField("tarball", tarballPath).Warn("error downloading package")
				//return Download{}, err
			} else { // Download successful
				downloadReader = resp.Body
				//downloadReaderForHash = resp.Body
			}

			urlHash := path.Base(urlObj.Path) //hashString(tarballPath)
			tmpTarballPath := filepath.Join(tmpDownloadDir, urlHash)
			tmpTarball, err := fs.Create(tmpTarballPath)
			if err != nil {
				log.WithFields(log.Fields{
					"err": err,
					"tmpTarballPath": tmpTarballPath,
					"tarballUrl": tarballPath,
				}).Error("error copying temporary tarball from remote location")
			}
			io.Copy(tmpTarball, downloadReader)
			downloadReader.Close()

			//fs.Open(tmpTarballPath)
			tarballHash, err := getHashedTarballName(tmpTarball)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"tarballURL": tarballPath,
				}).Errorf("error while creating hash for tarball in cache: %s", err)
				continue
			}
			tmpTarball.Close()
			unpackedDest := filepath.Join(cacheDir, tarballHash)


			tmpTarball2, err := fs.Open(tmpTarballPath)
			if err != nil {
				log.WithFields(log.Fields{
					"err": err,
					"tmpTarballPath": tmpTarballPath,
				}).Error("error accessing temporary tarball after download")
			}

			extractedDir := extractDirFromTar(tmpTarball2, tmpTarballPath, unpackedDest)
			untarredFolder = filepath.Join(unpackedDest, extractedDir.Name())

			tmpTarball2.Close()

			//defer resp.Body.Close()
			resp.Body.Close()
		} else {
		//End download code~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
			untarredFolder = untar(fs, tarballPath, cacheDir)
		}
		reader, err := fs.Open(filepath.Join(untarredFolder, "DESCRIPTION"))
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"file":  untarredFolder,
				"error": err,
			}).Fatal("error opening DESCRIPTION file for tarball package")
		}

		desc, err := desc.ParseDesc(reader)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"file":  untarredFolder,
				"error": err,
			}).Fatal("error parsing DESCRIPTION file for tarball package")
		}
		descriptions = append(descriptions, desc)
		untarredMap[desc.Package] = gpsr.AdditionalPkg{ InstallPath: untarredFolder, OriginPath: tarballPath, Type: "tarball"}
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
	extractedDir := extractDirFromTar(tgzFile, path, tarballDirectoryPath)

	return filepath.Join(tarballDirectoryPath, extractedDir.Name())
}

func extractDirFromTar(tgzReader io.Reader, tarballPath string, destDirectory string) os.FileInfo {
	gzipStream, err := gzip.NewReader(tgzReader)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"tarballPath": tarballPath,
			"err": err,
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
			fs.MkdirAll(filepath.Join(destDirectory, header.Name), 0755)
			break
		case tar.TypeReg:
			dstFile := filepath.Join(destDirectory, header.Name)
			extractFile(dstFile, tarStream)
			break
		default:
			logrus.WithFields(logrus.Fields{
				"type": header.Typeflag,
				"tarballPath": tarballPath,
			}).Error("unknown file type found while processing tarball")
			break
		}
	}
	dirEntries, err := afero.ReadDir(fs, destDirectory)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"directory": destDirectory,
			"tarball":   tarballPath,
			"error":     err,
		}).Error("error encountered while reading untarred directory")
	}
	var extractedDirs []os.FileInfo
	for _, entry := range dirEntries {
		if entry.IsDir() {
			extractedDirs = append(extractedDirs, entry)
		} else {
			logrus.WithFields(logrus.Fields{
				"tarball":           tarballPath,
				"untarredDirectory": destDirectory,
				"item":              entry.Name(),
			}).Trace("extraneous item found in untarred directory. If you have previously installed this tarball via pkgr, it may be a build artifact from the last installation")
		}
	}
	if len(extractedDirs) == 0 {
		logrus.WithFields(logrus.Fields{
			"directory": destDirectory,
			"tarball":   tarballPath,
		}).Fatal("untarred directory is empty -- cannot install tarball package")
	} else if len(extractedDirs) > 1 {
		var entries []string
		for _, entry := range extractedDirs {
			entries = append(entries, entry.Name())
		}
		logrus.WithFields(logrus.Fields{
			"directory": destDirectory,
			"tarball":   tarballPath,
			"items":     entries,
			"using":     extractedDirs[0].Name(),
		}).Warn("found more than one directory at top level in unarchived tarball -- assuming first alphabetical entry is package directory")
	}
	return extractedDirs[0]
}

func getHashedTarballName(tgzFile io.Reader) (string, error) {
	hash := md5.New()
	_, err := io.Copy(hash, tgzFile)
	hashInBytes := hash.Sum(nil)[:8]
	//Convert the bytes to a string, used as a directory name for the package.
	tarballDirectoryName := hex.EncodeToString(hashInBytes)
	// Hashing code adapted from https://mrwaggel.be/post/generate-md5-hash-of-a-file-in-golang/
	return tarballDirectoryName, err
}

//func getHashedTarballNameHttp(tgzFileBody io.ReadCloser) (string, error) {
//	hash := md5.New()
//	_, err := io.Copy(hash, tgzFileBody)
//	hashInBytes := hash.Sum(nil)[:8]
//	//Convert the bytes to a string, used as a directory name for the package.
//	tarballDirectoryName := hex.EncodeToString(hashInBytes)
//	// Hashing code adapted from https://mrwaggel.be/post/generate-md5-hash-of-a-file-in-golang/
//	return tarballDirectoryName, err
//}

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
