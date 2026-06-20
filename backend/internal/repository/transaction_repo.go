package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"kafe-omborxona/internal/domain"
)

type TransactionRepo struct {
	db *pgxpool.Pool
}

func NewTransactionRepo(db *pgxpool.Pool) *TransactionRepo {
	return &TransactionRepo{db: db}
}

func (r *TransactionRepo) GetAll(f domain.TransactionFilter) ([]domain.Transaction, error) {
	q := `SELECT t.id, t.product_id, COALESCE(p.name,''), t.supplier_id, COALESCE(s.name,''),
		t.user_id, COALESCE(u.full_name,''), t.type, t.quantity, t.unit_price, t.total_amount, COALESCE(t.note,''), t.expiry_date, t.created_at
		FROM transactions t
		LEFT JOIN products p ON p.id = t.product_id
		LEFT JOIN suppliers s ON s.id = t.supplier_id
		LEFT JOIN users u ON u.id = t.user_id`

	var conditions []string
	var args []interface{}
	i := 1

	if f.Type != "" {
		conditions = append(conditions, fmt.Sprintf("t.type = $%d", i))
		args = append(args, f.Type)
		i++
	}
	if f.ProductID > 0 {
		conditions = append(conditions, fmt.Sprintf("t.product_id = $%d", i))
		args = append(args, f.ProductID)
		i++
	}
	if f.DateFrom != "" {
		conditions = append(conditions, fmt.Sprintf("t.created_at >= $%d", i))
		args = append(args, f.DateFrom)
		i++
	}
	if f.DateTo != "" {
		conditions = append(conditions, fmt.Sprintf("t.created_at <= $%d::date + interval '1 day'", i))
		args = append(args, f.DateTo)
		i++
	}

	if len(conditions) > 0 {
		q += " WHERE " + strings.Join(conditions, " AND ")
	}
	q += " ORDER BY t.created_at DESC LIMIT 500"

	rows, err := r.db.Query(context.Background(), q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txns []domain.Transaction
	for rows.Next() {
		var t domain.Transaction
		if err := rows.Scan(&t.ID, &t.ProductID, &t.ProductName, &t.SupplierID, &t.SupplierName,
			&t.UserID, &t.UserName, &t.Type, &t.Quantity, &t.UnitPrice, &t.TotalAmount, &t.Note, &t.ExpiryDate, &t.CreatedAt); err != nil {
			return nil, err
		}
		txns = append(txns, t)
	}
	return txns, nil
}

func (r *TransactionRepo) GetByID(id int) (*domain.Transaction, error) {
	var t domain.Transaction
	err := r.db.QueryRow(context.Background(),
		`SELECT t.id, t.product_id, COALESCE(p.name,''), t.supplier_id, COALESCE(s.name,''),
		t.user_id, COALESCE(u.full_name,''), t.type, t.quantity, t.unit_price, t.total_amount, COALESCE(t.note,''), t.expiry_date, t.created_at
		FROM transactions t
		LEFT JOIN products p ON p.id = t.product_id
		LEFT JOIN suppliers s ON s.id = t.supplier_id
		LEFT JOIN users u ON u.id = t.user_id
		WHERE t.id = $1`, id).
		Scan(&t.ID, &t.ProductID, &t.ProductName, &t.SupplierID, &t.SupplierName,
			&t.UserID, &t.UserName, &t.Type, &t.Quantity, &t.UnitPrice, &t.TotalAmount, &t.Note, &t.ExpiryDate, &t.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TransactionRepo) Create(t *domain.Transaction) error {
	t.TotalAmount = t.Quantity * t.UnitPrice

	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	err = tx.QueryRow(context.Background(),
		`INSERT INTO transactions (product_id,supplier_id,user_id,type,quantity,unit_price,total_amount,note,expiry_date)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id,created_at`,
		t.ProductID, t.SupplierID, t.UserID, t.Type, t.Quantity, t.UnitPrice, t.TotalAmount, t.Note, t.ExpiryDate).
		Scan(&t.ID, &t.CreatedAt)
	if err != nil {
		return err
	}

	// Auto-deduct recipe ingredients if this is a sale
	if t.Type == domain.TransactionSale {
		q := `SELECT ri.ingredient_id, ri.quantity, p.cost_price 
			  FROM recipe_ingredients ri 
			  JOIN recipes rec ON rec.id = ri.recipe_id
			  JOIN products p ON p.id = ri.ingredient_id
			  WHERE rec.product_id = $1`
		rows, err := tx.Query(context.Background(), q, t.ProductID)
		if err == nil {
			type ingData struct {
				ID    int
				Qty   float64
				Price float64
			}
			var ings []ingData
			for rows.Next() {
				var d ingData
				if err := rows.Scan(&d.ID, &d.Qty, &d.Price); err == nil {
					ings = append(ings, d)
				}
			}
			rows.Close()

			for _, ing := range ings {
				writeOffQty := ing.Qty * t.Quantity
				writeOffAmount := writeOffQty * ing.Price
				note := fmt.Sprintf("Retsept: %d ID tranzaksiya uchun", t.ID)
				
				_, err = tx.Exec(context.Background(),
					`INSERT INTO transactions (product_id,user_id,type,quantity,unit_price,total_amount,note)
					 VALUES ($1,$2,'write_off',$3,$4,$5,$6)`,
					ing.ID, t.UserID, writeOffQty, ing.Price, writeOffAmount, note)
				if err != nil {
					return err
				}
			}
		}
	}

	return tx.Commit(context.Background())
}

func (r *TransactionRepo) Delete(id int) error {
	_, err := r.db.Exec(context.Background(), `DELETE FROM transactions WHERE id=$1`, id)
	return err
}

func (r *TransactionRepo) GetSummary(dateFrom, dateTo string) (*domain.DashboardSummary, error) {
	var s domain.DashboardSummary
	dateFilter := buildDateFilter(dateFrom, dateTo)

	q := fmt.Sprintf(`SELECT
		COALESCE(SUM(CASE WHEN type='purchase' THEN total_amount ELSE 0 END),0),
		COALESCE(SUM(CASE WHEN type='return' THEN total_amount ELSE 0 END),0),
		COALESCE(SUM(CASE WHEN type='sale' THEN total_amount ELSE 0 END),0),
		COALESCE(SUM(CASE WHEN type='write_off' THEN total_amount ELSE 0 END),0)
		FROM transactions %s`, dateFilter)

	err := r.db.QueryRow(context.Background(), q).
		Scan(&s.TotalPurchases, &s.TotalReturns, &s.TotalSales, &s.TotalWriteOffs)
	if err != nil {
		return nil, err
	}

	r.db.QueryRow(context.Background(), `SELECT COUNT(*) FROM products`).Scan(&s.TotalProducts)
	r.db.QueryRow(context.Background(), `SELECT COUNT(*) FROM categories`).Scan(&s.TotalCategories)
	r.db.QueryRow(context.Background(), `SELECT COUNT(*) FROM suppliers`).Scan(&s.TotalSuppliers)

	// Inventory value
	r.db.QueryRow(context.Background(),
		`SELECT COALESCE(SUM(stock * cost_price),0) FROM (
			SELECT p.cost_price,
			COALESCE(SUM(CASE WHEN t.type='purchase' THEN t.quantity ELSE 0 END),0)
			- COALESCE(SUM(CASE WHEN t.type='return' THEN t.quantity ELSE 0 END),0)
			- COALESCE(SUM(CASE WHEN t.type='sale' THEN t.quantity ELSE 0 END),0)
			- COALESCE(SUM(CASE WHEN t.type='write_off' THEN t.quantity ELSE 0 END),0) AS stock
			FROM products p LEFT JOIN transactions t ON t.product_id = p.id
			GROUP BY p.id) sub`).Scan(&s.InventoryValue)

	return &s, nil
}

func (r *TransactionRepo) GetProfitReport(dateFrom, dateTo string) (*domain.ProfitReport, error) {
	var p domain.ProfitReport
	p.PeriodFrom = dateFrom
	p.PeriodTo = dateTo
	dateFilter := buildDateFilter(dateFrom, dateTo)

	// Total sales revenue
	q := fmt.Sprintf(`SELECT COALESCE(SUM(total_amount),0) FROM transactions WHERE type='sale' %s`,
		strings.Replace(dateFilter, "WHERE", "AND", 1))
	if dateFilter == "" {
		q = `SELECT COALESCE(SUM(total_amount),0) FROM transactions WHERE type='sale'`
	}
	r.db.QueryRow(context.Background(), q).Scan(&p.TotalSalesRevenue)

	// Cost of sold goods (using product cost_price at time of sale)
	q2 := fmt.Sprintf(`SELECT COALESCE(SUM(t.quantity * pr.cost_price),0)
		FROM transactions t JOIN products pr ON pr.id = t.product_id
		WHERE t.type='sale' %s`,
		strings.Replace(dateFilter, "WHERE", "AND", 1))
	if dateFilter == "" {
		q2 = `SELECT COALESCE(SUM(t.quantity * pr.cost_price),0) FROM transactions t JOIN products pr ON pr.id = t.product_id WHERE t.type='sale'`
	}
	r.db.QueryRow(context.Background(), q2).Scan(&p.TotalCostOfSold)

	// Write-off losses
	q3 := fmt.Sprintf(`SELECT COALESCE(SUM(t.quantity * pr.cost_price),0)
		FROM transactions t JOIN products pr ON pr.id = t.product_id
		WHERE t.type='write_off' %s`,
		strings.Replace(dateFilter, "WHERE", "AND", 1))
	if dateFilter == "" {
		q3 = `SELECT COALESCE(SUM(t.quantity * pr.cost_price),0) FROM transactions t JOIN products pr ON pr.id = t.product_id WHERE t.type='write_off'`
	}
	r.db.QueryRow(context.Background(), q3).Scan(&p.TotalWriteOffLoss)

	p.NetProfit = p.TotalSalesRevenue - (p.TotalCostOfSold + p.TotalWriteOffLoss)
	return &p, nil
}

func (r *TransactionRepo) GetTopProducts(limit int, dateFrom, dateTo string) ([]domain.TopProduct, error) {
	dateFilter := buildDateFilter(dateFrom, dateTo)
	extra := ""
	if dateFilter != "" {
		extra = strings.Replace(dateFilter, "WHERE", "AND", 1)
	}
	q := fmt.Sprintf(`SELECT t.product_id, p.name, SUM(t.quantity), SUM(t.total_amount)
		FROM transactions t JOIN products p ON p.id = t.product_id
		WHERE t.type='sale' %s
		GROUP BY t.product_id, p.name
		ORDER BY SUM(t.total_amount) DESC LIMIT %d`, extra, limit)

	rows, err := r.db.Query(context.Background(), q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tops []domain.TopProduct
	for rows.Next() {
		var tp domain.TopProduct
		if err := rows.Scan(&tp.ProductID, &tp.ProductName, &tp.TotalSold, &tp.TotalRevenue); err != nil {
			return nil, err
		}
		tops = append(tops, tp)
	}
	return tops, nil
}

func (r *TransactionRepo) GetInventory() ([]domain.InventoryItem, error) {
	q := `SELECT p.id, p.name, p.unit, COALESCE(c.name,''), p.cost_price, p.sale_price,
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

	var items []domain.InventoryItem
	for rows.Next() {
		var item domain.InventoryItem
		if err := rows.Scan(&item.ProductID, &item.ProductName, &item.Unit, &item.CategoryName,
			&item.CostPrice, &item.SalePrice, &item.CurrentStock); err != nil {
			return nil, err
		}
		item.StockValue = item.CurrentStock * item.CostPrice
		items = append(items, item)
	}
	return items, nil
}

func buildDateFilter(dateFrom, dateTo string) string {
	var parts []string
	if dateFrom != "" {
		parts = append(parts, fmt.Sprintf("created_at >= '%s'", dateFrom))
	}
	if dateTo != "" {
		parts = append(parts, fmt.Sprintf("created_at <= '%s'::date + interval '1 day'", dateTo))
	}
	if len(parts) > 0 {
		return "WHERE " + strings.Join(parts, " AND ")
	}
	return ""
}

func (r *TransactionRepo) GetAlerts() ([]domain.DashboardAlert, error) {
	var alerts []domain.DashboardAlert

	// 1. Low stock alerts
	lowStockQuery := `SELECT id, name, unit, min_stock, stock FROM (
		SELECT p.id, p.name, p.unit, p.min_stock,
		COALESCE(SUM(CASE WHEN t.type='purchase' THEN t.quantity ELSE 0 END),0)
		- COALESCE(SUM(CASE WHEN t.type='return' THEN t.quantity ELSE 0 END),0)
		- COALESCE(SUM(CASE WHEN t.type='sale' THEN t.quantity ELSE 0 END),0)
		- COALESCE(SUM(CASE WHEN t.type='write_off' THEN t.quantity ELSE 0 END),0) AS stock
		FROM products p LEFT JOIN transactions t ON t.product_id = p.id
		GROUP BY p.id
	) sub WHERE min_stock > 0 AND stock <= min_stock ORDER BY stock ASC`

	rows, err := r.db.Query(context.Background(), lowStockQuery)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var id int
			var name, unit string
			var min, stock float64
			if err := rows.Scan(&id, &name, &unit, &min, &stock); err == nil {
				alerts = append(alerts, domain.DashboardAlert{
					Type:        "low_stock",
					ProductID:   id,
					ProductName: name,
					Message:     fmt.Sprintf("Qoldiq kam qoldi (Min: %g)", min),
					Value:       fmt.Sprintf("%g %s", stock, unit),
				})
			}
		}
	}

	// 2. Expiring soon alerts (within 7 days)
	expQuery := `SELECT p.id, p.name, t.expiry_date, t.quantity, p.unit 
		FROM transactions t JOIN products p ON p.id = t.product_id
		WHERE t.type = 'purchase' AND t.expiry_date IS NOT NULL 
		AND t.expiry_date <= CURRENT_DATE + INTERVAL '7 days'
		AND t.expiry_date >= CURRENT_DATE`
	
	rows2, err2 := r.db.Query(context.Background(), expQuery)
	if err2 == nil {
		defer rows2.Close()
		for rows2.Next() {
			var id int
			var name, unit, expDate string
			var qty float64
			if err := rows2.Scan(&id, &name, &expDate, &qty, &unit); err == nil {
				alerts = append(alerts, domain.DashboardAlert{
					Type:        "expiring_soon",
					ProductID:   id,
					ProductName: name,
					Message:     "Muddati tugamoqda",
					Value:       fmt.Sprintf("%s gacha", strings.Split(expDate, "T")[0]),
				})
			}
		}
	}

	return alerts, nil
}
