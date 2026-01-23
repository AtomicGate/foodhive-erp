-- =====================================================
-- FoodHive ERP System - Schema Updates & Additional Tables
-- Version: 2.0 ENHANCED (DocText.md Requirements)
-- =====================================================
-- EXECUTION ORDER:
-- 1. 001_core_tables.sql
-- 2. 002_transactions_tables.sql
-- 3. 003_schema_updates.sql (this file)
-- 4. 004_insert_roles.sql
-- 5. 005_insert_default_admin.sql

-- =====================================================
-- SECTION 1: AUDIT LOG TABLE (Suggested)
-- =====================================================

CREATE TABLE IF NOT EXISTS audit_log (
    id SERIAL PRIMARY KEY,
    table_name VARCHAR(100) NOT NULL,
    record_id INT NOT NULL,
    action VARCHAR(20) NOT NULL,
    old_values JSONB,
    new_values JSONB,
    changed_by INT REFERENCES employees(id),
    changed_at TIMESTAMP DEFAULT NOW(),
    ip_address VARCHAR(45),
    user_agent TEXT
);

CREATE INDEX idx_audit_log_table ON audit_log(table_name);
CREATE INDEX idx_audit_log_record ON audit_log(table_name, record_id);
CREATE INDEX idx_audit_log_date ON audit_log(changed_at);
CREATE INDEX idx_audit_log_user ON audit_log(changed_by);

-- =====================================================
-- SECTION 2: DOCUMENT ATTACHMENTS TABLE (Suggested)
-- =====================================================

CREATE TABLE IF NOT EXISTS document_attachments (
    id SERIAL PRIMARY KEY,
    entity_type VARCHAR(50) NOT NULL,
    entity_id INT NOT NULL,
    file_name TEXT NOT NULL,
    file_path TEXT NOT NULL,
    file_size INT,
    mime_type VARCHAR(100),
    document_type VARCHAR(50),
    description TEXT,
    uploaded_by INT REFERENCES employees(id),
    uploaded_at TIMESTAMP DEFAULT NOW(),
    is_active BOOLEAN DEFAULT true
);

CREATE INDEX idx_document_attachments_entity ON document_attachments(entity_type, entity_id);

-- =====================================================
-- SECTION 3: NOTIFICATIONS TABLE (Suggested)
-- =====================================================

CREATE TABLE IF NOT EXISTS notifications (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES employees(id) NOT NULL,
    title TEXT NOT NULL,
    message TEXT NOT NULL,
    notification_type VARCHAR(50),
    entity_type VARCHAR(50),
    entity_id INT,
    link_url TEXT,
    is_read BOOLEAN DEFAULT false,
    read_at TIMESTAMP,
    priority priority_level DEFAULT 'NORMAL',
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP
);

CREATE INDEX idx_notifications_user ON notifications(user_id);
CREATE INDEX idx_notifications_unread ON notifications(user_id, is_read) WHERE is_read = false;
CREATE INDEX idx_notifications_date ON notifications(created_at);

-- =====================================================
-- SECTION 4: APPROVAL WORKFLOWS (Suggested)
-- =====================================================

