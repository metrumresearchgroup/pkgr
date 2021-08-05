package baseline

import (
	"context"
	"github.com/metrumresearchgroup/command"
	. "github.com/metrumresearchgroup/pkgr/testhelper"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const (
	baselinePlanTest1="BSLNPLN-E2E-001"
	baselinePlanTest2="BSLNPLN-E2E-002"
	baselinePlanTest3="BSLNPLN-E2E-003"
)

// Golden Test File Names
const (
	baselinePlanGolden="baseline-plan-golden"
)

func TestPlan(t *testing.T) {



	t.Run(MakeTestName(baselinePlanTest1, "plan indicates packages to be installed, as well as version, source repo, and whether pkg is user-defined or a dependency"), func(t *testing.T) {
		DeleteTestFolder(t, "test-library")
		ctx := context.TODO()
		planCmd := command.New()

		capture, err := planCmd.Run(ctx, "pkgr", "plan", "--loglevel=debug", "--logjson")
		if err != nil {
			t.Fatalf("error running pkgr plan: %s", err)
		}
		//output := string(capture.Output)


		pkgRepoSetLogs := CollectPkgRepoSetLogs(t, capture)
		//assert.True(t, pkgRepoSetLogs.Contains("R6", "2.5.0", "CRAN", "user_defined"), "failed to find expected log message")
		//assert.True(t, pkgRepoSetLogs.Contains("pillar", "1.6.1", "CRAN", "user_defined"), "failed to find expected log message")
		//assert.True(t, pkgRepoSetLogs.Contains("glue", "1.4.2", "CRAN", "dependency"), "failed to find expected log message")
		//assert.True(t, pkgRepoSetLogs.Contains("fansi", "0.5.0", "CRAN", "dependency"), "failed to find expected log message")
		//assert.True(t, pkgRepoSetLogs.Contains("rlang", "0.4.11", "CRAN", "dependency"), "failed to find expected log message")
		//assert.True(t, pkgRepoSetLogs.Contains("utf8", "1.2.1", "CRAN", "dependency"), "failed to find expected log message")
		//assert.True(t, pkgRepoSetLogs.Contains("crayon", "1.4.1", "CRAN", "dependency"), "failed to find expected log message")
		//assert.True(t, pkgRepoSetLogs.Contains("lifecycle", "1.0.0", "CRAN", "dependency"), "failed to find expected log message")
		//assert.True(t, pkgRepoSetLogs.Contains("vctrs", "0.3.8", "CRAN", "dependency"), "failed to find expected log message")
		//assert.True(t, pkgRepoSetLogs.Contains("ellipsis", "0.3.2", "CRAN", "dependency"), "failed to find expected log message")
		//assert.True(t, pkgRepoSetLogs.Contains("cli", "2.5.0", "CRAN", "dependency"), "failed to find expected log message")

		//alternative
		g := goldie.New(t)

		g.Assert(t, baselinePlanGolden, (pkgRepoSetLogs.ToBytes()))
	})

	t.Run(MakeTestName(baselinePlanTest2, "number of workers (threads) can be set"), func(t *testing.T) {
		DeleteTestFolder(t, "test-library")
		ctx := context.TODO()
		planCmd := command.New()

		capture, err := planCmd.Run(ctx, "pkgr", "plan", "--threads=5", "--logjson")
		if err != nil {
			t.Fatalf("error running pkgr plan: %s", err)
		}
		output := string(capture.Output)

		jsonRegex := `\{"level":"info","msg":"Installation would launch 5 workers.*\}`
		assert.Regexp(t, jsonRegex, output)
	})

	t.Run(MakeTestName(baselinePlanTest3, "pkgr flags packages that were not installed by pkgr"), func(t *testing.T) {
		DeleteTestFolder(t, "test-library")
		err := os.MkdirAll("./test-library", 0777)
		if err != nil {
			t.Fatalf("error while creating empty test-library directory: %s", err)
		}

		ctx := context.TODO()
		rInstallCmd := command.New()
		planCmd := command.New()

		rScriptCapture, err := rInstallCmd.Run(ctx, "Rscript", "-e", "install.packages('ellipsis', lib='test-library', repos=c('https://mpn.metworx.com/snapshots/stable/2021-06-20'))")
		if err != nil {
			t.Fatalf("error while installing packages through non-pkgr means: %s\nOutput:\n%s", err, string(rScriptCapture.Output))
		}

		planCapture, err := planCmd.Run(ctx, "pkgr", "plan", "--logjson")
		if err != nil {
			t.Fatalf("error running pkgr plan: %s", err)
		}

		logs := CollectGenericLogs(t, planCapture, "Packages not installed by pkgr")
		assert.Len(t, logs, 1, "expected only one message containing all packages not installed by pkgr")
		assert.Len(t, logs[0].Packages, 2, "expected exactly two entries in the 'not installed by pkgr' log entry")
		assert.Contains(t, logs[0].Packages, "ellipsis", "ellipsis should have been listed under 'not installed by pkgr'")
		assert.Contains(t, logs[0].Packages, "rlang", "rlang should have been listed under 'not installed by pkgr'")
	})

}
