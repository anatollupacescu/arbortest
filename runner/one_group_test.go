package runner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSingleGroupSingleTest(t *testing.T) {
	var counter int

	g := New()

	mock := &fakeT{}
	at := NewT(mock)
	g.Group("testGroup")
	g.Append(at, "test", func(at *T) {
		counter++
	})

	assert.Equal(t, 1, counter)
	assert.False(t, at.Failed())
	assert.Len(t, g.groups, 1)
	assert.Equal(t, pass, g.groups.get("testGroup").status)
	assert.Len(t, g.groups.get("testGroup").tests, 1)

	real := g.groups.get("testGroup").tests[0]
	expected := test{
		name:   "test",
		status: pass,
	}
	assert.Equal(t, expected, real)
}

func TestOneTestFail(t *testing.T) {
	var counter int

	g := New()

	fakeT := &fakeT{}
	at := NewT(fakeT)
	g.Group("testGroup")
	g.Append(at, "test", func(at *T) {
		counter++
		at.Error("expected")
	})

	assert.Equal(t, 1, counter)
	assert.True(t, at.Failed())
	assert.Equal(t, fail, g.groups.get("testGroup").status)
	assert.Len(t, g.groups.get("testGroup").tests, 1)

	real := g.groups.get("testGroup").tests[0]
	expected := test{
		name:   "test",
		status: fail,
	}
	assert.Equal(t, expected, real)
}

func TestTwoTests(t *testing.T) {
	var counter int

	g := New()

	fakeT := &fakeT{}
	at := NewT(fakeT)
	g.Group("testGroup")
	g.Append(at, "test1", func(at *T) {
		counter++
	})
	g.Append(at, "test2", func(at *T) {
		counter++
	})

	assert.Equal(t, 2, counter)
	assert.False(t, at.Failed())
	assert.Equal(t, pass, g.groups.get("testGroup").status)
	assert.Len(t, g.groups.get("testGroup").tests, 2)

	real := g.groups.get("testGroup").tests[0]
	expected := test{
		name:   "test1",
		status: pass,
	}
	assert.Equal(t, expected, real)

	real = g.groups.get("testGroup").tests[1]
	expected = test{
		name:   "test2",
		status: pass,
	}
	assert.Equal(t, expected, real)
}

func TestFirstFailSkipsFollowingTests(t *testing.T) {
	var counter int

	g := New()

	fakeT := &fakeT{}
	at := NewT(fakeT)
	g.Group("testGroup")
	g.Append(at, "test1", func(at *T) {
		counter++
		at.Error()
	})
	g.Append(at, "test2", func(at *T) {
		counter++
	})

	assert.Equal(t, 1, counter)
	assert.True(t, at.Failed())
	assert.Equal(t, fail, g.groups.get("testGroup").status)

	real := g.groups.get("testGroup").tests[0]
	expected := test{
		name:   "test1",
		status: fail,
	}
	assert.Equal(t, expected, real)

	real = g.groups.get("testGroup").tests[1]
	expected = test{
		name:   "test2",
		status: skip,
	}
	assert.Equal(t, expected, real)
}

// props
type fakeT struct {
	failed bool
}

func (f *fakeT) Failed() bool {
	return f.failed
}

func (f *fakeT) Error(args ...interface{}) {
	f.failed = true
}

func (f *fakeT) Errorf(format string, args ...interface{}) {
	f.failed = true
}

func (f *fakeT) Log(args ...interface{}) {
}
