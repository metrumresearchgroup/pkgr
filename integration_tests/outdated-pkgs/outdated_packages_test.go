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
	OutdatedPackagesUpdateFlag = "outdated-packages-update-flag"
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
}

func TestOutdated(t *testing.T) {
	testCmd := command.New(command.WithDir("Rscripts"))
	planCmd := command.New()
	installCmd := command.New()
	ctx := context.TODO()


	t.Run("pkgr warns when updates are available and update setting false", func(t *testing.T) {
		setupBaseline(t)

		res, err := planCmd.Run(ctx, "pkgr", "plan", "--config=pkgr-no-update.yml", "--logjson")
		if err != nil {
			t.Fatalf("failed to successfully create pkgr plan: %s", err)
		}

		// Hmmm... I wonder if there's a way to go like "g.AssertContains()"
		jsonLine1 := `{"installed_version":"1.2.1","level":"warning","msg":"outdated package found","pkg":"crayon","update_version":"1.4.1"}`
		jsonLine2 := `{"installed_version":"2.0","level":"warning","msg":"outdated package found","pkg":"R6","update_version":"2.5.0"}`
		jsonLine3 := `{"installed_version":"1.2.1","level":"warning","msg":"outdated package found","pkg":"pillar","update_version":"1.6.1"}`
		jsonLine4 := `{"installed_version":"1.2-11","level":"warning","msg":"outdated package found","pkg":"Matrix","update_version":"1.3-4"}`

		assert.Contains(t, string(res.Output), jsonLine1, "missing expected warning message")
		assert.Contains(t, string(res.Output), jsonLine2, "missing expected warning message")
		assert.Contains(t, string(res.Output), jsonLine3, "missing expected warning message")
		assert.Contains(t, string(res.Output), jsonLine4, "missing expected warning message")
	})

	t.Run("pkgr install does not update when update setting false", func(t *testing.T) {
		setupBaseline(t)

		_, err := installCmd.Run(ctx, "pkgr", "install", "--config=pkgr-no-update.yml", "--logjson")
		if err != nil {
			t.Fatalf("failed to install updated packages: %s", err)
		}
		rScriptRes, err := testCmd.Run(ctx, "Rscript", "--quiet", "install_test.R")
		if err != nil {
			t.Fatalf("failed to run Rscript command with err: %s", err)
		}
		g := goldie.New(t)
		g.Assert(t, OutdatedPackagesNoUpdate, rScriptRes.Output)
	})

	t.Run("pkgr plan indicates that updates will be installed when update setting is true", func(t *testing.T) {
		setupBaseline(t)

		planRes, err := planCmd.Run(ctx, "pkgr", "plan", "--config=pkgr-update.yml", "--logjson")
		if err != nil {
			t.Fatalf("error running pkgr plan %s", err)
		}

		jsonLine1 := `{"installed_version":"1.2.1","level":"info","msg":"package will be updated","pkg":"pillar","update_version":"1.6.1"}`
		jsonLine2 := `{"installed_version":"1.2.1","level":"info","msg":"package will be updated","pkg":"crayon","update_version":"1.4.1"}`
		jsonLine3 := `{"installed_version":"2.0","level":"info","msg":"package will be updated","pkg":"R6","update_version":"2.5.0"}`
		jsonLine4 := `{"installed_version":"1.2-11","level":"info","msg":"package will be updated","pkg":"Matrix","update_version":"1.3-4"}`

		assert.Contains(t, string(planRes.Output), jsonLine1, "pkgr did not confirm package would be updated")
		assert.Contains(t, string(planRes.Output), jsonLine2, "pkgr did not confirm package would be updated")
		assert.Contains(t, string(planRes.Output), jsonLine3, "pkgr did not confirm package would be updated")
		assert.Contains(t, string(planRes.Output), jsonLine4, "pkgr did not confirm package would be updated")
	})

	t.Run("pkgr install installs updates when update setting is true", func(t *testing.T) {
		setupBaseline(t)

		_, err := installCmd.Run(ctx, "pkgr", "install", "--config=pkgr-update.yml")
		if err != nil {
			t.Fatalf("failed to install packages: %s", err)
		}
		rScriptRes, err := testCmd.Run(ctx, "Rscript", "--quiet", "install_test.R")
		if err != nil {
			t.Fatalf("failed to run Rscript command with err: %s", err)
		}
		g := goldie.New(t)
		g.Assert(t, OutdatedPackagesUpdate, rScriptRes.Output)
	})

	t.Run("pkgr install --update performs updates", func(t *testing.T) {
		setupBaseline(t)
		installRes, err := installCmd.Run(ctx, "pkgr", "install", "--config=pkgr-update-flag-unset.yml", "--update", "--logjson")
		if err != nil {
			t.Fatalf("failed to install packages: %s", err)
		}
		rScriptRes, err := testCmd.Run(ctx, "Rscript", "--quiet", "install_test.R")
		if err != nil {
			t.Fatalf("failed to run Rscript command with err: %s: ", err)
		}

		jsonLine1 := `{"installed_version":"1.2.1","level":"info","msg":"package will be updated","pkg":"pillar","update_version":"1.6.1"}`
		jsonLine2 := `{"installed_version":"1.2.1","level":"info","msg":"package will be updated","pkg":"crayon","update_version":"1.4.1"}`
		jsonLine3 := `{"installed_version":"2.0","level":"info","msg":"package will be updated","pkg":"R6","update_version":"2.5.0"}`
		jsonLine4 := `{"installed_version":"1.2-11","level":"info","msg":"package will be updated","pkg":"Matrix","update_version":"1.3-4"}`

		assert.Contains(t, string(installRes.Output), jsonLine1, "pkgr did not confirm package would be updated")
		assert.Contains(t, string(installRes.Output), jsonLine2, "pkgr did not confirm package would be updated")
		assert.Contains(t, string(installRes.Output), jsonLine3, "pkgr did not confirm package would be updated")
		assert.Contains(t, string(installRes.Output), jsonLine4, "pkgr did not confirm package would be updated")

		g := goldie.New(t)
		g.Assert(t, OutdatedPackagesUpdateFlag, rScriptRes.Output)
	})
}