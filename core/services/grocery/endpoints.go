package grocery

import (
	"cardamom/core/events"
	"cardamom/core/ext/gin_ext"
	"cardamom/core/ext/log_ext"
	"cardamom/core/models"
	"cardamom/core/services"
	"cardamom/core/services/auth"
	"cardamom/core/services/inventory"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddItem(c *gin.Context, r *AddItemRequest) {
	user := auth.GetActiveUserClaims(c)

	item := models.GroceryItem{}
	db := models.DB.Where(models.GroceryItem{
		Item:    r.Item,
		UserUid: user.Uid,
	}).FirstOrCreate(&item)

	if db.Error != nil {
		gin_ext.AbortNotFound(c, log_ext.Errorf("finding grocery item -- %w", db.Error))
		return
	}

	if r.Store != nil {
		item.Store = *r.Store
	}
	item.IsCollected = false
	if err := models.DB.Save(&item).Error; err != nil {
		gin_ext.ServerError(c, log_ext.Errorf("adding grocery item -- %w", db.Error))
	} else {
		c.JSON(http.StatusCreated, item)
	}

	events.Publish(&events.Event{
		Domain: "grocery",
		Type:   "add",
		Data:   map[string]string{"item": item.Item, "store": item.Store},
	})
}

func AddItems(c *gin.Context, r *AddItemsRequest) {
	user := auth.GetActiveUserClaims(c)

	items := make([]models.GroceryItem, len(r.Items))
	for i, itemValue := range r.Items {
		itemModel := &items[i]
		db := models.DB.Where(models.GroceryItem{
			Item:    itemValue,
			UserUid: user.Uid,
		}).FirstOrCreate(itemModel)

		if db.Error != nil {
			gin_ext.AbortNotFound(c, log_ext.Errorf("batch finding grocery item -- %w", db.Error))
			return
		}

		if r.Store != nil {
			itemModel.Store = *r.Store
		}
		itemModel.IsCollected = false
		if err := models.DB.Save(&itemModel).Error; err != nil {
			gin_ext.ServerError(c, log_ext.Errorf("batch adding grocery item -- %w", db.Error))
			return
		}

		events.Publish(&events.Event{
			Domain: "grocery",
			Type:   "add",
			Data:   map[string]string{"item": itemModel.Item, "store": itemModel.Store},
		})
	}

	c.JSON(http.StatusCreated, &items)
}

func ListItems(c *gin.Context) {
	user := auth.GetActiveUserClaims(c)
	var items []models.GroceryItem
	if err := models.DB.Where(&models.GroceryItem{UserUid: user.Uid}).Order("created_at desc").Find(&items).Error; err != nil {
		gin_ext.ServerError(c, log_ext.Errorf("finding grocery items -- %w", err))
	} else {
		c.JSON(http.StatusOK, items)
	}
}

func CollectItem(c *gin.Context, r *CollectItemRequest) {
	user := auth.GetActiveUserClaims(c)
	item, err := itemByUid(r.Uid, user.Uid)
	if err != nil {
		gin_ext.Abort(c, http.StatusBadRequest, log_ext.Errorf("attempt to update non existant item -- %w", err))
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
			gin_ext.ServerError(c, log_ext.Errorf("unable to update GroceryItem -- %w", err))
			return
		}
	}

	c.JSON(http.StatusOK, &item)
}

func UpdateItem(c *gin.Context, r *UpdateItemRequest) {
	user := auth.GetActiveUserClaims(c)
	item, err := itemByUid(r.Uid, user.Uid)
	if err != nil {
		gin_ext.Abort(c, http.StatusBadRequest, log_ext.Errorf("attempt to update non existant item -- %w", err))
		return
	}

	if r.Item != nil {
		item.Item = *r.Item
	}
	if r.Store != nil {
		item.Store = *r.Store
	}

	if err = models.DB.Save(&item).Error; err != nil {
		gin_ext.ServerError(c, log_ext.Errorf("unable to update GroceryItem -- %w", err))
	} else {
		c.JSON(http.StatusOK, &item)
	}

	events.Publish(&events.Event{
		Domain: "grocery",
		Type:   "update",
		Data:   map[string]string{"item": item.Item, "store": item.Store},
	})
}

func DeleteItem(c *gin.Context, r *services.ReadRequest) {
	user := auth.GetActiveUserClaims(c)
	if err := models.DB.Where(
		&models.GroceryItem{
			Uid:     r.Uid,
			UserUid: user.Uid,
		}).Delete(&models.GroceryItem{}).Error; err != nil {
		gin_ext.Abort(c, http.StatusBadRequest, log_ext.Errorf("attempt to delete non existant item -- %w", err))
		return
	}

	c.JSON(http.StatusOK, &gin.H{})
}

func ClearCollected(c *gin.Context) {
	user := auth.GetActiveUserClaims(c)
	err := models.DB.Where(
		&models.GroceryItem{
			UserUid:     user.Uid,
			IsCollected: true,
		}).Delete(&models.GroceryItem{}).Error

	if err != nil {
		gin_ext.ServerError(c, log_ext.Errorf("clearing collected -- %w", err))
	} else {
		c.JSON(http.StatusOK, gin.H{})
	}
}
