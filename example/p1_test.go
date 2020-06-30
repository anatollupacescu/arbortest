package example

func testOne() error {
	_ = providerOne()
	return nil
}

func providerOne() int {
	return 0
}
