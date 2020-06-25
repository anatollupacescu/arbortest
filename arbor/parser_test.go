package arbor_test

import (
	"testing"

	"github.com/anatollupacescu/sandbox/arbor"
	"github.com/stretchr/testify/assert"
)

type validators []func() error

func TestDisableRecipe(t *testing.T) {
	t.Run("given no tests", func(t *testing.T) {
		src := `package main

func providerOne() int {
	return 0
}`

		t.Run("should return warning", func(t *testing.T) {
			res := arbor.Parse(src)

			assert.Len(t, res.Warnings, 1)
			assert.Equal(t, "warning: no tests found", res.Warnings[0])
		})
	})

	t.Run("given a test and an unrelated provider", func(t *testing.T) {
		src := `package main

func testOne() error {
	return nil
}

func providerOne() int {
	return 0
}`

		t.Run("should return warning", func(t *testing.T) {
			res := arbor.Parse(src)
			assert.Equal(t, "warning: 'providerOne' declared but not used", res.Warnings[0])
			assert.Len(t, res.Tests, 1)
		})
	})

	t.Run("given a test and a related provider", func(t *testing.T) {
		src := `package main

func testOne() error {
	_ = providerOne()
	return nil
}

func providerOne() int {
	return 0
}`

		t.Run("should register test", func(t *testing.T) {
			res := arbor.Parse(src)
			assert.Len(t, res.Tests, 1)
			assert.Len(t, res.Warnings, 0)
		})
	})

	t.Run("given a test and an invalidated provider", func(t *testing.T) {
		src := `package main

func testOne() error {
	_ = providerOne()
	return nil
}

func providerOne() int {
	return 0
}

func testTwo() error {
	_ = providerOne()
	_ = providerTwo()
	return nil
}

func providerTwo() int {
	return 0
}`

		t.Run("should fail with message", func(t *testing.T) {
			res := arbor.Parse(src)
			assert.Len(t, res.Tests, 2)
			assert.Equal(t, "error: 'testTwo' calls invalid provider: 'providerTwo'", res.Error)
		})
	})
}
