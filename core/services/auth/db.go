package auth

import (
	"cardamom/core/models"
	"fmt"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

func generateUid() string {
	return gonanoid.Must(32)
}

func RegisterNewUser(email string, password string) (*models.User, error) {
	if user, err := models.GetUserByEmail(email); err != nil {
		return nil, err
	} else if user != nil {
		return nil, fmt.Errorf("user with email(%s) already exists", email)
	} else {
		newUser := &models.User{
			Uid:      generateUid(),
			Role:     models.USER,
			Email:    email,
			Password: models.HashPassword(password),
		}
		return newUser, models.DB.Create(newUser).Error
	}
}