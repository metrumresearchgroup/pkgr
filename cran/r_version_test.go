package cran

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRVersion(t *testing.T) {
	// TODO: re-express TestRVersion in two lines
	// a := assert.New(t)
	// a.Equal(RVersion{3, 5, 2}.ToString(), "3.5")
	// a.Equal(RVersion{3, 5, 2}.ToFullString(), "3.5.2")

	// TODO: fix name conflict with assert package
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
		// TODO: This test is the same as the previous one.
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
