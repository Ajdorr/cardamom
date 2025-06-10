package jwt_ext

import (
	cfg "cardamom/core/source/config"
	"cardamom/core/source/ext/rand_ext"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
)

var JWT_VALID_METHODS = jwt.WithValidMethods([]string{jwt.SigningMethodHS512.Name})

func CheckJWT[T jwt.Claims](tokenStr string, template T) (T, error) {
	token, err := jwt.ParseWithClaims(tokenStr, template, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.C.Server.JwtTokenSecret), nil
	}, JWT_VALID_METHODS)
	if err != nil {
		return template, err
	}

	claims, ok := token.Claims.(T)
	if !ok || !token.Valid {
		// Done to prevent cyclical import
		return template, errors.WithStack(fmt.Errorf("invalid claims"))
	}
	return claims, nil
}

func GetCSRFToken() string {
	return rand_ext.GetRandomString(24)
}

func CreateJWT(claims jwt.Claims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString([]byte(cfg.C.Server.JwtTokenSecret))
}
