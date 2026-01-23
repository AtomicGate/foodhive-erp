-- =====================================================
-- FoodHive ERP System - Sample Inventory Data
-- For Demo/Testing Purposes
-- =====================================================

-- Step 1: Insert Warehouses
INSERT INTO warehouses (warehouse_code, name, address_line1, city, country, phone, is_active) VALUES
('WH-MAIN', 'Main Warehouse', '123 Industrial Zone', 'Vientiane', 'Laos', '+856 21 555 0001', true),
('WH-COLD', 'Cold Storage', '45 Cold Chain Road', 'Vientiane', 'Laos', '+856 21 555 0002', true),
('WH-DRY', 'Dry Goods Warehouse', '78 Logistics Park', 'Vientiane', 'Laos', '+856 21 555 0003', true)
ON CONFLICT (warehouse_code) DO NOTHING;

-- Step 2: Insert Product Categories
INSERT INTO product_categories (code, name, description, is_active) VALUES
('MEAT', 'Meat Products', 'Fresh and frozen meat products', true),
('SEAFOOD', 'Seafood', 'Fresh and frozen seafood', true),
('DAIRY', 'Dairy Products', 'Milk, cheese, and dairy items', true),
('PRODUCE', 'Fresh Produce', 'Fruits and vegetables', true),
('DRY', 'Dry Goods', 'Rice, noodles, and dry products', true),
('FROZEN', 'Frozen Foods', 'Frozen prepared foods', true),
('BEVERAGE', 'Beverages', 'Drinks and beverages', true),
('SPICES', 'Spices & Seasonings', 'Cooking spices and seasonings', true)
ON CONFLICT (code) DO NOTHING;

-- Step 3: Insert Products
INSERT INTO products (sku, barcode, name, description, category_id, base_unit, is_lot_tracked, shelf_life_days, price_1, price_2, price_3, wholesale_price, retail_price, standard_cost, reorder_point, reorder_quantity, is_active) VALUES
-- Meat Products
('BEEF-001', '8851234567890', 'Premium Beef Tenderloin', 'High quality beef tenderloin, fresh', (SELECT id FROM product_categories WHERE code = 'MEAT'), 'KG', true, 7, 450000, 420000, 400000, 380000, 480000, 350000, 50, 100, true),
('PORK-001', '8851234567891', 'Pork Loin', 'Fresh pork loin cut', (SELECT id FROM product_categories WHERE code = 'MEAT'), 'KG', true, 5, 180000, 170000, 160000, 150000, 200000, 130000, 80, 150, true),
('CHKN-001', '8851234567892', 'Whole Chicken', 'Fresh whole chicken', (SELECT id FROM product_categories WHERE code = 'MEAT'), 'KG', true, 5, 95000, 90000, 85000, 80000, 100000, 70000, 100, 200, true),
('CHKN-002', '8851234567893', 'Chicken Breast', 'Boneless chicken breast', (SELECT id FROM product_categories WHERE code = 'MEAT'), 'KG', true, 5, 145000, 140000, 135000, 130000, 160000, 110000, 60, 120, true),

-- Seafood
('FISH-001', '8851234567894', 'Fresh Salmon Fillet', 'Norwegian salmon fillet', (SELECT id FROM product_categories WHERE code = 'SEAFOOD'), 'KG', true, 3, 850000, 820000, 800000, 780000, 900000, 700000, 30, 50, true),
('SHRP-001', '8851234567895', 'Tiger Prawns', 'Large tiger prawns, fresh', (SELECT id FROM product_categories WHERE code = 'SEAFOOD'), 'KG', true, 3, 550000, 530000, 510000, 490000, 600000, 420000, 40, 80, true),
('CRAB-001', '8851234567896', 'Blue Crab', 'Live blue crab', (SELECT id FROM product_categories WHERE code = 'SEAFOOD'), 'KG', true, 2, 380000, 360000, 350000, 340000, 400000, 280000, 25, 50, true),

-- Dairy
('MILK-001', '8851234567897', 'Fresh Milk 1L', 'Fresh pasteurized milk', (SELECT id FROM product_categories WHERE code = 'DAIRY'), 'BTL', true, 14, 25000, 24000, 23000, 22000, 28000, 18000, 200, 400, true),
('CHSE-001', '8851234567898', 'Cheddar Cheese 500g', 'Aged cheddar cheese block', (SELECT id FROM product_categories WHERE code = 'DAIRY'), 'PCS', true, 90, 85000, 82000, 80000, 78000, 95000, 60000, 50, 100, true),
('YOGT-001', '8851234567899', 'Greek Yogurt 500g', 'Plain Greek yogurt', (SELECT id FROM product_categories WHERE code = 'DAIRY'), 'PCS', true, 21, 35000, 33000, 32000, 30000, 40000, 25000, 100, 200, true),

