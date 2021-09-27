package version_test

import (
	"testing"

	"github.com/metrumresearchgroup/command"
	"github.com/stretchr/testify/assert"
)

const (
	versionE2ETest1 = "VSN-E2E-001"
	versionE2ETest2 = "VSN-E2E-002"
)

func TestVersion(t *testing.T) {
	t.Run("long flag --version works", func(t *testing.T) {
		testCmd := command.New("pkgr", "--version")
		res, err := testCmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}
		// this will always get the git tag for regular releases
		assert.Equal(t, "dev", string(res))
	})

	t.Run("short flag -v works", func(t *testing.T) {
		testCmd := command.New("pkgr", "-v")
		res, err := testCmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}
		// this will always get the git tag for regular releases
		assert.Equal(t, "dev", string(res))
	})

}
