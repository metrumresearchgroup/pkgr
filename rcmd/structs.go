package rcmd

// CmdResult stores information about the executed cmd
type CmdResult struct {
	Stdout   string `json:"stdout,omitempty"`
	Stderr   string `json:"stderr,omitempty"`
	ExitCode int    `json:"exit_code,omitempty"`
}

// ExecSettings controls settings related to R execution
type ExecSettings struct {
	WorkDir string `json:"work_dir,omitempty"`
}

// RSettings controls settings related to managing libraries
type RSettings struct {
	LibPaths []string          `json:"lib_paths,omitempty"`
	Rpath    string            `json:"rpath,omitempty"`
	EnvVars  map[string]string `json:"env_vars,omitempty"`
}

// InstallArgs represents the installation arguments R CMD INSTALL can consume
type InstallArgs struct {
	Clean          bool `rcmd:"clean"`
	Preclean       bool `rcmd:"preclean"`
	Debug          bool `rcmd:"debug"`
	NoConfigure    bool `rcmd:"no-configure"`
	Example        bool `rcmd:"example"`
	Fake           bool `rcmd:"fake"`
	Build          bool `rcmd:"build"`
	InstallTests   bool `rcmd:"install-tests"`
	NoMultiarch    bool `rcmd:"no-multiarch"`
	WithKeepSource bool `rcmd:"with-keep.source"`
	ByteCompile    bool `rcmd:"byte-compile"`
	NoTestLoad     bool `rcmd:"no-test-load"`
	NoCleanOnError bool `rcmd:"no-clean-on-error"`
	//set
	Library string `rcmd:"library=%s,fmt"`
}

// PackageCache provides metadata about the package cache
// Each repository should be a subfolder from the BaseDir
// with separate folders for binary and source packages
type PackageCache struct {
	BaseDir string
}
