package grocery_test

import (
	"cardamom/core/source/db/models"
	"cardamom/core/source/services/inventory"
	"cardamom/core/test/t_ext"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
)

var groceries = []string{
	"rice", "flour", "sugar",
	"baking soda", "baking powder",
	"vanilla extract",
}

var groceriesWithStore = []string{
	"pork", "beef", "chicken", "eggs",
}

func TestSimpleCreates(t *testing.T) {
	t_ext.Init(t)

	store := "Metro"
	rspBody := models.GroceryItem{}
	testCase := &t_ext.APITestCase{
		T:                    t,
		Method:               "POST",
		Endpoint:             "/api/grocery/create",
		ResponseBody:         &rspBody,
		ExpectedResponseCode: http.StatusCreated,
	}
	t_ext.AuthorizeTestCase(t, testCase)

	for _, item := range groceries {
		testCase.RequestBody = fmt.Sprintf(`{"item":"%s", "store": null}`, item)
		t_ext.API_Test(testCase)
		if rspBody.Item != item {
			t.Errorf("attempted to create item '%s' but found '%s'", item, rspBody.Item)
		}
	}

	for _, item := range groceriesWithStore {
		testCase.RequestBody = fmt.Sprintf(`{"item":"%s", "store": "%s"}`, item, store)

		t_ext.API_Test(testCase)
		if rspBody.Item != item {
			t.Errorf("attempted to create item '%s' but found '%s'", item, rspBody.Item)
		}
		if rspBody.Store != strings.ToLower(store) {
			t.Errorf("attempted to create item at store '%s' but found '%s'", store, rspBody.Store)
		}
	}

	var rspGroceries []models.GroceryItem
	testCase.Endpoint = "/api/grocery/list"
	testCase.ExpectedResponseCode = http.StatusOK
	testCase.RequestBody = "{}"
	testCase.ResponseBody = &rspGroceries
	t_ext.API_Test(testCase)
	listedGroceries := funk.Map(
		rspGroceries, func(i models.GroceryItem) string { return i.Item }).([]string)
	for _, item := range groceries {
		if !funk.Contains(listedGroceries, item) {
			t.Errorf("groceries on server did not contain %s", item)
		}
	}
	for _, item := range groceriesWithStore {
		if !funk.Contains(listedGroceries, item) {
			t.Errorf("groceries on server did not contain %s", item)
		}
	}
}

func TestDoubleCreate(t *testing.T) {
	t_ext.Init(t)

	rspBody := models.GroceryItem{}
	testCase := &t_ext.APITestCase{
		T:                    t,
		Method:               "POST",
		Endpoint:             "/api/grocery/create",
		ResponseBody:         &rspBody,
		RequestBody:          `{"item":"broccoli"}`,
		ExpectedResponseCode: http.StatusCreated,
	}

	t_ext.AuthorizeTestCase(t, testCase)
	t_ext.API_Test(testCase)
	uid := rspBody.Uid
	t_ext.API_Test(testCase)
	if uid != rspBody.Uid {
		t.Errorf("double create should return the same item, instead found differing uids(%s --- %s)", uid, rspBody.Uid)
	}
}

func TestCollectThenCreate(t *testing.T) {
	t_ext.Init(t)

	groceryItem := models.GroceryItem{}
	testCase := &t_ext.APITestCase{
		T:                    t,
		Method:               "POST",
		Endpoint:             "/api/grocery/create",
		ResponseBody:         &groceryItem,
		RequestBody:          `{"item":"ice"}`,
		ExpectedResponseCode: http.StatusCreated,
	}
	t_ext.AuthorizeTestCase(t, testCase)
	t_ext.API_Test(testCase)
	uid := groceryItem.Uid

	type CollectResponse struct {
		GroceryItem   models.GroceryItem   `json:"grocery_item"`
		InventoryItem models.InventoryItem `json:"inventory_item"`
	}
	collectRsp := CollectResponse{}
	testCase.Endpoint = "/api/grocery/collect"
	testCase.ExpectedResponseCode = http.StatusOK
	testCase.RequestBody = fmt.Sprintf(`{"uid": "%s", "is_collected": true}`, groceryItem.Uid)
	testCase.ResponseBody = &collectRsp
	t_ext.API_Test(testCase)
	if !collectRsp.GroceryItem.IsCollected {
		t.Errorf("item was not collected")
	}
	if collectRsp.GroceryItem.Uid != uid {
		t.Errorf("uid mismatch on collect")
	}

	testCase.Endpoint = "/api/grocery/create"
	testCase.ExpectedResponseCode = http.StatusCreated
	testCase.RequestBody = `{"item": "ice"}`
	testCase.ResponseBody = &groceryItem
	if groceryItem.Uid != uid {
		t.Errorf("uid mismatch on read")
	}

	testCase.Endpoint = "/api/grocery/collect"
	testCase.ExpectedResponseCode = http.StatusOK
	testCase.RequestBody = fmt.Sprintf(`{"uid": "%s", "is_collected": false}`, groceryItem.Uid)
	t_ext.API_Test(testCase)
	if groceryItem.IsCollected {
		t.Errorf("item was not uncollected")
	}
}

