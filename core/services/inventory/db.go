package inventory

import (
	"cardamom/core/models"
)

func ItemByUid(itemUid string, userUid string) (*models.InventoryItem, error) {
	item := models.InventoryItem{}
	err := models.DB.Where(&models.InventoryItem{
		Uid:     itemUid,
		UserUid: userUid,
	}).First(&item).Error
	return &item, err
}

func ItemByValue(itemValue, userUid string) (*models.InventoryItem, error) {
	item := models.InventoryItem{}
	err := models.DB.Where(&models.InventoryItem{
		Item:    itemValue,
		UserUid: userUid,
	}).First(&item).Error
	return &item, err
}

func FindOrCreateItem(item, userUid string) (*models.InventoryItem, error) {
	inventoryItem := &models.InventoryItem{}
	if err := models.DB.Where(&models.InventoryItem{
		Item:    item,
		UserUid: userUid,
	}).FirstOrCreate(inventoryItem).Error; err != nil {
		return nil, err
	}

	return inventoryItem, nil
}
