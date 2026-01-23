-- =====================================================
-- FoodHive ERP System - Admin Page Permissions
-- Run this AFTER creating the admin user
-- =====================================================

-- Step 1: Insert pages if they don't exist
INSERT INTO pages (page_name, route_name, icon, display_order) VALUES
('Dashboard', '/', 'LayoutDashboard', 1),
('Customers', '/customers', 'ShoppingBag', 2),
('Vendors', '/vendors', 'Truck', 3),
('Products', '/products', 'Package', 4),
('Employees', '/employees', 'Users', 5),
('Sales Orders', '/sales-orders', 'ClipboardList', 6),
('Purchase Orders', '/purchase-orders', 'ShoppingCart', 7),
('Inventory', '/inventory', 'Boxes', 8),
('AR Dashboard', '/financials/ar', 'Receipt', 9),
('AP Dashboard', '/financials/ap', 'CreditCard', 10),
('General Ledger', '/gl', 'BookOpen', 11),
('Pricing', '/pricing', 'Tag', 12),
('Catch Weight', '/operations/catch-weight', 'Scale', 13),
('Departments', '/admin/departments', 'Building2', 14),
('Roles', '/admin/roles', 'ShieldCheck', 15),
('Warehouses', '/admin/warehouses', 'Warehouse', 16)
ON CONFLICT (route_name) DO NOTHING;

-- Step 2: Give admin full permissions on ALL pages
INSERT INTO emp_page (user_id, page_id, can_create, can_update, can_delete, can_view)
SELECT e.id, p.id, true, true, true, true
FROM employees e
CROSS JOIN pages p
WHERE e.email = 'admin@foodhive.com'
ON CONFLICT (user_id, page_id) DO UPDATE SET 
  can_create = true,
  can_update = true,
  can_delete = true,
  can_view = true;

-- Step 3: Verify permissions were added
SELECT 
  e.email, 
  COUNT(ep.id) as total_pages,
  'Full access granted' as status
FROM emp_page ep 
JOIN employees e ON ep.user_id = e.id 
WHERE e.email = 'admin@foodhive.com'
GROUP BY e.email;
