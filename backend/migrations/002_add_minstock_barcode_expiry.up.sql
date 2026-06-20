ALTER TABLE products ADD COLUMN min_stock DECIMAL(14,3) NOT NULL DEFAULT 0;
ALTER TABLE products ADD COLUMN barcode VARCHAR(100);

CREATE INDEX idx_products_barcode ON products(barcode);

ALTER TABLE transactions ADD COLUMN expiry_date DATE;
