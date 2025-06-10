package inventory

import (
	"cardamom/core/source/db/models"
	"cardamom/core/source/ext/log_ext"
	"strings"

	"github.com/thoas/go-funk"
)

type AddItemRequest struct {
	Item     string                    `json:"item"`
	Category *models.InventoryCategory `json:"category"`
}

func (req *AddItemRequest) Validate() (string, error) {

	req.Item = strings.ToLower(strings.TrimSpace(req.Item))
	if len(req.Item) == 0 {
		return log_ext.ReturnBoth("item must not be empty")
	}

	if req.Category != nil && !funk.Contains(models.ValidCategories, *req.Category) {
		return log_ext.ReturnBoth("invalid category")
	}

	return "", nil
}

type AddItemsRequest struct {
	Items    []string                  `json:"items"`
	Category *models.InventoryCategory `json:"category"`
}

func (req *AddItemsRequest) Validate() (string, error) {

	if len(req.Items) == 0 {
		return log_ext.ReturnBoth("items must not be empty")
	}

	for i, item := range req.Items {
		req.Items[i] = strings.ToLower(strings.TrimSpace(item))
		if len(req.Items[i]) == 0 {
			return log_ext.ReturnBoth("item must not be empty")
		}
	}

	if req.Category != nil && !funk.Contains(models.ValidCategories, *req.Category) {
		return log_ext.ReturnBoth("invalid category")
	}

	return "", nil
}

type UpdateItemRequest struct {
	Uid      string                    `json:"uid"`
	Item     *string                   `json:"item,omitempty"`
	InStock  *bool                     `json:"in_stock,omitempty"`
	Category *models.InventoryCategory `json:"category,omitempty"`
}

func (req *UpdateItemRequest) Validate() (string, error) {

	req.Uid = strings.TrimSpace(req.Uid)
	if len(req.Uid) == 0 {
		return log_ext.ReturnBoth("uid must not be empty")
	}

	if req.Item != nil {
		*req.Item = strings.ToLower(strings.TrimSpace(*req.Item))
		if len(*req.Item) == 0 {
			return log_ext.ReturnBoth("item cannot not be empty")
		}
	}

	if req.Category != nil && !funk.Contains(models.ValidCategories, *req.Category) {
		return log_ext.ReturnBoth("invalid category")
	}

	return "", nil
}
