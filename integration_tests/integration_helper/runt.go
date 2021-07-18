package integration_helper

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/metrumresearchgroup/command"

	"github.com/metrumresearchgroup/pkgr/integration_tests/integration_helper/targets"
)

func RunTest(tgt targets.Target) error {
	if !isIntegrationTestDir() {
		return errors.New("you must run tests from the integration_tests directory")
	}

	dir, ok := targets.RunTargets[tgt]
	if !ok {
		return fmt.Errorf("target not found")
	}

	err := Reset(tgt)
	if err != nil {
		return err
	}

	r, err := command.New(command.WithDir(dir)).Run(context.Background(), "go", "test", "-json", ".")

	//goland:noinspection GoNilness
	fmt.Println(r.Output)

	if err != nil {
		return err
	}

	if strings.Contains(string(r.Output), "FAIL:") {
		return errors.New("found FAIL message")
	}

	return nil
}
