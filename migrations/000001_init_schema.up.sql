-- migrations/000001_init_schema.up.sql

-- Create accounts table
CREATE TABLE IF NOT EXISTS accounts (
    account_id SERIAL PRIMARY KEY,
    document_number VARCHAR(20) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create operation_types table
CREATE TABLE IF NOT EXISTS operation_types (
    operation_type_id INTEGER PRIMARY KEY,
    description VARCHAR(100) NOT NULL
);

-- Insert operation types
INSERT INTO operation_types (operation_type_id, description) VALUES
    (1, 'Normal Purchase'),
    (2, 'Purchase with installments'),
    (3, 'Withdrawal'),
    (4, 'Credit Voucher')
ON CONFLICT (operation_type_id) DO NOTHING;

-- Create transactions table
CREATE TABLE IF NOT EXISTS transactions (
    transaction_id SERIAL PRIMARY KEY,
    account_id INTEGER NOT NULL REFERENCES accounts(account_id) ON DELETE CASCADE,
    operation_type_id INTEGER NOT NULL REFERENCES operation_types(operation_type_id),
    amount DECIMAL(15, 2) NOT NULL,
    event_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT amount_check CHECK (
        (operation_type_id IN (1, 2, 3) AND amount < 0) OR
        (operation_type_id = 4 AND amount > 0)
    )
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_transactions_account_id ON transactions(account_id);
CREATE INDEX IF NOT EXISTS idx_transactions_event_date ON transactions(event_date);
CREATE INDEX IF NOT EXISTS idx_accounts_document_number ON accounts(document_number);

-- Add comments
COMMENT ON TABLE accounts IS 'Stores customer account information';
COMMENT ON TABLE transactions IS 'Stores all financial transactions';
COMMENT ON TABLE operation_types IS 'Lookup table for transaction types';