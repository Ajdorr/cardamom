package inventory

import "cardamom/core/models"

func UpdateItemQuantity(itemUid string, userUid string, delta int) error {
	if item, err := itemByUid(itemUid, userUid); err != nil {
		return err
	} else {
		if err = models.DB.Save(&item).Error; err != nil {
			return err
		} else {
			return nil
		}
	}
}

func CollectItem(groceryItem *models.GroceryItem, userUid string, isUndo bool) error {
	item := models.InventoryItem{}
	db := models.DB.Where(&models.InventoryItem{
		Item:    groceryItem.Item,
		UserUid: userUid,
	}).Attrs(&models.GroceryItem{Uid: generateUid()}).
		FirstOrCreate(&item)

	if db.Error != nil {
		return db.Error
	}

	item.InStock = !isUndo
	if err := models.DB.Save(&item).Error; err != nil {
		return err
	} else {
		return nil
	}
}

func GetInventory(userUid string) ([]models.InventoryItem, error) {
	var items []models.InventoryItem
	err := models.DB.Where(&models.InventoryItem{UserUid: userUid}).Find(&items).Error
	return items, err
}
