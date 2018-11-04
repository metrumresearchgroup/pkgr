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

// DownloadPackages downloads a set of packages concurrently
func DownloadPackages(fs afero.Fs, ds []PkgDl, st SourceType, baseDir string) (*PkgMap, error) {
	result := NewPkgMap()
	sem := make(chan struct{}, 10)
	wg := sync.WaitGroup{}
	var pkgType string
	var typeExt string
	switch st {
	case Binary:
		pkgType = "binary"
		typeExt = "tgz"
	case Source:
		pkgType = "src"
		typeExt = "tar.gz"
	default:
		pkgType = "src"
		typeExt = "tar.gz"
	}
	for _, d := range ds {
		wg.Add(1)
		go func(d PkgDl, wg *sync.WaitGroup) {
			sem <- struct{}{}
			defer func() {
				<-sem
				wg.Done()
			}()
			pkgdir := filepath.Join(baseDir, d.Repo.Name, pkgType)
			err := fs.MkdirAll(pkgdir, 0777)
			pkgFile := filepath.Join(pkgdir, fmt.Sprintf("%s_%s.%s", d.Package.Package, d.Package.Version, typeExt))
			if err != nil {
				fmt.Println("error creating package directory ", pkgdir)
				fmt.Println("error: ", err)
				// should this trigger something more impactful? probably since wouldn't download anything
				return
			}
			dl, err := DownloadPackage(fs, d, st, pkgFile)
			if err != nil {
				fmt.Println("downloading failed for package: ", d.Package.Package)
				return
			}
			result.Put(d.Package.Package, dl)
		}(d, &wg)
	}
	wg.Wait()
	fmt.Println("all packages downloaded")
	return result, nil
}

// DownloadPackage should download a package tarball if it doesn't exist and return
// the path to the downloaded tarball
func DownloadPackage(fs afero.Fs, d PkgDl, st SourceType, dest string) (Download, error) {
	if st != Source {
		panic("cannot download other than source for now")
	}
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
			Type:     st,
			Path:     dest,
			New:      false,
			Metadata: d,
		}, nil
	}
	pkgdl := fmt.Sprintf("%s/src/contrib/%s", strings.TrimSuffix(d.Repo.URL, "/"), filepath.Base(dest))
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
		Type:     st,
		Path:     dest,
		New:      true,
		Metadata: d,
	}, nil
}
