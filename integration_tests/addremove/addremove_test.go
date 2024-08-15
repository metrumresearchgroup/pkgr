package addremove

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/metrumresearchgroup/command"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"

	"github.com/metrumresearchgroup/pkgr/configlib"
)

func TestAddRemove(t *testing.T) {
	dir := t.TempDir()

	cfg := &configlib.PkgrConfig{
		Version: 1,
		Packages: []string{
			"foo",
			"bar",
		},
		Repos: []map[string]string{
			{
				"MPN": "mpn",
			},
		},
		Library: "lib",
		Cache:   "cache",
	}

	bs, err := yaml.Marshal(cfg)
	if err != nil {
		t.Fatal(err)
	}
	pkgrfile := filepath.Join(dir, "pkgr.yml")
	err = os.WriteFile(pkgrfile, bs, 0o666)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("add", func(t *testing.T) {
		cmd := command.New("pkgr", "add", "baz1", "baz2")
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("%s\n%s\n", out, err)
		}

		bsNew, err := os.ReadFile(pkgrfile)
		if err != nil {
			t.Fatal(err)
		}

		var cfgNew configlib.PkgrConfig
		err = yaml.Unmarshal(bsNew, &cfgNew)
		if err != nil {
			t.Fatal(err)
		}

		assert.Contains(t, cfgNew.Packages, "foo")
		assert.Contains(t, cfgNew.Packages, "bar")
		assert.Contains(t, cfgNew.Packages, "baz1")
		assert.Contains(t, cfgNew.Packages, "baz2")
	})

	t.Run("remove", func(t *testing.T) {
		cmd := command.New("pkgr", "remove", "foo", "baz1")
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("%s\n%s\n", out, err)
		}

		bsNew, err := os.ReadFile(pkgrfile)
		if err != nil {
			t.Fatal(err)
		}

		var cfgNew configlib.PkgrConfig
		err = yaml.Unmarshal(bsNew, &cfgNew)
		if err != nil {
			t.Fatal(err)
		}

		assert.NotContains(t, cfgNew.Packages, "foo")
		assert.Contains(t, cfgNew.Packages, "bar")
		assert.NotContains(t, cfgNew.Packages, "baz1")
		assert.Contains(t, cfgNew.Packages, "baz2")
	})
}

func TestAddInstall(t *testing.T) {
	dir := t.TempDir()

	repo, err := filepath.Abs("../../localrepos/simple")
	if err != nil {
		t.Fatal(err)
	}

	cfg := &configlib.PkgrConfig{
		Version: 1,
		Packages: []string{
			"R6",
		},
		Repos: []map[string]string{
			{
				"local": repo,
			},
		},
		Customizations: configlib.Customizations{
			Repos: []map[string]configlib.RepoConfig{
				{
					"local": {
						Type: "Source",
					},
				},
			},
		},
		Library: "lib",
		Cache:   "cache",
	}

	bs, err := yaml.Marshal(cfg)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(filepath.Join(dir, "pkgr.yml"), bs, 0o666)
	if err != nil {
		t.Fatal(err)
	}

	cmd := command.New("pkgr", "add", "--install", "cli")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("%s\n%s\n", out, err)
	}

	assert.DirExists(t, filepath.Join(dir, "lib", "R6"))
	assert.DirExists(t, filepath.Join(dir, "lib", "cli"))
}
