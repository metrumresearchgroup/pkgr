// Copyright 2024 Metrum Research Group
// SPDX-License-Identifier: MIT

package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

const usageMessage = `usage: fmttests [-allow-skips] [-subtests]

Read 'go test -json' lines from standard input and display a formatted output
line for each test result record (action of "pass", "fail", or "skip").

Passing subtests are not displayed unless the -subtests flag is specified.
Failed or skipped subtests are always displayed.

If any "fail" record is encountered, exit with status 1.  Any "skip" record will
also trigger an exit with status 1 unless the -allow-skips flag is specified.
`

var (
	subtests   = flag.Bool("subtests", false, "")
	allowSkips = flag.Bool("allow-skips", false, "")
)

func usage() {
	fmt.Fprint(flag.CommandLine.Output(), usageMessage)
}

// Modified from Go's src/cmd/internal/test2json/test2json.go.
type event struct {
	Time    time.Time `json:",omitempty"`
	Action  string
	Package string  `json:",omitempty"`
	Test    string  `json:",omitempty"`
	Elapsed float64 `json:",omitempty"`
	Output  string  `json:",omitempty"`
}

type summary struct {
	Failed  []string
	Passed  []string
	Skipped []string
}

func packageBaseName(s string) string {
	xs := strings.Split(s, "/")
	n := len(xs)
	if n == 0 {
		return s
	}

	return xs[n-1]
}

// processEvents reads a JSON record of a test event from r.  For each result
// record encountered (i.e. a record with an action of "pass", "fail", or
// "skip"), it writes a formatted line to w.  The subtests argument controls
// whether results for subtests are written.
//
// processEvents returns a summary instance that records the test names for each
// result record encountered.
func processEvents(r io.Reader, subtests bool, w io.Writer) (summary, error) {
	var res summary

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		var e event

		if err := json.Unmarshal(scanner.Bytes(), &e); err != nil {
			return res, err
		}
		if e.Test == "" {
			continue
		}

		var status string
		var alwaysShow bool
		switch e.Action {
		case "fail":
			res.Failed = append(res.Failed, e.Test)
			status = "failed"
			alwaysShow = true
		case "pass":
			res.Passed = append(res.Passed, e.Test)
			status = "passed"
		case "skip":
			res.Skipped = append(res.Skipped, e.Test)
			status = "skipped"
			alwaysShow = true
		case "bench":
			return res, errors.New("bench output is not supported")
		case "cont", "output", "pause", "run", "start":
			continue
		default:
			return res, fmt.Errorf("unknown action: %s", e.Action)
		}

		if subtests || alwaysShow || !strings.ContainsAny(e.Test, "/") {
			// TODO: Consider other approaches for formatting the package name
			// that avoid collisions (e.g., packageBaseName returns "cmd" for
			// both ".../foo/cmd" and ".../bar/cmd").
			fmt.Fprintf(w, "[%s] %s: %s\n",
				packageBaseName(e.Package), e.Test, status)
		}
	}

	if err := scanner.Err(); err != nil {
		return res, err
	}

	return res, nil
}

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		flag.CommandLine.SetOutput(os.Stdout)
	}
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 0 {
		flag.Usage()
		os.Exit(2)
	}

	res, err := processEvents(os.Stdin, *subtests, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(2)
	}

	if len(res.Failed) > 0 || (!*allowSkips && len(res.Skipped) > 0) {
		fmt.Fprintf(os.Stderr, "failed tests: %d, skipped tests: %d\n",
			len(res.Failed), len(res.Skipped))
		os.Exit(1)
	}
}
