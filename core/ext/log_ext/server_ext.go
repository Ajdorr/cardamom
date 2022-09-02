package log_ext

import (
	"cardamom/core/ext/gin_ext"
	"cardamom/core/models"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func ErrorHandler(c *gin.Context) {
	c.Next()

	statusCode := c.Writer.Status()
	errorLevel := "Warning"
	if c.Writer.Status() >= 500 {
		errorLevel = "Error"
	}

	reqObj := "<None>"
	if obj, ok := c.Get(gin_ext.REQUEST_OBJ); ok {
		if bin, err := json.Marshal(obj); err == nil {
			reqObj = string(bin)
		}
	}

	userID := "Anonymous"
	if userClaims := gin_ext.Get[models.AuthToken](c, gin_ext.JWT_ACCESS_CLAIMS_KEY); userClaims != nil {
		userID = fmt.Sprint(userClaims.Uid)
	}

	for _, err := range c.Errors {
		fmt.Fprintf(
			os.Stderr,
			"[%s] %s %s %d User(%s) %s: %v --- Request: %s\n",
			time.Now().Format(time.RFC3339), c.Request.Method, c.Request.URL.String(),
			statusCode, userID, errorLevel, err, reqObj,
		)
	}
}

func ReturnBoth(err string) (string, error) {
	return err, fmt.Errorf(err)
}

func ReturnBothErr(err error) (string, error) {
	return err.Error(), err
}
