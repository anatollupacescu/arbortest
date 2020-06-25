package arbor_test

import (
	"testing"

	"github.com/anatollupacescu/arbortest/arbor"
	"github.com/stretchr/testify/assert"
)

func TestSingleProvider(t *testing.T) {
	var src = `package sample

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
}`

	pr := arbor.Parse(src)

	output := arbor.GenerateSource(pr.Tests)

	expected := `package arbor

import (
	"testing"
	"github.com/anatollupacescu/arbortest/arbor"
)

func TestArbor(t *testing.T) {
	validators := map[string]string{
		"testOne": "providerOne", "testTwo": "providerTwo", 
	}
	dependencies := map[string][]string{
		"testMain": {"providerOne", "providerTwo", }, 
	}
	tests := map[string]func() error{
		"testMain": testMain, 
	}
	if r := arbor.Run(validators, dependencies, tests); r.Error != "" {
		t.Error(r.Error)
	}
}
`

	assert.Equal(t, expected, output)
}
