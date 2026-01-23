-- =====================================================
-- FoodHive ERP System - Transaction Tables
-- Version: 2.0 ENHANCED (DocText.md Requirements)
-- =====================================================

-- =====================================================
-- SECTION 1: ROUTES
-- =====================================================

CREATE TABLE IF NOT EXISTS routes (
    id SERIAL PRIMARY KEY,
    route_code VARCHAR(10) UNIQUE NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    warehouse_id INT REFERENCES warehouses(id),
    driver_id INT REFERENCES employees(id),
    vehicle_id INT,
    departure_time TIME,
    estimated_duration_minutes INT,
    distance_km DECIMAL(10,2),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS route_stops (
    id SERIAL PRIMARY KEY,
    route_id INT REFERENCES routes(id) ON DELETE CASCADE,
    customer_id INT REFERENCES customers(id),
    ship_to_id INT REFERENCES customer_ship_to(id),
    stop_sequence INT,
    estimated_arrival TIME,
    service_time_minutes INT DEFAULT 15,
    notes TEXT,
    UNIQUE(route_id, stop_sequence)
);

-- =====================================================
-- SECTION 2: SALES ORDERS
-- =====================================================

CREATE TABLE IF NOT EXISTS sales_orders (
    id SERIAL PRIMARY KEY,
    order_number VARCHAR(20) UNIQUE NOT NULL,
    customer_id INT REFERENCES customers(id),
    ship_to_id INT REFERENCES customer_ship_to(id),
    order_type order_type DEFAULT 'STANDARD',
    order_date DATE DEFAULT CURRENT_DATE,
    requested_ship_date DATE,
    actual_ship_date DATE,
    warehouse_id INT REFERENCES warehouses(id),
    route_id INT REFERENCES routes(id),
    status order_status DEFAULT 'DRAFT',
    subtotal DECIMAL(12,2) DEFAULT 0,
    tax_amount DECIMAL(12,2) DEFAULT 0,
    freight_amount DECIMAL(12,2) DEFAULT 0,
    discount_amount DECIMAL(12,2) DEFAULT 0,
    total_amount DECIMAL(12,2) DEFAULT 0,
    currency VARCHAR(3) DEFAULT 'LAK',
    exchange_rate DECIMAL(12,6) DEFAULT 1.00,
    notes TEXT,
    po_number VARCHAR(50),
    sales_rep_id INT REFERENCES employees(id),
    created_by INT REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS sales_order_lines (
    id SERIAL PRIMARY KEY,
    order_id INT REFERENCES sales_orders(id) ON DELETE CASCADE,
    line_number INT,
    product_id INT REFERENCES products(id),
    description TEXT,
    quantity_ordered DECIMAL(10,3),
    quantity_shipped DECIMAL(10,3) DEFAULT 0,
    unit_of_measure VARCHAR(10),
    unit_price DECIMAL(10,4),
    discount_percent DECIMAL(5,2) DEFAULT 0,
    tax_rate DECIMAL(5,2) DEFAULT 7.00,
    line_total DECIMAL(12,2),
    lot_number VARCHAR(50),
    expiry_date DATE,
    catch_weight DECIMAL(10,3),
    cost DECIMAL(10,4),
    notes TEXT,
    UNIQUE(order_id, line_number)
);

CREATE TABLE IF NOT EXISTS lost_sales (
    id SERIAL PRIMARY KEY,
    order_id INT REFERENCES sales_orders(id),
    product_id INT REFERENCES products(id),
    quantity_requested DECIMAL(10,3),
    quantity_available DECIMAL(10,3),
    reason TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- SECTION 3: PURCHASE ORDERS (ENHANCED)
-- =====================================================

CREATE TABLE IF NOT EXISTS purchase_orders (
    id SERIAL PRIMARY KEY,
    po_number VARCHAR(20) UNIQUE NOT NULL,
    vendor_id INT REFERENCES vendors(id),
    warehouse_id INT REFERENCES warehouses(id),
    order_date DATE DEFAULT CURRENT_DATE,
    expected_date DATE,
    received_date DATE,
    status po_status DEFAULT 'DRAFT',
    
    -- Amounts
    subtotal DECIMAL(12,2) DEFAULT 0,
    tax_amount DECIMAL(12,2) DEFAULT 0,
    freight_amount DECIMAL(12,2) DEFAULT 0,
    total_amount DECIMAL(12,2) DEFAULT 0,
    
    -- REQUIRED: DocText.md fields
    contact_person TEXT,
    customer_name TEXT,
    address TEXT,
    tel VARCHAR(50),
    email VARCHAR(255),
    document_ref VARCHAR(100),
    terms_of_payment TEXT,
    currency currency_code DEFAULT 'LAK',
    trade_discount DECIMAL(12,2) DEFAULT 0.00,
    amount_after_discount DECIMAL(12,2) DEFAULT 0.00,
    vat_rate DECIMAL(5,2) DEFAULT 7.00,
    
    -- REQUIRED: Approval chain
    prepared_by INT REFERENCES employees(id),
    prepared_date DATE,
    checked_by INT REFERENCES employees(id),
    checked_date DATE,
    authorized_by INT REFERENCES employees(id),
    authorized_date DATE,
    
    -- SUGGESTED: Shipping & logistics
    shipping_method VARCHAR(50),
    shipping_carrier VARCHAR(100),
    tracking_number VARCHAR(100),
    insurance_amount DECIMAL(12,2) DEFAULT 0.00,
    delivery_instructions TEXT,
    quality_requirements TEXT,
    
    -- SUGGESTED: Multi-currency
    exchange_rate DECIMAL(12,6) DEFAULT 1.00,
    base_currency VARCHAR(3) DEFAULT 'LAK',
    total_in_base_currency DECIMAL(12,2),
    
    -- SUGGESTED: Priority & linkage
    priority priority_level DEFAULT 'NORMAL',
    requisition_id INT,
    quotation_ref VARCHAR(50),
    contract_ref VARCHAR(50),
    
    -- SUGGESTED: Payment tracking
    payment_status VARCHAR(20) DEFAULT 'PENDING',
    amount_paid DECIMAL(12,2) DEFAULT 0.00,
    
    -- SUGGESTED: Cancellation
    cancelled_by INT REFERENCES employees(id),
    cancelled_date DATE,
    cancellation_reason TEXT,
    
    -- Standard fields
    notes TEXT,
    buyer_id INT REFERENCES employees(id),
    created_by INT REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS purchase_order_lines (
    id SERIAL PRIMARY KEY,
    po_id INT REFERENCES purchase_orders(id) ON DELETE CASCADE,
    line_number INT,
    product_id INT REFERENCES products(id),
    description TEXT,
    quantity_ordered DECIMAL(10,3),
    quantity_received DECIMAL(10,3) DEFAULT 0,
    unit_of_measure VARCHAR(10),
    unit_cost DECIMAL(10,4),
    line_total DECIMAL(12,2),
    expected_date DATE,
    
    -- REQUIRED: Line discount
    discount_percent DECIMAL(5,2) DEFAULT 0.00,
    discount_amount DECIMAL(12,2) DEFAULT 0.00,
    
    -- SUGGESTED: Tax per line
    tax_rate DECIMAL(5,2) DEFAULT 7.00,
    tax_amount DECIMAL(12,2) DEFAULT 0.00,
    
    -- SUGGESTED: Requested delivery date per line
    requested_date DATE,
    
    -- SUGGESTED: Line notes
    notes TEXT,
    
    UNIQUE(po_id, line_number)
);

CREATE TABLE IF NOT EXISTS receiving (
    id SERIAL PRIMARY KEY,
    receiving_number VARCHAR(20) UNIQUE NOT NULL,
    po_id INT REFERENCES purchase_orders(id),
    warehouse_id INT REFERENCES warehouses(id),
    vendor_id INT REFERENCES vendors(id),
    received_date DATE DEFAULT CURRENT_DATE,
    notes TEXT,
    received_by INT REFERENCES employees(id),
    checked_by INT REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS receiving_lines (
    id SERIAL PRIMARY KEY,
    receiving_id INT REFERENCES receiving(id) ON DELETE CASCADE,
    po_line_id INT REFERENCES purchase_order_lines(id),
    product_id INT REFERENCES products(id),
    quantity_received DECIMAL(10,3),
    quantity_rejected DECIMAL(10,3) DEFAULT 0,
    unit_of_measure VARCHAR(10),
    lot_number VARCHAR(50),
    production_date DATE,
    expiry_date DATE,
    location_code VARCHAR(50),
    unit_cost DECIMAL(10,4),
    is_short_shipment BOOLEAN DEFAULT false,
    rejection_reason TEXT,
    notes TEXT
);

-- =====================================================
-- SECTION 4: PURCHASE REQUISITIONS (NEW MODULE)
-- =====================================================

CREATE TABLE IF NOT EXISTS purchase_requisitions (
    id SERIAL PRIMARY KEY,
    requisition_number VARCHAR(20) UNIQUE NOT NULL,
    document_number VARCHAR(50),
    
    -- Requester Information
    requester_id INT REFERENCES employees(id),
    department_id INT REFERENCES departments(id),
    reason TEXT,
    
    -- Supplier Information
    supplier_id INT REFERENCES vendors(id),
    
    -- Dates
    document_date DATE DEFAULT CURRENT_DATE,
    request_date DATE DEFAULT CURRENT_DATE,
    required_date DATE,
    
    -- Amounts
    total_amount DECIMAL(12,2) DEFAULT 0.00,
    
    -- Notes
    remark TEXT,
    
    -- Approval Chain
    checked_by INT REFERENCES employees(id),
    checked_date DATE,
    authorized_by INT REFERENCES employees(id),
    authorized_date DATE,
    purchasing_employee_id INT REFERENCES employees(id),
    purchasing_date DATE,
    
    -- Linked PO
    converted_po_id INT REFERENCES purchase_orders(id),
    
    -- SUGGESTED: Enhanced tracking
    priority priority_level DEFAULT 'NORMAL',
    budget_code VARCHAR(50),
    cost_center VARCHAR(50),
    project_id INT,
    alternative_suppliers TEXT,
    single_source_justification TEXT,
    
    -- Rejection tracking
    rejected_by INT REFERENCES employees(id),
    rejected_date DATE,
    rejection_reason TEXT,
    
    -- Revision tracking
    revision_number INT DEFAULT 1,
    previous_version_id INT REFERENCES purchase_requisitions(id),
    
    -- Status & Audit
    status pr_status DEFAULT 'DRAFT',
    created_by INT REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Link PO to PR
ALTER TABLE purchase_orders ADD CONSTRAINT fk_po_requisition 
    FOREIGN KEY (requisition_id) REFERENCES purchase_requisitions(id);

CREATE TABLE IF NOT EXISTS purchase_requisition_lines (
    id SERIAL PRIMARY KEY,
    requisition_id INT REFERENCES purchase_requisitions(id) ON DELETE CASCADE,
    line_number INT NOT NULL,
    
    -- Product Information
    product_id INT REFERENCES products(id),
    product_code VARCHAR(50),
    description TEXT,
    
    -- Stock & Quantity
    stock_balance DECIMAL(12,3) DEFAULT 0.00,
    quantity_requested DECIMAL(10,3) NOT NULL,
    quantity_approved DECIMAL(10,3),
    
    -- Pricing
    unit_of_measure VARCHAR(10),
    unit_price DECIMAL(10,4) DEFAULT 0.00,
    amount DECIMAL(12,2) DEFAULT 0.00,
    
    -- Preferred vendor
    preferred_vendor_id INT REFERENCES vendors(id),
    
    -- Notes
    notes TEXT,
    
    UNIQUE(requisition_id, line_number)
);

-- =====================================================
-- SECTION 5: ADVANCE REQUESTS (NEW MODULE)
-- =====================================================

CREATE TABLE IF NOT EXISTS advance_requests (
    id SERIAL PRIMARY KEY,
    request_number VARCHAR(20) UNIQUE NOT NULL,
    
    -- Employee Information
    employee_id INT REFERENCES employees(id) NOT NULL,
    employee_name TEXT,
    position TEXT,
    
    -- Reference
    po_number VARCHAR(50),
    
    -- Dates
    request_date DATE DEFAULT CURRENT_DATE,
    expected_settlement_date DATE,
    
    -- Currency & Amount
    currency currency_code DEFAULT 'LAK',
    total_amount DECIMAL(12,2) DEFAULT 0.00,
    exchange_rate DECIMAL(12,6) DEFAULT 1.00,
    total_in_base_currency DECIMAL(12,2),
    
    -- Main Description
    description TEXT,
    purpose_category VARCHAR(50),
    
    -- Project/cost center
    project_id INT,
    cost_center VARCHAR(50),
    budget_code VARCHAR(50),
    
    -- Bank transfer details
    bank_account_number VARCHAR(50),
    bank_name VARCHAR(100),
    
    -- Approval Chain (3 signatures)
    approved_by_1 INT REFERENCES employees(id),
    approved_date_1 DATE,
    approved_by_2 INT REFERENCES employees(id),
    approved_date_2 DATE,
    approved_by_3 INT REFERENCES employees(id),
    approved_date_3 DATE,
    
    -- Rejection tracking
    rejected_by INT REFERENCES employees(id),
    rejected_date DATE,
    rejection_reason TEXT,
    
    -- Linked voucher
    voucher_id INT,
    
    -- Status & Audit
    status advance_request_status DEFAULT 'DRAFT',
    created_by INT REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS advance_request_lines (
    id SERIAL PRIMARY KEY,
    request_id INT REFERENCES advance_requests(id) ON DELETE CASCADE,
    line_number INT NOT NULL,
    
    -- Hierarchy support
    parent_line_id INT REFERENCES advance_request_lines(id),
    line_type VARCHAR(20) DEFAULT 'MAIN',
    
    -- Content
    description TEXT NOT NULL,
    amount DECIMAL(12,2) DEFAULT 0.00,
    expense_category VARCHAR(50),
    
    UNIQUE(request_id, line_number)
);

-- =====================================================
-- SECTION 6: ADVANCE VOUCHERS (NEW MODULE)
-- =====================================================

CREATE TABLE IF NOT EXISTS advance_vouchers (
    id SERIAL PRIMARY KEY,
    voucher_number VARCHAR(20) UNIQUE NOT NULL,
    
    -- Link to Advance Request
    advance_request_id INT REFERENCES advance_requests(id),
    
    -- Employee Information
    employee_id INT REFERENCES employees(id) NOT NULL,
    employee_name TEXT,
    position TEXT,
    
    -- Reference
    po_number VARCHAR(50),
    
    -- Dates
    voucher_date DATE DEFAULT CURRENT_DATE,
    
    -- Currency
    currency currency_code DEFAULT 'LAK',
    exchange_rate DECIMAL(12,6) DEFAULT 1.00,
    
    -- Amounts
    advance_amount DECIMAL(12,2) DEFAULT 0.00,
    expenditure_amount DECIMAL(12,2) DEFAULT 0.00,
    balance_amount DECIMAL(12,2) GENERATED ALWAYS AS (advance_amount - expenditure_amount) STORED,
    
    -- Approval
    accountant_id INT REFERENCES employees(id),
    accountant_date DATE,
    returned_by_id INT REFERENCES employees(id),
    returned_date DATE,
    
    -- Balance handling
    balance_returned BOOLEAN DEFAULT false,
    balance_return_date DATE,
    balance_received_by INT REFERENCES employees(id),
    
    -- Additional funds
    additional_amount_requested DECIMAL(12,2) DEFAULT 0.00,
    additional_approved_by INT REFERENCES employees(id),
    additional_approved_date DATE,
    
    -- Receipt tracking
    total_receipt_count INT DEFAULT 0,
    
    -- Verification
    verified_by INT REFERENCES employees(id),
    verified_date DATE,
    verification_notes TEXT,
    
    -- Rejection
    rejected_by INT REFERENCES employees(id),
    rejected_date DATE,
    rejection_reason TEXT,
    
    -- Status & Audit
    status advance_voucher_status DEFAULT 'DRAFT',
    created_by INT REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Link advance request to voucher
ALTER TABLE advance_requests ADD CONSTRAINT fk_advance_request_voucher 
    FOREIGN KEY (voucher_id) REFERENCES advance_vouchers(id);

CREATE TABLE IF NOT EXISTS advance_voucher_lines (
    id SERIAL PRIMARY KEY,
    voucher_id INT REFERENCES advance_vouchers(id) ON DELETE CASCADE,
    line_number INT NOT NULL,
    
    -- Line Type
    line_type VARCHAR(20) NOT NULL,
    
    -- Hierarchy
    parent_line_id INT REFERENCES advance_voucher_lines(id),
    
    -- Content
    description TEXT NOT NULL,
    amount DECIMAL(12,2) DEFAULT 0.00,
    
    -- Receipt reference
    receipt_number VARCHAR(50),
    receipt_date DATE,
    vendor_name TEXT,
    tax_amount DECIMAL(12,2) DEFAULT 0.00,
    
    UNIQUE(voucher_id, line_number)
);

CREATE TABLE IF NOT EXISTS advance_voucher_documents (
    id SERIAL PRIMARY KEY,
    voucher_id INT REFERENCES advance_vouchers(id) ON DELETE CASCADE,
    
    -- Document Type
    document_type voucher_document_type NOT NULL,
    
    -- Status
    has_document BOOLEAN DEFAULT false,
    
    -- File attachment
    document_path TEXT,
    document_name TEXT,
    file_size INT,
    mime_type VARCHAR(100),
    
    -- Verification
    verified BOOLEAN DEFAULT false,
    verified_by INT REFERENCES employees(id),
    verified_date DATE,
    
    created_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- SECTION 7: PICKING & SHIPPING
-- =====================================================

CREATE TABLE IF NOT EXISTS pick_lists (
    id SERIAL PRIMARY KEY,
    pick_number VARCHAR(20) UNIQUE NOT NULL,
    warehouse_id INT REFERENCES warehouses(id),
    route_id INT REFERENCES routes(id),
    pick_date DATE DEFAULT CURRENT_DATE,
    status pick_list_status DEFAULT 'PENDING',
    picker_id INT REFERENCES employees(id),
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_by INT REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS pick_list_lines (
    id SERIAL PRIMARY KEY,
    pick_list_id INT REFERENCES pick_lists(id) ON DELETE CASCADE,
    order_id INT REFERENCES sales_orders(id),
    order_line_id INT REFERENCES sales_order_lines(id),
    product_id INT REFERENCES products(id),
    location_code VARCHAR(50),
    lot_number VARCHAR(50),
    quantity_to_pick DECIMAL(10,3),
    quantity_picked DECIMAL(10,3) DEFAULT 0,
    picked_by INT REFERENCES employees(id),
    picked_at TIMESTAMP,
    notes TEXT
);

CREATE TABLE IF NOT EXISTS pick_zones (
    id SERIAL PRIMARY KEY,
    warehouse_id INT REFERENCES warehouses(id),
    zone_code VARCHAR(10) NOT NULL,
    name TEXT,
    pick_sequence INT,
    UNIQUE(warehouse_id, zone_code)
);

-- =====================================================
-- SECTION 8: INDEXES
-- =====================================================

-- Sales Orders
CREATE INDEX idx_sales_orders_customer ON sales_orders(customer_id);
CREATE INDEX idx_sales_orders_date ON sales_orders(order_date);
CREATE INDEX idx_sales_orders_status ON sales_orders(status);
CREATE INDEX idx_sales_orders_warehouse ON sales_orders(warehouse_id);
CREATE INDEX idx_sales_orders_route ON sales_orders(route_id);
CREATE INDEX idx_sales_order_lines_product ON sales_order_lines(product_id);

-- Purchase Orders
CREATE INDEX idx_purchase_orders_vendor ON purchase_orders(vendor_id);
CREATE INDEX idx_purchase_orders_status ON purchase_orders(status);
CREATE INDEX idx_purchase_orders_warehouse ON purchase_orders(warehouse_id);
CREATE INDEX idx_purchase_orders_date ON purchase_orders(order_date);
CREATE INDEX idx_purchase_orders_priority ON purchase_orders(priority);

-- Purchase Requisitions
CREATE INDEX idx_purchase_requisitions_status ON purchase_requisitions(status);
CREATE INDEX idx_purchase_requisitions_requester ON purchase_requisitions(requester_id);
CREATE INDEX idx_purchase_requisitions_date ON purchase_requisitions(document_date);
CREATE INDEX idx_purchase_requisitions_priority ON purchase_requisitions(priority);

-- Advance Requests
CREATE INDEX idx_advance_requests_status ON advance_requests(status);
CREATE INDEX idx_advance_requests_employee ON advance_requests(employee_id);
CREATE INDEX idx_advance_requests_date ON advance_requests(request_date);

-- Advance Vouchers
CREATE INDEX idx_advance_vouchers_status ON advance_vouchers(status);
CREATE INDEX idx_advance_vouchers_employee ON advance_vouchers(employee_id);
CREATE INDEX idx_advance_vouchers_date ON advance_vouchers(voucher_date);

-- Pick Lists
CREATE INDEX idx_pick_lists_date ON pick_lists(pick_date);
CREATE INDEX idx_pick_lists_status ON pick_lists(status);
CREATE INDEX idx_pick_lists_warehouse ON pick_lists(warehouse_id);
