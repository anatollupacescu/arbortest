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

	json := `{"nodes":[
		{"id": "group", "group":2, "status":"pass"},
		{"id": "test",  "group":2, "status":"pass"}
	], 
	"links":[
		{"source": "test","target": "group","value": 3}
	]}`

	json = strings.ReplaceAll(json, "\t", "")
	json = strings.ReplaceAll(json, "\n", "")
	json = strings.ReplaceAll(json, " ", "")
	assert.Equal(t, json, r.JSON())
}

func TestTwoAppend(t *testing.T) {
	rt := runner.NewT(t)
	r := runner.New()
	r.Group("group")
	r.Append(rt, "test", func(*runner.T) {})
	r.Append(rt, "test2", func(*runner.T) {})

	json := `{"nodes":[
		{"id": "group", "group":2, "status":"pass"},
		{"id": "test",  "group":2, "status":"pass"},
		{"id": "test2",  "group":2, "status":"pass"}
	], 
	"links":[
		{"source": "test","target": "group","value": 3},
		{"source": "test2","target": "group","value": 3}
	]}`

	json = strings.ReplaceAll(json, "\t", "")
	json = strings.ReplaceAll(json, "\n", "")
	json = strings.ReplaceAll(json, " ", "")
	assert.Equal(t, json, r.JSON())
}

func TestTwoGrops(t *testing.T) {
	rt := runner.NewT(t)
	r := runner.New()
	r.Group("group")
	r.Append(rt, "test", func(*runner.T) {})
	r.Append(rt, "test2", func(*runner.T) {})

	r.Group("group2")
	r.Append(rt, "test", func(*runner.T) {})
	r.Append(rt, "test2", func(*runner.T) {})

	json := `{"nodes":[
		{"id": "group", "group":2, "status":"pass"},
		{"id": "test",  "group":2, "status":"pass"},
		{"id": "test2",  "group":2, "status":"pass"},
		{"id": "group2", "group":2, "status":"pass"},
		{"id": "test",  "group":2, "status":"pass"},
		{"id": "test2",  "group":2, "status":"pass"}
	], 
	"links":[
		{"source": "test","target": "group","value": 3},
		{"source": "test2","target": "group","value": 3},
		{"source": "test","target": "group2","value": 3},
		{"source": "test2","target": "group2","value": 3}
	]}`

	json = strings.ReplaceAll(json, "\t", "")
	json = strings.ReplaceAll(json, "\n", "")
	json = strings.ReplaceAll(json, " ", "")
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

	json := `{"nodes":[
		{"id":"group","group":2,"status":"pass"},
		{"id":"test","group":2,"status":"pass"},
		{"id":"test2","group":2,"status":"pass"},
		{"id":"group2","group":2,"status":"pass"},
		{"id":"test","group":2,"status":"pass"},
		{"id":"test2","group":2,"status":"pass"},
		{"id":"group2-ext","group":2,"status":"pass"},
		{"id":"group-ext","group":2,"status":"pass"}],
	"links":[
		{"source":"test","target":"group","value":3},
		{"source":"test2","target":"group","value":3},
		{"source":"test","target":"group2","value":3},
		{"source":"test2","target":"group2","value":3},
		{"source":"group2-ext","target":"group2","value":3},
		{"source":"group2-ext","target":"group-ext","value":3},
		{"source":"group-ext","target":"group","value":3}]}`

	json = strings.ReplaceAll(json, "\t", "")
	json = strings.ReplaceAll(json, "\n", "")
	json = strings.ReplaceAll(json, " ", "")
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

	json := `{"nodes":[
		{"id":"group","group":1,"status":"fail"},
		{"id":"test","group":2,"status":"pass"},
		{"id":"test2","group":1,"status":"fail"},
		{"id":"group2","group":0,"status":"skip"},
		{"id":"test","group":0,"status":"skip"},
		{"id":"test2","group":0,"status":"skip"},
		{"id":"group2-ext","group":0,"status":"skip"},
		{"id":"group-ext","group":0,"status":"skip"}],
	"links":[
		{"source":"test","target":"group","value":3},
		{"source":"test2","target":"group","value":3},
		{"source":"test","target":"group2","value":3},
		{"source":"test2","target":"group2","value":3},
		{"source":"group2-ext","target":"group2","value":3},
		{"source":"group2-ext","target":"group-ext","value":3},
		{"source":"group-ext","target":"group","value":3}]}`

	json = strings.ReplaceAll(json, "\t", "")
	json = strings.ReplaceAll(json, "\n", "")
	json = strings.ReplaceAll(json, " ", "")
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
