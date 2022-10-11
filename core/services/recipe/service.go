package recipe

import (
	"cardamom/core/ext/log_ext"
	m "cardamom/core/models"

	"github.com/thoas/go-funk"
)

func filterRecipesByIngredients(
	inventoryItems []m.InventoryItem, recipes []m.Recipe) []m.Recipe {

	inventory := funk.Map(inventoryItems, func(i m.InventoryItem) string { return i.Item }).([]string)

	return funk.Filter(recipes,
		func(r m.Recipe) bool {
			return funk.Reduce(
				funk.Map(r.Ingredients, func(i m.RecipeIngredient) bool {
					return funk.Contains(inventory, i.Item)
				}).([]bool),
				func(a, b bool) bool { return a && b },
				true,
			).(bool)
		},
	).([]m.Recipe)
}

func resizeInstructions(user_uid string, instructions []string, recipe *m.Recipe) error {

	// resize instructions
	if len(instructions) < len(recipe.Instructions) {
		for i := len(instructions); i < len(recipe.Instructions); i++ {
			if err := m.DB.Delete(recipe.Instructions[i]).Error; err != nil {
				return log_ext.Errorf("deleting instructions -- %w", err)
			}
		}

		recipe.Instructions = recipe.Instructions[:len(instructions)]
	}

	// Copy, reorder and extend
	for i, instr := range instructions {
		if i >= len(recipe.Instructions) {
			recipe.Instructions = append(recipe.Instructions, m.RecipeInstruction{
				Uid:       generateInstrUid(),
				UserUid:   user_uid,
				RecipeUid: recipe.Uid,
				Meal:      recipe.Meal,
				SortOrder: i,
				Text:      instr,
			})
		} else {
			recipe.Instructions[i].Meal = recipe.Meal
			recipe.Instructions[i].Text = instr
			recipe.Instructions[i].SortOrder = i
		}

		if err := m.DB.Save(&recipe.Instructions[i]).Error; err != nil {
			return log_ext.Errorf("updating instruction(%d) -- %w", i, err)
		}
	}

	return nil
}

func resizeIngredients(user_uid string, ingredients []IngredientPart, recipe *m.Recipe) error {

	// Reize ingredients
	if len(ingredients) < len(recipe.Ingredients) {
		for i := len(ingredients); i < len(recipe.Ingredients); i++ {
			if err := m.DB.Delete(recipe.Ingredients[i]).Error; err != nil {
				return log_ext.Errorf("deleting instructions -- %w", err)
			}
		}

		recipe.Ingredients = recipe.Ingredients[:len(ingredients)]
	}

	for i, ingre := range ingredients {
		if i >= len(recipe.Ingredients) {
			recipe.Ingredients = append(recipe.Ingredients, m.RecipeIngredient{
				Uid:       generateIngreUid(),
				UserUid:   user_uid,
				RecipeUid: recipe.Uid,
				Meal:      recipe.Meal,
				SortOrder: i,
				Quantity:  ingre.Quantity,
				Unit:      ingre.Unit,
				Item:      ingre.Item,
			})
		} else {
			recipe.Ingredients[i].Meal = recipe.Meal
			recipe.Ingredients[i].Quantity = ingre.Quantity
			recipe.Ingredients[i].Unit = ingre.Unit
			recipe.Ingredients[i].Item = ingre.Item
			recipe.Ingredients[i].SortOrder = i
		}

		if err := m.DB.Save(&recipe.Ingredients[i]).Error; err != nil {
			return log_ext.Errorf("updating ingredient(%d) -- %w", i, err)
		}
	}

	return nil
}
