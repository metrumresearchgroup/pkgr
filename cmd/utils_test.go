package cmd

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetRestrictedWorkerCount(t *testing.T) {

	type testInput struct {
		context string
		threads int
		numCpus int
		expectedOutput int
	}

	inputs := []testInput {
		{"1 CPU, no thread set", 0, 1, 1},
		{"2 CPUs, no thread set", 0, 2, 2},
		{ "3 CPUs, no thread set", 0, 3, 3},
		{ "4 CPUs, no thread set", 0, 4, 4},
		{ "8 CPUs, no thread set", 0, 8, 8},
		{ "9 CPUs, no thread set", 0, 9, 8},
		{ "16 CPUs, no thread set", 0, 16, 8},
		{ "4 CPUs, 2 threads set", 2, 4, 2},
		{ "4 CPUs, 4 threads set", 4, 4, 4},
		{ "4 CPUs, 5 threads set", 5, 4, 5},
		{ "4 CPUs, 6 threads set", 6, 4, 6},
		{ "4 CPUs, 7 threads set", 7, 4, 7},
		{ "8 CPUs, 8 threads set", 8, 8, 8},
		{ "8 CPUs, 9 threads set", 9, 8, 9},
		{ "8 CPUs, 10 threads set", 10, 8, 10},
		{ "8 CPUs, 16 threads set", 16, 8, 16},
	}

	for _, test := range inputs {
		executeTestGetRestrictedWorkerCount(t, test.context, test.threads, test.numCpus, test.expectedOutput)
	}
}

func executeTestGetRestrictedWorkerCount(t *testing.T, context string, threads int, numCpus int, expectedOutput int) {
	t.Log(fmt.Sprintf("context: %s", context))
	actualResult := getWorkerCount(threads, numCpus)
	assert.Equal(t, expectedOutput, actualResult, "context: " + context)
}