package main

import (
	cfg "cardamom/core/config"
	"cardamom/core/models"
)

func main() {
	models.Migrate()

	if cfg.C.Env == "local" {
		if err := models.DB.Where(&models.User{
			Email: cfg.C.AdminUserEmail,
		}).Attrs(models.User{
			Uid:      "admin",
			Role:     "admin",
			Password: models.HashPassword(cfg.C.AdminUserPassword),
		}).FirstOrCreate(&models.User{}).Error; err != nil {
			panic(err)
		}
	}

}
