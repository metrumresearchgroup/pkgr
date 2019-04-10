package cmd

import (
	"github.com/metrumresearchgroup/pkgr/gpsr"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"
	"testing"
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
		Package: "CatsAndOranges",
		NewVersion: "2",
		OldVersion: "1",
	}

	tagOldInstallation(suite.FileSystem, "test-library", outdatedPackageFixture)

	_, nilError1 := suite.FileSystem.Stat("test-library/__OLD__CatsAndOranges")
	_, nilError2 := suite.FileSystem.Stat("test-library/__OLD__CatsAndOranges/DESCRIPTION")
	_, notNilError1 := suite.FileSystem.Stat("test-library/CatsAndOranges")
	_, notNilError2 := suite.FileSystem.Stat("test-library/CatsAndOranges/DESCRIPTION")

	suite.True(nilError1 == nil)
	suite.True(nilError2 == nil)
	suite.True(notNilError1 != nil)
	suite.True(notNilError2 != nil)
}

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
	restoreUnupdatedPackages(suite.FileSystem, updateAttemptFixture)

	_, nilError1 := suite.FileSystem.Stat("test-library/CatsAndOranges")
	_, nilError2 := suite.FileSystem.Stat("test-library/CatsAndOranges/DESCRIPTION")
	_, notNilError1 := suite.FileSystem.Stat("test-library/__OLD__CatsAndOranges")
	_, notNilError2 := suite.FileSystem.Stat("test-library/__OLD__CatsAndOranges/DESCRIPTION")

	suite.True(nilError1 == nil)
	suite.True(nilError2 == nil)
	suite.True(notNilError1 != nil)
	suite.True(notNilError2 != nil)

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
	restoreUnupdatedPackages(suite.FileSystem, updateAttemptFixture)

	_, nilError1 := suite.FileSystem.Stat("test-library/CatsAndOranges")
	_, nilError2 := suite.FileSystem.Stat("test-library/CatsAndOranges/DESCRIPTION")
	_, notNilError1 := suite.FileSystem.Stat("test-library/__OLD__CatsAndOranges")
	_, notNilError2 := suite.FileSystem.Stat("test-library/__OLD__CatsAndOranges/DESCRIPTION_OLD")
	_, notNilError3 := suite.FileSystem.Stat("test-library/CatsAndOranges/DESCRIPTION_OLD")

	suite.True(nilError1 == nil)
	suite.True(nilError2 == nil)
	suite.True(notNilError1 != nil)
	suite.True(notNilError2 != nil)
	suite.True(notNilError3 != nil)

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
			Package:                "CatsAndOranges",	///Not updated successfully
			BackupPackageDirectory: "test-library/__OLD__CatsAndOranges",
			ActivePackageDirectory: "test-library/CatsAndOranges",
			NewVersion:             "2",
			OldVersion:             "1",
		}, {
			Package: 				"DogsAndBananas",	///Updated successfully
			BackupPackageDirectory: "test-library/__OLD__DogsAndBananas",
			ActivePackageDirectory: "test-library/DogsAndBananas",
		NewVersion: 				"1.5",
			OldVersion: 			"1.2",
		},
	}
	restoreUnupdatedPackages(suite.FileSystem, updateAttemptFixture)

	//Expect to find
	_, nilError1 := suite.FileSystem.Stat("test-library/CatsAndOranges")
	_, nilError2 := suite.FileSystem.Stat("test-library/CatsAndOranges/DESCRIPTION_OLD")
	_, nilError3 := suite.FileSystem.Stat("test-library/DogsAndBananas")
	_, nilError4 := suite.FileSystem.Stat("test-library/DogsAndBananas/DESCRIPTION")

	//Expect not to find
	_, notNilError1 := suite.FileSystem.Stat("test-library/__OLD__CatsAndOranges")
	_, notNilError2 := suite.FileSystem.Stat("test-library/__OLD__CatsAndOranges/DESCRIPTION_OLD")
	_, notNilError3 := suite.FileSystem.Stat("test-library/__OLD__DogsAndBananas")
	_, notNilError4 := suite.FileSystem.Stat("test-library/__OLD__DogsAndBananas/DESCRIPTION_OLD")
	_, notNilError5 := suite.FileSystem.Stat("test-library/DogsAndBananas/DESCRIPTION_OLD")

	suite.True(nilError1 == nil)
	suite.True(nilError2 == nil)
	suite.True(nilError3 == nil)
	suite.True(nilError4 == nil)
	suite.True(notNilError1 != nil)
	suite.True(notNilError2 != nil)
	suite.True(notNilError3 != nil)
	suite.True(notNilError4 != nil)
	suite.True(notNilError5 != nil)

}



//////// Utility

/*
func installedPackagesAreEqual(expected, actual desc.Desc) bool {
	return expected.Package == actual.Package && expected.Version == actual.Version && expected.Repository == actual.Repository
}
*/
/*
func CopyDir(fs afero.Fs, src string, dst string) (string, error) {
	openedDir, err := fs.Open(src)
	if err != nil {
		return "", err
	}

	stuff, err := openedDir.Readdir(0)

	if err != nil {
		return "", err
	}

	for _, item := range(stuff) {
		srcSubPath := filepath.Join(src, item.Name())
		dstSubPath := filepath.Join(dst, item.Name())
		if item.IsDir() {
			fs.Mkdir(dstSubPath, item.Mode())
			CopyDir(fs, srcSubPath, dstSubPath)
		} else {
			_, err := goutils.CopyFS(fs, srcSubPath, dstSubPath)
			if err != nil {
				fmt.Print("Received error: ")
				fmt.Println(err)
			}
		}
	}

	return "Created " + dst + " from " + src, nil
}
*/
