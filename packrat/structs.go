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

// PackratMetadata contains the metadata about the packrat lock file
type PackratMetadata struct {
	Format   float32
	Version  string
	RVersion string
	Repos    []string
}
type LockFile struct {
	Metadata PackratMetadata
	CRANlike []PackageReqs
	Github   []GithubPackageReqs
}
