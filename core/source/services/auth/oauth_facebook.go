package auth

import (
	cfg "cardamom/core/source/config"
	"cardamom/core/source/db"
	"cardamom/core/source/db/models"
	"cardamom/core/source/ext/log_ext"
	"context"
	"sync"

	"github.com/parnurzeal/gorequest"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

var oa2Cfg_Facebook = oauthConfig_Facebook{}

type oauthConfig_Facebook struct {
	once sync.Once
	cfg  *oauth2.Config
}

func (c *oauthConfig_Facebook) get() *oauth2.Config {
	c.once.Do(func() {
		c.cfg = &oauth2.Config{
			ClientID:     cfg.C.OAuth2.Facebook.ClientId,
			ClientSecret: cfg.C.OAuth2.Facebook.ClientSecret,
			Scopes:       []string{"email"},
			Endpoint:     facebook.Endpoint,
			RedirectURL:  getOAuthRedirectURL("facebook"),
		}
	})
	return c.cfg
}

type facebookEmailResponse struct {
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
	var body facebookEmailResponse
	rsp, bodyRaw, errs := gorequest.New().Get(url).EndStruct(&body)
	if len(errs) > 0 || rsp.StatusCode < 200 || rsp.StatusCode >= 300 {
		errs = append(errs, log_ext.Errorf("user email response body -- %s", string(bodyRaw)))
		return nil, errs
	}

	user, err := createOrGetUser(body.Email)
	if err != nil {
		return nil, []error{err}
	}

	user.FacebookToken = &token.AccessToken
	if err = db.DB().Save(user).Error; err != nil {
		return nil, []error{err}
	}

	return user, []error{}
}
