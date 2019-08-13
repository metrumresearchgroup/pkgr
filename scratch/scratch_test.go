package main

import (
	"github.com/spf13/afero"
	"testing"
)

func Test_scratch(t *testing.T) {
	fs := afero.NewOsFs()
	fs.Rename("./dump/dump_dir", "./dump/dump_dir_old")
}