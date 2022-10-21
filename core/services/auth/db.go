package auth

import (
	"cardamom/core/ext/log_ext"
	"cardamom/core/models"
)

func RegisterNewUser(email string, password string) (*models.User, error) {
	if user, err := models.GetUserByEmail(email); err != nil {
		return nil, err
	} else if user != nil {
		return nil, log_ext.Errorf("user with email(%s) already exists", email)
	} else {
		newUser := &models.User{
			Role:     models.USER,
			Email:    email,
			Password: models.HashPassword(password),
		}
		return newUser, models.DB.Create(newUser).Error
	}
}
