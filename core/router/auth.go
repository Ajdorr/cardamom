package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterBody struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func register(c *gin.Context) {
	body := RegisterBody{}
	if err := c.BindJSON(body); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	c.JSON(http.StatusAccepted, &body)
}

type LoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func login(c *gin.Context) {
	body := LoginBody{}
	if err := c.BindJSON(body); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	c.JSON(http.StatusAccepted, &body)
}

var AuthBlueprint = Blueprint{
	path: "/auth",
	routes: map[string]Route{
		"/register": {
			method:  "POST",
			handler: register,
		},
		"/login": {
			method:  "POST",
			handler: login,
		},
	},
}
