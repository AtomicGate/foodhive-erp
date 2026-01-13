-- =====================================================
-- ERP System - Transaction Tables
-- Version: 1.0
-- =====================================================

-- =====================================================
-- SALES ORDERS
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
    notes TEXT,
    UNIQUE(route_id, stop_sequence)
);

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
    line_total DECIMAL(12,2),
    lot_number VARCHAR(50),
    expiry_date DATE,
    catch_weight DECIMAL(10,3),
    cost DECIMAL(10,4),
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
-- PURCHASE ORDERS
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
    subtotal DECIMAL(12,2) DEFAULT 0,
    tax_amount DECIMAL(12,2) DEFAULT 0,
    freight_amount DECIMAL(12,2) DEFAULT 0,
    total_amount DECIMAL(12,2) DEFAULT 0,
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
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS receiving_lines (
    id SERIAL PRIMARY KEY,
    receiving_id INT REFERENCES receiving(id) ON DELETE CASCADE,
    po_line_id INT REFERENCES purchase_order_lines(id),
    product_id INT REFERENCES products(id),
    quantity_received DECIMAL(10,3),
    unit_of_measure VARCHAR(10),
    lot_number VARCHAR(50),
    production_date DATE,
    expiry_date DATE,
    location_code VARCHAR(50),
    unit_cost DECIMAL(10,4),
    is_short_shipment BOOLEAN DEFAULT false,
    notes TEXT
);

-- =====================================================
-- PICKING & SHIPPING
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
-- INDEXES
-- =====================================================

CREATE INDEX idx_sales_orders_customer ON sales_orders(customer_id);
CREATE INDEX idx_sales_orders_date ON sales_orders(order_date);
CREATE INDEX idx_sales_orders_status ON sales_orders(status);
CREATE INDEX idx_sales_orders_warehouse ON sales_orders(warehouse_id);
CREATE INDEX idx_sales_orders_route ON sales_orders(route_id);
CREATE INDEX idx_sales_order_lines_product ON sales_order_lines(product_id);
CREATE INDEX idx_purchase_orders_vendor ON purchase_orders(vendor_id);
CREATE INDEX idx_purchase_orders_status ON purchase_orders(status);
CREATE INDEX idx_purchase_orders_warehouse ON purchase_orders(warehouse_id);
CREATE INDEX idx_pick_lists_date ON pick_lists(pick_date);
CREATE INDEX idx_pick_lists_status ON pick_lists(status);
CREATE INDEX idx_pick_lists_warehouse ON pick_lists(warehouse_id);

