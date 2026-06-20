package domain

import "time"

type DebtStatus string

const (
	DebtUnpaid  DebtStatus = "unpaid"
	DebtPartial DebtStatus = "partial"
	DebtPaid    DebtStatus = "paid"
)

type Debt struct {
	ID            int        `json:"id"`
	SupplierID    int        `json:"supplier_id"`
	SupplierName  string     `json:"supplier_name,omitempty"`
	TransactionID int        `json:"transaction_id"`
	TotalDebt     float64    `json:"total_debt"`
	PaidAmount    float64    `json:"paid_amount"`
	Status        DebtStatus `json:"status"`
	DueDate       *string    `json:"due_date,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

type CreateDebtRequest struct {
	SupplierID    int        `json:"supplier_id"`
	TransactionID int        `json:"transaction_id"`
	TotalDebt     float64    `json:"total_debt"`
	DueDate       *string    `json:"due_date,omitempty"`
}

type PayDebtRequest struct {
	Amount float64 `json:"amount"`
}

type DebtRepository interface {
	GetAll() ([]Debt, error)
	GetByID(id int) (*Debt, error)
	Create(d *Debt) error
	Pay(id int, amount float64) error
	Delete(id int) error
}
