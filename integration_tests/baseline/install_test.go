package baseline

import (
	"log"
	"os"
	"testing"

	"github.com/metrumresearchgroup/command"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"

	. "github.com/metrumresearchgroup/pkgr/testhelper"
)

// Test IDs
const (
	baselineInstallE2ETest1 = "BSLNINS-E2E-001"
	baselineInstallE2ETest2 = "BSLNINS-E2E-002"
	baselineInstallE2ETest3 = "BSLNINS-E2E-003"
	baselineInstallE2ETest4 = "BSLNINS-E2E-004"
	baselineInstallE2ETest5 = "BSLNINS-E2E-005"
	baselineInstallE2ETest6 = "BSLNINS-E2E-006"
)

// Golden file names
const (
	basicInstall        = "install"
	suggestsInstall     = "suggests-install"
	idempotenceInstall  = "idempotence-install"
	preinstalledInstall = "preinstalled-install"
)

func TestInstall(t *testing.T) {
	DeleteTestFolder(t, "test-library")
	installCmd := command.New("pkgr", "install")
	installRes, err := installCmd.CombinedOutput()
	if err != nil {
		panic(err)
	}

	t.Run(MakeTestName(baselineInstallE2ETest1, "should install 11"), func(t *testing.T) {
		assert.Contains(t, string(installRes), "to_install=11 to_update=0")
	})
	t.Run(MakeTestName(baselineInstallE2ETest2, "should install to the test library"), func(t *testing.T) {
		assert.Contains(t, string(installRes), "Library path to install packages: test-library")
	})

	// we could also suppress any global state by doing things like setting --vanilla
	// then setting the R_LIBS_SITE using command.WithEnv() however
	// this way it will be easy to interactively develop the script in the Rproj,
	// then invoke it whenever needed
	t.Run(MakeTestName(baselineInstallE2ETest3, "install the packages and dependencies"), func(t *testing.T) {
		testCmd := command.New("Rscript", "--quiet", "install_test.R")
		testCmd.Dir = "Rscripts"
		testRes, err := testCmd.CombinedOutput()
		if err != nil {
			t.Fatalf("failed to run Rscript command with err: %s", err)
		}
		g := goldie.New(t)
		g.Assert(t, basicInstall, testRes)
	})

	t.Run(MakeTestName(baselineInstallE2ETest4, "can install suggested dependencies for user packages"), func(t *testing.T) {
		DeleteTestFolder(t, "test-library")
		installCmd := command.New("pkgr", "install", "--config=pkgr-suggests.yml")

		err := installCmd.Run()
		if err != nil {
			t.Fatalf("error occurred while running pkgr install: %s", err)
		}

		rScriptCmd := command.New("Rscript", "--quiet", "install_test.R")
		rScriptCmd.Dir = "Rscripts"
		testCapture, err := rScriptCmd.CombinedOutput()
		if err != nil {
			log.Fatalf("error occurred while using R to check installed packages: %s", err)
		}

		g := goldie.New(t)
		g.Assert(t, suggestsInstall, testCapture)
	})
}

