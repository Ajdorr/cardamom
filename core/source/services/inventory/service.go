package inventory

import (
	"cardamom/core/source/db"
	"cardamom/core/source/db/models"
)

func UpdateItemQuantity(itemUid string, userUid string, delta int) error {
	if item, err := ItemByUid(itemUid, userUid); err != nil {
		return err
	} else {
		if err = db.DB().Save(&item).Error; err != nil {
			return err
		} else {
			return nil
		}
	}
}

func GetInventory(userUid string) ([]models.InventoryItem, error) {
	var items []models.InventoryItem
	err := db.DB().Where(&models.InventoryItem{UserUid: userUid, InStock: true}).Find(&items).Error
	return items, err
}
