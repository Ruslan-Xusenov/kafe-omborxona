package domain

import "time"

type TransactionType string

const (
	TransactionPurchase TransactionType = "purchase"
	TransactionReturn   TransactionType = "return"
	TransactionSale     TransactionType = "sale"
	TransactionWriteOff TransactionType = "write_off"
)

type Transaction struct {
	ID           int             `json:"id"`
	ProductID    int             `json:"product_id"`
	ProductName  string          `json:"product_name,omitempty"`
	SupplierID   *int            `json:"supplier_id,omitempty"`
	SupplierName string          `json:"supplier_name,omitempty"`
	UserID       int             `json:"user_id"`
	UserName     string          `json:"user_name,omitempty"`
	Type         TransactionType `json:"type"`
	Quantity     float64         `json:"quantity"`
	UnitPrice    float64         `json:"unit_price"`
	TotalAmount  float64         `json:"total_amount"`
	Note         string          `json:"note"`
	ExpiryDate   *string         `json:"expiry_date,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
}

type CreateTransactionRequest struct {
	ProductID  int             `json:"product_id"`
	SupplierID *int            `json:"supplier_id,omitempty"`
	Type       TransactionType `json:"type"`
	Quantity   float64         `json:"quantity"`
	UnitPrice  float64         `json:"unit_price"`
	Note       string          `json:"note"`
	ExpiryDate *string         `json:"expiry_date,omitempty"`
}

type TransactionFilter struct {
	Type      TransactionType `json:"type,omitempty"`
	ProductID int             `json:"product_id,omitempty"`
	DateFrom  string          `json:"date_from,omitempty"`
	DateTo    string          `json:"date_to,omitempty"`
}

type DashboardSummary struct {
	TotalPurchases    float64 `json:"total_purchases"`
	TotalReturns      float64 `json:"total_returns"`
	TotalSales        float64 `json:"total_sales"`
	TotalWriteOffs    float64 `json:"total_write_offs"`
	TotalProducts     int     `json:"total_products"`
	TotalCategories   int     `json:"total_categories"`
	TotalSuppliers    int     `json:"total_suppliers"`
	InventoryValue    float64 `json:"inventory_value"`
}

type ProfitReport struct {
	TotalSalesRevenue float64 `json:"total_sales_revenue"`
	TotalCostOfSold   float64 `json:"total_cost_of_sold"`
	TotalWriteOffLoss float64 `json:"total_write_off_loss"`
	NetProfit         float64 `json:"net_profit"`
	PeriodFrom        string  `json:"period_from"`
	PeriodTo          string  `json:"period_to"`
}

type TopProduct struct {
	ProductID   int     `json:"product_id"`
	ProductName string  `json:"product_name"`
	TotalSold   float64 `json:"total_sold"`
	TotalRevenue float64 `json:"total_revenue"`
}

type InventoryItem struct {
	ProductID    int     `json:"product_id"`
	ProductName  string  `json:"product_name"`
	Unit         string  `json:"unit"`
	CategoryName string  `json:"category_name"`
	CostPrice    float64 `json:"cost_price"`
	SalePrice    float64 `json:"sale_price"`
	CurrentStock float64 `json:"current_stock"`
	StockValue   float64 `json:"stock_value"`
}

type DashboardAlert struct {
	Type        string  `json:"type"` // "low_stock" | "expiring_soon"
	ProductID   int     `json:"product_id"`
	ProductName string  `json:"product_name"`
	Message     string  `json:"message"`
	Value       string  `json:"value"`
}

type TransactionRepository interface {
	GetAll(filter TransactionFilter) ([]Transaction, error)
	GetByID(id int) (*Transaction, error)
	Create(t *Transaction) error
	Delete(id int) error
	GetSummary(dateFrom, dateTo string) (*DashboardSummary, error)
	GetProfitReport(dateFrom, dateTo string) (*ProfitReport, error)
	GetTopProducts(limit int, dateFrom, dateTo string) ([]TopProduct, error)
	GetInventory() ([]InventoryItem, error)
	GetAlerts() ([]DashboardAlert, error)
}
