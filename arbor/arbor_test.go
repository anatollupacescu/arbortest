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

func TestInvalidSignature(t *testing.T) {
	t.Run("test name not following convention", func(t *testing.T) {
		var src = `package sample

func iDontMatter() error {
	return nil
}`

		var singleFile = TestFile(src)
		var singleFileDir = TestDir(func() []arbor.File {
			return []arbor.File{&singleFile}
		})

		t.Run("errors", func(t *testing.T) {
			err := arbor.Generate(&singleFileDir, nil, "test")
			assert.Equal(t, arbor.ErrNoTestsDeclared, err)
		})
	})

	t.Run("bad test parameters", func(t *testing.T) {

		var src = `package sample

// group:a
func testOld() error {
	return nil
}
`
		var wrongSignatureFile = TestFile(src)
		var singleFileDir = TestDir(func() []arbor.File {
			return []arbor.File{&wrongSignatureFile}
		})
		var outFile = &TestOutFile{}

		t.Run("errors", func(t *testing.T) {
			err := arbor.Generate(&singleFileDir, outFile, "ignored")
			assert.Equal(t, arbor.ErrNoTestsDeclared, err)
			assert.Equal(t, "", outFile.contents)
		})
	})
}

func TestSingleTestSingleGroup(t *testing.T) {
	var src = `package sample

import "github.com/anatollupacescu/arbortest/runner"

// group:one
func testOne(t *runner.T) {}`

	var (
		testProviderFile = TestFile(src)
		singleFileDir    = TestDir(func() []arbor.File {
			return []arbor.File{&testProviderFile}
		})
		outFile = &TestOutFile{}
	)

	t.Run("registers one group", func(t *testing.T) {
		err := arbor.Generate(&singleFileDir, outFile, "sample")
		assert.NoError(t, err)
		expected := `package sample

import (
	"testing"

	arbor "github.com/anatollupacescu/arbortest/runner"
)

func TestArbor(t *testing.T) {
	g := arbor.New()

	t.Run("one", func(t *testing.T) {
		at := arbor.NewT(t)
		g.Group("one")
		g.Append(at, "testOne", testOne)
	})

	output := g.JSON()

	arbor.Upload(output)
}
`
		assert.Equal(t, expected, outFile.contents)
	})
}

func TestTwoTestsSameGroup(t *testing.T) {
	var src = `package random

import (
	"testing"

	"github.com/anatollupacescu/arbortest/runner"
)

// group:one
func testOne(t *runner.T) {}

// group:one
func testTwo(t *runner.T) {}
`

	var (
		testProviderFile = TestFile(src)
		singleFileDir    = TestDir(func() []arbor.File {
			return []arbor.File{&testProviderFile}
		})
		outFile = &TestOutFile{}
	)

	t.Run("registers one group", func(t *testing.T) {
		err := arbor.Generate(&singleFileDir, outFile, "sample")
		assert.NoError(t, err)
		expected := `package sample

import (
	"testing"

	arbor "github.com/anatollupacescu/arbortest/runner"
)

func TestArbor(t *testing.T) {
	g := arbor.New()

	t.Run("one", func(t *testing.T) {
		at := arbor.NewT(t)
		g.Group("one")
		g.Append(at, "testOne", testOne)
		g.Append(at, "testTwo", testTwo)
	})

	output := g.JSON()

	arbor.Upload(output)
}
`
		assert.Equal(t, expected, outFile.contents)
	})
}

func TestTwoTestsTwoGroups(t *testing.T) {
	var src = `package random

import (
	"testing"

	"github.com/anatollupacescu/arbortest/runner"
)

// group:one
func testOne(t *runner.T) {}

// group:two
func testTwo(t *runner.T) {}
`

	var (
		testProviderFile = TestFile(src)
		singleFileDir    = TestDir(func() []arbor.File {
			return []arbor.File{&testProviderFile}
		})
		outFile = &TestOutFile{}
	)

	t.Run("registers one group", func(t *testing.T) {
		err := arbor.Generate(&singleFileDir, outFile, "sample")
		assert.NoError(t, err)
		expected := `package sample

import (
	"testing"

	arbor "github.com/anatollupacescu/arbortest/runner"
)

func TestArbor(t *testing.T) {
	g := arbor.New()

	t.Run("one", func(t *testing.T) {
		at := arbor.NewT(t)
		g.Group("one")
		g.Append(at, "testOne", testOne)
	})

	t.Run("two", func(t *testing.T) {
		at := arbor.NewT(t)
		g.Group("two")
		g.Append(at, "testTwo", testTwo)
	})

	output := g.JSON()

	arbor.Upload(output)
}
`
		assert.Equal(t, expected, outFile.contents)
	})
}

func TestTwoGroupsWithDependecies(t *testing.T) {
	var src = `package random

import (
	"testing"

	"github.com/anatollupacescu/arbortest/runner"
)

// group:one after:z
func testOne(t *runner.T) {}

// group:z
func testNotEmpty(t *runner.T) {}
`

	var (
		testProviderFile = TestFile(src)
		singleFileDir    = TestDir(func() []arbor.File {
			return []arbor.File{&testProviderFile}
		})
		outFile = &TestOutFile{}
	)

	t.Run("registers two groups in correct order", func(t *testing.T) {
		err := arbor.Generate(&singleFileDir, outFile, "sample")
		assert.NoError(t, err)
		expected := `package sample

import (
	"testing"

	arbor "github.com/anatollupacescu/arbortest/runner"
)

func TestArbor(t *testing.T) {
	g := arbor.New()

	t.Run("z", func(t *testing.T) {
		at := arbor.NewT(t)
		g.Group("z")
		g.Append(at, "testNotEmpty", testNotEmpty)
	})

	t.Run("one", func(t *testing.T) {
		at := arbor.NewT(t)
		g.Group("one")
		g.After(at, "z")
		g.Append(at, "testOne", testOne)
	})

	output := g.JSON()

	arbor.Upload(output)
}
`
		assert.Equal(t, expected, outFile.contents)
	})
}

