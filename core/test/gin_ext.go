package lib

import "github.com/gin-gonic/gin"

//go:generate mockery --name=MockResponseWriter
type MockResponseWriter interface {
	gin.ResponseWriter
}
