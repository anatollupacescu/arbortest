package arbor

import (
	"fmt"
)

type (
	// File contract for reading test source.
	File interface {
		Read() string
	}
	// Dir contract for listing parseable files from a given directory.
	Dir interface {
		List() []File
	}
	// OutFile contract for writing generated contents to output.
	OutFile interface {
		Write(contents string) error
	}
)

var (
	// ErrNoTestFilesFound returned when there are no test files in a given folder.
	ErrNoTestFilesFound = fmt.Errorf("no test files found")
	// ErrNoTestsDeclared returned when no testing code is found in the given test files.
	ErrNoTestsDeclared = fmt.Errorf("no tests found declared in the given folder files")
)

// Generate takes a folder and produces a test file for that package.
func Generate(dir Dir, out OutFile, pkg string) error {
	var testFiles = dir.List()
	if len(testFiles) == 0 {
		return ErrNoTestFilesFound
	}

	graph, err := buildGraph(testFiles)
	if err != nil {
		return err
	}

	if graph.isEmpty() {
		return ErrNoTestsDeclared
	}

	output := generateSource(pkg, graph)
	if err := out.Write(output); err != nil {
		return err
	}

	return nil
}

func buildGraph(testFiles []File) (graph, error) {
	bundles := make([]testBundle, 0)

	for _, fileName := range testFiles {
		source := fileName.Read()
		newBundles := parse(source)
		bundles = append(bundles, newBundles...)
	}

	built, err := build(bundles)

	return built, err
}
