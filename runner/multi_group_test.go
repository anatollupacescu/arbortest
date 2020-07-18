package runner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTwoIndependentGroups(t *testing.T) {
	var counter int

	g := New()

	mock := &fakeT{}
	at1 := NewT(mock)
	g.Group("testGroup")
	g.Append(at1, "test", func(at *T) {
		counter++
	})

	mock = &fakeT{}
	at2 := NewT(mock)
	g.Group("testGroup2")
	g.Append(at2, "test", func(at *T) {
		counter++
	})

	assert.False(t, at1.Failed())
	assert.False(t, at2.Failed())

	assert.Equal(t, 2, counter)
	assert.Len(t, g.groups, 2)

	assert.Equal(t, pass, g.groups.get("testGroup").status)
	assert.Equal(t, pass, g.groups.get("testGroup2").status)
}

func TestFailedDepSkipsDependant(t *testing.T) {
	var counter int

	g := New()

	mock := &fakeT{}
	at1 := NewT(mock)
	g.Group("testGroup")
	g.Append(at1, "test", func(at *T) {
		at.Error()
		counter++
	})

	mock = &fakeT{}
	at2 := NewT(mock)
	g.Group("testGroup2")
	g.After(at2, "testGroup")
	g.Append(at2, "test", func(at *T) {
		t.Error("should not have been called")
	})

	assert.Equal(t, 1, counter)
	assert.Equal(t, fail, g.groups.get("testGroup").status)
	assert.Equal(t, skip, g.groups.get("testGroup2").status)
}

func TestOneFailedDepSkipsDependant(t *testing.T) {
	var counter int

	g := New()

	mockT := &fakeT{}
	at1 := NewT(mockT)
	g.Group("testGroup1")
	g.Append(at1, "test", func(at *T) {
		at.Error()
		counter++
	})

	mockT = &fakeT{}
	at2 := NewT(mockT)
	g.Group("testGroup2")
	g.Append(at2, "test", func(at *T) {
		counter++
	})

	mockT = &fakeT{}
	at3 := NewT(mockT)
	g.Group("testGroup3")
	g.After(at3, "testGroup1", "testGroup2")
	g.Append(at3, "test", func(at *T) {
		t.Error("should not have been called")
	})

	assert.Equal(t, 2, counter)
	assert.Equal(t, fail, g.groups.get("testGroup1").status)
	assert.Equal(t, pass, g.groups.get("testGroup2").status)
	assert.Equal(t, skip, g.groups.get("testGroup3").status)
}
