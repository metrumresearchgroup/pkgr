package rollback

import (
	"fmt"
	"github.com/metrumresearchgroup/pkgr/testhelper"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"
	"path/filepath"
	"testing"
)

type OperationsTestSuite struct {
	suite.Suite
	FileSystem afero.Fs
	FilePrefix string
}

func (suite *OperationsTestSuite) SetupTest() {
	suite.FileSystem = afero.NewOsFs()
	suite.FilePrefix = "testsite/working"
}

func (suite *OperationsTestSuite) TearDownTest() {
	mm := afero.NewOsFs()
	mm.RemoveAll("testsite/working")
}

func TestOperationsTestSuite(t *testing.T) {
	suite.Run(t, new(OperationsTestSuite))
}

func InitializeTestEnvironment(fileSystem afero.Fs, goldenSet string) {
	goldenSetPath := filepath.Join("testsite", "golden", goldenSet)
	testWorkDir := filepath.Join("testsite", "working")
	fileSystem.MkdirAll(testWorkDir, 0755)

	err := testhelper.CopyDir(fileSystem, goldenSetPath, testWorkDir)

	if err != nil {
		panic(err)
	}
}

func (suite *OperationsTestSuite) TestRollbackPackageEnvironment_DeletesOnlyNewPackages() {
	InitializeTestEnvironment(suite.FileSystem, "basic-test1")

	cwd, _ := filepath.Abs(".")
	fmt.Println(fmt.Sprintf("Starting test with working directory %s", cwd))

	libraryPath, _ := filepath.Abs(filepath.Join("testsite", "working", "test-library"))

	rbpFixture := RollbackPlan{
		NewPackages: []string {"R6"},
		Library: libraryPath,
	}

	RollbackPackageEnvironment(suite.FileSystem, rbpFixture)

	suite.False(afero.DirExists(suite.FileSystem, filepath.Join(suite.FilePrefix, "test-library", "R6")))
	suite.True(afero.DirExists(suite.FileSystem, filepath.Join(suite.FilePrefix, "test-library", "crayon")))
}

func (suite *OperationsTestSuite) TestRollbackPackageEnvironment_DeletesMultiplePackages() {
	InitializeTestEnvironment(suite.FileSystem, "basic-test1")

	cwd, _ := filepath.Abs(".")
	fmt.Println(fmt.Sprintf("Starting test with working directory %s", cwd))

	libraryPath, _ := filepath.Abs(filepath.Join("testsite", "working", "test-library"))

	rbpFixture := RollbackPlan{
		NewPackages: []string {"R6", "crayon"},
		Library: libraryPath,
	}

	RollbackPackageEnvironment(suite.FileSystem, rbpFixture)

	suite.False(afero.DirExists(suite.FileSystem, filepath.Join(suite.FilePrefix, "test-library", "R6")))
	suite.False(afero.DirExists(suite.FileSystem, filepath.Join(suite.FilePrefix, "test-library", "crayon")))
}

func (suite *OperationsTestSuite) TestRollbackPackageEnvironment_HandlesEmptyListOfPackages() {
	InitializeTestEnvironment(suite.FileSystem, "basic-test1")

	cwd, _ := filepath.Abs(".")
	fmt.Println(fmt.Sprintf("Starting test with working directory %s", cwd))

	libraryPath, _ := filepath.Abs(filepath.Join(suite.FilePrefix, "basic-test1", "test-library"))

	rbpFixture := RollbackPlan{
		NewPackages: []string {},
		Library: libraryPath,
	}

	RollbackPackageEnvironment(suite.FileSystem, rbpFixture )

	suite.True(afero.DirExists(suite.FileSystem, filepath.Join(suite.FilePrefix, "test-library", "R6")))
	suite.True(afero.DirExists(suite.FileSystem, filepath.Join(suite.FilePrefix, "test-library", "crayon")))
}

func (suite *OperationsTestSuite) TestRollbackPackageEnvironment_PackagesAreCaseSensitive() {
	InitializeTestEnvironment(suite.FileSystem, "basic-test1")

	cwd, _ := filepath.Abs(".")
	fmt.Println(fmt.Sprintf("Starting test with working directory %s", cwd))

	libraryPath, _ := filepath.Abs(filepath.Join("testsite", "working", "test-library"))

	rbpFixture := RollbackPlan{
		NewPackages: []string {"R6", "CRAYON"},
		Library: libraryPath,
	}

	RollbackPackageEnvironment(suite.FileSystem, rbpFixture)

	suite.False(afero.DirExists(suite.FileSystem, filepath.Join(suite.FilePrefix, "test-library", "R6")))
	suite.True(afero.DirExists(suite.FileSystem, filepath.Join(suite.FilePrefix, "test-library", "crayon")))
}
