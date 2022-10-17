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
	"github.com/thoas/go-funk"
)

func CreateRecipe(c *gin.Context, r *CreateRecipeRequest) {
	user := auth.GetActiveUserClaims(c)
	recipe := &m.Recipe{
		Uid:          generateRecipeUid(),
		UserUid:      user.Uid,
		Name:         r.Name,
		Description:  r.Description,
		Meal:         r.Meal,
		Instructions: r.Instructions,
	}

	if err := m.DB.Create(recipe).Error; err != nil {
		gin_ext.ServerError(c, log_ext.Errorf("unable to create recipe -- %w", err))
		return
	}

	recipe.Ingredients = make([]m.RecipeIngredient, len(r.Ingredients))
	for i, ingre := range r.Ingredients {
		recipe.Ingredients[i] = m.RecipeIngredient{
			Uid:       generateIngreUid(),
			UserUid:   user.Uid,
			RecipeUid: recipe.Uid,
			SortOrder: i,
			Quantity:  ingre.Quantity,
			Unit:      ingre.Unit,
			Item:      ingre.Item,
			Modifier:  ingre.Modifier,
			Optional:  ingre.Optional,
		}
		err := m.DB.Create(&recipe.Ingredients[i]).Error
		if err != nil {
			gin_ext.ServerError(c, log_ext.Errorf("unable to create ingredient -- %w", err))
			return
		}
	}

	c.JSON(http.StatusCreated, recipe)
}

func ReadRecipe(c *gin.Context, r *services.ReadRequest) {
	user := auth.GetActiveUserClaims(c)
	recipe := &m.Recipe{}
	err := m.DB.Where(&m.Recipe{Uid: r.Uid, UserUid: user.Uid}).
		Preload("Ingredients").First(recipe).Error
	if err != nil {
		gin_ext.ServerError(c, log_ext.Errorf("listing recipes -- %w", err))
		return
	}
	sort.Slice(recipe.Ingredients, func(i, j int) bool {
		return recipe.Ingredients[i].SortOrder < recipe.Ingredients[j].SortOrder
	})

	c.JSON(http.StatusOK, recipe)
}

