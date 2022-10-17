package recipe_test

import (
	t_ext "cardamom/core/ext/testing_ext"
	"cardamom/core/ext/units"
	"cardamom/core/models"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
)

func getCreateTestCase(t *testing.T) *t_ext.APITestCase {
	testCase := &t_ext.APITestCase{
		T:                    t,
		Method:               "POST",
		Endpoint:             "/api/recipe/create",
		ResponseBody:         &models.Recipe{},
		ExpectedResponseCode: http.StatusCreated,
	}
	t_ext.AuthorizeTestCase(t, testCase)
	return testCase
}

func TestCreate(t *testing.T) {

	testCase := getCreateTestCase(t)
	rspBody := testCase.ResponseBody.(*models.Recipe)
	testCase.RequestBody = `{
    "name": "Bolonese", "description": "Italian meat sauce pasta", "meal": "dinner",
    "ingredients": [
      {"quantity": 200, "unit": "g", "item": "chicken", "modifier": "chopped"},
      {"quantity": "2 1/2", "unit": "cup", "item": "rice", "modifier": "", "optional": false},
      {"quantity": "1/2", "unit": "cup", "item": "tomato sauce", "modifier": "", "optional": false},
      {"quantity": "1", "unit": "tsp", "item": "black pepper", "modifier": "", "optional": false},
      {"quantity": "1/2", "unit": "tsp", "item": "salt", "modifier": "", "optional": false},
      {"quantity": "1/2", "unit": "tsp", "item": "cumin", "modifier": "", "optional": true}
    ],
    "instructions": "Grill chicken with cumin\nBoil pasta\nCombine"
   }`

	t_ext.API_Test(testCase)

	if rspBody.Name != "Bolonese" {
		t.Error("Incorrect name in response")
	}
	if rspBody.Description != "Italian meat sauce pasta" {
		t.Error("Incorrect description in response")
	}
	if rspBody.Meal != models.DINNER {
		t.Error("Incorrect meal in response")
	}
	if rspBody.Instructions != "Grill chicken with cumin\nBoil pasta\nCombine" {
		t.Error("Incorrect instructions in response")
	}

	if rspBody.Ingredients[0].Quantity.String() != "200" {
		t.Error("Incorrect ingredient quantity in response")
	}
	if rspBody.Ingredients[0].Unit.Name() != units.Gram.Name() {
		t.Error("Incorrect unit quantity in response")
	}
	if rspBody.Ingredients[0].Item != "chicken" {
		t.Error("Incorrect ingredient item in response")
	}
	if rspBody.Ingredients[0].Modifier != "chopped" {
		t.Error("Incorrect ingredient modifier in response")
	}
	if rspBody.Ingredients[0].Optional {
		t.Error("Incorrect ingredient optional in response")
	}

	if rspBody.Ingredients[1].Quantity.String() != "2 1/2" {
		t.Error("Incorrect ingredient quantity in response")
	}
	if rspBody.Ingredients[1].Optional {
		t.Error("Incorrect ingredient optional in response")
	}

	if rspBody.Ingredients[2].Quantity.String() != "1/2" {
		t.Error("Incorrect ingredient quantity in response")
	}
	if rspBody.Ingredients[3].Quantity.String() != "1" {
		t.Error("Incorrect ingredient quantity in response")
	}
	if !rspBody.Ingredients[5].Optional {
		t.Error("Incorrect ingredient optional in response")
	}
}