func TestThreeGroupsTwoConnected(t *testing.T) {
	var src = `package random

import (
	"testing"

	"github.com/anatollupacescu/arbortest/runner"
)

// group:one after:z
func testOne(t *runner.T) {}

// group:z
func testNotEmpty(t *runner.T) {}

// group:two
func testTwo(t *runner.T) {}
`

	var (
		testProviderFile = TestFile(src)
		singleFileDir    = TestDir(func() []arbor.File {
			return []arbor.File{&testProviderFile}
		})
		outFile = &TestOutFile{}
	)

	t.Run("registers two groups in correct order", func(t *testing.T) {
		err := arbor.Generate(&singleFileDir, outFile, "sample")
		assert.NoError(t, err)
		expected := `package sample

import (
	"testing"

	arbor "github.com/anatollupacescu/arbortest/runner"
)

func TestArbor(t *testing.T) {
	g := arbor.New()

	t.Run("two", func(t *testing.T) {
		at := arbor.NewT(t)
		g.Group("two")
		g.Append(at, "testTwo", testTwo)
	})

	t.Run("z", func(t *testing.T) {
		at := arbor.NewT(t)
		g.Group("z")
		g.Append(at, "testNotEmpty", testNotEmpty)
	})

	t.Run("one", func(t *testing.T) {
		at := arbor.NewT(t)
		g.Group("one")
		g.After(at, "z")
		g.Append(at, "testOne", testOne)
	})

	output := g.JSON()

	arbor.Upload(output)
}
`
		assert.Equal(t, expected, outFile.contents)
	})
}

func TestThreeGroupsAllConnected(t *testing.T) {
	var src = `package random

import (
	"testing"

	"github.com/anatollupacescu/arbortest/runner"
)

// group:one after:z
func testOne(t *runner.T) {}

// group:z
func testNotEmpty(t *runner.T) {}

// group:two after:one
func testTwo(t *runner.T) {}
`

	var (
		testProviderFile = TestFile(src)
		singleFileDir    = TestDir(func() []arbor.File {
			return []arbor.File{&testProviderFile}
		})
		outFile = &TestOutFile{}
	)

	t.Run("registers two groups in correct order", func(t *testing.T) {
		err := arbor.Generate(&singleFileDir, outFile, "sample")
		assert.NoError(t, err)
		expected := `package sample

import (
	"testing"

	arbor "github.com/anatollupacescu/arbortest/runner"
)

func TestArbor(t *testing.T) {
	g := arbor.New()

	t.Run("z", func(t *testing.T) {
		at := arbor.NewT(t)
		g.Group("z")
		g.Append(at, "testNotEmpty", testNotEmpty)
	})

	t.Run("one", func(t *testing.T) {
		at := arbor.NewT(t)
		g.Group("one")
		g.After(at, "z")
		g.Append(at, "testOne", testOne)
	})

	t.Run("two", func(t *testing.T) {
		at := arbor.NewT(t)
		g.Group("two")
		g.After(at, "one")
		g.Append(at, "testTwo", testTwo)
	})

	output := g.JSON()

	arbor.Upload(output)
}
`
		assert.Equal(t, expected, outFile.contents)
	})
}

func TestTwoFiles(t *testing.T) {
	var src1 = `package random

import (
	"testing"

	"github.com/anatollupacescu/arbortest/runner"
)

// group:one
func testOne(t *runner.T) {}

// group:two
func testTwo(t *runner.T) {}
`

	var src2 = `package random

import (
	"testing"

	"github.com/anatollupacescu/arbortest/runner"
)

// group:one after:two
func testThree(t *runner.T) {}
`

	var (
		file1      = TestFile(src1)
		file2      = TestFile(src2)
		twoFileDir = TestDir(func() []arbor.File {
			return []arbor.File{&file1, &file2}
		})
		outFile = &TestOutFile{}
	)

	t.Run("parses both files", func(t *testing.T) {
		err := arbor.Generate(&twoFileDir, outFile, "sample")
		assert.NoError(t, err)
		expected := `package sample

import (
	"testing"

	arbor "github.com/anatollupacescu/arbortest/runner"
)

func TestArbor(t *testing.T) {
	g := arbor.New()

	t.Run("two", func(t *testing.T) {
		at := arbor.NewT(t)
		g.Group("two")
		g.Append(at, "testTwo", testTwo)
	})

	t.Run("one", func(t *testing.T) {
		at := arbor.NewT(t)
		g.Group("one")
		g.After(at, "two")
		g.Append(at, "testOne", testOne)
		g.Append(at, "testThree", testThree)
	})

	output := g.JSON()

	arbor.Upload(output)
}
`
		assert.Equal(t, expected, outFile.contents)
	})
}

//helpers
type TestDir func() []arbor.File

func (t *TestDir) List() (out []arbor.File) {
	return (*t)()
}

type TestFile string

func (t *TestFile) Read() string {
	return string(*t)
}

type TestOutFile struct {
	contents string
	err      error
}

func (t *TestOutFile) Write(contents string) error {
	t.contents = contents
	return t.err
}
