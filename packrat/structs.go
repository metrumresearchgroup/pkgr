package packrat

// PackageReqs stores the package requirements per packrat
type PackageReqs struct {
	Package  string
	Source   string
	Version  string
	Hash     string
	Requires []string
}
