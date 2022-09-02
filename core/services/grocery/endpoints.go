package grocery

import (
	"cardamom/core/ext/gin_ext"
	"cardamom/core/models"
	"cardamom/core/services/auth"
	"cardamom/core/services/inventory"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddItem(c *gin.Context, r *AddItemRequest) {
	user := auth.GetActiveUserClaims(c)

	item := models.GroceryItem{}
	db := models.DB.Where(models.GroceryItem{
		Item:    r.Item,
		UserUid: user.Uid,
	}).Attrs(models.GroceryItem{Uid: generateUid()}).
		FirstOrCreate(&item)

	if db.Error != nil {
		gin_ext.AbortNotFound(c, fmt.Errorf("finding grocery item -- %w", db.Error))
		return
	}

	// If it does not already exists
	if db.RowsAffected == 1 {
		if r.Store != nil {
			item.Store = *r.Store
		}
	}
	item.IsCollected = false
	if err := models.DB.Save(&item).Error; err != nil {
		gin_ext.ServerError(c, fmt.Errorf("adding grocery item -- %w", db.Error))
	} else {
		c.JSON(http.StatusOK, item)
	}
}

func ListItems(c *gin.Context) {
	user := auth.GetActiveUserClaims(c)
	var items []models.GroceryItem
	if err := models.DB.Where(&models.GroceryItem{UserUid: user.Uid}).Order("created_at desc").Find(&items).Error; err != nil {
		gin_ext.ServerError(c, fmt.Errorf("finding grocery items -- %w", err))
	} else {
		c.JSON(http.StatusOK, items)
	}
}

func CollectItem(c *gin.Context, r *CollectItemRequest) {
	user := auth.GetActiveUserClaims(c)
	item, err := itemByUid(r.Uid, user.Uid)
	if err != nil {
		gin_ext.Abort(c, http.StatusBadRequest, fmt.Errorf("attempt to update non existant item -- %w", err))
		return
	}

	if item.IsCollected != r.IsCollected {
		if r.IsCollected {
			inventory.CollectItem(item, user.Uid, false)
		} else { // Undo
			inventory.CollectItem(item, user.Uid, true)
		}

		item.IsCollected = r.IsCollected
		if err = models.DB.Save(&item).Error; err != nil {
			gin_ext.ServerError(c, fmt.Errorf("unable to update GroceryItem -- %w", err))
			return
		}
	}

	c.JSON(http.StatusOK, &item)
}

func UpdateItem(c *gin.Context, r *UpdateItemRequest) {
	user := auth.GetActiveUserClaims(c)
	item, err := itemByUid(r.Uid, user.Uid)
	if err != nil {
		gin_ext.Abort(c, http.StatusBadRequest, fmt.Errorf("attempt to update non existant item -- %w", err))
		return
	}

	if r.Item != nil {
		item.Item = *r.Item
	}
	if r.Store != nil {
		item.Store = *r.Store
	}

	if err = models.DB.Save(&item).Error; err != nil {
		gin_ext.ServerError(c, fmt.Errorf("unable to update GroceryItem -- %w", err))
	} else {
		c.JSON(http.StatusOK, &item)
	}
}

func ClearCollected(c *gin.Context) {
	user := auth.GetActiveUserClaims(c)
	err := models.DB.Where(
		&models.GroceryItem{
			UserUid:     user.Uid,
			IsCollected: true,
		}).Delete(&models.GroceryItem{}).Error

	if err != nil {
		gin_ext.ServerError(c, fmt.Errorf("clearing collected -- %w", err))
	} else {
		c.JSON(http.StatusOK, gin.H{})
	}
}
