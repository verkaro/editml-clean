package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/verkaro/editml-go"
)

// AppVersion holds the version of the utility.
const AppVersion = "1.0.0"

// run contains the core processing logic for the application.
// It takes the raw input text, processes it through the EditML library,
// and returns the clean text and any issues encountered.
func run(inputText string) (cleanText string, allIssues []editml.Issue) {
	nodes, parseIssues := editml.Parse(inputText)
	// Even if parsing has issues, we attempt to transform to see if we get more issues or partial output.
	cleanText, transformIssues := editml.TransformCleanView(nodes)
	allIssues = append(parseIssues, transformIssues...)
	return cleanText, allIssues
}

// main handles command-line parsing, I/O, and coordinating the call to the
// core run function. It is responsible for exit codes and printing issues.
func main() {
	// --- 1. Flag Definition & Parsing ---
	versionFlag := flag.Bool("version", false, "Print version and exit.")
	outputFlag := flag.String("o", "", "Write output to the specified file instead of stdout. (shorthand)")
	outputFlagLong := flag.String("output", "", "Write output to the specified file instead of stdout.")
	debugFlag := flag.Bool("debug", false, "Emit parse/transform issues (warnings/errors) to stderr.")
	strictFlag := flag.Bool("strict", false, "Treat warnings as errors (exit non-zero on any issue).")

	flag.Parse()

	if *versionFlag {
		fmt.Println("editml-clean version", AppVersion)
		os.Exit(0)
	}

	// --- 2. Input Handling ---
	var inputReader io.Reader = os.Stdin
	var err error
	if flag.NArg() > 0 {
		inputFilename := flag.Arg(0)
		file, err := os.Open(inputFilename)
		if err != nil {
			log.Fatalf("Fatal: could not open input file %s: %v", inputFilename, err)
		}
		defer file.Close()
		inputReader = file
	}

	inputBytes, err := io.ReadAll(inputReader)
	if err != nil {
		log.Fatalf("Fatal: could not read input: %v", err)
	}
	inputText := string(inputBytes)

	// --- 3. Core Processing ---
	cleanText, allIssues := run(inputText)

	// --- 4. Issue and Error Handling ---
	hasError := false
	hasWarning := false
	if len(allIssues) > 0 {
		for _, issue := range allIssues {
			if issue.Severity == editml.SeverityError {
				hasError = true
			} else {
				hasWarning = true
			}
			if *debugFlag {
				severityStr := "Warning"
				if issue.Severity == editml.SeverityError {
					severityStr = "Error"
				}
				fmt.Fprintf(os.Stderr, "[%s] %s (L%d:C%d)\n", severityStr, issue.Message, issue.Line, issue.Column)
			}
		}
	}

	// --- 5. Exit Code and Output Logic ---
	if hasError {
		os.Exit(1)
	}
	if hasWarning && *strictFlag {
		os.Exit(2)
	}

	var outputWriter io.Writer = os.Stdout
	outputPath := *outputFlag
	if *outputFlagLong != "" {
		outputPath = *outputFlagLong
	}
	
	if outputPath != "" {
		file, err := os.Create(outputPath)
		if err != nil {
			log.Fatalf("Fatal: could not create output file %s: %v", outputPath, err)
		}
		defer file.Close()
		outputWriter = file
	}

	// FIX: Use fmt.Fprintln to ensure the output always ends with a newline,
	// which is standard convention for command-line tools.
	_, err = fmt.Fprintln(outputWriter, cleanText)
	if err != nil {
		log.Fatalf("Fatal: could not write to output: %v", err)
	}
}

