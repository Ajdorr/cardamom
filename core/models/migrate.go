package models

func Migrate() {
	DB.Migrator().RenameColumn(&GroceryItem{}, "update_at", "updated_at")
	DB.Migrator().RenameColumn(&InventoryItem{}, "update_at", "updated_at")
	DB.Migrator().RenameColumn(&Recipe{}, "update_at", "updated_at")
	DB.Migrator().RenameColumn(&RecipeIngredient{}, "update_at", "updated_at")
	DB.Migrator().RenameColumn(&RecipeInstruction{}, "update_at", "updated_at")

	DB.Migrator().AddColumn(&Recipe{}, "is_trashed")
	DB.Migrator().CreateIndex(&Recipe{}, "IsTrashed")
	DB.Migrator().AddColumn(&Recipe{}, "trash_at")
}