func TestCollectThenUndo(t *testing.T) {
	t_ext.Init(t)
	type CollectResponse struct {
		GroceryItem   *models.GroceryItem   `json:"grocery_item"`
		InventoryItem *models.InventoryItem `json:"inventory_item"`
	}
	groceryItem := &models.GroceryItem{}
	testCase := &t_ext.APITestCase{
		T:                    t,
		Method:               "POST",
		Endpoint:             "/api/grocery/create",
		ResponseBody:         groceryItem,
		RequestBody:          `{"item":"guava"}`,
		ExpectedResponseCode: http.StatusCreated,
	}
	t_ext.AuthorizeTestCase(t, testCase)
	t_ext.API_Test(testCase)
	uid := groceryItem.Uid

	testCase.Endpoint = "/api/grocery/collect"
	testCase.ExpectedResponseCode = http.StatusOK
	testCase.RequestBody = fmt.Sprintf(`{"uid": "%s", "is_collected": true}`, uid)
	inventoryItem := &models.InventoryItem{}
	testCase.ResponseBody = &CollectResponse{
		GroceryItem:   groceryItem,
		InventoryItem: inventoryItem,
	}
	t_ext.API_Test(testCase)
	if !groceryItem.IsCollected {
		t.Errorf("item was not collected")
	}
	if groceryItem.Uid != uid {
		t.Errorf("uid mismatch on collect")
	}
	inventoryItem, err := inventory.ItemByValue("guava", t_ext.GetTestUser().Uid)
	if err != nil {
		t.Fatal("unable to get inventory item")
	}
	t_ext.TestEqual(t, inventoryItem.InStock, true, "inventory item was not in stock")

	testCase.RequestBody = fmt.Sprintf(`{"uid": "%s", "is_collected": false}`, uid)
	t_ext.API_Test(testCase)

	inventoryItem, err = inventory.ItemByValue("guava", t_ext.GetTestUser().Uid)
	if err != nil {
		t.Fatal("unable to get inventory item")
	}
	t_ext.TestEqual(t, inventoryItem.InStock, false, "inventory item was not in stock")
}

func TestCreateThenDelete(t *testing.T) {
	t_ext.Init(t)

	rspBody := models.GroceryItem{}
	testCase := &t_ext.APITestCase{
		T:                    t,
		Method:               "POST",
		Endpoint:             "/api/grocery/create",
		ResponseBody:         &rspBody,
		RequestBody:          `{"item":"beer"}`,
		ExpectedResponseCode: http.StatusCreated,
	}
	t_ext.AuthorizeTestCase(t, testCase)
	t_ext.API_Test(testCase)
	// uid := rspBody.Uid

	testCase.Endpoint = "/api/grocery/delete"
	testCase.RequestBody = fmt.Sprintf(`{"uid": "%s"}`, rspBody.Uid)
	testCase.ExpectedResponseCode = http.StatusOK
	testCase.ResponseBody = &gin.H{}
	t_ext.API_Test(testCase)
}

func TestCreateGroceryBatch(t *testing.T) {
	t_ext.Init(t)

	initialNoStore := models.GroceryItem{}
	testCase := &t_ext.APITestCase{
		T:                    t,
		Method:               "POST",
		Endpoint:             "/api/grocery/create",
		ResponseBody:         &initialNoStore,
		ExpectedResponseCode: http.StatusCreated,
	}
	t_ext.AuthorizeTestCase(t, testCase)

	testCase.RequestBody = `{"item":"lays"}`
	t_ext.API_Test(testCase)

	initialWithStore := models.GroceryItem{}
	testCase.ResponseBody = &initialWithStore
	testCase.RequestBody = `{"item":"doritos", "store": "Metro"}`
	t_ext.API_Test(testCase)

	items := []models.GroceryItem{}
	testCase.ResponseBody = &items
	testCase.Endpoint = "/api/grocery/create-batch"
	testCase.RequestBody = `{"items":["lays", "pringles", "doritos"]}`
	t_ext.API_Test(testCase)

	t_ext.TestEqual(t, len(items), 3, "invalid number of returned items").FailNowIfUnsuccessful()
	t_ext.TestEqual(t, initialNoStore.Uid, items[0].Uid, "did not return existing item")
	t_ext.TestEqual(t, "pringles", items[1].Item, "incorrect item value")
	t_ext.TestEqual(t, initialWithStore.Uid, items[2].Uid, "did not return existing item")
	pringlesUid := items[1].Uid

	testCase.ResponseBody = &items
	testCase.RequestBody = `{"items":["lays", "pringles", "doritos", "snickers"], "store": "no frills"}`
	t_ext.API_Test(testCase)
	t_ext.TestEqual(t, len(items), 4, "invalid number of returned items").FailNowIfUnsuccessful()
	t_ext.TestEqual(t, initialNoStore.Uid, items[0].Uid, "did not return existing item")
	t_ext.TestEqual(t, pringlesUid, items[1].Uid, "did not return existing item")
	t_ext.TestEqual(t, initialWithStore.Uid, items[2].Uid, "did not return existing item")
	t_ext.TestEqual(t, "snickers", items[3].Item, "incorrect item value")
	for _, item := range items {
		t_ext.TestEqual(t, "no frills", item.Store, "incorrect store")
	}
}
