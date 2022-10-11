package grocery_test

import (
	t_ext "cardamom/core/ext/testing_ext"
	"cardamom/core/models"
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

	rspBody := models.GroceryItem{}
	testCase := &t_ext.APITestCase{
		T:                    t,
		Method:               "POST",
		Endpoint:             "/api/grocery/create",
		ResponseBody:         &rspBody,
		RequestBody:          `{"item":"ice"}`,
		ExpectedResponseCode: http.StatusCreated,
	}
	t_ext.AuthorizeTestCase(t, testCase)
	t_ext.API_Test(testCase)
	uid := rspBody.Uid

	testCase.Endpoint = "/api/grocery/collect"
	testCase.ExpectedResponseCode = http.StatusOK
	testCase.RequestBody = fmt.Sprintf(`{"uid": "%s", "is_collected": true}`, rspBody.Uid)
	t_ext.API_Test(testCase)
	if !rspBody.IsCollected {
		t.Errorf("item was not collected")
	}
	if rspBody.Uid != uid {
		t.Errorf("uid mismatch on collect")
	}

	testCase.Endpoint = "/api/grocery/create"
	testCase.ExpectedResponseCode = http.StatusCreated
	testCase.RequestBody = `{"item": "ice"}`
	if rspBody.Uid != uid {
		t.Errorf("uid mismatch on readd")
	}

	testCase.Endpoint = "/api/grocery/collect"
	testCase.ExpectedResponseCode = http.StatusOK
	testCase.RequestBody = fmt.Sprintf(`{"uid": "%s", "is_collected": false}`, rspBody.Uid)
	t_ext.API_Test(testCase)
	if rspBody.IsCollected {
		t.Errorf("item was not uncollected")
	}
}

func TestCreateThenDelete(t *testing.T) {

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

func init() {
	t_ext.EnsureTestUser()
}
