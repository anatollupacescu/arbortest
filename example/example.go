package example

func testOne() error {
	_ = providerOne()
	return nil
}

func providerOne() int {
	return 0
}

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
