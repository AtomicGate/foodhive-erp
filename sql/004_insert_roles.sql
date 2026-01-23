-- =====================================================
-- FoodHive ERP System - Default Roles
-- Version: 2.0 ENHANCED
-- =====================================================
-- Run this after creating the core tables

INSERT INTO roles (role_name, role_desc, is_active) VALUES
    -- System Administration
    ('Super Admin', 'Full system access with all permissions', true),
    ('Administrator', 'System administrator with most permissions', true),
    
    -- Management
    ('Manager', 'Department manager with oversight permissions', true),
    ('General Manager', 'Company-wide management access', true),
    
    -- Sales
    ('Sales Manager', 'Manages sales team and orders', true),
    ('Sales Rep', 'Sales representative - can create and manage sales orders', true),
    ('Sales Analyst', 'Sales reporting and analytics access', true),
    
    -- Purchasing
    ('Purchase Manager', 'Manages purchasing and vendors', true),
    ('Buyer', 'Can create and manage purchase orders', true),
    ('Purchase Analyst', 'Purchasing analytics and reporting', true),
    
    -- Warehouse
    ('Warehouse Manager', 'Manages inventory and warehouse operations', true),
    ('Warehouse Staff', 'Can update inventory and handle picking', true),
    ('Receiving Clerk', 'Handles goods receiving and inspection', true),
    ('Shipping Clerk', 'Handles order shipping and dispatch', true),
    
    -- Finance & Accounting
    ('Finance Manager', 'Full financial management access', true),
    ('Accountant', 'Financial management and reporting access', true),
    ('AR Manager', 'Accounts Receivable management', true),
    ('AR Clerk', 'AR data entry and collection', true),
    ('AP Manager', 'Accounts Payable management', true),
    ('AP Clerk', 'AP data entry and payment processing', true),
    ('GL Accountant', 'General Ledger and financial reporting', true),
    ('Cashier', 'Cash handling and payment collection', true),
    
    -- Human Resources
    ('HR Manager', 'Human resources and employee management', true),
    ('HR Staff', 'Employee data entry and basic HR tasks', true),
    ('Payroll Specialist', 'Payroll processing and management', true),
    
    -- Operations
    ('Operations Manager', 'Operations oversight and management', true),
    ('Customer Service', 'Customer support and order inquiries', true),
    ('Driver', 'Delivery driver - can view routes and update delivery status', true),
    ('Quality Control', 'Product quality inspection and control', true),
    
    -- NEW ROLES for DocText.md modules
    ('Requisition Approver', 'Can approve purchase requisitions', true),
    ('Advance Approver', 'Can approve advance requests and vouchers', true),
    ('Budget Controller', 'Budget management and cost center control', true),
    
    -- Basic Roles
    ('Employee', 'Standard employee with basic view permissions', true),
    ('Viewer', 'Read-only access to most modules', true),
    ('Guest', 'Limited read-only access', true)
ON CONFLICT DO NOTHING;

-- Verify roles were inserted
SELECT id, role_name, role_desc, is_active 
FROM roles 
ORDER BY id;

-- Count roles
SELECT COUNT(*) as total_roles FROM roles WHERE is_active = true;
