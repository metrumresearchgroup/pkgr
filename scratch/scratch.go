package main

import (
	"fmt"

	"github.com/dpastoor/rpackagemanager/rcmd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func main() {
	// ia := rcmd.InstallArgs{Library: "some/path"}
	ia := rcmd.NewDefaultInstallArgs()
	ia.Library = "../integration_tests/lib"
	appFS := afero.NewOsFs()
	lg := logrus.New()

	lg.Level = logrus.DebugLevel
	res, binaryPath, err := rcmd.InstallBinary(appFS,
		"../integration_tests/src/test1_0.0.1.tar.gz",
		ia,
		rcmd.RSettings{},
		rcmd.ExecSettings{},
		lg)

	fmt.Println("res: ", res)
	fmt.Println("binary path: ", binaryPath)
	fmt.Println("err: ", err)
}
