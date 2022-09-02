package recipe

import (
	"cardamom/core/ext/gin_ext"
	"cardamom/core/models"
	"cardamom/core/services"
	"cardamom/core/services/auth"
	"cardamom/core/services/inventory"
	"fmt"
	"math/rand"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
)

func CreateRecipe(c *gin.Context, r *CreateRecipeRequest) {
	user := auth.GetActiveUserClaims(c)
	recipe := models.Recipe{
		Uid:         generateRecipeUid(),
		UserUid:     user.Uid,
		Name:        r.Name,
		Description: r.Description,
		Meal:        r.Meal,
	}

	if err := models.DB.Create(&recipe).Error; err != nil {
		gin_ext.ServerError(c, fmt.Errorf("unable to create recipe -- %w", err))
		return
	}

	recipe.Instructions = make([]models.RecipeInstruction, len(r.Instructions))
	for i, instr := range r.Instructions {
		recipe.Instructions[i] = models.RecipeInstruction{
			Uid:       generateInstrUid(),
			UserUid:   user.Uid,
			RecipeUid: recipe.Uid,
			Meal:      r.Meal,
			SortOrder: i,
			Text:      instr,
		}
		err := models.DB.Create(&recipe.Instructions[i]).Error
		if err != nil {
			gin_ext.ServerError(c, fmt.Errorf("unable to create instruction -- %w", err))
			return
		}
	}

	recipe.Ingredients = make([]models.RecipeIngredient, len(r.Ingredients))
	for i, ingre := range r.Ingredients {
		recipe.Ingredients[i] = models.RecipeIngredient{
			Uid:       generateIngreUid(),
			UserUid:   user.Uid,
			RecipeUid: recipe.Uid,
			Meal:      r.Meal,
			SortOrder: i,
			Quantity:  ingre.Quantity,
			Unit:      ingre.Unit,
			Item:      ingre.Item,
		}
		err := models.DB.Create(&recipe.Ingredients[i]).Error
		if err != nil {
			gin_ext.ServerError(c, fmt.Errorf("unable to create ingredient -- %w", err))
			return
		}
	}

	c.JSON(http.StatusCreated, recipe)
}

func UpdateRecipe(c *gin.Context, r *UpdateRecipeRequest) {
	user := auth.GetActiveUserClaims(c)
	recipe := models.Recipe{}
	err := models.DB.Where(&models.Recipe{UserUid: user.Uid, Uid: r.Uid}).
		Preload("Instructions").Preload("Ingredients").First(&recipe).Error
	if err != nil {
		gin_ext.AbortNotFound(c, fmt.Errorf("finding recipe -- %w", err))
		return
	}
	if r.Name != nil {
		recipe.Name = *r.Name
	}
	if r.Description != nil {
		recipe.Description = *r.Description
	}
	if r.Meal != nil {
		recipe.Meal = *r.Meal
	}
	if err = models.DB.Save(&recipe).Error; err != nil {
		gin_ext.ServerError(c, fmt.Errorf("saving recipe -- %w", err))
		return
	}
	sort.Slice(recipe.Instructions, func(i, j int) bool {
		return recipe.Instructions[i].SortOrder < recipe.Instructions[j].SortOrder
	})
	sort.Slice(recipe.Ingredients, func(i, j int) bool {
		return recipe.Ingredients[i].SortOrder < recipe.Ingredients[j].SortOrder
	})

	if len(r.Instructions) < len(recipe.Instructions) {
		for i := len(r.Instructions); i < len(recipe.Instructions); i++ {
			if err = models.DB.Delete(recipe.Instructions[i]).Error; err != nil {
				gin_ext.ServerError(c, fmt.Errorf("deleting instructions -- %w", err))
				return
			}
		}

		recipe.Instructions = recipe.Instructions[:len(r.Instructions)]
	}
	for i, instr := range r.Instructions {
		if i >= len(recipe.Instructions) {
			recipe.Instructions = append(recipe.Instructions, models.RecipeInstruction{
				Uid:       generateInstrUid(),
				UserUid:   user.Uid,
				RecipeUid: recipe.Uid,
				Meal:      recipe.Meal,
				SortOrder: i,
				Text:      instr,
			})
		} else {
			recipe.Instructions[i].Meal = recipe.Meal
			recipe.Instructions[i].Text = instr
			recipe.Instructions[i].SortOrder = i
		}

		if err = models.DB.Save(&recipe.Instructions[i]).Error; err != nil {
			gin_ext.ServerError(c, fmt.Errorf("updating instruction(%d) -- %w", i, err))
			return
		}
	}

	if len(r.Ingredients) < len(recipe.Ingredients) {
		for i := len(r.Ingredients); i < len(recipe.Ingredients); i++ {
			if err = models.DB.Delete(recipe.Ingredients[i]).Error; err != nil {
				gin_ext.ServerError(c, fmt.Errorf("deleting instructions -- %w", err))
				return
			}
		}

		recipe.Ingredients = recipe.Ingredients[:len(r.Ingredients)]
	}
	for i, ingre := range r.Ingredients {
		if i >= len(recipe.Ingredients) {
			recipe.Ingredients = append(recipe.Ingredients, models.RecipeIngredient{
				Uid:       generateIngreUid(),
				UserUid:   user.Uid,
				RecipeUid: recipe.Uid,
				Meal:      recipe.Meal,
				SortOrder: i,
				Quantity:  ingre.Quantity,
				Unit:      ingre.Unit,
				Item:      ingre.Item,
			})
		} else {
			recipe.Ingredients[i].Meal = recipe.Meal
			recipe.Ingredients[i].Quantity = ingre.Quantity
			recipe.Ingredients[i].Unit = ingre.Unit
			recipe.Ingredients[i].Item = ingre.Item
			recipe.Ingredients[i].SortOrder = i
		}

		if err = models.DB.Save(&recipe.Ingredients[i]).Error; err != nil {
			gin_ext.ServerError(c, fmt.Errorf("updating ingredient(%d) -- %w", i, err))
			return
		}
	}

	c.JSON(http.StatusOK, recipe)
}