func UpdateRecipe(c *gin.Context, r *UpdateRecipeRequest) {
	user := auth.GetActiveUserClaims(c)
	recipe := &m.Recipe{}
	if err := m.DB.Where(&m.Recipe{UserUid: user.Uid, Uid: r.Uid}).First(recipe).Error; err != nil {
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
	if err := m.DB.Save(recipe).Error; err != nil {
		gin_ext.ServerError(c, log_ext.Errorf("saving recipe -- %w", err))
		return
	}

	c.JSON(http.StatusOK, recipe)
}

func CreateRecipeIngredient(c *gin.Context, r *CreateRecipeIngredientRequest) {
	user := auth.GetActiveUserClaims(c)
	recipe := &m.Recipe{}
	err := m.DB.Where(&m.Recipe{UserUid: user.Uid, Uid: r.RecipeUid}).First(recipe).Error
	if err != nil {
		gin_ext.Abort(c, http.StatusBadRequest, log_ext.Errorf("finding recipe ingredient -- %w", err))
		return
	}

	ingredient := &m.RecipeIngredient{
		Uid:       generateIngreUid(),
		UserUid:   user.Uid,
		RecipeUid: recipe.Uid,
		SortOrder: r.SortOrder,
		Quantity:  r.Quantity,
		Unit:      r.Unit,
		Item:      r.Item,
		Optional:  r.Optional,
		Modifier:  r.Modifier,
	}
	if err := m.DB.Save(ingredient).Error; err != nil {
		gin_ext.ServerError(c, log_ext.Errorf("creating new recipe -- %w", err))
	}

	c.JSON(http.StatusCreated, ingredient)
}

func UpdateRecipeIngredient(c *gin.Context, r *UpdateRecipeIngredientRequest) {
	user := auth.GetActiveUserClaims(c)
	ingredient := &m.RecipeIngredient{}
	err := m.DB.Where(&m.RecipeIngredient{UserUid: user.Uid, Uid: r.Uid}).First(ingredient).Error
	if err != nil {
		gin_ext.AbortNotFound(c, log_ext.Errorf("finding recipe ingredient -- %w", err))
		return
	}

	if r.Quantity != nil {
		ingredient.Quantity = *r.Quantity
	}
	if r.Unit != nil {
		ingredient.Unit = r.Unit
	}
	if r.Item != nil {
		ingredient.Item = *r.Item
	}
	if r.Optional != nil {
		ingredient.Optional = *r.Optional
	}
	if r.Modifier != nil {
		ingredient.Modifier = *r.Modifier
	}

	if err = m.DB.Save(ingredient).Error; err != nil {
		gin_ext.ServerError(c, log_ext.Errorf("saving ingredient -- %w", err))
		return
	}

	c.JSON(http.StatusOK, ingredient)
}

func ReorderRecipeIngredients(c *gin.Context, r *ReorderRecipeIngredientRequest) {
	user := auth.GetActiveUserClaims(c)
	recipe := &m.Recipe{}
	err := m.DB.Where(&m.Recipe{Uid: r.Uid, UserUid: user.Uid}).
		Preload("Ingredients").First(recipe).Error
	if err != nil {
		gin_ext.ServerError(c, log_ext.Errorf("reading recipe -- %w", err))
		return
	}

	for i, ingre := range recipe.Ingredients {
		newIndex := funk.IndexOf(r.Ingredients, ingre.Uid)
		if newIndex < 0 {
			gin_ext.ServerError(c, log_ext.Errorf("missing recipe uid(%s) in request", ingre.Uid))
			return
		}

		recipe.Ingredients[i].SortOrder = newIndex
	}

	for _, ingre := range recipe.Ingredients {
		if err := m.DB.Save(&ingre).Error; err != nil {
			gin_ext.ServerError(c, log_ext.Errorf("saving ingredient -- %w", err))
			return
		}
	}

	sort.Slice(recipe.Ingredients, func(i, j int) bool {
		return recipe.Ingredients[i].SortOrder < recipe.Ingredients[j].SortOrder
	})

	c.JSON(http.StatusOK, recipe)
}

func DeleteRecipeIngredient(c *gin.Context, r *services.ReadRequest) {
	user := auth.GetActiveUserClaims(c)
	ingredient := m.RecipeIngredient{}
	if err := m.DB.Where(&m.Recipe{Uid: r.Uid, UserUid: user.Uid}).First(&ingredient).Error; err != nil {
		gin_ext.ServerError(c, log_ext.Errorf("reading recipe ingredient -- %w", err))
		return
	}

	if err := m.DB.Delete(&ingredient).Error; err != nil {
		gin_ext.ServerError(c, log_ext.Errorf("deleting recipe ingredient -- %w", err))
		return
	}

	c.JSON(http.StatusOK, &gin.H{})
}

func ListRecipes(c *gin.Context) {
	user := auth.GetActiveUserClaims(c)
	var recipes []m.Recipe
	err := m.DB.Where("user_uid = ? and is_trashed = false", user.Uid).Find(&recipes).Error
	if err != nil {
		gin_ext.ServerError(c, log_ext.Errorf("listing recipes -- %w", err))
		return
	}

	c.JSON(http.StatusOK, &recipes)
}

func ListTrashedRecipes(c *gin.Context) {
	user := auth.GetActiveUserClaims(c)
	var recipes []m.Recipe
	err := m.DB.Where("user_uid = ? and is_trashed = true", user.Uid).Find(&recipes).Error
	if err != nil {
		gin_ext.ServerError(c, log_ext.Errorf("listing trashed recipes -- %w", err))
		return
	}

	c.JSON(http.StatusOK, &recipes)
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

func GetAvailableRecipes(c *gin.Context, r *GetAvailableRecipeRequest) {

	user := auth.GetActiveUserClaims(c)

	// Get inventory
	inventory, err := inventory.GetInventory(user.Uid)
	if err != nil {
		gin_ext.ServerError(c, log_ext.Errorf("getting groceries -- %w", err))
		return
	}

	// Get all recipes that belong to this user
	var recipes []m.Recipe
	db := m.DB.Where("user_uid = ? and is_trashed = false", user.Uid)
	if r.Meal != nil {
		db = db.Where(m.Recipe{Meal: *r.Meal})
	}

	db.Preload("Ingredients").Find(&recipes)
	if db.Error != nil {
		gin_ext.ServerError(c, log_ext.Errorf("getting recipes -- %w", err))
		return
	}

	c.JSON(http.StatusOK, filterRecipesByIngredients(inventory, recipes))
}
