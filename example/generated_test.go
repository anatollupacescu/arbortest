package example

import (
	"testing"

	"github.com/anatollupacescu/arbortest/arbor"
)

func TestArbor(t *testing.T) {
	validators := map[string]string{
		"testOne": "providerOne", "testTwo": "providerTwo",
	}
	dependencies := map[string][]string{
		"testMain": {"providerOne", "providerTwo"},
	}
	tests := map[string]func() error{
		"testOne": testOne, "testTwo": testTwo, "testMain": testMain,
	}
	r := arbor.Run(validators, dependencies, tests)
	t.Logf("output: %v", r.Output)
}
