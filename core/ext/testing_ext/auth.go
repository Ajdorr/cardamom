package testing_ext

import (
	cfg "cardamom/core/config"
	"cardamom/core/models"
	"cardamom/core/services/auth"
	"net/http"
	"testing"
	"time"
)

func GetTestUser() *models.User {
	return &models.User{
		Uid:   testUserUid,
		Role:  models.USER,
		Email: cfg.C.TestUserEmail,
	}
}

func AuthorizeTestCase(t *testing.T, testCase *APITestCase) {

	user := GetTestUser()
	token, csrf, err := user.GetAccessToken()
	if err != nil {
		t.Fatalf("unable to authorize test case: %s", err)
	}

	testCase.Cookies = []*http.Cookie{
		{
			Name:    auth.JWT_ACCESS_TOKEN_KEY,
			Value:   token,
			Expires: time.Now().Add(time.Hour * 1),
			Path:    "/api",
		},
	}
	testCase.Headers = map[string]string{"X-CSRF-TOKEN": csrf}
}
