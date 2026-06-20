ALTER TABLE transactions DROP CONSTRAINT transactions_product_id_fkey;
ALTER TABLE transactions ADD CONSTRAINT transactions_product_id_fkey FOREIGN KEY (product_id) REFERENCES products(id);

ALTER TABLE transactions DROP CONSTRAINT IF EXISTS transactions_user_id_fkey;
ALTER TABLE transactions ADD CONSTRAINT transactions_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id);
