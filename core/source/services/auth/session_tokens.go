package auth

import (
	cfg "cardamom/core/source/config"
	"cardamom/core/source/db"
	"cardamom/core/source/db/models"
	"cardamom/core/source/ext/gin_ext"
	"cardamom/core/source/ext/jwt_ext"
	"cardamom/core/source/ext/log_ext"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const JWT_ACCESS_TOKEN_KEY = "access_token"
const JWT_REFRESH_TOKEN_KEY = "refresh_token"

func GetActiveUser(c *gin.Context) *models.User {
	if cachedUser := gin_ext.GetKey[models.User](c, gin_ext.USER_KEY); cachedUser != nil {
		return cachedUser
	}

	user := models.User{}
	claims := GetActiveUserClaims(c)
	if err := db.DB().First(&user, claims.Uid).Error; err != nil {
		panic(log_ext.Errorf("unable to get active user from database -- %w", err))
	} else {
		c.Set(gin_ext.USER_KEY, &user)
		return &user
	}
}

func GetActiveUserClaims(c *gin.Context) *models.AuthToken {
	if claims := gin_ext.GetKey[models.AuthToken](c, gin_ext.JWT_ACCESS_CLAIMS_KEY); claims != nil {
		return claims
	} else {
		panic(log_ext.Errorf("missing active user claims"))
	}
}

func IsAuthenticated(c *gin.Context) error {
	if tokenCookie, err := c.Request.Cookie(JWT_ACCESS_TOKEN_KEY); err != nil {
		return log_ext.Errorf("unable to retrieve access token -- %w", err)
	} else if claims, err := jwt_ext.CheckJWT(tokenCookie.Value, &models.AuthToken{}); err != nil {
		return log_ext.Errorf("unable to parse token -- %w", err)
	} else if claims == nil {
		return log_ext.Errorf("empty claims")
	} else if claims.CSRF != c.GetHeader(CSRF_HEADER_KEY) {
		return log_ext.Errorf("%s header token mismatch", CSRF_HEADER_KEY)
	} else {
		c.Set(gin_ext.JWT_ACCESS_CLAIMS_KEY, claims)
		return nil
	}
}

func setAccessToken(c *gin.Context, token, csrf string, isSecure bool) {
	// FIXME domain should be an environment variable
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     JWT_ACCESS_TOKEN_KEY,
		Value:    token,
		Expires:  time.Now().Add(time.Hour),
		Path:     "/api",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   isSecure,
		// Domain:   cfg.C.Server.Domain,
	})
	c.Header(CSRF_HEADER_KEY, csrf)
}

func setRefreshToken(c *gin.Context, token string, isSecure bool) {
	// FIXME domain should be an environment variable
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     JWT_REFRESH_TOKEN_KEY,
		Value:    token,
		Path:     "/api/auth/refresh",
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   isSecure,
		// Domain:   cfg.C.Server.Domain,
	})
}

func sendAuthTokenResponse(c *gin.Context, user *models.User) {
	isSecure := !cfg.IsLocal()
	if accessToken, csrf, err := user.GetAccessToken(); err != nil {
		gin_ext.ServerError(c, log_ext.Errorf("creating access JWT -- %w", err))
	} else if refresh_token, _, err := user.GetRefreshToken(); err != nil {
		gin_ext.ServerError(c, log_ext.Errorf("creating refresh JWT -- %w", err))
	} else {
		setAccessToken(c, accessToken, csrf, isSecure)
		setRefreshToken(c, refresh_token, isSecure)
		c.JSON(http.StatusOK, &gin.H{})
	}
}
