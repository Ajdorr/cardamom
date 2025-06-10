package router_test

import (
	"cardamom/core/source/router"
	"cardamom/core/test/t_ext"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	t_ext.Init(t)

	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	if err != nil {
		t.Fatalf("unable build http request: %s", err)
	}

	w := httptest.NewRecorder()
	router.Engine.Handler().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
