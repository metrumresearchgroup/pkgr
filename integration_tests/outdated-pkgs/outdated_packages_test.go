package outdated_pkgs

import (
	"context"
	"github.com/metrumresearchgroup/command"
	"github.com/metrumresearchgroup/pkgr/testhelper"
	"github.com/sebdah/goldie/v2"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const (
	OutdatedPackagesUpdate = "outdated-packages-update"
	OutdatedPackagesNoUpdate = "outdated-packages-no-update"
)

func setupBaseline(t *testing.T) {
	err := os.RemoveAll("test-library")
	if err != nil {
		t.Fatalf("failed to cleanup test library: %s", err)
	}

	err = os.Mkdir("test-library", 0755)
	if err != nil {
		t.Fatalf("failed to create empty test library: %s", err )
	}

	fs := afero.NewOsFs()

	// I can't just install the packages with pkgr, because part of this test is
	// verifying that pkgr will still update tests that _weren't_ installed by
	// pkgr. For now, I'm re-using this helper function.
	err = testhelper.CopyDir(fs, "outdated-library", "test-library")
	if err != nil {
		t.Fatalf("error populating test-library with outdated packages: %s", err)
	}

	//ctx := context.TODO()
	//installCmd := command.New()
	//_, err = installCmd.Run(ctx, "pkgr", "install", "--config=pkgr-baseline.yml")
	//if err != nil {
	//	t.Fatalf("could not install baseline packages with err: %s")
	//}
}

func TestOutdated(t *testing.T) {
	testCmd := command.New(command.WithDir("Rscripts"))
	planCmd := command.New()
	installCmd := command.New()
	ctx := context.TODO()


	t.Run("pkgr warns when updates are available", func(t *testing.T) {
		setupBaseline(t)

		res, err := planCmd.Run(ctx, "pkgr", "plan", "--config=pkgr-no-update.yml", "--logjson")
		if err != nil {
			t.Fatalf("failed to successfully create pkgr plan: %s", err)
		}

		jsonLine1 := `{"installed_version":"1.2.1","level":"warning","msg":"outdated package found","pkg":"crayon","update_version":"1.4.1"}`
		jsonLine2 := `{"installed_version":"2.0","level":"warning","msg":"outdated package found","pkg":"R6","update_version":"2.5.0"}`
		jsonLine3 := `{"installed_version":"1.2.1","level":"warning","msg":"outdated package found","pkg":"pillar","update_version":"1.6.1"}`
		jsonLine4 := `{"installed_version":"1.2-11","level":"warning","msg":"outdated package found","pkg":"Matrix","update_version":"1.3-4"}`

		assert.Contains(t, string(res.Output), jsonLine1, "missing expected warning message")
		assert.Contains(t, string(res.Output), jsonLine2, "missing expected warning message")
		assert.Contains(t, string(res.Output), jsonLine3, "missing expected warning message")
		assert.Contains(t, string(res.Output), jsonLine4, "missing expected warning message")
	})

	t.Run("pkgr does not update when update setting false", func(t *testing.T) {
		setupBaseline(t)

		_, err := installCmd.Run(ctx, "pkgr", "install", "--config=pkgr-no-update.yml", "--logjson")
		if err != nil {
			t.Fatalf("failed to install updated packages: %s", err)
		}
		rScriptResults, err := testCmd.Run(ctx, "Rscript", "--quiet", "get_installed.R")
		if err != nil {
			t.Fatalf("failed to run Rscript command with err: %s", err)
		}
		g := goldie.New(t)
		g.Assert(t, OutdatedPackagesNoUpdate, rScriptResults.Output)
	})


}