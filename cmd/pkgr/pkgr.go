package main

import (
	"fmt"
	"os"

	"go.uber.org/automaxprocs/maxprocs"

	"github.com/metrumresearchgroup/pkgr/cmd"
)

func main() {
	setGOMAXPROCS()
	cmd.Execute()
}

func setGOMAXPROCS() {
	// Silently set the maxprocs
	discard := func(string, ...interface{}) {}
	if _, err := maxprocs.Set(maxprocs.Logger(discard)); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "unable to set maxprocs: %v", err)
	}
}
