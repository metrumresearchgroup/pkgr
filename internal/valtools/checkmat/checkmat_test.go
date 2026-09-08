// Copyright 2024 Metrum Research Group
// SPDX-License-Identifier: MIT

package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/goccy/go-yaml"
)

func assertCode(t *testing.T, output, code string, ntimes int) {
	t.Helper()
	found := strings.Count(output, "["+code+"]")
	if found != ntimes {
		t.Errorf("got %dx [%s] in output, want %dx\noutput: %q", found, code, ntimes, output)
	}
}

func TestCheckValidFileNamesBad(t *testing.T) {
	tests := []struct {
		name    string
		entries []entry
		want    int
	}{
		{
			name: "one entry",
			entries: []entry{
				{
					Entrypoint: "foo",
					Doc:        "/foo/bar",
					Tests:      []string{"foo/../bar", "/baz"},
				},
			},
			want: 3,
		},
		{
			name: "multiple entries",
			entries: []entry{
				{
					Entrypoint: "a",
					Doc:        "/foo/bar",
					Tests:      []string{"foo/../bar"},
				},
				{
					Entrypoint: "b",
					Doc:        "./foo/bar",
				},
				{
					Entrypoint: "c",
					Doc:        "good",
				},
			},
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			bad, err := checkValidFileNames(tt.entries, &buf)
			if err != nil {
				t.Fatal(err)
			}

			if bad != tt.want {
				t.Errorf("invalid file names: got %d, want %d", bad, tt.want)
			}
			out := buf.String()
			assertCode(t, out, "01", tt.want)
		})
	}
}

func TestCheckValidFileNamesGood(t *testing.T) {
	entries := []entry{
		{
			Entrypoint: "a",
			Doc:        "foo/bar",
			Tests:      []string{"foo/bar_test.go"},
		},
		{
			Entrypoint: "b",
			Doc:        "foo/bar",
		},
	}

	var buf bytes.Buffer
	bad, err := checkValidFileNames(entries, &buf)
	if err != nil {
		t.Fatal(err)
	}

	if bad != 0 {
		t.Errorf("got %d invalid file names, want none", bad)
	}

	if out := buf.String(); out != "" {
		t.Errorf("got %q for output, want empty", out)
	}
}

func TestCheckMissingFilesBad(t *testing.T) {
	dir := t.TempDir()
	tests := []struct {
		name    string
		entries []entry
		want    int
	}{
		{
			name: "one entry",
			entries: []entry{
				{
					Entrypoint: "foo",
					Code:       "cmd/foo.go",
					Doc:        "docs/commands/foo.md",
					Tests:      []string{"cmd/foo_test.go"},
				},
			},
			want: 3,
		},
		{
			name: "one entry without tests",
			entries: []entry{
				{
					Entrypoint: "foo",
					Code:       "cmd/foo.go",
					Doc:        "docs/commands/foo.md",
				},
			},
			want: 2,
		},
		{
			name: "two entries",
			entries: []entry{
				{
					Entrypoint: "foo",
					Code:       "cmd/foo.go",
					Doc:        "docs/commands/foo.md",
					Tests:      []string{"cmd/foo_test.go"},
				},
				{
					Entrypoint: "bar",
					Code:       "cmd/bar.go",
					Doc:        "docs/commands/bar.md",
					Tests:      []string{"cmd/bar_test.go", "another_test.go"},
				},
			},
			want: 7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			bad, err := checkMissingFiles(tt.entries, dir, &buf)
			if err != nil {
				t.Fatal(err)
			}

			if bad != tt.want {
				t.Errorf("missing files: got %d, want %d", bad, tt.want)
			}
			out := buf.String()
			assertCode(t, out, "02", tt.want)
		})
	}
}

