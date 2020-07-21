package example

import (
	"github.com/anatollupacescu/arbortest/runner"
)

// group:one after:two
func testOne(t *runner.T) {
}

// group:z after:one,two
func testNotEmpty(t *runner.T) {
}
