package baseline

import (
	"testing"

	"github.com/metrumresearchgroup/command"
	"github.com/stretchr/testify/assert"

	"github.com/metrumresearchgroup/pkgr/cmd"
	. "github.com/metrumresearchgroup/pkgr/testhelper"
)

const (
	baselineRootE2ETest1 = "BSLNRT-E2E-001"
)

func TestRoot(t *testing.T) {
	t.Run(MakeTestName(baselineRootE2ETest1, "-v flag prints version"), func(t *testing.T) {
		rootCmd := command.New("pkgr", "-v")
		capture, err := rootCmd.CombinedOutput()
		if err != nil {
			t.Fatalf("error occurred running pkgr -v: %s", err)
		}
		assert.Equal(t, cmd.VERSION, string(capture))
	})
}
