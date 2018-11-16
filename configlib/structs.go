package configlib

// PkgSettings provides information about custom settings during package installation
type PkgSettings struct {
	Suggests bool
	Env      []map[string]string
	Repo     string
}

// PkgrConfig provides a struct for all pkgr related configuration
type PkgrConfig struct {
	Version        int
	Packages       []string
	Suggests       bool
	Repos          []map[string]string
	Library        string
	LibPaths       []string
	Customizations map[string]PkgSettings
	Threads        int
	RPath          string
}
