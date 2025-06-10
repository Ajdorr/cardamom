package models

import (
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

type InventoryCategory string

const (
	COOKING         InventoryCategory = "cooking"
	SPICES          InventoryCategory = "spices"
	SAUCES          InventoryCategory = "sauces" // and condiments
	NON_PERISHABLES InventoryCategory = "non-perishables"
	NON_COOKING     InventoryCategory = "non-cooking"
)

var ValidCategories = []InventoryCategory{
	COOKING, SPICES, SAUCES, NON_PERISHABLES, NON_COOKING,
}

type InventoryItem struct {
	Uid       string            `gorm:"primaryKey;not null;default:null" json:"uid"`
	UpdatedAt time.Time         `json:"updated_at"`
	CreatedAt time.Time         `json:"created_at"`
	UserUid   string            `gorm:"index" json:"user_uid"`
	Item      string            `json:"item"`
	InStock   bool              `json:"in_stock"`
	Category  InventoryCategory `gorm:"default:cooking" json:"category"`
}

func (i *InventoryItem) BeforeCreate(tx *gorm.DB) error {
	i.Uid = gonanoid.Must()
	return nil
}
