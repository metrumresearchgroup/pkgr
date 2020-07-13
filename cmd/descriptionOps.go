package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	//"path/filepath"
	"github.com/metrumresearchgroup/pkgr/desc"
)

func unpackDescriptions(fs afero.Fs, descPaths []string) []desc.Desc {

	var descriptions []desc.Desc

	for _, descPath := range descPaths {
		//descPath = filepath.Join(descPath, "DESCRIPTION")
		reader, err := fs.Open(descPath)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"file":  descPath,
				"error": err,
			}).Fatal("error opening DESCRIPTION file specified in Descriptions")
		}

		desc, err := desc.ParseDesc(reader)
		reader.Close()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"file":  descPath,
				"error": err,
			}).Fatal("error parsing DESCRIPTION file specified in Descriptions")
		}

		descriptions = append(descriptions, desc)
	}

	return descriptions
}


