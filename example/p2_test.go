package example

import (
	"testing"
)

func testTwo(t *testing.T) {
	_ = providerTwo()
}

func providerTwo() int {
	return 0
}

func testMain(t *testing.T) {
	_ = providerOne()
	_ = providerTwo()
}
