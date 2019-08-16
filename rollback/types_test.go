package rollback

import (
	"github.com/metrumresearchgroup/pkgr/desc"
	"github.com/stretchr/testify/suite"
	"testing"
)

type TypesTestSuite struct {
	suite.Suite
}

func TestTypesTestSuite(t *testing.T) {
	suite.Run(t, new(TypesTestSuite))
}

func (suite *TypesTestSuite) TestDiscernNewPackages_ToInstallCanOutnumberPreinstalled() {
	crayon := desc.Desc{
		Package:    "crayon",
		Version:    "1.3.4",
		Repository: "CRAN",
	}

	r6 := desc.Desc{
		Package:    "R6",
		Version:    "2.4.0",
		Repository: "CRAN",
	}

	toInstallFixture := []string{"R6", "crayon", "shiny"}

	preinstalledPackagesFixture := make(map[string]desc.Desc)
	preinstalledPackagesFixture["crayon"] = crayon
	preinstalledPackagesFixture["R6"] = r6

	actual := DiscernNewPackages(toInstallFixture, preinstalledPackagesFixture)

	suite.Equal(1, len(actual))
	suite.Equal("shiny", actual[0])

}

func (suite *TypesTestSuite) TestDiscernNewPackages_AllPackagesPreinstalled() {
	crayon := desc.Desc{
		Package:    "crayon",
		Version:    "1.3.4",
		Repository: "CRAN",
	}

	r6 := desc.Desc{
		Package:    "R6",
		Version:    "2.4.0",
		Repository: "CRAN",
	}

	toInstallFixture := []string{"R6", "crayon"}

	preinstalledPackagesFixture := make(map[string]desc.Desc)
	preinstalledPackagesFixture["crayon"] = crayon
	preinstalledPackagesFixture["R6"] = r6

	actual := DiscernNewPackages(toInstallFixture, preinstalledPackagesFixture)

	suite.Equal(0, len(actual))
}

func (suite *TypesTestSuite) TestDiscernNewPackages_SomePackagesPreinstalled() {
	crayon := desc.Desc{
		Package:    "crayon",
		Version:    "1.3.4",
		Repository: "CRAN",
	}

	toInstallFixture := []string{"R6", "crayon"}

	preinstalledPackagesFixture := make(map[string]desc.Desc)
	preinstalledPackagesFixture["crayon"] = crayon

	actual := DiscernNewPackages(toInstallFixture, preinstalledPackagesFixture)

	suite.Equal(1, len(actual))
	suite.Equal("R6", actual[0])
}

func (suite *TypesTestSuite) TestDiscernNewPackages_SomePackagesPreinstalled2() {
	crayon := desc.Desc{
		Package:    "crayon",
		Version:    "1.3.4",
		Repository: "CRAN",
	}

	toInstallFixture := []string{"R6", "crayon", "RColorBrewer"}

	preinstalledPackagesFixture := make(map[string]desc.Desc)
	preinstalledPackagesFixture["crayon"] = crayon

	actual := DiscernNewPackages(toInstallFixture, preinstalledPackagesFixture)

	suite.Equal(2, len(actual))
	suite.Equal("R6", actual[0])
	suite.Equal("RColorBrewer", actual[1])
}

func (suite *TypesTestSuite) TestDiscernNewPackages_PackagesAreCaseSensitive() {
	crayon := desc.Desc{
		Package:    "CRAYON",
		Version:    "1.3.4",
		Repository: "CRAN",
	}

	r6 := desc.Desc{
		Package:    "R6",
		Version:    "2.4.0",
		Repository: "CRAN",
	}

	toInstallFixture := []string{"R6", "crayon"}

	preinstalledPackagesFixture := make(map[string]desc.Desc)
	preinstalledPackagesFixture["CRAYON"] = crayon
	preinstalledPackagesFixture["R6"] = r6

	actual := DiscernNewPackages(toInstallFixture, preinstalledPackagesFixture)

	suite.Equal(1, len(actual))
	suite.Equal("crayon", actual[0]) //"crayon" is considered a new package because all we can see is "CRAYON"
}
