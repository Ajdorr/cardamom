package models

import (
	m "cardamom/core/ext/math_ext"
	u "cardamom/core/ext/units"
	"time"
)

type MealType string

const (
	BREAKFAST MealType = "breakfast"
	LUNCH     MealType = "lunch"
	DINNER    MealType = "dinner"
	DESSERT   MealType = "desser"
)

var ValidMeals = []MealType{BREAKFAST, LUNCH, DINNER, DESSERT}

type Recipe struct {
	Uid          string             `gorm:"primaryKey;not null;default:null" json:"uid"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
	IsTrashed    bool               `gorm:"index;default:false" json:"is_trashed"`
	TrashAt      uint64             `json:"-"`
	UserUid      string             `gorm:"index" json:"user_uid"`
	Name         string             `json:"name"`
	Description  string             `json:"description"`
	Meal         MealType           `json:"meal"`
	Ingredients  []RecipeIngredient `gorm:"foreignKey:RecipeUid;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"ingredients"`
	Instructions string             `json:"instructions"`
}

type RecipeIngredient struct {
	Uid       string     `gorm:"primaryKey;not null;default:null" json:"uid"`
	RecipeUid string     `gorm:"index" json:"recipe_uid"`
	UserUid   string     `gorm:"index" json:"user_uid"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	SortOrder int        `json:"order"`
	Quantity  m.Rational `gorm:"decimal(30,2)" json:"quantity"`
	Unit      *u.Unit    `json:"unit"`
	Item      string     `gorm:"index" json:"item"`
	Optional  bool       `gorm:"default:false" json:"optional"`
	Modifier  string     `json:"modifier"`
}
