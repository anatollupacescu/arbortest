package runner_test

import (
	"strings"
	"testing"

	"github.com/anatollupacescu/arbortest/runner"
	"github.com/stretchr/testify/assert"
)

func TestSimple(t *testing.T) {
	t.Run("given a minimal dependency configuration", func(t *testing.T) {
		deps := map[string][]string{
			"test1": {"provider1"},
		}

		var count int
		tests := map[string]func(t *testing.T){
			"test1": incOK(&count),
		}

		output := runner.Run(t, deps, tests)

		t.Run("all tests pass", func(t *testing.T) {
			assert.Equal(t, 1, count)
		})
		t.Run("given", func(t *testing.T) {
			assert.False(t, t.Failed())
		})
		t.Run("correct output", func(t *testing.T) {
			expected := `
{
	"nodes": [
		{"id": "test1",	"group":2,"status":"pass"},
		{"id": "provider1","group":2,"status":"pass"}],
	"links": [{
		"source": "test1",
		"target": "provider1",
		"value": 3
	}]
}
`
			expected = strings.ReplaceAll(expected, "\t", "")
			expected = strings.ReplaceAll(expected, "\n", "")
			expected = strings.ReplaceAll(expected, " ", "")
			assert.Equal(t, expected, output)
		})
	})
}

func TestComplex(t *testing.T) {
	t.Run("given a minimal dependency configuration", func(t *testing.T) {
		deps := map[string][]string{
			"test1": {"provider1"},
			"test2": {"provider2"},
			"test3": {"provider1", "provider2"},
		}

		t.Run("when both providers are validated", func(t *testing.T) {
			var count int
			tests := map[string]func(t *testing.T){
				"test1": incOK(&count),
				"test2": incOK(&count),
				"test3": incOK(&count),
			}

			output := runner.Run(t, deps, tests)

			t.Run("all tests pass", func(t *testing.T) {
				assert.Equal(t, 3, count)
			})
			t.Run("correct output", func(t *testing.T) {
				expected := `
{
	"nodes": [
		{"id": "test1",	"group":2,"status":"pass"},
		{"id": "test2",	"group":2,"status":"pass"},
		{"id": "test3",	"group":2,"status":"pass"},
		{"id": "provider1","group":2,"status":"pass"},
		{"id": "provider2","group":2,"status":"pass"}],
	"links": [
		{"source": "test1","target": "provider1","value": 3},
		{"source": "test2","target": "provider2","value": 3},
		{"source": "test3","target": "provider1","value": 3},
		{"source": "test3","target": "provider2","value": 3}]
}`
				expected = strings.ReplaceAll(expected, "\t", "")
				expected = strings.ReplaceAll(expected, "\n", "")
				expected = strings.ReplaceAll(expected, " ", "")
				assert.Equal(t, expected, output)
			})
		})
	})
}

func incOK(counter *int) func(t *testing.T) {
	return func(t *testing.T) {
		*counter++
	}
}
