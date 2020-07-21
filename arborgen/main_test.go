package main_test

import (
	"os/exec"
	"testing"

	"github.com/anatollupacescu/arbortest/runner"
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

// group:one
func testOne(t *runner.T) {
}
