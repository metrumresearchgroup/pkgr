package version_test

import (
	"context"
	"github.com/metrumresearchgroup/command"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	versionE2ETest1 = "VSN-E2E-001"
	versionE2ETest2 = "VSN-E2E-002"
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
		assert.Equal(t, "dev", string(res.Output))
	})

	t.Run("short flag -v works", func(t *testing.T) {
		res, err := testCmd.Run(ctx, "pkgr", "--version")
		if err != nil {
			t.Fatal(err)
		}
		// this will always get the git tag for regular releases
		assert.Equal(t,"dev",  string(res.Output))
	})

}
