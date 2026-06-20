package domain

import "time"

type Supplier struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateSupplierRequest struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

type SupplierRepository interface {
	GetAll() ([]Supplier, error)
	GetByID(id int) (*Supplier, error)
	Create(sup *Supplier) error
	Update(sup *Supplier) error
	Delete(id int) error
}
