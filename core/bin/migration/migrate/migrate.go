package main

import (
	_ "cardamom/core/config"
	md "cardamom/core/models"

	"gorm.io/gorm"
)

func migrate(mg gorm.Migrator) {
	mg.AddColumn(&md.InventoryItem{}, "category")
	// Constraints aren't automatically dropped automatically, ensure you add them
}

func main() {
	migrate(md.DB.Migrator())
}
