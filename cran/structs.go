package cran

import (
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
	Dbs  map[SourceType]map[string]desc.Desc
	Time time.Time
	Repo RepoURL
}

//PkgConfig stores configuration information about a given package
type PkgConfig struct {
	Repo RepoURL
	Type SourceType
}

// PkgDb represents a package database
type PkgDb struct {
	Db                []*RepoDb
	Config            map[string]PkgConfig
	DefaultSourceType SourceType
}

// Download provides information about the package download
type Download struct {
	Path     string
	New      bool
	Metadata PkgDl
}

// PkgDl holds the metadata needed to download a package
type PkgDl struct {
	Config  PkgConfig
	Package desc.Desc
}

// AvailablePkgs provides information about the packages available in
// the package database from a set of requested packages
type AvailablePkgs struct {
	Packages []PkgDl
	Missing  []string
}
