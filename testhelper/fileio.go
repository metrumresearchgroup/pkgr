package testhelper

import (
	"fmt"
	"github.com/dpastoor/goutils"
	"github.com/spf13/afero"
	"path/filepath"
)

func CopyDir(fs afero.Fs, src string, dst string) error {

	err := fs.MkdirAll(dst, 0755)
	if err != nil {
		return err
	}

	openedDir, err := fs.Open(src)
	if err != nil {
		return err
	}

	directoryContents, err := openedDir.Readdir(0)
	openedDir.Close()
	if err != nil {
		return err
	}

	for _, item := range directoryContents {
		srcSubPath := filepath.Join(src, item.Name())
		dstSubPath := filepath.Join(dst, item.Name())
		if item.IsDir() {
			fs.Mkdir(dstSubPath, item.Mode())
			err := CopyDir(fs, srcSubPath, dstSubPath)
			if err != nil {
				return err
			}
		} else {
			_, err := goutils.CopyFS(fs, srcSubPath, dstSubPath)
			if err != nil {
				fmt.Print("Received error: ")
				fmt.Println(err)
				return err
			}
		}
	}
	return nil
}
