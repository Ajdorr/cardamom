package client_logging

import (
	"github.com/gin-gonic/gin"
)

type ErrorRequest struct {
	LogLevel string `json:"level"`
	Msg      string `json:"msg"`
	Data     *gin.H `json:"data"`
}

func (req *ErrorRequest) Validate() (string, error) {
	if req.LogLevel != "Warning" {
		req.LogLevel = "Error"
	}
	return "", nil
}

func (req *ErrorRequest) Level() string {
	return req.LogLevel
}
func (req *ErrorRequest) Message() string {
	return req.Msg
}
func (req *ErrorRequest) Payload() any {
	return req.Data
}