func createEmptyFile(t *testing.T, name string) {
	t.Helper()
	fh, err := os.Create(name)
	if err != nil {
		t.Fatal(err)
	}
	err = fh.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCheckMissingFilesGood(t *testing.T) {
	dir := t.TempDir()
	entries := []entry{
		{
			Entrypoint: "foo",
			Code:       "cmd/foo.go",
			Doc:        "docs/foo.md",
			Tests:      []string{"cmd/foo_test.go"},
		},
		{
			Entrypoint: "bar",
			Code:       "cmd/bar.go",
			Doc:        "docs/bar.md",
			Tests:      []string{"cmd/bar_test.go", "another_test.go"},
		},
	}

	err := os.MkdirAll(filepath.Join(dir, "cmd"), 0o777)
	if err != nil {
		t.Fatal(err)
	}
	err = os.MkdirAll(filepath.Join(dir, "docs"), 0o777)
	if err != nil {
		t.Fatal(err)
	}

	files := []string{
		filepath.Join("cmd", "foo.go"),
		filepath.Join("docs", "foo.md"),
		filepath.Join("cmd", "foo_test.go"),
		filepath.Join("cmd", "bar.go"),
		filepath.Join("docs", "bar.md"),
		filepath.Join("cmd", "bar_test.go"),
		"another_test.go",
	}
	for _, f := range files {
		createEmptyFile(t, filepath.Join(dir, f))
	}

	var buf bytes.Buffer
	bad, err := checkMissingFiles(entries, dir, &buf)
	if err != nil {
		t.Fatal(err)
	}

	if bad != 0 {
		t.Errorf("got %d missing files, want none", bad)
	}

	if out := buf.String(); out != "" {
		t.Errorf("got %q for output, want empty", out)
	}
}

func TestCheckEntrypointDocMismatchBad(t *testing.T) {
	tests := []struct {
		name    string
		entries []entry
		want    int
	}{
		{
			name: "one entry",
			entries: []entry{
				{
					Entrypoint: "foo",
					Doc:        "docs/commands/bar.md",
				},
			},
			want: 1,
		},
		{
			name: "two entries",
			entries: []entry{
				{
					Entrypoint: "foo",
					Doc:        "docs/commands/bar.md",
				},
				{
					Entrypoint: "bar",
					Doc:        "docs/commands/baz.md",
				},
			},
			want: 2,
		},
		{
			name: "skip",
			entries: []entry{
				{
					Entrypoint: "foo",
					Doc:        "docs/commands/bar.md",
					Skip:       true,
				},
				{
					Entrypoint: "bar",
					Doc:        "docs/commands/baz.md",
				},
			},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			bad, err := checkEntrypointDocMismatch(tt.entries, "", &buf)
			if err != nil {
				t.Fatal(err)
			}

			if bad != tt.want {
				t.Errorf("command/doc mismatches: got %d, want %d", bad, tt.want)
			}
			out := buf.String()
			assertCode(t, out, "03", tt.want)
		})
	}
}

func TestCheckEntrypointDocMismatchGood(t *testing.T) {
	entries := []entry{
		{
			Entrypoint: "foo bar",
			Doc:        "docs/commands/foo_bar.md",
		},
		{
			Entrypoint: "baz",
			Doc:        "docs/commands/baz.md",
		},
	}

	var buf bytes.Buffer
	bad, err := checkEntrypointDocMismatch(entries, "", &buf)
	if err != nil {
		t.Fatal(err)
	}

	if bad != 0 {
		t.Errorf("got %d mismatches, want none", bad)
	}

	if out := buf.String(); out != "" {
		t.Errorf("got %q for output, want empty", out)
	}
}

func TestCheckDupEntrypointsBad(t *testing.T) {
	entries := []entry{
		{
			Entrypoint: "foo",
		},
		{
			Entrypoint: "bar",
		},
		{
			Entrypoint: "foo",
		},
		{
			Entrypoint: "baz",
		},
	}

	var buf bytes.Buffer
	bad, err := checkDupEntrypoints(entries, &buf)
	if err != nil {
		t.Fatal(err)
	}

	if wantBad := 1; bad != wantBad {
		t.Errorf("duplicated entries: got %d, want %d", bad, wantBad)
	}
	out := buf.String()
	assertCode(t, out, "04", 1)
}

