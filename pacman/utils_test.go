package pacman

import (
	"testing"

	"github.com/dpastoor/goutils"
	"github.com/metrumresearchgroup/pkgr/cran"
	"github.com/metrumresearchgroup/pkgr/desc"
	"github.com/metrumresearchgroup/pkgr/gpsr"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"
)

type UtilsTestSuite struct {
	suite.Suite
	FileSystem afero.Fs
}

func (suite *UtilsTestSuite) SetupTest() {
	suite.FileSystem = afero.NewMemMapFs()
}

func (suite *UtilsTestSuite) TearDownTest() {
	mm := afero.NewOsFs()
	mm.RemoveAll("testsite/working")
}

func TestUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(UtilsTestSuite))
}

func (suite *UtilsTestSuite) TestTagOldInstallation_CreatesBackup() {
	_ = suite.FileSystem.MkdirAll("test-library/CatsAndOranges", 0755)
	_, _ = suite.FileSystem.Create("test-library/CatsAndOranges/DESCRIPTION")

	outdatedPackageFixture := gpsr.OutdatedPackage{
		Package:    "CatsAndOranges",
		NewVersion: "2",
		OldVersion: "1",
	}

	tagOldInstallation(suite.FileSystem, "test-library", outdatedPackageFixture)

	suite.True(afero.DirExists(suite.FileSystem, "test-library/__OLD__CatsAndOranges"))
	suite.True(afero.Exists(suite.FileSystem, "test-library/__OLD__CatsAndOranges/DESCRIPTION"))
	suite.False(afero.DirExists(suite.FileSystem, "test-library/CatsAndOranges"))
	suite.False(afero.Exists(suite.FileSystem, "test-library/CatsAndOranges/DESCRIPTION"))
}

