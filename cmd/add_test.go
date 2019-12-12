// Copyright © 2018 Devin Pastoor <devin.pastoor@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

type pkgData struct {
	name   string
	folder string
}

type testData struct {
	ymlfolder string
	data      []pkgData
}

func Test_rAddAndDelete(t *testing.T) {
	type args struct {
		ccmd *cobra.Command
		args []string
	}

	tests := []testData{
		{
			ymlfolder: "testsite/golden/simple",
			data: []pkgData{
				{
					name:   "R6",
					folder: "testsite/golden/simple/test-library/shiny",
				},
				{
					name:   "jsonlite",
					folder: "testsite/golden//simple/test-library/abc",
				},
			},
		},
		// Note: before adding a test, make sure the base test "pkgr install" works before testing "pkgr add --install"
	}
	fs := afero.NewOsFs()
	//fs := afero.NewMemMapFs()
	pkgrYamlContent := []byte(`
Version: 1

Packages:
  - fansi
Repos:
  - MPN: "https://mpn.metworx.com/snapshots/stable/2019-12-02"

Library: "test-library"
`)

	afero.WriteFile(fs, "testsite/golden/simple/pkgr.yml", pkgrYamlContent, 0755)

	defer fs.RemoveAll("testsite/golden/simple/test-library")
	for _, tt := range tests {

		ymlFilename := filepath.Join(tt.ymlfolder, "pkgr.yml")
		b, _ := afero.Exists(fs, ymlFilename)
		assert.Equal(t, true, b, fmt.Sprintf("yml file not found:%s", ymlFilename))
		ymlStart, _ := afero.ReadFile(fs, ymlFilename)

		t.Log("testing add ...")
		errstr, err := pkgrCommand("add", tt, "--install")
		assert.Equal(t, nil, err, fmt.Sprintf("Package add error. Check state of yml file <%s>. <%s> ", tt.ymlfolder, errstr))

		for _, d := range tt.data {
			b, _ := afero.FileContainsBytes(fs, filepath.Join(tt.ymlfolder, "pkgr.yml"), []byte(d.name))
			installed, _ := afero.DirExists(fs, filepath.Join(tt.ymlfolder, "test-library", d.name)) // make sure packages were installed with --install flag
			assert.Equal(t, true, b, fmt.Sprintf("Package not added:%s", d.name))
			assert.Equal(t, true, installed, fmt.Sprintf("Package not installed after being added:%s", d.name))
		}

		t.Log("testing remove ...")
		errstr, err = pkgrCommand("remove", tt, "")
		assert.Equal(t, nil, err, fmt.Sprintf("Package remove error. Check state of yml file <%s>. <%s> ", tt.ymlfolder, errstr))

		for _, d := range tt.data {
			b, _ := afero.FileContainsBytes(fs, filepath.Join(tt.ymlfolder, "pkgr.yml"), []byte(d.name))
			assert.Equal(t, !true, b, fmt.Sprintf("Package not removed:%s", d.name))
		}

		t.Log("testing pkgr.yml for difference ...")
		ymlEnd, _ := afero.ReadFile(fs, ymlFilename)
		n := bytes.Compare(cleanWs(ymlStart), cleanWs(ymlEnd))
		assert.Equal(t, 0, n, fmt.Sprintf("Yml file differs"))

		t.Log("restoring yml file ...")
		fi, _ := os.Stat(ymlFilename)
		err = afero.WriteFile(fs, ymlFilename, ymlStart, fi.Mode())
		assert.Equal(t, nil, err, "Error restoring yml file")
	}
}

func pkgrCommand(cmd string, test testData, flag string) (string, error) {
	var stderr bytes.Buffer
	pkgr := filepath.Join(os.Getenv("HOME"), "/go/bin/pkgr")
	args := []string{cmd}
	for _, d := range test.data {
		args = append(args, d.name)
	}

	if len(flag) > 0 {
		args = append(args, flag)
	}

	command := exec.Command(pkgr, args...)
	command.Dir = test.ymlfolder
	command.Stderr = &stderr

	err := command.Run()
	errstr := string(stderr.Bytes())
	return errstr, err
}

// cleanWs removes blank lines and leading and trailing white space
// it makes the above bytes.Compare(..) equivilent to: /usr/bin/diff -wB
func cleanWs(b []byte) []byte {
	var cb []byte
	lines := bytes.Split(b, []byte("\n"))
	for _, line := range lines {
		tl := bytes.Trim(line, " ")
		if len(tl) > 0 {
			cb = append(cb, line...)
		}
	}
	return cb
}

// // saving in case ...
// func diff(file1, file2 string) bool {
// 	var stdout bytes.Buffer
// 	diff := "/usr/bin/diff"
// 	args := []string{"-wB", file1, file2}

// 	command := exec.Command(diff, args...)
// 	command.Stdout = &stdout

// 	err := command.Run()
// 	outstr := string(stdout.Bytes())

// 	if err != nil || len(outstr) > 0 {
// 		return true
// 	}
// 	return false
// }
