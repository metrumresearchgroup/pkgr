package cmd

import (
	"fmt"
	"github.com/dpastoor/goutils"
	"github.com/spf13/afero"
	//"os"
	"path/filepath"
	"testing"
	"time"
)

func InitializeTestEnvironment(testName string) {
	//fsReal := afero.NewOsFs()
	fmt.Println(fmt.Sprintf(""))
	x, _ :=filepath.Abs(".")
	fmt.Println(fmt.Sprintf("start dir: %s", x))
	mm := afero.NewOsFs()//afero.NewMemMapFs()
	testWorkDir := filepath.Join("testsite", "working", testName)
	mm.MkdirAll(testWorkDir, 0755)
	st, err := CopyDir(mm, "testsite/golden/simple/", testWorkDir)
	if err != nil {
		panic(err)
	} else {
		fmt.Println(fmt.Sprintf("%s", st) )
	}
}

func DestroyTestEnvironment() {
	mm := afero.NewOsFs()
	mm.RemoveAll("testsite/working")
}

func TestGetPriorInstalledPackages_BasicTest (t *testing.T) {
	InitializeTestEnvironment("basictest")
	time.Sleep(1000000)
	//DestroyTestEnvironment()
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
		if(item.IsDir()) {
			//fmt.Println(fmt.Sprintf("Making dir %s", dstSubPath))
			fs.Mkdir(dstSubPath, item.Mode())
			CopyDir(fs, srcSubPath, dstSubPath)
		} else {
			//fmt.Println(fmt.Sprintf("Creating file %s", dstSubPath))
			//fmt.Println(fmt.Sprintf("Copying file %s from source %s", dstSubPath, srcSubPath))
			_, err := goutils.CopyFS(fs, srcSubPath, dstSubPath)
			if err != nil {
				fmt.Print("Received error: ")
				fmt.Println(err)
			}
		}
	}

	return "Created " + dst + " from " + src, nil
}