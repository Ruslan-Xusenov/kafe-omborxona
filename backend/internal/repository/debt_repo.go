package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"kafe-omborxona/internal/domain"
)

type DebtRepo struct {
	db *pgxpool.Pool
}

func NewDebtRepo(db *pgxpool.Pool) *DebtRepo {
	return &DebtRepo{db: db}
}

func (r *DebtRepo) GetAll() ([]domain.Debt, error) {
	q := `SELECT d.id, d.supplier_id, COALESCE(s.name,''), d.transaction_id, d.total_debt, d.paid_amount, d.status, d.due_date, d.created_at
		FROM debts d LEFT JOIN suppliers s ON s.id = d.supplier_id ORDER BY d.created_at DESC`
	rows, err := r.db.Query(context.Background(), q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var debts []domain.Debt
	for rows.Next() {
		var d domain.Debt
		if err := rows.Scan(&d.ID, &d.SupplierID, &d.SupplierName, &d.TransactionID, &d.TotalDebt, &d.PaidAmount, &d.Status, &d.DueDate, &d.CreatedAt); err != nil {
			return nil, err
		}
		debts = append(debts, d)
	}
	return debts, nil
}

func (r *DebtRepo) GetByID(id int) (*domain.Debt, error) {
	var d domain.Debt
	q := `SELECT d.id, d.supplier_id, COALESCE(s.name,''), d.transaction_id, d.total_debt, d.paid_amount, d.status, d.due_date, d.created_at
		FROM debts d LEFT JOIN suppliers s ON s.id = d.supplier_id WHERE d.id = $1`
	err := r.db.QueryRow(context.Background(), q, id).Scan(&d.ID, &d.SupplierID, &d.SupplierName, &d.TransactionID, &d.TotalDebt, &d.PaidAmount, &d.Status, &d.DueDate, &d.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *DebtRepo) Create(d *domain.Debt) error {
	d.Status = domain.DebtUnpaid
	d.PaidAmount = 0
	return r.db.QueryRow(context.Background(),
		`INSERT INTO debts (supplier_id, transaction_id, total_debt, due_date) VALUES ($1, $2, $3, $4) RETURNING id, created_at`,
		d.SupplierID, d.TransactionID, d.TotalDebt, d.DueDate).Scan(&d.ID, &d.CreatedAt)
}

func (r *DebtRepo) Pay(id int, amount float64) error {
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	var total, paid float64
	err = tx.QueryRow(context.Background(), `SELECT total_debt, paid_amount FROM debts WHERE id=$1 FOR UPDATE`, id).Scan(&total, &paid)
	if err != nil {
		return err
	}

	newPaid := paid + amount
	status := domain.DebtPartial
	if newPaid >= total {
		status = domain.DebtPaid
		newPaid = total
	}

	_, err = tx.Exec(context.Background(), `UPDATE debts SET paid_amount=$1, status=$2 WHERE id=$3`, newPaid, status, id)
	if err != nil {
		return err
	}

	return tx.Commit(context.Background())
}

func (r *DebtRepo) Delete(id int) error {
	_, err := r.db.Exec(context.Background(), `DELETE FROM debts WHERE id=$1`, id)
	return err
}
