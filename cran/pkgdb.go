package cran

import (
	"errors"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/metrumresearchgroup/pkgr/desc"
)

// NewPkgDb returns a new package database
func NewPkgDb(urls []RepoURL, dst SourceType, cfgdb *InstallConfig, rv RVersion) (*PkgDb, error) {
	db := PkgDb{
		Config:            cfgdb,
		DefaultSourceType: dst,
	}
	if len(urls) == 0 {
		return &db, errors.New("Package database must contain at least one RepoUrl")
	}
	type rd struct {
		Url       RepoURL
		Rdb       *RepoDb
		repoIndex int
		Err       error
	}
	rdbc := make(chan rd, len(urls))
	defer close(rdbc)
	ri := 0
	for _, url := range urls {
		db.Db = append(db.Db, nil)
		go func(url RepoURL, dst SourceType, ri int) {
			rdb, err := NewRepoDb(url, dst, cfgdb.Repos[url.Name], rv)
			rdbc <- rd{url, rdb, ri, err}
		}(url, dst, ri)
		ri++
	}
	var err error
	err = nil
	for i := 0; i < len(urls); i++ {
		result := <-rdbc
		if result.Err != nil {
			//
			log.Fatalf("error downloading repo information from: %s:%s\n", result.Url.Name, result.Url.URL)
			err = result.Err
		} else {
			db.Db[result.repoIndex] = result.Rdb
		}
	}
	return &db, err
}

// SetPackageRepo sets a package repository so querying the package will
// pull from that repo
func (p *PkgDb) SetPackageRepo(pkg string, repo string) error {
	for _, r := range p.Db {
		if r.Repo.Name == repo {
			cfg := p.Config.Packages[pkg]
			cfg.Repo = r.Repo
			p.Config.Packages[pkg] = cfg
			return nil
		}
	}
	return fmt.Errorf("no repo: %s, detected containing package: %s", repo, pkg)
}

// SetPackageType sets the package type (source/binary) for installation
func (p *PkgDb) SetPackageType(pkg string, t string) error {
	cfg := p.Config.Packages[pkg]
	if strings.EqualFold(t, "source") {
		cfg.Type = Source
	} else if strings.EqualFold(t, "binary") {
		cfg.Type = Binary
	} else {
		return fmt.Errorf("invalid source type: %s for pkg: %s", t, pkg)
	}
	p.Config.Packages[pkg] = cfg
	return nil
}

func pkgExists(pkg string, db map[string]desc.Desc) bool {
	_, exists := db[pkg]
	return exists
}
func pkgExistsInRepo(pkg string, dbs map[SourceType]map[string]desc.Desc) bool {
	exists := false
	for _, db := range dbs {
		_, exists = db[pkg]
		if exists {
			return exists
		}
	}
	return exists
}

func isCorrectRepo(pkg string, r RepoURL, cfg map[string]PkgConfig) bool {
	pkgcfg, exists := cfg[pkg]
	if exists && pkgcfg.Repo.Name != "" {
		if pkgcfg.Repo.Name == r.Name {
			return true
		} else {
			return false
		}
	}
	return true
}

// GetPackage gets a package from the package database, returning the first match
func (p *PkgDb) GetPackage(pkg string) (desc.Desc, PkgConfig, bool) {
	cfg, exists := p.Config.Packages[pkg]
	st := p.DefaultSourceType
	if exists && cfg.Type != Default {
		st = cfg.Type
	}
	for _, db := range p.Db {
		rst := st
		if db.DefaultSourceType != rst && !exists && db.DefaultSourceType != Default {
			rst = db.DefaultSourceType
		}
		// For now package existence is checked exactly as the package is specified
		// in the config. Eg, if specifies binary, will only check binary version
		// the checking if also exists as source or otherwise should occur upstream
		// then be set as part of the explicit configuration.
		if pkgExists(pkg, db.Dbs[rst]) && isCorrectRepo(pkg, db.Repo, p.Config.Packages) {
			return db.Dbs[rst][pkg], PkgConfig{Repo: db.Repo, Type: rst}, true
		}
	}
	return desc.Desc{}, PkgConfig{}, false
}

// GetPackageFromRepo gets a package from a repo in the package database
func (p *PkgDb) GetPackageFromRepo(pkg string, repo string) (desc.Desc, PkgConfig, bool) {
	st := p.Config.Packages[pkg].Type
	for _, db := range p.Db {
		if repo != "" && db.Repo.Name != repo {
			continue
		}
		if pkgExists(pkg, db.Dbs[st]) {
			return db.Dbs[st][pkg], PkgConfig{Repo: db.Repo, Type: st}, true
		}
	}
	return desc.Desc{}, PkgConfig{}, false
}

// GetPackages returns all packages and the repo that they
// will be acquired from, as well as any missing packages
func (p *PkgDb) GetPackages(pkgs []string) AvailablePkgs {
	ap := AvailablePkgs{}
	for _, pkg := range pkgs {
		pd, cfg, found := p.GetPackage(pkg)
		ap.Packages = append(ap.Packages, PkgDl{
			Package: pd,
			Config:  cfg,
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
		for st := range db.Dbs {
			for pkg := range db.Dbs[st] {
				pkgMap[pkg] = true
			}
		}
	}
	pkgs := []string{}
	for pkg := range pkgMap {
		pkgs = append(pkgs, pkg)
	}
	return pkgs
}
