package cran

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/dpastoor/goutils"
	"github.com/dpastoor/rpackagemanager/desc"
	"github.com/spf13/afero"
)

// DownloadPackages downloads a set of packages concurrently
func DownloadPackages(fs afero.OsFs, ds []desc.Desc, url string, dir string) error {
	// bound to
	sem := make(chan struct{}, 10)
	wg := sync.WaitGroup{}
	for _, d := range ds {
		wg.Add(1)
		go func(d desc.Desc) {
			<-sem
			fmt.Println("downloading package ", d.Package)
			DownloadPackage(fs, d, url, filepath.Join(dir, fmt.Sprintf("%s_%s.tar.gz", d.Package, d.Version)))
			sem <- struct{}{}
			wg.Done()
		}(d)
	}
	wg.Done()
	fmt.Println("all packages downloaded")
	return nil
}

func DownloadPackage(fs afero.OsFs, d desc.Desc, url string, dest string) error {

	ok, err := goutils.Exists(fs, dest)
	if err != nil {
		return err
	}
	if ok {
		fmt.Println("already have ", d.Package, " downloaded")
		return nil
	}
	pkgdl := filepath.Clean(fmt.Sprintf("%s/src/contrib/%s_%s.tar.gz", url, d.Package, d.Version))
	fmt.Println("downloading from: ", pkgdl)
	resp, err := http.Get(pkgdl)
	if err != nil {
		fmt.Println("error downloading package", d.Package)
		return err
	}
	defer resp.Body.Close()
	file, err := fs.Open(dest)
	if err != nil {
		fmt.Println("couldn't create tarball for ", d.Package)
		return err
	}
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
