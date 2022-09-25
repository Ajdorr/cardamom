package log_ext

import (
	"cardamom/core/ext/gin_ext"
	"cardamom/core/models"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

type errorFrame struct {
	Time       string            `json:"time"`
	Endpoint   string            `json:"endpoint"`
	Method     string            `json:"method"`
	Code       int               `json:"code"`
	Level      string            `json:"level"`
	User       string            `json:"user"`
	Message    string            `json:"msg"`
	Stacktrace errors.StackTrace `json:"stacktrace"`
	Request    any               `json:"request"`
}

func (ef *errorFrame) log() {
	if data, err := json.MarshalIndent(*ef, "", "  "); err != nil {
		fmt.Fprintf(os.Stderr, "[%s] %s %d %s %s: %s --- unable to marshal error to json",
			ef.Time, ef.Level, ef.Code, ef.Method, ef.Endpoint, ef.Message)
	} else {
		fmt.Fprintf(os.Stderr, "[%s] %s %d %s %s: %s",
			ef.Time, ef.Level, ef.Code, ef.Method, ef.Endpoint, data)
	}
}

func ErrorHandler(c *gin.Context) {
	c.Next()

	statusCode := c.Writer.Status()
	errorLevel := "Warning"
	if c.Writer.Status() >= 500 {
		errorLevel = "Error"
	}
	obj, _ := c.Get(gin_ext.REQUEST_OBJ)

	userID := "Anonymous"
	if userClaims := gin_ext.Get[models.AuthToken](c, gin_ext.JWT_ACCESS_CLAIMS_KEY); userClaims != nil {
		userID = fmt.Sprint(userClaims.Uid)
	}

	for _, err := range c.Errors {
		ef := errorFrame{
			Time:       time.Now().Format(time.RFC3339),
			Method:     c.Request.Method,
			Endpoint:   c.Request.URL.String(),
			Code:       statusCode,
			User:       userID,
			Level:      errorLevel,
			Message:    err.Error(),
			Request:    obj,
			Stacktrace: nil,
		}
		if st, ok := err.Err.(stackTracer); ok {
			ef.Stacktrace = st.StackTrace()
		}

		ef.log()
	}
}

func ReturnBoth(err string) (string, error) {
	return err, fmt.Errorf(err)
}

func ReturnBothErr(err error) (string, error) {
	return err.Error(), err
}
