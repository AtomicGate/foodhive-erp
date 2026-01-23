-- =====================================================
-- FoodHive ERP System - Core Tables
-- Version: 2.0 ENHANCED (DocText.md Requirements)
-- =====================================================

-- =====================================================
-- SECTION 1: ENUM TYPES
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

-- NEW ENUM TYPES FOR DocText.md REQUIREMENTS
CREATE TYPE price_type AS ENUM ('PRICE_1', 'PRICE_2', 'PRICE_3', 'WHOLESALE', 'RETAIL');
CREATE TYPE customer_level AS ENUM ('GENERAL', 'SILVER', 'GOLD', 'PLATINUM', 'VIP');
CREATE TYPE currency_code AS ENUM ('LAK', 'THB', 'USD', 'KIP', 'EUR', 'CNY');
CREATE TYPE pr_status AS ENUM ('DRAFT', 'SUBMITTED', 'CHECKED', 'AUTHORIZED', 'CONVERTED', 'CANCELLED');
CREATE TYPE advance_request_status AS ENUM ('DRAFT', 'SUBMITTED', 'APPROVED', 'REJECTED', 'CLEARED');
CREATE TYPE advance_voucher_status AS ENUM ('DRAFT', 'SUBMITTED', 'APPROVED', 'REJECTED', 'COMPLETED');
CREATE TYPE voucher_document_type AS ENUM ('RECEIPT', 'INVOICE', 'BANK_TRANSFER', 'PO_PR', 'TRANSPORTATION');
CREATE TYPE priority_level AS ENUM ('LOW', 'NORMAL', 'HIGH', 'URGENT', 'CRITICAL');
CREATE TYPE contact_preference AS ENUM ('PHONE', 'EMAIL', 'FAX', 'MOBILE', 'WHATSAPP', 'LINE');
CREATE TYPE payment_method AS ENUM ('CASH', 'BANK_TRANSFER', 'CHECK', 'CREDIT_CARD', 'CREDIT_TERM', 'COD');
CREATE TYPE approval_action AS ENUM ('PENDING', 'APPROVED', 'REJECTED', 'RETURNED', 'ESCALATED');
CREATE TYPE customer_type AS ENUM ('RETAIL', 'WHOLESALE', 'DISTRIBUTOR', 'CHAIN');
CREATE TYPE vendor_type AS ENUM ('SUPPLIER', 'MANUFACTURER', 'DISTRIBUTOR', 'IMPORTER');

