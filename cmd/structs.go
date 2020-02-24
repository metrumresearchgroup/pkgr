package cmd

//// Load report struct
type LoadReport struct {
	RMetadata   RSessionMetadata
	LoadResults map[string]LoadResult
	Failures    int
}

func InitLoadReport(rMetadata RSessionMetadata) LoadReport {
	return LoadReport{
		RMetadata:   rMetadata,
		LoadResults: make(map[string]LoadResult),
		Failures:    0,
	}
}

func (report *LoadReport) AddResult(pkg string, result LoadResult) {
	report.LoadResults[pkg] = result
	if !result.Success {
		report.Failures = report.Failures + 1
	}
}


//// Load result struct
type LoadResult struct {
	Package string
	Version string
	Path    string
	Exiterr string // could be equivalent to exit code.
	//ExitCode int
	Stdout  string
	Stderr  string
	Success bool
	// Can store information for JSON here
}

type pkgLoadMetadata struct { // Used to help create LoadResult
	pkgPath string
	pkgVersion string
}

func MakeLoadResult(pkg, version, path, outStr, errStr string, success bool, exitErr error) LoadResult {
	exitErrString := ""
	if exitErr != nil {
		exitErrString = exitErr.Error()
	}

	return LoadResult{
		Package: pkg,
		Version: version,
		Path:    path,
		Exiterr: exitErrString,
		Stdout:  outStr,
		Stderr:  errStr,
		Success: success,
	}
}

//// R Session Info Struct
type RSessionMetadata struct {
	LibPaths []string
	RPath    string
	RVersion string
}
