package env_vars

import (
	"context"
	"github.com/metrumresearchgroup/command"
	. "github.com/metrumresearchgroup/pkgr/testhelper"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Test IDs
const(
	envVarE2Etest1 = "ENV-E2E-001"
	envVarE2Etest2 = "ENV-E2E-002"
	envVarE2Etest3 = "ENV-E2E-003"
	envVarE2Etest4 = "ENV-E2E-004"
)

// Golden Test Files
const (
	EnvVarRpathInstall = "env-rpath-installed-pkgs"
)

func TestRpathEnvVar(t *testing.T) {
	t.Run(MakeTestName(envVarE2Etest1, "PKGR_RPATH environment variable is set in plan"), func(t *testing.T){

		// Setup
		rPathSymlink := GetTestSymlinkPath(t)
		DeleteTestFile(t, rPathSymlink)
		DeleteTestFolder(t, "test-library")
		validRPath := GetValidRPath(t)
		err := os.Symlink(validRPath, rPathSymlink)
		if err != nil {
			t.Fatalf("failed to create symlink to RPath for testing purposes: %s", err)
		}


		os.Setenv("PKGR_RPATH", rPathSymlink)

		// Run
		ctx := context.TODO()
		planCmd := command.New()
		capture, err := planCmd.Run(ctx, "pkgr", "plan", "--loglevel=trace", "--logjson")
		if err != nil {
			t.Fatalf("error when executing 'pkgr plan': %s", err)
		}
		logs := CollectGenericLogs(t, capture, "command args")

		// Validate
		for _, log := range logs {
			assert.Equal(t, rPathSymlink, log.RPath, "R Path to use does not match R Path set in env variable")
			assert.Equal(t, log.RSettings.Rpath, rPathSymlink)
		}
	})

	t.Run(MakeTestName(envVarE2Etest2, "installation successful with valid PKGR_RPATH environment variable set"), func(t *testing.T) {
		//Setup
		rPathSymlink := GetTestSymlinkPath(t)
		DeleteTestFile(t, rPathSymlink)
		DeleteTestFolder(t, "test-library")
		validRPath := GetValidRPath(t)
		err := os.Symlink(validRPath, rPathSymlink)
		if err != nil {
			t.Fatalf("failed to create symlink to RPath for testing purposes: %s", err)
		}

		os.Setenv("PKGR_RPATH", rPathSymlink)

		// Run
		ctx := context.TODO()
		installCmd := command.New()
		capture, err := installCmd.Run(ctx, "pkgr", "install", "--loglevel=trace", "--logjson")
		if err != nil {
			t.Fatalf("error when executing 'pkgr install': %s", err)
		}
		logs := CollectGenericLogs(t, capture, "command args")
		rCmd := command.New(command.WithDir("Rscripts"))
		rCmdCapture, err := rCmd.Run(ctx, "Rscript", "--quiet", "install_test.R")
		if err != nil {
			t.Fatalf("error when running Rscript to parse installed packages: %s", err)
		}

		// Validate
		for _, log := range logs {
			assert.Equal(t, rPathSymlink, log.RPath, "R Path to use does not match R Path set in env variable for log entry: %v", log)
			assert.Equal(t, log.RSettings.Rpath, rPathSymlink)
		}

		g := goldie.New(t)
		g.Assert(t, EnvVarRpathInstall, rCmdCapture.Output)
	})

	t.Run(MakeTestName(envVarE2Etest3, "installation fails with invalid PKGR_RPATH environment variable set"), func(t *testing.T) {
		//Setup
		DeleteTestFolder(t, "test-library")
		invalidRPath := "/home/FAKE_RPATH"

		os.Setenv("PKGR_RPATH", invalidRPath)
		// Run
		ctx := context.TODO()
		installCmd := command.New()
		capture, testError := installCmd.Run(ctx, "pkgr", "install", "--loglevel=trace", "--logjson")

		// Validate
		assert.Error(t, testError, "pkgr should have thrown an error when trying to install with an invalid RPath")
		assert.NoDirExists(t, "test-library", "test-library should not have been created since pkgr should have failed")
		logs := CollectGenericLogs(t, capture, "command args")
		for _, log := range logs {
			assert.Equal(t, invalidRPath, log.RPath, "R Path to use does not match R Path set in env variable")
			assert.Equal(t, log.RSettings.Rpath, invalidRPath)
		}
	})
}

func GetValidRPath(t *testing.T) string {
	ctx := context.TODO()
	utilCommand := command.New()
	rPathCapture, err := utilCommand.Run(ctx, "which", "R")
	if err != nil {
		t.Fatalf("error while looking for valid R installation on the PATH: %s", err)
	}
	validRPath := strings.TrimSpace(string(rPathCapture.Output))
	validRPath = filepath.Clean(validRPath)
	return validRPath
}

func GetTestSymlinkPath(t *testing.T) string {
	rpathSymlink, err := filepath.Abs("./rpath_symlink")
	if err != nil {
		t.Fatalf("could not get absolute path to local symlink for testing: %s", err)
	}
	return rpathSymlink
}