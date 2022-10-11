package main

import (
	_ "cardamom/core/config"
	t_ext "cardamom/core/ext/testing_ext"
	"cardamom/core/models"
)

func main() {
	models.AutoMigrate()
	t_ext.EnsureTestUser()
}
