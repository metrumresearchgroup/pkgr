package baseline

import (
	"context"
	"github.com/metrumresearchgroup/command"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInstall(t *testing.T) {
	installCmd := command.New()
	ctx := context.TODO()
	// should have some setup to make sure the test-library is cleared out
	installRes, err := installCmd.Run(ctx, "pkgr", "install")
	if err != nil {
		panic(err)
	}

	assert.Contains(t, installRes.Output, "to_install=11 to_update=0")
	assert.Contains(t, installRes.Output, "Library path to install packages: test-library")

	testCmd := command.New(command.WithDir("Rscripts"))
	// we could also suppress any global state by doing things like setting --vanilla
	// then setting the R_LIBS_SITE using command.WithEnv() however
	// this way it will be easy to interactively develop the script in the Rproj,
	// then invoke it whenever needed
	testRes, err := testCmd.Run(ctx, "Rscript", "--quiet", "install_test.R")
	if err != nil {
		panic(err)
	}
	g := goldie.New(t)
	g.Assert(t, "install", []byte(testRes.Output))
}
