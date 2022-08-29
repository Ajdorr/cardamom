package app

import (
	"cardamom/core/models"

	"github.com/joho/godotenv"
)

func Init() {
	godotenv.Load()
	models.Init()
}
