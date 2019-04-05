package cran

import (
	"time"

	"github.com/metrumresearchgroup/pkgr/desc"
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
	Dbs               map[SourceType]map[string]desc.Desc
	Time              time.Time
	Repo              RepoURL
	DefaultSourceType SourceType
}

// InstallConfig contains custom settings for a full install
type InstallConfig struct {
	Packages map[string]PkgConfig
	Repos    map[string]RepoConfig
}

// RepoConfig contains settings for a repo
type RepoConfig struct {
	DefaultSourceType SourceType
}

//PkgConfig stores configuration information about a given package
type PkgConfig struct {
	Repo RepoURL
	Type SourceType
}

// PkgNexus represents a sort of phone book of all available repositories and packages in those repositories
type PkgNexus struct {
	Db                []*RepoDb
	Config            *InstallConfig
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

// RVersion contains information about the R version
type RVersion struct {
	Major int
	Minor int
	Patch int
}
