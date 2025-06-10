package recipe

import (
	"cardamom/core/source/db"
	"cardamom/core/source/db/models"
	"cardamom/core/source/ext/gin_ext"
	"cardamom/core/source/ext/log_ext"
	"cardamom/core/source/services"
	"cardamom/core/source/services/auth"
	"cardamom/core/source/services/inventory"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
)

func CreateRecipe(c *gin.Context, r *CreateRecipeRequest) {
	user := auth.GetActiveUserClaims(c)
	recipe := &models.Recipe{
		Uid:          generateRecipeUid(),
		UserUid:      user.Uid,
		Name:         r.Name,
		Description:  r.Description,
		Meal:         r.Meal,
		Instructions: r.Instructions,
	}

	if err := db.DB().Create(recipe).Error; err != nil {
		gin_ext.ServerError(c, log_ext.Errorf("unable to create recipe -- %w", err))
		return
	}

	recipe.Ingredients = make([]models.RecipeIngredient, len(r.Ingredients))
	for i, ingre := range r.Ingredients {
		recipe.Ingredients[i] = models.RecipeIngredient{
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
		err := db.DB().Create(&recipe.Ingredients[i]).Error
		if err != nil {
			gin_ext.ServerError(c, log_ext.Errorf("unable to create ingredient -- %w", err))
			return
		}
	}

	c.JSON(http.StatusCreated, recipe)
}

func ReadRecipe(c *gin.Context, r *services.ReadRequest) {
	user := auth.GetActiveUserClaims(c)
	recipe := &models.Recipe{}
	err := db.DB().Where(&models.Recipe{Uid: r.Uid, UserUid: user.Uid}).
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
	recipe := &models.Recipe{}
	if err := db.DB().Where(&models.Recipe{UserUid: user.Uid, Uid: r.Uid}).First(recipe).Error; err != nil {
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
	if err := db.DB().Save(recipe).Error; err != nil {
		gin_ext.ServerError(c, log_ext.Errorf("saving recipe -- %w", err))
		return
	}

	c.JSON(http.StatusOK, recipe)
}

func CreateRecipeIngredient(c *gin.Context, r *CreateRecipeIngredientRequest) {
	user := auth.GetActiveUserClaims(c)
	recipe := &models.Recipe{}
	err := db.DB().Where(&models.Recipe{UserUid: user.Uid, Uid: r.RecipeUid}).First(recipe).Error
	if err != nil {
		gin_ext.Abort(c, http.StatusBadRequest, log_ext.Errorf("finding recipe ingredient -- %w", err))
		return
	}

	ingredient := &models.RecipeIngredient{
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
	if err := db.DB().Save(ingredient).Error; err != nil {
		gin_ext.ServerError(c, log_ext.Errorf("creating new recipe -- %w", err))
	}

	c.JSON(http.StatusCreated, ingredient)
}

func UpdateRecipeIngredient(c *gin.Context, r *UpdateRecipeIngredientRequest) {
	user := auth.GetActiveUserClaims(c)
	ingredient := &models.RecipeIngredient{}
	err := db.DB().Where(&models.RecipeIngredient{UserUid: user.Uid, Uid: r.Uid}).First(ingredient).Error
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

	if err = db.DB().Save(ingredient).Error; err != nil {
		gin_ext.ServerError(c, log_ext.Errorf("saving ingredient -- %w", err))
		return
	}

	c.JSON(http.StatusOK, ingredient)
}

func ReorderRecipeIngredients(c *gin.Context, r *ReorderRecipeIngredientRequest) {
	user := auth.GetActiveUserClaims(c)
	recipe := &models.Recipe{}
	err := db.DB().Where(&models.Recipe{Uid: r.Uid, UserUid: user.Uid}).
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
		if err := db.DB().Save(&ingre).Error; err != nil {
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
	ingredient := models.RecipeIngredient{}
	if err := db.DB().Where(&models.Recipe{Uid: r.Uid, UserUid: user.Uid}).First(&ingredient).Error; err != nil {
		gin_ext.ServerError(c, log_ext.Errorf("reading recipe ingredient -- %w", err))
		return
	}

	if err := db.DB().Delete(&ingredient).Error; err != nil {
		gin_ext.ServerError(c, log_ext.Errorf("deleting recipe ingredient -- %w", err))
		return
	}

	c.JSON(http.StatusOK, &gin.H{})
}

func ListRecipes(c *gin.Context) {
	user := auth.GetActiveUserClaims(c)
	var recipes []models.Recipe
	err := db.DB().Where("user_uid = ? and is_trashed = false", user.Uid).Find(&recipes).Error
	if err != nil {
		gin_ext.ServerError(c, log_ext.Errorf("listing recipes -- %w", err))
		return
	}

	c.JSON(http.StatusOK, &recipes)
}

func ListTrashedRecipes(c *gin.Context) {
	user := auth.GetActiveUserClaims(c)
	var recipes []models.Recipe
	err := db.DB().Where("user_uid = ? and is_trashed = true", user.Uid).Find(&recipes).Error
	if err != nil {
		gin_ext.ServerError(c, log_ext.Errorf("listing trashed recipes -- %w", err))
		return
	}

	c.JSON(http.StatusOK, &recipes)
}

func SearchRecipe(c *gin.Context, r *SearchRecipeRequest) {
	user := auth.GetActiveUserClaims(c)
	var recipes []models.Recipe
	query := db.DB()

	if r.Ingredient != nil {
		query = query.Table("recipes").Select("recipes.*").
			Joins("inner join recipe_ingredients on recipes.uid = recipe_ingredients.recipe_uid").
			Where("lower(recipe_ingredients.item) = lower(?)", *r.Ingredient)
	}
	if r.Name != nil {
		query = query.Where("lower(recipes.name) like lower(?)", "%"+*r.Name+"%")
	}
	if r.Meal != nil {
		query = query.Where("recipes.meal = ?", *r.Meal)
	}
	if r.Description != nil {
		query = query.Where("lower(recipes.description) like lower(?)", "%"+*r.Description+"%")
	}
	query.Where(&models.Recipe{UserUid: user.Uid})

	err := query.Find(&recipes).Error
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
	var recipes []models.Recipe
	query := db.DB().Where("user_uid = ? and is_trashed = false", user.Uid)
	if r.Meal != nil {
		query = query.Where(models.Recipe{Meal: *r.Meal})
	}

	query.Preload("Ingredients").Find(&recipes)
	if query.Error != nil {
		gin_ext.ServerError(c, log_ext.Errorf("getting recipes -- %w", err))
		return
	}

	c.JSON(http.StatusOK, filterRecipesByIngredients(inventory, recipes))
}
