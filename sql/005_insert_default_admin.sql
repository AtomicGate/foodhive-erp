-- =====================================================
-- FoodHive ERP System - Default Admin Employee
-- Version: 2.0 ENHANCED
-- =====================================================
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
    department_id,
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
    1,  -- Default department
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

-- Create employee details for admin
INSERT INTO employee_details (employee_id, job_title, notes)
SELECT 1, 'System Administrator', 'Default admin account - Please update password!'
WHERE EXISTS (SELECT 1 FROM employees WHERE id = 1)
  AND NOT EXISTS (SELECT 1 FROM employee_details WHERE employee_id = 1);

-- Create employee finances record for admin
INSERT INTO employee_finances (employee_id)
SELECT 1
WHERE EXISTS (SELECT 1 FROM employees WHERE id = 1)
  AND NOT EXISTS (SELECT 1 FROM employee_finances WHERE employee_id = 1);

-- Verify admin employee was created
SELECT 
    e.id, 
    e.email, 
    e.english_name, 
    r.role_name,
    e.status,
    e.account_status,
    e.security_level
FROM employees e
LEFT JOIN roles r ON e.role_id = r.id
WHERE e.id = 1;

-- Give admin full permissions on all pages
INSERT INTO emp_page (user_id, page_id, can_create, can_update, can_delete, can_view)
SELECT 1, p.id, true, true, true, true
FROM pages p
WHERE EXISTS (SELECT 1 FROM employees WHERE id = 1)
  AND NOT EXISTS (
    SELECT 1 FROM emp_page ep WHERE ep.user_id = 1 AND ep.page_id = p.id
);

-- Show admin permissions
SELECT 
    p.page_name,
    p.route_name,
    ep.can_create,
    ep.can_update,
    ep.can_delete,
    ep.can_view
FROM emp_page ep
JOIN pages p ON ep.page_id = p.id
WHERE ep.user_id = 1
ORDER BY p.page_name;
