package rp

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLineScanning(t *testing.T) {
	assert := assert.New(t)

	var installArgsTests = []struct {
		in       []byte
		expected []string
		context  string
	}{
		{
			[]byte("[1] line 1 info"),
			[]string{"line 1 info"},
			"simplest",
		},
		{
			[]byte(`[1] line 1  
[2]	line 2  	`),
			[]string{
				"line 1",
				"line 2",
			},
			"two lines with whitespace",
		},
		{
			[]byte(`[1] line 1  
[2]	line 2  	
[3]
`),
			[]string{
				"line 1",
				"line 2",
			},
			"two lines with trailing new lines",
		},
	}
	for i, tt := range installArgsTests {
		actual := ScanLines(tt.in)
		assert.Equal(actual, tt.expected, fmt.Sprintf("context: %s, test num: %v", tt.context, i+1))

	}
}
