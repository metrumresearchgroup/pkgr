package cran

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/dpastoor/goutils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// SourceType represents the type of package to download
type SourceType int

func (s SourceType) String() string {
	if s == Default {
		s = DefaultType()
	}
	if s == Binary {
		return "binary"
	}
	return "source"
}

// Constraints on package deps
// Least to most constraining
const (
	// Use Default so can actually control default behavior in different
	// contexts rather than defaulting to Source as the 0 type
	// for example on Windows/Mac if Default then should use Binary
	Default SourceType = iota
	Source
	Binary
)

func getRepos(ds []PkgDl) map[string]RepoURL {
	rpm := make(map[string]RepoURL)
	for _, d := range ds {
		rpm[d.Config.GetOrigin().Name] = d.Config.GetOrigin()
	}
	return rpm
}

// DownloadPackages downloads a set of packages concurrently
func DownloadPackages(fs afero.Fs, ds []PkgDl, baseDir string, rv RVersion) (*PkgMap, error) {
	startTime := time.Now()
	result := NewPkgMap()
	sem := make(chan struct{}, 10)
	wg := sync.WaitGroup{}
	rpm := getRepos(ds)
	for _, r := range rpm {
		urlHash := RepoURLHash(r)
		for _, pt := range []string{"src", filepath.Join("binary", rv.ToString())} {
			pkgdir := filepath.Join(baseDir, urlHash, pt)
			err := fs.MkdirAll(pkgdir, 0777)
			if err != nil {
				log.WithField("dir", pkgdir).WithField("error", err.Error()).Fatal("error creating package directory ")
				return result, err
			}
		}
	}
	log.WithField("dir", baseDir).Info("downloading required packages within directory ")
	for _, d := range ds {
		wg.Add(1)
		go func(d PkgDl, wg *sync.WaitGroup) {
			var pkgType string
			if d.Config.GetSourceType2() == Default {
				d.Config.SetSourceType(DefaultType().String())
			}
			switch d.Config.GetSourceType2() {
			case Binary:
				pkgType = "binary"
			case Source:
				pkgType = "src"
			default:
				pkgType = "src"
			}
			sem <- struct{}{}
			defer func() {
				<-sem
				wg.Done()
			}()
			// TODO: should potentially provide a lookup mapping
			// but would want to do this outside the goroutine that downloads\
			// the package so didn't get invoked multiple times
			urlHash := RepoURLHash(d.Config.GetOrigin())
			pkgdir := filepath.Join(baseDir, urlHash, pkgType)

			var pkgFile string

			if d.Config.IsDownloadable() {
				if d.Config.GetSourceType2() == Binary {
					pkgFile = filepath.Join(pkgdir, rv.ToString(), binaryName(d.Package.Package, d.Package.Version))
				} else {
					pkgFile = filepath.Join(pkgdir, fmt.Sprintf("%s_%s.tar.gz", d.Package.Package, d.Package.Version))
				}
			} else {
				pkgFile = d.Config.GetOrigin().URL
			}


			startDl := time.Now()
			dl, err := DownloadPackage(fs, d, pkgFile, rv)
			if err != nil {
				// TODO:  should this cause a failure downstream rather than just printing
				// as right now it keeps running and just doesn't install that package?
				log.WithField("package", d.Package.Package).Warn("downloading failed")
				return
			}

			if dl.New {
				log.WithFields(log.Fields{
					"package": d.Package.Package,
					"dltime":  time.Since(startDl),
					"size":    fmt.Sprintf("%.2f MB", dl.GetMegabytes()),
				}).Debug("download successful")
			}
			result.Put(d.Package.Package, dl)
		}(d, &wg)
	}
	wg.Wait()
	log.WithField("duration", time.Since(startTime)).Info("all packages downloaded")
	return result, nil
}

// DownloadPackage should download a package tarball if it doesn't exist and return
// the path to the downloaded tarball
func DownloadPackage(fs afero.Fs, d PkgDl, dest string, rv RVersion) (Download, error) {
	if !filepath.IsAbs(dest) {
		cwd, _ := os.Getwd()
		// turn to absolute
		dest = filepath.Clean(filepath.Join(cwd, dest))
	}
	exists, err := goutils.Exists(fs, dest)
	if err != nil {
		return Download{}, err
	}
	if exists {
		log.WithField("package", d.Package.Package).Debug("package already downloaded ")
		return Download{
			Path:     dest,
			New:      false,
			Metadata: d,
			Size:     0,
		}, nil
	}
	var pkgdl string
	if d.Config.GetSourceType2() == Source || !SupportsCranBinary() {
		d.Config.SetSourceType(Source.String())// in case was originally set to binary
		if !SupportsCranBinary() {
			log.WithField("package", d.Package.Package).Debug("CRAN binary not supported, downloading source instead ")
		}
		pkgdl = fmt.Sprintf("%s/src/contrib/%s", strings.TrimSuffix(d.Config.GetOrigin().URL, "/"), filepath.Base(dest))
	} else {
		pkgdl = fmt.Sprintf("%s/bin/%s/contrib/%s/%s",
			strings.TrimSuffix(d.Config.GetOrigin().URL, "/"),
			cranBinaryURL(),
			rv.ToString(),
			filepath.Base(dest))
	}

	log.WithField("package", d.Package.Package).Info("downloading package ")
	var from io.ReadCloser
	if strings.HasPrefix(pkgdl, "http") {
		resp, err := http.Get(pkgdl)
		if resp.StatusCode != 200 {
			log.WithFields(log.Fields{
				"package":     d.Package.Package,
				"url":         pkgdl,
				"status":      resp.Status,
				"status_code": resp.StatusCode,
			}).Warn("bad server response")
			b, _ := ioutil.ReadAll(resp.Body)
			log.WithField("package", d.Package.Package).Println("body: ", string(b))
			return Download{}, err
		}
		// TODO: the response will often be valid, but return like a server 404 or other error
		if err != nil {
			log.WithField("package", d.Package).Warn("error downloading package")
			return Download{}, err
		}
		from = resp.Body
		defer resp.Body.Close()
	} else {
		from, err = fs.Open(pkgdl)
		if err != nil {
			log.WithFields(log.Fields{
				"package": d.Package.Package,
				"path":    pkgdl,
			}).Fatal("missing package")
		}
	}

	file, err := fs.Create(dest)
	if err != nil {
		log.WithFields(log.Fields{
			"package": d.Package,
			"err":     err,
		}).Warn("error downloading package, no tarball created")
		return Download{}, err
	}
	defer file.Close()
	size, err := io.Copy(file, from)
	if err != nil {
		return Download{}, err
	}

	return Download{
		Path:     dest,
		New:      true,
		Metadata: d,
		Size:     size,
	}, nil
}
