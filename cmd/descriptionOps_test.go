package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestUnpackDescriptions(t *testing.T) {
	// This is just a light test that unpackDescriptions is wired up.  The desc
	// package has more detailed tests of DESCRIPTION parsing.
	descFoo := []byte(`Package: foo
Version: 0.4.0
`)
	descBar := []byte(`Package: bar
Version: 1.0.0
License: GPL (>=2)
`)

	dir := t.TempDir()
	err := os.Mkdir(filepath.Join(dir, "foo"), 0o777)
	if err != nil {
		t.Fatal(err)
	}
	err = os.Mkdir(filepath.Join(dir, "bar"), 0o777)
	if err != nil {
		t.Fatal(err)
	}

	foo := filepath.Join(dir, "foo", "DESCRIPTION")
	bar := filepath.Join(dir, "bar", "DESCRIPTION")

	err = os.WriteFile(foo, descFoo, 0o666)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(bar, descBar, 0o666)
	if err != nil {
		t.Fatal(err)
	}

	res := unpackDescriptions(afero.NewOsFs(), []string{foo, bar})
	assert.Equal(t, res[0].Package, "foo")
	assert.Equal(t, res[0].Version, "0.4.0")
	assert.Equal(t, res[1].Package, "bar")
	assert.Equal(t, res[1].Version, "1.0.0")
	assert.Equal(t, res[1].License, "GPL (>=2)")
}
