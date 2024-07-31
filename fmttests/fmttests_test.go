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
			var buf bytes.Buffer
			s, err := processEvents(makeEventReader(t, tt.events), tt.subtests, &buf)
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

			out := buf.String()
			outWant := strings.Join(tt.lines, "\n") + "\n"
			if out != outWant {
				t.Errorf("stdout:\n  want %q,\n   got %q", outWant, out)
			}
		})
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
			var buf bytes.Buffer
			_, err := processEvents(makeEventReader(t, tt.events), false, &buf)
			if err == nil {
				t.Errorf("processEvents unexpectedly passed")
			}
		})
	}
}
