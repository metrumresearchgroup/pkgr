package baseline

import (
	"context"
	"github.com/metrumresearchgroup/command"
	. "github.com/metrumresearchgroup/pkgr/testhelper"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

// Test IDs
const(
	baselineInstallE2ETest1 = "BSLNINS-E2E-001"
	baselineInstallE2ETest2 = "BSLNINS-E2E-002"
	baselineInstallE2ETest3 = "BSLNINS-E2E-003"
	baselineInstallE2ETest4 = "BSLNINS-E2E-004"
	baselineInstallE2ETest5 = "BSLNINS-E2E-005"
	baselineInstallE2ETest6 = "BSLNINS-E2E-006"
)

// Golden file names
const(
	basicInstall = "install"
	suggestsInstall = "suggests-install"
	idempotenceInstall = "idempotence-install"
)

func TestInstall(t *testing.T) {
	DeleteTestFolder(t, "test-library")
	installCmd := command.New()
	ctx := context.TODO()
	installRes, err := installCmd.Run(ctx, "pkgr", "install")
	if err != nil {
		panic(err)
	}

	t.Run(MakeTestName(baselineInstallE2ETest1, "should install 11"), func(t *testing.T) {
		assert.Contains(t, string(installRes.Output), "to_install=11 to_update=0")
	})
	t.Run(MakeTestName(baselineInstallE2ETest2, "should install to the test library"), func(t *testing.T) {
		assert.Contains(t, string(installRes.Output), "Library path to install packages: test-library")
	})


	testCmd := command.New(command.WithDir("Rscripts"))
	// we could also suppress any global state by doing things like setting --vanilla
	// then setting the R_LIBS_SITE using command.WithEnv() however
	// this way it will be easy to interactively develop the script in the Rproj,
	// then invoke it whenever needed
	t.Run(MakeTestName(baselineInstallE2ETest3, "install the packages and dependencies"), func(t *testing.T) {
		testRes, err := testCmd.Run(ctx, "Rscript", "--quiet", "install_test.R")
		if err != nil {
		  	t.Fatalf("failed to run Rscript command with err: %s", err)
		}
		g := goldie.New(t)
		g.Assert(t, basicInstall, testRes.Output)
	})

	t.Run(MakeTestName(baselineInstallE2ETest1, "can install suggested dependencies for user packages"), func(t *testing.T){
		DeleteTestFolder(t, "test-library")
		installCmd := command.New()
		ctx := context.TODO()

		_, err := installCmd.Run(ctx, "pkgr", "install", "--config=pkgr-suggests.yml")
		if err != nil {
			t.Fatalf("error occurred while running pkgr install: %s", err)
		}

		rScriptCmd := command.New(command.WithDir("Rscripts"))
		testCapture, err := rScriptCmd.Run(ctx, "Rscript", "--quiet", "install_test.R")
		if err != nil {
			log.Fatalf("error occurred while using R to check installed packages: %s", err)
		}

		g := goldie.New(t)
		g.Assert(t, suggestsInstall, testCapture.Output)
	})
}

// Just making a separate function until I can refactor the first one into shape
func TestInstall2(t *testing.T) {
	t.Run(MakeTestName(baselineInstallE2ETest4, "installs are idempotent"), func(t *testing.T) {

		//Setup
		DeleteTestFolder(t, "test-cache")
		SetupEndToEndWithInstall(t, "pkgr.yml", "test-library")
		assert.DirExists(t, "test-library/fansi", "expected fansi to be installed, but couldn't find folder in test-library")
		DeleteTestFolder(t, "test-library/fansi") // giving pkgr the need to install this package again


		// Execute
		ctx := context.TODO()
		installCmd := command.New()
		rScriptCommand := command.New(command.WithDir("Rscripts"))

		capture, err := installCmd.Run(ctx, "pkgr", "install", "--config=pkgr.yml", "--logjson")
		if err != nil {
			t.Fatalf("error occurred running pkgr install: %s", err)
		}
		installPlanLogs := CollectGenericLogs(t, capture, "package installation plan")

		rScriptCapture, err := rScriptCommand.Run(ctx, "Rscript", "--quiet", "install_test.R")
		if err != nil {
			log.Fatalf("error occurred while using R to check installed packages: %s", err)
		}

		// Validate
		assert.Len(t, installPlanLogs, 1)
		assert.Equal(t, 1, installPlanLogs[0].ToInstall)
		assert.Equal(t, 0, installPlanLogs[0].ToUpdate)

		g := goldie.New(t)
		g.Assert(t, idempotenceInstall, rScriptCapture.Output)
	})
}
