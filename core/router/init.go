package router

import (
	cfg "cardamom/core/config"
	"cardamom/core/ext/gin_ext"

	"github.com/gin-gonic/gin"
)

var Engine *gin.Engine

func init() {

	Engine = gin.New()
	Engine.Use(gin.Logger())   // TODO replace?
	Engine.Use(gin.Recovery()) // TODO add custom logic
	Engine.Use(gin_ext.ErrorHandler)
	if cfg.IsLocal() {
		Engine.SetTrustedProxies([]string{"localhost"})
	}

	RegisterEndpoints(Engine)

}
