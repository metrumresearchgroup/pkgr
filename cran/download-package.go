package cran

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/dpastoor/goutils"
	"github.com/dpastoor/rpackagemanager/desc"
	"github.com/spf13/afero"
)

// DownloadPackages downloads a set of packages concurrently
func DownloadPackages(fs afero.Fs, ds []desc.Desc, url string, dir string) error {
	// bound to
	sem := make(chan struct{}, 10)
	wg := sync.WaitGroup{}
	for _, d := range ds {
		wg.Add(1)
		go func(d desc.Desc, wg *sync.WaitGroup) {
			sem <- struct{}{}
			defer func() {
				<-sem
				wg.Done()
			}()
			pkgdir := filepath.Join(dir, fmt.Sprintf("%s_%s.tar.gz", d.Package, d.Version))
			DownloadPackage(fs, d, url, pkgdir)
		}(d, &wg)
	}
	wg.Wait()
	fmt.Println("all packages downloaded")
	return nil
}

func DownloadPackage(fs afero.Fs, d desc.Desc, url string, dest string) error {

	ok, err := goutils.Exists(fs, dest)
	if err != nil {
		return err
	}
	if ok {
		fmt.Println("already have ", d.Package, " downloaded")
		return nil
	}
	pkgdl := fmt.Sprintf("%s/src/contrib/%s_%s.tar.gz", strings.TrimSuffix(url, "/"), d.Package, d.Version)
	fmt.Println("downloading from: ", pkgdl)
	resp, err := http.Get(pkgdl)
	if err != nil {
		fmt.Println("error downloading package", d.Package)
		return err
	}
	defer resp.Body.Close()
	file, err := fs.Create(dest)
	if err != nil {
		fmt.Println("couldn't create tarball for ", d.Package)
		fmt.Println(err)
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