func ReadRecipe(c *gin.Context, r *services.ReadRequest) {
	user := auth.GetActiveUserClaims(c)
	recipe := models.Recipe{}
	err := models.DB.Where(&models.Recipe{Uid: r.Uid, UserUid: user.Uid}).
		Preload("Instructions").Preload("Ingredients").First(&recipe).Error
	if err != nil {
		gin_ext.ServerError(c, fmt.Errorf("listing recipes -- %w", err))
		return
	}

	c.JSON(http.StatusOK, recipe)
}

func ListRecipes(c *gin.Context) {
	user := auth.GetActiveUserClaims(c)
	var recipes []models.Recipe
	err := models.DB.Where(&models.Recipe{UserUid: user.Uid}).Find(&recipes).Error
	if err != nil {
		gin_ext.ServerError(c, fmt.Errorf("listing recipes -- %w", err))
		return
	}

	c.JSON(http.StatusOK, recipes)
}

func GetRandomAvailableRecipes(c *gin.Context) {
	user := auth.GetActiveUserClaims(c)
	inventory, err := inventory.GetInventory(user.Uid)
	if err != nil {
		gin_ext.ServerError(c, fmt.Errorf("getting items -- %w", err))
		return
	}
	items := funk.Map(inventory, func(i models.InventoryItem) string { return i.Item }).([]string)

	randItem := items[rand.Intn(len(items))]
	var ingredients []models.RecipeIngredient
	err = models.DB.
		Where(&models.RecipeIngredient{UserUid: user.Uid, Item: randItem}).
		Find(&ingredients).Error
	if err != nil {
		gin_ext.ServerError(c, fmt.Errorf("finding recipe ingredients -- %w", err))
		return
	}

	recipes := make([]models.Recipe, len(ingredients))
	for i, ingre := range ingredients {
		err = models.DB.
			Where(&models.Recipe{UserUid: user.Uid, Uid: ingre.RecipeUid}).
			Preload("Instructions").Preload("Ingredients").
			Find(&recipes[i]).Error
		if err != nil {
			gin_ext.ServerError(c, fmt.Errorf("finding related recipes -- %w", err))
			return
		}
	}

	filteredRecipes := funk.Filter(recipes, func(r models.Recipe) bool {
		for _, ingre := range r.Ingredients {
			if !funk.Contains(items, ingre.Item) {
				return false
			}
		}
		return true
	}).([]models.Recipe)
	rand.Shuffle(len(filteredRecipes), func(i, j int) {
		filteredRecipes[i], filteredRecipes[j] = filteredRecipes[j], filteredRecipes[i]
	})

	c.JSON(http.StatusOK, filteredRecipes)
}
