package configlib

// PkgConfig provides information about custom settings during package installation
type PkgConfig struct {
	Suggests bool
	Env      map[string]string
	Repo     string
	Type     string
}

// RepoConfig provides information about custom repository settings
type RepoConfig struct {
	//Suggests bool
	Type string
}

// LogConfig stores information for logging purposes
type LogConfig struct {
	All     string
	Install string
	Level   string
	Overwrite bool
}

// Customizations contains various custom configurations
type Customizations struct {
	Packages map[string]PkgConfig
	Repos    map[string]RepoConfig
}

// PkgrConfig provides a struct for all pkgr related configuration
type PkgrConfig struct {
	Version        int
	Packages       []string
	Suggests       bool
	Repos          []map[string]string
	Library        string
	LibPaths       []string
	Customizations Customizations
	Threads        int
	RPath          string
	Cache          string
	Logging        LogConfig
}
