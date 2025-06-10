package models

import (
	"time"
)

// FIXME convert to redis server
type OAuthState struct {
	IPAddress string    `gorm:"primaryKey;not null;default:null" json:"ip"`
	TTL       time.Time `gorm:"not null" json:"ttl"`
	Provider  string    `gorm:"not null" json:"provider"`
	State     string    `gorm:"not null" json:"state"`
}
