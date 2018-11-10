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

// GetPackages returns all packages and the repo that they
// will be acquired from, as well as any missing packages
func (p *PkgDb) GetPackages(pkgs []string) AvailablePkgs {
	ap := AvailablePkgs{}
	for _, pkg := range pkgs {
		found := false
		for _, db := range p.Db {
			if pkgExists(pkg, db.Db) {
				ap.Packages = append(ap.Packages, PkgDl{
					Package: db.Db[pkg],
					Repo:    db.Repo,
				})
				found = true
				break
			}
		}
		if !found {
			ap.Missing = append(ap.Missing, pkg)
		}
	}
	return ap
}

// CheckAllAvailable returns whether all requested packages
// are available in the package database. It is a simple wrapper
// around GetPackages
func (p *PkgDb) CheckAllAvailable(pkgs []string) bool {
	ap := p.GetPackages(pkgs)
	if len(ap.Missing) > 0 {
		return false
	}
	return true
}

// GetAllPkgsByName returns all packages in the database
func (p *PkgDb) GetAllPkgsByName() []string {
	// use map so will remove duplicate packages
	pkgMap := make(map[string]bool)
	for _, db := range p.Db {
		for pkg := range db.Db {
			pkgMap[pkg] = true
		}
	}
	pkgs := []string{}
	for pkg := range pkgMap {
		pkgs = append(pkgs, pkg)
	}
	return pkgs
}
