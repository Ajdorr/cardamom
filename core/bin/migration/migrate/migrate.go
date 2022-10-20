package main

import (
	_ "cardamom/core/config"
	md "cardamom/core/models"

	"gorm.io/gorm"
)

func migrate(mg gorm.Migrator) {

	md.DB.Exec(`update recipes set meal = 'dessert' where meal = 'dessert'`)
	// Constraints aren't automatically dropped automatically, ensure you add them
}

func main() {
	migrate(md.DB.Migrator())
}
