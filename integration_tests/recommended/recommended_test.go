package recommended

import (
	"context"
	"github.com/metrumresearchgroup/command"
	. "github.com/metrumresearchgroup/pkgr/testhelper"
	"github.com/stretchr/testify/assert"
	"testing"
)

const(
	recommendedE2ETest1 = "REC-E2E-001"
	recommendedE2ETest2 = "REC-E2E-002"
	recommendedE2ETest3 = "REC-E2E-003"
)

func TestRecommended(t *testing.T) {
	t.Run(MakeTestName(recommendedE2ETest1, "plan includes recommended packages by default"), func(t *testing.T) {
		DeleteTestFolder(t, "test-library")
		ctx := context.TODO()
		planCmd := command.New()

		capture, err := planCmd.Run(ctx, "pkgr", "plan", "--config=pkgr.yml", "--loglevel=debug", "--logjson")
		if err != nil {
			t.Fatalf("error running pkgr plan: %s", err)
		}

		pkgRepoLogs := CollectPkgRepoSetLogs(t, capture)
		assert.True(
			t,
			pkgRepoLogs.Contains("survival", "3.1-8", "MPN", "dependency"),
			"did not find recommended package 'survival' marked a dependency",
		)
	})

	t.Run(MakeTestName(recommendedE2ETest2, "plan excludes recommended packages when NoRecommended is true"), func(t *testing.T){
		DeleteTestFolder(t, "test-library")
		ctx := context.TODO()
		planCmd := command.New()

		capture, err := planCmd.Run(ctx, "pkgr", "plan", "--config=pkgr-no-recommended.yml", "--loglevel=debug", "--logjson")
		if err != nil {
			t.Fatalf("error running pkgr plan: %s", err)
		}

		pkgRepoLogs := CollectPkgRepoSetLogs(t, capture)
		assert.False(
			t,
			pkgRepoLogs.Contains("survival", "3.1-8", "MPN", "dependency"),
			"found recommended package 'survival' marked a dependency when NoRecommended was set",
		)

	})

	t.Run(MakeTestName(recommendedE2ETest3, "specifying recommeded pkg as a user package overrides NoRecommeded option"), func(t *testing.T){
		DeleteTestFolder(t, "test-library")
		ctx := context.TODO()
		planCmd := command.New()

		capture, err := planCmd.Run(ctx, "pkgr", "plan", "--config=pkgr-survival-direct.yml", "--loglevel=debug", "--logjson")
		if err != nil {
			t.Fatalf("error running pkgr plan: %s", err)
		}

		pkgRepoLogs := CollectPkgRepoSetLogs(t, capture)
		assert.True(
			t,
			pkgRepoLogs.Contains("survival", "3.1-8", "MPN", "user_defined"),
			"did not find recommended package 'survival' in install plan",
		)

	})
}