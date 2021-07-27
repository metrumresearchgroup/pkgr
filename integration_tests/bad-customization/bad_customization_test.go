package bad_customization

import (
	"context"
	"github.com/metrumresearchgroup/command"
	. "github.com/metrumresearchgroup/pkgr/testhelper"
	"github.com/stretchr/testify/assert"
	"testing"
)

// test IDs
const (
	badCustomizationE2ETest1 = "BDCST-E2E-001"
)


func TestBadCustomization(t *testing.T) {
	t.Run(MakeTestName(badCustomizationE2ETest1, "repo from customization not specified under repos"), func(t *testing.T) {
		DeleteTestFolder(t, "test-library")

		ctx := context.TODO()
		installCmd := command.New()

		capture, err := installCmd.Run(ctx, "pkgr", "install", "--logjson")
		assert.Error(t, err, "expected an error to be thrown, but got none")

		logs := CollectGenericLogs(t, capture, "error finding custom repo to set")
		assert.Len(t, logs, 1, )
		errorMessage := logs[0]
		assert.Equal(t, "R6", errorMessage.Pkg)
		assert.Equal(t, "DoesNotExist", errorMessage.Repo)
	})
}