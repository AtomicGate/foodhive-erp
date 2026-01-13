-- =====================================================
-- SCHEMA UPDATES - Add missing columns and tables
-- =====================================================
-- EXECUTION ORDER:
-- 1. 001_core_tables.sql
-- 2. 002_transactions_tables.sql
-- 3. 004_insert_roles.sql (roles must exist first)
-- 4. 005_insert_default_admin.sql (creates employee id=1)
-- 5. 003_schema_updates.sql (this file - assigns permissions)

-- Add missing columns to employees table
ALTER TABLE employees ADD COLUMN IF NOT EXISTS department_id INT;
ALTER TABLE employees ADD COLUMN IF NOT EXISTS warehouse_id INT REFERENCES warehouses(id);
ALTER TABLE employees ADD COLUMN IF NOT EXISTS created_by INT REFERENCES employees(id);

-- Add missing columns to employee_details table
ALTER TABLE employee_details ADD COLUMN IF NOT EXISTS national_id TEXT;

-- Rename emp_id to employee_id in employee_details if needed
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'employee_details' AND column_name = 'emp_id') THEN
        ALTER TABLE employee_details RENAME COLUMN emp_id TO employee_id;
    END IF;
END $$;

-- Add missing columns to contracts table
ALTER TABLE contracts ADD COLUMN IF NOT EXISTS notes TEXT;

-- =====================================================
-- DEPARTMENTS TABLE
-- =====================================================

CREATE TABLE IF NOT EXISTS departments (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    manager_id INT REFERENCES employees(id),
    parent_id INT REFERENCES departments(id),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Add foreign key for department_id in employees
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'employees_department_id_fkey'
    ) THEN
        ALTER TABLE employees ADD CONSTRAINT employees_department_id_fkey 
        FOREIGN KEY (department_id) REFERENCES departments(id);
    END IF;
END $$;

-- =====================================================
-- EMPLOYEE FINANCES TABLE
-- =====================================================

CREATE TABLE IF NOT EXISTS employee_finances (
    id SERIAL PRIMARY KEY,
    employee_id INT UNIQUE REFERENCES employees(id) ON DELETE CASCADE,
    base_salary DECIMAL(12,2) DEFAULT 0,
    years_of_service INT DEFAULT 0,
    academic_allowance DECIMAL(10,2) DEFAULT 0,
    degree_allowance DECIMAL(10,2) DEFAULT 0,
    position_allowance DECIMAL(10,2) DEFAULT 0,
    profession_allowance DECIMAL(10,2) DEFAULT 0,
    transport_allowance DECIMAL(10,2) DEFAULT 0,
    housing_allowance DECIMAL(10,2) DEFAULT 0,
    bank_account_number TEXT,
    bank_name TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- ADDRESSES TABLE
-- =====================================================

CREATE TABLE IF NOT EXISTS addresses (
    id SERIAL PRIMARY KEY,
    entity_type TEXT NOT NULL,
    entity_id INT NOT NULL,
    address_line1 TEXT,
    address_line2 TEXT,
    city TEXT,
    state TEXT,
    country TEXT,
    postal_code TEXT,
    house TEXT,
    avenue TEXT,
    neighborhood TEXT,
    emergency_phone_number TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_addresses_entity ON addresses(entity_type, entity_id);

-- =====================================================
-- INDEXES
-- =====================================================

CREATE INDEX IF NOT EXISTS idx_employees_department ON employees(department_id);
CREATE INDEX IF NOT EXISTS idx_employees_warehouse ON employees(warehouse_id);
CREATE INDEX IF NOT EXISTS idx_departments_parent ON departments(parent_id);

-- =====================================================
-- INSERT DEFAULT DATA
-- =====================================================

-- Insert default department if not exists
INSERT INTO departments (id, name, description) 
VALUES (1, 'Administration', 'Administrative department')
ON CONFLICT (id) DO NOTHING;

-- Reset the sequence for departments
SELECT setval('departments_id_seq', COALESCE((SELECT MAX(id) FROM departments), 1));

-- =====================================================
-- ADD MISSING PAGES FOR PERMISSIONS
-- =====================================================

-- Add unique constraint to pages route_name if not exists
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'pages_route_name_key'
    ) THEN
        ALTER TABLE pages ADD CONSTRAINT pages_route_name_key UNIQUE (route_name);
    END IF;
EXCEPTION
    WHEN duplicate_table THEN NULL;
END $$;

-- Insert pages for all modules
INSERT INTO pages (page_name, route_name, icon) VALUES 
    ('Dashboard', '/dashboard', 'home'),
    ('Sales Orders', '/sales-orders', 'shopping-cart'),
    ('Purchase Orders', '/purchase-orders', 'truck'),
    ('Inventory', '/inventory', 'package'),
    ('Products', '/products', 'box'),
    ('Customers', '/customers', 'users'),
    ('Vendors', '/vendors', 'building'),
    ('Employees', '/employees', 'user'),
    ('Departments', '/departments', 'briefcase'),
    ('Roles', '/roles', 'shield'),
    ('Warehouses', '/warehouses', 'warehouse'),
    ('AR', '/ar', 'dollar-sign'),
    ('AP', '/ap', 'credit-card'),
    ('GL', '/gl', 'book'),
    ('Pricing', '/pricing', 'tag'),
    ('Picking', '/picking', 'clipboard'),
    ('Catch Weight', '/catch-weight', 'scale')
ON CONFLICT (route_name) DO NOTHING;

-- Give admin user (id=1) full permissions on all pages (only if employee exists)
INSERT INTO emp_page (user_id, page_id, can_create, can_update, can_delete, can_view)
SELECT 1, p.id, true, true, true, true
FROM pages p
WHERE EXISTS (SELECT 1 FROM employees WHERE id = 1)
  AND NOT EXISTS (
    SELECT 1 FROM emp_page ep WHERE ep.user_id = 1 AND ep.page_id = p.id
);
