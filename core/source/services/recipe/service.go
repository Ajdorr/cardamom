package recipe

import (
	m "cardamom/core/source/db/models"

	"github.com/thoas/go-funk"
)

func addAlwaysAvailableIngredients(inventory []string) []string {
	return append(inventory, "water")
}

func filterRecipesByIngredients(
	inventoryItems []m.InventoryItem, recipes []m.Recipe) []m.Recipe {

	completeInventory := addAlwaysAvailableIngredients(
		funk.Map(inventoryItems, func(i m.InventoryItem) string { return i.Item }).([]string))

	return funk.Filter(recipes,
		func(r m.Recipe) bool {
			return funk.Reduce(
				funk.Map(r.Ingredients, func(i m.RecipeIngredient) bool {
					return i.Optional || funk.Contains(completeInventory, i.Item)
				}).([]bool),
				func(a, b bool) bool { return a && b },
				true,
			).(bool)
		},
	).([]m.Recipe)
}
