package domain

import "time"

type Category struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

type CategoryRepository interface {
	GetAll() ([]Category, error)
	GetByID(id int) (*Category, error)
	Create(cat *Category) error
	Update(cat *Category) error
	Delete(id int) error
}
