package testhelper

import (
	"context"
	"fmt"
	"github.com/metrumresearchgroup/command"
	"os"
	"testing"
)

func MakeTestName(testId, testName string) string {
	return(fmt.Sprintf("[%s] %s", testId, testName))
}

func MakeSubtestName(testId, subtestId, testName string) string {
	return(fmt.Sprintf("[%s-%s] %s", testId, subtestId, testName))
}

func SetupEndToEndWithInstall(t *testing.T, pkgrConfig, testLibrary string) {
	err := os.RemoveAll(testLibrary)
	if err != nil {
		t.Fatalf("failed to clean up test library at '%s'. Error: %s", testLibrary, err)
	}
	ctx := context.TODO()

	installCmd := command.New()
	_, err = installCmd.Run(ctx, "pkgr", "install", fmt.Sprintf("--config=%s", pkgrConfig))
	if err != nil {
		t.Fatalf("could not install baseline packages to '%s' with config file '%s'. Error: %s", testLibrary, pkgrConfig, err)
	}
}