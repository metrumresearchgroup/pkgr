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
	DeleteTestLibrary(testLibrary, t)
	ctx := context.TODO()
	installCmd := command.New()
	_, err := installCmd.Run(ctx, "pkgr", "install", fmt.Sprintf("--config=%s", pkgrConfig))
	if err != nil {
		t.Fatalf("could not install baseline packages to '%s' with config file '%s'. Error: %s", testLibrary, pkgrConfig, err)
	}
}

func DeleteTestLibrary(testLibrary string, t *testing.T) {
	err := os.RemoveAll(testLibrary)
	if err != nil {
		t.Fatalf("failed to clean up test library at '%s'. Error: %s", testLibrary, err)
	}
}

// Returns errors
// I don't like this because part of the reason I wanted to pull this into functions was to reduce boilerplate
// "if err != nil { t.Fatalf(...) } code, and returning the err sort of defeats that purpose.
// However, I'm including it for considation.
//func SetupEndToEndWithInstall(t *testing.T, pkgrConfig, testLibrary string) error {
//	err := DeleteTestLibrary(testLibrary, t)
//	if err != nil {
//		return err
//	}
//	ctx := context.TODO()
//	installCmd := command.New()
//	_, err = installCmd.Run(ctx, "pkgr", "install", fmt.Sprintf("--config=%s", pkgrConfig))
//	if err != nil {
//		t.Errorf("could not install baseline packages to '%s' with config file '%s'. Error: %s", testLibrary, pkgrConfig, err)
//		return err
//	}
//	return nil
//
//
//}
//
//func DeleteTestLibrary(testLibrary string, t *testing.T) error {
//	err := os.RemoveAll(testLibrary)
//	if err != nil {
//		t.Errorf("failed to clean up test library at '%s'. Error: %s", testLibrary, err)
//	}
//	return err
//}