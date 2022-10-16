package main

import (
	_ "cardamom/core/config"
	md "cardamom/core/models"
	"strings"

	"github.com/thoas/go-funk"
	"gorm.io/gorm"
)

type RecipeResult struct {
	Text string
}

func migrate(mg gorm.Migrator) {
	mg.AddColumn(&md.Recipe{}, "instructions")

	var recipes []md.Recipe
	md.DB.Find(&recipes)
	for _, recipe := range recipes {
		var results []RecipeResult
		md.DB.Raw("select text from recipe_instructions where recipe_uid = ?", recipe.Uid).Scan(&results)

		recipe.Instructions = strings.Join(funk.Map(
			results, func(r RecipeResult) string { return r.Text }).([]string), "\n")
		md.DB.Save(&recipe)
	}
	mg.DropTable("recipe_instructions")
	// Constraints aren't automatically dropped automatically, ensure you add them
}

func main() {
	migrate(md.DB.Migrator())
}
