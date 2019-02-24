// +build integration

package rcmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRVersionExecution(t *testing.T) {
	assert := assert.New(t)
	rs := NewRSettings()
	expected := RVersion{3, 5, 2}
	assert.Equal(rs.Version, RVersion{}, "unitialized to no R version")
	actual := GetRVersion(&rs)
	assert.Equal(expected, actual, "returns the R version")
	assert.Equal(expected, rs.Version, "after initialization")
}
