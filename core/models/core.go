package models

import "time"

type GroceryItem struct {
	Uid         string    `gorm:"primaryKey;not null;default:null" json:"uid"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UserUid     string    `gorm:"index" json:"user_uid"`
	Item        string    `json:"item"`
	Store       string    `json:"store"`
	IsCollected bool      `json:"is_collected"`
}

type InventoryCategory string

const (
	COOKING     InventoryCategory = "cooking"
	SPICES      InventoryCategory = "spices"
	SAUCES      InventoryCategory = "sauces" // and condiments
	NON_COOKING InventoryCategory = "non-cooking"
)

var ValidCategories = []InventoryCategory{
	COOKING, SPICES, SAUCES, NON_COOKING,
}

type InventoryItem struct {
	Uid       string    `gorm:"primaryKey;not null;default:null" json:"uid"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
	UserUid   string    `gorm:"index" json:"user_uid"`
	Item      string    `json:"item"`
	InStock   bool      `json:"in_stock"`
	Category  string    `gorm:"default:cooking" json:"category"`
}
