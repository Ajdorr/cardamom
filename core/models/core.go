package models

import "time"

type GroceryItem struct {
	Uid         string    `gorm:"primaryKey;not null;default:null" json:"uid"`
	CreatedAt   time.Time `json:"created_at"`
	UpdateAt    time.Time `json:"update_at"`
	UserUid     string    `gorm:"index" json:"user_uid"`
	Item        string    `json:"item"`
	Store       string    `json:"store"`
	IsCollected bool      `json:"is_collected"`
}

type InventoryItem struct {
	Uid       string    `gorm:"primaryKey;not null;default:null" json:"uid"`
	UpdateAt  time.Time `json:"update_at"`
	CreatedAt time.Time `json:"created_at"`
	UserUid   string    `gorm:"index" json:"user_uid"`
	Item      string    `json:"item"`
	InStock   bool      `json:"in_stock"`
}
