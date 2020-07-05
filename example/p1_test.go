package example

import (
	"testing"
)

func testOne(t *testing.T) {
	_ = providerOne()
}

func providerOne() int {
	return 0
}
