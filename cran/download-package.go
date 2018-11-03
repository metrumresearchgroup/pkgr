package cran

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/dpastoor/goutils"
	"github.com/dpastoor/rpackagemanager/desc"
	"github.com/spf13/afero"
)

// SourceType represents the type of package to download
type SourceType int

// Constraints on package deps
// Least to most constraining
const (
	Source SourceType = iota
	Binary
)

// Download provides information about the package download
type Download struct {
	Type SourceType
	Path string
	New  bool
}

// PkgDl holds the metadata needed to download a package
type PkgDl struct {
	Package desc.Desc
	Repo    RepoURL
}

// DownloadPackages downloads a set of packages concurrently
func DownloadPackages(fs afero.Fs, ds []PkgDl, st SourceType, dir string) (map[string](Download), error) {
	result := NewSyncMap()
	sem := make(chan struct{}, 10)
	wg := sync.WaitGroup{}
	for _, d := range ds {
		wg.Add(1)
		go func(d PkgDl, wg *sync.WaitGroup) {
			sem <- struct{}{}
			defer func() {
				<-sem
				wg.Done()
			}()
			pkgdir := filepath.Join(dir, fmt.Sprintf("%s_%s.tar.gz", d.Package.Package, d.Package.Version))
			dl, err := DownloadPackage(fs, d.Package, d.Repo.URL, st, pkgdir)
			if err != nil {
				fmt.Println("downloading failed for package: ", d.Package.Package)
				return
			}
			result.Put(d.Package.Package, dl)
		}(d, &wg)
	}
	wg.Wait()
	fmt.Println("all packages downloaded")
	return result.Map, nil
}

// DownloadPackage should download a package tarball if it doesn't exist and return
// the path to the downloaded tarball
func DownloadPackage(fs afero.Fs, d desc.Desc, url string, st SourceType, dest string) (Download, error) {
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
		fmt.Println("already have ", d.Package, " downloaded")
		return Download{
			Type: st,
			Path: dest,
			New:  false,
		}, nil
	}
	pkgdl := fmt.Sprintf("%s/src/contrib/%s_%s.tar.gz", strings.TrimSuffix(url, "/"), d.Package, d.Version)
	fmt.Println("downloading from: ", pkgdl)
	resp, err := http.Get(pkgdl)
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
		Type: st,
		Path: dest,
		New:  true,
	}, nil
}
