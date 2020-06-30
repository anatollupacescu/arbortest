package arbor

import (
	"fmt"
	"log"
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
func Generate(dir Dir, out OutFile) []error {
	var testFiles = dir.ListTestFiles()
	if len(testFiles) == 0 {
		return []error{ErrNoTestFilesFound}
	}

	var callSuite, errs = parseFiles(testFiles)

	if len(errs) > 0 {
		return errs
	}

	output := generateSource("main", callSuite)
	if err := out.WriteContents(output); err != nil {
		return []error{err}
	}

	return nil
}

func parseFiles(testFiles []File) (map[string][]string, []error) {
	var s suite

	s.calls = make(map[string][]string)

	for _, fileName := range testFiles {
		fileContents := fileName.ReadContents()
		result := parseSource(fileContents)
		s.providers = append(s.providers, result.providers...)

		for k, v := range result.calls {
			s.calls[k] = v
		}
	}

	warnings := validateWarns(s)
	errs := validateErr(s.calls)

	for _, warn := range warnings {
		log.Printf("\u26a1 warning: %s", warn)
	}

	return s.calls, errs
}
