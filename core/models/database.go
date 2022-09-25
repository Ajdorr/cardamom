package models

import (
	cfg "cardamom/core/config"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Migrate() {
	DB.AutoMigrate(
		User{}, OAuthState{},
		GroceryItem{}, InventoryItem{},
		Recipe{}, RecipeIngredient{}, RecipeInstruction{})
}

func init() {

	if cfg.C.Env == "local" {
		if db, err := gorm.Open(sqlite.Open("core.db"), &gorm.Config{}); err != nil {
			panic(fmt.Errorf("failed to connect to database -- %w", err))
		} else {
			DB = db
		}
	} else if db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.C.DB_Host, cfg.C.DB_Port, cfg.C.DB_Username, cfg.C.DB_Password, cfg.C.DB_Name),
		PreferSimpleProtocol: true,
	}), &gorm.Config{}); err != nil {
		panic(fmt.Errorf("failed to connect to database -- %w", err))
	} else {
		DB = db
	}

}