-- Produce
('VEG-001', '8851234567900', 'Fresh Tomatoes', 'Ripe red tomatoes', (SELECT id FROM product_categories WHERE code = 'PRODUCE'), 'KG', false, 7, 25000, 23000, 22000, 20000, 28000, 15000, 150, 300, true),
('VEG-002', '8851234567901', 'Lettuce Head', 'Fresh iceberg lettuce', (SELECT id FROM product_categories WHERE code = 'PRODUCE'), 'PCS', false, 5, 18000, 17000, 16000, 15000, 20000, 12000, 100, 200, true),
('VEG-003', '8851234567902', 'Carrots 1kg', 'Fresh orange carrots', (SELECT id FROM product_categories WHERE code = 'PRODUCE'), 'KG', false, 14, 15000, 14000, 13000, 12000, 18000, 10000, 100, 250, true),
('FRT-001', '8851234567903', 'Banana Bunch', 'Fresh ripe bananas', (SELECT id FROM product_categories WHERE code = 'PRODUCE'), 'KG', false, 5, 20000, 19000, 18000, 17000, 22000, 14000, 150, 300, true),

-- Dry Goods
('RICE-001', '8851234567904', 'Jasmine Rice 5kg', 'Premium Thai jasmine rice', (SELECT id FROM product_categories WHERE code = 'DRY'), 'BAG', false, 365, 75000, 72000, 70000, 68000, 80000, 55000, 200, 500, true),
('NDLS-001', '8851234567905', 'Rice Noodles 500g', 'Dried rice noodles', (SELECT id FROM product_categories WHERE code = 'DRY'), 'PKT', false, 180, 15000, 14000, 13000, 12000, 18000, 9000, 300, 600, true),
('OIL-001', '8851234567906', 'Vegetable Oil 1L', 'Premium cooking oil', (SELECT id FROM product_categories WHERE code = 'DRY'), 'BTL', false, 365, 35000, 33000, 32000, 30000, 40000, 25000, 150, 300, true),

-- Beverages
('BEV-001', '8851234567907', 'Mineral Water 1.5L', 'Natural mineral water', (SELECT id FROM product_categories WHERE code = 'BEVERAGE'), 'BTL', false, 365, 8000, 7500, 7000, 6500, 10000, 5000, 500, 1000, true),
('BEV-002', '8851234567908', 'Orange Juice 1L', 'Fresh orange juice', (SELECT id FROM product_categories WHERE code = 'BEVERAGE'), 'BTL', true, 30, 45000, 43000, 42000, 40000, 50000, 32000, 100, 200, true),
('BEV-003', '8851234567909', 'Green Tea 500ml', 'Bottled green tea', (SELECT id FROM product_categories WHERE code = 'BEVERAGE'), 'BTL', false, 180, 12000, 11000, 10500, 10000, 15000, 7000, 300, 600, true)
ON CONFLICT (sku) DO NOTHING;

-- Step 4: Insert Inventory Records
-- Main Warehouse inventory
INSERT INTO inventory (product_id, warehouse_id, location_code, lot_number, production_date, expiry_date, quantity_on_hand, quantity_allocated, quantity_on_order, last_cost, average_cost)
SELECT 
    p.id,
    (SELECT id FROM warehouses WHERE warehouse_code = 'WH-MAIN'),
    'A-01-01',
    'LOT' || TO_CHAR(NOW(), 'YYYYMMDD') || '-' || LPAD(p.id::text, 3, '0'),
    CURRENT_DATE - (RANDOM() * 5)::INT,
    CURRENT_DATE + p.shelf_life_days,
    FLOOR(RANDOM() * 200 + 50)::DECIMAL,
    FLOOR(RANDOM() * 20)::DECIMAL,
    FLOOR(RANDOM() * 50)::DECIMAL,
    p.standard_cost,
    p.standard_cost * (0.95 + RANDOM() * 0.1)
FROM products p
WHERE p.category_id IN (SELECT id FROM product_categories WHERE code IN ('MEAT', 'PRODUCE', 'DRY'))
ON CONFLICT (product_id, warehouse_id, location_code, lot_number) DO NOTHING;

