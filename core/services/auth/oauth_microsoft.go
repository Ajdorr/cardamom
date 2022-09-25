package auth

import (
	cfg "cardamom/core/config"
	"cardamom/core/models"
	"context"
	"fmt"
	"sync"

	"github.com/parnurzeal/gorequest"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

var oa2Cfg_Microsoft = oauthConfig_Microsoft{}

type oauthConfig_Microsoft struct {
	once sync.Once
	cfg  *oauth2.Config
}

func (c *oauthConfig_Microsoft) get() *oauth2.Config {
	c.once.Do(func() {
		c.cfg = &oauth2.Config{
			ClientID:     cfg.C.OAuthMicrosoftClientId,
			ClientSecret: cfg.C.OAuthMicrosoftClientSecret,
			Scopes:       []string{"User.Read"},
			Endpoint:     microsoft.AzureADEndpoint("common"),
			RedirectURL:  getOAuthRedirectURL("microsoft"),
		}
	})
	return c.cfg
}

type microsoftEmailResponse struct {
	ID        string `json:"id"`
	Name      string `json:"displayName"`
	FirstName string `json:"givenName"`
	LastName  string `json:"surname"`
	Email     string `json:"userPrincipalName"`
}

func completeOAuth2Microsoft(code string) (*models.User, []error) {

	token, err := oa2Cfg_Microsoft.get().Exchange(context.Background(), code)
	if err != nil {
		return nil, []error{err}
	}

	var body microsoftEmailResponse
	rsp, bodyRaw, errs := gorequest.New().Get("https://graph.microsoft.com/v1.0/me").
		Set("Authorization", "Bearer "+token.AccessToken).
		EndStruct(&body)
	if len(errs) > 0 || rsp.StatusCode < 200 || rsp.StatusCode >= 300 {
		errs = append(errs, fmt.Errorf("user email response body -- %s", string(bodyRaw)))
		return nil, errs
	}

	user, err := createOrGetUser(body.Email)
	if err != nil {
		return nil, []error{err}
	}

	user.MicrosoftToken = &token.AccessToken
	if err = models.DB.Save(user).Error; err != nil {
		return nil, []error{err}
	}

	return user, []error{}
}
