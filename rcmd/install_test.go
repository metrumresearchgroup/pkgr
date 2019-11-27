package rcmd

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestInstallArgs(t *testing.T) {
	assert := assert.New(t)

	var installArgsTests = []struct {
		in       InstallArgs
		expected []string
	}{
		{
			InstallArgs{},
			[]string{},
		},
		{
			InstallArgs{Clean: true},
			[]string{"--clean"},
		},
		{
			InstallArgs{Library: "path/to/lib"},
			[]string{"--library=path/to/lib"},
		},
		{
			NewDefaultInstallArgs(),
			[]string{"--build", "--install-tests", "--no-multiarch", "--with-keep.source"},
		},
	}
	for i, tt := range installArgsTests {
		actual := tt.in.CliArgs()
		assert.Equal(tt.expected, actual, fmt.Sprintf("test num: %v", i+1))
	}
}

//func TestUpdateDescriptionInfoByLines(t *testing.T) {
//	var tests = []struct {
//		filename     			string
//		version      			string
//		installType  			string
//		repoURL      			string
//		repo         			string
//		expectedRepo			string
//		expectedOriginalRepo	string
//		message      			string
//	}{
//		{
//			filename:     "testsite/golden/simple/test-library/R6/Description",
//			version:      "version",
//			installType:  "binary",
//			repoURL:      "myURL",
//			repo:         "CRAN",
//			expectedRepo: "CRAN",
//			expectedOriginalRepo: "",
//			message:      "R6 test",
//		},
//		{
//			filename:     "testsite/golden/simple/test-library/pillar/Description",
//			version:      "1.2.3",
//			installType:  "binary",
//			repoURL:      "www.myURL.com",
//			repo:         "AlCran_Mandragoran",
//			expectedRepo: "AlCran_Mandragoran",
//			expectedOriginalRepo: "CRAN",
//			message:      "pillar test",
//		},
//	}
//
//	for tt := range tests {
//		result := updateDescriptionInfoByLines()
//	}
//}

//func TestUpdateDcfFile(t *testing.T) {
//	var tests = []struct {
//		filename     			string
//		version      			string
//		installType  			string
//		repoURL      			string
//		repo         			string
//		expectedRepo			string
//		expectedOriginalRepo	string
//		message      			string
//	}{
//		{
//			filename:     "testsite/golden/simple/test-library/R6/Description",
//			version:      "version",
//			installType:  "binary",
//			repoURL:      "myURL",
//			repo:         "CRAN",
//			expectedRepo: "CRAN",
//			expectedOriginalRepo: "",
//			message:      "R6 test",
//		},
//		{
//			filename:     "testsite/golden/simple/test-library/pillar/Description",
//			version:      "1.2.3",
//			installType:  "binary",
//			repoURL:      "www.myURL.com",
//			repo:         "AlCran_Mandragoran",
//			expectedRepo: "AlCran_Mandragoran",
//			expectedOriginalRepo: "CRAN",
//			message:      "pillar test",
//		},
//	}
//
//	for _, tt := range tests {
//
//		dcf, err := updateDescriptionInfo(tt.filename, tt.version, tt.installType, tt.repoURL, tt.repo)
//
//		var dcfBytes [][]byte
//		for _, s := range dcf {
//			dcfBytes = append(dcfBytes, []byte(s))
//		}
//
//		installedPackage, _ := desc.ParseDesc((dcf))
//
//		assert.Equal(t, nil, err, fmt.Sprintf("Error: %s", err))
//		assert.Equal(t, tt.expectedRepo, installedPackage.Repository, fmt.Sprintf("Failed: %s", tt.message))
//		assert.Equal(t, tt.version, installedPackage.PkgrVersion, fmt.Sprintf("Failed: %s", tt.message))
//		assert.Equal(t, tt.repoURL, installedPackage.PkgrRepositoryURL, fmt.Sprintf("Failed: %s", tt.message))
//		assert.Equal(t, tt.installType, installedPackage.PkgrInstallType, fmt.Sprintf("Failed: %s", tt.message))
//		assert.Equal(t, tt.expectedOriginalRepo, installedPackage.OriginalRepository, fmt.Sprintf("Failed: %s", tt.message))
//	}
//}

