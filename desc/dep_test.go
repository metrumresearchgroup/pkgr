package desc

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type DepTestSuite struct {
	suite.Suite
	versionFixture Version
}

func (suite *DepTestSuite) SetupTest() {
	suite.versionFixture = Version{
		Major: 2,
		Minor: 3,
		Patch: 1,
		String: "2.3.1",
	}
}

func TestDepTestSuite(t *testing.T) {
	suite.Run(t, new(DepTestSuite))
}

func (suite *DepTestSuite) TestDepToString_GTConstraint() {

	fixture := Dep{
		Version: suite.versionFixture,
		Constraint: GT,
		Name: "CatsAndOranges",
	}

	expected := "CatsAndOranges (> 2.3.1)"
	actual := fixture.ToString()

	suite.Equal(expected, actual)
}

func (suite *DepTestSuite) TestDepToString_GTEConstraint() {

	fixture := Dep{
		Version: suite.versionFixture,
		Constraint: GTE,
		Name: "CatsAndOranges",
	}

	expected := "CatsAndOranges (>= 2.3.1)"
	actual := fixture.ToString()

	suite.Equal(expected, actual)
}

func (suite *DepTestSuite) TestDepToString_LTConstaint() {

	fixture := Dep{
		Version: suite.versionFixture,
		Constraint: LT,
		Name: "CatsAndOranges",
	}

	expected := "CatsAndOranges (< 2.3.1)"
	actual := fixture.ToString()

	suite.Equal(expected, actual)
}

func (suite *DepTestSuite) TestDepToString_LTEConstaint() {

	fixture := Dep{
		Version: suite.versionFixture,
		Constraint: LTE,
		Name: "CatsAndOranges",
	}

	expected := "CatsAndOranges (<= 2.3.1)"
	actual := fixture.ToString()

	suite.Equal(expected, actual)
}


func (suite *DepTestSuite) TestDepToString_EqualsConstaint() {

	fixture := Dep{
		Version: suite.versionFixture,
		Constraint: Equals,
		Name: "CatsAndOranges",
	}

	expected := "CatsAndOranges (= 2.3.1)"
	actual := fixture.ToString()

	suite.Equal(expected, actual)
}

func (suite *DepTestSuite) TestDepToString_VersionContainsDev() {

	versionFixtureWithDev := suite.versionFixture
	versionFixtureWithDev.Dev = 11
	versionFixtureWithDev.String = "2.3.1.11"

	fixture := Dep{
		Version: versionFixtureWithDev,
		Constraint: GTE,
		Name: "CatsAndOranges",
	}

	expected := "CatsAndOranges (>= 2.3.1.11)"
	actual := fixture.ToString()

	suite.Equal(expected, actual)
}