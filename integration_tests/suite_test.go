// +build linux darwin

package integration_tests__test

import (
	"testing"

	"github.com/metrumresearchgroup/pkgr/integration_tests/integration_helper"
	"github.com/metrumresearchgroup/pkgr/integration_tests/integration_helper/targets"
)

func Test(t *testing.T) {
	err := integration_helper.RunTest(targets.RepoLocal)
	if err != nil {
		t.Fatalf("%v", err)
	}
}
