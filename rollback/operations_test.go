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

	libraryPath, _ := filepath.Abs(filepath.Join(suite.FilePrefix, "test-library"))

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

	libraryPath, _ := filepath.Abs(filepath.Join(suite.FilePrefix, "test-library"))

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

	libraryPath, _ := filepath.Abs(filepath.Join(suite.FilePrefix, "test-library"))

	rbpFixture := RollbackPlan{
		NewPackages: []string {},
		Library: libraryPath,
	}

	RollbackPackageEnvironment(suite.FileSystem, rbpFixture )

	suite.True(afero.DirExists(suite.FileSystem, filepath.Join(suite.FilePrefix, "test-library", "R6")))
	suite.True(afero.DirExists(suite.FileSystem, filepath.Join(suite.FilePrefix, "test-library", "crayon")))
}

/*
// Possible invalid for MacOS
func (suite *OperationsTestSuite) TestRollbackPackageEnvironment_PackagesAreCaseSensitive() {
	InitializeTestEnvironment(suite.FileSystem, "basic-test1")

	cwd, _ := filepath.Abs(".")
	fmt.Println(fmt.Sprintf("Starting test with working directory %s", cwd))

	libraryPath, _ := filepath.Abs(filepath.Join(suite.FilePrefix, "test-library"))

	rbpFixture := RollbackPlan{
		NewPackages: []string {"R6", "CRAYON"},
		Library: libraryPath,
	}

	RollbackPackageEnvironment(suite.FileSystem, rbpFixture)

	suite.False(afero.DirExists(suite.FileSystem, filepath.Join(suite.FilePrefix, "test-library", "R6")))
	suite.True(afero.DirExists(suite.FileSystem, filepath.Join(suite.FilePrefix, "test-library", "crayon")))
}
*/

// Have to use an OS filesystem for this because of this issue: https://github.com/spf13/afero/issues/141
func (suite *OperationsTestSuite) TestRollbackUpdatePackages_RestoresWhenNoActiveInstallation() {

	_ = suite.FileSystem.MkdirAll(filepath.Join(suite.FilePrefix, "test-library", "__OLD__CatsAndOranges"), 0755)
	_, _ = suite.FileSystem.Create(filepath.Join(suite.FilePrefix, "test-library", "__OLD__CatsAndOranges", "DESCRIPTION"))

	updateAttemptFixture := []UpdateAttempt{
		UpdateAttempt{
			Package:                "CatsAndOranges",
			BackupPackageDirectory: filepath.Join(suite.FilePrefix, "test-library", "__OLD__CatsAndOranges"),
			ActivePackageDirectory: filepath.Join(suite.FilePrefix, "test-library", "CatsAndOranges"),
			NewVersion:             "2",
			OldVersion:             "1",
		},
	}
	RollbackUpdatePackages(suite.FileSystem, updateAttemptFixture)

	suite.True(afero.DirExists(suite.FileSystem, filepath.Join(suite.FilePrefix, "test-library", "CatsAndOranges")))
	suite.True(afero.Exists(suite.FileSystem, filepath.Join(suite.FilePrefix, "test-library", "CatsAndOranges", "DESCRIPTION")))
	suite.False(afero.DirExists(suite.FileSystem, filepath.Join(suite.FilePrefix, "test-library", "__OLD__CatsAndOranges")))
	suite.False(afero.Exists(suite.FileSystem, filepath.Join(suite.FilePrefix, "test-library" ,"__OLD__CatsAndOranges","DESCRIPTION")))

}

// Have to use an OS filesystem for this because of this issue: https://github.com/spf13/afero/issues/141
func (suite *OperationsTestSuite) TestRollbackUpdatePackages_OverwritesFreshInstallation() {

	_ = suite.FileSystem.MkdirAll(filepath.Join(suite.FilePrefix, "test-library", "CatsAndOranges"), 0755)
	_, _ = suite.FileSystem.Create(filepath.Join(suite.FilePrefix, "test-library", "CatsAndOranges", "DESCRIPTION_New"))
	_ = suite.FileSystem.MkdirAll(filepath.Join(suite.FilePrefix, "test-library", "__OLD__CatsAndOranges"), 0755)
	_, _ = suite.FileSystem.Create(filepath.Join(suite.FilePrefix, "test-library", "__OLD__CatsAndOranges", "DESCRIPTION"))

	updateAttemptFixture := []UpdateAttempt{
		UpdateAttempt{
			Package:                "CatsAndOranges",
			BackupPackageDirectory: filepath.Join(suite.FilePrefix, "test-library", "__OLD__CatsAndOranges"),
			ActivePackageDirectory: filepath.Join(suite.FilePrefix, "test-library", "CatsAndOranges"),
			NewVersion:             "2",
			OldVersion:             "1",
		},
	}
	RollbackUpdatePackages(suite.FileSystem, updateAttemptFixture)

	suite.True(afero.DirExists(suite.FileSystem, filepath.Join(suite.FilePrefix, "test-library", "CatsAndOranges")))
	suite.True(afero.Exists(suite.FileSystem, filepath.Join(suite.FilePrefix, "test-library", "CatsAndOranges", "DESCRIPTION")))
	suite.False(afero.Exists(suite.FileSystem, filepath.Join(suite.FilePrefix, "test-library", "CatsAndOranges", "DESCRIPTION_New")))
	suite.False(afero.DirExists(suite.FileSystem, filepath.Join(suite.FilePrefix, "test-library", "__OLD__CatsAndOranges")))
	suite.False(afero.Exists(suite.FileSystem, filepath.Join(suite.FilePrefix, "test-library" ,"__OLD__CatsAndOranges","DESCRIPTION")))

}