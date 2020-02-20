package cmd

import "testing"

// This test is just to easily allow me to set breakpoints to see what's going on.
// This should not be considered part of the functional test suite, but I'm leaving
// the test here as a development tool (the mocked object is useful).
func TestDebugPrintLoadReport(t *testing.T) {
	t.Log("Quick debugging test.")
	mockRpt := loadReport{
		rMetadata: rSessionMetadata{
			libPaths: append([]string{"path1", "path2", "path3"}),
			rVersion: "rVersionString",
			rPath: "rPathString",
		},
		results : map[string]loadResult {
			"lr1" : loadResult {
				stdout: "stdout",
				stderr: "stderr",
				success: true,
				exiterr: nil,
				path: "pathpathpath",
				pkg: "pkg1",
				version: "pkg1version",
			},
			"lr2" : loadResult {
				stdout: "stdout",
				stderr: "stderr",
				success: true,
				exiterr: nil,
				path: "pathpathpath",
				pkg: "pkg2",
				version: "pkg2version",
			},
		},
	}

	logLoadReport(mockRpt)
}