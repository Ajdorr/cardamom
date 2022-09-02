package jwt_ext

import (
	cfg "cardamom/core/config"
	"cardamom/core/ext/rand_ext"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var JWT_VALID_METHODS = jwt.WithValidMethods([]string{jwt.SigningMethodHS512.Name})

func CheckJWT[T jwt.Claims](tokenStr string, template T) (T, error) {
	token, err := jwt.ParseWithClaims(tokenStr, template, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.C.JwtTokenSecret), nil
	}, JWT_VALID_METHODS)
	if err != nil {
		return template, err
	}

	claims, ok := token.Claims.(T)
	if !ok || !token.Valid {
		return template, fmt.Errorf("invalid claims")
	}
	return claims, nil
}

func GetBaseClaims(subject string, expiry time.Time) jwt.RegisteredClaims {
	return jwt.RegisteredClaims{
		Issuer:    cfg.C.Domain,
		Subject:   subject,
		Audience:  []string{cfg.C.Domain},
		ExpiresAt: jwt.NewNumericDate(expiry),
		NotBefore: jwt.NewNumericDate(time.Now()),
	}
}

func GetCSRFToken() string {
	return rand_ext.GetRandomString(24)
}

func CreateJWT(claims jwt.Claims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString([]byte(cfg.C.JwtTokenSecret))
}
