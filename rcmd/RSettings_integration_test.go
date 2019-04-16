// +build R

package rcmd

import (
	"testing"

	"github.com/metrumresearchgroup/pkgr/cran"
	"github.com/stretchr/testify/assert"
)

func TestRVersionExecution(t *testing.T) {
	assert := assert.New(t)
	rs := NewRSettings("")
	// this test expects a machine with R 3.5.2 available on the default System Path
	// TODO: refactor to make more generalized or mock
	expected := cran.RVersion{3, 5, 2}
	assert.Equal(rs.Version, cran.RVersion{}, "unitialized to no R version")
	actual := GetRVersion(&rs)
	assert.Equal(expected, actual, "returns the R version")
	assert.Equal(expected, rs.Version, "after initialization")
}
