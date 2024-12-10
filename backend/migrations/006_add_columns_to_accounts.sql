ALTER TABLE transactions
ADD CONSTRAINT fk_transactions_categories
FOREIGN KEY (category_id) REFERENCES categories (id)
ON DELETE SET NULL;
