package pacman

import (
	"testing"

	"github.com/dpastoor/goutils"
	"github.com/metrumresearchgroup/pkgr/cran"
	"github.com/metrumresearchgroup/pkgr/desc"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"
)

type UtilsTestSuite struct {
	suite.Suite
	FileSystem afero.Fs
	FileSystemOs afero.Fs
	FilePrefix string
}

func (suite *UtilsTestSuite) SetupTest() {
	suite.FileSystem = afero.NewMemMapFs()
	suite.FileSystemOs = afero.NewOsFs()
	suite.FileSystem.MkdirAll("testsite/working", 0755)
	suite.FilePrefix = "testsite/working"
}

func (suite *UtilsTestSuite) TearDownTest() {
	mm := afero.NewOsFs()
	mm.RemoveAll("testsite/working")
}

func TestUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(UtilsTestSuite))
}

func (suite *UtilsTestSuite) TestScanInstalledPackage_ScansReleventFieldsForOutdatedComparison() {
	_ = suite.FileSystem.MkdirAll("test-library/CatsAndOranges", 0755)
	_, _ = suite.FileSystem.Create("test-library/CatsAndOranges/DESCRIPTION")

	descriptionFileContents := []string{
		"Package: releasy",
		"Title: Simple Functions to Make and Download GitHub Releases",
		"Version: 0.0.0.9000",
		"Authors@R:",
		"    person(given = \"CatsAndOranges\",",
		"           family = \"Guy\",",
		"           role = c(\"aut\", \"cre\"),",
		"           email = \"animalsAndFruitsMakeGreatNonsense@hotmail.com\")",
		"Description: A description, please",
		"License: What license it uses",
		"Encoding: UTF-8",
		"LazyData: true",
		"Imports:",
		"    httr,",
		"    gh,",
		"    purrr,",
		"    yaml,",
		"   config",
		"RoxygenNote: 6.1.1",
		"Suggests:",
		"    testthat",
	}
	/*
		expectedResults := desc.Desc {
			Package     : "releasy",
			//Source      : "",
			Version     : "0.0.0.9000",
			//Maintainer  : "",
			Description : "A description, please",
			//MD5sum      : "",
			//Remotes     : "",
			//Repository  : "",
			//Imports     :
			//Suggests    :
			//Depends     :
			//LinkingTo   :
		}
	*/
	goutils.WriteLinesFS(suite.FileSystem, descriptionFileContents, "test-library/CatsAndOranges/DESCRIPTION")

	actualResults, err := scanInstalledPackage("test-library/CatsAndOranges/DESCRIPTION", suite.FileSystem)

	suite.Nil(err)
	suite.Equal(actualResults.Version, "0.0.0.9000")
	suite.Equal(actualResults.Package, "releasy")
}

func (suite *UtilsTestSuite) TestScanInstalledPackage_ReturnsNilWhenNoDescriptionFileFound() {
	_ = suite.FileSystem.MkdirAll("test-library/CatsAndOranges", 0755)
	//_, _ = suite.FileSystem.Create("test-library/CatsAndOranges/DESCRIPTION")

	actualResults, err := scanInstalledPackage("test-library/CatsAndOranges/DESCRIPTION", suite.FileSystem)

	suite.NotNil(err)
	suite.Equal(actualResults.Version, "")
	suite.Equal(actualResults.Package, "")
}

func (suite *UtilsTestSuite) TestGetOutdatedPackages_FindsOutdatedPackage() {
	outdatedDescFixture := desc.Desc{
		Package: "CatsAndOranges",
		Version: "1.0.1",
	}
	installedFixture := make(map[string]desc.Desc)
	installedFixture["CatsAndOranges"] = outdatedDescFixture

	updatedDescFixture := desc.Desc{
		Package: "CatsAndOranges",
		Version: "1.0.2",
	}

	var availablePackagesFixture []cran.PkgDl
	updatedPkgDlFixture := cran.PkgDl{Package: updatedDescFixture}
	availablePackagesFixture = append(availablePackagesFixture, updatedPkgDlFixture)

	actualResults := GetOutdatedPackages(installedFixture, availablePackagesFixture)

	suite.Equal(1, len(actualResults))
	suite.Equal("CatsAndOranges", actualResults[0].Package)
	suite.Equal("1.0.1", actualResults[0].OldVersion)
	suite.Equal("1.0.2", actualResults[0].NewVersion)

}

func (suite *UtilsTestSuite) TestGetOutdatedPackages_DoesNotFlagOlderPackage() {
	outdatedDescFixture := desc.Desc{
		Package: "CatsAndOranges",
		Version: "1.0.1",
	}
	installedFixture := make(map[string]desc.Desc)
	installedFixture["CatsAndOranges"] = outdatedDescFixture

	olderDescFixture := desc.Desc{
		Package: "CatsAndOranges",
		Version: "1.0.0",
	}

	var availablePackagesFixture []cran.PkgDl
	olderPkgDlFixture := cran.PkgDl{Package: olderDescFixture}
	availablePackagesFixture = append(availablePackagesFixture, olderPkgDlFixture)

	actualResults := GetOutdatedPackages(installedFixture, availablePackagesFixture)

	suite.Equal(0, len(actualResults))
}
