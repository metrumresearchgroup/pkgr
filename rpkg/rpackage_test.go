package rpkg

import (
	"fmt"
	"testing"

	"github.com/spf13/afero"

	"github.com/stretchr/testify/assert"
)

func TestHashing(t *testing.T) {
	assert := assert.New(t)
	appFS := afero.NewOsFs()
	var data = []struct {
		in       string
		expected string
	}{
		{
			"../integration_tests/src/test1_0.0.1.tar.gz",
			"b28ba6e911e86ae4e682f834741e85e0",
		},
	}
	for i, tt := range data {
		actual, err := Hash(appFS, tt.in)
		if err != nil {
			assert.FailNowf("error hashing", "error: %s", err)
		}
		assert.Equal(tt.expected, actual, fmt.Sprintf("test num: %v", i+1))

	}
}
