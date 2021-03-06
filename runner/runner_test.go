package runner_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/anatollupacescu/arbortest/runner"
)

func TestSingleAppend(t *testing.T) {
	rt := runner.NewT(t)
	r := runner.New()
	r.Group("group")
	r.Append(rt, "test", func(*runner.T) {})

	json := `{
		"commit":"test",
		"message":"test",
		"nodes":[
		{"id": "group", "status":"pass"},
		{"id": "test",  "status":"pass"}
	], 
	"links":[
		{"source": "test","target": "group","value": 1}
	]}`

	json = strings.ReplaceAll(json, "\t", "")
	json = strings.ReplaceAll(json, "\n", "")
	json = strings.ReplaceAll(json, " ", "")

	r.CommitInfoProvider(func() (string, string) {
		return "test", "test"
	})
	assert.Equal(t, json, r.JSON())
}

func TestTwoAppend(t *testing.T) {
	rt := runner.NewT(t)
	r := runner.New()
	r.Group("group")
	r.Append(rt, "test", func(*runner.T) {})
	r.Append(rt, "test2", func(*runner.T) {})

	json := `{
		"commit":"test",
		"message":"test",
		"nodes":[
		{"id": "group", "status":"pass"},
		{"id": "test",  "status":"pass"},
		{"id": "test2",  "status":"pass"}
	], 
	"links":[
		{"source": "test","target": "group","value": 1},
		{"source": "test2","target": "group","value": 1}
	]}`

	json = strings.ReplaceAll(json, "\t", "")
	json = strings.ReplaceAll(json, "\n", "")
	json = strings.ReplaceAll(json, " ", "")

	r.CommitInfoProvider(func() (string, string) {
		return "test", "test"
	})
	assert.Equal(t, json, r.JSON())
}

func TestTwoGrops(t *testing.T) {
	rt := runner.NewT(t)
	r := runner.New()
	r.Group("group")
	r.Append(rt, "test1", func(*runner.T) {})
	r.Append(rt, "test2", func(*runner.T) {})

	r.Group("group2")
	r.Append(rt, "test3", func(*runner.T) {})
	r.Append(rt, "test4", func(*runner.T) {})

	json := `{
		"commit":"test",
		"message":"test",
		"nodes":[
		{"id": "group", "status":"pass"},
		{"id": "test1", "status":"pass"},
		{"id": "test2", "status":"pass"},
		{"id": "group2","status":"pass"},
		{"id": "test3", "status":"pass"},
		{"id": "test4", "status":"pass"}
	], 
	"links":[
		{"source": "test1","target": "group","value": 1},
		{"source": "test2","target": "group","value": 1},
		{"source": "test3","target": "group2","value": 1},
		{"source": "test4","target": "group2","value": 1}
	]}`

	json = strings.ReplaceAll(json, "\t", "")
	json = strings.ReplaceAll(json, "\n", "")
	json = strings.ReplaceAll(json, " ", "")

	r.CommitInfoProvider(func() (string, string) {
		return "test", "test"
	})
	assert.Equal(t, json, r.JSON())
}

func TestAfterCreatesLink(t *testing.T) {
	rt := runner.NewT(t)
	r := runner.New()
	r.Group("group")
	r.Append(rt, "test", func(*runner.T) {})
	r.Append(rt, "test2", func(*runner.T) {})

	r.Group("group2")
	r.After(rt, "group")
	r.Append(rt, "test", func(*runner.T) {})
	r.Append(rt, "test2", func(*runner.T) {})

	json := `{
		"commit":"test",
		"message":"test",
		"nodes":[
		{"id":"group",	"status":"pass"},
		{"id":"test",	"status":"pass"},
		{"id":"test2",	"status":"pass"},
		{"id":"group2",	"status":"pass"},
		{"id":"test",	"status":"pass"},
		{"id":"test2",	"status":"pass"}
	],
	"links":[
		{"source":"test","target":"group","value":1},
		{"source":"test2","target":"group","value":1},
		{"source":"test","target":"group2","value":1},
		{"source":"test2","target":"group2","value":1},
		{"source":"group2","target":"group","value":3}
	]}`

	json = strings.ReplaceAll(json, "\t", "")
	json = strings.ReplaceAll(json, "\n", "")
	json = strings.ReplaceAll(json, " ", "")

	r.CommitInfoProvider(func() (string, string) {
		return "test", "test"
	})
	assert.Equal(t, json, r.JSON())
}

func TestAfterFailedGroup(t *testing.T) {
	r := runner.New()

	mock := &fakeT{}
	rt := runner.NewT(mock)
	r.Group("group")
	r.Append(rt, "test", func(*runner.T) {})
	r.Append(rt, "test2", func(at *runner.T) { at.Error("stop here") })

	mock = &fakeT{}
	rt = runner.NewT(mock)
	r.Group("group2")
	r.After(rt, "group")
	r.Append(rt, "test", func(*runner.T) {})
	r.Append(rt, "test2", func(*runner.T) {})

	json := `{
		"commit":"test",
		"message":"test",
		"nodes":[
		{"id":"group",	"status":"fail"},
		{"id":"test",	"status":"pass"},
		{"id":"test2",	"status":"fail"},
		{"id":"group2",	"status":"skip"},
		{"id":"test",	"status":"skip"},
		{"id":"test2",	"status":"skip"}
	],
	"links":[
		{"source":"test","target":"group","value":1},
		{"source":"test2","target":"group","value":1},
		{"source":"test","target":"group2","value":1},
		{"source":"test2","target":"group2","value":1},
		{"source":"group2","target":"group","value":3}
	]}`

	json = strings.ReplaceAll(json, "\t", "")
	json = strings.ReplaceAll(json, "\n", "")
	json = strings.ReplaceAll(json, " ", "")

	r.CommitInfoProvider(func() (string, string) {
		return "test", "test"
	})
	assert.Equal(t, json, r.JSON())
}

type fakeT struct {
	fail bool
}

func (f *fakeT) Failed() bool {
	return f.fail
}

func (f *fakeT) Error(args ...interface{}) {
	f.fail = true
}

func (f *fakeT) Errorf(format string, args ...interface{}) {
	f.fail = true
}

func (f *fakeT) Log(args ...interface{}) {
}

func (f *fakeT) Run(_ string, _ func(t *testing.T)) bool {
	return true
}
