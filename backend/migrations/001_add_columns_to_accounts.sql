-- 001_add_columns_to_accounts.sql
ALTER TABLE accounts
ADD COLUMN IF NOT EXISTS account_type VARCHAR(50),
ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT TRUE;
