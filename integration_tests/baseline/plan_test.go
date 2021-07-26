package baseline

import (
	"context"
	"github.com/metrumresearchgroup/command"
	. "github.com/metrumresearchgroup/pkgr/testhelper"
	"github.com/stretchr/testify/assert"
	"os"
	"regexp"
	"testing"
)

const (
	baselinePlanTest1="BSLNPLN-E2E-001"
	baselinePlanTest2="BSLNPLN-E2E-002"
	baselinePlanTest3="BSLNPLN-E2E-003"
)

func TestPlan(t *testing.T) {



	t.Run(MakeTestName(baselinePlanTest1, "plan indicates packages to be installed, as well as version, source repo, and whether pkg is user-defined or a dependency"), func(t *testing.T) {
		DeleteTestLibrary("test-library", t)
		ctx := context.TODO()
		planCmd := command.New()

		outputBytes, err := planCmd.Run(ctx, "pkgr", "plan", "--loglevel=debug", "--logjson")
		if err != nil {
			t.Fatalf("error running pkgr plan: %s", err)
		}
		output := string(outputBytes.Output)

		jsonRegexR6 := `\{"level":"debug","msg":"package repository set","pkg":"R6","relationship":"user_defined","repo":"CRAN","type":2,"version":"2.5.0"\}`
		jsonRegexPillar := `\{"level":"debug","msg":"package repository set","pkg":"pillar","relationship":"user_defined","repo":"CRAN","type":2,"version":"1.6.1"\}`
		jsonRegexGlue := `\{"level":"debug","msg":"package repository set","pkg":"glue","relationship":"dependency","repo":"CRAN","type":2,"version":"1.4.2"\}`
		jsonRegexFansi := `\{"level":"debug","msg":"package repository set","pkg":"fansi","relationship":"dependency","repo":"CRAN","type":2,"version":"0.5.0"\}`
		jsonRegexRlang := `\{"level":"debug","msg":"package repository set","pkg":"rlang","relationship":"dependency","repo":"CRAN","type":2,"version":"0.4.11"\}`
		jsonRegexUtf8 := `\{"level":"debug","msg":"package repository set","pkg":"utf8","relationship":"dependency","repo":"CRAN","type":2,"version":"1.2.1"\}`
		jsonRegexCrayon := `\{"level":"debug","msg":"package repository set","pkg":"crayon","relationship":"dependency","repo":"CRAN","type":2,"version":"1.4.1"\}`
		jsonRegexLifecycle := `\{"level":"debug","msg":"package repository set","pkg":"lifecycle","relationship":"dependency","repo":"CRAN","type":2,"version":"1.0.0"\}`
		jsonRegexVctrs := `\{"level":"debug","msg":"package repository set","pkg":"vctrs","relationship":"dependency","repo":"CRAN","type":2,"version":"0.3.8"\}`
		jsonRegexEllipsis := `\{"level":"debug","msg":"package repository set","pkg":"ellipsis","relationship":"dependency","repo":"CRAN","type":2,"version":"0.3.2"\}`
		jsonRegexCli := `\{"level":"debug","msg":"package repository set","pkg":"cli","relationship":"dependency","repo":"CRAN","type":2,"version":"2.5.0"\}`

		assert.Regexp(t, jsonRegexR6, output)
		assert.Regexp(t, jsonRegexPillar, output)
		assert.Regexp(t, jsonRegexGlue, output)
		assert.Regexp(t, jsonRegexFansi, output)
		assert.Regexp(t, jsonRegexRlang, output)
		assert.Regexp(t, jsonRegexUtf8, output)
		assert.Regexp(t, jsonRegexCrayon, output)
		assert.Regexp(t, jsonRegexLifecycle, output)
		assert.Regexp(t, jsonRegexVctrs, output)
		assert.Regexp(t, jsonRegexEllipsis, output)
		assert.Regexp(t, jsonRegexCli, output)
	})

	t.Run(MakeTestName(baselinePlanTest2, "number of workers (threads) can be set"), func(t *testing.T) {
		DeleteTestLibrary("test-library", t)
		ctx := context.TODO()
		planCmd := command.New()

		outputBytes, err := planCmd.Run(ctx, "pkgr", "plan", "--threads=5", "--logjson")
		if err != nil {
			t.Fatalf("error running pkgr plan: %s", err)
		}
		output := string(outputBytes.Output)

		jsonRegex := `\{"level":"info","msg":"Installation would launch 5 workers.*\}`
		assert.Regexp(t, jsonRegex, output)
	})

	// We can't determine the specific global cache location up front because it varies by OS and might be a tempdir.
	// So this test is only going to verify that it is printed out and used.
	t.Run(MakeTestName(baselinePlanTest3, "pkgr will indicate the location of the global cache and use it for installation"), func(t *testing.T) {
		DeleteTestLibrary("test-library", t)
		ctx := context.TODO()
		installCmd := command.New()

		outputBytes, err := installCmd.Run(ctx, "pkgr", "plan", "--config=pkgr-global-cache.yml", "--logjson")
		if err != nil {
			t.Fatalf("error running pkgr plan: %s", err)
		}
		output := string(outputBytes.Output)

		pkgDbRegex := `\{"level":"info","msg":"Database cache directory: .*\}` // We can't determine this up front because it varies by OS and might be a temp dir.
		pkgCacheRegex:= `\{"level":"info","msg":"Package installation cache directory: .*\}` // We can't determine this up front because it varies by OS and might be a temp dir.
		assert.Regexp(t, pkgDbRegex, output)
		assert.Regexp(t, pkgCacheRegex, output)

		// If that passed, we can assume this will work:
		rPkgDbs := regexp.MustCompile(`\{"level":"info","msg":"Database cache directory: (.*)"\}`)
		pkgDbDir := rPkgDbs.FindStringSubmatch(output)[1]
		rPkgCache := regexp.MustCompile(`\{"level":"info","msg":"Package installation cache directory: (.*)"\}`)
		pkgCacheDir := rPkgCache.FindStringSubmatch(output)[1]

		//t.Log(pkgDbDir)
		//t.Log(pkgCacheDir)
		pkgDbDirContents, err := os.ReadDir(pkgDbDir)
		if err != nil {
			t.Fatalf("error attempting to read global pkgDb dir: %s", err)
		}
		assert.NotEmpty(t, pkgDbDirContents)


		pkgCacheDirContents, err := os.ReadDir(pkgCacheDir)
		if err != nil {
			t.Fatalf("error attempting to read global pkgDb dir: %s", err)
		}

		assert.NotEmpty(t, pkgCacheDirContents)
	})


}
