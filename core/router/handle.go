package router

import (
	"cardamom/core/ext/gin_ext"
	"cardamom/core/services/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Request[R any] interface {
	Validate() (string, error)
	*R
}

func handleRequest(handler func(c *gin.Context), isProtected bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if isProtected {
			if err := auth.IsAuthenticated(c); err != nil {
				if strings.HasPrefix(c.Request.URL.Path, "/api/auth") {
					c.Error(err)
				}
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
				return
			}
		}

		handler(c)
	}
}

func handlePost[R any, T Request[R]](handler func(c *gin.Context, r T), isProtected bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if isProtected {
			if err := auth.IsAuthenticated(c); err != nil {
				if strings.HasPrefix(c.Request.URL.Path, "/api/auth") {
					c.Error(err)
				}
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
				return
			}
		}

		rPtr := new(T)
		if err := c.ShouldBindJSON(rPtr); err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
			return
		}

		r := *rPtr
		if msg, err := r.Validate(); err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		c.Set(gin_ext.REQUEST_OBJ, *r)
		handler(c, r)
	}
}

func handleGet[R any, T Request[R]](handler func(c *gin.Context, r T), isProtected bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if isProtected {
			if err := auth.IsAuthenticated(c); err != nil {
				if strings.HasPrefix(c.Request.URL.Path, "/api/auth") {
					c.Error(err)
				}
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
				return
			}
		}

		ptr := new(T)
		r := *ptr
		if err := c.ShouldBind(r); err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		} else if msg, err := r.Validate(); err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		c.Set(gin_ext.REQUEST_OBJ, *r)
		handler(c, r)
	}
}
