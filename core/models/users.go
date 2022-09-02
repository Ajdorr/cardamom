package models

import (
	cfg "cardamom/core/config"
	"cardamom/core/ext/jwt_ext"
	"crypto/sha512"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/exp/slices"
)

type UserRole string

const (
	USER  UserRole = "user"
	ADMIN UserRole = "admin"
)

type User struct {
	Uid            string    `gorm:"primaryKey;not null;default:null" json:"uid"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Role           UserRole  `json:"role"`
	Email          string    `gorm:"unique;index" json:"email"`
	Password       []byte    `json:"-"`
	GithubToken    *string   `gorm:"default:null" json:"-"`
	GoogleToken    *string   `gorm:"default:null" json:"-"`
	FacebookToken  *string   `gorm:"default:null" json:"-"`
	MicrosoftToken *string   `gorm:"default:null" json:"-"`
}

func getSalt() []byte {
	return []byte(cfg.C.PasswordSalt)
}

func HashPassword(password string) []byte {
	// Iteration count "should" double ever two years
	// https://en.wikipedia.org/wiki/PBKDF2#:~:text=In%202021%2C%20OWASP%20recommended%20to,for%20PBKDF2%2DHMAC%2DSHA512.
	return pbkdf2.Key([]byte(password), getSalt(), 120000, 64, sha512.New)
}

func (u *User) IsPasswordMatch(password string) bool {
	return slices.Equal(u.Password, HashPassword(password))
}

func (u *User) SetPassword(password string) {
	u.Password = HashPassword(password)
}

func GetUserByEmail(email string) (*User, error) {
	user := User{}
	if db := DB.Where(&User{Email: email}).Limit(1).First(&user); db.Error != nil {
		return nil, db.Error
	} else if db.RowsAffected > 0 {
		return &user, nil
	} else {
		return nil, nil
	}
}

func (u User) GetAccessToken() (string, string, error) {
	csrf := jwt_ext.GetCSRFToken()
	token, err := jwt_ext.CreateJWT(AuthToken{
		BaseClaims: jwt.RegisteredClaims{
			Issuer:    cfg.C.Domain,
			Subject:   "access_token",
			Audience:  []string{cfg.C.Domain},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
		Uid:   u.Uid,
		Role:  u.Role,
		Email: u.Email,
		CSRF:  csrf,
	})
	return token, csrf, err
}

func (u User) GetRefreshToken() (string, string, error) {
	csrf := jwt_ext.GetCSRFToken()
	token, err := jwt_ext.CreateJWT(AuthToken{
		BaseClaims: jwt.RegisteredClaims{
			Issuer:    cfg.C.Domain,
			Subject:   "refresh_token",
			Audience:  []string{cfg.C.Domain},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 30)),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
		Uid:   u.Uid,
		Role:  u.Role,
		Email: u.Email,
		CSRF:  csrf,
	})
	return token, csrf, err
}

type AuthToken struct {
	BaseClaims jwt.RegisteredClaims
	Uid        string   `json:"uid"`
	Role       UserRole `json:"role"`
	Email      string   `json:"email"`
	CSRF       string   `json:"csrf"`
}

func (t AuthToken) Valid() error {
	return t.BaseClaims.Valid()
}
