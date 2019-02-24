package rcmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRVersion(t *testing.T) {
	assert := assert.New(t)

	var installArgsTests = []struct {
		in             RVersion
		expectedString string
		expectedFull   string
	}{
		{
			RVersion{3, 5, 2},
			"3.5",
			"3.5.2",
		},
		{
			RVersion{2, 1, 4},
			"2.1",
			"2.1.4",
		},
	}
	for i, tt := range installArgsTests {
		actual := tt.in.ToString()
		assert.Equal(tt.expectedString, actual, fmt.Sprintf("test num: %v", i+1))
		actual = tt.in.ToFullString()
		assert.Equal(tt.expectedFull, actual, fmt.Sprintf("test num: %v", i+1))

	}
}
