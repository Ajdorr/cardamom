package models

type GroceryListItem struct {
	ID          string `json:"id"`
	CreatedAt   string `json:"created_at"`
	UpdateAt    string `json:"update_at"`
	Value       string `json:"value"`
	UserID      string `json:"user_id"`
	IsCollected bool   `json:"is_collected"`
}

type InventoryItem struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UserID    string `json:"user_id"`
	Value     string `json:"value"`
}

type Recipe struct {
	ID          string `json:"id"`
	CreatedAt   string `json:"created_at"`
	UpdateAt    string `json:"update_at"`
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"desc"`
}

type RecipeInstruction struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdateAt  string `json:"update_at"`
	RecipieID string `json:"recipie_id"`
	SortOrder string `json:"order"`
	Value     string `json:"value"`
}

type RecipeIngredient struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdateAt  string `json:"update_at"`
	RecipieID string `json:"recipie_id"`
	SortOrder string `json:"order"`
	Value     string `json:"value"`
}
