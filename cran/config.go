package cran

// NewPkgConfigDB initializes a PkgConfig map
func NewInstallConfig() *InstallConfig {
	return &InstallConfig{
		Packages: make(map[string]PkgConfig),
		Repos:    make(map[string]RepoConfig),
	}
}
