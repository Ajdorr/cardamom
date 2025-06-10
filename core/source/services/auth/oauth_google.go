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
)

var oa2Cfg_Google = oauth2Config_Google{}

type oauth2Config_Google struct {
	once sync.Once
	cfg  *oauth2.Config
}

func (c *oauth2Config_Google) get() *oauth2.Config {
	c.once.Do(func() {
		c.cfg = &oauth2.Config{
			ClientID:     cfg.C.OAuth2.Google.ClientId,
			ClientSecret: cfg.C.OAuth2.Google.ClientSecret,
			RedirectURL:  cfg.C.OAuth2.Google.RedirectURIs[0],
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "openid", "profile", "email"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  cfg.C.OAuth2.Google.AuthURI,
				TokenURL: cfg.C.OAuth2.Google.TokenURI,
			},
		}
	})

	return c.cfg
}

type googleUserResponse struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	FirstName     string `json:"given_name"`
	LastName      string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

func completeOAuth2Google(code string) (*models.User, []error) {

	token, err := oa2Cfg_Google.get().Exchange(context.Background(), code)
	if err != nil {
		return nil, []error{err}
	}

	var body googleUserResponse
	rsp, bodyRaw, errs := gorequest.New().Get("https://www.googleapis.com/oauth2/v2/userinfo").
		Set("Authorization", "Bearer "+token.AccessToken).
		EndStruct(&body)
	if len(errs) > 0 || rsp.StatusCode < 200 || rsp.StatusCode >= 300 {
		errs = append(errs, log_ext.Errorf("user email response body -- %s", string(bodyRaw)))
		return nil, errs
	}

	user, err := createOrGetUser(body.Email)
	if err != nil {
		return nil, []error{err}
	}

	user.GoogleToken = &token.AccessToken
	if err = db.DB().Save(user).Error; err != nil {
		return nil, []error{err}
	}

	return user, []error{}

}
