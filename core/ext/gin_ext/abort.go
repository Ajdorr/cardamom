package gin_ext

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ServerErrors(c *gin.Context, errs []error) {
	for _, err := range errs {
		c.Error(err)
	}
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{})
}

func ServerError(c *gin.Context, err error) {
	Abort(c, http.StatusInternalServerError, err)
}

func AbortNotFound(c *gin.Context, err error) {
	Abort(c, http.StatusNotFound, err)
}

func Abort(c *gin.Context, code int, err error) {
	c.Error(err)
	c.AbortWithStatusJSON(code, gin.H{})
}
