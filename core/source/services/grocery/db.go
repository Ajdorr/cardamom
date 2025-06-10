package grocery

import (
	"cardamom/core/source/db"
	"cardamom/core/source/db/models"
)

func itemByUid(itemUid string, userUid string) (*models.GroceryItem, error) {
	item := models.GroceryItem{}
	err := db.DB().Where(&models.GroceryItem{
		Uid:     itemUid,
		UserUid: userUid,
	}).First(&item).Error
	return &item, err
}
