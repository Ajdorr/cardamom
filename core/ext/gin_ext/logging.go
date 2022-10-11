package gin_ext

import (
	cfg "cardamom/core/config"
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
	HttpCode   int               `json:"httpCode"`
	LogLevel   string            `json:"level"`
	User       string            `json:"user"`
	Agent      string            `json:"agent"`
	IP         string            `json:"ip"`
	Message    string            `json:"msg"`
	Stacktrace errors.StackTrace `json:"stacktrace"`
	Request    any               `json:"request"`
}

func (ef *errorFrame) log() {
	indent := ""
	if cfg.IsLocal() {
		indent = "  "
	}

	if data, err := json.MarshalIndent(*ef, "", indent); err != nil {
		fmt.Fprintf(os.Stderr, "[%s] %s %d %s %s: %s --- unable to marshal error to json",
			ef.Time, ef.LogLevel, ef.HttpCode, ef.Method, ef.Endpoint, ef.Message)
	} else {
		fmt.Fprintf(os.Stderr, "[%s] %s %d %s %s: %s",
			ef.Time, ef.LogLevel, ef.HttpCode, ef.Method, ef.Endpoint, data)
	}
}

func ErrorHandler(c *gin.Context) {
	c.Next()

	statusCode := c.Writer.Status()
	errorLevel := "Warning"
	if c.Writer.Status() >= 500 {
		errorLevel = "Error"
	}
	obj, _ := c.Get(REQUEST_OBJ)

	userID := "Anonymous"
	if userClaims := GetKey[models.AuthToken](c, JWT_ACCESS_CLAIMS_KEY); userClaims != nil {
		userID = fmt.Sprint(userClaims.Uid)
	}

	for _, err := range c.Errors {
		ef := errorFrame{
			Time:       time.Now().Format(time.RFC3339),
			Method:     c.Request.Method,
			Endpoint:   c.Request.URL.String(),
			HttpCode:   statusCode,
			User:       userID,
			Agent:      c.Request.UserAgent(),
			IP:         c.RemoteIP(),
			LogLevel:   errorLevel,
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

type ClientError interface {
	Level() string
	Message() string
	Payload() any
}

type clientErrorFrame struct {
	Data     *gin.H
	Time     string `json:"time"`
	LogLevel string `json:"level"`
	User     string `json:"user"`
	Agent    string `json:"agent"`
	IP       string `json:"ip"`
	Message  string `json:"msg"`
	Payload  any    `json:"payload"`
}

func (ef *clientErrorFrame) log() {
	indent := ""
	if cfg.IsLocal() {
		indent = "  "
	}

	if data, err := json.MarshalIndent(*ef, "", indent); err != nil {
		fmt.Fprintf(os.Stderr, "[%s] %s 600 client null: %s --- unable to marshal error to json",
			ef.Time, ef.LogLevel, ef.Message)
	} else {
		fmt.Fprintf(os.Stderr, "[%s] %s 600 client null: %s",
			ef.Time, ef.LogLevel, data)
	}
}

func LogClient(c *gin.Context, claims *models.AuthToken, err ClientError) {
	ef := clientErrorFrame{
		Time:     time.Now().Format(time.RFC3339),
		LogLevel: err.Level(),
		User:     claims.Uid,
		Agent:    c.Request.UserAgent(),
		IP:       c.RemoteIP(),
		Message:  err.Message(),
		Payload:  err.Payload(),
	}
	ef.log()
}
