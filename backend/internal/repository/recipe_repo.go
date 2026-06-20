package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"kafe-omborxona/internal/domain"
)

type RecipeRepo struct {
	db *pgxpool.Pool
}

func NewRecipeRepo(db *pgxpool.Pool) *RecipeRepo {
	return &RecipeRepo{db: db}
}

func (r *RecipeRepo) GetAll() ([]domain.Recipe, error) {
	q := `SELECT r.id, r.product_id, COALESCE(p.name,''), r.name, r.created_at 
		FROM recipes r LEFT JOIN products p ON p.id = r.product_id ORDER BY r.name`
	rows, err := r.db.Query(context.Background(), q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []domain.Recipe
	for rows.Next() {
		var rec domain.Recipe
		if err := rows.Scan(&rec.ID, &rec.ProductID, &rec.ProductName, &rec.Name, &rec.CreatedAt); err != nil {
			return nil, err
		}
		
		// Load ingredients
		ings, err := r.getIngredients(rec.ID)
		if err == nil {
			rec.Ingredients = ings
		}
		recipes = append(recipes, rec)
	}
	return recipes, nil
}

func (r *RecipeRepo) getIngredients(recipeID int) ([]domain.RecipeIngredient, error) {
	q := `SELECT ri.id, ri.recipe_id, ri.ingredient_id, COALESCE(p.name,''), COALESCE(p.unit,''), ri.quantity 
		FROM recipe_ingredients ri LEFT JOIN products p ON p.id = ri.ingredient_id WHERE ri.recipe_id = $1`
	rows, err := r.db.Query(context.Background(), q, recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ings []domain.RecipeIngredient
	for rows.Next() {
		var ri domain.RecipeIngredient
		if err := rows.Scan(&ri.ID, &ri.RecipeID, &ri.IngredientID, &ri.IngredientName, &ri.Unit, &ri.Quantity); err != nil {
			return nil, err
		}
		ings = append(ings, ri)
	}
	return ings, nil
}

func (r *RecipeRepo) GetByID(id int) (*domain.Recipe, error) {
	var rec domain.Recipe
	q := `SELECT r.id, r.product_id, COALESCE(p.name,''), r.name, r.created_at 
		FROM recipes r LEFT JOIN products p ON p.id = r.product_id WHERE r.id = $1`
	err := r.db.QueryRow(context.Background(), q, id).Scan(&rec.ID, &rec.ProductID, &rec.ProductName, &rec.Name, &rec.CreatedAt)
	if err != nil {
		return nil, err
	}
	ings, err := r.getIngredients(rec.ID)
	if err == nil {
		rec.Ingredients = ings
	}
	return &rec, nil
}

func (r *RecipeRepo) GetByProductID(productID int) (*domain.Recipe, error) {
	var rec domain.Recipe
	q := `SELECT r.id, r.product_id, COALESCE(p.name,''), r.name, r.created_at 
		FROM recipes r LEFT JOIN products p ON p.id = r.product_id WHERE r.product_id = $1`
	err := r.db.QueryRow(context.Background(), q, productID).Scan(&rec.ID, &rec.ProductID, &rec.ProductName, &rec.Name, &rec.CreatedAt)
	if err != nil {
		return nil, err
	}
	ings, err := r.getIngredients(rec.ID)
	if err == nil {
		rec.Ingredients = ings
	}
	return &rec, nil
}

func (r *RecipeRepo) Create(rec *domain.Recipe) error {
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	err = tx.QueryRow(context.Background(), 
		`INSERT INTO recipes (product_id, name) VALUES ($1, $2) RETURNING id, created_at`,
		rec.ProductID, rec.Name).Scan(&rec.ID, &rec.CreatedAt)
	if err != nil {
		return err
	}

	for _, ing := range rec.Ingredients {
		_, err = tx.Exec(context.Background(), 
			`INSERT INTO recipe_ingredients (recipe_id, ingredient_id, quantity) VALUES ($1, $2, $3)`,
			rec.ID, ing.IngredientID, ing.Quantity)
		if err != nil {
			return err
		}
	}

	return tx.Commit(context.Background())
}

func (r *RecipeRepo) Update(rec *domain.Recipe) error {
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), `UPDATE recipes SET product_id=$1, name=$2 WHERE id=$3`, rec.ProductID, rec.Name, rec.ID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(context.Background(), `DELETE FROM recipe_ingredients WHERE recipe_id=$1`, rec.ID)
	if err != nil {
		return err
	}

	for _, ing := range rec.Ingredients {
		_, err = tx.Exec(context.Background(), 
			`INSERT INTO recipe_ingredients (recipe_id, ingredient_id, quantity) VALUES ($1, $2, $3)`,
			rec.ID, ing.IngredientID, ing.Quantity)
		if err != nil {
			return err
		}
	}

	return tx.Commit(context.Background())
}

func (r *RecipeRepo) Delete(id int) error {
	_, err := r.db.Exec(context.Background(), `DELETE FROM recipes WHERE id=$1`, id)
	return err
}
