package cran

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dpastoor/rpackagemanager/desc"
)

// NewRepoDb returns a new Repo database
func NewRepoDb(url RepoURL) (*RepoDb, error) {
	ddb := &RepoDb{Db: make(map[string]desc.Desc), Time: time.Now(), Repo: url}
	return ddb, ddb.GetPackages()
}

// Decode decodes the package database
func (r *RepoDb) Decode(file string) error {
	f, err := os.Open(file)
	if err != nil {
		fmt.Println("problem opening crandb", file)
		return err
	}
	d := gob.NewDecoder(f)
	return d.Decode(&r.Db)
}

// Encode encodes the PackageDatabase
func (r *RepoDb) Encode(file string) error {
	err := os.MkdirAll(filepath.Dir(file), 0777)
	if err != nil {
		return err
	}
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	e := gob.NewEncoder(f)

	// Encoding the map
	err = e.Encode(r.Db)
	if err != nil {
		return err
	}
	return nil
}

// GetPackages gets the packages for  RepoDb
// R_AVAILABLE_PACKAGES_CACHE_CONTROL_MAX_AGE controls the timing to requery the cache in R
func (r *RepoDb) GetPackages() error {
	// just get src versions for now
	pkgURL := strings.TrimSuffix(r.Repo.URL, "/") + "/src/contrib/PACKAGES"
	cdir, err := os.UserCacheDir()
	if err != nil {
		fmt.Println("could not use user cache dir, using temp dir")
		cdir = os.TempDir()
	}
	h := md5.New()
	io.WriteString(h, r.Repo.URL+r.Repo.Name)
	pkgHash := fmt.Sprintf("%x", h.Sum(nil))
	pkgFile := filepath.Join(cdir, "r_package_caches", pkgHash)
	if _, err := os.Stat(pkgFile); !os.IsNotExist(err) {
		return r.Decode(pkgFile)
	}

	res, err := http.Get(pkgURL)
	if err != nil {
		return fmt.Errorf("problem getting packages from url %s", pkgURL)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	cb := bytes.Split(body, []byte("\n\n"))
	fmt.Println("fetched, decoding pkgs...")
	for _, p := range cb {
		reader := bytes.NewReader(p)
		d, err := desc.ParseDesc(reader)
		r.Db[d.Package] = d
		if err != nil {
			fmt.Println("problem parsing package with info ", string(p))
			panic(err)
		}
	}
	fmt.Println("encoding...")
	err = r.Encode(pkgFile)
	fmt.Println("done encoding...")
	return err
}
