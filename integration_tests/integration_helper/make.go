package integration_helper

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/metrumresearchgroup/command"

	"github.com/metrumresearchgroup/pkgr/integration_tests/integration_helper/targets"
)

// limitations for running the makefile:
// 1. it appears to need to be run in this directory, so command has been modified to take an optional WithDir input.

func isIntegrationTestDir() bool {
	wd, _ := os.Getwd()
	return path.Base(wd) == "integration_tests"
}

func Reset(tgt targets.Target) error {
	if !isIntegrationTestDir() {
		return errors.New("you must run make functions from the integration_tests directory")
	}
	r, err := command.New().Run(nil, "make", targets.ResetTargets[tgt])
	if err != nil {
		//goland:noinspection GoNilness
		fmt.Println(string(r.Output))
	}
	return err
}
