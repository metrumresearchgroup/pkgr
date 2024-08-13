// Copyright 2024 Metrum Research Group
// SPDX-License-Identifier: MIT

package main

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"testing"
	"time"
)

type testEvent struct {
	Action  string
	Package string
	Test    string
	Output  string
}

func makeEventReader(t *testing.T, tes []testEvent) io.Reader {
	t.Helper()
	var buf bytes.Buffer
	var e event
	for _, te := range tes {
		e = event{
			Time:    time.Time{},
			Action:  te.Action,
			Package: te.Package,
			Test:    te.Test,
			Output:  te.Output,
		}
		bs, err := json.Marshal(&e)
		if err != nil {
			t.Fatal(err)
		}
		buf.Write(bs)
		buf.Write([]byte("\n"))
	}

	return bytes.NewReader(buf.Bytes())
}

func TestProcessEvents(t *testing.T) {
	var tests = []struct {
		name     string
		events   []testEvent
		subtests bool
		lines    []string
		nfailed  int
		npassed  int
		nskipped int
	}{
		{
			name: "base",
			events: []testEvent{
				{
					Action:  "pass",
					Package: "example.com/abc/def",
					Test:    "TestFoo",
				},
				{
					Action:  "pass",
					Package: "example.com/abc/def",
					Test:    "TestFoo/Bar",
				},
				{
					Action:  "pass",
					Package: "example.com/ghi/cmd",
					Test:    "TestBaz",
				},
			},
			subtests: false,
			lines: []string{
				"[def] TestFoo: passed",
				"[cmd] TestBaz: passed",
			},
			nfailed:  0,
			npassed:  3,
			nskipped: 0,
		},
		{
			name: "subtests via flag",
			events: []testEvent{
				{
					Action:  "pass",
					Package: "example.com/abc/def",
					Test:    "TestFoo",
				},
				{
					Action:  "pass",
					Package: "example.com/abc/def",
					Test:    "TestFoo/Bar",
				},
				{
					Action:  "pass",
					Package: "example.com/ghi/cmd",
					Test:    "TestBaz",
				},
			},
			subtests: true,
			lines: []string{
				"[def] TestFoo: passed",
				"[def] TestFoo/Bar: passed",
				"[cmd] TestBaz: passed",
			},
			nfailed:  0,
			npassed:  3,
			nskipped: 0,
		},
		{
			name: "include skipped subtests",
			events: []testEvent{
				{
					Action:  "pass",
					Package: "example.com/abc/def",
					Test:    "TestFoo",
				},
				{
					Action:  "skip",
					Package: "example.com/abc/def",
					Test:    "TestFoo/Bar",
				},
				{
					Action:  "pass",
					Package: "example.com/ghi/cmd",
					Test:    "TestBaz",
				},
			},
			subtests: false,
			lines: []string{
				"[def] TestFoo: passed",
				"[def] TestFoo/Bar: skipped",
				"[cmd] TestBaz: passed",
			},
			nfailed:  0,
			npassed:  2,
			nskipped: 1,
		},
		{
			name: "no package",
			events: []testEvent{
				{
					Action:  "pass",
					Package: "",
					Test:    "TestFoo",
				},
			},
			subtests: false,
			lines: []string{
				"[] TestFoo: passed",
			},
			nfailed:  0,
			npassed:  1,
			nskipped: 0,
		},
		{
			name: "include failed subtests",
			events: []testEvent{
				{
					Action:  "pass",
					Package: "example.com/abc/def",
					Test:    "TestFoo",
				},
				{
					Action:  "fail",
					Package: "example.com/abc/def",
					Test:    "TestFoo/Bar",
				},
				{
					Action:  "pass",
					Package: "example.com/ghi/cmd",
					Test:    "TestBaz",
				},
			},
			subtests: false,
			lines: []string{
				"[def] TestFoo: passed",
				"[def] TestFoo/Bar: failed",
				"[cmd] TestBaz: passed",
			},
			nfailed:  1,
			npassed:  2,
			nskipped: 0,
		},
		{
			name: "filtered",
			events: []testEvent{
				{
					Action:  "start",
					Package: "example.com/abc/def",
				},
				{
					Action:  "run",
					Package: "example.com/abc/def",
					Test:    "TestFoo",
				},
				{
					Action:  "pass",
					Package: "example.com/abc/def",
					Test:    "TestFoo",
				},
			},
			subtests: false,
			lines: []string{
				"[def] TestFoo: passed",
			},
			nfailed:  0,
			npassed:  1,
			nskipped: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bufStdout bytes.Buffer
			var bufStderr bytes.Buffer
			s, err := processEvents(makeEventReader(t, tt.events),
				tt.subtests, &bufStdout, &bufStderr)
			if err != nil {
				t.Fatal(err)
			}

			if len(s.Failed) != tt.nfailed {
				t.Errorf("failed: want %d, got %d", tt.nfailed, len(s.Failed))
			}
			if len(s.Passed) != tt.npassed {
				t.Errorf("passed: want %d, got %d", tt.npassed, len(s.Passed))
			}
			if len(s.Skipped) != tt.nskipped {
				t.Errorf("skipped: want %d, got %d", tt.nskipped, len(s.Skipped))
			}

			stdout := bufStdout.String()
			stdoutWant := strings.Join(tt.lines, "\n") + "\n"
			if stdout != stdoutWant {
				t.Errorf("stdout:\n  want %q,\n   got %q", stdoutWant, stdout)
			}

			stderr := bufStderr.String()
			if stderr != "" {
				t.Errorf("expected empty stderr, got %q", stderr)
			}
		})
	}
}