/*
func (suite *UtilsTestSuite) TestRestoreUnupdatedPackages_RestoresWhenNoActiveInstallation() {
	_ = suite.FileSystem.MkdirAll("test-library/__OLD__CatsAndOranges", 0755)
	_, _ = suite.FileSystem.Create("test-library/__OLD__CatsAndOranges/DESCRIPTION")

	updateAttemptFixture := []UpdateAttempt{
		UpdateAttempt{
			Package:                "CatsAndOranges",
			BackupPackageDirectory: "test-library/__OLD__CatsAndOranges",
			ActivePackageDirectory: "test-library/CatsAndOranges",
			NewVersion:             "2",
			OldVersion:             "1",
		},
	}
	RestoreUnupdatedPackages(suite.FileSystem, updateAttemptFixture)

	suite.True(afero.DirExists(suite.FileSystem, "test-library/CatsAndOranges"))
	suite.True(afero.Exists(suite.FileSystem, "test-library/CatsAndOranges/DESCRIPTION"))
	suite.False(afero.DirExists(suite.FileSystem, "test-library/__OLD__CatsAndOranges"))
	suite.False(afero.Exists(suite.FileSystem, "test-library/__OLD__CatsAndOranges/DESCRIPTION"))

}

func (suite *UtilsTestSuite) TestRestoreUnupdatedPackages_DoesNotRestoreWhenProperlyInstalled() {
	_ = suite.FileSystem.MkdirAll("test-library/__OLD__CatsAndOranges", 0755)
	_, _ = suite.FileSystem.Create("test-library/__OLD__CatsAndOranges/DESCRIPTION_OLD")
	_ = suite.FileSystem.MkdirAll("test-library/CatsAndOranges", 0755)
	_, _ = suite.FileSystem.Create("test-library/CatsAndOranges/DESCRIPTION")

	updateAttemptFixture := []UpdateAttempt{
		UpdateAttempt{
			Package:                "CatsAndOranges",
			BackupPackageDirectory: "test-library/__OLD__CatsAndOranges",
			ActivePackageDirectory: "test-library/CatsAndOranges",
			NewVersion:             "2",
			OldVersion:             "1",
		},
	}
	RestoreUnupdatedPackages(suite.FileSystem, updateAttemptFixture)

	suite.True(afero.DirExists(suite.FileSystem, "test-library/CatsAndOranges"))
	suite.True(afero.Exists(suite.FileSystem, "test-library/CatsAndOranges/DESCRIPTION"))
	suite.False(afero.DirExists(suite.FileSystem, "test-library/__OLD__CatsAndOranges"))
	suite.False(afero.Exists(suite.FileSystem, "test-library/__OLD__CatsAndOranges/DESCRIPTION_OLD"))
	suite.False(afero.Exists(suite.FileSystem, "test-library/CatsAndOranges/DESCRIPTION_OLD"))
}

func (suite *UtilsTestSuite) TestRestoreUnupdatedPackages_RestoresOneAndAllowsAnother() {
	_ = suite.FileSystem.MkdirAll("test-library/__OLD__CatsAndOranges", 0755)
	_, _ = suite.FileSystem.Create("test-library/__OLD__CatsAndOranges/DESCRIPTION_OLD")
	_ = suite.FileSystem.MkdirAll("test-library/__OLD__DogsAndBananas", 0755)
	_, _ = suite.FileSystem.Create("test-library/__OLD__DogsAndBananas/DESCRIPTION_OLD")
	_ = suite.FileSystem.MkdirAll("test-library/DogsAndBananas", 0755)
	_, _ = suite.FileSystem.Create("test-library/DogsAndBananas/DESCRIPTION")

	//_, _ = suite.FileSystem.Create("test-library/DogsAndBananas/DESCRIPTION_OLD")

	updateAttemptFixture := []UpdateAttempt{
		{
			Package:                "CatsAndOranges", ///Not updated successfully
			BackupPackageDirectory: "test-library/__OLD__CatsAndOranges",
			ActivePackageDirectory: "test-library/CatsAndOranges",
			NewVersion:             "2",
			OldVersion:             "1",
		}, {
			Package:                "DogsAndBananas", ///Updated successfully
			BackupPackageDirectory: "test-library/__OLD__DogsAndBananas",
			ActivePackageDirectory: "test-library/DogsAndBananas",
			NewVersion:             "1.5",
			OldVersion:             "1.2",
		},
	}
	RestoreUnupdatedPackages(suite.FileSystem, updateAttemptFixture)

	suite.True(afero.DirExists(suite.FileSystem, "test-library/CatsAndOranges"))
	suite.True(afero.Exists(suite.FileSystem, "test-library/CatsAndOranges/DESCRIPTION_OLD"))
	suite.True(afero.DirExists(suite.FileSystem, "test-library/DogsAndBananas"))
	suite.True(afero.Exists(suite.FileSystem, "test-library/DogsAndBananas/DESCRIPTION"))
	suite.False(afero.DirExists(suite.FileSystem, "test-library/__OLD__CatsAndOranges"))
	suite.False(afero.Exists(suite.FileSystem, "test-library/__OLD__CatsAndOranges/DESCRIPTION_OLD"))
	suite.False(afero.DirExists(suite.FileSystem, "test-library/__OLD__DogsAndBananas"))
	suite.False(afero.Exists(suite.FileSystem, "test-library/__OLD__DogsAndBananas/DESCRIPTION_OLD"))
	suite.False(afero.Exists(suite.FileSystem, "test-library/DogsAndBananas/DESCRIPTION_OLD"))
}
*/
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
/*
func (suite *UtilsTestSuite) TestStringInSlice_FindsStringInSlice() {
	sliceFixture := []string{"Cats", "And", "Oranges"}
	suite.True(stringInSlice("Cats", sliceFixture))
}

func (suite *UtilsTestSuite) TestStringInSlice_DoesNotFindStringNotInSlice() {
	sliceFixture := []string{"Cats", "And", "Oranges"}
	suite.False(stringInSlice("Orangutans", sliceFixture))
}

func (suite *UtilsTestSuite) TestStringInSlice_IsCaseSensitive() {
	sliceFixture := []string{"Cats", "And", "Oranges"}
	suite.False(stringInSlice("cats", sliceFixture))
}
*/