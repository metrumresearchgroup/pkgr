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

func Test_rAdd(t *testing.T) {
	type args struct {
		ccmd *cobra.Command
		args []string
	}

	tests := []testData{
		{
			ymlfolder: "../integration_tests/simple",
			data: []pkgData{
				{
					name:   "shiny",
					folder: "../integration_tests/simple/test-library/shiny",
				},
				{
					name:   "abc",
					folder: "../integration_tests/simple/test-library/abc",
				},
			},
		},
	}
	fs := afero.NewOsFs()
	for _, tt := range tests {

		t.Log("testing add ...")
		errstr, err := pkgrCommand("add", tt, "--install")
		assert.Equal(t, nil, err, fmt.Sprintf("Package add error. Check state of yml file <%s>. <%s> ", tt.ymlfolder, errstr))

		for _, d := range tt.data {
			b, _ := afero.FileContainsBytes(fs, filepath.Join(tt.ymlfolder, "pkgr.yml"), []byte(d.name))
			assert.Equal(t, true, b, fmt.Sprintf("Package not added:%s", d.name))
		}

		t.Log("testing remove ...")
		errstr, err = pkgrCommand("remove", tt, "")
		assert.Equal(t, nil, err, fmt.Sprintf("Package remove error. Check state of yml file <%s>. <%s> ", tt.ymlfolder, errstr))

		for _, d := range tt.data {
			b, _ := afero.FileContainsBytes(fs, filepath.Join(tt.ymlfolder, "pkgr.yml"), []byte(d.name))
			assert.Equal(t, !true, b, fmt.Sprintf("Package not removed:%s", d.name))
		}

		t.Log("testing pkgr.yml for difference ...")
		b, msg := isDiff(filepath.Join(tt.ymlfolder, "pkgr.yml"))
		assert.Equal(t, false, b, fmt.Sprintf("Yml file differs:%s", msg))

	}
}

func pkgrCommand(cmd string, test testData, flag string) (string, error) {

	var stderr bytes.Buffer
	pkgr := os.Getenv("HOME") + "/go/bin/pkgr"
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

func isDiff(filename string) (bool, string) {
	var stdout, stderr bytes.Buffer
	git := "/usr/bin/git"
	args := []string{"diff"}
	args = append(args, filename)

	command := exec.Command(git, args...)
	command.Stdout = &stdout
	command.Stderr = &stderr

	err := command.Run()
	if err != nil {
		errstr := string(stderr.Bytes())
		return true, errstr
	}

	outstr := string(stdout.Bytes())
	if len(outstr) > 0 {
		return true, outstr
	}
	return false, ""
}
