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
		// TODO: update filesystem with files and uncomment, or delete this test
		// {
		// 	"../integration_tests/src/test1_0.0.1.tar.gz",
		// 	"b28ba6e911e86ae4e682f834741e85e0",
		// },
		// {
		// 	"../integration_tests/testsets/testset2/packrat/src/test1/test1_0.0.1.tar.gz",
		// 	"43ff4f49dfe4c9628c9b48160264dfb2",
		// },
	}
	for i, tt := range data {
		actual, err := Hash(appFS, tt.in)
		if err != nil {
			assert.FailNowf("error hashing", "error: %s", err)
		}
		assert.Equal(tt.expected, actual, fmt.Sprintf("test num: %v", i+1))

	}
}
