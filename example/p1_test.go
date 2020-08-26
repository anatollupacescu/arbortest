package example

import (
	"github.com/anatollupacescu/arbortest/runner"
)

// group:recipe after:ingredient
func testUpdateSuccess(t *runner.T) {
}

// group:order after:recipe,ingredient
func testRejectsEmptyQuantity(t *runner.T) {
}
