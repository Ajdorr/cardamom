package auth

import (
	cfg "cardamom/core/config"
	"cardamom/core/models"
	"context"
	"fmt"
	"sync"

	"github.com/parnurzeal/gorequest"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

var oa2Cfg_Facebook = oauth2Config_Facebook{}

type oauth2Config_Facebook struct {
	once sync.Once
	cfg  *oauth2.Config
}

func (c *oauth2Config_Facebook) get() *oauth2.Config {
	c.once.Do(func() {
		c.cfg = &oauth2.Config{
			ClientID:     cfg.C.OAuthFacebookClientId,
			ClientSecret: cfg.C.OAuthFacebookClientSecret,
			Scopes:       []string{"email"},
			Endpoint:     facebook.Endpoint,
			RedirectURL:  "http://localhost:3000/auth/oauth-return/facebook",
		}
	})
	return c.cfg
}

type oauth2FacebookEmailResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func completeOAuth2Facebook(code string) (*models.User, []error) {

	token, err := oa2Cfg_Facebook.get().Exchange(context.Background(), code)
	if err != nil {
		return nil, []error{err}
	}

	url := "https://graph.facebook.com/v15.0/me?locale=en_US&fields=name,email&access_token=" + token.AccessToken
	var body oauth2FacebookEmailResponse
	rsp, bodyRaw, errs := gorequest.New().Get(url).EndStruct(&body)
	if len(errs) > 0 || rsp.StatusCode < 200 || rsp.StatusCode >= 300 {
		errs = append(errs, fmt.Errorf("user email response body -- %s", string(bodyRaw)))
		return nil, errs
	}

	user, err := createOrGetUser(body.Email)
	if err != nil {
		return nil, []error{err}
	}

	user.FacebookToken = &token.AccessToken
	if err = models.DB.Save(user).Error; err != nil {
		return nil, []error{err}
	}

	return user, []error{}
}
