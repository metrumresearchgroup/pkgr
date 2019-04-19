package gpsr

import "github.com/metrumresearchgroup/pkgr/cran"

//InstallPlan provides metadata around an installation plan
type InstallPlan struct {
	StartingPackages []string
	DepDb            map[string][]string
	PackageDownloads []cran.PkgDl
	OutdatedPackages []OutdatedPackage
}

// PkgDeps contains which dependencies should be installed
// for a particular package
type PkgDeps struct {
	Depends   bool
	Imports   bool
	Suggests  bool
	LinkingTo bool
}

// InstallDeps contains the information about dependencies to be installed
type InstallDeps struct {
	Deps    map[string]PkgDeps
	Default PkgDeps
}


type OutdatedPackage struct {
	Package string
	OldVersion string
	NewVersion string
}

