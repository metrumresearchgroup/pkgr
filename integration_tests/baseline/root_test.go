package baseline

import (
	"context"
	"github.com/metrumresearchgroup/command"
	"github.com/metrumresearchgroup/pkgr/cmd"
	. "github.com/metrumresearchgroup/pkgr/testhelper"
	"github.com/stretchr/testify/assert"
	"testing"
)

const(
	baselineRootE2ETest1 = "BSLNRT-E2E-001"
)

func TestRoot(t *testing.T) {
	t.Run(MakeTestName(baselineRootE2ETest1, "-v flag prints version"), func(t *testing.T) {
		ctx := context.TODO()
		rootCmd := command.New()
		capture, err := rootCmd.Run(ctx, "pkgr", "-v")
		if err != nil {
			t.Fatalf("error occurred running pkgr -v: %s", err)
		}
		assert.Equal(t, cmd.VERSION, string(capture.Output))
	})
}