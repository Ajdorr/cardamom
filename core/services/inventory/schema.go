package inventory

import (
	"cardamom/core/ext/log_ext"
	"strings"
)

type AddItemRequest struct {
	Item string `json:"item"`
}

func (req *AddItemRequest) Validate() (string, error) {

	req.Item = strings.ToLower(strings.TrimSpace(req.Item))
	if len(req.Item) == 0 {
		return log_ext.ReturnBoth("item must not be empty")
	}

	return "", nil
}

type UpdateItemRequest struct {
	Uid     string  `json:"uid"`
	Item    *string `json:"item,omitempty"`
	InStock *bool   `json:"in_stock,omitempty"`
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

	return "", nil
}