func TestCheckDupEntrypointsGood(t *testing.T) {
	entries := []entry{
		{
			Entrypoint: "foo",
		},
		{
			Entrypoint: "bar",
		},
		{
			Entrypoint: "baz",
		},
	}

	var buf bytes.Buffer
	bad, err := checkDupEntrypoints(entries, &buf)
	if err != nil {
		t.Fatal(err)
	}

	if bad != 0 {
		t.Errorf("got %d duplicated entries, want none", bad)
	}

	if out := buf.String(); out != "" {
		t.Errorf("got %q for output, want empty", out)
	}
}

func TestCheckMissingEntriesBad(t *testing.T) {
	dir := t.TempDir()
	files := []string{
		"foo_bar.md",
		"baz.md",
	}
	for _, f := range files {
		fname := filepath.Join(dir, f)
		createEmptyFile(t, fname)
	}

	tests := []struct {
		name    string
		entries []entry
		want    int
	}{
		{
			name:    "no entries",
			entries: []entry{},
			want:    2,
		},
		{
			name: "one entry",
			entries: []entry{
				{
					Entrypoint: "foo bar",
				},
			},
			want: 1,
		},
		{
			name: "other entry",
			entries: []entry{
				{
					Entrypoint: "baz",
				},
			},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			bad, err := checkMissingEntries(tt.entries, dir, "", &buf)
			if err != nil {
				t.Fatal(err)
			}

			if bad != tt.want {
				t.Errorf("missing entries: got %d, want %d", bad, tt.want)
			}
			out := buf.String()
			assertCode(t, out, "05", tt.want)
		})
	}
}

func TestCheckMissingEntriesGood(t *testing.T) {
	dir := t.TempDir()
	files := []string{
		"foo_bar.md",
		"baz.md",
	}
	for _, f := range files {
		fname := filepath.Join(dir, f)
		createEmptyFile(t, fname)
	}

	entries := []entry{
		{
			Entrypoint: "foo bar",
		},
		{
			Entrypoint: "baz",
		},
	}

	var buf bytes.Buffer
	bad, err := checkMissingEntries(entries, dir, "", &buf)
	if err != nil {
		t.Fatal(err)
	}

	if bad != 0 {
		t.Errorf("got %d missing entries, want none", bad)
	}

	if out := buf.String(); out != "" {
		t.Errorf("got %q for output, want empty", out)
	}
}

