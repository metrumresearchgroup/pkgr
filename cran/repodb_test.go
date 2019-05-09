package cran

import (
	"github.com/metrumresearchgroup/pkgr/desc"
	"github.com/stretchr/testify/suite"
	"testing"
)

type RepoDbTestSuite struct {
	suite.Suite
	highVersionFixture desc.Version
	lowVersionFixture desc.Version
}

func (suite *RepoDbTestSuite) SetupTest() {
	suite.highVersionFixture = desc.Version{
		Major: 99,
		Minor: 99,
		Patch: 99,
		String: "99.99.99",
	}
	suite.lowVersionFixture = desc.Version{
		Major: 0,
		Minor: 1,
		Patch: 0,
		String: "0.1.0",
	}
}

func TestRepoDbTestSuite(t *testing.T) {
	suite.Run(t, new(RepoDbTestSuite))
}

func (suite *RepoDbTestSuite) TestIsRVersionCompatible_VersionIsGreaterThanRequired() {

	rVersionFixture := RVersion {
		Major: suite.highVersionFixture.Major,
		Minor: suite.highVersionFixture.Minor,
		Patch: suite.highVersionFixture.Patch,
	}

	depFixture := desc.Dep{
		Version: suite.lowVersionFixture,
		Name: "CatsAndOranges",
		Constraint: desc.GTE,
	}

	// rVersion RVersion desc.Dep
	actual := checkRVersionCompatibility(rVersionFixture, depFixture)

	suite.True(actual, "R package is valid")


}