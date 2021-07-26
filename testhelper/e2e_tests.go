package testhelper

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/metrumresearchgroup/command"
	"os"
	"strings"
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

type PkgrJsonLogOutput struct {
	Logs []string
}

func (lo PkgrJsonLogOutput) Collect(results command.Capture) {
	stringOutput := string(results.Output)
	lo.Logs = strings.Split(stringOutput, "\n")
}

func (lo PkgrJsonLogOutput) Contains(search string) ([]string, bool, error) {
	type msgField struct {
		msg string
	}

	matched := []string{}
	found := false

	for _, jsonLog := range lo.Logs {
		var msg msgField
		err := json.Unmarshal([]byte(jsonLog), &msg)
		if err != nil {
			return []string{}, found, err
		}

		if strings.Contains(msg.msg, search) {
			matched = append(matched, msg.msg)
			found = true
		}
	}
	return matched, found, nil
}




// Helper objects  for "package repository set" logs. ------------------------------------------------------------------
type PkgRepoSetMsg struct {
	//level string
	Msg string
	Pkg string
	Relationship string
	Repo string
	//TypE int
	Version string
}

type PkgRepoSetMsgCollection struct {
	Logs []PkgRepoSetMsg
}

func (prsmc PkgRepoSetMsgCollection) Contains(pkg, version, repo, relationship string) bool {
	for _, log := range prsmc.Logs {
		if log.Pkg == pkg && log.Version == version && log.Repo == repo && log.Relationship == relationship {
			return true
		}
	}

	return false
}

func CollectPkgRepoSetLogs(t *testing.T, capture command.Capture) PkgRepoSetMsgCollection {

	parsedLines := []PkgRepoSetMsg{}

	msgKey := "package repository set"
	outputLines := strings.Split(string(capture.Output), "\n")

	for _, line := range outputLines {
		if strings.Contains(line, msgKey) {
			var parsedLine PkgRepoSetMsg
			err := json.Unmarshal([]byte(line), &parsedLine)
			if err != nil {
				t.Fatalf("error unmarshalling the following JSON line: '%s'. error was: %s", line, err)
			}
			parsedLines = append(parsedLines, parsedLine)
		}
	}

	return PkgRepoSetMsgCollection{
			Logs : parsedLines,
		}


}
// ---------------------------------------------------------------------------------------------------------------------




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