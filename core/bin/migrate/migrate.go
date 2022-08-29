package main

import (
	"cardamom/core/app"
	"cardamom/core/models"
	"os"
)

func main() {
	app.Init()
	models.Migrate()

	admin := models.User{
		FirstName: os.Getenv("ADMIN_USER_FIRST_NAME"),
		LastName:  os.Getenv("ADMIN_USER_LAST_NAME"),
		Email:     os.Getenv("ADMIN_USER_EMAIL"),
		Role:      "admin",
		Password:  models.HashPassword(os.Getenv("ADMIN_USER_PASSWORD")),
	}

	models.DB.Create(&admin)
}
