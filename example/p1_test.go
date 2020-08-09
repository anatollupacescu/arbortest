package example

import (
	"github.com/anatollupacescu/arbortest/runner"
)

// group:one after:two
func testOne(t *runner.T) {
}

// group:zorro after:one,two
func testNotEmpty(t *runner.T) {
}
