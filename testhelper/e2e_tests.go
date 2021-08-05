package testhelper

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/metrumresearchgroup/command"
	"github.com/metrumresearchgroup/pkgr/cran"
	"github.com/stretchr/testify/assert"
	"os"
	"regexp"
	"sort"
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
	DeleteTestFolder(t, testLibrary)
	ctx := context.TODO()
	installCmd := command.New()
	_, err := installCmd.Run(ctx, "pkgr", "install", fmt.Sprintf("--config=%s", pkgrConfig))
	if err != nil {
		t.Fatalf("could not install baseline packages to '%s' with config file '%s'. Error: %s", testLibrary, pkgrConfig, err)
	}
}

func DeleteTestFolder(t *testing.T, folderToDelete string) {
	err := os.RemoveAll(folderToDelete)
	if err != nil {
		t.Fatalf("failed to clean up test library at '%s'. Error: %s", folderToDelete, err)
	}
}

func DeleteTestFile(t *testing.T, rpathSymlink string) {
	err := os.Remove(rpathSymlink)
	if err != nil {
		// Just fail the test if we can't setup properly.
		// It's a little more reliable than writing out own code to check for file existence.
		assert.NoFileExistsf(
			t,
			rpathSymlink,
			fmt.Sprintf("could not cleanup file (%s) from previous test run: %s", rpathSymlink, err),
		)
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
	Msg string `json:"msg,omitempty"`
	Pkg string `json:"pkg,omitempty"`
	Relationship string `json:"relationship,omitempty"`
	Repo string `json:"repo,omitempty"`
	Version string `json:"version,omitempty"`
}

func (prsm PkgRepoSetMsg) ToString() string {
	return fmt.Sprintf(
		"msg: '%s', pkg: '%s', version: '%s', repo: '%s', relationship: %s'",
		prsm.Msg,
		prsm.Pkg,
		prsm.Version,
		prsm.Repo,
		prsm.Relationship,
	)
}

// Returns -1 if A < B, 0 if A==B, and 1 if A > B
func ComparePkgRepoSetMsg(a, b PkgRepoSetMsg) int {
	if a.Pkg != b.Pkg {
		return strings.Compare(a.Pkg, b.Pkg)
	} else if a.Version != b.Version {
		return strings.Compare(a.Version, b.Version)
	} else if a.Repo != b.Repo {
		return strings.Compare(a.Repo, b.Repo)
	} else if a.Relationship != b.Relationship {
		return strings.Compare(a.Relationship, b.Relationship)
	} else {
		return 0
	}
}

//type PkgRepoSetMsgCollection struct {
//	Logs []PkgRepoSetMsg
//}

type PkgRepoSetMsgCollection []PkgRepoSetMsg

func (prsmc PkgRepoSetMsgCollection) Contains(pkg, version, repo, relationship string) bool {
	for _, log := range prsmc {
		if log.Pkg == pkg && log.Version == version && log.Repo == repo && log.Relationship == relationship {
			return true
		}
	}

	return false
}

func (prsmc PkgRepoSetMsgCollection) ToString() string {
	cleanStrings := []string{}
	for _, log := range prsmc {
		cleanStrings = append(cleanStrings, log.ToString())
	}
	return strings.Join(cleanStrings, "\n")
}

func (prsmc PkgRepoSetMsgCollection) ToBytes() []byte {
	return []byte(prsmc.ToString())
}

// PkgRepoSetMsgCollection returned will be sorted for the purposes of making golden files.
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

	sort.Slice(parsedLines, func(i, j int) bool {
		return ComparePkgRepoSetMsg(parsedLines[i], parsedLines[j]) < 0
	})

	return parsedLines
}

