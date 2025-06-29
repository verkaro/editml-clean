// main_test.go
package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/verkaro/editml-go"
)

// TestRun is a table-driven test for the core `run` function.
// It verifies that the supported EditML inputs are correctly transformed.
func TestRun(t *testing.T) {
	// (Unit tests remain the same as the last working version)
	testCases := []struct {
		name           string
		input          string
		expectedText   string
		expectedIssues []editml.Issue
	}{
		{ name: "Simple Addition", input: "Hello{+ world+}.", expectedText: "Hello world." },
		{ name: "Simple Deletion", input: "This is{- not-} good.", expectedText: "This is good." },
		{ name: "Comment Removal", input: "A key point.{>Remember to check this later.<}", expectedText: "A key point." },
		{ name: "Highlight Removal", input: "The answer is {=42=}.", expectedText: "The answer is 42." },
		{ name: "Combination of Edits", input: "Please {-review-}{+read+} this {=document=} carefully.", expectedText: "Please read this document carefully." },
		{ name: "Official Move Syntax", input: "Let's put B here: {move:word}. And here is {move~A~word}.", expectedText: "Let's put B here: A. And here is ."},
		{ name: "Official Copy Syntax", input: "Here is the original: {copy~A~word}. And here is a copy: {copy:word}.", expectedText: "Here is the original: A. And here is a copy: A."},
		{ name: "Shorthand Move Syntax (mv)", input: "Third.{mv~First~t1}Second.{mv:t1}", expectedText: "Third.Second.First" },
		{ name: "Shorthand Copy Syntax (cp)", input: "{cp~A~t1}BC {cp:t1}", expectedText: "ABC A" },
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualText, actualIssues := run(tc.input)
			if actualText != tc.expectedText {
				t.Errorf("Mismatched text:\nexpected: %q\nactual:   %q", tc.expectedText, actualText)
			}
			if len(actualIssues) == 0 && len(tc.expectedIssues) == 0 {
			} else if !reflect.DeepEqual(actualIssues, tc.expectedIssues) {
				t.Errorf("Mismatched issues:\nexpected: %+v\nactual:   %+v", tc.expectedIssues, actualIssues)
			}
		})
	}
}


// --- Integration Tests ---

// TestIntegration runs end-to-end tests on the compiled binary by reading
// test cases from the './testdata' directory.
func TestIntegration(t *testing.T) {
	// Build the binary once for all integration tests.
	tempDir := t.TempDir()
	binaryPath := filepath.Join(tempDir, "editml-clean")
	if err := exec.Command("go", "build", "-o", binaryPath, ".").Run(); err != nil {
		t.Fatalf("Failed to build binary for integration test: %v", err)
	}

	// Walk the testdata directory to find all .editml files.
	testdataDir := "testdata"
	filepath.Walk(testdataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".editml") {
			// Found a test case, run it.
			testName := strings.TrimSuffix(info.Name(), ".editml")
			t.Run(testName, func(t *testing.T) {
				inputFile := path
				goldenFile := strings.Replace(path, ".editml", ".golden.txt", 1)
				outputFile := filepath.Join(tempDir, testName+".output.txt")

				// All flags must come BEFORE positional arguments.
				cmd := exec.Command(binaryPath, "--output", outputFile, inputFile)
				
				var stderr strings.Builder
				cmd.Stderr = &stderr
				
				err := cmd.Run()
				if err != nil {
					t.Fatalf("Command execution failed: %v\nstderr: %s", err, stderr.String())
				}

				actualBytes, err := os.ReadFile(outputFile)
				if err != nil {
					t.Fatalf("Failed to read actual output file '%s': %v", outputFile, err)
				}
				
				expectedBytes, err := os.ReadFile(goldenFile)
				if err != nil {
					t.Fatalf("Failed to read golden file '%s': %v", goldenFile, err)
				}

				// REVERTED: Perform a strict, byte-for-byte comparison.
				// The application is now responsible for adding the trailing newline.
				if string(actualBytes) != string(expectedBytes) {
					t.Errorf("Mismatched output:\nexpected: %q\nactual:   %q", string(expectedBytes), string(actualBytes))
				}
			})
		}
		return nil
	})
}

