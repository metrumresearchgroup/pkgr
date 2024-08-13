// Copyright 2024 Metrum Research Group
// SPDX-License-Identifier: MIT

package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const usageMessage = `usage: checkmat <yaml> <dir>

Run various checks on the traceability matrix defined in <yaml>.  Each file in
<dir> is taken as the documentation for an entry point, where the base name maps
to an entry point in <yaml> once the file extension is removed and underscores
are substituted for spaces.

Verify that

 [01] each file name is valid

      To be considered "valid", a file name must be relative, must use "/" as
      the path separator, and must not contain "." or ".." elements.

 [02] each named file exists

      The current working directory is taken as the top-level project directory,
      and the files should be relative to this.

 [03] the base file name for the documentation matches the entry point name

 [04] no entry point has more than one entry in <yaml>

 [05] for each file in <dir>, an entry with a matching entry point is found in
      <yaml>

      This check is skipped for any entries with a true "skip" value.

Exit with status 1 if any issues are found.
`

func usage() {
	fmt.Fprint(flag.CommandLine.Output(), usageMessage)
}

type entry struct {
	Entrypoint string
	Code       string
	Doc        string
	Tests      []string
	Skip       bool
}

func readEntries(f string) ([]entry, error) {
	bs, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}
	var es []entry
	if err := yaml.Unmarshal(bs, &es); err != nil {
		return nil, err
	}

	return es, nil
}

func entryFiles(e entry) []string {
	var fs []string
	if e.Code != "" {
		fs = append(fs, e.Code)
	}
	if e.Doc != "" {
		fs = append(fs, e.Doc)
	}
	fs = append(fs, e.Tests...)

	return fs
}

func checkValidFileNames(es []entry, w io.Writer) (int, error) {
	var bad int

	for _, e := range es {
		for _, f := range entryFiles(e) {
			if !fs.ValidPath(f) {
				bad++
				fmt.Fprintln(w, "[01] file name is invalid:", f)
			}
		}
	}

	return bad, nil
}

func checkMissingFiles(es []entry, topdir string, w io.Writer) (int, error) {
	var bad int

	for _, e := range es {
		for _, f := range entryFiles(e) {
			// Note: Join() takes care of replacing slashes with
			// os.PathSeparator.
			f = filepath.Join(topdir, f)
			_, err := os.Stat(f)
			if errors.Is(err, os.ErrNotExist) {
				bad++
				fmt.Fprintln(w, "[02] file does not exist:", f)
			} else if err != nil {
				return bad, err
			}
		}
	}

	return bad, nil
}

// TODO: Make it possible to customize how to map the doc file to the entry
// point name.  As is it is only useful for command-line subcommands or one-off
// top-level commands, and it assumes that none of the commands have an
// underscore in their name.

func docToEntrypoint(f string) string {
	base := filepath.Base(f)
	name := strings.TrimSuffix(base, filepath.Ext(base))

	return strings.ReplaceAll(name, "_", " ")
}

func checkEntrypointDocMismatch(es []entry, w io.Writer) (int, error) {
	var bad int

	for _, e := range es {
		if e.Skip {
			continue
		}
		if docToEntrypoint(e.Doc) != e.Entrypoint {
			bad++
			fmt.Fprintf(w, "[03] entry point and doc file mismatch: %q != %q\n",
				e.Entrypoint, e.Doc)
		}
	}

	return bad, nil
}

func checkDupEntrypoints(es []entry, w io.Writer) (int, error) {
	var bad int

	cmds := make(map[string]bool)
	for _, e := range es {
		if _, found := cmds[e.Entrypoint]; found {
			fmt.Fprintf(w, "[04] entry point %q defined more than once\n",
				e.Entrypoint)
			bad++
		} else {
			cmds[e.Entrypoint] = true
		}
	}

	return bad, nil
}

func checkMissingEntries(es []entry, docdir string, w io.Writer) (int, error) {
	var bad int

	fh, err := os.Open(docdir)
	if err != nil {
		return bad, err
	}
	defer fh.Close()

	fnames, err := fh.Readdirnames(-1)
	if err != nil {
		return bad, err
	}

	cmds := make(map[string]bool)
	for _, e := range es {
		cmds[e.Entrypoint] = true
	}

	for _, f := range fnames {
		if _, found := cmds[docToEntrypoint(f)]; !found {
			fmt.Fprintf(w, "[05] No yaml entry for %q\n", filepath.Join(docdir, f))
			bad++
		}
	}

	return bad, nil
}

// check runs all the check functions on the traceability matrix defined in file
// yaml and returns the total number of issues found.  docdir points to a
// directory containing the documentation files.  topdir is an absolute path to
// top-level directory to which files in `yaml` are specified as relative.
//
// For each issue found, a message is written to w.
func check(yaml string, docdir string, topdir string, w io.Writer) (int, error) {
	var bad int

	entries, err := readEntries(yaml)
	if err != nil {
		return bad, err
	}

	type check func([]entry, io.Writer) (int, error)
	checks := []check{
		checkValidFileNames,
		func(es []entry, w io.Writer) (int, error) {
			return checkMissingFiles(es, topdir, w)
		},
		checkEntrypointDocMismatch,
		checkDupEntrypoints,
		func(es []entry, w io.Writer) (int, error) {
			return checkMissingEntries(es, docdir, w)
		},
	}

	for _, f := range checks {
		n, err := f(entries, w)
		if err != nil {
			return bad, err
		}
		bad += n
	}

	return bad, nil
}

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		flag.CommandLine.SetOutput(os.Stdout)
	}
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		flag.Usage()
		os.Exit(2)
	}

	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(2)
	}

	bad, err := check(args[0], args[1], wd, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(2)
	}

	if bad > 0 {
		fmt.Printf("\nproblems found: %d\n", bad)
		os.Exit(1)
	}
}
