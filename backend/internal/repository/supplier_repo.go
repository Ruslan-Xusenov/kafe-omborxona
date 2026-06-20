package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"kafe-omborxona/internal/domain"
)

type SupplierRepo struct {
	db *pgxpool.Pool
}

func NewSupplierRepo(db *pgxpool.Pool) *SupplierRepo {
	return &SupplierRepo{db: db}
}

func (r *SupplierRepo) GetAll() ([]domain.Supplier, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, name, phone, address, created_at FROM suppliers ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sups []domain.Supplier
	for rows.Next() {
		var s domain.Supplier
		if err := rows.Scan(&s.ID, &s.Name, &s.Phone, &s.Address, &s.CreatedAt); err != nil {
			return nil, err
		}
		sups = append(sups, s)
	}
	return sups, nil
}

func (r *SupplierRepo) GetByID(id int) (*domain.Supplier, error) {
	var s domain.Supplier
	err := r.db.QueryRow(context.Background(),
		`SELECT id, name, phone, address, created_at FROM suppliers WHERE id = $1`, id).
		Scan(&s.ID, &s.Name, &s.Phone, &s.Address, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SupplierRepo) Create(s *domain.Supplier) error {
	return r.db.QueryRow(context.Background(),
		`INSERT INTO suppliers (name, phone, address) VALUES ($1, $2, $3) RETURNING id, created_at`,
		s.Name, s.Phone, s.Address).Scan(&s.ID, &s.CreatedAt)
}

func (r *SupplierRepo) Update(s *domain.Supplier) error {
	_, err := r.db.Exec(context.Background(),
		`UPDATE suppliers SET name = $1, phone = $2, address = $3 WHERE id = $4`,
		s.Name, s.Phone, s.Address, s.ID)
	return err
}

func (r *SupplierRepo) Delete(id int) error {
	_, err := r.db.Exec(context.Background(), `DELETE FROM suppliers WHERE id = $1`, id)
	return err
}
