package configlib

// PkgConfig provides information about custom settings during package installation
type PkgConfig struct {
	Suggests bool              `yaml:"Suggests,omitempty"`
	Env      map[string]string `yaml:"Env,omitempty"`
	Repo     string            `yaml:"Repo,omitempty"`
	Type     string            `yaml:"Type,omitempty"`
}

// RepoConfig provides information about custom repository settings
type RepoConfig struct {
	//Suggests bool
	Type string `yaml:"Type,omitempty"`
}

// LogConfig stores information for logging purposes
type LogConfig struct {
	All       string `yaml:"All,omitempty"`
	Install   string `yaml:"Install,omitempty"`
	Level     string `yaml:"Level,omitempty"`
	Overwrite bool   `yaml:"Overwrite,omitempty"`
}

// Customizations contains various custom configurations
type Customizations struct {
	Packages map[string]PkgConfig  `yaml:"Packages,omitempty"`
	Repos    map[string]RepoConfig `yaml:"Repos,omitempty"`
}

// PkgrConfig provides a struct for all pkgr related configuration
type PkgrConfig struct {
	Version        int                 `yaml:"Version,omitempty"`
	Packages       []string            `yaml:"Packages,omitempty"`
	Suggests       bool                `yaml:"Suggests,omitempty"`
	Repos          []map[string]string `yaml:"Repos,omitempty"`
	Library        string              `yaml:"Library,omitempty"`
	LibPaths       []string            `yaml:"LibPaths,omitempty"`
	Customizations Customizations      `yaml:"Customizations,omitempty"`
	Threads        int                 `yaml:"Threads,omitempty"`
	RPath          string              `yaml:"RPath,omitempty"`
	Cache          string              `yaml:"Cache,omitempty"`
	Logging        LogConfig           `yaml:"Logging,omitempty"`
	Update         bool                `yaml:"Update,omitempty"`
}
