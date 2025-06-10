package t_ext

import (
	"cardamom/core/source/services/auth"
	"net/http"
	"testing"
	"time"
)

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
