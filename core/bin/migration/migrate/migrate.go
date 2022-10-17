package main

import (
	_ "cardamom/core/config"
	md "cardamom/core/models"

	"gorm.io/gorm"
)

func migrate(mg gorm.Migrator) {
	mg.AddColumn(&md.RecipeIngredient{}, "optional")
	mg.AddColumn(&md.RecipeIngredient{}, "modifier")

	md.DB.Exec(`update recipe_ingredients set unit = '' where unit is null`)
	// Constraints aren't automatically dropped automatically, ensure you add them
}

func main() {
	migrate(md.DB.Migrator())
}
