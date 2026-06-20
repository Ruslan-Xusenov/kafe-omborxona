package domain

import "time"

type RecipeIngredient struct {
	ID             int     `json:"id"`
	RecipeID       int     `json:"recipe_id"`
	IngredientID   int     `json:"ingredient_id"`
	IngredientName string  `json:"ingredient_name,omitempty"`
	Unit           string  `json:"unit,omitempty"`
	Quantity       float64 `json:"quantity"`
}

type Recipe struct {
	ID          int                `json:"id"`
	ProductID   int                `json:"product_id"`
	ProductName string             `json:"product_name,omitempty"`
	Name        string             `json:"name"`
	Ingredients []RecipeIngredient `json:"ingredients"`
	CreatedAt   time.Time          `json:"created_at"`
}

type CreateRecipeRequest struct {
	ProductID   int                `json:"product_id"`
	Name        string             `json:"name"`
	Ingredients []RecipeIngredient `json:"ingredients"`
}

type RecipeRepository interface {
	GetAll() ([]Recipe, error)
	GetByProductID(productID int) (*Recipe, error)
	GetByID(id int) (*Recipe, error)
	Create(r *Recipe) error
	Update(r *Recipe) error
	Delete(id int) error
}
