package packrat

import "github.com/dpastoor/rpackagemanager/gpsr"

// NewLockFileDb initializes a new lockfile
func NewLockFileDb() *LockFileDb {
	return &LockFileDb{
		CRANlike: make(map[string]PackageReqs),
		Github:   make(map[string]GithubPackageReqs),
	}
}

// GetPackage returns an interface for either cran or github package
// that can be used downstream with a type switch
func (ldb LockFileDb) GetPackage(pkg string) (bool, interface{}) {

	p, ok := ldb.CRANlike[pkg]
	if ok {
		return true, p
	}
	pg, ok := ldb.Github[pkg]
	if ok {
		return true, pg
	}
	return false, nil
}

func (ldb LockFileDb) GetPackageReqs(pkg string) (bool, PackageReqs) {

	p, ok := ldb.CRANlike[pkg]
	if ok {
		return true, p
	}
	pg, ok := ldb.Github[pkg]
	if ok {
		return true, pg.Reqs
	}
	return false, PackageReqs{}
}

// SolveLockfile provides a solution give a lockfile spec
func SolveLockfile(ldb LockFileDb) ([][]string, error) {
	workingGraph := gpsr.NewGraph()

	for _, p := range ldb.CRANlike {
		workingGraph[p.Package] = gpsr.NewNode(p.Package, p.Requires)
	}
	for _, p := range ldb.Github {
		workingGraph[p.Reqs.Package] = gpsr.NewNode(p.Reqs.Package, p.Reqs.Requires)
	}

	resolved, err := gpsr.ResolveGraph(workingGraph)

	return resolved, err
}
