package recipe

import (
	"cardamom/core/ext/log_ext"
	m "cardamom/core/ext/math_ext"
	u "cardamom/core/ext/units"
	"cardamom/core/models"
	"strings"
)

type IngredientPart struct {
	Quantity m.Rational `json:"quantity"`
	Unit     *u.Unit    `json:"unit"`
	Item     string     `json:"item"`
}

type CreateRecipeRequest struct {
	Name         string           `json:"name"`
	Description  string           `json:"description,omitempty"`
	Meal         models.MealType  `json:"meal"`
	Ingredients  []IngredientPart `json:"ingredients"`
	Instructions string           `json:"instructions"`
}

func (req *CreateRecipeRequest) Validate() (string, error) {

	req.Name = strings.TrimSpace(req.Name)
	if len(req.Name) == 0 {
		return log_ext.ReturnBoth("name must not be empty in request")
	}

	if len(req.Ingredients) == 0 {
		return log_ext.ReturnBoth("ingredients must not be empty in request")
	}

	if len(req.Instructions) == 0 {
		return log_ext.ReturnBoth("instructions must not be empty in request")
	}

	return "", nil
}

type UpdateRecipeRequest struct {
	Uid          string           `json:"uid"`
	IsTrashed    *bool            `json:"is_trashed,omitempty"`
	Name         *string          `json:"name,omitempty"`
	Description  *string          `json:"description,omitempty"`
	Meal         *models.MealType `json:"meal,omitempty"`
	Ingredients  []IngredientPart `json:"ingredients,omitempty"`
	Instructions *string          `json:"instructions,omitempty"`
}

func (req *UpdateRecipeRequest) Validate() (string, error) {

	req.Uid = strings.TrimSpace(req.Uid)
	if len(req.Uid) == 0 {
		return log_ext.ReturnBoth("name must not be empty in request")
	}

	if req.Name != nil {
		*req.Name = strings.TrimSpace(*req.Name)
		if len(*req.Name) == 0 {
			return log_ext.ReturnBoth("name must not be empty")
		}
	}

	return "", nil
}

type UpdateRecipeResponse struct {
	Uid          string           `json:"uid"`
	Name         string           `json:"name,omitempty"`
	Description  string           `json:"description,omitempty"`
	Meal         models.MealType  `json:"meal,omitempty"`
	Ingredients  []IngredientPart `json:"ingredients,omitempty"`
	Instructions []string         `json:"instructions,omitempty"`
}

type SearchRecipeRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Meal        *string `json:"meal,omitempty"`
	Ingredient  *string `json:"ingredient,omitempty"`
}

func (req *SearchRecipeRequest) Validate() (string, error) {

	if req.Name == nil && req.Description == nil && req.Meal == nil && req.Ingredient == nil {
		return log_ext.ReturnBoth("must have at least one search criteria")
	}

	return "", nil
}
