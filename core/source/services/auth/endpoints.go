package auth

import (
	"cardamom/core/source/db"
	"cardamom/core/source/db/models"
	"cardamom/core/source/ext/gin_ext"
	"cardamom/core/source/ext/jwt_ext"
	"cardamom/core/source/ext/log_ext"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const CSRF_HEADER_KEY = "X-CSRF-TOKEN"

func Login(c *gin.Context, req *LoginRequest) {
	if user, err := models.GetUserByEmail(req.Email); err != nil {
		gin_ext.Abort(c, http.StatusBadRequest, err)
	} else if user == nil {
		gin_ext.Abort(c, http.StatusBadRequest, log_ext.Errorf("attempt to log in to user(%s) that does not exist", req.Email))
	} else if !user.IsPasswordMatch(req.Password) {
		gin_ext.Abort(c, http.StatusBadRequest, log_ext.Errorf("login with bad password"))
	} else {
		sendAuthTokenResponse(c, user)
	}
}

func Logout(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{Name: JWT_ACCESS_TOKEN_KEY, MaxAge: -1, Expires: time.Now()})
	http.SetCookie(c.Writer, &http.Cookie{Name: JWT_REFRESH_TOKEN_KEY, MaxAge: -1, Expires: time.Now()})
	c.JSON(http.StatusOK, &gin.H{})
}

func Refresh(c *gin.Context) {
	user := models.User{}
	// Check refresh cookie exists
	if cookie, err := c.Cookie(JWT_REFRESH_TOKEN_KEY); err != nil {
		gin_ext.Abort(c, http.StatusUnauthorized, log_ext.Errorf("refresh token validation -- %w", err))
		// Check refresh token is the correct format, parse claims
	} else if claims, err := jwt_ext.CheckJWT(cookie, &models.AuthToken{}); err != nil {
		gin_ext.Abort(c, http.StatusUnauthorized, log_ext.Errorf("refresh token validation -- %w", err))
		// Find the user in the database
	} else if err = db.DB().First(&user, "uid = ?", claims.Uid).Error; err != nil {
		gin_ext.Abort(c, http.StatusUnauthorized, log_ext.Errorf("refresh token validation -- %w", err))
		// Create the access token
	} else {
		// TODO Invalidate refresh token and set it
		sendAuthTokenResponse(c, &user)
	}
}

// Initialized by init()
var StartOAuthGithub func(*gin.Context)
var CompleteOAuthGithub func(*gin.Context, *CompleteOAuth2Request)
var StartOAuthGoogle func(*gin.Context)
var CompleteOAuthGoogle func(*gin.Context, *CompleteOAuth2Request)
var StartOAuthFacebook func(*gin.Context)
var CompleteOAuthFacebook func(*gin.Context, *CompleteOAuth2Request)
var StartOAuthMicrosoft func(*gin.Context)
var CompleteOAuthMicrosoft func(*gin.Context, *CompleteOAuth2Request)

func init() {
	StartOAuthGithub = StartOAuth("github", oa2Cfg_Github.get())
	CompleteOAuthGithub = CompleteOAuth("github", oa2Cfg_Github.get(), completeOAuth2Github)
	StartOAuthGoogle = StartOAuth("google", oa2Cfg_Google.get())
	CompleteOAuthGoogle = CompleteOAuth("google", oa2Cfg_Google.get(), completeOAuth2Google)
	StartOAuthFacebook = StartOAuth("facebook", oa2Cfg_Facebook.get())
	CompleteOAuthFacebook = CompleteOAuth("facebook", oa2Cfg_Facebook.get(), completeOAuth2Facebook)
	StartOAuthMicrosoft = StartOAuth("microsoft", oa2Cfg_Microsoft.get())
	CompleteOAuthMicrosoft = CompleteOAuth("microsoft", oa2Cfg_Microsoft.get(), completeOAuth2Microsoft)
}
