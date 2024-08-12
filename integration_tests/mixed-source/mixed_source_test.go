package mixed_source

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"text/template"

	"github.com/metrumresearchgroup/command"
	"github.com/sebdah/goldie/v2"
	//"github.com/stretchr/testify/assert"

	. "github.com/metrumresearchgroup/pkgr/testhelper"
)

const (
	mixedSourceE2ETest1 = "MXSRC-E2E-001"
	mixedSourceE2ETest2 = "MXSRC-E2E-002"
)

const (
	goldenPlan              = "plan-customizations-linux"
	goldenInstall           = "install-customizations-linux"
)

func skipIfNotUbunutu(t *testing.T) {
	t.Helper()

	prog := "lsb_release"
	if _, err := exec.LookPath(prog); err != nil {
		t.Skip("lsb_release not available; test requires Ubuntu")
	}
	cmd := command.New(prog, "-i")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("calling lsb_release failed: %s", err)
	}
	if !bytes.HasSuffix(out, []byte("Ubuntu\n")) {
		t.Skip("test requires Ubuntu")
	}
}

func generateConfig() (string, error) {
	outfile := "pkgr-linux.yml"
	tfile := outfile + ".template"

	t, err := template.New(tfile).ParseFiles(tfile)
	if err != nil {
		return "", err
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	cmd := command.New(filepath.Join(wd, "setup-binary-repo.sh"))
	cmdout, err := cmd.Output()
	if err != nil {
		return "", err
	}
	repo := strings.TrimSuffix(string(cmdout), "\n")

	values := map[string]string{
		"BinaryRepo": repo,
	}

	f, err := os.Create(outfile)
	if err != nil {
		return "", err
	}

	if err = t.Execute(f, values); err != nil {
		f.Close()
		return "", err
	}

	return outfile, f.Close()
}

func TestMixedSource(t *testing.T) {
	skipIfNotUbunutu(t)
	t.Run(MakeTestName(mixedSourceE2ETest1, "pkgr can install both source and binary packages"), func(t *testing.T) {
		DeleteTestFolder(t, "test-library")
		DeleteTestFolder(t, "test-cache")
		g := goldie.New(t)

		pkgrConfig, err := generateConfig()
		if err != nil {
			t.Fatalf("generating pkgr.yml failed: %s", err)
		}

		t.Run(MakeSubtestName(mixedSourceE2ETest1, "A", "repo and package customizations are set properly"), func(t *testing.T) {
			planCmd := command.New("pkgr", "plan", fmt.Sprintf("--config=%s", pkgrConfig), "--loglevel=debug", "--logjson")
			planCapture, err := planCmd.CombinedOutput()
			if err != nil {
				t.Fatalf("error running pkgr plan: %s\noutput:\n%s", err, string(planCapture))
			}
			pkgRepoSettings := CollectPkgRepoSetLogs(t, planCapture)

			// This should verify that:
			// Packages coming from MPNJuly2020 install from binaries, except digest, which should builid from source
			// Packages coming from MPNJune2021 install from source. Yaml should come from this repo
			// Yaml should also install its suggested packages (RUnit)
			g.Assert(t, goldenPlan, pkgRepoSettings.ToBytesWithType())
		})

		t.Run(MakeSubtestName(mixedSourceE2ETest1, "B", "pkgr can install from both source and binary files"), func(t *testing.T) {
			installCmd := command.New("pkgr", "install", fmt.Sprintf("--config=%s", pkgrConfig), "--logjson")
			installCapture, err := installCmd.CombinedOutput()
			if err != nil {
				t.Fatalf("error running pkgr install: %s\n output:\n%s", err, string(installCapture))
			}

			testCmd := command.New("Rscript", "--quiet", "install_test.R")
			testCmd.Dir = "Rscripts"
			testCapture, err := testCmd.CombinedOutput()
			if err != nil {
				t.Fatalf("error running R script to scan installed packages: %s\noutput:\n%s", err, string(testCapture))
			}
			g.Assert(t, goldenInstall, testCapture)
		})
	})
	t.Skip(MakeTestName(mixedSourceE2ETest2, "repo and package customizations synchronize when compatible. SKIPPING till issue #329 fixes this bug."))
	/*t.Run(MakeTestName(mixedSourceE2ETest2, "repo and package customizations synchronize when compatible"), func(t *testing.T) {
		DeleteTestFolder(t, "test-library")
		DeleteTestFolder(t, "test-cache")
		planCmd := command.New("pkgr", "plan", "--config=pkgr-issue-329.yml", "--loglevel=debug", "--logjson")

		planCapture, err := planCmd.CombinedOutput()
		if err != nil {
			t.Fatalf("error running pkgr plan: %s\noutput:\n%s", err, string(planCapture))
		}
		pkgRepoSettings := CollectPkgRepoSetLogs(t, planCapture)
		assert.True(t, pkgRepoSettings.ContainsWithType("R6", "2.5.0", "MPNSource", "user_defined", "source"), "expected 'R6' version 2.5.0 installed from source.\nActual pkg plan:\n%s", pkgRepoSettings.ToStringWithType())
		assert.True(t, pkgRepoSettings.ContainsWithType("digest", "0.6.25", "MPNBinary", "user_defined", "binary"), "expected 'digest' version 0.6.25 installed from binary\nActual pkg plan:\n%s", pkgRepoSettings.ToStringWithType())
	})*/
}
