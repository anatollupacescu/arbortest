package main_test

import (
	"io/ioutil"
	"log"
	"os/exec"
	"testing"
)

var (
	sampleTestFileName = "sample_test.go"
	generatedFileName  = "generated_test.go"
)

func TestSmoke(t *testing.T) {
	setUp()
	defer tearDown()

	cmd := exec.Command("bash", "-c", "go run .")
	if out, err := cmd.Output(); err != nil {
		t.Log(out)
		t.Errorf("go run: %s", err)
	}

	cmd = exec.Command("go", "test", "-v", sampleTestFileName, generatedFileName)
	if _, err := cmd.Output(); err != nil {
		t.Errorf("go test: %s", err)
	}
}

func setUp() {
	var src = `package main

func testOne() error {
	_ = providerOne()
	return nil
}

func providerOne() int {
	return 0
}`

	err := ioutil.WriteFile(sampleTestFileName, []byte(src), 0600)

	if err != nil {
		log.Fatalf("write sample test file: %s", err)
	}
}

func tearDown() {
	cmd := exec.Command("rm", sampleTestFileName, generatedFileName)
	if out, err := cmd.Output(); err != nil {
		log.Println("tear down", out)
		log.Fatal(out)
	}
}
