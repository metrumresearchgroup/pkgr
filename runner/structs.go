package runner

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
