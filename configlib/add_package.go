package configlib

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/dpastoor/goutils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/thoas/go-funk"
	"gopkg.in/yaml.v3"
	"strings"
)

// AddPackage add a package to the Package section of the yml config file
func AddPackages(pkgs []string) error {
	var pkgsToInstall []string
	ymlfile := viper.ConfigFileUsed()
	appFS := afero.NewOsFs()
	yf, err := afero.ReadFile(appFS, ymlfile)
	if err != nil {
		return err
	}
	var lines []string
	scanner := bufio.NewScanner(bytes.NewReader(yf))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	var pc PkgrConfig
	err = yaml.Unmarshal(yf, &pc)
	if err != nil {
		return err
	}

	for _, p := range pkgs {
		if funk.ContainsString(pc.Packages, p) {
			log.Debug(fmt.Sprintf("Package <%s> already found in <%s>", p, ymlfile))
			continue
		}
		pkgsToInstall = append(pkgsToInstall, p)
	}
	newYml := insertPackages(lines, pkgsToInstall...)
	err = goutils.WriteLines(newYml, ymlfile)
	return err
}


// insert packages will look for Packages: and add packages beneth
// if no Packages: line is found (eg if only tarballs/descriptions
// then one will be added. It will also preserve the indentation
// given the following are both valid:
//
// 	Packages:
// 	- dplyr
//
// 	Packages:
//    - dplyr
//
// however if both indentations are found, the yaml parser will error
// 	Packages:
// 	- dplyr
// 	  - ggplot2
func insertPackages(lines []string, pkgs ...string) []string {
	var packageStartIndex int
	var packagesSet bool
	var outputLines []string
	for i, line := range lines {
		if strings.HasPrefix(line, "Packages:")	{
			packagesSet = true
			packageStartIndex = i
		}
	}
	if packageStartIndex == 0 && !packagesSet {
		outputLines = lines
		outputLines = append(outputLines, "", "Packages:")
		for _, pkg := range pkgs {
			outputLines = append(outputLines, "- " + pkg)
		}
	} else {
		outputLines = append(outputLines, lines[:packageStartIndex+1]...)
		// given comments might precede packages, we need to find the first non-comment line
		var firstPackageIndex int
		for i, l := range lines[packageStartIndex:] {
			if strings.HasPrefix(strings.TrimSpace(l), "-") {
				firstPackageIndex = i
				break
			}
		}
		numPadding := strings.Index(lines[packageStartIndex+firstPackageIndex], "-")

		padding := strings.Repeat(" ", numPadding)
		for _, pkg := range pkgs {
			outputLines = append(outputLines, padding + "- " + pkg)
		}
		outputLines = append(outputLines, lines[packageStartIndex+1:]...)
	}

	return outputLines

}