-- Cold Storage inventory
INSERT INTO inventory (product_id, warehouse_id, location_code, lot_number, production_date, expiry_date, quantity_on_hand, quantity_allocated, quantity_on_order, last_cost, average_cost)
SELECT 
    p.id,
    (SELECT id FROM warehouses WHERE warehouse_code = 'WH-COLD'),
    'C-01-01',
    'LOT' || TO_CHAR(NOW(), 'YYYYMMDD') || '-C' || LPAD(p.id::text, 3, '0'),
    CURRENT_DATE - (RANDOM() * 3)::INT,
    CURRENT_DATE + p.shelf_life_days,
    FLOOR(RANDOM() * 100 + 30)::DECIMAL,
    FLOOR(RANDOM() * 15)::DECIMAL,
    FLOOR(RANDOM() * 30)::DECIMAL,
    p.standard_cost,
    p.standard_cost * (0.95 + RANDOM() * 0.1)
FROM products p
WHERE p.category_id IN (SELECT id FROM product_categories WHERE code IN ('SEAFOOD', 'DAIRY', 'FROZEN'))
ON CONFLICT (product_id, warehouse_id, location_code, lot_number) DO NOTHING;

-- Dry Goods Warehouse
INSERT INTO inventory (product_id, warehouse_id, location_code, lot_number, production_date, expiry_date, quantity_on_hand, quantity_allocated, quantity_on_order, last_cost, average_cost)
SELECT 
    p.id,
    (SELECT id FROM warehouses WHERE warehouse_code = 'WH-DRY'),
    'D-01-01',
    'LOT' || TO_CHAR(NOW(), 'YYYYMMDD') || '-D' || LPAD(p.id::text, 3, '0'),
    CURRENT_DATE - (RANDOM() * 30)::INT,
    CURRENT_DATE + p.shelf_life_days,
    FLOOR(RANDOM() * 500 + 100)::DECIMAL,
    FLOOR(RANDOM() * 50)::DECIMAL,
    FLOOR(RANDOM() * 100)::DECIMAL,
    p.standard_cost,
    p.standard_cost * (0.95 + RANDOM() * 0.1)
FROM products p
WHERE p.category_id IN (SELECT id FROM product_categories WHERE code IN ('DRY', 'BEVERAGE', 'SPICES'))
ON CONFLICT (product_id, warehouse_id, location_code, lot_number) DO NOTHING;

-- Add some expiring soon items for testing
INSERT INTO inventory (product_id, warehouse_id, location_code, lot_number, production_date, expiry_date, quantity_on_hand, quantity_allocated, last_cost, average_cost)
SELECT 
    p.id,
    (SELECT id FROM warehouses WHERE warehouse_code = 'WH-COLD'),
    'C-02-01',
    'EXPIRING-' || LPAD(p.id::text, 3, '0'),
    CURRENT_DATE - 5,
    CURRENT_DATE + 3, -- Expiring in 3 days!
    FLOOR(RANDOM() * 30 + 10)::DECIMAL,
    0,
    p.standard_cost,
    p.standard_cost
FROM products p
WHERE p.category_id IN (SELECT id FROM product_categories WHERE code IN ('MEAT', 'SEAFOOD'))
LIMIT 5
ON CONFLICT (product_id, warehouse_id, location_code, lot_number) DO NOTHING;

-- Add some low stock items
INSERT INTO inventory (product_id, warehouse_id, location_code, lot_number, production_date, expiry_date, quantity_on_hand, quantity_allocated, last_cost, average_cost)
SELECT 
    p.id,
    (SELECT id FROM warehouses WHERE warehouse_code = 'WH-MAIN'),
    'A-03-01',
    'LOWSTOCK-' || LPAD(p.id::text, 3, '0'),
    CURRENT_DATE - 2,
    CURRENT_DATE + p.shelf_life_days,
    FLOOR(RANDOM() * 10 + 5)::DECIMAL, -- Low stock: 5-15 units
    FLOOR(RANDOM() * 3)::DECIMAL,
    p.standard_cost,
    p.standard_cost
FROM products p
WHERE p.is_active = true
LIMIT 5
ON CONFLICT (product_id, warehouse_id, location_code, lot_number) DO NOTHING;

-- Step 5: Verify data was inserted
SELECT 'Warehouses:' as info, COUNT(*) as count FROM warehouses;
SELECT 'Products:' as info, COUNT(*) as count FROM products;
SELECT 'Inventory Records:' as info, COUNT(*) as count FROM inventory;

-- Show inventory summary
SELECT 
    w.name as warehouse,
    COUNT(DISTINCT i.product_id) as products,
    SUM(i.quantity_on_hand) as total_qty,
    SUM(i.quantity_on_hand * i.average_cost) as total_value
FROM inventory i
JOIN warehouses w ON i.warehouse_id = w.id
GROUP BY w.name
ORDER BY w.name;
