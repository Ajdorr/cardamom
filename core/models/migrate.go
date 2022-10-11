package models

func AutoMigrate() {
	DB.AutoMigrate(
		User{}, OAuthState{},
		GroceryItem{}, InventoryItem{},
		Recipe{}, RecipeIngredient{}, RecipeInstruction{})
}
