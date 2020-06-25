package main

import (
	"testing"

	"github.com/anatollupacescu/arbortest/arbor"
)

func TestImports(t *testing.T) {
	src := `package sample

	func providerOne() int {
		return 1
	}
	
	func testOne() error {
		_ = providerOne()
		return nil
	}`

	r := arbor.Parse(src)

	if len(r.Tests) != 1 {
		t.Fail()
	}
}
