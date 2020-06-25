package arbor_test

import (
	"errors"
	"testing"
)

func providerOne() int {
	return 1
}

func testOne() error {
	return errors.New("lol")
}

func testTwo() error {
	_ = providerOne()
	return nil
}

func Test(t *testing.T) {

}
