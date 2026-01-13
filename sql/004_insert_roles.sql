-- =====================================================
-- INSERT DEFAULT ROLES
-- =====================================================
-- This script creates standard roles for the ERP system
-- Run this after creating the roles table

-- Insert roles (skip if already exists)
INSERT INTO roles (role_name, role_desc, is_active) VALUES
    ('Super Admin', 'Full system access with all permissions', true),
    ('Administrator', 'System administrator with most permissions', true),
    ('Manager', 'Department manager with oversight permissions', true),
    ('Sales Manager', 'Manages sales team and orders', true),
    ('Sales Rep', 'Sales representative - can create and manage sales orders', true),
    ('Purchase Manager', 'Manages purchasing and vendors', true),
    ('Buyer', 'Can create and manage purchase orders', true),
    ('Warehouse Manager', 'Manages inventory and warehouse operations', true),
    ('Warehouse Staff', 'Can update inventory and handle picking', true),
    ('Accountant', 'Financial management and reporting access', true),
    ('AR Manager', 'Accounts Receivable management', true),
    ('AP Manager', 'Accounts Payable management', true),
    ('GL Accountant', 'General Ledger and financial reporting', true),
    ('HR Manager', 'Human resources and employee management', true),
    ('HR Staff', 'Employee data entry and basic HR tasks', true),
    ('Customer Service', 'Customer support and order inquiries', true),
    ('Driver', 'Delivery driver - can view routes and update delivery status', true),
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
