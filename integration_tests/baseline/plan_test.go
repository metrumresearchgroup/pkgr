package baseline

import (
	"context"
	"github.com/metrumresearchgroup/command"
	. "github.com/metrumresearchgroup/pkgr/testhelper"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	baselinePlanTest1="BSLNPLN-E2E-001"
	baselinePlanTest2="BSLNPLN-E2E-002"
	baselinePlanTest3="BSLNPLN-E2E-003"
	baselinePlanTest4="BSLNPLN-E2E-004"
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

	t.Run(MakeTestName(baselinePlanTest3, "pkgr accurately reports package installation status and packages that were not installed by pkgr"), func(t *testing.T) {
		t.Skip("Test is currently skipped because it is failing. See issue #380")

		// Installs an outdated version of ellipsis, rlang
		DeleteTestFolder(t, "test-cache")
		SetupEndToEndWithInstall(t, "pkgr-preinstalled-setup.yml", "test-library")

		ctx := context.TODO()
		rInstallCmd := command.New()
		planCmd := command.New()

		rInstallCapture, err := rInstallCmd.Run(ctx, "Rscript", "-e", "install.packages(c('digest', 'R6'), lib='test-library', repos=c('https://mpn.metworx.com/snapshots/stable/2021-06-20'))")
		if err != nil {
			t.Fatalf("error while installing packages through non-pkgr means: %s\nOutput:\n%s", err, string(rInstallCapture.Output))
		}

		planCapture, err := planCmd.Run(ctx, "pkgr", "plan", "--config=pkgr-preinstalled.yml", "--logjson")
		if err != nil {
			t.Fatalf("error running pkgr plan: %s", err)
		}

		notInstalledByPkgrListLogs := CollectGenericLogs(t, planCapture, "Packages not installed by pkgr")
		assert.Len(t, notInstalledByPkgrListLogs, 1, "expected only one message containing all packages not installed by pkgr")
		assert.Len(t, notInstalledByPkgrListLogs[0].Packages, 2, "expected exactly two entries in the 'not installed by pkgr' log entry")
		assert.Contains(t, notInstalledByPkgrListLogs[0].Packages, "digest", "digest should have been listed under 'not installed by pkgr'")
		assert.Contains(t, notInstalledByPkgrListLogs[0].Packages, "R6", "R6 should have been listed under 'not installed by pkgr'")

		installationStatusLogs := CollectGenericLogs(t, planCapture, "package installation status")
		assert.Len(t, installationStatusLogs, 1, "expected exactly one message containing `package installation status` metadata")
		assert.Equal(t, 4, installationStatusLogs[0].Installed) // rlang and ellipsis from outdated pkgr.yml, digest (extraneous) and R6 from install.packages.
		assert.Equal(t, 2, installationStatusLogs[0].NotFromPkgr)
		assert.Equal(t, 2, installationStatusLogs[0].Outdated)
		assert.Equal(t, 11, installationStatusLogs[0].TotalPackagesRequired)

		packageInstallationPlanLogs := CollectGenericLogs(t, planCapture, "package installation plan")
		assert.Len(t, packageInstallationPlanLogs, 1, "expected exactly one 'package installation plan' message")
		assert.Equal(t, 8, packageInstallationPlanLogs[0].ToInstall, "expected plan to claim 8 package would be installed")
		assert.Equal(t, 2, packageInstallationPlanLogs[0].ToUpdate, "expected plan to claim 2 packages would be updated")

	})

}
