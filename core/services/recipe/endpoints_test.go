package recipe_test

import (
	t_ext "cardamom/core/ext/testing_ext"
	"cardamom/core/models"
	"encoding/json"
	"net/http"
	"testing"
)

var createRecipes = []string{
	`{
    "name": "French toast", "description": "Egg soaked grilled toast", "meal": "breakfast",
    "ingredients": [
      {"quantity": 5, "unit": null, "item": "bread"},
      {"quantity": 2, "unit": null, "item": "eggs"},
      {"quantity": "1/4", "unit": "tsp", "item": "cloves"}
    ],
    "instructions": ["In a large bowl, whisk eggs with cloves", "soak bread in egg mixture", "grill bread in skillet"]
   }`,
	`{
    "name": "Bangers and Mash", "description": "Traditional Irish dinner", "meal": "dinner",
    "ingredients": [
      {"quantity": 250, "unit": "g", "item": "sausage"},
      {"quantity": "500", "unit": "g", "item": "potato"},
      {"quantity": "1 1/2", "unit": "tsp", "item": "black pepper"},
      {"quantity": "1 1/2", "unit": "tsp", "item": "salt"}
    ],
    "instructions": ["boil potatoes", "fry sausage in skillet", "mash potato with salt and pepper", "combine"]
   }`,
	`{
    "name": "Bolonese", "description": "Italian meat sauce pasta", "meal": "dinner",
    "ingredients": [
      {"quantity": 200, "unit": "g", "item": "chicken"},
      {"quantity": "2 1/2", "unit": "cup", "item": "rice"},
      {"quantity": "1/2", "unit": "tsp", "item": "cumin"}
    ],
    "instructions": ["Grill chicken with cumin", "Cook rice", "Combine"]
   }`,
}

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

func TestCreates(t *testing.T) {

	testCase := getCreateTestCase(t)
	for _, recipe := range createRecipes {
		testCase.RequestBody = recipe
		t_ext.API_Test(testCase)
	}
}

var searchRecipes = []string{
	`{
    "name": "Croissant", "description": "French baked good", "meal": "breakfast",
    "ingredients": [
      {"quantity": "1 1/2", "unit": "cup", "item": "butter"},
      {"quantity": 2, "unit": "cup", "item": "flour"},
      {"quantity": 1, "unit": "tsp", "item": "yeast"}
    ],
    "instructions": ["With flour, yeast and butter make puff pastry", "flatten then cut into triangles, then roll", "bake at 350F for 15 minutes"]
   }`,
	`{
    "name": "French Onion Soup", "description": "French classic food", "meal": "lunch",
    "ingredients": [
      {"quantity": 4, "unit": "cup", "item": "water"},
      {"quantity": "4", "unit": null, "item": "bread"},
      {"quantity": "1", "unit": "cup", "item": "onion"},
      {"quantity": "1", "unit": "cup", "item": "cheese"}
    ],
    "instructions": ["Boil water, slice onions", "Combine"]
   }`,
	`{
    "name": "Escargot", "description": "French delicacy", "meal": "dinner",
    "ingredients": [
      {"quantity": 1, "unit": "cup", "item": "snails"},
      {"quantity": 2, "unit": "Tbsp", "item": "butter"},
      {"quantity": 1, "unit": "Tbsp", "item": "parsely"},
      {"quantity": 1, "unit": "tsp", "item": "salt"}
    ],
    "instructions": ["Fry snails, then combine"]
   }`,
	`{
    "name": "Chicken Confit", "description": "French oven roasted chicken", "meal": "dinner",
    "ingredients": [
      {"quantity": 2, "unit": null, "item": "chicken thighs"},
      {"quantity": 2, "unit": "Tbsp", "item": "butter"},
      {"quantity": 1, "unit": "Tbsp", "item": "garlic"},
      {"quantity": "1 1/2", "unit": "Tbsp", "item": "lemon juice"}
    ],
    "instructions": ["Marinate chicken", "Cook in oven for 25 minutes"]
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
    "instructions": ["Grease heavy based pot with butter", "Cook onions until golden brown", "Add beef, carrots and red wine", "stew for at least 30 minutes"]
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
			Preload("Instructions").Preload("Ingredients").
			Find(&resultRecipe).Error; err != nil {
			panic(err)
		}

		if len(resultRecipe) == 0 {
			testCase.RequestBody = recipeRaw
			t_ext.API_Test(testCase)
		}
	}
}

func init() {
	t_ext.EnsureTestUser()
}
