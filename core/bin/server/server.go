package main

import (
	cfg "cardamom/core/config"
	"cardamom/core/ext/log_ext"
	"cardamom/core/router"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.New()
	r.Use(gin.Logger())   // TODO replace?
	r.Use(gin.Recovery()) // TODO add custom logic
	r.Use(log_ext.ErrorHandler)
	router.RegisterEndpoints(r)

	r.Run(fmt.Sprintf("%s:%s", cfg.C.Host, cfg.C.Port))
}
