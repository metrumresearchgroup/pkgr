package cmd

import (
	"fmt"
	"github.com/dpastoor/goutils"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"

	"path/filepath"
	"testing"
)

type PlanTestSuite struct {
	suite.Suite
	FileSystem afero.Fs
}

func (suite *PlanTestSuite) SetupTest() {
	suite.FileSystem = afero.NewOsFs()
}

func (suite *PlanTestSuite) TearDownTest() {
	mm := afero.NewOsFs()
	mm.RemoveAll("testsite/working")
}


func TestPlanTestSuite(t *testing.T) {
	suite.Run(t, new(PlanTestSuite))
}

func InitializeTestEnvironment(fileSystem afero.Fs, goldenSet, testName string) {
	goldenSetPath := filepath.Join("testsite", "golden", goldenSet)
	testWorkDir := filepath.Join("testsite", "working", testName)
	fileSystem.MkdirAll(testWorkDir, 0755)

	st, err := CopyDir(fileSystem, goldenSetPath, testWorkDir)

	if err != nil {
		panic(err)
	} else {
		fmt.Println(fmt.Sprintf("%s", st) )
	}
}

func (suite *PlanTestSuite) TestGetPriorInstalledPackages_BasicTest () {
	InitializeTestEnvironment(suite.FileSystem, "basic-test1", "basic-test1")


	cwd, _ := filepath.Abs(".")
	fmt.Println(fmt.Sprintf("Starting test with working directory %s", cwd ))

	expectedResult1 := InstalledPackage {
		Name: "crayon",
		Version: "1.3.4",
		Repo: "CRAN",
	}

	expectedResult2 := InstalledPackage {
		Name: "R6",
		Version: "2.4.0",
		Repo: "CRAN",
	}

	libraryPath, _ := filepath.Abs(filepath.Join("testsite", "working", "basic-test1", "test-library"))

	actual := GetPriorInstalledPackages(suite.FileSystem, libraryPath)

	suite.Equal(2, len(actual))
	suite.True(installedPackagesAreEqual(expectedResult1, actual[expectedResult1.Name]))
	suite.True(installedPackagesAreEqual(expectedResult2, actual[expectedResult2.Name]))


}

func (suite *PlanTestSuite) TestGetPriorInstalledPackages_NoPreinstalledPackages() {
	InitializeTestEnvironment(suite.FileSystem, "null-test", "null-test")

	cwd, _ := filepath.Abs(".")
	fmt.Println(fmt.Sprintf("Starting test with working directory %s", cwd ))


	libraryPath, _ := filepath.Abs(filepath.Join("testsite", "working", "null-test", "test-library"))

	actual := GetPriorInstalledPackages(suite.FileSystem, libraryPath)

	suite.Equal(0, len(actual))
}




//////// Utility

func installedPackagesAreEqual(expected, actual InstalledPackage) bool {
	return expected.Name == actual.Name && expected.Version == actual.Version && expected.Repo == actual.Repo
}


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
