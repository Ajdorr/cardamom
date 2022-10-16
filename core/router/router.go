package router

import (
	"cardamom/core/services/auth"
	"cardamom/core/services/grocery"
	"cardamom/core/services/inventory"
	"cardamom/core/services/recipe"

	"github.com/gin-gonic/gin"
)

func RegisterEndpoints(r *gin.Engine) {

	app := r.Group("/api")

	g := app.Group("/auth")
	{
		g.POST("/login", handlePost(auth.Login, false))
		g.POST("/logout", handleRequest(auth.Logout, true))
		g.POST("/refresh", handleRequest(auth.Refresh, false))
		// g.POST("/start-register", handlePost(auth.StartRegister, false))
		// g.GET("/complete-register", handleGet(auth.CompleteRegister, false))
		g.POST("/set-password", handlePost(auth.SetPassword, true))
		// g.POST("/start-reset-password", handlePost(auth.ResetPassword, true))
		// g.POST("/complete-reset-password", handleRoute(auth.ResetPassword, &auth.ResetPasswordRequest{}, true))

		g.GET("/oauth-start/github", handleRequest(auth.StartOAuthGithub, false))
		g.POST("/oauth-complete/github", handlePost(auth.CompleteOAuthGithub, false))
		g.GET("/oauth-start/google", handleRequest(auth.StartOAuthGoogle, false))
		g.POST("/oauth-complete/google", handlePost(auth.CompleteOAuthGoogle, false))
		g.GET("/oauth-start/facebook", handleRequest(auth.StartOAuthFacebook, false))
		g.POST("/oauth-complete/facebook", handlePost(auth.CompleteOAuthFacebook, false))
		g.GET("/oauth-start/microsoft", handleRequest(auth.StartOAuthMicrosoft, false))
		g.POST("/oauth-complete/microsoft", handlePost(auth.CompleteOAuthMicrosoft, false))
	}

	g = app.Group("/grocery")
	{
		g.POST("/create", handlePost(grocery.AddItem, true))
		g.POST("/list", handleRequest(grocery.ListItems, true))
		g.POST("/update", handlePost(grocery.UpdateItem, true))
		g.POST("/delete", handlePost(grocery.DeleteItem, true))
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
		g.POST("/search", handlePost(recipe.SearchRecipe, true))
		g.POST("/trash", handleRequest(recipe.ListTrashedRecipes, true))
		g.POST("/available", handlePost(recipe.GetAvailableRecipes, true))
	}

}
