package auth

import (
	cfg "cardamom/core/config"
	"cardamom/core/ext/gin_ext"
	"cardamom/core/ext/jwt_ext"
	"cardamom/core/ext/log_ext"
	"cardamom/core/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

const CSRF_HEADER_KEY = "X-CSRF-TOKEN"

func StartRegister(c *gin.Context, req *StartRegisterRequest) {
	claims := RegisterToken{
		BaseClaims: jwt_ext.GetBaseClaims("register-token", time.Now().Add(time.Hour*24)),
		Email:      req.Email,
		Password:   req.Password,
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString(
		[]byte(cfg.C.JwtTokenSecret))
	if err != nil {
		gin_ext.ServerError(c, err)
		return
	}

	// TODO Email the token
	fmt.Print("Register token: " + token)
	c.JSON(http.StatusOK, &gin.H{})
}

func CompleteRegister(c *gin.Context, req *CompleteRegisterRequest) {
	if claims, err := jwt_ext.CheckJWT(req.RegisterToken, &RegisterToken{}); err != nil {
		gin_ext.ServerError(c, log_ext.Errorf("invaid jwt -- %w", err))

	} else if user, err := RegisterNewUser(claims.Email, claims.Password); err != nil {
		gin_ext.ServerError(c, log_ext.Errorf("unable to create user -- %w", err))

	} else if user != nil {
		gin_ext.ServerError(c, log_ext.Errorf("user with email(%s) already exists", user.Email))

	} else {
		c.JSON(http.StatusCreated, &gin.H{})
	}
}

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
	} else if err = models.DB.First(&user, "uid = ?", claims.Uid).Error; err != nil {
		gin_ext.Abort(c, http.StatusUnauthorized, log_ext.Errorf("refresh token validation -- %w", err))
		// Create the access token
	} else {
		// TODO Invalidate refresh token and set it
		sendAuthTokenResponse(c, &user)
	}
}

func SetPassword(c *gin.Context, req *SetPasswordRequest) {
	user := GetActiveUser(c)
	if !user.IsPasswordMatch(req.CurrentPassword) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}

	user.SetPassword(req.NewPassword)
	c.JSON(http.StatusOK, &gin.H{})
}

func ResetPassword(c *gin.Context, req *ResetPasswordRequest) {

	claims := ResetPasswordToken{
		BaseClaims: jwt_ext.GetBaseClaims("reset-password-token", time.Now().Add(time.Minute*5)),
		Email:      req.Email,
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString(
		[]byte(cfg.C.JwtTokenSecret))
	if err != nil {
		gin_ext.ServerError(c, err)
		return
	}

	// TODO Email the token
	fmt.Print("Reset password token: " + token)
	c.JSON(http.StatusOK, &gin.H{})
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
