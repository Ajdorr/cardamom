package recipe

import (
	"cardamom/core/ext/gin_ext"
	m "cardamom/core/models"
	"cardamom/core/services"
	"cardamom/core/services/auth"
	"cardamom/core/services/inventory"
	"fmt"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
)

func CreateRecipe(c *gin.Context, r *CreateRecipeRequest) {
	user := auth.GetActiveUserClaims(c)
	recipe := m.Recipe{
		Uid:         generateRecipeUid(),
		UserUid:     user.Uid,
		Name:        r.Name,
		Description: r.Description,
		Meal:        r.Meal,
	}

	if err := m.DB.Create(&recipe).Error; err != nil {
		gin_ext.ServerError(c, fmt.Errorf("unable to create recipe -- %w", err))
		return
	}

	recipe.Instructions = make([]m.RecipeInstruction, len(r.Instructions))
	for i, instr := range r.Instructions {
		recipe.Instructions[i] = m.RecipeInstruction{
			Uid:       generateInstrUid(),
			UserUid:   user.Uid,
			RecipeUid: recipe.Uid,
			Meal:      r.Meal,
			SortOrder: i,
			Text:      instr,
		}
		err := m.DB.Create(&recipe.Instructions[i]).Error
		if err != nil {
			gin_ext.ServerError(c, fmt.Errorf("unable to create instruction -- %w", err))
			return
		}
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
			gin_ext.ServerError(c, fmt.Errorf("unable to create ingredient -- %w", err))
			return
		}
	}

	c.JSON(http.StatusCreated, recipe)
}

func UpdateRecipe(c *gin.Context, r *UpdateRecipeRequest) {
	user := auth.GetActiveUserClaims(c)
	recipe := m.Recipe{}
	err := m.DB.Where(&m.Recipe{UserUid: user.Uid, Uid: r.Uid}).
		Preload("Instructions").Preload("Ingredients").First(&recipe).Error
	if err != nil {
		gin_ext.AbortNotFound(c, fmt.Errorf("finding recipe -- %w", err))
		return
	}
	if r.Name != nil {
		recipe.Name = *r.Name
	}
	if r.IsTrashed != nil {
		recipe.IsTrashed = *r.IsTrashed
	}
	if r.Description != nil {
		recipe.Description = *r.Description
	}
	if r.Meal != nil {
		recipe.Meal = *r.Meal
	}
	if err = m.DB.Save(&recipe).Error; err != nil {
		gin_ext.ServerError(c, fmt.Errorf("saving recipe -- %w", err))
		return
	}
	sort.Slice(recipe.Instructions, func(i, j int) bool {
		return recipe.Instructions[i].SortOrder < recipe.Instructions[j].SortOrder
	})
	sort.Slice(recipe.Ingredients, func(i, j int) bool {
		return recipe.Ingredients[i].SortOrder < recipe.Ingredients[j].SortOrder
	})

	if len(r.Instructions) > 0 {
		if err = resizeInstructions(user.Uid, r.Instructions, &recipe); err != nil {
			gin_ext.ServerError(c, err)
			return
		}
	}

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
		Preload("Instructions").Preload("Ingredients").First(&recipe).Error
	if err != nil {
		gin_ext.ServerError(c, fmt.Errorf("listing recipes -- %w", err))
		return
	}

	c.JSON(http.StatusOK, recipe)
}

func ListRecipes(c *gin.Context) {
	user := auth.GetActiveUserClaims(c)
	var recipes []m.Recipe
	err := m.DB.Where("user_uid = ? and is_trashed = false", user.Uid).Find(&recipes).Error
	if err != nil {
		gin_ext.ServerError(c, fmt.Errorf("listing recipes -- %w", err))
		return
	}

	c.JSON(http.StatusOK, recipes)
}

func ListTrashedRecipes(c *gin.Context) {
	user := auth.GetActiveUserClaims(c)
	var recipes []m.Recipe
	err := m.DB.Where("user_uid = ? and is_trashed = true", user.Uid).Find(&recipes).Error
	if err != nil {
		gin_ext.ServerError(c, fmt.Errorf("listing trashed recipes -- %w", err))
		return
	}

	c.JSON(http.StatusOK, recipes)
}

func GetAvailableRecipes(c *gin.Context) {

	user := auth.GetActiveUserClaims(c)

	// Get inventory
	inventory, err := inventory.GetInventory(user.Uid)
	if err != nil {
		gin_ext.ServerError(c, fmt.Errorf("getting groceries -- %w", err))
		return
	}

	// Get all recipes that belong to this user
	var recipes []m.Recipe
	err = m.DB.Where("user_uid = ? and is_trashed = false", user.Uid).
		Preload("Instructions").Preload("Ingredients").
		Find(&recipes).Error
	if err != nil {
		gin_ext.ServerError(c, fmt.Errorf("getting recipes -- %w", err))
		return
	}

	c.JSON(http.StatusOK, filterRecipesByIngredients(inventory, recipes))
}
