package arbor

import (
	"fmt"
)

type (
	// Dir contract for listing parseable files from a given directory.
	Dir interface {
		ListTestFiles() (out []File)
	}
	// File contract for reading test source.
	File interface {
		ReadContents() string
	}
	// OutFile contract for writing generated contents to output.
	OutFile interface {
		WriteContents(contents string) error
	}
)

// ErrNoTestFilesFound returned when there are no test files in a given folder.
var ErrNoTestFilesFound = fmt.Errorf("no test files found")

// Generate takes a folder and produces a file capable of running the tests in that location.
func Generate(dir Dir, out OutFile, pkg string) error {
	var testFiles = dir.ListTestFiles()
	if len(testFiles) == 0 {
		return ErrNoTestFilesFound
	}

	var callSuite = buildSuite(testFiles)

	printWarnings(callSuite)

	if err := validateErr(callSuite); err != nil {
		return err
	}

	output := generateSource(pkg, callSuite)
	if err := out.WriteContents(output); err != nil {
		return err
	}

	return nil
}

func buildSuite(testFiles []File) suite {
	var s = make(suite)

	for _, fileName := range testFiles {
		fileContents := fileName.ReadContents()

		result := parseSource(fileContents)
		for k, v := range result {
			s[k] = v
		}
	}

	return s
}
