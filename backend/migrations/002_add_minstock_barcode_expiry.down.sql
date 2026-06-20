ALTER TABLE transactions DROP COLUMN expiry_date;

DROP INDEX idx_products_barcode;
ALTER TABLE products DROP COLUMN barcode;
ALTER TABLE products DROP COLUMN min_stock;
