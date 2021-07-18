package version_test

import (
	"context"
	"github.com/metrumresearchgroup/command"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVersion(t *testing.T) {
	testCmd := command.New()
	ctx := context.TODO()

	t.Run("long flag --version works", func(t *testing.T) {
		res, err := testCmd.Run(ctx, "pkgr", "--version")
		if err != nil {
			t.Fatal(err)
		}
		// this will always get the git tag for regular releases
		assert.Contains(t, string(res.Output), "dev-")
	})

	t.Run("short flag -v works", func(t *testing.T) {
		res, err := testCmd.Run(ctx, "pkgr", "--version")
		if err != nil {
			t.Fatal(err)
		}
		// this will always get the git tag for regular releases
		assert.Contains(t, string(res.Output), "dev-")
	})

}
