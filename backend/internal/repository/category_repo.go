package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"kafe-omborxona/internal/domain"
)

type CategoryRepo struct {
	db *pgxpool.Pool
}

func NewCategoryRepo(db *pgxpool.Pool) *CategoryRepo {
	return &CategoryRepo{db: db}
}

func (r *CategoryRepo) GetAll() ([]domain.Category, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, name, created_at FROM categories ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []domain.Category
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.CreatedAt); err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	return cats, nil
}

func (r *CategoryRepo) GetByID(id int) (*domain.Category, error) {
	var c domain.Category
	err := r.db.QueryRow(context.Background(),
		`SELECT id, name, created_at FROM categories WHERE id = $1`, id).
		Scan(&c.ID, &c.Name, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepo) Create(c *domain.Category) error {
	return r.db.QueryRow(context.Background(),
		`INSERT INTO categories (name) VALUES ($1) RETURNING id, created_at`,
		c.Name).Scan(&c.ID, &c.CreatedAt)
}

func (r *CategoryRepo) Update(c *domain.Category) error {
	_, err := r.db.Exec(context.Background(),
		`UPDATE categories SET name = $1 WHERE id = $2`, c.Name, c.ID)
	return err
}

func (r *CategoryRepo) Delete(id int) error {
	_, err := r.db.Exec(context.Background(), `DELETE FROM categories WHERE id = $1`, id)
	return err
}
