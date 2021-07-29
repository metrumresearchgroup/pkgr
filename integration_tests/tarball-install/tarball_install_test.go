package tarball_install

import (
	"context"
	"github.com/metrumresearchgroup/command"
	. "github.com/metrumresearchgroup/pkgr/testhelper"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// Test IDs
const (
	tarballInstallE2ETest1 = "TRB-E2E-001"
	tarballInstallE2ETest2 = "TRB-E2E-002"
	tarballInstallE2ETest3 = "TRB-E2E-003"
	tarballInstallE2ETest4 = "TRB-E2E-004"
	tarballInstallE2ETest5 = "TRB-E2E-005"
)

// Golden file IDs
const (
	tarballBasicInstall = "tarball-basic-install"
)

//TODO
func TestTarballInstall(t *testing.T) {

	t.Run(MakeTestName(tarballInstallE2ETest1, "plan includes message about tarball installation"), func(t *testing.T) {
		DeleteTestFolder(t,"test-library")
		DeleteTestFolder(t, "test-cache")

		ctx := context.TODO()
		planCmd := command.New()

		capture, err := planCmd.Run(ctx, "pkgr", "plan", "--loglevel=debug", "--logjson")
		if err != nil {
			t.Fatalf("error running pkgr plan: %s", err)
		}
		logs1 := CollectGenericLogs(t, capture, "additional installation set")
		logs2 := CollectGenericLogs(t, capture, "package installation sources")

		assert.Len(t, logs1, 1, "expected exactly one 'additional installation set' log entry")

		assert.Equal(t, "test-cache/b64f2bee6438d2bd/R6", logs1[0].InstallFrom)
		assert.Equal(t, "tarball", logs1[0].Method)
		assert.Equal(t, "./tarballs/R6_2.4.0.tar.gz", logs1[0].Origin)
		assert.Equal(t, "R6", logs1[0].Pkg)

		assert.Len(t, logs2, 1, "expected exactly one 'package installation sources' log entry")
		assert.Equal(t, 7, logs2[0].LocalRepo)
		assert.Equal(t, 1, logs2[0].Tarballs)
	})

	t.Run(MakeTestName(tarballInstallE2ETest2, "installs from tarballs directly"), func(t *testing.T){
		DeleteTestFolder(t,"test-library")
		DeleteTestFolder(t, "test-cache")

		ctx := context.TODO()
		installCmd := command.New()
		rScriptCmd := command.New(command.WithDir("Rscripts"))

		_, err := installCmd.Run(ctx, "pkgr", "install")
		if err != nil {
			t.Fatalf("error during pkgr install: %s", err)
		}

		testCapture, err := rScriptCmd.Run(ctx,"Rscript", "--quiet", "install_test.R")
		if err != nil {
			t.Fatalf("error running Rscript to collect installed packages: %s", err)
		}
		g := goldie.New(t)
		g.Assert(t, tarballBasicInstall, testCapture.Output)
	})

	t.Run(MakeTestName(tarballInstallE2ETest3, "clean cleans the local cache"), func(t *testing.T) {
		DeleteTestFolder(t, "test-cache")
		SetupEndToEndWithInstall(t, "pkgr.yml", "test-library")

		ctx := context.TODO()
		cleanCmd := command.New()

		_, err := cleanCmd.Run(ctx, "pkgr", "clean", "cache")
		if err != nil {
			t.Fatalf("error running 'pkgr clean cache': %s", err)
		}

		assert.DirExists(t, "test-cache", "entire cache folder was removed instead of just contents")
		contents, err := os.ReadDir("test-cache")
		if err != nil {
			t.Fatalf("error while checking directory contents: %s", err)
		}
		assert.Len(t, contents, 0, "there are files remaining in the pkg cache after cleaning")

		//rScriptCmd := command.New(command.WithDir("Rscripts"))
	})
}