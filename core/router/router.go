package router

import (
	"cardamom/core/ext/gin_ext"
	"cardamom/core/services/auth"
	"cardamom/core/services/grocery"
	"cardamom/core/services/inventory"
	"cardamom/core/services/recipe"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Request[R any] interface {
	Validate() (string, error)
	*R
}

func handleRequest(handler func(c *gin.Context), isProtected bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if isProtected {
			if err := auth.IsAuthenticated(c); err != nil {
				c.Error(err)
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
				return
			}
		}

		handler(c)
	}
}

func handlePost[R any, T Request[R]](handler func(c *gin.Context, r T), isProtected bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if isProtected {
			if err := auth.IsAuthenticated(c); err != nil {
				c.Error(err)
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
				return
			}
		}

		rPtr := new(T)
		if err := c.ShouldBindJSON(rPtr); err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
			return
		}

		r := *rPtr
		if msg, err := r.Validate(); err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		c.Set(gin_ext.REQUEST_OBJ, *r)
		handler(c, r)
	}
}

func handleGet[R any, T Request[R]](handler func(c *gin.Context, r T), isProtected bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if isProtected {
			if err := auth.IsAuthenticated(c); err != nil {
				gin_ext.Abort(c, http.StatusUnauthorized, err)
				return
			}
		}

		ptr := new(T)
		r := *ptr
		if err := c.ShouldBind(r); err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		} else if msg, err := r.Validate(); err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		c.Set(gin_ext.REQUEST_OBJ, *r)
		handler(c, r)
	}
}

func RegisterEndpoints(r *gin.Engine) {

	app := r.Group("/api")

	g := app.Group("/auth")
	{
		g.POST("/login", handlePost(auth.Login, false))
		g.POST("/logout", handleRequest(auth.Logout, true))
		g.POST("/refresh", handleRequest(auth.Refresh, true))
		// g.POST("/start-register", handlePost(auth.StartRegister, false))
		// g.GET("/complete-register", handleGet(auth.CompleteRegister, false))
		g.POST("/set-password", handlePost(auth.SetPassword, true))
		// g.POST("/start-reset-password", handlePost(auth.ResetPassword, true))
		// g.POST("/complete-reset-password", handleRoute(auth.ResetPassword, &auth.ResetPasswordRequest{}, true))

		g.GET("/oauth2-start/github", handleRequest(auth.StartOAuth2Github, false))
		g.POST("/oauth2-complete/github", handlePost(auth.CompleteOAuth2Github, false))
		g.GET("/oauth2-start/google", handleRequest(auth.StartOAuth2Google, false))
		g.POST("/oauth2-complete/google", handlePost(auth.CompleteOAuth2Google, false))
		g.GET("/oauth2-start/facebook", handleRequest(auth.StartOAuth2Facebook, false))
		g.POST("/oauth2-complete/facebook", handlePost(auth.CompleteOAuth2Facebook, false))
		g.GET("/oauth2-start/microsoft", handleRequest(auth.StartOAuth2Microsoft, false))
		g.POST("/oauth2-complete/microsoft", handlePost(auth.CompleteOAuth2Microsoft, false))
	}

	g = app.Group("/grocery")
	{
		g.POST("/create", handlePost(grocery.AddItem, true))
		g.POST("/list", handleRequest(grocery.ListItems, true))
		g.POST("/update", handlePost(grocery.UpdateItem, true))
		g.POST("/collect", handlePost(grocery.CollectItem, true))
		g.POST("/clear", handleRequest(grocery.ClearCollected, true))
	}

	g = app.Group("/inventory")
	{
		g.POST("/create", handlePost(inventory.AddItem, true))
		g.POST("/list", handleRequest(inventory.ListItems, true))
		g.POST("/update", handlePost(inventory.UpdateItem, true))
	}

	g = app.Group("/recipe")
	{
		g.POST("/create", handlePost(recipe.CreateRecipe, true))
		g.POST("/read", handlePost(recipe.ReadRecipe, true))
		g.POST("/update", handlePost(recipe.UpdateRecipe, true))
		g.POST("/list", handleRequest(recipe.ListRecipes, true))
		g.POST("/random-available-recipes", handleRequest(recipe.GetRandomAvailableRecipes, true))
	}

}
