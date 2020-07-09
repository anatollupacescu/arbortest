package arbor_test

import (
	"testing"

	"github.com/anatollupacescu/arbortest/arbor"
	"github.com/stretchr/testify/assert"
)

func TestNoFiles(t *testing.T) {
	var emptyDir = TestDir(func() []arbor.File {
		return nil
	})

	t.Run("errors", func(t *testing.T) {
		err := arbor.Generate(&emptyDir, nil, "test")
		assert.Equal(t, arbor.ErrNoTestFilesFound, err)
	})
}

func TestNoValidTestsDeclared(t *testing.T) {
	var src = `package sample

func testOld() error {
	return nil
}
`
	var noTestsFile = TestFile(src)
	var emptyDir = TestDir(func() []arbor.File {
		return []arbor.File{&noTestsFile}
	})
	var outFile = &TestOutFile{}

	t.Run("errors", func(t *testing.T) {
		err := arbor.Generate(&emptyDir, outFile, "ignored")
		assert.Equal(t, arbor.ErrNoTestsDeclared, err)
		assert.Equal(t, "", outFile.contents)
	})
}

func TestSingle(t *testing.T) {
	var src = `package sample

import "testing"

func testOne(t *testing.T) {
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

	t.Run("errors", func(t *testing.T) {
		err := arbor.Generate(&singleFileDir, outFile, "sample")
		assert.NoError(t, err)
		expected := `package sample

import (
	"testing"

	"github.com/anatollupacescu/arbortest/runner"
)

func TestArbor(t *testing.T) {
	dependencies := map[string][]string{
		"testOne": {"providerOne"},
	}
	tests := map[string]func(*testing.T) {
		"testOne": testOne,
	}

	output := runner.Run(t, dependencies, tests)

	if t.Failed() {
		t.Log("FAIL")
	}

	runner.Upload(output)
}
`
		assert.Equal(t, expected, outFile.contents)
	})
}

func TestCallToInvalidatedProvider(t *testing.T) {
	var src = `package random

import "testing"

func testOne(t *testing.T) {
	_ = providerOne()
	return nil
}

func providerOne() int {
	return 0
}

func testTwo(t *testing.T) {
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
		err := arbor.Generate(&singleFileDir, outFile, "random")
		err2, ok := err.(*arbor.ErrInvalidProviderCall)
		assert.True(t, ok)
		assert.EqualError(t, err2, `"testTwo" calls invalid provider: "providerTwo"`)
		assert.Equal(t, "", outFile.contents)
	})
}

func TestValidCase(t *testing.T) {
	var src = `package main

import "testing"

func testOne(t *testing.T) {
	_ = providerOne()
	return nil
}

func providerOne() int {
	return 0
}

func testTwo(t *testing.T) {
	_ = providerTwo()
	return nil
}

func providerTwo() int {
	return 0
}

func testMain(t *testing.T) {
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
		err := arbor.Generate(&singleFileDir, outFile, "main")
		assert.NoError(t, err)
		expected := `package main

import (
	"testing"

	"github.com/anatollupacescu/arbortest/runner"
)

func TestArbor(t *testing.T) {
	dependencies := map[string][]string{
		"testMain": {"providerOne", "providerTwo"}, "testOne": {"providerOne"}, "testTwo": {"providerTwo"},
	}
	tests := map[string]func(*testing.T) {
		"testMain": testMain, "testOne": testOne, "testTwo": testTwo,
	}

	output := runner.Run(t, dependencies, tests)

	if t.Failed() {
		t.Log("FAIL")
	}

	runner.Upload(output)
}
`
		assert.Equal(t, expected, outFile.contents)
	})
}

func TestValidCaseTwoFiles(t *testing.T) {
	var src1 = `package main

func testOne(t *testing.T) {
	_ = providerOne()
	return nil
}

func providerOne() int {
	return 0
}

func testTwo(t *testing.T) {
	_ = providerTwo()
	return nil
}
`

	var src2 = `package main

func providerTwo() int {
	return 0
}

func testMain(t *testing.T) {
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
		err := arbor.Generate(&singleFileDir, outFile, "main")
		assert.NoError(t, err)
		expected := `package main

import (
	"testing"

	"github.com/anatollupacescu/arbortest/runner"
)

func TestArbor(t *testing.T) {
	dependencies := map[string][]string{
		"testMain": {"providerOne", "providerTwo"}, "testOne": {"providerOne"}, "testTwo": {"providerTwo"},
	}
	tests := map[string]func(*testing.T) {
		"testMain": testMain, "testOne": testOne, "testTwo": testTwo,
	}

	output := runner.Run(t, dependencies, tests)

	if t.Failed() {
		t.Log("FAIL")
	}

	runner.Upload(output)
}
`
		assert.Equal(t, expected, outFile.contents)
	})
}

//helpers
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
