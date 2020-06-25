package arbor

func providerOne() int {
	return 1
}

func testOne() error {
	return nil
}

func testTwo() error {
	_ = providerOne()
	return nil
}
