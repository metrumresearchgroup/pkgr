package cran

import (
	"errors"
	"fmt"

	"github.com/dpastoor/rpackagemanager/desc"
)

// NewPkgDb returns a new package database
func NewPkgDb(urls []RepoURL) (*PkgDb, error) {
	db := PkgDb{Config: make(map[string]RepoURL)}
	if len(urls) == 0 {
		return &db, errors.New("Package database must contain at least one RepoUrl")
	}
	for _, url := range urls {
		rdb, err := NewRepoDb(url)
		if err != nil {
			return &db, err
		}
		db.Db = append(db.Db, rdb)
	}
	return &db, nil
}

// SetPackageRepo sets a package repository so querying the package will
// pull from that repo
func (p *PkgDb) SetPackageRepo(pkg string, repo string) error {
	for _, r := range p.Db {
		if r.Repo.Name == repo {
			p.Config[pkg] = r.Repo
			return nil
		}
	}
	return fmt.Errorf("no repo: %s, detected containing package: %s", repo, pkg)
}

func pkgExists(pkg string, db map[string]desc.Desc) bool {
	_, exists := db[pkg]
	return exists
}

func isCorrectRepo(pkg string, r RepoURL, cfg map[string]RepoURL) bool {
	repo, exists := cfg[pkg]
	if exists {
		if repo.Name == r.Name {
			return true
		} else {
			return false
		}
	}
	return true
}

// GetPackage gets a package from the package database, returning the first match
func (p *PkgDb) GetPackage(pkg string) (desc.Desc, RepoURL, bool) {
	for _, db := range p.Db {
		if pkgExists(pkg, db.Db) && isCorrectRepo(pkg, db.Repo, p.Config) {
			return db.Db[pkg], db.Repo, true
		}
	}
	return desc.Desc{}, RepoURL{}, false
}

// GetPackageFromRepo gets a package from a repo in the package database
func (p *PkgDb) GetPackageFromRepo(pkg string, repo string) (desc.Desc, RepoURL, bool) {
	for _, db := range p.Db {
		if repo != "" && db.Repo.Name != repo {
			continue
		}
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
		pd, repo, found := p.GetPackage(pkg)
		ap.Packages = append(ap.Packages, PkgDl{
			Package: pd,
			Repo:    repo,
		})
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
