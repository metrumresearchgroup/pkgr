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
			pillarLoadRegex := `\{"level":"debug","msg":"Package loaded successfully","pkg":"pillar","r_dir":".*"\}`
			allPackagesLoadRegex := `\{"dependencies_attempted":false,"level":"info","msg":"all packages loaded successfully","user_packages_attempted":"true","working_directory":".*"\}`


			// One option of doing this
			//match, err := MatchString(jsonRegex1, output)
			//if err != nil {
			//	t.Fatalf("error checking regex against logs: %s", err)
			//}
			//assert.True(t, match, fmt.Sprintf("did not find regex in load output: %s", jsonRegex1))


			// Built in option
			assert.Regexp(t,  r6LoadRegex, output)
			assert.Regexp(t,  pillarLoadRegex, output)
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
			//r6LoadRegex := `\{"level":"debug","msg":"Package loaded successfully","pkg":"R6","r_dir":".*"\}`
			pillarLoadRegex := `\{"level":"debug","msg":"Package loaded successfully","pkg":"pillar".*\}`
			utf8LoadRegex := `\{"level":"debug","msg":"Package loaded successfully","pkg":"utf8".*\}`
			fansiLoadRegex := `\{"level":"debug","msg":"Package loaded successfully","pkg":"fansi".*\}`
			cliLoadRegex := `\{"level":"debug","msg":"Package loaded successfully","pkg":"cli".*\}`
			assertthatLoadRegex := `\{"level":"debug","msg":"Package loaded successfully","pkg":"assertthat".*\}`
			crayonLoadRegex := `\{"level":"debug","msg":"Package loaded successfully","pkg":"crayon".*\}`
			rlangLoadRegex := `\{"level":"debug","msg":"Package loaded successfully","pkg":"rlang".*\}`
			allPackagesLoadRegex := `\{"dependencies_attempted":true,"level":"info","msg":"all packages loaded successfully","user_packages_attempted":"true","working_directory":".*"\}`

			assert.Regexp(t,  r6LoadRegex, output)
			assert.Regexp(t,  pillarLoadRegex, output)
			assert.Regexp(t,  utf8LoadRegex, output)
			assert.Regexp(t,  fansiLoadRegex, output)
			assert.Regexp(t,  cliLoadRegex, output)
			assert.Regexp(t,  assertthatLoadRegex, output)
			assert.Regexp(t,  crayonLoadRegex, output)
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

		json.Unmarshal(outputBytes.Output, &loadReport)

		// Just trying this out for now.
		assert.Len(t, loadReport.LoadResults, 12)

	})


}