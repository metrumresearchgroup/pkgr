package multi_repo

import (
	"context"
	"github.com/metrumresearchgroup/command"
	. "github.com/metrumresearchgroup/pkgr/testhelper"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func setupMultiRepoTest(t *testing.T) {
	err := os.RemoveAll("test-library")
	if err != nil {
		t.Fatalf("failed to remove test library at beginning of test, %s", err)
	}
}

// Test IDs
const(
	multiRepoE2ETest1 = "MRPO-E2E-001"
	multiRepoE2ETest2 = "MRPO-E2E-002"
	multiRepoE2ETest3 = "MRPO-E2E-003"
	multiRepoE2ETest4 = "MRPO-E2E-004"
	multiRepoE2ETest5 = "MRPO-E2E-005"
)

// Golden file names
const (
	multiRepoInstallation = "multi-repo-installed-packages"
)

func TestMultiRepoInstall(t *testing.T) {
	t.Run(MakeTestName(multiRepoE2ETest1, "pkgr plan takes packages from both local and remote repos in the order listed in pkgr.yml"), func(t *testing.T) {
		setupMultiRepoTest(t)

		ctx := context.TODO()
		planCmd := command.New()

		outputBytes, err := planCmd.Run(ctx, "pkgr", "plan", "--loglevel=debug", "--logjson")
		if err != nil {
			t.Fatalf("error occurred when installing packages: %s", err)
		}
		output := string(outputBytes.Output)

		// Check repositories set correctly
		//# R6 -- Install from local. Used to ensure a user package can be installed from local, ensure repo-order is respected when installing user packages.
		//# ellipsis -- Install from remote. Used to ensure a user package can be installed from remote. All of this package's deps will come from remote as well.
		//# rlang Install from remote. Used to ensure a dependency package can be installed from remote (dep. of ellipsis)
		//# openssl -- Install from remote. Used to ensure that a user package can be installed from remote
		//# askpass -- Install from local. Used to ensure that repo order is respected when installing dependencies (askpass is a dependency of openssl. Openssl is installed from second repo, askpass should be installed from first repo, i.e. local.)
		//# sys -- Install from local. Used to ensure a dependency package can be installed from local (dep. of askpass).
		r6Regex := `\{"level":"debug","msg":"package repository set","pkg":"R6","relationship":"user_defined","repo":"LOCALREPO","type":1,"version":"2.5.0"\}`
		ellipsisRegex := `\{"level":"debug","msg":"package repository set","pkg":"ellipsis","relationship":"user_defined","repo":"REMOTEREPO","type":1,"version":"0.3.2"\}`
		rlangRegex := `\{"level":"debug","msg":"package repository set","pkg":"rlang","relationship":"dependency","repo":"REMOTEREPO","type":1,"version":"0.4.11"\}`
		opensslRegex := `\{"level":"debug","msg":"package repository set","pkg":"openssl","relationship":"user_defined","repo":"REMOTEREPO","type":1,"version":"1.4.4"\}`
		askpassRegex := `\{"level":"debug","msg":"package repository set","pkg":"askpass","relationship":"dependency","repo":"LOCALREPO","type":1,"version":"1.1"\}`
		sysRegex := `\{"level":"debug","msg":"package repository set","pkg":"sys","relationship":"dependency","repo":"LOCALREPO","type":1,"version":"3.4"\}`

		assert.Regexp(t, r6Regex, output)
		assert.Regexp(t, ellipsisRegex, output)
		assert.Regexp(t, rlangRegex, output)
		assert.Regexp(t, opensslRegex, output)
		assert.Regexp(t, askpassRegex, output)
		assert.Regexp(t, sysRegex, output)
	})

	t.Run(MakeTestName(multiRepoE2ETest2, "pkgr install can install packages from multiple repositories"), func(t *testing.T) {

		ctx := context.TODO()
		installCmd := command.New()
		rScriptCmd := command.New(command.WithDir("Rscripts"))
		_, err := installCmd.Run(ctx, "pkgr", "install", "--loglevel=debug", "--logjson")
		if err != nil {
			t.Fatalf("error occurred when installing packages: %s", err)
		}

		rScriptOutputBytes, err := rScriptCmd.Run(ctx, "Rscript", "--quiet", "install_test.R")
		//t.Log(string(rScriptOutputBytes.Output))
		if err != nil {
			t.Fatalf("error occurred while detecting installed packages: %s", err)
		}

		g := goldie.New(t)
		g.Assert(t, multiRepoInstallation, rScriptOutputBytes.Output)
	})

}