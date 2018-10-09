package packrat

import (
	"bytes"
	"fmt"
)

// ChunkLockfile breaks a packrat lockfile into chunks
func ChunkLockfile(b []byte) *LockFileDb {
	lf := NewLockFileDb()
	cb := bytes.Split(b, []byte("\n\n"))
	for _, p := range cb {
		p = CollapseIndentation(p)
		if bytes.HasPrefix(p, []byte("PackratFormat")) {
			fmt.Printf("%s\n", p)
		} else if bytes.Contains(p, []byte("GithubRepo")) {
			gpkg := ParsePackageReqsGH(p)
			lf.Github[gpkg.Reqs.Package] = gpkg
		} else {
			pkg := ParsePackageReqs(p)
			lf.CRANlike[pkg.Package] = pkg
		}
	}
	return lf
}
