package mixed_source

import (
	"context"
	"github.com/metrumresearchgroup/command"
	. "github.com/metrumresearchgroup/pkgr/testhelper"
	"github.com/sebdah/goldie/v2"
	"testing"
)

const (
	mixedSourceE2ETest1 = "MXSRC-E2E-001"
)

const(
	goldenPlanWithCustomizations    = "plan-customizations"
	goldenInstallWithCustomizations = "install-customizations"
)

func TestMixedSource(t *testing.T) {
	t.Run(MakeTestName(mixedSourceE2ETest1, "pkgr can install both source and binary packages"), func(t *testing.T) {
		DeleteTestFolder(t, "test-library")
		DeleteTestFolder(t, "test-cache")
		ctx := context.TODO()
		planCmd := command.New()
		installCmd := command.New()
		testCmd := command.New(command.WithDir("Rscripts"))
		g := goldie.New(t)

		t.Run(MakeSubtestName(mixedSourceE2ETest1, "A", "repo and package customizations are set properly"), func(t *testing.T) {
			planCapture, err := planCmd.Run(ctx, "pkgr", "plan", "--loglevel=debug", "--logjson")
			if err != nil {
				t.Fatalf("error running pkgr plan: %s\noutput:\n%s", err, string(planCapture.Output))
			}
			pkgRepoSettings := CollectPkgRepoSetLogs(t, planCapture)

			// This should verify that:
			// Packages coming from MPNJuly2020 install from binaries, except digest, which should builid from source
			// Packages coming from MPNJune2021 install from source. Yaml should come from this repo
			// Yaml should also install its suggested packages (RUnit)
			g.Assert(t, goldenPlanWithCustomizations, pkgRepoSettings.ToBytes())
		})

		t.Run(MakeSubtestName(mixedSourceE2ETest1, "B", "pkgr can install from both source and binary files"), func(t *testing.T){
			installCapture, err := installCmd.Run(ctx, "pkgr", "install", "--logjson")
			if err != nil {
				t.Fatalf("error running pkgr install: %s\n output:\n%s", err, string(installCapture.Output))
			}
			testCapture, err := testCmd.Run(ctx, "Rscript", "--quiet", "install_test.R")
			if err != nil {
				t.Fatalf("error running R script to scan installed packages: %s\noutput:\n%s", err, string(testCapture.Output))
			}
			g.Assert(t, goldenInstallWithCustomizations, testCapture.Output)
		})
	})
}