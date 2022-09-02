package inventory

import (
	"cardamom/core/models"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

func itemByUid(itemUid string, userUid string) (*models.InventoryItem, error) {
	item := models.InventoryItem{}
	err := models.DB.Where(&models.InventoryItem{
		Uid:     itemUid,
		UserUid: userUid,
	}).First(&item).Error
	return &item, err
}

func generateUid() string {
	return gonanoid.Must()
}
