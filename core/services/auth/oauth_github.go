package auth

import (
	cfg "cardamom/core/config"
	"cardamom/core/ext/log_ext"
	"cardamom/core/models"
	"context"
	"sync"

	"github.com/parnurzeal/gorequest"
	"github.com/thoas/go-funk"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var oa2Cfg_Github = oauthConfig_Github{}

type oauthConfig_Github struct {
	once sync.Once
	cfg  *oauth2.Config
}

func (c *oauthConfig_Github) get() *oauth2.Config {
	c.once.Do(func() {
		c.cfg = &oauth2.Config{
			ClientID:     cfg.C.OAuthGithubClientId,
			ClientSecret: cfg.C.OAuthGithubClientSecret,
			Scopes:       []string{"user:email"},
			Endpoint:     github.Endpoint,
			RedirectURL:  getOAuthRedirectURL("github"),
		}
	})
	return c.cfg
}

type githubEmailResponse struct {
	Email      string `json:"email"`
	Primary    bool   `json:"primary"`
	Verified   bool   `json:"verified"`
	Visibility string `json:"visibility"`
}

func completeOAuth2Github(code string) (*models.User, []error) {

	token, err := oa2Cfg_Github.get().Exchange(context.Background(), code)
	if err != nil {
		return nil, []error{err}
	}

	var body []githubEmailResponse
	rsp, bodyRaw, errs := gorequest.New().Get("https://api.github.com/user/emails").
		Set("Authorization", "Bearer "+token.AccessToken).
		EndStruct(&body)
	if len(errs) > 0 || rsp.StatusCode < 200 || rsp.StatusCode >= 300 {
		errs = append(errs, log_ext.Errorf("user email response body -- %s", string(bodyRaw)))
		return nil, errs
	}

	emailRsp := funk.Filter(body,
		func(b githubEmailResponse) bool { return b.Primary }).([]githubEmailResponse)[0]
	user, err := createOrGetUser(emailRsp.Email)
	if err != nil {
		return nil, []error{err}
	}

	user.GithubToken = &token.AccessToken
	if err = models.DB.Save(user).Error; err != nil {
		return nil, []error{err}
	}

	return user, []error{}
}
