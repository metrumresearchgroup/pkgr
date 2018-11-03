package cran

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/dpastoor/rpackagemanager/desc"
)

// RepoURL represents the URL and name for a repo
// to match the R convention of specifying a repository name
// CRAN = https://cran.rstudio.com would be
// RepoUrl{URL: "https://cran.rstudio.com", Name: "CRAN"}
type RepoURL struct {
	URL  string
	Name string
}

// RepoDb represents a Db
type RepoDb struct {
	Db   map[string]desc.Desc
	Time time.Time
	Repo RepoURL
}

// PkgDb represents a package database
type PkgDb struct {
	Db []*RepoDb
}

// NewRepoDb returns a new Repo database
func NewRepoDb(url RepoURL) (*RepoDb, error) {
	ddb := &RepoDb{Db: make(map[string]desc.Desc), Time: time.Now(), Repo: url}
	return ddb, ddb.GetPackages()
}

// NewPkgDb returns a new package database
func NewPkgDb(urls []RepoURL) (*PkgDb, error) {
	db := &PkgDb{}
	if len(urls) == 0 {
		return db, errors.New("Package database must contain at least one RepoUrl")
	}
	for _, url := range urls {
		rdb, err := NewRepoDb(url)
		if err != nil {
			return db, err
		}
		db.Db = append(db.Db, rdb)
	}
	return db, nil
}

// Decode decodes the package database
func (r *RepoDb) Decode(file string) error {
	f, err := os.Open(file)
	if err != nil {
		fmt.Println("problem opening crandb", file)
		return err
	}
	d := gob.NewDecoder(f)
	return d.Decode(r.Db)
}

// Encode encodes the PackageDatabase
func (r *RepoDb) Encode(file string) error {
	err := os.MkdirAll(filepath.Dir(file), 0644)
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

// // GetPackages gets all packages from the specified RepoUrls
// func (p *PkgDb) GetPackages() error {
// 	if len(p.Urls) == 0 {
// 		return errors.New("no RepoUrls to query")
// 	}
// 	for _, url := range p.Urls {

// 	}
// 	return nil
// }

// GetPackages gets the packages for  RepoDb
// R_AVAILABLE_PACKAGES_CACHE_CONTROL_MAX_AGE controls the timing to requery the cache in R
func (r *RepoDb) GetPackages() error {
	// just get src versions for now
	pkgURL := filepath.Join(r.Repo.URL, "src/contrib/PACKAGES")
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
	for _, p := range cb {
		reader := bytes.NewReader(p)
		d, err := desc.ParseDesc(reader)
		r.Db[d.Package] = d
		if err != nil {
			fmt.Println("problem parsing package with info ", string(p))
			panic(err)
		}
		err = r.Encode(pkgFile)
	}
	return err
}
