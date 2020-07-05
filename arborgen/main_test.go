package main_test

import (
	"os/exec"
	"testing"
)

const generatedFileName = "generated_test.go"

func TestSmoke(t *testing.T) {
	cmd := exec.Command("bash", "-c", "go run . -pkg=main_test")
	if _, err := cmd.Output(); err != nil {
		t.Errorf("go run: %s", err)
		return
	}

	defer tearDown(t)

	cmd = exec.Command("go", "test", "-v", "-count=1", "-run", "^(TestArbor)$")
	if _, err := cmd.Output(); err != nil {
		t.Errorf("go test: %s", err)
	}
}

func tearDown(t *testing.T) {
	cmd := exec.Command("rm", generatedFileName)
	if _, err := cmd.Output(); err != nil {
		t.Errorf("tear down: %s", err)
	}
}

func testOne(t *testing.T) {
	_ = providerOne()
}

func providerOne() int {
	return 0
}

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

func TestIntegrationListTestFiles(t *testing.T) {
	t.Skip() //TODO implementation pending
}

func TestIntegrationReadContents(t *testing.T) {
	t.Skip() //TODO implementation pending
}

func TestIntegrationWriteContents(t *testing.T) {
	t.Skip() //TODO implementation pending
}
