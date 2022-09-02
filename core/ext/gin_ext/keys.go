package gin_ext

import (
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
)

const JWT_ACCESS_CLAIMS_KEY = "jwt.access.claims"
const USER_KEY = "auth.user"

const REQUEST_OBJ = "request.obj"

func Get[T any](c *gin.Context, key string) *T {
	if v, ok := c.Get(key); !ok {
		return nil
	} else if value, ok := v.(*T); !ok {
		c.Error(fmt.Errorf(
			"found invalid type(%s) in conext with key(%s)",
			key, reflect.TypeOf(v).String()))
		return nil
	} else {
		return value
	}
}
