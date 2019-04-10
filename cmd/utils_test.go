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
	_, notNilError := suite.FileSystem.Stat("test-library/CatsAndOranges")

	suite.True(nilError1 == nil)
	suite.True(nilError2 == nil)
	suite.True(notNilError != nil)


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
