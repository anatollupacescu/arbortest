package example

import (
	"testing"

	arbor "github.com/anatollupacescu/arbortest/runner"
)

func TestArbor(t *testing.T) {
	g := arbor.New()

	t.Run("ingredient", func(t *testing.T) {
		at := arbor.NewT(t)
		g.Group("ingredient")
		g.Append(at, "CreateSuccess", testCreateSuccess)
		g.Append(at, "UpdateWithNegativeValue", testUpdateWithNegativeValue)
	})

	t.Run("recipe", func(t *testing.T) {
		at := arbor.NewT(t)
		g.Group("recipe")
		g.After(at, "ingredient")
		g.Append(at, "UpdateSuccess", testUpdateSuccess)
	})

	t.Run("order", func(t *testing.T) {
		at := arbor.NewT(t)
		g.Group("order")
		g.After(at, "recipe","ingredient")
		g.Append(at, "RejectsEmptyQuantity", testRejectsEmptyQuantity)
	})

	output := g.JSON()

	arbor.Upload(output)
}
