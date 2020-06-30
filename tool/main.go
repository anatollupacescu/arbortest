package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/anatollupacescu/arbortest/arbor"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("» ")
}

//nolint:gochecknoglobals	//idiomatic way of working with flags in Go
var (
	path = flag.String("dir", "./", "the path to the folder containing tests")
	name = flag.String("filename", "generated_test.go", "full generated file name")
)

func main() {
	flag.Parse()

	fsDir := FsDir(*path)
	outFile := FsOutFile(*name)

	err := arbor.Generate(&fsDir, &outFile)
	for _, e := range err {
		log.Printf("error: %s", e)
	}

	if err != nil {
		log.Fatalf("\u274C%d error(s) found, aborting ", len(err))
	}
}

// FsDir represents a filesystem directory.
type FsDir string

// ListTestFiles ioutil based implementation for listing 'test' files in a directory.
func (d *FsDir) ListTestFiles() (out []arbor.File) {
	files, err := ioutil.ReadDir(string(*d))
	if err != nil {
		log.Fatalf("list test files in current directory: %v", err)
	}

	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), "_test.go") {
			fsFile := FsFile(f.Name())
			out = append(out, &fsFile)
		}
	}

	return
}

// FsFile represents a filesystem file.
type FsFile string

// ReadContents ioutil based implementation for reading contents of a file from the disk.
func (f *FsFile) ReadContents() string {
	name := string(*f)

	fileBytes, err := ioutil.ReadFile(filepath.Clean(name))
	if err != nil {
		log.Fatalf("read contents, file: %s, error: %v", name, err)
	}

	return string(fileBytes)
}

// FsOutFile represents an output file.
type FsOutFile string

// WriteContents ioutil based implementation for writing contents to disk.
func (f *FsOutFile) WriteContents(contents string) error {
	name := string(*f)
	if err := ioutil.WriteFile(name, []byte(contents), 0600); err != nil {
		return fmt.Errorf("generate test file %q: %w", name, err)
	}

	return nil
}