CREATE TABLE IF NOT EXISTS approval_workflows (
    id SERIAL PRIMARY KEY,
    workflow_name VARCHAR(100) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS approval_workflow_steps (
    id SERIAL PRIMARY KEY,
    workflow_id INT REFERENCES approval_workflows(id) ON DELETE CASCADE,
    step_number INT NOT NULL,
    step_name VARCHAR(100) NOT NULL,
    approver_role_id INT REFERENCES roles(id),
    approver_department_id INT REFERENCES departments(id),
    specific_approver_id INT REFERENCES employees(id),
    min_amount DECIMAL(12,2),
    max_amount DECIMAL(12,2),
    can_skip BOOLEAN DEFAULT false,
    auto_approve_after_days INT,
    UNIQUE(workflow_id, step_number)
);

CREATE TABLE IF NOT EXISTS approval_history (
    id SERIAL PRIMARY KEY,
    workflow_id INT REFERENCES approval_workflows(id),
    entity_type VARCHAR(50) NOT NULL,
    entity_id INT NOT NULL,
    step_number INT NOT NULL,
    action approval_action NOT NULL,
    approved_by INT REFERENCES employees(id),
    approved_at TIMESTAMP DEFAULT NOW(),
    comments TEXT,
    next_step INT
);

CREATE INDEX idx_approval_history_entity ON approval_history(entity_type, entity_id);

-- =====================================================
-- SECTION 5: INSERT DEFAULT DATA
-- =====================================================

-- Insert default department
INSERT INTO departments (id, name, description) 
VALUES (1, 'Administration', 'Administrative department')
ON CONFLICT (id) DO NOTHING;

SELECT setval('departments_id_seq', COALESCE((SELECT MAX(id) FROM departments), 1));

-- =====================================================
-- SECTION 6: INSERT PAGES FOR ALL MODULES
-- =====================================================

INSERT INTO pages (page_name, route_name, icon) VALUES 
    -- Core modules
    ('Dashboard', '/dashboard', 'home'),
    ('Sales Orders', '/sales-orders', 'shopping-cart'),
    ('Purchase Orders', '/purchase-orders', 'truck'),
    ('Inventory', '/inventory', 'package'),
    ('Products', '/products', 'box'),
    ('Product Categories', '/product-categories', 'list'),
    ('Product Groups', '/product-groups', 'layers'),
    ('Brands', '/brands', 'tag'),
    ('Customers', '/customers', 'users'),
    ('Customer Groups', '/customer-groups', 'user-group'),
    ('Vendors', '/vendors', 'building'),
    ('Employees', '/employees', 'user'),
    ('Departments', '/departments', 'briefcase'),
    ('Roles', '/roles', 'shield'),
    ('Warehouses', '/warehouses', 'warehouse'),
    
    -- Financial modules
    ('AR', '/ar', 'dollar-sign'),
    ('AP', '/ap', 'credit-card'),
    ('GL', '/gl', 'book'),
    
    -- Operations
    ('Pricing', '/pricing', 'tag'),
    ('Picking', '/picking', 'clipboard'),
    ('Catch Weight', '/catch-weight', 'scale'),
    ('Routes', '/routes', 'map'),
    
    -- NEW MODULES from DocText.md
    ('Purchase Requisitions', '/purchase-requisitions', 'file-text'),
    ('Advance Requests', '/advance-requests', 'wallet'),
    ('Advance Vouchers', '/advance-vouchers', 'receipt'),
    
    -- Suggested modules
    ('Audit Log', '/audit-log', 'history'),
    ('Notifications', '/notifications', 'bell'),
    ('Approval Workflows', '/approval-workflows', 'git-branch'),
    ('Documents', '/documents', 'folder'),
    ('Settings', '/settings', 'settings'),
    ('Reports', '/reports', 'file-chart')
ON CONFLICT (route_name) DO NOTHING;

-- =====================================================
-- SECTION 7: VIEWS FOR ENHANCED DISPLAY
-- =====================================================

-- Inventory Display View
CREATE OR REPLACE VIEW inventory_display AS
SELECT 
    i.id,
    i.product_id,
    p.sku as product_code,
    p.barcode,
    p.name as description,
    i.quantity_on_hand as stock_balance,
    p.base_unit as unit,
    p.price_1,
    p.price_2,
    p.price_3,
    p.wholesale_price,
    p.retail_price,
    p.vat_percent,
    i.expiry_date,
    i.lot_number,
    i.warehouse_id,
    w.name as warehouse_name,
    pc.name as category_name,
    pg.name as product_group_name,
    b.name as brand_name,
    i.average_cost,
    (i.quantity_on_hand * COALESCE(i.average_cost, 0)) as inventory_value,
    CASE 
        WHEN i.expiry_date IS NOT NULL 
        THEN i.expiry_date - CURRENT_DATE 
        ELSE NULL 
    END as days_to_expiry
FROM inventory i
JOIN products p ON i.product_id = p.id
LEFT JOIN warehouses w ON i.warehouse_id = w.id
LEFT JOIN product_categories pc ON p.category_id = pc.id
LEFT JOIN product_groups pg ON p.product_group_id = pg.id
LEFT JOIN brands b ON p.brand_id = b.id;

-- Customer Summary View
CREATE OR REPLACE VIEW customer_summary AS
SELECT 
    c.id,
    c.customer_code,
    c.name,
    c.default_contact_name,
    c.tel,
    c.mobile,
    c.email,
    c.full_address,
    c.price_type,
    c.customer_level,
    c.customer_type,
    c.credit_limit,
    c.current_balance,
    c.credit_days,
    c.is_blocked,
    c.is_active,
    e.english_name as sales_rep_name,
    w.name as warehouse_name,
    cg.name as customer_group_name,
    COALESCE(c.total_orders, 0) as order_count,
    COALESCE(c.total_spent, 0) as lifetime_value,
    c.credit_limit - c.current_balance as available_credit,
    c.last_order_date,
    c.loyalty_points
FROM customers c
LEFT JOIN employees e ON c.sales_rep_id = e.id
LEFT JOIN warehouses w ON c.default_warehouse_id = w.id
LEFT JOIN customer_groups cg ON c.customer_group_id = cg.id;

-- Vendor Summary View
CREATE OR REPLACE VIEW vendor_summary AS
SELECT 
    v.id,
    v.vendor_code,
    v.name,
    v.contact_name,
    v.phone,
    v.mobile,
    v.email,
    v.vendor_type,
    v.vendor_rating,
    v.quality_score,
    v.on_time_delivery_rate,
    v.credit_days,
    v.is_blocked,
    v.is_active,
    e.english_name as buyer_name,
    COALESCE(v.total_po_count, 0) as order_count,
    COALESCE(v.total_po_value, 0) as total_purchases,
    v.last_po_date
FROM vendors v
LEFT JOIN employees e ON v.buyer_id = e.id;

-- Purchase Requisition Summary View
CREATE OR REPLACE VIEW purchase_requisition_summary AS
SELECT 
    pr.id,
    pr.requisition_number,
    pr.document_number,
    pr.document_date,
    pr.request_date,
    pr.required_date,
    pr.status,
    pr.priority,
    pr.total_amount,
    pr.reason,
    pr.remark,
    e.english_name as requester_name,
    d.name as department_name,
    v.name as supplier_name,
    ce.english_name as checked_by_name,
    ae.english_name as authorized_by_name,
    po.po_number as converted_po_number
FROM purchase_requisitions pr
LEFT JOIN employees e ON pr.requester_id = e.id
LEFT JOIN departments d ON pr.department_id = d.id
LEFT JOIN vendors v ON pr.supplier_id = v.id
LEFT JOIN employees ce ON pr.checked_by = ce.id
LEFT JOIN employees ae ON pr.authorized_by = ae.id
LEFT JOIN purchase_orders po ON pr.converted_po_id = po.id;

-- Advance Request Summary View
CREATE OR REPLACE VIEW advance_request_summary AS
SELECT 
    ar.id,
    ar.request_number,
    ar.request_date,
    ar.status,
    ar.currency,
    ar.total_amount,
    ar.exchange_rate,
    ar.total_in_base_currency,
    ar.description,
    ar.purpose_category,
    ar.po_number,
    e.english_name as employee_name,
    ar.position,
    a1.english_name as approved_by_1_name,
    ar.approved_date_1,
    a2.english_name as approved_by_2_name,
    ar.approved_date_2,
    a3.english_name as approved_by_3_name,
    ar.approved_date_3,
    av.voucher_number as linked_voucher
FROM advance_requests ar
LEFT JOIN employees e ON ar.employee_id = e.id
LEFT JOIN employees a1 ON ar.approved_by_1 = a1.id
LEFT JOIN employees a2 ON ar.approved_by_2 = a2.id
LEFT JOIN employees a3 ON ar.approved_by_3 = a3.id
LEFT JOIN advance_vouchers av ON ar.voucher_id = av.id;

-- Advance Voucher Summary View
CREATE OR REPLACE VIEW advance_voucher_summary AS
SELECT 
    av.id,
    av.voucher_number,
    av.voucher_date,
    av.status,
    av.currency,
    av.advance_amount,
    av.expenditure_amount,
    av.balance_amount,
    av.po_number,
    e.english_name as employee_name,
    av.position,
    ar.request_number as advance_request_number,
    acc.english_name as accountant_name,
    av.accountant_date,
    ret.english_name as returned_by_name,
    av.returned_date,
    av.total_receipt_count,
    av.balance_returned
FROM advance_vouchers av
LEFT JOIN employees e ON av.employee_id = e.id
LEFT JOIN advance_requests ar ON av.advance_request_id = ar.id
LEFT JOIN employees acc ON av.accountant_id = acc.id
LEFT JOIN employees ret ON av.returned_by_id = ret.id;

-- Purchase Order Summary View
CREATE OR REPLACE VIEW purchase_order_summary AS
SELECT 
    po.id,
    po.po_number,
    po.order_date,
    po.expected_date,
    po.status,
    po.priority,
    po.currency,
    po.subtotal,
    po.tax_amount,
    po.freight_amount,
    po.trade_discount,
    po.total_amount,
    po.payment_status,
    po.amount_paid,
    v.name as vendor_name,
    v.contact_name as vendor_contact,
    w.name as warehouse_name,
    pe.english_name as prepared_by_name,
    po.prepared_date,
    ce.english_name as checked_by_name,
    po.checked_date,
    ae.english_name as authorized_by_name,
    po.authorized_date,
    pr.requisition_number as linked_requisition
FROM purchase_orders po
LEFT JOIN vendors v ON po.vendor_id = v.id
LEFT JOIN warehouses w ON po.warehouse_id = w.id
LEFT JOIN employees pe ON po.prepared_by = pe.id
LEFT JOIN employees ce ON po.checked_by = ce.id
LEFT JOIN employees ae ON po.authorized_by = ae.id
LEFT JOIN purchase_requisitions pr ON po.requisition_id = pr.id;

-- =====================================================
-- SECTION 8: GIVE ADMIN PERMISSIONS
-- =====================================================

-- Give admin user (id=1) full permissions on all pages
INSERT INTO emp_page (user_id, page_id, can_create, can_update, can_delete, can_view)
SELECT 1, p.id, true, true, true, true
FROM pages p
WHERE EXISTS (SELECT 1 FROM employees WHERE id = 1)
  AND NOT EXISTS (
    SELECT 1 FROM emp_page ep WHERE ep.user_id = 1 AND ep.page_id = p.id
);