func TestProcessEventsFailureOutput(t *testing.T) {
	events := []testEvent{
		{
			Action:  "output",
			Package: "example.com/o/foo",
			Test:    "TestFoo1",
			Output:  "foo1 output\n",
		},
		{
			Action:  "pass",
			Package: "example.com/o/foo",
			Test:    "TestFoo1",
		},
		{
			Action:  "output",
			Package: "example.com/o/foo",
			Test:    "TestFoo2",
			Output:  "TestFoo2 failed\n",
		},
		{
			Action:  "output",
			Package: "example.com/o/foo",
			Test:    "TestFoo2",
			Output:  "error was ...\n",
		},
		{
			Action:  "fail",
			Package: "example.com/o/foo",
			Test:    "TestFoo2",
		},
		{
			Action:  "pass",
			Package: "example.com/o/k/bar",
			Test:    "TestBar",
		},
		{
			Action:  "output",
			Package: "example.com/o/k/baz",
			Test:    "TestBaz/1",
			Output:  "TestBaz/1 failed\n",
		},
		{
			Action:  "pass",
			Package: "example.com/o/k/baz",
			Test:    "TestBaz/2",
		},
		{
			Action:  "fail",
			Package: "example.com/o/k/baz",
			Test:    "TestBaz/1",
		},
		{
			Action:  "fail",
			Package: "example.com/o/k/baz",
			Test:    "TestBaz",
		},
	}
	var bufStdout bytes.Buffer
	var bufStderr bytes.Buffer
	s, err := processEvents(makeEventReader(t, events),
		false, &bufStdout, &bufStderr)
	if err != nil {
		t.Fatal(err)
	}

	if nfailedWant := 3; len(s.Failed) != nfailedWant {
		t.Errorf("failed: want %d, got %d", nfailedWant, len(s.Failed))
	}

	if npassedWant := 3; len(s.Passed) != npassedWant {
		t.Errorf("passed: want %d, got %d", npassedWant, len(s.Passed))
	}

	if nskippedWant := 0; len(s.Skipped) != nskippedWant {
		t.Errorf("skipped: want %d, got %d", nskippedWant, len(s.Skipped))
	}

	stdoutWant := strings.Join([]string{
		"[foo] TestFoo1: passed",
		"[foo] TestFoo2: failed",
		"[bar] TestBar: passed",
		"[baz] TestBaz/1: failed",
		"[baz] TestBaz: failed",
	}, "\n") + "\n"
	if stdout := bufStdout.String(); stdout != stdoutWant {
		t.Errorf("stdout:\n  want %q,\n   got %q", stdoutWant, stdout)
	}

	stderrWant := strings.Join([]string{
		"TestFoo2 failed",
		"error was ...",
		"TestBaz/1 failed",
	}, "\n") + "\n"
	if stderr := bufStderr.String(); stderr != stderrWant {
		t.Errorf("stderr: want %q, got %q", stderrWant, stderr)
	}
}

func TestProcessEventsErrors(t *testing.T) {
	var tests = []struct {
		name     string
		events   []testEvent
		subtests bool
		lines    []string
		nfailed  int
		npassed  int
		nskipped int
	}{
		{
			name: "bench",
			events: []testEvent{
				{
					Action:  "pass",
					Package: "example.com/abc/def",
					Test:    "TestFoo",
				},
				{
					Action:  "bench",
					Package: "example.com/abc/def",
					Test:    "TestFoo/Bar",
				},
			},
		},
		{
			name: "unknown",
			events: []testEvent{
				{
					Action:  "pass",
					Package: "example.com/abc/def",
					Test:    "TestFoo",
				},
				{
					Action:  "unknown",
					Package: "example.com/abc/def",
					Test:    "TestFoo/Bar",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bufStdout bytes.Buffer
			var bufStderr bytes.Buffer
			_, err := processEvents(makeEventReader(t, tt.events),
				false, &bufStdout, &bufStderr)
			if err == nil {
				t.Errorf("processEvents unexpectedly passed")
			}
		})
	}
}
