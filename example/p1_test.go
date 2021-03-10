package example

import (
	"testing"

	"github.com/anatollupacescu/arbortest/runner"
)

// group:recipe after:ingredient
func testUpdateSuccess(t *runner.T) {
	t.Run("ignored atm", func(*testing.T) {
		t.Error("no")
	})
}

// group:order after:recipe,ingredient
func testRejectsEmptyQuantity(t *runner.T) {
	t.Run("ignored atm", func(*testing.T) {
	})
}