-- =====================================================
-- SECTION 2: USERS, ROLES & EMPLOYEES
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
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS departments (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    manager_id INT,
    parent_id INT REFERENCES departments(id),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
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
    department_id INT REFERENCES departments(id),
    warehouse_id INT,
    account_status TEXT DEFAULT 'active',
    is_sales_rep BOOLEAN DEFAULT false,
    is_buyer BOOLEAN DEFAULT false,
    is_driver BOOLEAN DEFAULT false,
    commission_rate DECIMAL(5,2) DEFAULT 0,
    security_level INT DEFAULT 1,
    created_by INT REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

ALTER TABLE departments ADD CONSTRAINT fk_departments_manager 
    FOREIGN KEY (manager_id) REFERENCES employees(id);

CREATE TABLE IF NOT EXISTS employee_details (
    id SERIAL PRIMARY KEY,
    employee_id INT UNIQUE REFERENCES employees(id) ON DELETE CASCADE,
    gender gender,
    job_title TEXT,
    major_study TEXT,
    notes TEXT,
    passport_number TEXT,
    national_id TEXT,
    is_retired BOOLEAN DEFAULT false,
    is_married BOOLEAN DEFAULT false,
    number_of_children SMALLINT DEFAULT 0
);

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

CREATE TABLE IF NOT EXISTS pages (
    id SERIAL PRIMARY KEY,
    page_name TEXT,
    route_name TEXT NOT NULL UNIQUE,
    icon TEXT,
    parent_id INT REFERENCES pages(id),
    display_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS emp_page (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    page_id INT NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
    can_create BOOLEAN DEFAULT FALSE,
    can_update BOOLEAN DEFAULT FALSE,
    can_delete BOOLEAN DEFAULT FALSE,
    can_view BOOLEAN DEFAULT FALSE,
    UNIQUE (user_id, page_id)
);

-- =====================================================
-- SECTION 3: WAREHOUSES & LOCATIONS
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
    phone VARCHAR(50),
    email VARCHAR(255),
    manager_id INT REFERENCES employees(id),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW()
);

ALTER TABLE employees ADD CONSTRAINT fk_employees_warehouse 
    FOREIGN KEY (warehouse_id) REFERENCES warehouses(id);

CREATE TABLE IF NOT EXISTS warehouse_zones (
    id SERIAL PRIMARY KEY,
    warehouse_id INT REFERENCES warehouses(id) ON DELETE CASCADE,
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
    warehouse_id INT REFERENCES warehouses(id) ON DELETE CASCADE,
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
-- SECTION 4: PRODUCTS & CATEGORIES
-- =====================================================

CREATE TABLE IF NOT EXISTS product_categories (
    id SERIAL PRIMARY KEY,
    code VARCHAR(20) UNIQUE,
    name TEXT NOT NULL,
    description TEXT,
    parent_id INT REFERENCES product_categories(id),
    gl_sales_account_id INT,
    gl_cogs_account_id INT,
    gl_inventory_account_id INT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS product_groups (
    id SERIAL PRIMARY KEY,
    code VARCHAR(20) UNIQUE NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    parent_id INT REFERENCES product_groups(id),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS brands (
    id SERIAL PRIMARY KEY,
    code VARCHAR(20) UNIQUE NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    logo_url TEXT,
    website VARCHAR(255),
    country_of_origin VARCHAR(100),
    manufacturer TEXT,
    is_active BOOLEAN DEFAULT true,
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
    product_group_id INT REFERENCES product_groups(id),
    brand_id INT REFERENCES brands(id),
    base_unit VARCHAR(10) DEFAULT 'EA',
    is_catch_weight BOOLEAN DEFAULT false,
    catch_weight_unit VARCHAR(10),
    is_lot_tracked BOOLEAN DEFAULT false,
    is_serialized BOOLEAN DEFAULT false,
    country_of_origin VARCHAR(3),
    shelf_life_days INT,
    min_shelf_life_days INT,
    haccp_category VARCHAR(50),
    qc_required BOOLEAN DEFAULT false,
    brand TEXT,
    product_group TEXT,
    price_1 DECIMAL(12,4) DEFAULT 0.00,
    price_2 DECIMAL(12,4) DEFAULT 0.00,
    price_3 DECIMAL(12,4) DEFAULT 0.00,
    wholesale_price DECIMAL(12,4) DEFAULT 0.00,
    retail_price DECIMAL(12,4) DEFAULT 0.00,
    vat_percent DECIMAL(5,2) DEFAULT 7.00,
    standard_cost DECIMAL(12,4) DEFAULT 0.00,
    last_purchase_cost DECIMAL(12,4) DEFAULT 0.00,
    average_cost DECIMAL(12,4) DEFAULT 0.00,
    reorder_point DECIMAL(12,3) DEFAULT 0.00,
    reorder_quantity DECIMAL(12,3) DEFAULT 0.00,
    safety_stock DECIMAL(12,3) DEFAULT 0.00,
    lead_time_days INT DEFAULT 7,
    avg_daily_sales DECIMAL(12,3) DEFAULT 0.00,
    last_sale_date DATE,
    total_quantity_sold DECIMAL(15,3) DEFAULT 0.00,
    is_seasonal BOOLEAN DEFAULT false,
    season_start_month INT,
    season_end_month INT,
    allergen_info TEXT,
    nutritional_info TEXT,
    storage_temp_min DECIMAL(5,2),
    storage_temp_max DECIMAL(5,2),
    storage_instructions TEXT,
    weight_kg DECIMAL(10,4),
    length_cm DECIMAL(10,2),
    width_cm DECIMAL(10,2),
    height_cm DECIMAL(10,2),
    volume_cm3 DECIMAL(12,2),
    image_url TEXT,
    thumbnail_url TEXT,
    primary_vendor_id INT,
    internal_notes TEXT,
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
-- SECTION 5: CUSTOMER GROUPS & CUSTOMERS
-- =====================================================

CREATE TABLE IF NOT EXISTS customer_groups (
    id SERIAL PRIMARY KEY,
    code VARCHAR(20) UNIQUE NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    default_price_type price_type DEFAULT 'PRICE_1',
    discount_percent DECIMAL(5,2) DEFAULT 0.00,
    default_payment_days INT DEFAULT 30,
    credit_limit_default DECIMAL(12,2) DEFAULT 0.00,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS customers (
    id SERIAL PRIMARY KEY,
    customer_code VARCHAR(20) UNIQUE NOT NULL,
    name TEXT NOT NULL,
    default_contact_name TEXT,
    full_address TEXT,
    delivery_name TEXT,
    tel VARCHAR(50),
    fax VARCHAR(50),
    mobile VARCHAR(50),
    email VARCHAR(255),
    web_email VARCHAR(255),
    tax_code VARCHAR(50),
    passport_no VARCHAR(50),
    country_code VARCHAR(10),
    gender gender,
    birthday DATE,
    price_type price_type DEFAULT 'PRICE_1',
    customer_level customer_level DEFAULT 'GENERAL',
    bill_discount DECIMAL(5,2) DEFAULT 0.00,
    credit_limit DECIMAL(12,2) DEFAULT 0,
    credit_days INT DEFAULT 0,
    current_balance DECIMAL(12,2) DEFAULT 0,
    payment_terms_days INT DEFAULT 30,
    currency VARCHAR(3) DEFAULT 'LAK',
    member_expire DATE,
    use_promotion BOOLEAN DEFAULT true,
    use_pmt_bill BOOLEAN DEFAULT false,
    collect_point BOOLEAN DEFAULT false,
    payee TEXT,
    rd_branch TEXT,
    barcode VARCHAR(100),
    notes TEXT,
    last_order_date DATE,
    total_orders INT DEFAULT 0,
    total_spent DECIMAL(15,2) DEFAULT 0.00,
    loyalty_points DECIMAL(12,2) DEFAULT 0.00,
    credit_score INT DEFAULT 100,
    preferred_delivery_time VARCHAR(50),
    delivery_instructions TEXT,
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),
    referred_by_id INT REFERENCES customers(id),
    referral_code VARCHAR(20),
    acquisition_source VARCHAR(50),
    preferred_contact contact_preference DEFAULT 'PHONE',
    secondary_contact_name TEXT,
    secondary_contact_phone VARCHAR(50),
    line_id VARCHAR(100),
    whatsapp VARCHAR(50),
    customer_type customer_type DEFAULT 'RETAIL',
    company_registration_no VARCHAR(50),
    vat_registration_no VARCHAR(50),
    account_manager_id INT REFERENCES employees(id),
    customer_group_id INT REFERENCES customer_groups(id),
    is_blocked BOOLEAN DEFAULT false,
    block_reason TEXT,
    blocked_date DATE,
    blocked_by INT REFERENCES employees(id),
    billing_address_id INT,
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
    contact_name TEXT,
    address_line1 TEXT,
    address_line2 TEXT,
    city VARCHAR(100),
    state VARCHAR(100),
    postal_code VARCHAR(20),
    country VARCHAR(100),
    phone VARCHAR(50),
    mobile VARCHAR(50),
    email VARCHAR(255),
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),
    delivery_instructions TEXT,
    is_default BOOLEAN DEFAULT false,
    warehouse_id INT REFERENCES warehouses(id),
    route_id INT,
    is_active BOOLEAN DEFAULT true,
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
-- SECTION 6: VENDORS
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
    contact_name TEXT,
    phone VARCHAR(50),
    mobile VARCHAR(50),
    fax VARCHAR(50),
    email VARCHAR(255),
    tax_code VARCHAR(50),
    rd_branch TEXT,
    picture_name TEXT,
    payment_terms_days INT DEFAULT 30,
    credit_days INT DEFAULT 0,
    discount_percent DECIMAL(5,2) DEFAULT 0.00,
    po_credit_money DECIMAL(12,2) DEFAULT 0.00,
    currency VARCHAR(3) DEFAULT 'LAK',
    lead_time_days INT DEFAULT 7,
    minimum_order DECIMAL(12,2),
    buyer_id INT REFERENCES employees(id),
    vendor_rating DECIMAL(3,2) DEFAULT 5.00,
    quality_score DECIMAL(5,2) DEFAULT 100.00,
    on_time_delivery_rate DECIMAL(5,2) DEFAULT 100.00,
    last_po_date DATE,
    total_po_value DECIMAL(15,2) DEFAULT 0.00,
    total_po_count INT DEFAULT 0,
    website VARCHAR(255),
    bank_name VARCHAR(100),
    bank_account_number VARCHAR(50),
    bank_account_name VARCHAR(100),
    bank_branch VARCHAR(100),
    swift_code VARCHAR(20),
    preferred_payment payment_method DEFAULT 'BANK_TRANSFER',
    certifications TEXT,
    certification_expiry DATE,
    notes TEXT,
    preferred_contact contact_preference DEFAULT 'EMAIL',
    secondary_contact_name TEXT,
    secondary_contact_phone VARCHAR(50),
    vendor_type vendor_type DEFAULT 'SUPPLIER',
    return_policy TEXT,
    return_window_days INT DEFAULT 30,
    is_blocked BOOLEAN DEFAULT false,
    block_reason TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

ALTER TABLE products ADD CONSTRAINT fk_products_primary_vendor 
    FOREIGN KEY (primary_vendor_id) REFERENCES vendors(id);

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
    last_purchase_date DATE,
    last_purchase_price DECIMAL(10,4),
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
-- SECTION 7: INVENTORY
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
-- SECTION 8: ADDRESSES
-- =====================================================

CREATE TABLE IF NOT EXISTS addresses (
    id SERIAL PRIMARY KEY,
    entity_type TEXT NOT NULL,
    entity_id INT NOT NULL,
    address_type VARCHAR(20) DEFAULT 'PRIMARY',
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
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),
    is_default BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- SECTION 9: INDEXES
-- =====================================================

CREATE INDEX idx_employees_email ON employees(email);
CREATE INDEX idx_employees_role ON employees(role_id);
CREATE INDEX idx_employees_department ON employees(department_id);
CREATE INDEX idx_employees_warehouse ON employees(warehouse_id);
CREATE INDEX idx_departments_parent ON departments(parent_id);
CREATE INDEX idx_products_sku ON products(sku);
CREATE INDEX idx_products_barcode ON products(barcode);
CREATE INDEX idx_products_category ON products(category_id);
CREATE INDEX idx_products_brand ON products(brand_id);
CREATE INDEX idx_products_group ON products(product_group_id);
CREATE INDEX idx_products_active ON products(is_active);
CREATE INDEX idx_customers_code ON customers(customer_code);
CREATE INDEX idx_customers_sales_rep ON customers(sales_rep_id);
CREATE INDEX idx_customers_level ON customers(customer_level);
CREATE INDEX idx_customers_price_type ON customers(price_type);
CREATE INDEX idx_customers_group ON customers(customer_group_id);
CREATE INDEX idx_customers_blocked ON customers(is_blocked);
CREATE INDEX idx_customers_active ON customers(is_active);
CREATE INDEX idx_vendors_code ON vendors(vendor_code);
CREATE INDEX idx_vendors_buyer ON vendors(buyer_id);
CREATE INDEX idx_vendors_type ON vendors(vendor_type);
CREATE INDEX idx_vendors_blocked ON vendors(is_blocked);
CREATE INDEX idx_vendors_active ON vendors(is_active);
CREATE INDEX idx_inventory_product ON inventory(product_id);
CREATE INDEX idx_inventory_warehouse ON inventory(warehouse_id);
CREATE INDEX idx_inventory_lot ON inventory(lot_number);
CREATE INDEX idx_inventory_expiry ON inventory(expiry_date);
CREATE INDEX idx_inventory_transactions_product ON inventory_transactions(product_id);
CREATE INDEX idx_inventory_transactions_date ON inventory_transactions(created_at);
CREATE INDEX idx_inventory_transactions_type ON inventory_transactions(transaction_type);
CREATE INDEX idx_addresses_entity ON addresses(entity_type, entity_id);
