package domain

import "time"

type Product struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Unit         string    `json:"unit"`
	CategoryID   int       `json:"category_id"`
	CategoryName string    `json:"category_name,omitempty"`
	CostPrice    float64   `json:"cost_price"`
	SalePrice    float64   `json:"sale_price"`
	MinStock     float64   `json:"min_stock"`
	Barcode      *string   `json:"barcode,omitempty"`
	CurrentStock float64   `json:"current_stock,omitempty"`
	StockValue   float64   `json:"stock_value,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateProductRequest struct {
	Name       string  `json:"name"`
	Unit       string  `json:"unit"`
	CategoryID int     `json:"category_id"`
	CostPrice  float64 `json:"cost_price"`
	SalePrice  float64 `json:"sale_price"`
	MinStock   float64 `json:"min_stock"`
	Barcode    *string `json:"barcode,omitempty"`
}

type ProductRepository interface {
	GetAll() ([]Product, error)
	GetByID(id int) (*Product, error)
	GetWithStock() ([]Product, error)
	Create(p *Product) error
	Update(p *Product) error
	Delete(id int) error
}
