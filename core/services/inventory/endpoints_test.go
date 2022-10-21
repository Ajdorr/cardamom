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
		if rspBody.Category != models.COOKING {
			t.Errorf("bad category returned: %s", rspBody.Category)
		}
		if rspBody.UserUid != t_ext.GetTestUser().Uid {
			t.Errorf("mismatch between user uids(%s::%s)", rspBody.Uid, t_ext.GetTestUser().Uid)
		}
	}
}

func TestCategoryCreates(t *testing.T) {

	rspBody := models.InventoryItem{}
	testCase := &t_ext.APITestCase{
		T:                    t,
		Method:               "POST",
		Endpoint:             "/api/inventory/create",
		ResponseBody:         &rspBody,
		ExpectedResponseCode: http.StatusCreated,
	}
	t_ext.AuthorizeTestCase(t, testCase)

	testCase.RequestBody = `{"item": "noodles", "category": "cooking"}`
	t_ext.API_Test(testCase)
	if rspBody.Item != "noodles" {
		t.Errorf("mismatch between request and response item(%s::noodles)", rspBody.Item)
	}
	if rspBody.Category != "cooking" {
		t.Errorf("mismatch between request and response category(%s::cooking)", rspBody.Category)
	}

	testCase.RequestBody = `{"item": "paprika", "category": "spices"}`
	t_ext.API_Test(testCase)
	if rspBody.Item != "paprika" {
		t.Errorf("mismatch between request and response item(%s::paprika)", rspBody.Item)
	}
	if rspBody.Category != "spices" {
		t.Errorf("mismatch between request and response category(%s::spices)", rspBody.Category)
	}

	testCase.RequestBody = `{"item": "hp", "category": "sauces"}`
	t_ext.API_Test(testCase)
	if rspBody.Item != "hp" {
		t.Errorf("mismatch between request and response item(%s::hp)", rspBody.Item)
	}
	if rspBody.Category != "sauces" {
		t.Errorf("mismatch between request and response category(%s::sauces)", rspBody.Category)
	}

	testCase.RequestBody = `{"item": "coffee", "category": "non-perishables"}`
	t_ext.API_Test(testCase)
	if rspBody.Item != "coffee" {
		t.Errorf("mismatch between request and response item(%s::coffee)", rspBody.Item)
	}
	if rspBody.Category != "non-perishables" {
		t.Errorf("mismatch between request and response category(%s::non-perishables)", rspBody.Category)
	}

	testCase.RequestBody = `{"item": "froot loops", "category": "non-cooking"}`
	t_ext.API_Test(testCase)
	if rspBody.Item != "froot loops" {
		t.Errorf("mismatch between request and response item(%s::froot loops)", rspBody.Item)
	}
	if rspBody.Category != "non-cooking" {
		t.Errorf("mismatch between request and response category(%s::non-cooking)", rspBody.Category)
	}
}

func TestCreateBatch(t *testing.T) {
	initialNoCategory := models.InventoryItem{}
	testCase := &t_ext.APITestCase{
		T:                    t,
		Method:               "POST",
		Endpoint:             "/api/inventory/create",
		ResponseBody:         &initialNoCategory,
		ExpectedResponseCode: http.StatusCreated,
	}
	t_ext.AuthorizeTestCase(t, testCase)

	testCase.RequestBody = `{"item":"ruffles"}`
	t_ext.API_Test(testCase)

	initialWithCategory := models.InventoryItem{}
	testCase.ResponseBody = &initialWithCategory
	testCase.RequestBody = `{"item":"tostitos", "category": "non-perishables"}`
	t_ext.API_Test(testCase)

	items := []models.InventoryItem{}
	testCase.ResponseBody = &items
	testCase.Endpoint = "/api/inventory/create-batch"
	testCase.RequestBody = `{"items":["ruffles", "lays stacks", "tostitos"]}`
	t_ext.API_Test(testCase)

	t_ext.TestEqual(t, len(items), 3, "invalid number of returned items").FailNowIfUnsuccessful()
	t_ext.TestEqual(t, initialNoCategory.Uid, items[0].Uid, "did not return existing item")
	t_ext.TestEqual(t, "lays stacks", items[1].Item, "incorrect item value")
	t_ext.TestEqual(t, initialWithCategory.Uid, items[2].Uid, "did not return existing item")
	laysUid := items[1].Uid

	testCase.ResponseBody = &items
	testCase.RequestBody = `{"items":["ruffles", "lays stacks", "tostitos", "mars"], "category": "non-cooking"}`
	t_ext.API_Test(testCase)
	t_ext.TestEqual(t, len(items), 4, "invalid number of returned items").FailNowIfUnsuccessful()
	t_ext.TestEqual(t, initialNoCategory.Uid, items[0].Uid, "did not return existing item")
	t_ext.TestEqual(t, laysUid, items[1].Uid, "did not return existing item")
	t_ext.TestEqual(t, initialWithCategory.Uid, items[2].Uid, "did not return existing item")
	t_ext.TestEqual(t, "mars", items[3].Item, "incorrect item value")
	for _, item := range items {
		t_ext.TestEqual(t, models.NON_COOKING, item.Category, "incorrect category")
	}
}

func init() {
	t_ext.EnsureTestUser()
}
