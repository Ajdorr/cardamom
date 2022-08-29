package models

import (
	"crypto/sha512"
	"os"

	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
)

type role string

const (
	USER  role = "user"
	ADMIN role = "admin"
)

type User struct {
	gorm.Model
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdateAt  string `json:"update_at"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      role   `json:"role"`
	Email     string `json:"email"`
	Password  []byte `json:"password"`
}

func getSalt() []byte {
	return []byte(os.Getenv("PASSWORD_SALT"))
}

func HashPassword(password string) []byte {
	// Iteration count "should" double ever two years
	// https://en.wikipedia.org/wiki/PBKDF2#:~:text=In%202021%2C%20OWASP%20recommended%20to,for%20PBKDF2%2DHMAC%2DSHA512.
	return pbkdf2.Key([]byte(password), getSalt(), 120000, 64, sha512.New)
}

func IsPasswordMatch(user *User, password string) bool {
	return slices.Equal(user.Password, HashPassword(password))
}
