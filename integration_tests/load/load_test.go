package load

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/metrumresearchgroup/command"
	"github.com/metrumresearchgroup/pkgr/cmd"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func setupEndToEndWithInstall(t *testing.T, pkgrConfig, testLibrary string) {
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

func makeTestName(testId, testName string) string {
	return(fmt.Sprintf("[%s] %s", testId, testName))
}

func makeSubtestName(testId, subtestId, testName string) string {
	return(fmt.Sprintf("[%s-%s] %s", testId, subtestId, testName))
}

// Test IDs
const (
	loadE2ETest1 = "LOAD-E2E-001"
	loadE2ETest2 = "LOAD-E2E-002"
	loadE2ETest3 = "LOAD-E2E-003"
	loadE2ETest4 = "LOAD-E2E-004"
	loadE2ETest5 = "LOAD-E2E-005"
	loadE2ETest6 = "LOAD-E2E-006"
	loadE2ETest7 = "LOAD-E2E-007"
)



func TestLoad(t *testing.T) {

	t.Run(makeTestName(loadE2ETest1, "load indicates that packages load successfully"), func(t *testing.T) {
		t.Run(makeSubtestName(loadE2ETest1, "A", "user packages only"), func(t *testing.T) {
			setupEndToEndWithInstall(t, "pkgr-load.yml", "test-library")

			loadCmd := command.New()
			ctx := context.TODO()

			outputBytes, err := loadCmd.Run(ctx, "pkgr", "load", "--config=pkgr-load.yml", "--loglevel=debug", "--logjson")
			if err != nil {
				t.Fatalf("error executing pkgr load command: %s", err)
			}

			output := string(outputBytes.Output)

			r6LoadRegex := `\{"level":"debug","msg":"Package loaded successfully","pkg":"R6","r_dir":".*"\}`
			ellipsisLoadRegex := `\{"level":"debug","msg":"Package loaded successfully","pkg":"ellipsis","r_dir":".*"\}`
			allPackagesLoadRegex := `\{"dependencies_attempted":false,"level":"info","msg":"all packages loaded successfully","user_packages_attempted":"true","working_directory":".*"\}`


			// One option of doing this
			//match, err := MatchString(jsonRegex1, output)
			//if err != nil {
			//	t.Fatalf("error checking regex against logs: %s", err)
			//}
			//assert.True(t, match, fmt.Sprintf("did not find regex in load output: %s", jsonRegex1))


			// Built in option
			assert.Regexp(t,  r6LoadRegex, output)
			assert.Regexp(t,  ellipsisLoadRegex, output)
			assert.Regexp(t,  allPackagesLoadRegex, output)

			// Lazy option
			//assert.Contains(t, output, "Package R6 loaded successfully")
			//assert.Contains(t, output, "Package pillar loaded successfully")
			//assert.Contains(t, output, "all packages loaded successfully")
		})

		t.Run(makeSubtestName(loadE2ETest1, "B", "user packages and deps"), func(t *testing.T) {
			setupEndToEndWithInstall(t, "pkgr-load.yml", "test-library")

			loadCmd := command.New()
			ctx := context.TODO()

			outputBytes, err := loadCmd.Run(ctx, "pkgr", "load", "--config=pkgr-load.yml", "--all", "--loglevel=debug", "--logjson")
			if err != nil {
				t.Fatalf("error executing pkgr load command: %s", err)
			}

			output := string(outputBytes.Output)

			r6LoadRegex := `\{"level":"debug","msg":"Package loaded successfully","pkg":"R6".*\}`
			ellipsisLoadRegex := `\{"level":"debug","msg":"Package loaded successfully","pkg":"ellipsis".*\}`
			rlangLoadRegex := `\{"level":"debug","msg":"Package loaded successfully","pkg":"rlang".*\}`
			allPackagesLoadRegex := `\{"dependencies_attempted":true,"level":"info","msg":"all packages loaded successfully","user_packages_attempted":"true","working_directory":".*"\}`

			assert.Regexp(t,  r6LoadRegex, output)
			assert.Regexp(t,  ellipsisLoadRegex, output)
			assert.Regexp(t,  rlangLoadRegex, output)
			assert.Regexp(t,  allPackagesLoadRegex, output)
		})


	})

	t.Run(makeTestName(loadE2ETest2, "Load outputs a JSON report"), func(t *testing.T) {
		setupEndToEndWithInstall(t, "pkgr-load.yml", "test-library")

		loadCmd := command.New()
		ctx := context.TODO()

		outputBytes, err := loadCmd.Run(ctx, "pkgr", "load", "--config=pkgr-load.yml", "--all", "--json", "--loglevel=debug", "--logjson")
		if err != nil {
			t.Fatalf("error executing pkgr load command: %s", err)
		}

		//output := string(outputBytes.Output)
		var loadReport *cmd.LoadReport

		err = json.Unmarshal(outputBytes.Output, &loadReport)
		if err != nil {
			t.Fatalf("could not unmarshal JSON into expected format: %s", err)
		}

		//loadResultsMap := loadReport.LoadResults
		//r6, found := loadResultsMap["R6"]
		//assert.True(t, found, "failed to find load results for R6 in report")
		//assert.True(t, r6.Success, "R6 did not load successfully in report")
		checkLoadReport(t, loadReport, "R6", true)
		checkLoadReport(t, loadReport, "ellipsis", true)
		checkLoadReport(t, loadReport, "rlang", true)
	})


}

func checkLoadReport(t *testing.T, report *cmd.LoadReport, pkg string, expectSuccess bool) {
	pkgLoadResults, found := report.LoadResults[pkg]
	assert.True(t, found, fmt.Sprintf("failed to find load results for %s in report", pkg))
	if expectSuccess {
		assert.True(t, pkgLoadResults.Success, fmt.Sprintf("Report does not indicate that pkg '%s' loaded successfully", pkg))
	} else {
		assert.False(t, pkgLoadResults.Success, fmt.Sprintf("Report indicates that pkg '%s' loaded successfully, but we expected it to fail", pkg))
	}
}