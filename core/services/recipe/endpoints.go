package recipe

import (
	"cardamom/core/ext/gin_ext"
	"cardamom/core/ext/log_ext"
	m "cardamom/core/models"
	"cardamom/core/services"
	"cardamom/core/services/auth"
	"cardamom/core/services/inventory"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateRecipe(c *gin.Context, r *CreateRecipeRequest) {
	user := auth.GetActiveUserClaims(c)
	recipe := m.Recipe{
		Uid:          generateRecipeUid(),
		UserUid:      user.Uid,
		Name:         r.Name,
		Description:  r.Description,
		Meal:         r.Meal,
		Instructions: r.Instructions,
	}

	if err := m.DB.Create(&recipe).Error; err != nil {
		gin_ext.ServerError(c, log_ext.Errorf("unable to create recipe -- %w", err))
		return
	}

	recipe.Ingredients = make([]m.RecipeIngredient, len(r.Ingredients))
	for i, ingre := range r.Ingredients {
		recipe.Ingredients[i] = m.RecipeIngredient{
			Uid:       generateIngreUid(),
			UserUid:   user.Uid,
			RecipeUid: recipe.Uid,
			Meal:      r.Meal,
			SortOrder: i,
			Quantity:  ingre.Quantity,
			Unit:      ingre.Unit,
			Item:      ingre.Item,
		}
		err := m.DB.Create(&recipe.Ingredients[i]).Error
		if err != nil {
			gin_ext.ServerError(c, log_ext.Errorf("unable to create ingredient -- %w", err))
			return
		}
	}

	c.JSON(http.StatusCreated, recipe)
}

func UpdateRecipe(c *gin.Context, r *UpdateRecipeRequest) {
	user := auth.GetActiveUserClaims(c)
	recipe := m.Recipe{}
	err := m.DB.Where(&m.Recipe{UserUid: user.Uid, Uid: r.Uid}).
		Preload("Ingredients").First(&recipe).Error
	if err != nil {
		gin_ext.AbortNotFound(c, log_ext.Errorf("finding recipe -- %w", err))
		return
	}
	if r.Name != nil {
		recipe.Name = *r.Name
	}
	if r.IsTrashed != nil {
		recipe.IsTrashed = *r.IsTrashed
		if recipe.IsTrashed {
			recipe.TrashAt = uint64(time.Now().UTC().AddDate(0, 0, 30).Unix())
		}
	}
	if r.Description != nil {
		recipe.Description = *r.Description
	}
	if r.Meal != nil {
		recipe.Meal = *r.Meal
	}
	if r.Instructions != nil {
		recipe.Instructions = *r.Instructions
	}
	if err = m.DB.Save(&recipe).Error; err != nil {
		gin_ext.ServerError(c, log_ext.Errorf("saving recipe -- %w", err))
		return
	}

	sort.Slice(recipe.Ingredients, func(i, j int) bool {
		return recipe.Ingredients[i].SortOrder < recipe.Ingredients[j].SortOrder
	})

	if len(r.Ingredients) > 0 {
		if err = resizeIngredients(user.Uid, r.Ingredients, &recipe); err != nil {
			gin_ext.ServerError(c, err)
			return
		}
	}

	c.JSON(http.StatusOK, recipe)
}

func ReadRecipe(c *gin.Context, r *services.ReadRequest) {
	user := auth.GetActiveUserClaims(c)
	recipe := m.Recipe{}
	err := m.DB.Where(&m.Recipe{Uid: r.Uid, UserUid: user.Uid}).
		Preload("Ingredients").First(&recipe).Error
	if err != nil {
		gin_ext.ServerError(c, log_ext.Errorf("listing recipes -- %w", err))
		return
	}

	c.JSON(http.StatusOK, recipe)
}

func ListRecipes(c *gin.Context) {
	user := auth.GetActiveUserClaims(c)
	var recipes []m.Recipe
	err := m.DB.Where("user_uid = ? and is_trashed = false", user.Uid).Find(&recipes).Error
	if err != nil {
		gin_ext.ServerError(c, log_ext.Errorf("listing recipes -- %w", err))
		return
	}

	c.JSON(http.StatusOK, recipes)
}

func ListTrashedRecipes(c *gin.Context) {
	user := auth.GetActiveUserClaims(c)
	var recipes []m.Recipe
	err := m.DB.Where("user_uid = ? and is_trashed = true", user.Uid).Find(&recipes).Error
	if err != nil {
		gin_ext.ServerError(c, log_ext.Errorf("listing trashed recipes -- %w", err))
		return
	}

	c.JSON(http.StatusOK, recipes)
}

func GetAvailableRecipes(c *gin.Context) {

	user := auth.GetActiveUserClaims(c)

	// Get inventory
	inventory, err := inventory.GetInventory(user.Uid)
	if err != nil {
		gin_ext.ServerError(c, log_ext.Errorf("getting groceries -- %w", err))
		return
	}

	// Get all recipes that belong to this user
	var recipes []m.Recipe
	err = m.DB.Where("user_uid = ? and is_trashed = false", user.Uid).
		Preload("Ingredients").Find(&recipes).Error
	if err != nil {
		gin_ext.ServerError(c, log_ext.Errorf("getting recipes -- %w", err))
		return
	}

	c.JSON(http.StatusOK, filterRecipesByIngredients(inventory, recipes))
}

func SearchRecipe(c *gin.Context, r *SearchRecipeRequest) {
	user := auth.GetActiveUserClaims(c)
	var recipes []m.Recipe
	db := m.DB

	if r.Ingredient != nil {
		db = db.Table("recipes").Select("recipes.*").
			Joins("inner join recipe_ingredients on recipes.uid = recipe_ingredients.recipe_uid").
			Where("lower(recipe_ingredients.item) = lower(?)", *r.Ingredient)
	}
	if r.Name != nil {
		db = db.Where("lower(recipes.name) like lower(?)", "%"+*r.Name+"%")
	}
	if r.Meal != nil {
		db = db.Where("recipes.meal = ?", *r.Meal)
	}
	if r.Description != nil {
		db = db.Where("lower(recipes.description) like lower(?)", "%"+*r.Description+"%")
	}
	db.Where(&m.Recipe{UserUid: user.Uid})

	err := db.Find(&recipes).Error
	if err != nil {
		gin_ext.ServerError(c, log_ext.Errorf("searching for recipes -- %w", err))
		return
	}

	c.JSON(http.StatusOK, &recipes)
}
