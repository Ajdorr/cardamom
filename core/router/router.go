package router

import (
	"github.com/gin-gonic/gin"
)

func RegisterEndpoints(r *gin.Engine) {

	rootBlueprint := Blueprint{
		path:   "/app",
		routes: map[string]Route{},
		subroutes: []*Blueprint{
			&AuthBlueprint,
		},
	}
	RegisterBlueprint(r, &rootBlueprint)
}