func TestModify(t *testing.T) {
	testCase := getCreateTestCase(t)
	recipeCreateRsp := testCase.ResponseBody.(*models.Recipe)
	t_ext.AuthorizeTestCase(t, testCase)
	testCase.RequestBody = `{
    "name": "French toast", "description": "Egg soaked grilled toast", "meal": "breakfast",
    "ingredients": [
      {"quantity": 5, "unit": null, "item": "bread", "modifier": "", "optional": false},
      {"quantity": 2, "unit": null, "item": "eggs", "modifier": "", "optional": false}
    ],
    "instructions": "In a large bowl, whisk eggs with cloves\nsoak bread in egg mixture\ngrill bread in skillet"
   }`
	t_ext.API_Test(testCase)

	recipeUpdateRsp := models.Recipe{}
	testCase.ResponseBody = &recipeUpdateRsp
	testCase.ExpectedResponseCode = http.StatusOK
	testCase.Endpoint = "/api/recipe/update"
	testCase.RequestBody = fmt.Sprintf(`{
    "uid": "%s", "is_trashed": true,
    "name": "Bangers and Mash", "description": "Traditional Irish dinner", "meal": "dinner",
    "instructions": "boil potatoes\nfry sausage in skillet\nmash potato with salt and pepper\ncombine"
   }`, recipeCreateRsp.Uid)
	t_ext.API_Test(testCase)
	if recipeUpdateRsp.Name != "Bangers and Mash" {
		t.Error("Incorrect name in response")
	}
	if recipeUpdateRsp.Description != "Traditional Irish dinner" {
		t.Error("Incorrect description in response")
	}
	if recipeUpdateRsp.Meal != models.DINNER {
		t.Error("Incorrect meal in response")
	}
	if !recipeUpdateRsp.IsTrashed {
		t.Error("Incorrect trashed status in response")
	}

	ingredientRsp := models.RecipeIngredient{}
	testCase.ResponseBody = &ingredientRsp
	testCase.Endpoint = "/api/recipe/ingredient/update"
	testCase.RequestBody = fmt.Sprintf(`{ "uid": "%s", "quantity": 250 }`, recipeCreateRsp.Ingredients[0].Uid)
	t_ext.API_Test(testCase)
	if ingredientRsp.Quantity.String() != "250" {
		t.Error("Incorrect ingredient quantity in response")
	}

	testCase.RequestBody = fmt.Sprintf(`{ "uid": "%s", "unit": "g" }`, recipeCreateRsp.Ingredients[0].Uid)
	t_ext.API_Test(testCase)
	if ingredientRsp.Unit.Name() != units.Gram.Name() {
		t.Error("Incorrect unit quantity in response")
	}

	testCase.RequestBody = fmt.Sprintf(`{ "uid": "%s", "item": "sausage" }`, recipeCreateRsp.Ingredients[0].Uid)
	t_ext.API_Test(testCase)
	if ingredientRsp.Item != "sausage" {
		t.Error("Incorrect ingredient item in response")
	}

	testCase.RequestBody = fmt.Sprintf(`{ "uid": "%s", "modifier": "chopped" }`, recipeCreateRsp.Ingredients[0].Uid)
	t_ext.API_Test(testCase)
	if ingredientRsp.Modifier != "chopped" {
		t.Error("Incorrect ingredient modifier in response")
	}

	testCase.RequestBody = fmt.Sprintf(`{ "uid": "%s", "optional": true }`, recipeCreateRsp.Ingredients[0].Uid)
	t_ext.API_Test(testCase)
	if !ingredientRsp.Optional {
		t.Error("Incorrect ingredient optional in response")
	}

	testCase.RequestBody = fmt.Sprintf(`{
    "uid": "%s", "quantity": "2 1/2",
    "unit": "cup", "item": "potato",
    "modifier": null, "optional": null
  }`, recipeCreateRsp.Ingredients[1].Uid)
	t_ext.API_Test(testCase)
	if ingredientRsp.Quantity.String() != "2 1/2" {
		t.Error("Incorrect ingredient quantity in response")
	}
	if ingredientRsp.Unit.Name() != units.Cup.Name() {
		t.Error("Incorrect unit in response")
	}
	if ingredientRsp.Item != "potato" {
		t.Error("Incorrect ingredient item in response")
	}
	if ingredientRsp.Modifier != "" {
		t.Error("Incorrect ingredient modifier in response")
	}
	if ingredientRsp.Optional {
		t.Error("Incorrect ingredient optional in response")
	}

	testCase.Endpoint = "/api/recipe/ingredient/create"
	testCase.ExpectedResponseCode = http.StatusCreated
	testCase.ResponseBody = &ingredientRsp
	testCase.RequestBody = fmt.Sprintf(`{
    "recipe_uid": "%s",
    "quantity": "1/2", "unit": "tsp", "item": "black pepper",
    "order": 2, "modifier": "", "optional": true}`, recipeCreateRsp.Uid)
	t_ext.API_Test(testCase)
	if ingredientRsp.Quantity.String() != "1/2" {
		t.Error("Incorrect ingredient quantity in response")
	}
	if ingredientRsp.Unit.Name() != units.Teaspoon.Name() {
		t.Error("Incorrect unit quantity in response")
	}
	if ingredientRsp.Item != "black pepper" {
		t.Error("Incorrect ingredient item in response")
	}
	if ingredientRsp.Modifier != "" {
		t.Error("Incorrect ingredient modifier in response")
	}
	if !ingredientRsp.Optional {
		t.Error("Incorrect ingredient optional in response")
	}

	testCase.ResponseBody = &recipeCreateRsp
	testCase.ExpectedResponseCode = http.StatusOK
	testCase.Endpoint = "/api/recipe/read"
	testCase.RequestBody = fmt.Sprintf(`{ "uid": "%s"}`, recipeCreateRsp.Uid)
	t_ext.API_Test(testCase)
	if len(recipeCreateRsp.Ingredients) != 3 {
		t.Fatalf("incorrect number of ingredients returned")
	}

	testCase.ResponseBody = &gin.H{}
	testCase.Endpoint = "/api/recipe/ingredient/delete"
	testCase.RequestBody = fmt.Sprintf(`{"uid": "%s"}`, ingredientRsp.Uid)
	t_ext.API_Test(testCase)

	testCase.ResponseBody = &recipeCreateRsp
	testCase.Endpoint = "/api/recipe/read"
	testCase.RequestBody = fmt.Sprintf(`{ "uid": "%s"}`, recipeCreateRsp.Uid)
	t_ext.API_Test(testCase)
	if len(recipeCreateRsp.Ingredients) != 2 {
		t.Fatalf("incorrect number of ingredients returned")
	}
}

