// Copyright 2024 Metrum Research Group
// SPDX-License-Identifier: MIT

package main

import (
	"bytes"
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"golang.org/x/tools/cover"
)

func assertNearEqual(t *testing.T, a float64, b float64) {
	t.Helper()
	if math.IsNaN(a) || math.IsNaN(b) {
		t.Fatal("assertNearEqual: NaN values are not allowed")
	}

	var tol float64 = 1e-8
	d := math.Abs(a - b)
	if d >= tol {
		t.Errorf("absolute difference exceeds tolerance (%e)\na=%f\nb=%f", tol, a, b)
	}
}

func assertFileCoverage(t *testing.T, got []*fileCoverage, want []*fileCoverage) {
	t.Helper()

	ngot := len(got)
	nwant := len(want)
	if ngot != nwant {
		t.Errorf("expected coverage for %d files, got %d", nwant, ngot)

		return
	}

	mapWant := make(map[string]float64, nwant)
	for _, fc := range want {
		if _, found := mapWant[fc.File]; found {
			t.Error("repeated coverage for", fc.File)
		}
		mapWant[fc.File] = fc.Coverage
	}

	for _, fc := range got {
		perc, found := mapWant[fc.File]
		if found {
			assertNearEqual(t, fc.Coverage, perc)
		} else {
			t.Error("no coverage measurement for", fc.File)
		}
	}
}

// See cover.ParseProfilesFromReader for a description of the format.
var testProfile string = `mode: set
example.com/tmod/foo.go:3.16,5.2 1 1
example.com/tmod/foo.go:7.21,9.2 1 0
example.com/tmod/foo.go:11.21,13.2 1 0
example.com/tmod/foo.go:15.19,17.2 1 0
example.com/tmod/cmd/bar.go:3.18,5.2 1 1
example.com/tmod/cmd/bar.go:7.25,9.2 1 0
example.com/tmod/cmd/main.go:7.13,10.2 2 0
`

func TestPercentCovered(t *testing.T) {
	r := strings.NewReader(testProfile)
	profiles, err := cover.ParseProfilesFromReader(r)
	if err != nil {
		t.Fatal(err)
	}
	cov := percentCovered(profiles)
	assertNearEqual(t, cov.Overall, 25.0)
	assertFileCoverage(t, cov.Files,
		[]*fileCoverage{
			{
				File:     "example.com/tmod/foo.go",
				Coverage: 25.0,
			},
			{
				File:     "example.com/tmod/cmd/bar.go",
				Coverage: 50.0,
			},
			{
				File:     "example.com/tmod/cmd/main.go",
				Coverage: 0.0,
			},
		},
	)
}

func TestPercentCoveredZero(t *testing.T) {
	re := regexp.MustCompile(`(?m) 1$`)
	r := strings.NewReader(re.ReplaceAllString(testProfile, " 0"))
	profiles, err := cover.ParseProfilesFromReader(r)
	if err != nil {
		t.Fatal(err)
	}
	cov := percentCovered(profiles)
	assertNearEqual(t, cov.Overall, 0.0)
	assertFileCoverage(t, cov.Files,
		[]*fileCoverage{
			{
				File:     "example.com/tmod/foo.go",
				Coverage: 0.0,
			},
			{
				File:     "example.com/tmod/cmd/bar.go",
				Coverage: 0.0,
			},
			{
				File:     "example.com/tmod/cmd/main.go",
				Coverage: 0.0,
			},
		},
	)
}

func TestPercentCoveredNoFiles(t *testing.T) {
	r := strings.NewReader("mode: set\n")
	profiles, err := cover.ParseProfilesFromReader(r)
	if err != nil {
		t.Fatal(err)
	}
	cov := percentCovered(profiles)
	assertNearEqual(t, cov.Overall, 0.0)

	if len(cov.Files) != 0 {
		t.Errorf("expected coverage for 0 files, got %d", len(cov.Files))
	}
}

func TestWrite(t *testing.T) {
	r := strings.NewReader(testProfile)
	profiles, err := cover.ParseProfilesFromReader(r)
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	err = write(&buf, profiles, "", "")
	if err != nil {
		t.Fatal(err)
	}

	var cov coverage
	err = json.Unmarshal(buf.Bytes(), &cov)
	if err != nil {
		t.Fatal(err)
	}

	assertNearEqual(t, cov.Overall, 25.0)
	assertFileCoverage(t, cov.Files,
		[]*fileCoverage{
			{
				File:     "example.com/tmod/foo.go",
				Coverage: 25.0,
			},
			{
				File:     "example.com/tmod/cmd/bar.go",
				Coverage: 50.0,
			},
			{
				File:     "example.com/tmod/cmd/main.go",
				Coverage: 0.0,
			},
		},
	)
}

