package load

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/metrumresearchgroup/command"
	"github.com/metrumresearchgroup/pkgr/cmd"
	"github.com/metrumresearchgroup/pkgr/testhelper"
	"github.com/stretchr/testify/assert"
	"os"
	. "path"
	"testing"
)


func setupCorruptedPackageEnvironment(t *testing.T) {
	testLibPath := "test-library"
	testhelper.SetupEndToEndWithInstall(t, "pkgr-load-fail-setup.yml", testLibPath)

	// Remove "R" folder from R6
	err := os.RemoveAll(Join(testLibPath, "R6", "R"))
	if err != nil {
		t.Fatalf("error while corrupting R6 package during test setup: %s", err)
	}

	// Remove "R" folder from rlang, which ellipsis depends on.
	err = os.RemoveAll(Join(testLibPath, "rlang", "R"))
	if err != nil {
		t.Fatalf("error while corrupting rlang package during test setup: %s", err)
	}

	// Desired end state:
	// R6 is installed but can't be loaded
	// rlang is installed but can't be loaded
	// ellipsis is installed properly, but can't be loaded only because rlang can't be loaded
	// fansi is not installed and therefore can't be loaded.
	// glue is installed and can be loaded.

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

	t.Run(testhelper.MakeTestName(loadE2ETest1, "load indicates that packages load successfully"), func(t *testing.T) {
		t.Run(testhelper.MakeSubtestName(loadE2ETest1, "A", "user packages only"), func(t *testing.T) {
			testhelper.SetupEndToEndWithInstall(t, "pkgr-load.yml", "test-library")

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

		t.Run(testhelper.MakeSubtestName(loadE2ETest1, "B", "user packages and deps"), func(t *testing.T) {
			testhelper.SetupEndToEndWithInstall(t, "pkgr-load.yml", "test-library")

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

	t.Run(testhelper.MakeTestName(loadE2ETest2, "Load outputs a JSON report"), func(t *testing.T) {
		testhelper.SetupEndToEndWithInstall(t, "pkgr-load.yml", "test-library")

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
		assert.Equal(t, 0, loadReport.Failures)
	})

	t.Run(testhelper.MakeTestName(loadE2ETest3, "Load fails when a package can't load"), func(t *testing.T) {
		setupCorruptedPackageEnvironment(t)

		loadCmd := command.New()
		ctx := context.TODO()

		outputBytes, err := loadCmd.Run(ctx, "pkgr", "load", "--config=pkgr-load-fail-test.yml", "--all", "--loglevel=debug", "--logjson")
		if err != nil {
			t.Fatalf("error when executing pkgr load: %s", err)
		}
		output := string(outputBytes.Output)

		// Note: Failure messages include backticks around library calls, which we are replacing with Regex periods for convenience.
		r6FailRegex := `\{"go_error":"exit status 1","level":"error","msg":"error loading package via .library\(R6\)..*`
		rlangFailRegex := `\{"go_error":"exit status 1","level":"error","msg":"error loading package via .library\(rlang\)..*`
		ellipsisFailRegex := `\{"go_error":"exit status 1","level":"error","msg":"error loading package via .library\(ellipsis\)..*`
		fansiFailRegex := `\{"go_error":"exit status 1","level":"error","msg":"error loading package via .library\(fansi\)..*`

		glueSucceedRegex := `\{"level":"debug","msg":"Package loaded successfully","pkg":"glue".*\}`

		overallFailRegex := `\{"failures":4,"level":"error","msg":"some packages failed to load.","working_directory":".*"\}`

		assert.Regexp(t, r6FailRegex, output)
		assert.Regexp(t, rlangFailRegex, output)
		assert.Regexp(t, ellipsisFailRegex, output)
		assert.Regexp(t, fansiFailRegex, output)
		assert.Regexp(t, glueSucceedRegex, output)
		assert.Regexp(t, overallFailRegex, output)

	})

	t.Run(testhelper.MakeTestName(loadE2ETest4, "Load JSON report captures failures"), func(t *testing.T) {
		setupCorruptedPackageEnvironment(t)

		loadCmd := command.New()
		ctx := context.TODO()

		outputBytes, err := loadCmd.Run(ctx, "pkgr", "load", "--config=pkgr-load-fail-test.yml", "--all", "--json")
		if err != nil {
			t.Fatalf("error when executing pkgr load: %s", err)
		}

		//output := string(outputBytes.Output)
		var loadReport *cmd.LoadReport

		err = json.Unmarshal(outputBytes.Output, &loadReport)
		if err != nil {
			t.Fatalf("could not unmarshal JSON into expected format: %s", err)
		}

		checkLoadReport(t, loadReport, "R6", false)
		checkLoadReport(t, loadReport, "rlang", false)
		checkLoadReport(t, loadReport, "ellipsis", false)
		checkLoadReport(t, loadReport, "fansi", false)
		checkLoadReport(t, loadReport, "glue", true)
		assert.Equal(t, 4, loadReport.Failures, "incorrect number of failed loads")
	})


}

func TestLoadFail(t *testing.T) {

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