var searchRecipes = []string{
	`{
    "name": "Croissant", "description": "French baked good", "meal": "breakfast",
    "ingredients": [
      {"quantity": "1 1/2", "unit": "cup", "item": "butter"},
      {"quantity": 2, "unit": "cup", "item": "flour"},
      {"quantity": "1/2", "unit": "cup", "item": "chocolate", "optional": true},
      {"quantity": 1, "unit": "tsp", "item": "yeast"}
    ],
    "instructions": "With flour, yeast and butter make puff pastry\nflatten then cut into triangles, then roll\nbake at 350F for 15 minutes"
   }`,
	`{
    "name": "Mille Feuille", "description": "Light and flaky french baked good", "meal": "breakfast",
    "ingredients": [
      {"quantity": "1 1/2", "unit": "cup", "item": "butter"},
      {"quantity": 2, "unit": "cup", "item": "flour"},
      {"quantity": "1/2", "unit": "cup", "item": "corn starch"},
      {"quantity": 1, "unit": "tsp", "item": "yeast"}
    ],
    "instructions": "With flour, yeast and butter make puff pastry\nflatten then cut into triangles, then roll\nbake at 350F for 15 minutes"
   }`,
	`{
    "name": "French Onion Soup", "description": "French classic food", "meal": "lunch",
    "ingredients": [
      {"quantity": 4, "unit": "cup", "item": "water"},
      {"quantity": "4", "unit": null, "item": "bread"},
      {"quantity": "1", "unit": "cup", "item": "onion"},
      {"quantity": "1", "unit": "cup", "item": "cheese"}
    ],
    "instructions": "Boil water, slice onions\nCombine"
   }`,
	`{
    "name": "Escargot", "description": "French delicacy", "meal": "dinner",
    "ingredients": [
      {"quantity": 1, "unit": "cup", "item": "snails"},
      {"quantity": 2, "unit": "Tbsp", "item": "butter"},
      {"quantity": 1, "unit": "Tbsp", "item": "parsely", "optional": true},
      {"quantity": 1, "unit": "tsp", "item": "salt"}
    ],
    "instructions": "Fry snails, then combine"
   }`,
	`{
    "name": "Chicken Confit", "description": "French oven roasted chicken", "meal": "dinner",
    "ingredients": [
      {"quantity": 2, "unit": null, "item": "chicken thighs"},
      {"quantity": 2, "unit": "Tbsp", "item": "butter"},
      {"quantity": 1, "unit": "Tbsp", "item": "garlic"},
      {"quantity": "1 1/2", "unit": "Tbsp", "item": "lemon juice"}
    ],
    "instructions": "Marinate chicken\nCook in oven for 25 minutes"
   }`,
	`{
    "name": "Beef Bourguignon", "description": "French beef stew", "meal": "dinner",
    "ingredients": [
      {"quantity": "1 1/2", "unit": "kg", "item": "beef"},
      {"quantity": "1/2", "unit": "cup", "item": "onion"},
      {"quantity": "1/2", "unit": "cup", "item": "carrots"},
      {"quantity": 1, "unit": "Tbsp", "item": "butter"},
      {"quantity": "1/4", "unit": "cup", "item": "red wine"}
    ],
    "instructions": "Grease heavy based pot with butter\nCook onions until golden brown\nAdd beef, carrots and red wine\nstew for at least 30 minutes"
   }`,
}

