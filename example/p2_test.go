package main

func testTwo() error {
	_ = providerTwo()
	return nil
}

func providerTwo() int {
	return 0
}

func testMain() error {
	_ = providerOne()
	_ = providerTwo()
	return nil
}
