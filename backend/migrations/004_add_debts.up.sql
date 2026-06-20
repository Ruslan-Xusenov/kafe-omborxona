CREATE TABLE debts (
    id SERIAL PRIMARY KEY,
    supplier_id INT REFERENCES suppliers(id) ON DELETE CASCADE,
    transaction_id INT REFERENCES transactions(id) ON DELETE CASCADE,
    total_debt DECIMAL(14,2) NOT NULL,
    paid_amount DECIMAL(14,2) NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'unpaid' CHECK (status IN ('unpaid', 'partial', 'paid')),
    due_date DATE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_debts_supplier ON debts(supplier_id);
CREATE INDEX idx_debts_status ON debts(status);
