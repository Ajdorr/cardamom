package recipe

import (
	"cardamom/core/ext/log_ext"
	m "cardamom/core/ext/math_ext"
	u "cardamom/core/ext/units"
	"cardamom/core/models"
	"strings"

	"github.com/thoas/go-funk"
)

type IngredientPart struct {
	Quantity m.Rational `json:"quantity"`
	Unit     *u.Unit    `json:"unit"`
	Item     string     `json:"item"`
	Optional bool       `json:"optional"`
	Modifier string     `json:"modifier"`
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

	for i, ingre := range req.Ingredients {
		req.Ingredients[i].Item = strings.TrimSpace(strings.ToLower(ingre.Item))
		req.Ingredients[i].Modifier = strings.TrimSpace(strings.ToLower(ingre.Modifier))
	}

	return "", nil
}

type UpdateRecipeRequest struct {
	Uid          string           `json:"uid"`
	IsTrashed    *bool            `json:"is_trashed,omitempty"`
	Name         *string          `json:"name,omitempty"`
	Description  *string          `json:"description,omitempty"`
	Meal         *models.MealType `json:"meal,omitempty"`
	Instructions *string          `json:"instructions,omitempty"`
}

func (req *UpdateRecipeRequest) Validate() (string, error) {

	req.Uid = strings.TrimSpace(req.Uid)
	if len(req.Uid) == 0 {
		return log_ext.ReturnBoth("uid must not be empty in request")
	}

	if req.Name != nil {
		*req.Name = strings.TrimSpace(*req.Name)
		if len(*req.Name) == 0 {
			return log_ext.ReturnBoth("name must not be empty")
		}
	}

	return "", nil
}

type CreateRecipeIngredientRequest struct {
	RecipeUid string     `json:"recipe_uid"`
	Quantity  m.Rational `json:"quantity,omitempty"`
	Unit      *u.Unit    `json:"unit,omitempty"`
	Item      string     `json:"item,omitempty"`
	SortOrder int        `json:"order"`
	Optional  bool       `json:"optional,omitempty"`
	Modifier  string     `json:"modifier,omitempty"`
}

func (req *CreateRecipeIngredientRequest) Validate() (string, error) {

	req.RecipeUid = strings.TrimSpace(req.RecipeUid)
	if len(req.RecipeUid) == 0 {
		return log_ext.ReturnBoth("uid must not be empty in request")
	}
	req.Item = strings.TrimSpace(strings.ToLower(req.Item))
	req.Modifier = strings.TrimSpace(strings.ToLower(req.Modifier))

	return "", nil
}

type UpdateRecipeIngredientRequest struct {
	Uid      string      `json:"uid"`
	Quantity *m.Rational `json:"quantity,omitempty"`
	Unit     *u.Unit     `json:"unit,omitempty"`
	Item     *string     `json:"item,omitempty"`
	Optional *bool       `json:"optional,omitempty"`
	Modifier *string     `json:"modifier,omitempty"`
}

func (req *UpdateRecipeIngredientRequest) Validate() (string, error) {

	req.Uid = strings.TrimSpace(req.Uid)
	if len(req.Uid) == 0 {
		return log_ext.ReturnBoth("uid must not be empty in request")
	}

	if req.Item != nil {
		*req.Item = strings.TrimSpace(strings.ToLower(*req.Item))
	}

	if req.Modifier != nil {
		*req.Modifier = strings.TrimSpace(strings.ToLower(*req.Modifier))
	}

	return "", nil
}

type ReorderRecipeIngredientRequest struct {
	Uid         string   `json:"uid"`
	Ingredients []string `json:"ingredient_uids"`
}

func (req *ReorderRecipeIngredientRequest) Validate() (string, error) {

	req.Uid = strings.TrimSpace(req.Uid)
	if len(req.Uid) == 0 {
		return log_ext.ReturnBoth("uid must not be empty in request")
	}
	return "", nil
}

type SearchRecipeRequest struct {
	Name        *string          `json:"name,omitempty"`
	Description *string          `json:"description,omitempty"`
	Meal        *models.MealType `json:"meal,omitempty"`
	Ingredient  *string          `json:"ingredient,omitempty"`
}

func (req *SearchRecipeRequest) Validate() (string, error) {

	if req.Name == nil && req.Description == nil && req.Meal == nil && req.Ingredient == nil {
		return log_ext.ReturnBoth("must have at least one search criteria")
	}

	return "", nil
}

type GetAvailableRecipeRequest struct {
	Meal *models.MealType `json:"meal,omitempty"`
}

func (req *GetAvailableRecipeRequest) Validate() (string, error) {

	if req.Meal != nil && !funk.Contains(models.ValidMeals, *req.Meal) {
		return log_ext.ReturnBoth("invalid meal type")
	}

	return "", nil
}
