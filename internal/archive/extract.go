package archive

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/mholt/archives"
)

// Extract writes each file from the archive at filename beneath the directory
// dir.  If the directory already exists, it is removed before extracting the
// archive.
//
// The archive can be in any format supported by the archives package, which
// covers the types of compression supported by 'R CMD build' (as of R 4.4.3).
//
// This function is for extracting R package source tarballs rather than
// general-purpose extraction.
//
//   - Only regular files and directories are supported.  An error is returned if
//     entries for any other mode are encountered, including a symbolic link.
//
//     Supporting symbolic links is not necessary because, even if a package's
//     source tree has a symbolic link, 'R CMD build' dereferences the link when
//     preparing the tarball.
//
//   - An error is returned if any entry creates a file that points outside of
//     the target directory (e.g., a top-level entry for "../foo" or an entry
//     for an absolute path).
//
//   - The full permission bits from the archive entries are not mirrored.
//     Regular files are created with mode 0o666 or 0o777 (before umask),
//     depending on whether the entry in the archive has execute bits.
//     Directories are created with mode 0o777 (before umask).
func Extract(filename string, dir string) error {
	format, _, err := archives.Identify(context.Background(), filename, nil)
	if err != nil {
		return fmt.Errorf("unable to identify format: %q: %w", filename, err)
	}
	ex, ok := format.(archives.Extractor)
	if !ok {
		return fmt.Errorf("no extractor found for %q", filename)
	}

	r, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer r.Close()

	// Clear any existing destination.  In the context of R source tarballs, the
	// set of files in the directory should match the set in the archive, with
	// no extra files.
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o777); err != nil {
		return err
	}
	root, err := os.OpenRoot(dir)
	if err != nil {
		return err
	}
	defer root.Close()

	return ex.Extract(context.Background(), r, func(ctx context.Context, fi archives.FileInfo) error {
		return extractFile(ctx, fi, root)
	})
}

func extractFile(ctx context.Context, fi archives.FileInfo, root *os.Root) error {
	// root is responsible for ensuring that created files don't escape the
	// target directory, but this guard avoids calling MkdirAll with a path that
	// escapes root (e.g., if NameInArchive starts with "../").
	if !filepath.IsLocal(fi.NameInArchive) {
		return fmt.Errorf("path not local to archive: %q", fi.NameInArchive)
	}

	filename := filepath.Join(root.Name(), fi.NameInArchive)
	mode := fi.Mode()
	switch {
	case mode.IsDir():
		return os.MkdirAll(filename, 0o777)
	case mode.IsRegular():
		// Don't assume that an entry for each parent directory was already
		// encountered or that an entry even exists (issue 259).
		if err := os.MkdirAll(filepath.Dir(filename), 0o777); err != nil {
			return err
		}

		r, err := fi.Open()
		if err != nil {
			return err
		}
		defer r.Close()

		w, err := root.OpenFile(fi.NameInArchive, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o666|mode.Perm())
		if err != nil {
			w.Close()
			return err
		}

		if _, err := io.Copy(w, r); err != nil {
			w.Close()
			return err
		}

		return w.Close()
	default:
		return fmt.Errorf("unsupported mode: %v", mode)
	}
}
