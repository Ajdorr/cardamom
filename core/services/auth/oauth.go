package auth

import (
	"cardamom/core/ext/gin_ext"
	"cardamom/core/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func StartOAuth2(provider string, cfg *oauth2.Config) func(*gin.Context) {

	return func(c *gin.Context) {
		if state, err := SetOAuthState(c.ClientIP(), provider); err != nil {
			gin_ext.ServerError(c, err)
		} else {
			c.JSON(http.StatusOK, &gin.H{"redirect_url": cfg.AuthCodeURL(state)})
		}
	}
}

func CompleteOAuth2(
	provider string, cfg *oauth2.Config,
	providerCompleteFunc func(code string) (*models.User, []error),
) func(*gin.Context, *CompleteOAuth2Request) {

	return func(c *gin.Context, req *CompleteOAuth2Request) {
		if state, err := GetOAuthState(c.ClientIP(), provider); err != nil {
			gin_ext.ServerError(c, err)
		} else if time.Now().After(state.TTL) || req.State != state.State {
			c.AbortWithStatusJSON(http.StatusNotFound, &gin.H{"error": "Token not valid"})
		} else if user, errs := providerCompleteFunc(req.Code); len(errs) > 0 {
			gin_ext.ServerErrors(c, errs)
		} else {
			sendAuthTokenResponse(c, user)
		}
	}

}
