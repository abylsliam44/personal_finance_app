-- 004_remove_date_column.sql
ALTER TABLE transactions
DROP COLUMN IF EXISTS date;
