package rpkg

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	Log "github.com/metrumresearchgroup/pkgr/logger"
	"github.com/spf13/afero"
)

// Hash a tarball
func Hash(fs afero.Fs, tbp string) (string, error) {
	f, err := fs.Open(tbp)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		Log.Log.Fatal(err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), err
}
