package packrat

import (
	"bytes"
	"fmt"
)

// ChunkLockfile breaks a packrat lockfile into chunks
func ChunkLockfile(b []byte) LockFile {
	lf := LockFile{}
	cb := bytes.Split(b, []byte("\n\n"))
	for _, p := range cb {
		p = CollapseIndentation(p)
		if bytes.HasPrefix(p, []byte("PackratFormat")) {
			fmt.Printf("%s\n", p)
		} else if bytes.Contains(p, []byte("GithubRepo")) {
			lf.Github = append(lf.Github, ParsePackageReqsGH(p))
		} else {
			lf.CRANlike = append(lf.CRANlike, ParsePackageReqs(p))
		}
	}
	return lf
}
