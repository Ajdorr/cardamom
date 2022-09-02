package auth

import (
	"cardamom/core/ext/rand_ext"
	"cardamom/core/models"
	"fmt"
	"regexp"
	"time"
)

const (
	MIN_PASSWORD_LENGTH = 8
)

var RE_ALPHA_UPPER = regexp.MustCompile("[A-Z]")
var RE_ALPHA_LOWER = regexp.MustCompile("[a-z]")
var RE_NUMBER = regexp.MustCompile("[0-9]")
var RE_SPECIAL_CHARACTERS = regexp.MustCompile("[^a-zA-Z0-9]")

func validatePassword(password string) error {

	if len(password) < MIN_PASSWORD_LENGTH {
		return fmt.Errorf("password too short")
	} else if RE_ALPHA_UPPER.FindStringIndex(password) == nil {
		return fmt.Errorf("missing upper case characters")
	} else if RE_ALPHA_LOWER.FindStringIndex(password) == nil {
		return fmt.Errorf("missing lower case characters")
	} else if RE_NUMBER.FindStringIndex(password) == nil {
		return fmt.Errorf("missing numbers")
	} else if RE_SPECIAL_CHARACTERS.FindStringIndex(password) == nil {
		return fmt.Errorf("missing special characters")
	}

	return nil
}

func createOrGetUser(email string) (*models.User, error) {

	user := &models.User{}
	err := models.DB.Where(&models.User{Email: email}).
		Attrs(&models.User{Uid: generateUid()}).
		FirstOrCreate(&user).Error
	if err != nil {
		return nil, err
	}
	return user, err
}

func SetOAuthState(ip string, provider string) (string, error) {
	state := models.OAuthState{}
	err := models.DB.Where(&models.OAuthState{IPAddress: ip}).FirstOrCreate(&state).Error
	if err != nil {
		return "", err
	}

	state.TTL = time.Now().Add(time.Minute * 10)
	state.Provider = provider
	state.State = rand_ext.GetRandomString(8)
	if err := models.DB.Save(&state).Error; err != nil {
		return "", err
	}

	return state.State, nil
}

func GetOAuthState(ip string, provider string) (*models.OAuthState, error) {
	state := &models.OAuthState{}
	err := models.DB.Where(&models.OAuthState{
		IPAddress: ip, Provider: provider}).First(state).Error
	if err != nil {
		return nil, err
	}

	return state, nil
}
