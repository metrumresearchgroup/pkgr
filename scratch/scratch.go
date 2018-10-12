package main

import (
	"encoding/json"
	"fmt"

	"github.com/dpastoor/rpackagemanager/desc"
	"github.com/dpastoor/rpackagemanager/rpkg"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func main() {
	// ia := rcmd.InstallArgs{Library: "some/path"}
	appFS := afero.NewOsFs()
	lg := logrus.New()
	lg.Level = logrus.DebugLevel

	// b28ba6e911e86ae4e682f834741e85e0
	hash, _ := rpkg.Hash(appFS, "../integration_tests/src/test1_0.0.1.tar.gz")
	fmt.Println(hash)

	d, err := desc.ReadDesc("../desc/testdata/D2")
	if err != nil {
		panic(err)
	}
	PrettyPrint(d)
}

func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}
