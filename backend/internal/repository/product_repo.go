package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"kafe-omborxona/internal/domain"
)

type ProductRepo struct {
	db *pgxpool.Pool
}

func NewProductRepo(db *pgxpool.Pool) *ProductRepo {
	return &ProductRepo{db: db}
}

func (r *ProductRepo) GetAll() ([]domain.Product, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT p.id, p.name, p.unit, p.category_id, COALESCE(c.name,''), p.cost_price, p.sale_price, p.min_stock, p.barcode, p.created_at
		 FROM products p LEFT JOIN categories c ON c.id = p.category_id ORDER BY p.name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Unit, &p.CategoryID, &p.CategoryName, &p.CostPrice, &p.SalePrice, &p.MinStock, &p.Barcode, &p.CreatedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *ProductRepo) GetByID(id int) (*domain.Product, error) {
	var p domain.Product
	err := r.db.QueryRow(context.Background(),
		`SELECT p.id, p.name, p.unit, p.category_id, COALESCE(c.name,''), p.cost_price, p.sale_price, p.min_stock, p.barcode, p.created_at
		 FROM products p LEFT JOIN categories c ON c.id = p.category_id WHERE p.id = $1`, id).
		Scan(&p.ID, &p.Name, &p.Unit, &p.CategoryID, &p.CategoryName, &p.CostPrice, &p.SalePrice, &p.MinStock, &p.Barcode, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepo) GetWithStock() ([]domain.Product, error) {
	q := `SELECT p.id, p.name, p.unit, p.category_id, COALESCE(c.name,''), p.cost_price, p.sale_price, p.min_stock, p.barcode, p.created_at,
		COALESCE(SUM(CASE WHEN t.type='purchase' THEN t.quantity ELSE 0 END),0)
		- COALESCE(SUM(CASE WHEN t.type='return' THEN t.quantity ELSE 0 END),0)
		- COALESCE(SUM(CASE WHEN t.type='sale' THEN t.quantity ELSE 0 END),0)
		- COALESCE(SUM(CASE WHEN t.type='write_off' THEN t.quantity ELSE 0 END),0) AS stock
	FROM products p
	LEFT JOIN categories c ON c.id = p.category_id
	LEFT JOIN transactions t ON t.product_id = p.id
	GROUP BY p.id, c.name ORDER BY p.name`

	rows, err := r.db.Query(context.Background(), q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Unit, &p.CategoryID, &p.CategoryName, &p.CostPrice, &p.SalePrice, &p.MinStock, &p.Barcode, &p.CreatedAt, &p.CurrentStock); err != nil {
			return nil, err
		}
		p.StockValue = p.CurrentStock * p.CostPrice
		products = append(products, p)
	}
	return products, nil
}

func (r *ProductRepo) Create(p *domain.Product) error {
	return r.db.QueryRow(context.Background(),
		`INSERT INTO products (name,unit,category_id,cost_price,sale_price,min_stock,barcode) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id,created_at`,
		p.Name, p.Unit, p.CategoryID, p.CostPrice, p.SalePrice, p.MinStock, p.Barcode).Scan(&p.ID, &p.CreatedAt)
}

func (r *ProductRepo) Update(p *domain.Product) error {
	_, err := r.db.Exec(context.Background(),
		`UPDATE products SET name=$1,unit=$2,category_id=$3,cost_price=$4,sale_price=$5,min_stock=$6,barcode=$7 WHERE id=$8`,
		p.Name, p.Unit, p.CategoryID, p.CostPrice, p.SalePrice, p.MinStock, p.Barcode, p.ID)
	return err
}

func (r *ProductRepo) Delete(id int) error {
	_, err := r.db.Exec(context.Background(), `DELETE FROM products WHERE id=$1`, id)
	return err
}
