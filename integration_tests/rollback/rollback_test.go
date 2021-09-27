package rollback_test

import (
	"os"
	"testing"

	"github.com/metrumresearchgroup/command"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"

	"github.com/metrumresearchgroup/pkgr/testhelper"
)

const (
	rollbackE2ETest1 = "RBK-E2E-001"
	rollbackE2ETest2 = "RBK-E2E-002"
	rollbackE2ETest3 = "RBK-E2E-003"
	rollbackE2ETest4 = "RBK-E2E-004"
	rollbackE2ETest5 = "RBK-E2E-005"
	rollbackE2ETest6 = "RBK-E2E-006"
)

const (
	BaselineInstalled         = "baseline-installed"
	RollbackDisabledInstalled = "rollback-disabled-installed"
)

func setupBaseline(t *testing.T) {
	err := os.RemoveAll("test-library")
	if err != nil {
		t.Fatalf("failed to cleanup test-library")
	}

	installCmd := command.New("pkgr", "install", "--config=pkgr-baseline.yml")
	err = installCmd.Run()
	if err != nil {
		t.Fatalf("could not install baseline packages with err: %s", err)
	}
}

func TestRollback(t *testing.T) {
	// should have some setup to make sure the test-library is cleared out

	setupBaseline(t)

	// this test is really just a sanity check to make sure the baseline is really set up properly to reflect the
	// future tests. This feels like a reasonable middleground to checking
	// its set up like that after each setupBaseline() call, especially since
	// setupBaseline doesn't do any assertions/checks within it
	t.Run(testhelper.MakeTestName(rollbackE2ETest1, "the baseline package was installed"), func(t *testing.T) {
		testCmd := command.New("Rscript", "--quiet", "install_test.R")
		testCmd.Dir = "Rscripts"

		testRes, err := testCmd.CombinedOutput()
		if err != nil {
			t.Fatalf("failed to run Rscript command with err: %s", err)
		}
		g := goldie.New(t)
		g.Assert(t, "baseline-installed", testRes)
	})

	t.Run(testhelper.MakeTestName(rollbackE2ETest2, "will rollback on failure to install tarball"), func(t *testing.T) {
		installCmd := command.New("pkgr", "install", "--config=pkgr-rollback-tarball.yml", "--logjson")
		res, err := installCmd.CombinedOutput()
		if err != nil {
			t.Fatalf("could not install baseline packages with err: %s", err)
		}
		rollbackOutputCheckHelper(t, res, BaselineInstalled)
	})

	t.Run(testhelper.MakeTestName(rollbackE2ETest3, "will not rollback on failure to install tarball when rollback disabled"), func(t *testing.T) {
		t.Run(testhelper.MakeSubtestName(rollbackE2ETest3, "A", "in configuration file"), func(t *testing.T) {
			setupBaseline(t)
			installCmd := command.New("pkgr", "install", "--config=pkgr-no-rollback-tarball.yml", "--logjson")
			res, err := installCmd.CombinedOutput()
			if err != nil {
				t.Fatalf("could not install baseline packages with err: %s", err)
			}
			rollbackOutputCheckHelper(t, res, RollbackDisabledInstalled)
		})

		t.Run(testhelper.MakeSubtestName(rollbackE2ETest3, "B", "as CLI flag"), func(t *testing.T) {
			setupBaseline(t)
			installCmd := command.New("pkgr", "install", "--config=pkgr-rollback-tarball.yml", "--logjson", "--no-rollback")
			res, err := installCmd.CombinedOutput()
			if err != nil {
				t.Fatalf("could not install baseline packages with err: %s", err)
			}
			rollbackOutputCheckHelper(t, res, RollbackDisabledInstalled)
		})
	})
}

func rollbackOutputCheckHelper(t *testing.T, capture []byte, goldenName string) {
	jsonRes := `{"level":"info","msg":"Successfully Installed.","package":"ps","remaining":1,"repo":"MPN","version":"1.6.0"}`
	assert.Contains(t, string(capture), jsonRes)

	testCmd := command.New("Rscript", "--quiet", "install_test.R")
	testCmd.Dir = "Rscripts"

	testRes, err := testCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to run Rscript command with err: %s", err)
	}
	g := goldie.New(t)
	// should be the same values as baseline
	g.Assert(t, goldenName, testRes)
}
