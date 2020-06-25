package example

import (
	"log"
)

func testOne() error {
	_ = providerOne()
	return nil
}

func providerOne() int {
	return 0
}

func testTwo() error {
	log.Println("called")
	_ = providerOne()
	_ = providerTwo()
	return nil
}

func validateTwo() error {
	_ = providerTwo()
	return nil
}

func providerTwo() int {
	return 0
}
