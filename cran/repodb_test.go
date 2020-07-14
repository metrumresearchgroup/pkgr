package cran

import (
	"github.com/metrumresearchgroup/pkgr/desc"
	"github.com/stretchr/testify/suite"
	"testing"
)

type RepoDbTestSuite struct {
	suite.Suite
	highVersionFixture desc.Version
	lowVersionFixture  desc.Version
}

func (suite *RepoDbTestSuite) SetupTest() {
	suite.highVersionFixture = desc.Version{
		Major:  99,
		Minor:  99,
		Patch:  99,
		String: "99.99.99",
	}
	suite.lowVersionFixture = desc.Version{
		Major:  0,
		Minor:  1,
		Patch:  0,
		String: "0.1.0",
	}
}

func TestRepoDbTestSuite(t *testing.T) {
	suite.Run(t, new(RepoDbTestSuite))
}

func (suite *RepoDbTestSuite) TestIsRVersionCompatible_GTEInstalledVersionGreaterThan() {

	rVersionFixture := RVersion{
		Major: suite.highVersionFixture.Major,
		Minor: suite.highVersionFixture.Minor,
		Patch: suite.highVersionFixture.Patch,
	}

	depFixture := desc.Dep{
		Version:    suite.lowVersionFixture,
		Name:       "CatsAndOranges",
		Constraint: desc.GTE,
	}

	// rVersion RVersion desc.Dep
	actual := checkRVersionCompatibility(rVersionFixture, depFixture)

	suite.True(actual, "R package is valid")
}

func (suite *RepoDbTestSuite) TestIsRVersionCompatible_GTEInstalledVersionEquals() {

	rVersionFixture := RVersion{
		Major: suite.highVersionFixture.Major,
		Minor: suite.highVersionFixture.Minor,
		Patch: suite.highVersionFixture.Patch,
	}

	depFixture := desc.Dep{
		Version:    suite.highVersionFixture,
		Name:       "CatsAndOranges",
		Constraint: desc.GTE,
	}

	// rVersion RVersion desc.Dep
	actual := checkRVersionCompatibility(rVersionFixture, depFixture)

	suite.True(actual, "R package is valid")
}

func (suite *RepoDbTestSuite) TestIsRVersionCompatible_LTEInstalledVersionEquals() {

	rVersionFixture := RVersion{
		Major: suite.highVersionFixture.Major,
		Minor: suite.highVersionFixture.Minor,
		Patch: suite.highVersionFixture.Patch,
	}

	depFixture := desc.Dep{
		Version:    suite.highVersionFixture,
		Name:       "CatsAndOranges",
		Constraint: desc.LTE,
	}

	// rVersion RVersion desc.Dep
	actual := checkRVersionCompatibility(rVersionFixture, depFixture)

	suite.True(actual, "R package is valid")
}

func (suite *RepoDbTestSuite) TestIsRVersionCompatible_LTInstalledVersionLessThan() {

	rVersionFixture := RVersion{
		Major: suite.lowVersionFixture.Major,
		Minor: suite.lowVersionFixture.Minor,
		Patch: suite.lowVersionFixture.Patch,
	}

	depFixture := desc.Dep{
		Version:    suite.highVersionFixture,
		Name:       "CatsAndOranges",
		Constraint: desc.LT,
	}

	// rVersion RVersion desc.Dep
	actual := checkRVersionCompatibility(rVersionFixture, depFixture)

	suite.True(actual, "R package is valid")
}

func (suite *RepoDbTestSuite) TestIsRVersionCompatible_EqualsInstalledVersionEquals() {

	rVersionFixture := RVersion{
		Major: suite.lowVersionFixture.Major,
		Minor: suite.lowVersionFixture.Minor,
		Patch: suite.lowVersionFixture.Patch,
	}

	depFixture := desc.Dep{
		Version:    suite.lowVersionFixture,
		Name:       "CatsAndOranges",
		Constraint: desc.Equals,
	}

	// rVersion RVersion desc.Dep
	actual := checkRVersionCompatibility(rVersionFixture, depFixture)

	suite.True(actual, "R package is valid")
}

/// Invalid cases

