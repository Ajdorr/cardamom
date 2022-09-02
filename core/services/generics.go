package services

import (
	"cardamom/core/ext/gin_ext"
	"cardamom/core/ext/log_ext"
	"cardamom/core/models"
	"cardamom/core/services/auth"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ReadRequest struct {
	Uid string `json:"uid"`
}

func (req *ReadRequest) Validate() (string, error) {

	req.Uid = strings.TrimSpace(req.Uid)
	if len(req.Uid) == 0 {
		return log_ext.ReturnBoth("uid cannot be empty")
	}
	return "", nil

}

func ListRecords[T any](c *gin.Context, getKeyModel func(token *models.AuthToken) *T) {
	user := auth.GetActiveUserClaims(c)
	modelKey := getKeyModel(user)
	var records []T
	if err := models.DB.Where(modelKey).Find(&records).Error; err != nil {
		gin_ext.ServerError(c, fmt.Errorf("listing records -- %w", err))
	} else {
		c.JSON(http.StatusOK, records)
	}
}
