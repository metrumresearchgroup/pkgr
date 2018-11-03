package cran

import (
	"errors"

	"github.com/dpastoor/rpackagemanager/desc"
)

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

func pkgExists(pkg string, db map[string]desc.Desc) bool {
	_, exists := db[pkg]
	return exists
}

// GetPackage gets a package from the package database, returning the first match
func (p *PkgDb) GetPackage(pkg string) (desc.Desc, RepoURL, bool) {
	for _, db := range p.Db {
		if pkgExists(pkg, db.Db) {
			return db.Db[pkg], db.Repo, true
		}
	}
	return desc.Desc{}, RepoURL{}, false
}
