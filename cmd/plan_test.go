package cmd

import (
	"fmt"
	"github.com/metrumresearchgroup/pkgr/desc"
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

	err := CopyDir(fileSystem, goldenSetPath, testWorkDir)

	if err != nil {
		panic(err)
	}
}

func (suite *PlanTestSuite) TestGetPriorInstalledPackages_BasicTest () {
	InitializeTestEnvironment(suite.FileSystem, "basic-test1", "basic-test1")


	cwd, _ := filepath.Abs(".")
	fmt.Println(fmt.Sprintf("Starting test with working directory %s", cwd ))

	expectedResult1 := desc.Desc {
		Package: "crayon",
		Version: "1.3.4",
		Repository: "CRAN",
	}

	expectedResult2 := desc.Desc {
		Package: "R6",
		Version: "2.4.0",
		Repository: "CRAN",
	}

	libraryPath, _ := filepath.Abs(filepath.Join("testsite", "working", "basic-test1", "test-library"))

	actual := GetPriorInstalledPackages(suite.FileSystem, libraryPath)

	suite.Equal(2, len(actual))
	suite.True(installedPackagesAreEqual(expectedResult1, actual[expectedResult1.Package]))
	suite.True(installedPackagesAreEqual(expectedResult2, actual[expectedResult2.Package]))


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

func installedPackagesAreEqual(expected, actual desc.Desc) bool {
	return expected.Package == actual.Package && expected.Version == actual.Version && expected.Repository == actual.Repository
}

