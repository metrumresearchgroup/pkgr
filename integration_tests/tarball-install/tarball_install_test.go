package tarball_install

import (
	"os"
	"testing"

	"github.com/metrumresearchgroup/command"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"

	. "github.com/metrumresearchgroup/pkgr/testhelper"
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
	tarballBasicInstall           = "tarball-basic-install"
	tarballOverwriteInstallBefore = "tarball-overwrite-install-before"
	tarballOverwriteInstallAfter  = "tarball-overwrite-install-after"
)

// TODO
func TestTarballInstall(t *testing.T) {

	t.Run(MakeTestName(tarballInstallE2ETest1, "plan includes message about tarball installation and gets tarball deps from repo"), func(t *testing.T) {
		DeleteTestFolder(t, "test-library")
		DeleteTestFolder(t, "test-cache")

		planCmd := command.New("pkgr", "plan", "--loglevel=debug", "--logjson")

		capture, err := planCmd.CombinedOutput()
		if err != nil {
			t.Fatalf("error running pkgr plan: %s", err)
		}
		pkgRepoLogs := CollectPkgRepoSetLogs(t, capture)
		// assert.True(t, pkgRepoLogs.Contains("ellipsis", "0.3.2", "local_tarball", "user_defined")) // This is covered in "additional installation set" logs.

		// Note: even though rlang is a dependency of ellipsis, pkgr currently flags it as "user-defined."
		// this is the result of a "shortcut" we took to ensure that tarball deps. were installed before the actual tarballs
		// While this is not the "desired" behavior, for now, it is the "expected" behavior.
		assert.True(t, pkgRepoLogs.Contains("rlang", "0.3.4", "LOCALREPO", "user_defined")) // Dependency of ellipsis
		assert.True(t, pkgRepoLogs.Contains("crayon", "1.3.4", "LOCALREPO", "user_defined"))

		logs1 := CollectGenericLogs(t, capture, "additional installation set")
		assert.Len(t, logs1, 1, "expected exactly one 'additional installation set' log entry")
		assert.Equal(t, "test-cache/c453b4514da3717a/ellipsis", logs1[0].InstallFrom)
		assert.Equal(t, "tarball", logs1[0].Method)
		assert.Equal(t, "./tarballs/ellipsis_0.3.2.tar.gz", logs1[0].Origin)
		assert.Equal(t, "ellipsis", logs1[0].Pkg)

		logs2 := CollectGenericLogs(t, capture, "package installation sources")
		assert.Len(t, logs2, 1, "expected exactly one 'package installation sources' log entry")
		assert.Equal(t, 2, logs2[0].LocalRepo)
		assert.Equal(t, 1, logs2[0].Tarballs)
	})

	t.Run(MakeTestName(tarballInstallE2ETest2, "installs from tarball directly"), func(t *testing.T) {
		DeleteTestFolder(t, "test-library")
		DeleteTestFolder(t, "test-cache")

		installCmd := command.New("pkgr", "install")

		err := installCmd.Run()
		if err != nil {
			t.Fatalf("error during pkgr install: %s", err)
		}

		rScriptCmd := command.New("Rscript", "--quiet", "install_test.R")
		rScriptCmd.Dir = "Rscripts"

		testCapture, err := rScriptCmd.CombinedOutput()
		if err != nil {
			t.Fatalf("error running Rscript to collect installed packages: %s", err)
		}
		g := goldie.New(t)
		g.Assert(t, tarballBasicInstall, testCapture)
	})

	t.Run(MakeTestName(tarballInstallE2ETest3, "overwrites existing package with tarball"), func(t *testing.T) {
		DeleteTestFolder(t, "test-cache")
		SetupEndToEndWithInstall(t, "pkgr-setup-old-ellipsis.yml", "test-library")

		g := goldie.New(t)

		rScriptCmdBefore := command.New("Rscript", "--quiet", "install_test.R")
		rScriptCmdBefore.Dir = "Rscripts"
		preCapture, err := rScriptCmdBefore.CombinedOutput()
		if err != nil {
			t.Fatalf("error running Rscript to collect installed packages: %s", err)
		}
		g.Assert(t, tarballOverwriteInstallBefore, preCapture)

		installCmd := command.New("pkgr", "install")
		err = installCmd.Run()
		if err != nil {
			t.Fatalf("error during pkgr install: %s", err)
		}

		rScriptCmdAfter := command.New("Rscript", "--quiet", "install_test.R")
		rScriptCmdAfter.Dir = "Rscripts"

		postCapture, err := rScriptCmdAfter.CombinedOutput()
		if err != nil {
			t.Fatalf("error running Rscript to collect installed packages: %s", err)
		}

		g.Assert(t, tarballOverwriteInstallAfter, postCapture)
	})

	t.Run(MakeTestName(tarballInstallE2ETest4, "clean cleans the local cache"), func(t *testing.T) {
		DeleteTestFolder(t, "test-cache")
		SetupEndToEndWithInstall(t, "pkgr.yml", "test-library")

		cleanCmd := command.New("pkgr", "clean", "cache")

		err := cleanCmd.Run()
		if err != nil {
			t.Fatalf("error running 'pkgr clean cache': %s", err)
		}

		assert.DirExists(t, "test-cache", "entire cache folder was removed instead of just contents")
		contents, err := os.ReadDir("test-cache")
		if err != nil {
			t.Fatalf("error while checking directory contents: %s", err)
		}
		assert.Len(t, contents, 0, "there are files remaining in the pkg cache after cleaning")
	})
}
