package outdated_pkgs

import (
	"encoding/json"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/metrumresearchgroup/command"
	"github.com/sebdah/goldie/v2"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	"github.com/metrumresearchgroup/pkgr/testhelper"
)

const (
	outdatedPackagesE2ETest1 = "UPDT-E2E-001"
	outdatedPackagesE2ETest2 = "UPDT-E2E-002"
	outdatedPackagesE2ETest3 = "UPDT-E2E-003"
	outdatedPackagesE2ETest4 = "UPDT-E2E-004"
	outdatedPackagesE2ETest5 = "UPDT-E2E-005"
)

const (
	OutdatedPackagesUpdate     = "outdated-packages-update"
	OutdatedPackagesNoUpdate   = "outdated-packages-no-update"
	OutdatedPackagesUpdateFlag = "outdated-packages-update-flag"
)

func setupBaseline(t *testing.T) {
	err := os.RemoveAll("test-library")
	if err != nil {
		t.Fatalf("failed to cleanup test library: %s", err)
	}

	err = os.Mkdir("test-library", 0755)
	if err != nil {
		t.Fatalf("failed to create empty test library: %s", err)
	}

	fs := afero.NewOsFs()

	// I can't just install the packages with pkgr, because part of this test is
	// verifying that pkgr will still update tests that _weren't_ installed by
	// pkgr. For now, I'm re-using this helper function.
	err = testhelper.CopyDir(fs, "outdated-library", "test-library")
	if err != nil {
		t.Fatalf("error populating test-library with outdated packages: %s", err)
	}
}

// Utilities for this test suite only:
type UpdateMsg struct {
	InstalledVersion string `json:"installed_version,omitempty"`
	UpdateVersion    string `json:"update_version,omitempty"`
	Package          string `json:"pkg,omitempty"`
	Msg              string `json:"msg, omitempty"`
	Level            string `json:"level,omitempty"`
}

type UpdateLogs []UpdateMsg

func (ul UpdateLogs) Contains(pkg, installedVersion, updateVersion, loglevel string) bool {
	for _, msg := range ul {
		if msg.Package == pkg && msg.InstalledVersion == installedVersion && msg.UpdateVersion == updateVersion && msg.Level == loglevel {
			return true
		}
	}
	return false
}

type UpdateMsgType int

const (
	OutdatedMsg UpdateMsgType = iota
	ToUpdateMsg UpdateMsgType = iota
)

func CompareUpdateMsgs(a, b UpdateMsg) int {
	if a.Level != b.Level {
		return (strings.Compare(a.Level, b.Level))
	} else if a.Msg != b.Msg {
		return strings.Compare(a.Msg, b.Msg)
	} else if a.Package != b.Package {
		return strings.Compare(a.Package, b.Package)
	} else if a.InstalledVersion != b.InstalledVersion {
		return strings.Compare(a.InstalledVersion, b.InstalledVersion)
	} else if a.UpdateVersion != b.UpdateVersion {
		return strings.Compare(a.UpdateVersion, b.UpdateVersion)
	} else {
		return 0
	}
}

func CollectUpdateLogs(t *testing.T, capture []byte, msgType UpdateMsgType) UpdateLogs {
	parsedLines := []UpdateMsg{}

	var msgKey string
	if msgType == OutdatedMsg {
		msgKey = "outdated package found"
	} else if msgType == ToUpdateMsg {
		msgKey = "package will be updated"
	} else {
		t.Fatalf("invalid msgType passed to 'CollectUpdateLogs'. MsgType: %d", msgType)
	}

	outputLines := strings.Split(string(capture), "\n")

	for _, line := range outputLines {
		if strings.Contains(line, msgKey) {
			var parsedLine UpdateMsg
			err := json.Unmarshal([]byte(line), &parsedLine)
			if err != nil {
				t.Fatalf("error unmarshalling the following JSON line: '%s'. error was: %s", line, err)
			}
			parsedLines = append(parsedLines, parsedLine)
		}
	}

	sort.Slice(parsedLines, func(i, j int) bool {
		return CompareUpdateMsgs(parsedLines[i], parsedLines[j]) < 0
	})

	return parsedLines
}

// end utililities

