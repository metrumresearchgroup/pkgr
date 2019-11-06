package cmd

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
)

//type UtilsTestSuite struct {
//	suite.Suite
//}
//
//func TestUtilsTestSuite(t *testing.T) {
//	suite.Run(t, new(UtilsTestSuite))
//}


func TestGetRestrictedWorkerCount(t *testing.T) {

	type testInput struct {
		context string
		threads int
		numCpus int
		expectedOutput int
	}

	inputs := []testInput {
		{"2 CPUs, no thread set", 0, 2, 2},
		{ "3 CPUs, no thread set", 0, 3, 2},
		{ "4 CPUs, no thread set", 0, 4, 3},
		{ "8 CPUs, no thread set", 0, 8, 7},
		{ "9 CPUs, no thread set", 0, 9, 7},
		{ "9 CPUs, no thread set", 0, 9, 7},
		{ "4 CPUs, 2 threads set", 2, 4, 2},
		{ "4 CPUs, 4 threads set", 4, 4, 4},
		{ "4 CPUs, 5 threads set", 5, 4, 5},
		{ "4 CPUs, 6 threads set", 6, 4, 6},
		{ "4 CPUs, 7 threads set", 7, 4, 7},
		{ "8 CPUs, 8 threads set", 8, 8, 8},
		{ "8 CPUs, 16 threads set", 16, 8, 16},
	}

	for _, test := range inputs {
		executeTestGetRestrictedWorkerCount(t, test.context, test.threads, test.numCpus, test.expectedOutput)
	}
}

func executeTestGetRestrictedWorkerCount(t *testing.T, context string, threads int, numCpus int, expectedOutput int) {
	t.Log(fmt.Sprintf("context: %s", context))
	actualResult := GetRestrictedWorkerCount(threads, numCpus)
	assert.Equal(t, expectedOutput, actualResult)
	t.Log("PASS")
}