package example

import (
	"testing"

	"github.com/anatollupacescu/sandbox/arbor"
	"github.com/stretchr/testify/assert"
)

func TestArbor(t *testing.T) {
	validators := map[string]string{
		"testOne":     "providerOne",
		"validateTwo": "providerTwo",
	}

	dependencies := map[string][]string{
		"testTwo": {"providerOne", "providerTwo"},
	}

	tests := map[string]func() error{
		"testOne":     testOne,
		"validateTwo": validateTwo,
		"testTwo":     testTwo,
	}

	res := arbor.Run(validators, dependencies, tests)

	expected := `{"nodes":[{"id":"testOne","group":2,"status":"failed"}],"links":[]}`
	assert.Equal(t, expected, res.Output)
}
