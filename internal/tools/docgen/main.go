package main

import (
	"fmt"
	"os"

	"github.com/metrumresearchgroup/pkgr/cmd"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func gen(cmd *cobra.Command, dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0777); err != nil {
		return err
	}

	return doc.GenMarkdownTree(cmd, dir)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <directory>\n", os.Args[0])
		os.Exit(2)
	}

	outdir := os.Args[1]

	root := cmd.RootCmd
	root.Long = "" // Disable Synopsis section with version.

	root.DisableAutoGenTag = true
	for _, c := range root.Commands() {
		// Reduce update noise.
		c.DisableAutoGenTag = true
	}

	if err := gen(root, outdir); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