// Just making a separate function until I can refactor the first one into shape
func TestInstall2(t *testing.T) {
	t.Run(MakeTestName(baselineInstallE2ETest5, "installs are idempotent"), func(t *testing.T) {

		//Setup
		DeleteTestFolder(t, "test-cache")
		SetupEndToEndWithInstall(t, "pkgr.yml", "test-library")
		assert.DirExists(t, "test-library/fansi", "expected fansi to be installed, but couldn't find folder in test-library")
		DeleteTestFolder(t, "test-library/fansi") // giving pkgr the need to install this package again

		// Execute
		installCmd := command.New("pkgr", "install", "--config=pkgr.yml", "--logjson")
		installCmd.Dir = "Rscripts"

		capture, err := installCmd.CombinedOutput()
		if err != nil {
			t.Fatalf("error occurred running pkgr install: %s", err)
		}
		installPlanLogs := CollectGenericLogs(t, capture, "package installation plan")

		rScriptCommand := command.New("Rscript", "--quiet", "install_test.R")
		rScriptCommand.Dir = "Rscripts"
		rScriptCapture, err := rScriptCommand.CombinedOutput()
		if err != nil {
			log.Fatalf("error occurred while using R to check installed packages: %s", err)
		}

		// Validate
		assert.Len(t, installPlanLogs, 1)
		assert.Equal(t, 1, installPlanLogs[0].ToInstall)
		assert.Equal(t, 0, installPlanLogs[0].ToUpdate)

		g := goldie.New(t)
		g.Assert(t, idempotenceInstall, rScriptCapture)
	})

	t.Run(MakeTestName(baselineInstallE2ETest6, "installs do not overwrite packages not installed by pkgr"), func(t *testing.T) {
		//Setup
		DeleteTestFolder(t, "test-cache")
		DeleteTestFolder(t, "test-library")
		err := os.MkdirAll("./test-library", 0777)
		if err != nil {
			t.Fatalf("error while creating empty test-library directory: %s", err)
		}

		// Execute
		rInstallCmd := command.New("Rscript", "-e", "install.packages(c('ellipsis', 'digest'), lib='test-library', repos=c('https://mpn.metworx.com/snapshots/stable/2021-06-20'))")
		rInstallCapture, err := rInstallCmd.CombinedOutput()
		if err != nil {
			t.Fatalf("error while installing packages through non-pkgr means: %s\nOutput:\n%s", err, string(rInstallCapture))
		}

		installCmd := command.New("pkgr", "install", "--config=pkgr.yml", "--logjson")
		pkgrCapture, err := installCmd.CombinedOutput()
		if err != nil {
			t.Fatalf("error occurred running pkgr install: %s", err)
		}

		// Validate
		notInstalledByPkgLogs := CollectGenericLogs(t, pkgrCapture, "Packages not installed by pkgr")
		assert.Len(t, notInstalledByPkgLogs, 1, "expected only one message containing all packages not installed by pkgr")
		assert.Len(t, notInstalledByPkgLogs[0].Packages, 3, "expected exactly three entries in the 'not installed by pkgr' log entry")
		assert.Contains(t, notInstalledByPkgLogs[0].Packages, "ellipsis", "ellipsis should have been listed under 'not installed by pkgr'")
		assert.Contains(t, notInstalledByPkgLogs[0].Packages, "rlang", "rlang should have been listed under 'not installed by pkgr'")
		assert.Contains(t, notInstalledByPkgLogs[0].Packages, "digest", "digest should have been listed under 'not installed by pkgr'") // Extraneous, shouldn't matter to pkgr.yml

		installPlanLogs := CollectGenericLogs(t, pkgrCapture, "package installation plan")
		assert.Len(t, installPlanLogs, 1, "expected exactly one message containing 'package installation plan' metadata")
		assert.Equal(t, 9, installPlanLogs[0].ToInstall, "since two packages were already installed from pkgr's plan, we expect only 9 packages to be installed")
		assert.Equal(t, 0, installPlanLogs[0].ToUpdate)

		installationStatusLogs := CollectGenericLogs(t, pkgrCapture, "package installation status")
		assert.Len(t, installationStatusLogs, 1, "expected exactly one message containing `package installation status` metadata")
		assert.Equal(t, 3, installationStatusLogs[0].Installed)
		assert.Equal(t, 3, installationStatusLogs[0].NotFromPkgr)
		assert.Equal(t, 0, installationStatusLogs[0].Outdated)
		assert.Equal(t, 11, installationStatusLogs[0].TotalPackagesRequired)

		verifyInstalledCommand := command.New("Rscript", "--quiet", "install_test.R")
		verifyInstalledCommand.Dir = "Rscripts"
		rScriptCapture, err := verifyInstalledCommand.CombinedOutput()
		if err != nil {
			log.Fatalf("error occurred while using R to check installed packages: %s", err)
		}

		g := goldie.New(t)
		g.Assert(t, preinstalledInstall, rScriptCapture)
	})

}
