-- =====================================================
-- ERP System - Core Tables
-- Version: 1.0
-- =====================================================

-- =====================================================
-- ENUM TYPES
-- =====================================================

CREATE TYPE order_type AS ENUM (
    'STANDARD', 'ADVANCE', 'PRE_PAID', 'ON_HOLD', 
    'QUOTE', 'CREDIT_MEMO', 'PICK_UP'
);

CREATE TYPE order_status AS ENUM (
    'DRAFT', 'CONFIRMED', 'PICKING', 'SHIPPED', 
    'DELIVERED', 'INVOICED', 'CANCELLED'
);

CREATE TYPE po_status AS ENUM (
    'DRAFT', 'SUBMITTED', 'CONFIRMED', 'PARTIAL', 
    'RECEIVED', 'CLOSED', 'CANCELLED'
);

CREATE TYPE ar_invoice_status AS ENUM (
    'DRAFT', 'POSTED', 'PARTIAL', 'PAID', 'OVERDUE', 'VOID'
);

CREATE TYPE ap_invoice_status AS ENUM (
    'PENDING', 'APPROVED', 'PARTIAL', 'PAID', 'VOID'
);

CREATE TYPE inventory_transaction_type AS ENUM (
    'RECEIVE', 'SHIP', 'ADJUST_IN', 'ADJUST_OUT',
    'TRANSFER_IN', 'TRANSFER_OUT', 'RETURN', 'DISPOSE', 'CYCLE_COUNT'
);

CREATE TYPE wms_task_status AS ENUM (
    'PENDING', 'ASSIGNED', 'IN_PROGRESS', 'COMPLETE', 'CANCELLED'
);

CREATE TYPE pick_list_status AS ENUM (
    'PENDING', 'IN_PROGRESS', 'COMPLETE', 'CANCELLED'
);

CREATE TYPE bank_transaction_type AS ENUM (
    'DEPOSIT', 'WITHDRAWAL', 'CHECK', 'TRANSFER',
    'FEE', 'INTEREST', 'ADJUSTMENT'
);

CREATE TYPE gl_account_type AS ENUM (
    'ASSET', 'LIABILITY', 'EQUITY', 'REVENUE', 'EXPENSE'
);

CREATE TYPE emp_status AS ENUM ('RESIGN', 'CONTINUED');
CREATE TYPE gender AS ENUM ('MALE', 'FEMALE');
CREATE TYPE contract_type AS ENUM ('YEARLY', 'FULL');
CREATE TYPE payroll_type AS ENUM ('LECTURER', 'EMPLOYEE', 'PROFESSOR');

-- =====================================================
-- USERS & EMPLOYEES
-- =====================================================

CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    role_name TEXT NOT NULL,
    role_desc TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS contracts (
    id SERIAL PRIMARY KEY,
    start_date DATE,
    end_date DATE,
    contract_type contract_type,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS employees (
    id SERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    english_name TEXT,
    arabic_name TEXT,
    nationality TEXT,
    phone TEXT,
    date_of_birth DATE,
    status emp_status DEFAULT 'CONTINUED',
    contract_id INT REFERENCES contracts(id),
    role_id INT REFERENCES roles(id),
    account_status TEXT DEFAULT 'active',
    is_sales_rep BOOLEAN DEFAULT false,
    is_buyer BOOLEAN DEFAULT false,
    is_driver BOOLEAN DEFAULT false,
    commission_rate DECIMAL(5,2) DEFAULT 0,
    security_level INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS employee_details (
    id SERIAL PRIMARY KEY,
    emp_id INT UNIQUE REFERENCES employees(id),
    gender gender,
    job_title TEXT,
    major_study TEXT,
    notes TEXT,
    passport_number TEXT,
    is_retired BOOLEAN DEFAULT false,
    is_married BOOLEAN DEFAULT false,
    number_of_children SMALLINT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS pages (
    id SERIAL PRIMARY KEY,
    page_name TEXT,
    route_name TEXT NOT NULL,
    icon TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS emp_page (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES employees(id),
    page_id INT NOT NULL REFERENCES pages(id),
    can_create BOOLEAN DEFAULT FALSE,
    can_update BOOLEAN DEFAULT FALSE,
    can_delete BOOLEAN DEFAULT FALSE,
    can_view BOOLEAN DEFAULT FALSE,
    UNIQUE (user_id, page_id)
);

-- =====================================================
-- WAREHOUSES
-- =====================================================

CREATE TABLE IF NOT EXISTS warehouses (
    id SERIAL PRIMARY KEY,
    warehouse_code VARCHAR(10) UNIQUE NOT NULL,
    name TEXT NOT NULL,
    address_line1 TEXT,
    address_line2 TEXT,
    city VARCHAR(100),
    state VARCHAR(100),
    postal_code VARCHAR(20),
    country VARCHAR(100),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS warehouse_zones (
    id SERIAL PRIMARY KEY,
    warehouse_id INT REFERENCES warehouses(id),
    zone_code VARCHAR(10) NOT NULL,
    name TEXT,
    zone_type VARCHAR(20),
    temperature_controlled BOOLEAN DEFAULT false,
    min_temperature DECIMAL(5,2),
    max_temperature DECIMAL(5,2),
    UNIQUE(warehouse_id, zone_code)
);

CREATE TABLE IF NOT EXISTS warehouse_locations (
    id SERIAL PRIMARY KEY,
    warehouse_id INT REFERENCES warehouses(id),
    zone_id INT REFERENCES warehouse_zones(id),
    location_code VARCHAR(50) UNIQUE NOT NULL,
    aisle VARCHAR(10),
    rack VARCHAR(10),
    shelf VARCHAR(10),
    bin VARCHAR(10),
    location_type VARCHAR(20),
    max_weight DECIMAL(10,2),
    max_volume DECIMAL(10,2),
    is_active BOOLEAN DEFAULT true,
    pick_sequence INT
);

-- =====================================================
-- PRODUCTS
-- =====================================================

CREATE TABLE IF NOT EXISTS product_categories (
    id SERIAL PRIMARY KEY,
    code VARCHAR(20) UNIQUE,
    name TEXT NOT NULL,
    parent_id INT REFERENCES product_categories(id),
    gl_sales_account_id INT,
    gl_cogs_account_id INT,
    gl_inventory_account_id INT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    sku VARCHAR(50) UNIQUE NOT NULL,
    barcode VARCHAR(50),
    upc VARCHAR(50),
    name TEXT NOT NULL,
    description TEXT,
    category_id INT REFERENCES product_categories(id),
    base_unit VARCHAR(10) DEFAULT 'EA',
    is_catch_weight BOOLEAN DEFAULT false,
    catch_weight_unit VARCHAR(10),
    country_of_origin VARCHAR(3),
    shelf_life_days INT,
    min_shelf_life_days INT,
    is_lot_tracked BOOLEAN DEFAULT false,
    is_serialized BOOLEAN DEFAULT false,
    haccp_category VARCHAR(50),
    qc_required BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS product_units (
    id SERIAL PRIMARY KEY,
    product_id INT REFERENCES products(id) ON DELETE CASCADE,
    unit_name VARCHAR(10) NOT NULL,
    description TEXT,
    conversion_factor DECIMAL(10,6) NOT NULL,
    barcode VARCHAR(50),
    weight DECIMAL(10,4),
    is_purchase_unit BOOLEAN DEFAULT false,
    is_sales_unit BOOLEAN DEFAULT true,
    UNIQUE(product_id, unit_name)
);

-- =====================================================
-- CUSTOMERS
-- =====================================================

CREATE TABLE IF NOT EXISTS customers (
    id SERIAL PRIMARY KEY,
    customer_code VARCHAR(20) UNIQUE NOT NULL,
    name TEXT NOT NULL,
    billing_address_id INT,
    credit_limit DECIMAL(12,2) DEFAULT 0,
    current_balance DECIMAL(12,2) DEFAULT 0,
    payment_terms_days INT DEFAULT 30,
    currency VARCHAR(3) DEFAULT 'USD',
    sales_rep_id INT REFERENCES employees(id),
    default_route_id INT,
    default_warehouse_id INT REFERENCES warehouses(id),
    tax_exempt BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    created_by INT REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS customer_ship_to (
    id SERIAL PRIMARY KEY,
    customer_id INT REFERENCES customers(id) ON DELETE CASCADE,
    ship_to_code VARCHAR(20),
    name TEXT,
    address_line1 TEXT,
    address_line2 TEXT,
    city VARCHAR(100),
    state VARCHAR(100),
    postal_code VARCHAR(20),
    country VARCHAR(100),
    phone VARCHAR(50),
    is_default BOOLEAN DEFAULT false,
    warehouse_id INT REFERENCES warehouses(id),
    route_id INT,
    UNIQUE(customer_id, ship_to_code)
);

CREATE TABLE IF NOT EXISTS customer_order_guides (
    id SERIAL PRIMARY KEY,
    customer_id INT REFERENCES customers(id) ON DELETE CASCADE,
    product_id INT REFERENCES products(id) ON DELETE CASCADE,
    default_quantity DECIMAL(10,3),
    last_ordered_date DATE,
    last_ordered_quantity DECIMAL(10,3),
    avg_weekly_quantity DECIMAL(10,3),
    times_ordered INT DEFAULT 0,
    is_push_item BOOLEAN DEFAULT false,
    custom_price DECIMAL(10,4),
    UNIQUE(customer_id, product_id)
);

-- =====================================================
-- VENDORS
-- =====================================================

CREATE TABLE IF NOT EXISTS vendors (
    id SERIAL PRIMARY KEY,
    vendor_code VARCHAR(20) UNIQUE NOT NULL,
    name TEXT NOT NULL,
    address_line1 TEXT,
    address_line2 TEXT,
    city VARCHAR(100),
    state VARCHAR(100),
    postal_code VARCHAR(20),
    country VARCHAR(100),
    phone VARCHAR(50),
    email VARCHAR(255),
    payment_terms_days INT DEFAULT 30,
    currency VARCHAR(3) DEFAULT 'USD',
    lead_time_days INT DEFAULT 7,
    minimum_order DECIMAL(12,2),
    buyer_id INT REFERENCES employees(id),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS vendor_products (
    id SERIAL PRIMARY KEY,
    vendor_id INT REFERENCES vendors(id) ON DELETE CASCADE,
    product_id INT REFERENCES products(id) ON DELETE CASCADE,
    vendor_sku VARCHAR(50),
    vendor_description TEXT,
    unit_of_measure VARCHAR(10),
    unit_cost DECIMAL(10,4),
    minimum_order_qty DECIMAL(10,3),
    lead_time_days INT,
    is_preferred BOOLEAN DEFAULT false,
    UNIQUE(vendor_id, product_id)
);

CREATE TABLE IF NOT EXISTS vendor_discounts (
    id SERIAL PRIMARY KEY,
    vendor_id INT REFERENCES vendors(id) ON DELETE CASCADE,
    discount_days INT,
    discount_percent DECIMAL(5,2),
    UNIQUE(vendor_id, discount_days)
);

-- =====================================================
-- INVENTORY
-- =====================================================

CREATE TABLE IF NOT EXISTS inventory (
    id SERIAL PRIMARY KEY,
    product_id INT REFERENCES products(id),
    warehouse_id INT REFERENCES warehouses(id),
    location_code VARCHAR(50),
    lot_number VARCHAR(50),
    production_date DATE,
    expiry_date DATE,
    quantity_on_hand DECIMAL(12,3) DEFAULT 0,
    quantity_allocated DECIMAL(12,3) DEFAULT 0,
    quantity_on_order DECIMAL(12,3) DEFAULT 0,
    quantity_available DECIMAL(12,3) GENERATED ALWAYS AS (quantity_on_hand - quantity_allocated) STORED,
    last_cost DECIMAL(10,4),
    average_cost DECIMAL(10,4),
    last_counted_date DATE,
    last_movement_date DATE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(product_id, warehouse_id, location_code, lot_number)
);

CREATE TABLE IF NOT EXISTS inventory_transactions (
    id SERIAL PRIMARY KEY,
    product_id INT REFERENCES products(id),
    warehouse_id INT REFERENCES warehouses(id),
    location_code VARCHAR(50),
    transaction_type inventory_transaction_type,
    quantity DECIMAL(12,3),
    lot_number VARCHAR(50),
    unit_cost DECIMAL(10,4),
    reference_type VARCHAR(20),
    reference_id INT,
    reference_number VARCHAR(50),
    notes TEXT,
    created_by INT REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS product_reservations (
    id SERIAL PRIMARY KEY,
    product_id INT REFERENCES products(id),
    warehouse_id INT REFERENCES warehouses(id),
    reserved_for_type VARCHAR(20),
    reserved_for_id INT,
    quantity DECIMAL(10,3),
    expiry_date DATE,
    created_by INT REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- INDEXES
-- =====================================================

CREATE INDEX idx_employees_email ON employees(email);
CREATE INDEX idx_employees_role ON employees(role_id);
CREATE INDEX idx_products_sku ON products(sku);
CREATE INDEX idx_products_barcode ON products(barcode);
CREATE INDEX idx_products_category ON products(category_id);
CREATE INDEX idx_customers_code ON customers(customer_code);
CREATE INDEX idx_customers_sales_rep ON customers(sales_rep_id);
CREATE INDEX idx_vendors_code ON vendors(vendor_code);
CREATE INDEX idx_vendors_buyer ON vendors(buyer_id);
CREATE INDEX idx_inventory_product ON inventory(product_id);
CREATE INDEX idx_inventory_warehouse ON inventory(warehouse_id);
CREATE INDEX idx_inventory_lot ON inventory(lot_number);
CREATE INDEX idx_inventory_expiry ON inventory(expiry_date);
CREATE INDEX idx_inventory_transactions_product ON inventory_transactions(product_id);
CREATE INDEX idx_inventory_transactions_date ON inventory_transactions(created_at);

