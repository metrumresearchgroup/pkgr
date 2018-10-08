package packrat

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
