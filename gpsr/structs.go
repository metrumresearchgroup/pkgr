package gpsr

//InstallPlan provides metadata around an installation plan
type InstallPlan struct {
	StartingPackages []string
	DepDb            map[string][]string
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
