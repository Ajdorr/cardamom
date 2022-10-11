package inventory_test

import (
	t_ext "cardamom/core/ext/testing_ext"
	"cardamom/core/models"
	"fmt"
	"net/http"
	"testing"
)

var inventoryItems = []string{
	"apples", "bananas", "peaches",
	"pears", "grapes",
}

func TestCreates(t *testing.T) {

	rspBody := models.InventoryItem{}
	testCase := &t_ext.APITestCase{
		T:                    t,
		Method:               "POST",
		Endpoint:             "/api/inventory/create",
		ResponseBody:         &rspBody,
		ExpectedResponseCode: http.StatusCreated,
	}
	t_ext.AuthorizeTestCase(t, testCase)

	for _, item := range inventoryItems {
		testCase.RequestBody = fmt.Sprintf(`{"item": "%s"}`, item)
		t_ext.API_Test(testCase)
		if rspBody.Item != item {
			t.Errorf("mismatch between request and response item(%s::%s)", rspBody.Item, item)
		}
		if rspBody.UserUid != t_ext.GetTestUser().Uid {
			t.Errorf("mismatch between user uids(%s::%s)", rspBody.Uid, t_ext.GetTestUser().Uid)
		}
	}
}

func init() {
	t_ext.EnsureTestUser()
}
