package cran

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDownloadFMethod_GetMegabytes(t *testing.T) {
	testFixture := Download{
		Size: 1048576, //1024 * 1024 -- 1 MB
	}

	assert.Equal(t, 1.0, testFixture.GetMegabytes())
}
