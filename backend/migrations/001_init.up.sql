CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'warehouse_manager'
        CHECK (role IN ('admin', 'warehouse_manager')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE suppliers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    phone VARCHAR(20),
    address TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    unit VARCHAR(20) NOT NULL DEFAULT 'dona',
    category_id INT NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    cost_price DECIMAL(14,2) NOT NULL DEFAULT 0,
    sale_price DECIMAL(14,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_products_category ON products(category_id);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    product_id INT NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    supplier_id INT REFERENCES suppliers(id) ON DELETE SET NULL,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    type VARCHAR(20) NOT NULL
        CHECK (type IN ('purchase', 'return', 'sale', 'write_off')),
    quantity DECIMAL(14,3) NOT NULL CHECK (quantity > 0),
    unit_price DECIMAL(14,2) NOT NULL CHECK (unit_price >= 0),
    total_amount DECIMAL(14,2) NOT NULL,
    note TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_transactions_product ON transactions(product_id);
CREATE INDEX idx_transactions_type ON transactions(type);
CREATE INDEX idx_transactions_date ON transactions(created_at);
