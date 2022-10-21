package grocery

import (
	"cardamom/core/models"
)

func itemByUid(itemUid string, userUid string) (*models.GroceryItem, error) {
	item := models.GroceryItem{}
	err := models.DB.Where(&models.GroceryItem{
		Uid:     itemUid,
		UserUid: userUid,
	}).First(&item).Error
	return &item, err
}
