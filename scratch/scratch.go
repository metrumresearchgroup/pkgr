package main

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/metrumresearchgroup/pkgr/rcmd"
	"github.com/spf13/afero"
)

func main() {
	appFS := afero.NewOsFs()
	log := log.New()
	log.Level = log.DebugLevel
	// log.SetFormatter(&log.JSONFormatter{})
	// appFS.Remove("logfile.txt")
	// logf, _ := appFS.OpenFile("logfile.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// log.SetOutput(logf)

	// fmt.Println("library: ", viper.GetString("library"))
}
func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}

func WhatsTheCache() string {
	pc := rcmd.NewPackageCache(cmd.userCache(cfg.Cache), false)
	return ""
}