func TestUpdateDescriptionInfoByLines_RepoUpdated(t *testing.T) {
	tests := map[string]struct{
		startingLines []string
		version string
		installType string
		repoURL string
		repo string
	}{
		"Repository Upated": {
			startingLines: []string{"Package: R6", "Version: 2.4.0", "Repository: CRAN"},
			version: "pkgr0.0.test",
			installType: "binary",
			repoURL: "https://www.fakecranrepos.org",
			repo: "AlCRAN_Mandragoran",
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			results := updateDescriptionInfoByLines(test.startingLines, test.version, test.installType, test.repoURL, test.repo)
			assert.Len(t, results, 7)
			assert.Equal(t, results[0], "Package: R6")
			assert.Equal(t, results[1], "Version: 2.4.0")
			assert.Equal(t, results[2], "OriginalRepository: CRAN")
			assert.Equal(t, results[3], "Repository: " + test.repo)
			assert.Equal(t, results[4], "PkgrVersion: " + test.version)
			assert.Equal(t, results[5], "PkgrInstallType: " + test.installType)
			assert.Equal(t, results[6], "PkgrRepositoryURL: " + test.repoURL)
		})

	}


}

func TestUpdateDescriptionInfoByLines_RepoTheSame(t *testing.T) {
	tests := map[string]struct{
		startingLines []string
		version string
		installType string
		repoURL string
		repo string
	}{
		"Repository Upated": {
			startingLines: []string{"Package: R6", "Version: 2.4.0", "Repository: CRAN"},
			version: "pkgr0.0.test",
			installType: "binary",
			repoURL: "https://www.fakecranrepos.org",
			repo: "CRAN",
		},
		"Pkgr Info Updated": {
			startingLines: []string{
				"Package: R6",
				"Version: 2.4.0",
				"Repository: CRAN",
				"PkgrVersion: pkgr_older_version",
				"PkgrInstallType: source",
				"PkgrRepositoryURL: https://cran.r-project.org/",
			},
			version: "pkgr0.0.test",
			installType: "binary",
			repoURL: "https://www.fakecranrepos.org",
			repo: "CRAN",
		},
		"Pkgr Info Partially Updated": {
			startingLines: []string{
				"Package: R6",
				"Version: 2.4.0",
				"Repository: CRAN",
				"PkgrVersion: pkgr_older_version",
				"PkgrInstallType: binary", // matches final result
				"PkgrRepositoryURL: https://cran.r-project.org/",
			},
			version: "pkgr0.0.test",
			installType: "binary",
			repoURL: "https://www.fakecranrepos.org",
			repo: "CRAN",
		},
		"Pkgr Info Not Upated": {
			startingLines: []string{
				"Package: R6",
				"Version: 2.4.0",
				"Repository: CRAN",
				"PkgrVersion: pkgr0.0.test",
				"PkgrInstallType: binary", // matches final result
				"PkgrRepositoryURL: https://www.fakecranrepos.org",
			},
			version: "pkgr0.0.test",
			installType: "binary",
			repoURL: "https://www.fakecranrepos.org",
			repo: "CRAN",
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			results := updateDescriptionInfoByLines(test.startingLines, test.version, test.installType, test.repoURL, test.repo)
			assert.Len(t, results, 6)
			assert.Equal(t, results[0], "Package: R6")
			assert.Equal(t, results[1], "Version: 2.4.0")
			//assert.Equal(t, results[2], "OriginalRepository: CRAN")
			assert.Equal(t, results[2], "Repository: " + test.repo)
			assert.Equal(t, results[3], "PkgrVersion: " + test.version)
			assert.Equal(t, results[4], "PkgrInstallType: " + test.installType)
			assert.Equal(t, results[5], "PkgrRepositoryURL: " + test.repoURL)
		})

	}
}