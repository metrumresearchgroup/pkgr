package gpsr

// NewDefaultInstallDeps provides the default install deps
// Depends/Imports/LinkingTo, not suggests
func NewDefaultInstallDeps() InstallDeps {
	return InstallDeps{
		Depends:   true,
		Imports:   true,
		LinkingTo: true,
		Suggests:  false,
	}
}
