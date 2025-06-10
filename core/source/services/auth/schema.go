package auth

import (
	"cardamom/core/source/ext/log_ext"
	"fmt"
	"net/mail"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type StartRegisterRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (req *StartRegisterRequest) Validate() (string, error) {
	req.FirstName = strings.TrimSpace(req.FirstName)
	if len(req.FirstName) == 0 {
		return log_ext.ReturnBoth("empty first name in request")
	}

	req.LastName = strings.TrimSpace(req.LastName)
	if len(req.LastName) == 0 {
		return log_ext.ReturnBoth("empty last name in request")
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if len(req.Email) == 0 {
		return log_ext.ReturnBoth("empty last name in request")
	} else if _, err := mail.ParseAddress(req.Email); err != nil {
		return log_ext.ReturnBoth(fmt.Sprintf("invalid email address %s", req.Email))
	}

	if err := validatePassword(req.Password); err != nil {
		return log_ext.ReturnBothErr(log_ext.Errorf("invalid password -- %w", err))
	}

	return "", nil
}

type RegisterToken struct {
	BaseClaims jwt.RegisteredClaims
	Email      string `json:"email"`
	Password   string `json:"password"`
}

func (t RegisterToken) Valid() error {
	return t.BaseClaims.Valid()
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (req *LoginRequest) Validate() (string, error) {

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if len(req.Email) == 0 {
		return log_ext.ReturnBoth("empty email in request")
	}

	req.Password = strings.TrimSpace(req.Password)
	if len(req.Password) == 0 {
		return log_ext.ReturnBoth("empty password in request")
	}
	return "", nil
}

type SetPasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

func (req *SetPasswordRequest) Validate() (string, error) {

	req.CurrentPassword = strings.TrimSpace(req.CurrentPassword)
	if len(req.CurrentPassword) == 0 {
		return log_ext.ReturnBoth("empty current password in request")
	}
	req.NewPassword = strings.TrimSpace(req.NewPassword)
	if len(req.NewPassword) == 0 {
		return log_ext.ReturnBoth("empty new password in request")
	} else if err := validatePassword(req.NewPassword); err != nil {
		return log_ext.ReturnBothErr(err)
	}
	return "", nil
}

type ResetPasswordRequest struct {
	Email string `json:"email"`
}

func (req *ResetPasswordRequest) Validate() (string, error) {
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if len(req.Email) == 0 {
		return log_ext.ReturnBoth("empty email in request")
	}
	return "", nil
}

type ResetPasswordToken struct {
	BaseClaims jwt.RegisteredClaims
	Email      string `json:"email"`
}

func (t ResetPasswordToken) Valid() error {
	return t.BaseClaims.Valid()
}

type CompleteOAuth2Request struct {
	State string `json:"state"`
	Code  string `json:"code"`
}

func (req *CompleteOAuth2Request) Validate() (string, error) {
	if len(req.State) == 0 {
		return log_ext.ReturnBoth("state must not be empty ")
	} else if len(req.Code) == 0 {
		return log_ext.ReturnBoth("code must not be empty ")
	}

	return "", nil
}