func TestOutdated(t *testing.T) {

	t.Run(testhelper.MakeTestName(outdatedPackagesE2ETest1, "pkgr warns when updates are available and NoUpdate setting true"), func(t *testing.T) {
		setupBaseline(t)

		planCmd := command.New("pkgr", "plan", "--config=pkgr-no-update.yml", "--logjson")
		res, err := planCmd.CombinedOutput()
		if err != nil {
			t.Fatalf("failed to successfully create pkgr plan: %s", err)
		}

		outdatedLogs := CollectUpdateLogs(t, res, OutdatedMsg)
		assert.Len(t, outdatedLogs, 3, "expected outdated messages for three packages")

		assert.True(t, outdatedLogs.Contains("pillar", "1.2.1", "1.6.1", "warning"), "pkgr did not notify that package was outdated")
		assert.True(t, outdatedLogs.Contains("crayon", "1.2.1", "1.4.1", "warning"), "pkgr did not notify that package was outdated")
		assert.True(t, outdatedLogs.Contains("R6", "2.0", "2.5.0", "warning"), "pkgr did not notify that package was outdated")
	})

	t.Run(testhelper.MakeTestName(outdatedPackagesE2ETest2, "pkgr install does not update when NoUpdate setting true"), func(t *testing.T) {
		setupBaseline(t)
		installCmd := command.New("pkgr", "install", "--config=pkgr-no-update.yml", "--logjson")

		err := installCmd.Run()
		if err != nil {
			t.Fatalf("failed to install updated packages: %s", err)
		}
		testCmd := command.New("Rscript", "--quiet", "install_test.R")
		testCmd.Dir = "Rscripts"

		rScriptRes, err := testCmd.CombinedOutput()
		if err != nil {
			t.Fatalf("failed to run Rscript command with err: %s", err)
		}
		g := goldie.New(t)
		g.Assert(t, OutdatedPackagesNoUpdate, rScriptRes)
	})

	t.Run(testhelper.MakeTestName(outdatedPackagesE2ETest3, "pkgr plan indicates that updates will be installed when NoUpdate setting is false"), func(t *testing.T) {
		setupBaseline(t)

		planCmd := command.New("pkgr", "plan", "--config=pkgr-update.yml", "--logjson")
		planRes, err := planCmd.CombinedOutput()
		if err != nil {
			t.Fatalf("error running pkgr plan %s", err)
		}

		updateLogs := CollectUpdateLogs(t, planRes, ToUpdateMsg)
		assert.Len(t, updateLogs, 3, "expected updated messages for three packages")
		assert.True(t, updateLogs.Contains("pillar", "1.2.1", "1.6.1", "info"), "pkgr did not confirm package would be updated")
		assert.True(t, updateLogs.Contains("crayon", "1.2.1", "1.4.1", "info"), "pkgr did not confirm package would be updated")
		assert.True(t, updateLogs.Contains("R6", "2.0", "2.5.0", "info"), "pkgr did not confirm package would be updated")

	})

	t.Run(testhelper.MakeTestName(outdatedPackagesE2ETest4, "pkgr install installs updates when NoUpdate setting is false"), func(t *testing.T) {
		setupBaseline(t)

		installCmd := command.New("pkgr", "install", "--config=pkgr-update.yml")
		err := installCmd.Run()
		if err != nil {
			t.Fatalf("failed to install packages: %s", err)
		}

		testCmd := command.New("Rscript", "--quiet", "install_test.R")
		testCmd.Dir = "Rscripts"
		rScriptRes, err := testCmd.CombinedOutput()
		if err != nil {
			t.Fatalf("failed to run Rscript command with err: %s", err)
		}
		g := goldie.New(t)
		g.Assert(t, OutdatedPackagesUpdate, rScriptRes)
	})

	// this maintains legacy behavior of using --update in case people
	// have scripts or just by habit
	t.Run(testhelper.MakeTestName(outdatedPackagesE2ETest5, "pkgr install --update flag overrides other config"), func(t *testing.T) {
		setupBaseline(t)

		installCmd := command.New("pkgr", "install", "--config=pkgr-no-update.yml", "--update", "--logjson")
		installRes, err := installCmd.CombinedOutput()
		if err != nil {
			t.Fatalf("failed to install packages: %s", err)
		}

		testCmd := command.New("Rscript", "--quiet", "install_test.R")
		testCmd.Dir = "Rscripts"
		rScriptRes, err := testCmd.CombinedOutput()
		if err != nil {
			t.Fatalf("failed to run Rscript command with err: %s: ", err)
		}

		updateLogs := CollectUpdateLogs(t, installRes, ToUpdateMsg)
		assert.Len(t, updateLogs, 3, "expected exactly three update messages")
		assert.True(t, updateLogs.Contains("pillar", "1.2.1", "1.6.1", "info"), "pkgr did not confirm package would be updated")
		assert.True(t, updateLogs.Contains("crayon", "1.2.1", "1.4.1", "info"), "pkgr did not confirm package would be updated")
		assert.True(t, updateLogs.Contains("R6", "2.0", "2.5.0", "info"), "pkgr did not confirm package would be updated")

		g := goldie.New(t)
		g.Assert(t, OutdatedPackagesUpdateFlag, rScriptRes)
	})
}
