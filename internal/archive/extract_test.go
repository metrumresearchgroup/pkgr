package archive

import (
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/klauspost/compress/zstd"
	"github.com/ulikunitz/xz"
)

func compress(t *testing.T, r io.Reader, format string, filename string) {
	w, err := os.Create(filename)
	if err != nil {
		t.Fatal(err)
	}

	var wc io.WriteCloser
	switch format {
	case "gz":
		wc = gzip.NewWriter(w)
	case "xz":
		wc, err = xz.NewWriter(w)
		if err != nil {
			t.Fatal(err)
		}
	case "zst":
		wc, err = zstd.NewWriter(w)
		if err != nil {
			t.Fatal(err)
		}
	default:
		t.Fatalf("unknown format: %s", format)
	}

	if _, err := io.Copy(wc, r); err != nil {
		t.Fatal(err)
	}
	if err = wc.Close(); err != nil {
		t.Fatal(err)
	}
	if err = w.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestExtractCompressed(t *testing.T) {
	r, err := os.Open("testdata/nested.tar")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()

	dir := t.TempDir()
	formats := []string{
		"gz",
		"xz",
		"zst",
	}
	for _, format := range formats {
		t.Run(format, func(t *testing.T) {
			dirOut := filepath.Join(dir, "out-"+format)
			filename := filepath.Join(dir, "nested.tar."+format)
			_, err := r.Seek(0, 0)
			if err != nil {
				t.Fatal(err)
			}
			compress(t, r, format, filename)

			if err = Extract(filename, dirOut); err != nil {
				t.Fatal(err)
			}

			subpaths := map[string]string{
				"foo":             "A\n",
				"sub/bar":         "B\n",
				"sub/subsub/baz":  "C\n",
				"sub/subsub/baz2": "D\n",
			}

			for subpath, want := range subpaths {
				bs, err := os.ReadFile(filepath.Join(dirOut, subpath))
				if err != nil {
					t.Fatal(err)
				}
				got := string(bs)
				if got != want {
					t.Errorf("unexpected content for %q: got %q, want %q", subpath, got, want)
				}
			}
		})
	}
}

func TestExtractPermissions(t *testing.T) {
	dir := t.TempDir()
	err := Extract("testdata/perm.tar", dir)
	if err != nil {
		t.Fatal(err)
	}

	fi, err := os.Stat(filepath.Join(dir, "foo"))
	if err != nil {
		t.Fatal(err)
	}
	mode := fi.Mode()
	if !mode.IsRegular() || mode&0o700 != 0o600 {
		t.Errorf("foo should be regular file with user read and write bits: %v", mode)
	}

	fi, err = os.Stat(filepath.Join(dir, "bar"))
	if err != nil {
		t.Fatal(err)
	}
	mode = fi.Mode()
	if !mode.IsRegular() || mode&0o700 != 0o700 {
		t.Errorf("bar should be regular file with user read, write, and execute bits: %v", mode)
	}
}

func TestExtractClearsTarget(t *testing.T) {
	dir := t.TempDir()
	fname := filepath.Join(dir, "preexisting")
	w, err := os.Create(fname)
	if err != nil {
		t.Fatal(err)
	}
	if err = w.Close(); err != nil {
		t.Fatal(err)
	}

	err = Extract("testdata/perm.tar", dir)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(fname); err == nil {
		t.Errorf("file should not exist: %q", fname)
	}
}

func TestExtractCreatesImplicitSubdir(t *testing.T) {
	dir := t.TempDir()
	err := Extract("testdata/subdir.tar", dir)
	if err != nil {
		t.Fatal(err)
	}

	// sub/ is created despite not having an explicit entry in subdir.tar.
	fname := filepath.Join(dir, "sub", "foo")
	if _, err := os.Stat(fname); err != nil {
		t.Errorf("file should exist: %q", fname)
	}
}

func TestExtractAbortUnsupported(t *testing.T) {
	dir := t.TempDir()

	tests := []struct {
		filename       string
		errorSubstring string
	}{
		{
			filename:       "abs.tar",
			errorSubstring: "path not local",
		},
		{
			filename:       "rel.tar",
			errorSubstring: "path not local",
		},
		{
			filename:       "symlink.tar",
			errorSubstring: "unsupported mode",
		},
		{
			filename:       "symlink-outside.tar",
			errorSubstring: "unsupported mode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			err := Extract(filepath.Join("testdata", tt.filename), dir)
			if err == nil {
				t.Fatalf("call with %s should error", tt.filename)
			}
			if !strings.Contains(err.Error(), tt.errorSubstring) {
				t.Errorf("expected error to contain %s: %v", tt.errorSubstring, err)
			}
		})
	}
}

func TestExtractUnknownFormat(t *testing.T) {
	err := Extract("foo", "")
	if err == nil {
		t.Fatal("call with file without extension should error")
	}
	if !strings.Contains(err.Error(), "identify format") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestExtractBadFormat(t *testing.T) {
	dir := t.TempDir()

	r, err := os.Open("testdata/nested.tar")
	if err != nil {
		t.Fatal(err)
	}
	fname := filepath.Join(dir, "nested.tar.gz")
	w, err := os.Create(fname)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.Copy(w, r); err != nil {
		t.Fatal(err)
	}
	if err = r.Close(); err != nil {
		t.Fatal(err)
	}
	if err = w.Close(); err != nil {
		t.Fatal(err)
	}

	if err := Extract(fname, filepath.Join(dir, "out")); err == nil {
		t.Fatalf("expected error when uncompressed input has .gz extension")
	}
}
