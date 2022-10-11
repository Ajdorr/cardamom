package client_logging

import (
	"cardamom/core/ext/gin_ext"
	"cardamom/core/services/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddItem(c *gin.Context, r *ErrorRequest) {
	gin_ext.LogClient(c, auth.GetActiveUserClaims(c), r)
	c.JSON(http.StatusOK, &gin.H{})
}
