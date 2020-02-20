package cmd

import "testing"

// This test is just to easily allow me to set breakpoints to see what's going on.
// This should not be considered part of the functional test suite, but I'm leaving
// the test here as a development tool (the mocked object is useful).
func TestDebugPrintLoadReport(t *testing.T) {
	t.Log("Quick debugging test.")
	mockRpt := loadReport{
		RMetadata: rSessionMetadata{
			LibPaths: append([]string{"path1", "path2", "path3"}),
			RVersion: "rVersionString",
			RPath:    "rPathString",
		},
		LoadResults: map[string]loadResult {
			"lr1" : loadResult {
				Stdout:  "Stdout",
				Stderr:  "Stderr",
				Success: true,
				Exiterr: nil,
				Path:    "pathpathpath",
				Package: "pkg1",
				Version: "pkg1version",
			},
			"lr2" : loadResult {
				Stdout:  "Stdout",
				Stderr:  "Stderr",
				Success: true,
				Exiterr: nil,
				Path:    "pathpathpath",
				Package: "pkg2",
				Version: "pkg2version",
			},
		},
	}

	logLoadReport(mockRpt)
}