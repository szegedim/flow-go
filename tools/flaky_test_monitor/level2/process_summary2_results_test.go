package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/onflow/flow-go/tools/flaky_test_monitor/common"
	"github.com/onflow/flow-go/tools/flaky_test_monitor/common/testdata"
)

func TestProcessSummary2TestRun(t *testing.T) {
	testDataMap := map[string]testdata.TestDataLevel2{
		"1 level 1 summary, 1 failure the rest pass": {
			Directory:        "test1-1package-1failure",
			HasFailures:      true,
			HasNoResultTests: false,
			TestRuns:         testdata.GetTestData_Level2_1FailureRestPass(),
		},

		"1 level 1 summary, 1 no-result test, no other tests": {
			Directory:        "test2-1-no-result-test",
			HasFailures:      false,
			HasNoResultTests: true,
			TestRuns:         testdata.GetTestsData_Level2_1NoResultNoOtherTests(),
		},
		// "many level 1 summaries, many no-result tests":                             {"test3-multi-no-result-tests", false, true},
		// "many level 1 summaries, many failures, many passes":                       {"test4-multi-failures", true, false},
		// "many level 1 summaries, many failures, many passes, many no-result tests": {"test5-multi-failures-multi-no-result-tests", true, true},
	}

	for k, testData := range testDataMap {
		t.Run(k, func(t *testing.T) {
			setUp(t)
			runProcessSummary2TestRun(t, testData)
			tearDown(t)
		})
	}
}

func setUp(t *testing.T) {
	deleteMessagesDir(t)
}

func tearDown(t *testing.T) {
	deleteMessagesDir(t)
}

// HELPERS - UTILITIES

const actualFailureMessagesPath = "./failures"
const actualNoResultMessagesPath = "./no-results"

func deleteMessagesDir(t *testing.T) {
	// delete failure test dir that stores failure messages
	err := os.RemoveAll(actualFailureMessagesPath)
	require.Nil(t, err)

	// delete no-result test dir that stores no-result messages
	err = os.RemoveAll(actualNoResultMessagesPath)
	require.Nil(t, err)
}

func runProcessSummary2TestRun(t *testing.T, testData testdata.TestDataLevel2) {

	inputTestDataPath := filepath.Join("../testdata/summary2", testData.Directory, "input")

	expectedOutputTestDataPath := filepath.Join("../testdata/summary2", testData.Directory, "expected-output", testData.Directory+".json")
	expectedFailureMessagesPath := filepath.Join("../testdata/summary2", testData.Directory, "expected-output/failures")
	expectedNoResultMessagesPath := filepath.Join("../testdata/summary2", testData.Directory, "expected-output/no-results")

	// **************************************************************
	// can run the test from expected structs of level 1 TestRuns or from generated level 1 JSON files
	var actualTestLevel2Summary common.TestsLevel2Summary
	if len(testData.TestRuns) > 0 {
		actualTestLevel2Summary = processSummary2TestRunFromStructs(testData.TestRuns)
	} else {
		actualTestLevel2Summary = processSummary2TestRun(inputTestDataPath)
	}
	// **************************************************************

	// read in expected summary level 2
	var expectedTestsLevel2Summary common.TestsLevel2Summary
	expectedTestSummary2JsonBytes, err := os.ReadFile(expectedOutputTestDataPath)
	require.Nil(t, err)
	require.NotEmpty(t, expectedTestSummary2JsonBytes)
	err = json.Unmarshal(expectedTestSummary2JsonBytes, &expectedTestsLevel2Summary)
	require.Nil(t, err)

	//check all details of test summary level 2 between expected and actual
	require.Equal(t, len(expectedTestsLevel2Summary.TestResults), len(actualTestLevel2Summary.TestResults))

	for expectedTestResultSummaryKey := range expectedTestsLevel2Summary.TestResults {
		expectedTestSummary, isFoundExpected := expectedTestsLevel2Summary.TestResults[expectedTestResultSummaryKey]
		actualTestSummary, isFoundActual := actualTestLevel2Summary.TestResults[expectedTestResultSummaryKey]

		require.True(t, isFoundExpected)
		require.True(t, isFoundActual, "not found: "+expectedTestResultSummaryKey)

		common.AssertTestSummariesEqual(t, *expectedTestSummary, *actualTestSummary)
	}

	// make sure calculated summary level 2 is what we expected
	require.Equal(t, expectedTestsLevel2Summary, actualTestLevel2Summary)

	checkFailureMessages(t, testData.HasFailures, expectedFailureMessagesPath)
	checkNoResultMessages(t, testData.HasNoResultTests, expectedNoResultMessagesPath)
}

