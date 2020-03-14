package gpsr

// NewDefaultInstallDeps provides the default install deps
// Depends/Imports/LinkingTo, not suggests
func NewDefaultInstallDeps() InstallDeps {
	return InstallDeps{
		Deps: make(map[string]PkgDeps),
		Default: PkgDeps{
			Depends:   true,
			Imports:   true,
			LinkingTo: true,
			Suggests:  false,
			NoRecommended: false,
		}}
}

// AllPkgDeps returns PkgDeps with all set to true
func AllPkgDeps() PkgDeps {
	return PkgDeps{
		Depends:   true,
		Imports:   true,
		LinkingTo: true,
		Suggests:  true,
	}
}
