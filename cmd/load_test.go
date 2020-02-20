package cmd

import (
	"regexp"
	"testing"
)

//// This test is just to easily allow me to set breakpoints to see what's going on.
//// This should not be considered part of the functional test suite, but I'm leaving
//// the test here as a development tool (the mocked object is useful).
//func TestDebugPrintLoadReport(t *testing.T) {
//	t.Log("Quick debugging test.")
//	mockRpt := loadReport{
//		RMetadata: rSessionMetadata{
//			LibPaths: append([]string{"path1", "path2", "path3"}),
//			RVersion: "rVersionString",
//			RPath:    "rPathString",
//		},
//		LoadResults: map[string]loadResult {
//			"lr1" : loadResult {
//				Stdout:  "Stdout",
//				Stderr:  "Stderr",
//				Success: true,
//				Exiterr: "",
//				Path:    "pathpathpath",
//				Package: "pkg1",
//				Version: "pkg1version",
//			},
//			"lr2" : loadResult {
//				Stdout:  "Stdout",
//				Stderr:  "Stderr",
//				Success: true,
//				Exiterr: "",
//				Path:    "pathpathpath",
//				Package: "pkg2",
//				Version: "pkg2version",
//			},
//		},
//	}
//
//	logLoadReport(mockRpt)
//}
//
//func TestStringSlice(t *testing.T) {
//	sss := "/Users/johncarlos/go/src/github.com/metrumresearchgroup/pkgr/integration_tests/load-fail/test-library"
//	t.Log(sss[5:])
//	t.Log(sss[0])
//	t.Log(string(sss[0]))
//
//
//	//// match regexp as in question
//	//	//pat := regexp.MustCompile(`https?://.*\.txt`)
//	//	//s := pat.FindString(myString)
//	ss := "[12323] \"/Users/johncarlos/go/src/github.com/metrumresearchgroup/pkgr/integration_tests/load-fail/test-library\""
//	ss2 := ">"
//	pattern := regexp.MustCompile(`\[\d*\] \"(.*)\"`)
//	submatch := pattern.FindStringSubmatch(ss)[1]
//	t.Log(submatch)
//	submatch2 := pattern.FindStringSubmatch(ss2)
//	t.Log(len(submatch2))
//}