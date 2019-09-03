package cmd

import (
	"github.com/metrumresearchgroup/pkgr/testhelper"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"
	"path/filepath"
	"testing"
)

type CleanCacheSuite struct {
	suite.Suite
	FileSystemOs afero.Fs
	FilePrefix string
}

func (suite *CleanCacheSuite) SetupTest() {
	//suite.FileSystem = afero.NewMemMapFs()
	suite.FileSystemOs = afero.NewOsFs()
	suite.FileSystemOs.MkdirAll("testsite/working", 0755)
	suite.FilePrefix = "testsite/working"
}

func (suite *CleanCacheSuite) TearDownTest() {
	mm := afero.NewOsFs()
	mm.RemoveAll("testsite/working")
}

func TestCleanCacheTestSuite(t *testing.T) {
	suite.Run(t, new(CleanCacheSuite))
}

func InitializeTestEnvironment2(fileSystem afero.Fs, goldenSet, testName string) {
	goldenSetPath := filepath.Join("testsite", "golden", goldenSet)
	testWorkDir := filepath.Join("testsite", "working", testName)
	fileSystem.MkdirAll(testWorkDir, 0755)

	err := testhelper.CopyDir(fileSystem, goldenSetPath, testWorkDir)

	if err != nil {
		panic(err)
	}
}

func (suite *CleanCacheSuite) TestCleanCache_CleansRepoFoldersWhenEmpty() {
	//InitializeTestEnvironment(suite.FileSystemOs, "basic-test1", "basic-test1")
	InitializeTestEnvironment2(suite.FileSystemOs, "cache", "cache")
	cacheDirectory := filepath.Join(suite.FilePrefix, "cache")

	repo1, _ := filepath.Abs(filepath.Join(suite.FilePrefix, "cache", "CRAN"))
	repo2, _ := filepath.Abs(filepath.Join(suite.FilePrefix, "cache", "CRAN-Micro"))
	repo3, _ := filepath.Abs(filepath.Join(suite.FilePrefix, "cache", "r_validated"))
	repo4, _ := filepath.Abs(filepath.Join(suite.FilePrefix, "cache", "empty_repo"))

	actual1 := deleteCacheSubfolders(suite.FileSystemOs, nil, "src", cacheDirectory)
	actual2 := deleteCacheSubfolders(suite.FileSystemOs, nil, "binary", cacheDirectory)

	suite.True(actual1 == nil)
	suite.True(actual2 == nil)
	suite.False(afero.Exists(suite.FileSystemOs, repo1))
	suite.False(afero.Exists(suite.FileSystemOs, repo2))
	suite.False(afero.Exists(suite.FileSystemOs, repo3))
	suite.False(afero.Exists(suite.FileSystemOs, repo4))

}

func (suite *CleanCacheSuite) TestCleanCache_DoesNotDeleteNonEmptyRepoFolders() {
	InitializeTestEnvironment2(suite.FileSystemOs, "cache", "cache")
	cacheDirectory := filepath.Join(suite.FilePrefix, "cache")


	repo1, _ := filepath.Abs(filepath.Join(suite.FilePrefix, "cache", "CRAN"))
	repo2, _ := filepath.Abs(filepath.Join(suite.FilePrefix, "cache", "CRAN-Micro"))
	repo3, _ := filepath.Abs(filepath.Join(suite.FilePrefix, "cache", "r_validated"))
	repo4, _ := filepath.Abs(filepath.Join(suite.FilePrefix, "cache", "empty_repo"))

	actual := deleteCacheSubfolders(suite.FileSystemOs, nil, "binary", cacheDirectory)

	suite.True(actual == nil)
	suite.False(afero.Exists(suite.FileSystemOs, repo1)) // Only has binary, should be deleted
	suite.True(afero.Exists(suite.FileSystemOs, repo2)) // Has binary and src, should remain
	suite.True(afero.Exists(suite.FileSystemOs, repo3)) // Has binary and src, should remain.
	suite.False(afero.Exists(suite.FileSystemOs, repo4)) // Empty, should get removed as part of this process

}

func (suite *CleanCacheSuite) TestCleanCache_DeletesSpecificRepos() {
	InitializeTestEnvironment2(suite.FileSystemOs, "cache", "cache")
	cacheDirectory := filepath.Join(suite.FilePrefix, "cache")

	reposFixture := []string{"CRAN", "r_validated"}

	repo1, _ := filepath.Abs(filepath.Join(suite.FilePrefix, "cache", "CRAN"))
	repo2, _ := filepath.Abs(filepath.Join(suite.FilePrefix, "cache", "CRAN-Micro"))
	repo3, _ := filepath.Abs(filepath.Join(suite.FilePrefix, "cache", "r_validated"))
	repo4, _ := filepath.Abs(filepath.Join(suite.FilePrefix, "cache", "empty_repo"))

	actual1 := deleteCacheSubfolders(suite.FileSystemOs, reposFixture, "binary", cacheDirectory)
	actual2 := deleteCacheSubfolders(suite.FileSystemOs, reposFixture, "src", cacheDirectory)

	suite.True(actual1 == nil)
	suite.True(actual2 == nil)
	suite.False(afero.Exists(suite.FileSystemOs, repo1)) //CRAN should be deleted
	suite.True(afero.Exists(suite.FileSystemOs, repo2))
	suite.False(afero.Exists(suite.FileSystemOs, repo3)) //r_validated should be deleted
	suite.True(afero.Exists(suite.FileSystemOs, repo4)) //Empty repo should be ignored in this case

}

