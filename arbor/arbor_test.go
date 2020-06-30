package arbor_test

import (
	"sort"
	"testing"

	"github.com/anatollupacescu/arbortest/arbor"
	"github.com/stretchr/testify/assert"
)

type TestDir func() []arbor.File

func (t *TestDir) ListTestFiles() (out []arbor.File) {
	return (*t)()
}

type TestFile string

func (t *TestFile) ReadContents() string {
	return string(*t)
}

type TestOutFile struct {
	contents string
	err      error
}

func (t *TestOutFile) WriteContents(contents string) error {
	t.contents = contents
	return t.err
}

func TestArbor(t *testing.T) {
	t.Run("given a folder with no test files", func(t *testing.T) {
		var emptyDir = TestDir(func() []arbor.File {
			return nil
		})

		t.Run("errors", func(t *testing.T) {
			errors := arbor.Generate(&emptyDir, nil, "test")
			if assert.Len(t, errors, 1) {
				expected := arbor.ErrNoTestFilesFound
				assert.Equal(t, expected, errors[0])
			}
		})
	})
	t.Run("given a file without tests", func(t *testing.T) {
		var noTestsFile = TestFile("package test")
		var emptyDir = TestDir(func() []arbor.File {
			return []arbor.File{&noTestsFile}
		})
		var outFile = &TestOutFile{}

		t.Run("errors", func(t *testing.T) {
			errors := arbor.Generate(&emptyDir, outFile, "ignored")
			if assert.Len(t, errors, 1) {
				expected := arbor.ErrNoTestsDeclared
				assert.Equal(t, expected, errors[0])
				assert.Equal(t, "", outFile.contents)
			}
		})
	})
	t.Run("given a file with one test", func(t *testing.T) {
		var src = `package testing

func testOne() error {
	return nil
}`

		var noTestsFile = TestFile(src)

		var emptyDir = TestDir(func() []arbor.File {
			return []arbor.File{&noTestsFile}
		})
		var outFile = &TestOutFile{}

		t.Run("registers one test", func(t *testing.T) {
			errors := arbor.Generate(&emptyDir, outFile, "testing")
			assert.Len(t, errors, 0)

			expected := `package testing

import (
	"testing"

	"github.com/anatollupacescu/arbortest/runner"
)

func TestArbor(t *testing.T) {
	var validators map[string]string
	dependencies := map[string][]string{
		"testOne": {},
	}
	tests := map[string]func() error{
		"testOne": testOne,
	}
	if r := runner.Run(validators, dependencies, tests); r.Error != "" {
		t.Error(r.Error)
	}
}
`
			assert.Equal(t, expected, outFile.contents)
		})
	})

	t.Run("given a test and a related provider", func(t *testing.T) {
		var src = `package sample

func testOne() error {
	_ = providerOne()
	return nil
}

func providerOne() int {
	return 0
}`

		var (
			testProviderFile = TestFile(src)
			singleFileDir    = TestDir(func() []arbor.File {
				return []arbor.File{&testProviderFile}
			})
			outFile = &TestOutFile{}
		)

		t.Run("should register test", func(t *testing.T) {
			errors := arbor.Generate(&singleFileDir, outFile, "sample")
			assert.Len(t, errors, 0)

			expected := `package sample

import (
	"testing"

	"github.com/anatollupacescu/arbortest/runner"
)

func TestArbor(t *testing.T) {
	var validators map[string]string
	var dependencies map[string][]string
	tests := map[string]func() error{
		"testOne": testOne,
	}
	if r := runner.Run(validators, dependencies, tests); r.Error != "" {
		t.Error(r.Error)
	}
}
`
			assert.Equal(t, expected, outFile.contents)
		})
	})

	t.Run("given a test and two invalidated providers", func(t *testing.T) {
		var src = `package random

func testOne() error {
	_ = providerOne()
	return nil
}

func providerOne() int {
	return 0
}

func testTwo() error {
	_ = providerOne()
	_ = providerTwo()
	_ = providerThree()
	return nil
}

func providerTwo() int {
	return 0
}

func providerThree() int {
	return 0
}`

		var (
			testProviderFile = TestFile(src)
			singleFileDir    = TestDir(func() []arbor.File {
				return []arbor.File{&testProviderFile}
			})
			outFile = &TestOutFile{}
		)

		t.Run("returns two errors", func(t *testing.T) {
			errors := arbor.Generate(&singleFileDir, outFile, "random")
			if assert.Len(t, errors, 2) {
				sort.Slice(errors, func(i, j int) bool {
					return errors[i].Error() < errors[j].Error()
				})
				err1, ok1 := errors[0].(*arbor.ErrInvalidProviderCall)
				err2, ok2 := errors[1].(*arbor.ErrInvalidProviderCall)
				assert.True(t, ok1 && ok2, "errors are of the correct type")

				assert.EqualError(t, err1, `"testTwo" calls invalid provider: "providerThree"`)
				assert.EqualError(t, err2, `"testTwo" calls invalid provider: "providerTwo"`)
			}
			assert.Equal(t, "", outFile.contents)
		})
	})
	t.Run("given a valid test", func(t *testing.T) {
		var src = `package main

func testOne() error {
	_ = providerOne()
	return nil
}

func providerOne() int {
	return 0
}

func testTwo() error {
	_ = providerTwo()
	return nil
}

func providerTwo() int {
	return 0
}

func testMain() error {
	_ = providerOne()
	_ = providerTwo()
	return nil
}
`
		var (
			testProviderFile = TestFile(src)
			singleFileDir    = TestDir(func() []arbor.File {
				return []arbor.File{&testProviderFile}
			})
			outFile = &TestOutFile{}
		)

		t.Run("correct file is generated", func(t *testing.T) {
			errors := arbor.Generate(&singleFileDir, outFile, "main")
			assert.Len(t, errors, 0)
			expected := `package main

import (
	"testing"

	"github.com/anatollupacescu/arbortest/runner"
)

func TestArbor(t *testing.T) {
	validators := map[string]string{
		"testOne": "providerOne", "testTwo": "providerTwo",
	}
	dependencies := map[string][]string{
		"testMain": {"providerOne", "providerTwo"},
	}
	tests := map[string]func() error{
		"testMain": testMain, "testOne": testOne, "testTwo": testTwo,
	}
	if r := runner.Run(validators, dependencies, tests); r.Error != "" {
		t.Error(r.Error)
	}
}
`
			assert.Equal(t, expected, outFile.contents)
		})
	})
	t.Run("given two valid test files", func(t *testing.T) {
		var src1 = `package main

func testOne() error {
	_ = providerOne()
	return nil
}

func providerOne() int {
	return 0
}

func testTwo() error {
	_ = providerTwo()
	return nil
}
`

		var src2 = `package main

func providerTwo() int {
	return 0
}

func testMain() error {
	_ = providerOne()
	_ = providerTwo()
	return nil
}
`
		var (
			testProviderFile1 = TestFile(src1)
			testProviderFile2 = TestFile(src2)
			singleFileDir     = TestDir(func() []arbor.File {
				return []arbor.File{&testProviderFile1, &testProviderFile2}
			})
			outFile = &TestOutFile{}
		)

		t.Run("it parses both as one", func(t *testing.T) {
			errors := arbor.Generate(&singleFileDir, outFile, "main")
			assert.Len(t, errors, 0)
			expected := `package main

import (
	"testing"

	"github.com/anatollupacescu/arbortest/runner"
)

func TestArbor(t *testing.T) {
	validators := map[string]string{
		"testOne": "providerOne", "testTwo": "providerTwo",
	}
	dependencies := map[string][]string{
		"testMain": {"providerOne", "providerTwo"},
	}
	tests := map[string]func() error{
		"testMain": testMain, "testOne": testOne, "testTwo": testTwo,
	}
	if r := runner.Run(validators, dependencies, tests); r.Error != "" {
		t.Error(r.Error)
	}
}
`
			assert.Equal(t, expected, outFile.contents)
		})
	})
}