func TestSearch(t *testing.T) {

	ensureSearchData(t)
	results := []models.Recipe{}
	testCase := &t_ext.APITestCase{
		T:                    t,
		Method:               "POST",
		Endpoint:             "/api/recipe/search",
		ResponseBody:         &results,
		ExpectedResponseCode: http.StatusOK,
	}
	t_ext.AuthorizeTestCase(t, testCase)

	testCase.RequestBody = `{ "name": "french" }`
	t_ext.API_Test(testCase)
	if len(results) < 1 {
		t.Error("Unable to find appropriate results")
	}

	testCase.RequestBody = `{ "description": "french" }`
	t_ext.API_Test(testCase)
	if len(results) < 5 {
		t.Error("Unable to find appropriate results")
	}

	testCase.RequestBody = `{ "meal": "lunch" }`
	t_ext.API_Test(testCase)
	if len(results) < 1 {
		t.Error("Unable to find appropriate results")
	}

	testCase.RequestBody = `{ "ingredient": "butter" }`
	t_ext.API_Test(testCase)
	if len(results) < 4 {
		t.Error("Unable to find appropriate results")
	}

	testCase.RequestBody = `{ "ingredient": "butter", "name": "beef" }`
	t_ext.API_Test(testCase)
	if len(results) < 1 {
		t.Error("Unable to find appropriate results")
	}

	testCase.RequestBody = `{ "meal": "lunch", "description": "french" }`
	t_ext.API_Test(testCase)
	if len(results) < 1 {
		t.Error("Unable to find appropriate results")
	}

	testCase.RequestBody = `{ "meal": "lunch", "ingredient": "water" }`
	t_ext.API_Test(testCase)
	if len(results) < 1 {
		t.Error("Unable to find appropriate results")
	}
}

func TestGetAvailableRecipes(t *testing.T) {

	ensureSearchData(t)
	results := []models.Recipe{}
	testCase := &t_ext.APITestCase{
		T:                    t,
		Method:               "POST",
		Endpoint:             "/api/recipe/available",
		RequestBody:          `{"meal": "breakfast"}`,
		ResponseBody:         &results,
		ExpectedResponseCode: http.StatusOK,
	}
	t_ext.AuthorizeTestCase(t, testCase)
	t_ext.API_Test(testCase)

	resultNames := funk.Map(results, func(r models.Recipe) string { return r.Name }).([]string)
	if !funk.Contains(resultNames, "Croissant") {
		t.Error("Did not find Croissant, which was expected")
	}
	if funk.Contains(resultNames, "Mille Feuille") {
		t.Error("Found Mille Feuille, which should not be available")
	}
}

func ensureSearchData(t *testing.T) {

	user := t_ext.GetTestUser()
	testCase := getCreateTestCase(t)

	for _, recipeRaw := range searchRecipes {

		var recipe map[string]any
		if err := json.Unmarshal([]byte(recipeRaw), &recipe); err != nil {
			panic(err)
		}

		var resultRecipe []models.Recipe
		if err := models.DB.Where(&models.Recipe{Name: recipe["name"].(string), UserUid: user.Uid}).
			Preload("Ingredients").Find(&resultRecipe).Error; err != nil {
			panic(err)
		}

		if len(resultRecipe) == 0 {
			testCase.RequestBody = recipeRaw
			t_ext.API_Test(testCase)
		}
	}

	inventoryItems := []string{"butter", "flour", "yeast"}
	inventoryCase := &t_ext.APITestCase{
		T:                    t,
		Method:               "POST",
		Endpoint:             "/api/inventory/create",
		ResponseBody:         &models.InventoryItem{},
		ExpectedResponseCode: http.StatusCreated,
	}
	t_ext.AuthorizeTestCase(t, inventoryCase)
	for _, item := range inventoryItems {
		inventoryCase.RequestBody = fmt.Sprintf(`{"item": "%s"}`, item)
		t_ext.API_Test(inventoryCase)
	}
}

func init() {
	t_ext.EnsureTestUser()
}
