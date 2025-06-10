package t_ext

import (
	"cardamom/core/source/db"
	"cardamom/core/source/db/models"
)

const TestUserUid = "test"
const TestUserEmail string = "admin@cardamom.com"
const TestUserPassword string = "1234"

func GetTestUser() *models.User {
	return &models.User{
		Uid:   TestUserUid,
		Role:  models.USER,
		Email: TestUserEmail,
	}
}

func EnsureTestUser() {
	if err := db.DB().Where(&models.User{
		Uid:   TestUserUid,
		Email: TestUserEmail,
	}).Attrs(models.User{
		Role:     models.USER,
		Password: models.HashPassword(TestUserPassword),
	}).FirstOrCreate(&models.User{}).Error; err != nil {
		panic(err)
	}
}
