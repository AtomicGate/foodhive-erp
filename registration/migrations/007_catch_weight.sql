-- ============================================
-- Catch Weight Tables for Food Distribution
-- Handles variable weight products (meat, cheese, produce, etc.)
-- ============================================

-- Catch Weight Entries (Header)
-- One entry per product per receiving/sales/picking reference
CREATE TABLE IF NOT EXISTS catch_weight_entries (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES products(id),
    reference_type VARCHAR(20) NOT NULL, -- 'RECEIVING', 'SALES', 'PICKING', 'ADJUSTMENT'
    reference_id INTEGER NOT NULL,
    reference_number VARCHAR(50),
    lot_number VARCHAR(50),
    expected_weight DECIMAL(12, 4) NOT NULL DEFAULT 0, -- What was ordered/expected
    actual_weight DECIMAL(12, 4) NOT NULL DEFAULT 0,   -- What was actually weighed
    weight_uom VARCHAR(5) NOT NULL DEFAULT 'KG',       -- KG, LB, GR, OZ
    piece_count INTEGER NOT NULL DEFAULT 0,
    variance DECIMAL(12, 4) GENERATED ALWAYS AS (actual_weight - expected_weight) STORED,
    variance_percent DECIMAL(8, 4) GENERATED ALWAYS AS (
        CASE WHEN expected_weight > 0 THEN ((actual_weight - expected_weight) / expected_weight) * 100 
        ELSE 0 END
    ) STORED,
    is_billed BOOLEAN DEFAULT FALSE,
    captured_by INTEGER REFERENCES employees(id),
    captured_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    notes TEXT,
    
    -- Ensure unique entry per reference + product
    CONSTRAINT uq_catch_weight_reference UNIQUE (reference_type, reference_id, product_id)
);

-- Individual Piece Weights
-- Each box may contain multiple pieces with different weights
CREATE TABLE IF NOT EXISTS catch_weight_pieces (
    id SERIAL PRIMARY KEY,
    entry_id INTEGER NOT NULL REFERENCES catch_weight_entries(id) ON DELETE CASCADE,
    piece_number INTEGER NOT NULL, -- 1, 2, 3, etc. within the entry
    weight DECIMAL(10, 4) NOT NULL,
    weight_uom VARCHAR(5) NOT NULL DEFAULT 'KG',
    barcode VARCHAR(50),          -- Individual piece barcode (if any)
    tag_number VARCHAR(50),       -- Physical tag number
    quality_grade VARCHAR(10),    -- A, B, C grade
    temperature DECIMAL(5, 2),    -- For cold chain products
    captured_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    notes TEXT
);

-- Product Catch Weight Configuration
-- How to handle catch weight for each product
CREATE TABLE IF NOT EXISTS catch_weight_config (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL UNIQUE REFERENCES products(id),
    standard_piece_weight DECIMAL(10, 4) DEFAULT 0, -- Expected weight per piece
    weight_uom VARCHAR(5) NOT NULL DEFAULT 'KG',
    min_weight DECIMAL(10, 4) DEFAULT 0,            -- Minimum acceptable weight
    max_weight DECIMAL(10, 4) DEFAULT 0,            -- Maximum acceptable weight
    variance_tolerance DECIMAL(5, 2) DEFAULT 5.00,  -- Allowed variance % (e.g., 5%)
    require_piece_weights BOOLEAN DEFAULT FALSE,    -- Must capture each piece?
    pricing_method VARCHAR(20) DEFAULT 'ACTUAL_WEIGHT', -- ACTUAL_WEIGHT, STANDARD_WEIGHT, CATCH_UP
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_cw_entries_product ON catch_weight_entries(product_id);
CREATE INDEX idx_cw_entries_reference ON catch_weight_entries(reference_type, reference_id);
CREATE INDEX idx_cw_entries_lot ON catch_weight_entries(lot_number) WHERE lot_number IS NOT NULL;
CREATE INDEX idx_cw_entries_captured_at ON catch_weight_entries(captured_at);
CREATE INDEX idx_cw_pieces_entry ON catch_weight_pieces(entry_id);

-- Example: How data flows
-- ============================================
-- 
-- RECEIVING FLOW:
-- 1. PO arrives for "10 boxes of Chicken Breast" (expected: 17kg per box = 170kg total)
-- 2. During receiving, operator weighs each box:
--    - Box 1: 17.55kg (3 pieces: 5.8kg, 5.9kg, 5.85kg)
--    - Box 2: 16.80kg (3 pieces: 5.5kg, 5.7kg, 5.6kg)
--    - etc.
-- 3. System records:
--    - catch_weight_entries: expected=170kg, actual=173.5kg, variance=+3.5kg (+2.06%)
--    - catch_weight_pieces: individual piece weights for each box
-- 
-- SALES FLOW:
-- 1. Customer orders "5 pieces of Chicken Breast"
-- 2. During picking, each piece is weighed:
--    - Piece 1: 5.6kg
--    - Piece 2: 5.8kg
--    - etc.
-- 3. Invoice is generated based on ACTUAL weight, not standard weight
-- 
-- BILLING:
-- - Standard weight: 5 pieces × 5.5kg (standard) = 27.5kg @ $10/kg = $275
-- - Actual weight: 5 pieces totaling 28.3kg @ $10/kg = $283
-- - Adjustment: +$8.00
-- ============================================

COMMENT ON TABLE catch_weight_entries IS 'Tracks actual weights for variable weight products';
COMMENT ON TABLE catch_weight_pieces IS 'Individual piece weights within a catch weight entry';
COMMENT ON TABLE catch_weight_config IS 'Per-product configuration for catch weight handling';
COMMENT ON COLUMN catch_weight_entries.reference_type IS 'Type of transaction: RECEIVING, SALES, PICKING, ADJUSTMENT';
COMMENT ON COLUMN catch_weight_entries.variance IS 'Difference between actual and expected weight';
COMMENT ON COLUMN catch_weight_config.pricing_method IS 'ACTUAL_WEIGHT=bill by actual, STANDARD_WEIGHT=bill by standard, CATCH_UP=reconcile periodically';
