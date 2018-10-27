package desc

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	assert := assert.New(t)
	var data = []struct {
		in       string
		expected Version
	}{
		{
			"1.0.0",
			Version{1, 0, 0, 0, 0, "1.0.0"},
		},
		{
			"1.0-0",
			Version{1, 0, 0, 0, 0, "1.0-0"},
		},
		{
			"0.1.2.9000",
			Version{0, 1, 2, 9000, 0, "0.1.2.9000"},
		},
		{
			"1.2.3.4",
			Version{1, 2, 3, 4, 0, "1.2.3.4"},
		},
		{
			"1-2-3-4",
			Version{1, 2, 3, 4, 0, "1-2-3-4"},
		},
	}
	for i, tt := range data {
		actual := ParseVersion(tt.in)
		assert.Equal(tt.expected, actual, fmt.Sprintf("test num: %v", i+1))
	}
}

func TestComparisonString(t *testing.T) {
	assert := assert.New(t)
	var data = []struct {
		in       []string
		expected int
	}{
		{
			[]string{"0.1.0", "0.1.1"},
			-1,
		},
		{
			[]string{"0.1.2", "0.1.1"},
			1,
		},
		{
			[]string{"0.1.1", "0.1.1"},
			0,
		},
		{
			[]string{"1.0.0", "0.1.1"},
			1,
		},
		{
			[]string{"0.0.9", "0.1.0"},
			-1,
		},
		{
			[]string{"0.0.0.1", "0.0.0.0"},
			1,
		},
		{
			[]string{"0.0.0.0", "0.0.0.0"},
			0,
		},
	}
	for i, tt := range data {
		actual := CompareVersionStrings(tt.in[0], tt.in[1])
		assert.Equal(tt.expected, actual, fmt.Sprintf("test num: %v", i+1))
	}
}