func TestWriteShortenNames(t *testing.T) {
	r := strings.NewReader(testProfile)
	profiles, err := cover.ParseProfilesFromReader(r)
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	err = write(&buf, profiles, "example.com/tmod", "")
	if err != nil {
		t.Fatal(err)
	}

	var cov coverage
	err = json.Unmarshal(buf.Bytes(), &cov)
	if err != nil {
		t.Fatal(err)
	}

	assertNearEqual(t, cov.Overall, 25.0)
	assertFileCoverage(t, cov.Files,
		[]*fileCoverage{
			{
				File:     "foo.go",
				Coverage: 25.0,
			},
			{
				File:     "cmd/bar.go",
				Coverage: 50.0,
			},
			{
				File:     "cmd/main.go",
				Coverage: 0.0,
			},
		},
	)
}

func chdir(t *testing.T, dir string) {
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	err = os.Chdir(dir)
	if err != nil {
		_ = os.Chdir(old)
		t.Fatal(err)
	}

	t.Cleanup(func() {
		_ = os.Chdir(old)
	})
}

func setupRunDir(t *testing.T) string {
	dir := t.TempDir()
	realmodPath := filepath.Join(dir, "realmod")
	err := os.MkdirAll(filepath.Join(realmodPath, "cmd"), 0777)
	if err != nil {
		t.Fatal(err)
	}

	modPath := filepath.Join(dir, "mod")
	err = os.Symlink(realmodPath, modPath)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(
		filepath.Join(modPath, "go.mod"),
		[]byte("module example.com/tmod\n\ngo 1.22.5"),
		0666)
	if err != nil {
		t.Fatal(err)
	}

	tfiles := []string{
		"foo.go",
		"cmd/bar.go",
		"cmd/main.go",
	}
	for _, f := range tfiles {
		fh, err := os.Create(filepath.Join(modPath, f))
		if err != nil {
			t.Fatal(err)
		}
		err = fh.Close()
		if err != nil {
			t.Fatal(err)
		}
	}

	chdir(t, modPath)

	return modPath
}

func TestRun(t *testing.T) {
	modPath := setupRunDir(t)
	pcontent := strings.ReplaceAll(testProfile,
		"example.com/tmod/cmd/main.go",
		modPath+"/"+"cmd/main.go")
	profPath := filepath.Join(modPath, "coverage.out")
	err := os.WriteFile(profPath, []byte(pcontent), 0666)
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	err = run(profPath, "go.mod", &buf)
	if err != nil {
		t.Fatal(err)
	}

	var cov coverage
	err = json.Unmarshal(buf.Bytes(), &cov)
	if err != nil {
		t.Fatal(err)
	}

	assertNearEqual(t, cov.Overall, 25.0)
	assertFileCoverage(t, cov.Files,
		[]*fileCoverage{
			{
				File:     "foo.go",
				Coverage: 25.0,
			},
			{
				File:     "cmd/bar.go",
				Coverage: 50.0,
			},
			{
				File:     "cmd/main.go",
				Coverage: 0.0,
			},
		},
	)
}

func TestRunNoGoMod(t *testing.T) {
	modPath := setupRunDir(t)
	profPath := filepath.Join(modPath, "coverage.out")
	err := os.WriteFile(profPath, []byte(testProfile), 0666)
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	err = run(profPath, "", &buf)
	if err != nil {
		t.Fatal(err)
	}

	var cov coverage
	err = json.Unmarshal(buf.Bytes(), &cov)
	if err != nil {
		t.Fatal(err)
	}

	assertNearEqual(t, cov.Overall, 25.0)
	assertFileCoverage(t, cov.Files,
		[]*fileCoverage{
			{
				File:     "example.com/tmod/foo.go",
				Coverage: 25.0,
			},
			{
				File:     "example.com/tmod/cmd/bar.go",
				Coverage: 50.0,
			},
			{
				File:     "example.com/tmod/cmd/main.go",
				Coverage: 0.0,
			},
		},
	)
}
