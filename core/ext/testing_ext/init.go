package testing_ext

import (
	cfg "cardamom/core/config"
	"cardamom/core/models"
)

const testUserUid = "test"

func EnsureTestUser() {
	if cfg.IsLocal() {
		if err := models.DB.Where(&models.User{
			Uid:   testUserUid,
			Email: cfg.C.TestUserEmail,
		}).Attrs(models.User{
			Role:     models.USER,
			Password: models.HashPassword(cfg.C.TestUserPassword),
		}).FirstOrCreate(&models.User{}).Error; err != nil {
			panic(err)
		}
	}
}
