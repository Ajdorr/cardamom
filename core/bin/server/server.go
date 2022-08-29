package main

import (
	"cardamom/core/app"
	"cardamom/core/router"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	app.Init()

	r := gin.Default()
	router.RegisterEndpoints(r)

	r.Run(fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")))
}
