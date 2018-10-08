package packrat

// PackageReqs stores the package requirements per packrat
type PackageReqs struct {
	Package  string
	Source   string
	Version  string
	Hash     string
	Requires []string
}

// GithubPackageReqs also identifies the
type GithubPackageReqs struct {
	Reqs           PackageReqs
	GithubRepo     string
	GithubUsername string
	GithubRef      string
	GithubSha1     string
}

// Metadata contains the metadata about the packrat lock file
type Metadata struct {
	Format   float32
	Version  string
	RVersion string
	Repos    []string
}

// LockFileDb contains information from the lockfile
type LockFileDb struct {
	Metadata Metadata
	CRANlike map[string]PackageReqs
	Github   map[string]GithubPackageReqs
}
