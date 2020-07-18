package runner_test

import (
	"testing"

	"github.com/anatollupacescu/arbortest/runner"
)

func TestExp(t *testing.T) {
	t.SkipNow()
	g := runner.New()

	t.Run("inventory item", func(t *testing.T) {
		at := runner.NewT(t)
		g.Group("inventory item")
		g.Append(at, "rejects empty", pass)
		g.Append(at, "rejects duplicate", fail)
	})

	t.Run("recipe", func(t *testing.T) {
		at := runner.NewT(t)
		g.Group("recipe")
		g.After(at, "inventory item")
		g.Append(at, "rejects no ingredients", pass)
		g.Append(at, "reject duplicate name", pass)
	})

	t.Run("order", func(t *testing.T) {
		at := runner.NewT(t)
		g.Group("order")
		g.After(at, "inventory item")
		g.Append(at, "reject zero quantity", pass)
		g.Append(at, "reject missing name", pass)
	})

	data := g.JSON()

	runner.Upload(data)
}

func pass(t *runner.T) {
}

func fail(t *runner.T) {
	t.Error("stop")
}
