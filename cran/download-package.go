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
		rpm[d.Config.Repo.Name] = d.Config.Repo
	}
	return rpm
}

// DownloadPackages downloads a set of packages concurrently
func DownloadPackages(fs afero.Fs, ds []PkgDl, baseDir string) (*PkgMap, error) {
	startTime := time.Now()
	result := NewPkgMap()
	sem := make(chan struct{}, 10)
	wg := sync.WaitGroup{}
	rpm := getRepos(ds)
	for _, r := range rpm {
		urlHash := RepoURLHash(r)
		for _, pt := range []string{"src", "binary"} {
			pkgdir := filepath.Join(baseDir, urlHash, pt)
			err := fs.MkdirAll(pkgdir, 0777)
			if err != nil {
				fmt.Println("error creating package directory ", pkgdir)
				fmt.Println("error: ", err)
				return result, err
			}
		}
	}
	for _, d := range ds {
		wg.Add(1)
		go func(d PkgDl, wg *sync.WaitGroup) {
			var pkgType string
			if d.Config.Type == Default {
				d.Config.Type = DefaultType()
			}
			switch d.Config.Type {
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
			urlHash := RepoURLHash(d.Config.Repo)
			pkgdir := filepath.Join(baseDir, urlHash, pkgType)
			var pkgFile string
			if d.Config.Type == Binary {
				pkgFile = filepath.Join(pkgdir, binaryName(d.Package.Package, d.Package.Version))
			} else {
				pkgFile = filepath.Join(pkgdir, fmt.Sprintf("%s_%s.tar.gz", d.Package.Package, d.Package.Version))
			}
			dl, err := DownloadPackage(fs, d, pkgFile)
			if err != nil {
				// TODO: this should cause a failure downstream rather than just printing
				// as right now it keeps running and just doesn't install that package
				fmt.Println("downloading failed for package: ", d.Package.Package)
				return
			}
			result.Put(d.Package.Package, dl)
		}(d, &wg)
	}
	wg.Wait()
	fmt.Printf("all packages downloaded in %s\n\n", time.Since(startTime))
	return result, nil
}

// DownloadPackage should download a package tarball if it doesn't exist and return
// the path to the downloaded tarball
func DownloadPackage(fs afero.Fs, d PkgDl, dest string) (Download, error) {
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
		fmt.Println("already have ", d.Package.Package, " downloaded")
		return Download{
			Path:     dest,
			New:      false,
			Metadata: d,
		}, nil
	}
	var pkgdl string
	if d.Config.Type == Source || !SupportsCranBinary() {
		d.Config.Type = Source // in case was originally set to binary
		pkgdl = fmt.Sprintf("%s/src/contrib/%s", strings.TrimSuffix(d.Config.Repo.URL, "/"), filepath.Base(dest))
	} else {
		// TODO: fix so isn't hard coded to 3.5 binaries
		pkgdl = fmt.Sprintf("%s/bin/%s/contrib/%s/%s", strings.TrimSuffix(d.Config.Repo.URL, "/"), cranBinaryURL(), "3.5", filepath.Base(dest))
	}
	resp, err := http.Get(pkgdl)
	if resp.StatusCode != 200 {
		fmt.Println("error downloading package", d.Package.Package)
		fmt.Println("from URL:", pkgdl)
		fmt.Println("status: ", resp.Status)
		fmt.Println("status code: ", resp.StatusCode)
		b, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("body: ", string(b))
		return Download{}, err
	}
	// TODO: the response will often be valid, but return like a server 404 or other error
	if err != nil {
		fmt.Println("error downloading package", d.Package)
		return Download{}, err
	}
	defer resp.Body.Close()
	file, err := fs.Create(dest)
	if err != nil {
		fmt.Println("couldn't create tarball for ", d.Package)
		fmt.Println(err)
		return Download{}, err
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return Download{}, err
	}

	return Download{
		Path:     dest,
		New:      true,
		Metadata: d,
	}, nil
}
