package main

import (
	gorm_ext "cardamom/core/ext/gorm_ext"
	m "cardamom/core/models"
	"fmt"
	"os"
	"strings"
)

var fileTemplate = `
  package main

import (
	_ "cardamom/core/config"
	ge "cardamom/core/ext/gorm_ext"
	md "cardamom/core/models"

	"gorm.io/gorm"
)

func migrate(mg gorm.Migrator) {
  %s
	// Constraints aren't automatically dropped automatically, ensure you add them
}

func main() {
	migrate(md.DB.Migrator())
}`

func Generate() {
	changes, err := gorm_ext.GenerateMigration(
		m.User{}, m.GroceryItem{}, m.InventoryItem{}, m.OAuthState{},
		m.Recipe{}, m.RecipeIngredient{},
	)
	if err != nil {
		panic(err)
	}

	fileText := fmt.Sprintf(fileTemplate, strings.Join(changes, "\n  "))
	if fp := os.Getenv("CARDAMOM_CORE_HOME"); len(fp) > 0 {
		os.WriteFile(fp+"/bin/migration/migrate/migrate.go", []byte(fileText), os.ModePerm)
	} else {
		fmt.Print(fileText)
	}
}

func main() {
	Generate()
}
