package rcmd

// ExecSettings controls settings related to R execution
type ExecSettings struct {
	WorkDir string
}

// RSettings controls settings related to managing libraries
type RSettings struct {
	LibPaths []string          `json:"lib_paths,omitempty"`
	Rpath    string            `json:"rpath,omitempty"`
	EnvVars  map[string]string `json:"env_vars,omitempty"`
}

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
