package cmd

import (
	"fmt"
	"github.com/dpastoor/goutils"
	"github.com/spf13/afero"
	//"os"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func InitializeTestEnvironment(goldenSet, testName string) afero.Fs{
	goldenSetPath := filepath.Join("testsite", "golden", goldenSet)
	testWorkDir := filepath.Join("testsite", "working", testName)

	fileSystem := afero.NewOsFs()// afero.NewMemMapFs() //afero.NewOsFs()

	fileSystem.MkdirAll(testWorkDir, 0755)

	st, err := CopyDir(fileSystem, goldenSetPath, testWorkDir)

	if err != nil {
		panic(err)
	} else {
		fmt.Println(fmt.Sprintf("%s", st) )
	}

	return fileSystem
}

func DestroyTestEnvironment() {
	mm := afero.NewOsFs()
	mm.RemoveAll("testsite/working")
}







func TestGetPriorInstalledPackages_BasicTest (t *testing.T) {
	fileSystem := InitializeTestEnvironment("basic-test1", "basic-test1")

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

	actual := GetPriorInstalledPackages(fileSystem, libraryPath)

	assert.Equal(t, 2, len(actual))
	assert.True(t, installedPackagesAreEqual(expectedResult1, actual[expectedResult1.Name]))
	assert.True(t, installedPackagesAreEqual(expectedResult2, actual[expectedResult2.Name]))

	DestroyTestEnvironment()
}

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