func writeEntries(t *testing.T, es []entry, outfile string) {
	t.Helper()
	bs, err := yaml.Marshal(&es)
	if err != nil {
		t.Fatal(err)
	}
	fd, err := os.Create(outfile)
	if err != nil {
		fd.Close()
		t.Fatal(err)
	}
	_, err = fd.Write(bs)
	if err != nil {
		t.Fatal(err)
	}

	err = fd.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCheckAll(t *testing.T) {
	dir := t.TempDir()
	docdir := filepath.Join(dir, "docs", "commands")
	err := os.MkdirAll(docdir, 0o777)
	if err != nil {
		t.Fatal(err)
	}

	codedir := filepath.Join(dir, "cmd")
	err = os.MkdirAll(codedir, 0o777)
	if err != nil {
		t.Fatal(err)
	}

	files := []string{
		filepath.Join(docdir, "README.md"),
		filepath.Join(docdir, "foo_bar.md"),
		filepath.Join(docdir, "baz.md"),
		filepath.Join(docdir, "skip.md"),
		filepath.Join(codedir, "foobar.go"),
		filepath.Join(codedir, "baz.go"),
		filepath.Join(codedir, "baz_test.go"),
	}
	for _, f := range files {
		createEmptyFile(t, f)
	}

	goodEntries := []entry{
		{
			Entrypoint: "foo bar",
			Code:       "cmd/foobar.go",
			Doc:        "docs/commands/foo_bar.md",
		},
		{
			Entrypoint: "baz",
			Code:       "cmd/baz.go",
			Doc:        "docs/commands/baz.md",
			Tests:      []string{"cmd/baz_test.go"},
		},
		{
			Entrypoint: "skip",
			Skip:       true,
		},
	}

	yfile := filepath.Join(dir, "docs", "matrix.yaml")
	writeEntries(t, goodEntries, yfile)

	t.Run("all good", func(t *testing.T) {
		toskip := make(map[string]bool)

		var buf bytes.Buffer
		bad, err := check(yfile, docdir, dir, "", toskip, &buf)
		if err != nil {
			t.Fatal(err)
		}

		if bad != 0 {
			t.Errorf("got %d missing entries, want none", bad)
		}

		out := buf.String()
		if out != "" {
			t.Errorf("got %q for output, want empty", out)
		}
	})

	// Now force each category of failure.

	badEntries := []entry{
		// Non-existent files (02), command/doc mismatch (03).
		{
			Entrypoint: "notthere",
			Code:       "cmd/../notthere.go",
			Doc:        "docs/commands/nt.md",
			Tests:      []string{"cmd/notthere_test.go"},
		},
		// Repeated entry (04).
		{
			Entrypoint: "foo bar",
			Code:       "cmd/foobar.go",
			Doc:        "docs/commands/foo_bar.md",
		},
	}

	writeEntries(t, append(goodEntries, badEntries...), yfile)
	// Documentation file without entry (05).
	createEmptyFile(t, filepath.Join(docdir, "noentry.md"))

	t.Run("some bad", func(t *testing.T) {
		toskip := make(map[string]bool)

		var buf bytes.Buffer
		bad, err := check(yfile, docdir, dir, "", toskip, &buf)
		if err != nil {
			t.Fatal(err)
		}

		wantBad := 7
		if bad != wantBad {
			t.Errorf("got %d missing entries, want %d", bad, wantBad)
		}

		out := buf.String()
		assertCode(t, out, "01", 1)
		assertCode(t, out, "02", 3)
		assertCode(t, out, "03", 1)
		assertCode(t, out, "04", 1)
		assertCode(t, out, "05", 1)

		buf.Reset()

		toskip["03"] = true
		toskip["05"] = true
		bad, err = check(yfile, docdir, dir, "", toskip, &buf)
		if err != nil {
			t.Fatal(err)
		}

		wantBad = 5
		if bad != wantBad {
			t.Errorf("got %d missing entries, want %d", bad, wantBad)
		}

		out = buf.String()
		assertCode(t, out, "01", 1)
		assertCode(t, out, "02", 3)
		assertCode(t, out, "03", 0)
		assertCode(t, out, "04", 1)
		assertCode(t, out, "05", 0)
	})
}

func TestCheckPrefix(t *testing.T) {
	dir := t.TempDir()
	docdir := filepath.Join(dir, "docs", "commands")
	err := os.MkdirAll(docdir, 0o777)
	if err != nil {
		t.Fatal(err)
	}

	codedir := filepath.Join(dir, "cmd")
	err = os.MkdirAll(codedir, 0o777)
	if err != nil {
		t.Fatal(err)
	}

	files := []string{
		filepath.Join(docdir, "foo.md"),
		filepath.Join(docdir, "bar.md"),
		filepath.Join(docdir, "baz.md"),
		filepath.Join(codedir, "foo.go"),
		filepath.Join(codedir, "foo_test.go"),
		filepath.Join(codedir, "bar.go"),
		filepath.Join(codedir, "bar_test.go"),
	}
	for _, f := range files {
		createEmptyFile(t, f)
	}

	entries := []entry{
		{
			Entrypoint: "foo",
			Code:       "cmd/foo.go",
			Doc:        "docs/commands/foo.md",
			Tests:      []string{"cmd/foo_test.go"},
		},
		{
			Entrypoint: "foo bar",
			Code:       "cmd/bar.go",
			Doc:        "docs/commands/bar.md",
			Tests:      []string{"cmd/bar_test.go"},
		},
	}

	yfile := filepath.Join(dir, "docs", "matrix.yaml")
	writeEntries(t, entries, yfile)

	toskip := make(map[string]bool)

	var buf bytes.Buffer
	bad, err := check(yfile, docdir, dir, "foo", toskip, &buf)
	if err != nil {
		t.Fatal(err)
	}

	if want := 1; bad != want {
		t.Errorf("got %d missing entries, want %d", bad, want)
	}

	out := buf.String()
	assertCode(t, out, "05", 1)
}
