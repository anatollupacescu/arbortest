package arbor_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/anatollupacescu/sandbox/arbor"
	"github.com/stretchr/testify/assert"
)

func TestNoProvider(t *testing.T) {
	all := []struct {
		name         string
		validators   map[string]string
		dependencies map[string][]string
		tests        map[string]func() error
		output       string
	}{
		{
			name: "single failing test",
			tests: map[string]func() error{
				"test1": func() error {
					return errors.New("fail")
				},
			},
			output: `{"nodes":[{"id":"test1","group":2,"status":"failed"}],"links":[]}`,
		}, {
			name: "single successful test",
			tests: map[string]func() error{
				"test1": func() error {
					return nil
				},
			},
			output: `{"nodes":[{"id":"test1","group":1,"status":"passed"}],"links":[]}`,
		},
	}

	for _, test := range all {
		t.Run(test.name, func(t *testing.T) {
			res := arbor.Run(test.validators, test.dependencies, test.tests)
			assert.Equal(t, test.output, res.Output)
		})
	}
}

func TestInvalidProvider(t *testing.T) {
	validators := map[string]string{
		"test1": "provider1",
		"test2": "provider2",
	}
	dependencies := map[string][]string{
		"test3": {"provider1", "provider2"},
	}

	var count int

	tests := map[string]func() error{
		"test1": func() error {
			return nil
		},
		"test2": func() error {
			return errors.New("provider invalid")
		},
		"test3": func() error {
			count++
			return nil
		},
	}

	output := `{"nodes":[
		{"id":"test1","group":1,"status":"passed"}, 
		{"id":"provider1","group":1,"status":"passed"},
		{"id":"test2","group":2,"status":"failed"}, 
		{"id":"provider2","group":2,"status":"failed"},
		{"id":"test3","group":0,"status":"skipped"}],
	"links":[
		{"source":"test1","target":"provider1","value":0}, 
		{"source":"test2","target":"provider2","value":0},
		{"source":"test3","target":"provider1","value":0},
		{"source":"test3","target":"provider2","value":0}]}`

	res := arbor.Run(validators, dependencies, tests)

	output = strings.ReplaceAll(output, "\n", "")
	output = strings.ReplaceAll(output, "\t", "")
	output = strings.ReplaceAll(output, " ", "")

	assert.Equal(t, output, res.Output)
	assert.Equal(t, 0, count)
}

func TestValidProvider(t *testing.T) {
	validators := map[string]string{
		"test1": "provider1",
		"test2": "provider2",
	}
	dependencies := map[string][]string{
		"test3": {"provider1", "provider2"},
	}

	var count int

	tests := map[string]func() error{
		"test1": func() error {
			return nil
		},
		"test2": func() error {
			return nil
		},
		"test3": func() error {
			count++
			return nil
		},
	}

	output := `{"nodes":[
		{"id":"test1","group":1,"status":"passed"}, 
		{"id":"provider1","group":1,"status":"passed"},
		{"id":"test2","group":1,"status":"passed"}, 
		{"id":"provider2","group":1,"status":"passed"},
		{"id":"test3","group":1,"status":"passed"}],
	"links":[
		{"source":"test1","target":"provider1","value":0}, 
		{"source":"test2","target":"provider2","value":0},
		{"source":"test3","target":"provider1","value":0},
		{"source":"test3","target":"provider2","value":0}]}`

	res := arbor.Run(validators, dependencies, tests)

	output = strings.ReplaceAll(output, "\n", "")
	output = strings.ReplaceAll(output, "\t", "")
	output = strings.ReplaceAll(output, " ", "")

	assert.Equal(t, output, res.Output)
	assert.Equal(t, 1, count)
}
