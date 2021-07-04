// +build linux darwin

package repo_local

import (
	"context"
	"os"
	"testing"

	"github.com/metrumresearchgroup/command"
)

func TestPlanInstall(t *testing.T) {
	env := os.Environ()

	// this is where we manipulate the environment

	r, err := command.New(command.WithEnv(env)).Run(context.Background(), "pkgr", "--logjson", "plan")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	t.Logf("first run output: %s", r.Output)

	r, err = command.New(command.WithEnv(env)).Run(context.Background(), "pkgr", "--logjson", "install")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	t.Logf("second run output: %s", r.Output)
}
