package auth_test

import (
	"cardamom/core/source/services/auth"
	"cardamom/core/test/t_ext"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestLogin(t *testing.T) {
	t_ext.Init(t)

	reqBody := fmt.Sprintf(`{"email":"%s","password":"%s"}`, t_ext.TestUserEmail, t_ext.TestUserPassword)
	testCase := &t_ext.APITestCase{
		T:                    t,
		Method:               "POST",
		Endpoint:             "/api/auth/login",
		RequestBody:          reqBody,
		ResponseBody:         &gin.H{},
		ExpectedResponseCode: http.StatusOK,
	}
	t_ext.AuthorizeTestCase(t, testCase)
	rsp := t_ext.API_Test(testCase)

	cookies := rsp.Cookies()
	if len(cookies) != 2 {
		t.Errorf("incorrect number of cookies returned on login(%d)", len(cookies))
	} else if cookies[0].Name != "access_token" {
		t.Errorf("missing access token: %s", cookies[0].Name)
	} else if cookies[1].Name != "refresh_token" {
		t.Errorf("missing refresh token: %s", cookies[1].Name)
	}

	if len(rsp.Header.Get("X-CSRF-TOKEN")) == 0 {
		t.Errorf("no csrf returned on login")
	}
}
func TestRefesh(t *testing.T) {
	t_ext.Init(t)

	testCase := &t_ext.APITestCase{
		T:                    t,
		Method:               "POST",
		Endpoint:             "/api/auth/refresh",
		RequestBody:          "{}",
		ResponseBody:         &gin.H{},
		ExpectedResponseCode: http.StatusOK,
	}
	user := t_ext.GetTestUser()
	token, _, err := user.GetRefreshToken()
	if err != nil {
		t.Fatalf("unable to get user refresh tokens test case: %s", err)
	}
	testCase.Cookies = []*http.Cookie{
		{
			Name:    auth.JWT_REFRESH_TOKEN_KEY,
			Value:   token,
			Expires: time.Now().Add(time.Hour * 1),
			Path:    "/api/auth/refresh",
		},
	}
	rsp := t_ext.API_Test(testCase)

	cookies := rsp.Cookies()
	if len(cookies) != 2 {
		t.Errorf("incorrect number of cookies returned on login(%d)", len(cookies))
	} else {
		if cookies[0].Name != "access_token" {
			t.Errorf("missing access token: %s", cookies[0].Name)
		}
		if cookies[1].Name != "refresh_token" {
			t.Errorf("missing refresh token: %s", cookies[0].Name)
		}
	}

	if len(rsp.Header.Get("X-CSRF-TOKEN")) == 0 {
		t.Errorf("no csrf returned on login")
	}
}

func TestLogout(t *testing.T) {
	t_ext.Init(t)

	testCase := &t_ext.APITestCase{
		T:                    t,
		Method:               "POST",
		Endpoint:             "/api/auth/logout",
		RequestBody:          "{}",
		ResponseBody:         &gin.H{},
		ExpectedResponseCode: http.StatusOK,
	}
	t_ext.AuthorizeTestCase(t, testCase)
	t_ext.API_Test(testCase)
}
