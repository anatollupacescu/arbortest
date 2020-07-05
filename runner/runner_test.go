package runner_test

import (
	"testing"

	"github.com/anatollupacescu/arbortest/runner"
	"github.com/stretchr/testify/assert"
)

func TestRunner(t *testing.T) {
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

			runner.Run(t, deps, tests)

			t.Run("all tests pass", func(t *testing.T) {
				assert.Equal(t, 3, count)
			})
		})
	})
}

func incOK(counter *int) func(t *testing.T) {
	return func(t *testing.T) {
		*counter++
	}
}
