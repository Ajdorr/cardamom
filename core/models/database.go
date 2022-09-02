package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	db, err := gorm.Open(sqlite.Open("core.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	DB = db
}

func Migrate() {
	DB.AutoMigrate(
		User{}, OAuthState{},
		GroceryItem{}, InventoryItem{},
		Recipe{}, RecipeIngredient{}, RecipeInstruction{})
}
