package rcmd

import (
	"github.com/metrumresearchgroup/pkgr/cran"
	"github.com/spf13/afero"
)

// CmdResult stores information about the executed cmd
type CmdResult struct {
	Stdout   string `json:"stdout,omitempty"`
	Stderr   string `json:"stderr,omitempty"`
	ExitCode int    `json:"exit_code,omitempty"`
}

// ExecSettings controls settings related to R execution
type ExecSettings struct {
	WorkDir     string `json:"work_dir,omitempty"`
	PkgrVersion string `json:"pkgr_version,omitempty"`
}

// RSettings controls settings related to managing libraries
type RSettings struct {
	Version       cran.RVersion                `json:"r_version,omitempty"`
	LibPaths      []string                     `json:"lib_paths,omitempty"`
	Rpath         string                       `json:"rpath,omitempty"`
	GlobalEnvVars map[string]string            `json:"global_env_vars,omitempty"`
	PkgEnvVars    map[string]map[string]string `json:"pkg_env_vars,omitempty"`
	Platform      string                       `json:"platform,omitempty"`
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

// InstallRequest provides information about the installation request
type InstallRequest struct {
	Package      string
	Metadata     cran.Download
	Cache        PackageCache
	InstallArgs  InstallArgs
	ExecSettings ExecSettings
	RSettings    RSettings
}

// InstallUpdate provides information about the Job in the queue
type InstallUpdate struct {
	Result       CmdResult
	Package      string
	BinaryPath   string
	Msg          string
	Err          error
	ShouldUpdate bool
}

// Worker does work
type Worker struct {
	ID          int
	WorkQueue   <-chan InstallRequest
	UpdateQueue chan<- InstallUpdate
	Quit        chan bool
	InstallFunc func(fs afero.Fs,
		ir InstallRequest,
		pc PackageCache) (CmdResult, string, error)
}

// InstallQueue represents a new install queue
type InstallQueue struct {
	WorkQueue   chan InstallRequest
	UpdateQueue chan InstallUpdate
	Workers     []Worker
}
