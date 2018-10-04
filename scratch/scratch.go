package main

import (
	"fmt"

	"github.com/dpastoor/rpackagemanager/rcmd"
)

func main() {
	// ia := rcmd.InstallArgs{Library: "some/path"}
	ia := rcmd.NewDefaultInstallArgs()

	fmt.Println(ia.CliArgs())
	ia.Clean = true
	fmt.Println(ia.CliArgs())
	ia.Library = "path/to/lib"
	fmt.Println(ia.CliArgs())
}
