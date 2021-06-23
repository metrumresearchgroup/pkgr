package main

import (
	// TODO: unless there's a reason not to, move these features into a cmds sub-directory"
	"github.com/metrumresearchgroup/pkgr/cmd"
)

// buildTime  can be set from LDFLAGS during development
var buildTime string

// TODO: remove comment
// if want to generate docs
//	import "github.com/spf13/cobra/doc"
//	err := doc.GenMarkdownTree(cmd.RootCmd, "../../docs/bbi")
//	if err != nil {
//		panic(err)
//	}
func main() {
	cmd.Execute(buildTime)
}
