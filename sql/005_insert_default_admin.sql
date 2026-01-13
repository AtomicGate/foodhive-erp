-- =====================================================
-- INSERT DEFAULT ADMIN EMPLOYEE
-- =====================================================
-- This script creates a default admin employee for initial setup
-- Run this after 004_insert_roles.sql
-- IMPORTANT: Change the password after first login!

-- Insert default admin employee (only if doesn't exist)
-- Password: 'admin123' (CHANGE THIS AFTER FIRST LOGIN!)
-- Note: In production, passwords should be hashed with bcrypt
INSERT INTO employees (
    id,
    email, 
    password, 
    english_name, 
    arabic_name,
    role_id,
    status,
    account_status,
    security_level,
    is_sales_rep,
    is_buyer,
    is_driver
) 
SELECT 
    1,
    'admin@foodhive.com',
    'admin123',  -- TODO: Hash this with bcrypt in production
    'System Administrator',
    'مدير النظام',
    (SELECT id FROM roles WHERE role_name = 'Super Admin' LIMIT 1),
    'CONTINUED',
    'active',
    10,  -- Highest security level
    false,
    false,
    false
WHERE NOT EXISTS (SELECT 1 FROM employees WHERE id = 1)
  AND EXISTS (SELECT 1 FROM roles WHERE role_name = 'Super Admin');

-- Reset the sequence for employees to ensure next insert gets id=2
SELECT setval('employees_id_seq', COALESCE((SELECT MAX(id) FROM employees), 1), true);

-- Verify admin employee was created
SELECT 
    id, 
    email, 
    english_name, 
    role_id,
    status,
    account_status
FROM employees 
WHERE id = 1;
