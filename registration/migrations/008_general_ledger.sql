-- ============================================
-- General Ledger (GL) Tables
-- Complete double-entry accounting system
-- ============================================

-- GL Accounts (Chart of Accounts)
CREATE TABLE IF NOT EXISTS gl_accounts (
    id SERIAL PRIMARY KEY,
    account_code VARCHAR(20) NOT NULL UNIQUE,
    account_name VARCHAR(100) NOT NULL,
    account_type VARCHAR(20) NOT NULL, -- ASSET, LIABILITY, EQUITY, REVENUE, EXPENSE
    account_sub_type VARCHAR(30),      -- CASH, BANK, RECEIVABLES, INVENTORY, etc.
    parent_id INTEGER REFERENCES gl_accounts(id),
    description TEXT,
    currency VARCHAR(3) DEFAULT 'USD',
    is_active BOOLEAN DEFAULT TRUE,
    is_postable BOOLEAN DEFAULT TRUE,  -- Can post transactions to this account
    is_bank_account BOOLEAN DEFAULT FALSE,
    bank_account_id INTEGER,           -- Link to bank_accounts table if bank account
    normal_balance VARCHAR(6) NOT NULL, -- DEBIT or CREDIT
    opening_balance DECIMAL(15,2) DEFAULT 0,
    current_balance DECIMAL(15,2) DEFAULT 0,
    budget_amount DECIMAL(15,2) DEFAULT 0,
    department_id INTEGER REFERENCES departments(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- GL Fiscal Years
CREATE TABLE IF NOT EXISTS gl_fiscal_years (
    id SERIAL PRIMARY KEY,
    year_code VARCHAR(20) NOT NULL UNIQUE,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_current BOOLEAN DEFAULT FALSE,
    is_closed BOOLEAN DEFAULT FALSE,
    closed_by INTEGER REFERENCES employees(id),
    closed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- GL Periods (Monthly periods within fiscal year)
CREATE TABLE IF NOT EXISTS gl_periods (
    id SERIAL PRIMARY KEY,
    fiscal_year_id INTEGER NOT NULL REFERENCES gl_fiscal_years(id),
    period_number INTEGER NOT NULL,     -- 1-12 for normal periods, 13 for adjustments
    period_name VARCHAR(50) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    status VARCHAR(10) DEFAULT 'OPEN', -- OPEN, CLOSED, LOCKED
    is_adjustment BOOLEAN DEFAULT FALSE, -- True for period 13
    closed_by INTEGER REFERENCES employees(id),
    closed_at TIMESTAMP,
    UNIQUE(fiscal_year_id, period_number)
);

-- GL Journal Entries (Header)
CREATE TABLE IF NOT EXISTS gl_journal_entries (
    id SERIAL PRIMARY KEY,
    journal_number VARCHAR(30) NOT NULL UNIQUE,
    entry_date DATE NOT NULL,
    posting_date DATE,
    period_id INTEGER NOT NULL REFERENCES gl_periods(id),
    entry_type VARCHAR(20) NOT NULL,    -- MANUAL, AR, AP, INVENTORY, PAYROLL, BANK, RECURRING, ADJUSTMENT, CLOSING
    status VARCHAR(15) DEFAULT 'DRAFT', -- DRAFT, PENDING, POSTED, REVERSED, VOIDED
    description TEXT NOT NULL,
    reference VARCHAR(100),
    source_document VARCHAR(50),        -- e.g., "AR-INV-001"
    source_module VARCHAR(20),          -- e.g., "AR", "AP"
    source_id INTEGER,                  -- ID of source document
    total_debit DECIMAL(15,2) NOT NULL DEFAULT 0,
    total_credit DECIMAL(15,2) NOT NULL DEFAULT 0,
    currency VARCHAR(3) DEFAULT 'USD',
    exchange_rate DECIMAL(10,6) DEFAULT 1.0,
    is_recurring BOOLEAN DEFAULT FALSE,
    recurring_id INTEGER,               -- Link to recurring template
    reversed_entry_id INTEGER REFERENCES gl_journal_entries(id),
    auto_reverse BOOLEAN DEFAULT FALSE,
    auto_reverse_date DATE,
    created_by INTEGER NOT NULL REFERENCES employees(id),
    posted_by INTEGER REFERENCES employees(id),
    posted_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- GL Journal Lines (Detail)
CREATE TABLE IF NOT EXISTS gl_journal_lines (
    id SERIAL PRIMARY KEY,
    journal_id INTEGER NOT NULL REFERENCES gl_journal_entries(id) ON DELETE CASCADE,
    line_number INTEGER NOT NULL,
    account_id INTEGER NOT NULL REFERENCES gl_accounts(id),
    description TEXT,
    debit_amount DECIMAL(15,2) DEFAULT 0,
    credit_amount DECIMAL(15,2) DEFAULT 0,
    department_id INTEGER REFERENCES departments(id),
    project_id INTEGER,
    reference VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT chk_debit_or_credit CHECK (
        (debit_amount = 0 AND credit_amount > 0) OR 
        (debit_amount > 0 AND credit_amount = 0) OR
        (debit_amount = 0 AND credit_amount = 0)
    )
);

-- GL Recurring Entry Templates
CREATE TABLE IF NOT EXISTS gl_recurring_entries (
    id SERIAL PRIMARY KEY,
    template_name VARCHAR(100) NOT NULL,
    description TEXT,
    frequency VARCHAR(20) NOT NULL,     -- MONTHLY, QUARTERLY, YEARLY
    next_run_date DATE NOT NULL,
    end_date DATE,
    last_run_date DATE,
    total_runs INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    auto_reverse BOOLEAN DEFAULT FALSE,
    days_to_reverse INTEGER DEFAULT 0,
    created_by INTEGER NOT NULL REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW()
);

-- GL Recurring Entry Lines
CREATE TABLE IF NOT EXISTS gl_recurring_lines (
    id SERIAL PRIMARY KEY,
    recurring_id INTEGER NOT NULL REFERENCES gl_recurring_entries(id) ON DELETE CASCADE,
    line_number INTEGER NOT NULL,
    account_id INTEGER NOT NULL REFERENCES gl_accounts(id),
    description TEXT,
    debit_amount DECIMAL(15,2) DEFAULT 0,
    credit_amount DECIMAL(15,2) DEFAULT 0
);

-- GL Budgets
CREATE TABLE IF NOT EXISTS gl_budgets (
    id SERIAL PRIMARY KEY,
    fiscal_year_id INTEGER NOT NULL REFERENCES gl_fiscal_years(id),
    budget_name VARCHAR(100) NOT NULL,
    description TEXT,
    is_approved BOOLEAN DEFAULT FALSE,
    approved_by INTEGER REFERENCES employees(id),
    approved_at TIMESTAMP,
    created_by INTEGER NOT NULL REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW()
);

-- GL Budget Lines
CREATE TABLE IF NOT EXISTS gl_budget_lines (
    id SERIAL PRIMARY KEY,
    budget_id INTEGER NOT NULL REFERENCES gl_budgets(id) ON DELETE CASCADE,
    account_id INTEGER NOT NULL REFERENCES gl_accounts(id),
    period_id INTEGER NOT NULL REFERENCES gl_periods(id),
    budget_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    notes TEXT,
    UNIQUE(budget_id, account_id, period_id)
);

-- Indexes
CREATE INDEX idx_gl_accounts_type ON gl_accounts(account_type);
CREATE INDEX idx_gl_accounts_parent ON gl_accounts(parent_id);
CREATE INDEX idx_gl_accounts_code ON gl_accounts(account_code);
CREATE INDEX idx_gl_periods_fiscal_year ON gl_periods(fiscal_year_id);
CREATE INDEX idx_gl_periods_dates ON gl_periods(start_date, end_date);
CREATE INDEX idx_gl_journal_entries_period ON gl_journal_entries(period_id);
CREATE INDEX idx_gl_journal_entries_date ON gl_journal_entries(entry_date);
CREATE INDEX idx_gl_journal_entries_status ON gl_journal_entries(status);
CREATE INDEX idx_gl_journal_entries_type ON gl_journal_entries(entry_type);
CREATE INDEX idx_gl_journal_lines_journal ON gl_journal_lines(journal_id);
CREATE INDEX idx_gl_journal_lines_account ON gl_journal_lines(account_id);

-- Insert default Chart of Accounts
INSERT INTO gl_accounts (account_code, account_name, account_type, account_sub_type, normal_balance, is_postable) VALUES
-- Assets
('1000', 'Assets', 'ASSET', NULL, 'DEBIT', FALSE),
('1100', 'Cash and Bank', 'ASSET', 'CASH', 'DEBIT', FALSE),
('1110', 'Cash on Hand', 'ASSET', 'CASH', 'DEBIT', TRUE),
('1120', 'Checking Account', 'ASSET', 'BANK', 'DEBIT', TRUE),
('1130', 'Savings Account', 'ASSET', 'BANK', 'DEBIT', TRUE),
('1200', 'Accounts Receivable', 'ASSET', 'RECEIVABLES', 'DEBIT', TRUE),
('1300', 'Inventory', 'ASSET', 'INVENTORY', 'DEBIT', TRUE),
('1400', 'Prepaid Expenses', 'ASSET', 'OTHER_ASSETS', 'DEBIT', TRUE),
('1500', 'Fixed Assets', 'ASSET', 'FIXED_ASSETS', 'DEBIT', FALSE),
('1510', 'Equipment', 'ASSET', 'FIXED_ASSETS', 'DEBIT', TRUE),
('1520', 'Vehicles', 'ASSET', 'FIXED_ASSETS', 'DEBIT', TRUE),
('1590', 'Accumulated Depreciation', 'ASSET', 'FIXED_ASSETS', 'CREDIT', TRUE),

-- Liabilities
('2000', 'Liabilities', 'LIABILITY', NULL, 'CREDIT', FALSE),
('2100', 'Accounts Payable', 'LIABILITY', 'PAYABLES', 'CREDIT', TRUE),
('2200', 'Accrued Liabilities', 'LIABILITY', 'ACCRUED_LIABILITIES', 'CREDIT', TRUE),
('2300', 'Sales Tax Payable', 'LIABILITY', 'ACCRUED_LIABILITIES', 'CREDIT', TRUE),
('2400', 'Payroll Liabilities', 'LIABILITY', 'ACCRUED_LIABILITIES', 'CREDIT', TRUE),
('2500', 'Long Term Debt', 'LIABILITY', 'LONG_TERM_DEBT', 'CREDIT', TRUE),

-- Equity
('3000', 'Equity', 'EQUITY', NULL, 'CREDIT', FALSE),
('3100', 'Owner Capital', 'EQUITY', 'CAPITAL', 'CREDIT', TRUE),
('3200', 'Retained Earnings', 'EQUITY', 'RETAINED_EARNINGS', 'CREDIT', TRUE),

-- Revenue
('4000', 'Revenue', 'REVENUE', NULL, 'CREDIT', FALSE),
('4100', 'Sales Revenue', 'REVENUE', 'SALES', 'CREDIT', TRUE),
('4200', 'Service Revenue', 'REVENUE', 'SALES', 'CREDIT', TRUE),
('4300', 'Other Income', 'REVENUE', 'OTHER_INCOME', 'CREDIT', TRUE),

-- Expenses
('5000', 'Cost of Goods Sold', 'EXPENSE', 'COGS', 'DEBIT', TRUE),
('6000', 'Operating Expenses', 'EXPENSE', NULL, 'DEBIT', FALSE),
('6100', 'Salaries & Wages', 'EXPENSE', 'PAYROLL_EXPENSES', 'DEBIT', TRUE),
('6200', 'Rent Expense', 'EXPENSE', 'OPERATING_EXPENSES', 'DEBIT', TRUE),
('6300', 'Utilities Expense', 'EXPENSE', 'OPERATING_EXPENSES', 'DEBIT', TRUE),
('6400', 'Insurance Expense', 'EXPENSE', 'OPERATING_EXPENSES', 'DEBIT', TRUE),
('6500', 'Depreciation Expense', 'EXPENSE', 'OPERATING_EXPENSES', 'DEBIT', TRUE),
('6600', 'Supplies Expense', 'EXPENSE', 'OPERATING_EXPENSES', 'DEBIT', TRUE),
('6700', 'Bank Charges', 'EXPENSE', 'OTHER_EXPENSES', 'DEBIT', TRUE),
('6800', 'Miscellaneous Expense', 'EXPENSE', 'OTHER_EXPENSES', 'DEBIT', TRUE)
ON CONFLICT (account_code) DO NOTHING;

-- Update parent IDs
UPDATE gl_accounts SET parent_id = (SELECT id FROM gl_accounts WHERE account_code = '1000') WHERE account_code IN ('1100', '1200', '1300', '1400', '1500');
UPDATE gl_accounts SET parent_id = (SELECT id FROM gl_accounts WHERE account_code = '1100') WHERE account_code IN ('1110', '1120', '1130');
UPDATE gl_accounts SET parent_id = (SELECT id FROM gl_accounts WHERE account_code = '1500') WHERE account_code IN ('1510', '1520', '1590');
UPDATE gl_accounts SET parent_id = (SELECT id FROM gl_accounts WHERE account_code = '2000') WHERE account_code IN ('2100', '2200', '2300', '2400', '2500');
UPDATE gl_accounts SET parent_id = (SELECT id FROM gl_accounts WHERE account_code = '3000') WHERE account_code IN ('3100', '3200');
UPDATE gl_accounts SET parent_id = (SELECT id FROM gl_accounts WHERE account_code = '4000') WHERE account_code IN ('4100', '4200', '4300');
UPDATE gl_accounts SET parent_id = (SELECT id FROM gl_accounts WHERE account_code = '6000') WHERE account_code IN ('6100', '6200', '6300', '6400', '6500', '6600', '6700', '6800');

COMMENT ON TABLE gl_accounts IS 'Chart of Accounts - all GL accounts';
COMMENT ON TABLE gl_fiscal_years IS 'Fiscal years for financial reporting';
COMMENT ON TABLE gl_periods IS 'Monthly periods within each fiscal year';
COMMENT ON TABLE gl_journal_entries IS 'Journal entry headers';
COMMENT ON TABLE gl_journal_lines IS 'Journal entry line items';
COMMENT ON COLUMN gl_accounts.normal_balance IS 'DEBIT for assets/expenses, CREDIT for liabilities/equity/revenue';
COMMENT ON COLUMN gl_accounts.is_postable IS 'If false, this is a summary account';
