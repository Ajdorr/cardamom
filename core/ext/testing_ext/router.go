package testing_ext

import (
	"bytes"
	"cardamom/core/router"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type APITestCase struct {
	T                    *testing.T
	Method               string
	Endpoint             string
	Headers              map[string]string
	Cookies              []*http.Cookie
	RequestBody          string
	ResponseBody         any
	ResponseHttp         *http.Response
	ExpectedResponseCode int
}

func API_Test(testCase *APITestCase) *http.Response {

	data := bytes.NewReader([]byte(testCase.RequestBody))
	req, err := http.NewRequest(testCase.Method, testCase.Endpoint, data)
	if err != nil {
		testCase.T.Fatalf("unable build http request: %s", err)
	}
	for k, v := range testCase.Headers {
		req.Header.Add(k, v)
	}
	for _, c := range testCase.Cookies {
		req.AddCookie(c)
	}

	w := httptest.NewRecorder()
	router.Engine.Handler().ServeHTTP(w, req)

	if testCase.ExpectedResponseCode != w.Code {
		testCase.T.Fatalf("unexpected http response code (expected: %d, found: %d)", testCase.ExpectedResponseCode, w.Code)
	}

	if err := json.Unmarshal(w.Body.Bytes(), testCase.ResponseBody); err != nil {
		testCase.T.Fatalf("unable to unmarshal response(%s)", w.Body.Bytes())
	}

	return w.Result()
}