// check failure messages created
// there are 2 types of scenarios:
// 1. test summaries with no failures - these will not have a `failures` sub-directory and no more checking is needed
// 2. test summaries with failures - these will have a `failures` sub-directory with failure messages saved
//    in text files (1 file/failure under separate sub-directory for each test that has failures)
func checkFailureMessages(t *testing.T, hasFailures bool, expectedFailureMessagesPath string) {
	if !hasFailures {
		return
	}
	checkMessagesHelper(t, expectedFailureMessagesPath, actualFailureMessagesPath)
}

// check no-result messages created - for tests that generated no pass / fail
// there are 2 types of scenarios:
// 1. test summaries with no "no-result" - these will not have a `no-result` sub-directory and no more checking is needed
// 2. test summaries with no-results - these will have a `no-result` sub-directory with output messages saved
//    in text files (1 file/no-result under separate sub-directory for each test that has no-results)
func checkNoResultMessages(t *testing.T, hasNoResultTests bool, expectedNoResultMessagesPath string) {
	if !hasNoResultTests {
		return
	}
	checkMessagesHelper(t, expectedNoResultMessagesPath, actualNoResultMessagesPath)
}

// helps check for both failures and no-result messages since they are very similar, just under
// different directories
func checkMessagesHelper(t *testing.T, expectedMessagesPath string, actualMessagesPath string) {

	// count expected failure / no-result directories (1 directory/test)
	expectedMessageDirs, err := os.ReadDir(expectedMessagesPath)
	require.Nil(t, err)

	// count actual failure / no-result directories
	actualMessageDirs, err := os.ReadDir(actualMessagesPath)
	require.Nil(t, err)

	// expected test summary has at least 1 failure / no-result
	require.Equal(t, len(expectedMessageDirs), len(actualMessageDirs))

	// compare expected vs actual messages
	for expectedMessageDirIndex, expectedMessageDir := range expectedMessageDirs {

		// sub-directory names should be the same - each sub directory corresponds to a failed / no-result test name
		require.Equal(t, expectedMessageDir.Name(), actualMessageDirs[expectedMessageDirIndex].Name())

		// under each sub-directory, there should be 1 or more text files (failure1.txt/no-result1.txt, failure2.txt/no-result2.txt, etc)
		// that holds the raw failure / no-result message for that test
		expectedMessagesDirFiles, err := os.ReadDir(filepath.Join(expectedMessagesPath, expectedMessageDir.Name()))
		require.Nil(t, err)

		actualMessageDirFiles, err := os.ReadDir(filepath.Join(actualMessagesPath, actualMessageDirs[expectedMessageDirIndex].Name()))
		require.Nil(t, err)

		// make sure there are the expected number of failed / no-result text files in the sub-directory
		require.Equal(t, len(expectedMessagesDirFiles), len(actualMessageDirFiles))

		// check contents of each text file for expected failure / no-result message
		// for every test that has failures / no-result, there should be 1 text file per failure / no-result

		// if test has failures / no-results, there should be directory of failure / no-result messages text files
		// a sub-directory of the test name will hold all test failure / no-result messages

		for expectedMessageFileIndex, expectedMessageFileDirEntry := range expectedMessagesDirFiles {
			expectedMessageFilePath := filepath.Join(expectedMessagesPath, expectedMessageDir.Name(), expectedMessageFileDirEntry.Name())
			expectedMessageFileBytes, err := os.ReadFile(expectedMessageFilePath)
			require.Nil(t, err)

			actualMessageFilePath := filepath.Join(actualMessagesPath, actualMessageDirs[expectedMessageDirIndex].Name(), actualMessageDirFiles[expectedMessageFileIndex].Name())
			actualMessageFileBytes, err := os.ReadFile(actualMessageFilePath)
			require.Nil(t, err)

			// read expected and actual text files as bytes and compare them all at once
			require.Equal(t, expectedMessageFileBytes, actualMessageFileBytes)
		}
	}
}
