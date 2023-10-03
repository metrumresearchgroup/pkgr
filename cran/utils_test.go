package cran

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUrlHash(t *testing.T) {
	assert := assert.New(t)

	var installArgsTests = []struct {
		in       RepoURL
		expected string
	}{
		{
			RepoURL{
				Name: "MPN",
				URL:  "https://mpn.metworx.com/snapshots/stable/2023-05-14",
			},
			"MPN-8520d4ecc108",
		},
		{
			RepoURL{
				Name: "CRAN",
				URL:  "https://cran.rstudio.com",
			},
			"CRAN-739227e5b53e",
		},
		{
			RepoURL{
				Name: "gh_dev",
				URL:  "https://metrumresearchgroup.github.io/rpkgs/gh_dev",
			},
			"gh_dev-a1f00a415a5e",
		},
	}
	for i, tt := range installArgsTests {
		actual := RepoURLHash(tt.in)
		assert.Equal(actual, tt.expected, fmt.Sprintf("test num: %v", i+1))

	}
}
