package gpsr

//InstallPlan provides metadata around an installation plan
type InstallPlan struct {
	StartingPackages []string
	DepDb            map[string][]string
}

// InstallDeps contains which dependencies should be installed
type InstallDeps struct {
	Depends   bool
	Imports   bool
	Suggests  bool
	LinkingTo bool
}
