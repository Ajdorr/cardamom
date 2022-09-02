package inventory

import (
	"cardamom/core/ext/gin_ext"
	"cardamom/core/models"
	"cardamom/core/services/auth"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddItem(c *gin.Context, r *AddItemRequest) {
	claims := auth.GetActiveUserClaims(c)

	item := models.InventoryItem{}
	db := models.DB.Where(&models.InventoryItem{
		Item:    r.Item,
		UserUid: claims.Uid,
	}).Attrs(&models.GroceryItem{Uid: generateUid()}).
		FirstOrCreate(&item)

	if db.Error != nil {
		gin_ext.ServerError(c, fmt.Errorf("adding item to database -- %w", db.Error))
		return
	}

	item.InStock = true

	if err := models.DB.Save(&item).Error; err != nil {
		gin_ext.ServerError(c, fmt.Errorf("updating quantity to database -- %w", err))
		return
	}

	c.JSON(http.StatusOK, &item)
}

func ListItems(c *gin.Context) {
	user := auth.GetActiveUserClaims(c)
	var inventoryItems []models.InventoryItem
	err := models.DB.Where(models.InventoryItem{
		UserUid: user.Uid,
		InStock: true,
	}).Find(&inventoryItems).Error

	if err != nil {
		gin_ext.ServerError(c, fmt.Errorf("finding grocery items -- %w", err))
	} else {
		c.JSON(http.StatusOK, inventoryItems)
	}
}

func UpdateItem(c *gin.Context, r *UpdateItemRequest) {
	user := auth.GetActiveUserClaims(c)
	if item, err := itemByUid(r.Uid, user.Uid); err != nil {
		gin_ext.AbortNotFound(c, fmt.Errorf("attempt to update non existant item(%s) -- %w", r.Uid, err))
	} else {

		if r.Item != nil {
			item.Item = *r.Item
		}
		if r.InStock != nil {
			item.InStock = *r.InStock
		}

		if err = models.DB.Save(&item).Error; err != nil {
			gin_ext.ServerError(c, fmt.Errorf("unable to update GroceryItem -- %w", err))
		} else {
			c.JSON(http.StatusOK, &item)
		}
	}
}
