package cran

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/dpastoor/goutils"
	"github.com/dpastoor/rpackagemanager/desc"
	"github.com/spf13/afero"
)

// DownloadPackage downloads a package
func DownloadPackage(fs afero.OsFs, d desc.Desc, url string, dest string) error {

	ok, err := goutils.Exists(fs, dest)
	if err != nil {
		return err
	}
	if ok {
		fmt.Println("already have ", d.Package, " downloaded")
		return nil
	}
	resp, err := http.Get(filepath.Clean(fmt.Sprintf("%s/src/contrib/%s_%s.tar.gz", url, d.Package, d.Version)))
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
