package grocery

import (
	"cardamom/core/source/ext/log_ext"
	"strings"
)

type AddItemRequest struct {
	Item  string  `json:"item"`
	Store *string `json:"store,omitempty"`
}

func (req *AddItemRequest) Validate() (string, error) {

	req.Item = strings.ToLower(strings.TrimSpace(req.Item))
	if len(req.Item) == 0 {
		return log_ext.ReturnBoth("item must be non-empty")
	}

	if req.Store != nil {
		*req.Store = strings.ToLower(strings.TrimSpace(*req.Store))
		if len(*req.Store) == 0 {
			return log_ext.ReturnBoth("store must not be an empty string")
		}
	}

	return "", nil
}

type AddItemsRequest struct {
	Items []string `json:"items"`
	Store *string  `json:"store,omitempty"`
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

	if req.Store != nil {
		*req.Store = strings.ToLower(strings.TrimSpace(*req.Store))
		if len(*req.Store) == 0 {
			return log_ext.ReturnBoth("store must not be an empty string")
		}
	}

	return "", nil
}

type UpdateItemRequest struct {
	Uid   string  `json:"uid"`
	Item  *string `json:"item,omitempty"`
	Store *string `json:"store,omitempty"`
}

func (req *UpdateItemRequest) Validate() (string, error) {

	req.Uid = strings.TrimSpace(req.Uid)
	if len(req.Uid) == 0 {
		return log_ext.ReturnBoth("uid must be non-empty")
	}

	if req.Item == nil && req.Store == nil {
		return log_ext.ReturnBoth("must specify a item or store")
	}

	if req.Item != nil {
		*req.Item = strings.ToLower(strings.TrimSpace(*req.Item))
		if len(*req.Item) == 0 {
			return log_ext.ReturnBoth("item must not be empty")
		}
	}

	return "", nil
}

type CollectItemRequest struct {
	Uid         string `json:"uid"`
	IsCollected bool   `json:"is_collected"`
}

func (req *CollectItemRequest) Validate() (string, error) {

	req.Uid = strings.TrimSpace(req.Uid)
	if len(req.Uid) == 0 {
		return log_ext.ReturnBoth("uid must be non-empty")
	}

	return "", nil
}
