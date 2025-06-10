package models

import (
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

type GroceryItem struct {
	Uid         string    `gorm:"primaryKey;not null;default:null" json:"uid"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UserUid     string    `gorm:"index" json:"user_uid"`
	Item        string    `json:"item"`
	Store       string    `json:"store"`
	IsCollected bool      `json:"is_collected"`
}

func (g *GroceryItem) BeforeCreate(tx *gorm.DB) error {
	g.Uid = gonanoid.Must()
	return nil
}
