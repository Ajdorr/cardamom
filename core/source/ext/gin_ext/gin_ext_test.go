package gin_ext

import (
	"bytes"
	gin_mock "cardamom/core/test/mock/gin"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestServerErrors(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	ServerErrors(c, append(
		make([]error, 0),
		errors.Errorf(""),
		errors.Errorf(""),
	))
}

func TestServerError(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	ServerError(c, errors.Errorf(""))
}

func TestAbortNotFound(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	AbortNotFound(c, errors.Errorf(""))
}

func TestGetKey(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Keys = make(map[string]any)
	c.Keys["test"] = "test"
	v := GetKey[int](c, "test")
	assert.Nil(t, v)
}

func TestLog(t *testing.T) {
	ef := errorFrame{}
	ef.log()
	ef.Request = func() {}
	ef.log()
}

func TestErrorHandler(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	responseWriterMock := gin_mock.NewMockResponseWriter(t)
	responseWriterMock.Mock.On("Status").Return(500)
	c.Writer = responseWriterMock
	c.Errors = append(make([]*gin.Error, 0), &gin.Error{Err: errors.Errorf("")})
	c.Request = httptest.NewRequest("POST", "/url", io.NopCloser(bytes.NewReader(nil)))
	ErrorHandler(c)
}
