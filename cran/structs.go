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
	DescriptionsBySourceType map[SourceType]map[string]desc.Desc
	Time                     time.Time
	Repo                     RepoURL
	DefaultSourceType        SourceType
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

// PkgNexus represents a package database
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
	Size     int64
}

func (d Download) GetMegabytes() float64 {
	return float64(d.Size) / (1024 * 1024)
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

type OutdatedPackage struct {
	Package    string
	OldVersion string
	NewVersion string
}

type OsRelease struct {
	Name            string `mapstructure:"NAME"`
	Version         string `mapstructure:"VERSION"`
	Id              string `mapstructure:"ID"`
	IdLike          string `mapstructure:"ID_LIKE"`
	LtsRelease      string
	PrettyName      string `mapstructure:"PRETTY_NAME"`
	VersionId       string `mapstructure:"VERSION_ID"`
	VersionCodename string `mapstructure:"VERSION_CODENAME"`
	UbuntuCodename  string `mapstructure:"UBUNTU_CODENAME"`
}