// This type will be more "flexible" than PkgRepoSet messages and should only be used for object comparisons for which
// it doesn't make sense to make a bunch of helper functions. Most of these fields will come back blank in most
// test runs.
// Basically, don't use this if you need golden files, use a more robust object with line-sorting.
// This object is meant to remain flexible with whatever fields we may need for a given test.
type GenericLog struct {
	Level string `json:"level,omitempty"`
	Msg string `json:"msg,omitempty"`
	Pkg string `json:"pkg,omitempty"`
	Package string `json:"package,omitempty"`
	Relationship string `json:"relationship,omitempty"`
	Repo string `json:"repo,omitempty"`
	InstallType int `json:"type,omitempty"`
	InstallFrom string `json:"install_from,omitempty"`
	Method string `json:"method,omitempty"`
	Origin string `json:"origin,omitempty"`
	Version string `json:"version,omitempty"`
	RSettings RSettingsDup `json:"RSettings,omitempty"` // "Do you want ciruclar dependencies? Because this is how you get circular dependencies."
	CmdArgs []string `json:"cmdArgs,omitempty"`
	RPath string `json:"rpath,omitempty"`
	LocalRepo int `json:"LOCALREPO,omitempty"` // Very specific to certain test cases
	Tarballs int `json:"tarballs"` // Very specific to certain test cases
	Library string `json:"library,omitempty"`
	ToInstall int `json:"to_install,omitempty"`
	ToUpdate int `json:"to_update,omitempty"`
	InstalledVersion string `json:"installed_version,omitempty"`
	UpdateVersion string `json:"update_version,omitempty"`
	Packages []string `json:"packages,omitempty"`
	Installed int `json:"installed,omitempty"`
	NotFromPkgr int `json:"not_from_pkgr,omitempty"`
	Outdated int `json:"outdated,omitempty"`
	TotalPackagesRequired int `json:"total_packages_required,omitempty"`
}



type GenericLogsCollection []GenericLog

// Only write these filter functions as needed
func (glc GenericLogsCollection) FilterByPackageTag(pkg string) GenericLogsCollection {
	matched := []GenericLog{}
	for _, obj := range glc {
		if(obj.Pkg == pkg) {
			matched = append(matched, obj)
		}
	}
	return matched
}

func CollectGenericLogs(t *testing.T, capture command.Capture, messageRegex string) GenericLogsCollection {
	re := regexp.MustCompile(messageRegex)

	parsedLines := []GenericLog{}

	outputLines := strings.Split(string(capture.Output), "\n")
	for _, line := range outputLines {
		if re.MatchString(line) {
			var parsedLine GenericLog
			err := json.Unmarshal([]byte(line), &parsedLine)
			if err != nil {
				t.Fatalf("error unmarshalling the following JSON line: '%s'. error was: %s", line, err)
			}
			parsedLines = append(parsedLines, parsedLine)
		}
	}

	// Where sorting would happen, except don't, because it's not worth maintaining a sort function for a generic object.

	return parsedLines
}
// ---------------------------------------------------------------------------------------------------------------------

// Need to duplicate these structs from Rcmd to avoid circular dependencies in Rcmd's unit tests.-----------------------
// Since we are pulling out the Rcmd package anyway, I think this is an okay shortcut to take for now.
type RSettingsDup struct {
	Version       cran.RVersion                `json:"r_version,omitempty"`
	LibPaths      []string                     `json:"lib_paths,omitempty"`
	Rpath         string                       `json:"rpath,omitempty"`
	GlobalEnvVars NvpListDup                      `json:"global_env_vars,omitempty"`
	PkgEnvVars    map[string]map[string]string `json:"pkg_env_vars,omitempty"`
	Platform      string                       `json:"platform,omitempty"`
}

// Nvp name-value pair, each of type string
type NvpDup struct {
	Name  string `json:"global_env_vars_name,omitempty"`
	Value string `json:"global_env_vars_value,omitempty"`
}

// NvpList is a slice of Nvp. The slice maintains consistent ordering of the Nvp objects
type NvpListDup struct {
	Pairs []NvpDup `json:"global_env_vars_pairs,omitempty"`
}
// ---------------------------------------------------------------------------------------------------------------------