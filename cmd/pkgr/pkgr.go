package main

import (
	"github.com/metrumresearchgroup/pkgr/cmd"
)

// if want to generate docs
//	import "github.com/spf13/cobra/doc"
//	err := doc.GenMarkdownTree(cmd.RootCmd, "../../docs/bbi")
//	if err != nil {
//		panic(err)
//	}
func main() {
	cmd.Execute()
}