func (suite *RepoDbTestSuite) TestIsRVersionCompatible_InvalidGTEInstalledVersionLower() {

	rVersionFixture := RVersion{
		Major: suite.lowVersionFixture.Major,
		Minor: suite.lowVersionFixture.Minor,
		Patch: suite.lowVersionFixture.Patch,
	}

	depFixture := desc.Dep{
		Version:    suite.highVersionFixture,
		Name:       "CatsAndOranges",
		Constraint: desc.GTE,
	}

	// rVersion RVersion desc.Dep
	actual := checkRVersionCompatibility(rVersionFixture, depFixture)

	suite.False(actual, "R package is invalid")
}

func (suite *RepoDbTestSuite) TestIsRVersionCompatible_InvalidGTInstalledVersionLower() {

	rVersionFixture := RVersion{
		Major: suite.lowVersionFixture.Major,
		Minor: suite.lowVersionFixture.Minor,
		Patch: suite.lowVersionFixture.Patch,
	}

	depFixture := desc.Dep{
		Version:    suite.highVersionFixture,
		Name:       "CatsAndOranges",
		Constraint: desc.GT,
	}

	// rVersion RVersion desc.Dep
	actual := checkRVersionCompatibility(rVersionFixture, depFixture)

	suite.False(actual, "R package is invalid")
}

func (suite *RepoDbTestSuite) TestIsRVersionCompatible_InvalidGTInstalledVersionEquals() {

	rVersionFixture := RVersion{
		Major: suite.lowVersionFixture.Major,
		Minor: suite.lowVersionFixture.Minor,
		Patch: suite.lowVersionFixture.Patch,
	}

	depFixture := desc.Dep{
		Version:    suite.lowVersionFixture,
		Name:       "CatsAndOranges",
		Constraint: desc.GT,
	}

	// rVersion RVersion desc.Dep
	actual := checkRVersionCompatibility(rVersionFixture, depFixture)

	suite.False(actual, "R package is invalid")
}

func (suite *RepoDbTestSuite) TestIsRVersionCompatible_InvalidLTEInstalledVersionHigher() {

	rVersionFixture := RVersion{
		Major: suite.highVersionFixture.Major,
		Minor: suite.highVersionFixture.Minor,
		Patch: suite.highVersionFixture.Patch,
	}

	depFixture := desc.Dep{
		Version:    suite.lowVersionFixture,
		Name:       "CatsAndOranges",
		Constraint: desc.LTE,
	}

	// rVersion RVersion desc.Dep
	actual := checkRVersionCompatibility(rVersionFixture, depFixture)

	suite.False(actual, "R package is invalid")
}

func (suite *RepoDbTestSuite) TestIsRVersionCompatible_InvalidLTInstalledVersionEqual() {

	rVersionFixture := RVersion{
		Major: suite.highVersionFixture.Major,
		Minor: suite.highVersionFixture.Minor,
		Patch: suite.highVersionFixture.Patch,
	}

	depFixture := desc.Dep{
		Version:    suite.highVersionFixture,
		Name:       "CatsAndOranges",
		Constraint: desc.LT,
	}

	// rVersion RVersion desc.Dep
	actual := checkRVersionCompatibility(rVersionFixture, depFixture)

	suite.False(actual, "R package is invalid")
}

func (suite *RepoDbTestSuite) TestIsRVersionCompatible_InvalidLTInstalledVersionHigher() {

	rVersionFixture := RVersion{
		Major: suite.highVersionFixture.Major,
		Minor: suite.highVersionFixture.Minor,
		Patch: suite.highVersionFixture.Patch,
	}

	depFixture := desc.Dep{
		Version:    suite.lowVersionFixture,
		Name:       "CatsAndOranges",
		Constraint: desc.LT,
	}

	// rVersion RVersion desc.Dep
	actual := checkRVersionCompatibility(rVersionFixture, depFixture)

	suite.False(actual, "R package is invalid")
}

func (suite *RepoDbTestSuite) TestIsRVersionCompatible_InvalidEqualsInstalledVersionDifferent() {

	rVersionFixture := RVersion{
		Major: suite.highVersionFixture.Major,
		Minor: suite.highVersionFixture.Minor,
		Patch: suite.highVersionFixture.Patch,
	}

	depFixture := desc.Dep{
		Version:    suite.lowVersionFixture,
		Name:       "CatsAndOranges",
		Constraint: desc.Equals,
	}

	// rVersion RVersion desc.Dep
	actual := checkRVersionCompatibility(rVersionFixture, depFixture)

	suite.False(actual, "R package is invalid")
}
