# FoodHive ERP - Enhanced Implementation Plan
## Complete System Implementation with Professional Suggestions

This document provides a comprehensive implementation plan that includes:
- ✅ **Required Fields** - From DocText.md requirements
- 💡 **SUGGESTED** - Professional recommendations for a complete ERP system (marked with comments)

---

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [Phase 1: Database Schema Changes](#phase-1-database-schema-changes)
3. [Phase 2: Backend Models (Go)](#phase-2-backend-models-go)
4. [Phase 3: Backend Services & Routes](#phase-3-backend-services--routes)
5. [Phase 4: Frontend Services](#phase-4-frontend-services)
6. [Phase 5: Frontend Components](#phase-5-frontend-components)
7. [Phase 6: System Enhancements (Suggested)](#phase-6-system-enhancements-suggested)
8. [Implementation Priority & Timeline](#implementation-priority--timeline)

---

## Legend

| Symbol | Meaning |
|--------|---------|
| ✅ | Required from DocText.md |
| 💡 | Suggested enhancement (with reason) |
| 🔄 | Existing field being modified |

---

## Phase 1: Database Schema Changes

### File: `sql/006_requirements_gap_updates.sql`

```sql
-- =====================================================
-- PHASE 1: Database Schema Updates for DocText.md Requirements
-- Version: 2.0 ENHANCED
-- =====================================================

-- =====================================================
-- 1.1 ENUM TYPES FOR NEW MODULES
-- =====================================================

-- Price Type Enum ✅ Required
CREATE TYPE price_type AS ENUM ('PRICE_1', 'PRICE_2', 'PRICE_3', 'WHOLESALE', 'RETAIL');

-- Customer Level Enum ✅ Required
CREATE TYPE customer_level AS ENUM ('GENERAL', 'SILVER', 'GOLD', 'PLATINUM', 'VIP');

-- Currency Enum ✅ Required (extend for multi-currency)
CREATE TYPE currency_code AS ENUM ('LAK', 'THB', 'USD', 'KIP');

-- Purchase Requisition Status ✅ Required
CREATE TYPE pr_status AS ENUM ('DRAFT', 'SUBMITTED', 'CHECKED', 'AUTHORIZED', 'CONVERTED', 'CANCELLED');

-- Advance Request Status ✅ Required
CREATE TYPE advance_request_status AS ENUM ('DRAFT', 'SUBMITTED', 'APPROVED', 'REJECTED', 'CLEARED');

-- Advance Voucher Status ✅ Required
CREATE TYPE advance_voucher_status AS ENUM ('DRAFT', 'SUBMITTED', 'APPROVED', 'REJECTED', 'COMPLETED');

-- Document Type for Advance Voucher ✅ Required
CREATE TYPE voucher_document_type AS ENUM ('RECEIPT', 'INVOICE', 'BANK_TRANSFER', 'PO_PR', 'TRANSPORTATION');

-- 💡 SUGGESTED: Priority levels for requisitions and requests
-- Reason: Helps prioritize urgent purchases and advance requests
CREATE TYPE priority_level AS ENUM ('LOW', 'NORMAL', 'HIGH', 'URGENT', 'CRITICAL');

-- 💡 SUGGESTED: Contact preference for vendors/customers
-- Reason: Improves communication efficiency
CREATE TYPE contact_preference AS ENUM ('PHONE', 'EMAIL', 'FAX', 'MOBILE', 'WHATSAPP', 'LINE');

-- 💡 SUGGESTED: Payment method enum
-- Reason: Standardizes payment tracking across modules
CREATE TYPE payment_method AS ENUM ('CASH', 'BANK_TRANSFER', 'CHECK', 'CREDIT_CARD', 'CREDIT_TERM', 'COD');

-- 💡 SUGGESTED: Document status for approvals
-- Reason: Provides clearer workflow states
CREATE TYPE approval_action AS ENUM ('PENDING', 'APPROVED', 'REJECTED', 'RETURNED', 'ESCALATED');

-- =====================================================
-- 1.2 CUSTOMER TABLE UPDATES
-- =====================================================

-- ✅ REQUIRED FIELDS FROM DocText.md
ALTER TABLE customers ADD COLUMN IF NOT EXISTS default_contact_name TEXT;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS full_address TEXT;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS delivery_name TEXT;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS tel VARCHAR(50);
ALTER TABLE customers ADD COLUMN IF NOT EXISTS fax VARCHAR(50);
ALTER TABLE customers ADD COLUMN IF NOT EXISTS mobile VARCHAR(50);
ALTER TABLE customers ADD COLUMN IF NOT EXISTS tax_code VARCHAR(50);
ALTER TABLE customers ADD COLUMN IF NOT EXISTS birthday DATE;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS web_email VARCHAR(255);
ALTER TABLE customers ADD COLUMN IF NOT EXISTS email VARCHAR(255);
ALTER TABLE customers ADD COLUMN IF NOT EXISTS price_type price_type DEFAULT 'PRICE_1';
ALTER TABLE customers ADD COLUMN IF NOT EXISTS customer_level customer_level DEFAULT 'GENERAL';
ALTER TABLE customers ADD COLUMN IF NOT EXISTS bill_discount DECIMAL(5,2) DEFAULT 0.00;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS member_expire DATE;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS credit_days INT DEFAULT 0;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS passport_no VARCHAR(50);
ALTER TABLE customers ADD COLUMN IF NOT EXISTS country_code VARCHAR(10);
ALTER TABLE customers ADD COLUMN IF NOT EXISTS gender gender;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS use_promotion BOOLEAN DEFAULT true;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS use_pmt_bill BOOLEAN DEFAULT false;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS collect_point BOOLEAN DEFAULT false;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS payee TEXT;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS rd_branch TEXT;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS barcode VARCHAR(100);
ALTER TABLE customers ADD COLUMN IF NOT EXISTS notes TEXT;

-- 💡 SUGGESTED: Customer analytics and engagement tracking
-- Reason: Helps with sales analysis, marketing, and customer relationship management

-- 💡 Last order date tracking
-- Reason: Quick reference for customer activity without querying orders table
ALTER TABLE customers ADD COLUMN IF NOT EXISTS last_order_date DATE;

-- 💡 Total order count
-- Reason: Identify VIP customers and order frequency patterns
ALTER TABLE customers ADD COLUMN IF NOT EXISTS total_orders INT DEFAULT 0;

-- 💡 Total lifetime spending
-- Reason: Customer value analysis for loyalty programs and discounts
ALTER TABLE customers ADD COLUMN IF NOT EXISTS total_spent DECIMAL(15,2) DEFAULT 0.00;

-- 💡 Loyalty points balance
-- Reason: Supports loyalty/reward programs (collect_point checkbox implies this)
ALTER TABLE customers ADD COLUMN IF NOT EXISTS loyalty_points DECIMAL(12,2) DEFAULT 0.00;

-- 💡 Credit risk score
-- Reason: Helps manage credit risk and automate credit decisions
ALTER TABLE customers ADD COLUMN IF NOT EXISTS credit_score INT DEFAULT 100; -- 0-100 scale

-- 💡 Preferred delivery time window
-- Reason: Improves delivery scheduling and customer satisfaction
ALTER TABLE customers ADD COLUMN IF NOT EXISTS preferred_delivery_time VARCHAR(50);

-- 💡 Special delivery instructions
-- Reason: Reduces delivery errors and improves service
ALTER TABLE customers ADD COLUMN IF NOT EXISTS delivery_instructions TEXT;

-- 💡 Referral tracking
-- Reason: Supports referral programs and tracks customer acquisition
ALTER TABLE customers ADD COLUMN IF NOT EXISTS referred_by_id INT REFERENCES customers(id);
ALTER TABLE customers ADD COLUMN IF NOT EXISTS referral_code VARCHAR(20);

-- 💡 Customer source/channel
-- Reason: Track where customers come from for marketing analysis
ALTER TABLE customers ADD COLUMN IF NOT EXISTS acquisition_source VARCHAR(50);

-- 💡 Preferred contact method
-- Reason: Improves communication efficiency
ALTER TABLE customers ADD COLUMN IF NOT EXISTS preferred_contact contact_preference DEFAULT 'PHONE';

-- 💡 Secondary contact person
-- Reason: Backup contact for critical communications
ALTER TABLE customers ADD COLUMN IF NOT EXISTS secondary_contact_name TEXT;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS secondary_contact_phone VARCHAR(50);

-- 💡 Social media / messaging IDs
-- Reason: Modern communication channels widely used in Southeast Asia
ALTER TABLE customers ADD COLUMN IF NOT EXISTS line_id VARCHAR(100);
ALTER TABLE customers ADD COLUMN IF NOT EXISTS whatsapp VARCHAR(50);

-- 💡 GPS coordinates for delivery
-- Reason: Precise delivery location for mapping and route optimization
ALTER TABLE customers ADD COLUMN IF NOT EXISTS latitude DECIMAL(10,8);
ALTER TABLE customers ADD COLUMN IF NOT EXISTS longitude DECIMAL(11,8);

-- 💡 Customer classification for B2B
-- Reason: Distinguish between B2B and B2C customers for different pricing/terms
ALTER TABLE customers ADD COLUMN IF NOT EXISTS customer_type VARCHAR(20) DEFAULT 'RETAIL'; -- RETAIL, WHOLESALE, DISTRIBUTOR, CHAIN

-- 💡 Company registration info for B2B
-- Reason: Required for invoicing and tax compliance in B2B transactions
ALTER TABLE customers ADD COLUMN IF NOT EXISTS company_registration_no VARCHAR(50);
ALTER TABLE customers ADD COLUMN IF NOT EXISTS vat_registration_no VARCHAR(50);

-- 💡 Account manager assignment (different from sales_rep)
-- Reason: Separates order-taking from relationship management
ALTER TABLE customers ADD COLUMN IF NOT EXISTS account_manager_id INT REFERENCES employees(id);

-- 💡 Customer group for pricing
-- Reason: Allows group-based pricing and promotions
ALTER TABLE customers ADD COLUMN IF NOT EXISTS customer_group_id INT;

-- 💡 Blocked status with reason
-- Reason: Prevents orders for problematic customers with audit trail
ALTER TABLE customers ADD COLUMN IF NOT EXISTS is_blocked BOOLEAN DEFAULT false;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS block_reason TEXT;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS blocked_date DATE;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS blocked_by INT REFERENCES employees(id);

-- =====================================================
-- 1.3 VENDOR TABLE UPDATES
-- =====================================================

-- ✅ REQUIRED FIELDS FROM DocText.md
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS contact_name TEXT;
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS rd_branch TEXT;
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS picture_name TEXT;
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS mobile VARCHAR(50);
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS fax VARCHAR(50);
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS tax_code VARCHAR(50);
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS credit_days INT DEFAULT 0;
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS discount_percent DECIMAL(5,2) DEFAULT 0.00;
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS po_credit_money DECIMAL(12,2) DEFAULT 0.00;

-- 💡 SUGGESTED: Vendor performance and relationship management
-- Reason: Enables vendor evaluation, quality control, and strategic sourcing

-- 💡 Vendor rating/score
-- Reason: Track vendor performance for purchasing decisions
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS vendor_rating DECIMAL(3,2) DEFAULT 5.00; -- 0-5 scale

-- 💡 Quality score
-- Reason: Track product quality from this vendor
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS quality_score DECIMAL(5,2) DEFAULT 100.00; -- percentage

-- 💡 On-time delivery rate
-- Reason: Critical for supply chain reliability assessment
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS on_time_delivery_rate DECIMAL(5,2) DEFAULT 100.00;

-- 💡 Last PO date
-- Reason: Quick reference for vendor activity
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS last_po_date DATE;

-- 💡 Total PO value lifetime
-- Reason: Vendor relationship value assessment
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS total_po_value DECIMAL(15,2) DEFAULT 0.00;

-- 💡 Total PO count
-- Reason: Order frequency tracking
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS total_po_count INT DEFAULT 0;

-- 💡 Website
-- Reason: Quick access to vendor information
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS website VARCHAR(255);

-- 💡 Bank account info for payments
-- Reason: Required for AP payments
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS bank_name VARCHAR(100);
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS bank_account_number VARCHAR(50);
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS bank_account_name VARCHAR(100);
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS bank_branch VARCHAR(100);
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS swift_code VARCHAR(20);

-- 💡 Preferred payment method
-- Reason: Streamlines AP processing
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS preferred_payment payment_method DEFAULT 'BANK_TRANSFER';

-- 💡 Certifications and compliance
-- Reason: Food industry requires tracking of vendor certifications (HACCP, ISO, etc.)
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS certifications TEXT; -- JSON or comma-separated
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS certification_expiry DATE;

-- 💡 Notes/comments
-- Reason: Store important vendor-specific information
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS notes TEXT;

-- 💡 Preferred contact method
-- Reason: Improves ordering efficiency
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS preferred_contact contact_preference DEFAULT 'EMAIL';

-- 💡 Secondary contact
-- Reason: Backup contact for urgent orders
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS secondary_contact_name TEXT;
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS secondary_contact_phone VARCHAR(50);

-- 💡 Vendor type classification
-- Reason: Different vendor types may have different terms and handling
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS vendor_type VARCHAR(30) DEFAULT 'SUPPLIER'; -- SUPPLIER, MANUFACTURER, DISTRIBUTOR, IMPORTER

-- 💡 Return policy info
-- Reason: Important for handling defective/expired goods
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS return_policy TEXT;
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS return_window_days INT DEFAULT 30;

-- 💡 Blocked status
-- Reason: Prevent POs to problematic vendors
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS is_blocked BOOLEAN DEFAULT false;
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS block_reason TEXT;

-- =====================================================
-- 1.4 PURCHASE ORDER TABLE UPDATES
-- =====================================================

-- ✅ REQUIRED FIELDS FROM DocText.md
ALTER TABLE purchase_orders ADD COLUMN IF NOT EXISTS contact_person TEXT;
ALTER TABLE purchase_orders ADD COLUMN IF NOT EXISTS customer_name TEXT;
ALTER TABLE purchase_orders ADD COLUMN IF NOT EXISTS address TEXT;
ALTER TABLE purchase_orders ADD COLUMN IF NOT EXISTS tel VARCHAR(50);
ALTER TABLE purchase_orders ADD COLUMN IF NOT EXISTS email VARCHAR(255);
ALTER TABLE purchase_orders ADD COLUMN IF NOT EXISTS document_ref VARCHAR(100);
ALTER TABLE purchase_orders ADD COLUMN IF NOT EXISTS terms_of_payment TEXT;
ALTER TABLE purchase_orders ADD COLUMN IF NOT EXISTS currency currency_code DEFAULT 'LAK';
ALTER TABLE purchase_orders ADD COLUMN IF NOT EXISTS trade_discount DECIMAL(12,2) DEFAULT 0.00;
ALTER TABLE purchase_orders ADD COLUMN IF NOT EXISTS amount_after_discount DECIMAL(12,2) DEFAULT 0.00;
ALTER TABLE purchase_orders ADD COLUMN IF NOT EXISTS vat_rate DECIMAL(5,2) DEFAULT 7.00;
ALTER TABLE purchase_orders ADD COLUMN IF NOT EXISTS prepared_by INT REFERENCES employees(id);
ALTER TABLE purchase_orders ADD COLUMN IF NOT EXISTS authorized_by INT REFERENCES employees(id);
ALTER TABLE purchase_orders ADD COLUMN IF NOT EXISTS checked_by INT REFERENCES employees(id);
ALTER TABLE purchase_orders ADD COLUMN IF NOT EXISTS prepared_date DATE;
ALTER TABLE purchase_orders ADD COLUMN IF NOT EXISTS authorized_date DATE;
ALTER TABLE purchase_orders ADD COLUMN IF NOT EXISTS checked_date DATE;

-- ✅ REQUIRED: Line item discount
ALTER TABLE purchase_order_lines ADD COLUMN IF NOT EXISTS discount_percent DECIMAL(5,2) DEFAULT 0.00;
ALTER TABLE purchase_order_lines ADD COLUMN IF NOT EXISTS discount_amount DECIMAL(12,2) DEFAULT 0.00;

-- 💡 SUGGESTED: Enhanced PO tracking and logistics
-- Reason: Improves receiving process and supply chain visibility

-- 💡 Shipping/logistics info
-- Reason: Track delivery method and costs
ALTER TABLE purchase_orders ADD COLUMN IF NOT EXISTS shipping_method VARCHAR(50);
ALTER TABLE purchase_orders ADD COLUMN IF NOT EXISTS shipping_carrier VARCHAR(100);
ALTER TABLE purchase_orders ADD COLUMN IF NOT EXISTS tracking_number VARCHAR(100);
ALTER TABLE purchase_orders ADD COLUMN IF NOT EXISTS insurance_amount DECIMAL(12,2) DEFAULT 0.00;

-- 💡 Delivery instructions
-- Reason: Special handling requirements for cold storage items
ALTER TABLE purchase_orders ADD COLUMN IF NOT EXISTS delivery_instructions TEXT;

-- 💡 Quality notes
-- Reason: Special quality requirements for this order
ALTER TABLE purchase_orders ADD COLUMN IF NOT EXISTS quality_requirements TEXT;

-- 💡 Exchange rate at time of PO
-- Reason: Critical for multi-currency POs
ALTER TABLE purchase_orders ADD COLUMN IF NOT EXISTS exchange_rate DECIMAL(12,6) DEFAULT 1.00;
ALTER TABLE purchase_orders ADD COLUMN IF NOT EXISTS base_currency VARCHAR(3) DEFAULT 'LAK';
ALTER TABLE purchase_orders ADD COLUMN IF NOT EXISTS total_in_base_currency DECIMAL(12,2);

-- 💡 Priority level
-- Reason: Helps prioritize receiving and processing
ALTER TABLE purchase_orders ADD COLUMN IF NOT EXISTS priority priority_level DEFAULT 'NORMAL';

-- 💡 Linked documents
-- Reason: Connect PO to related documents (requisitions, quotes, contracts)
ALTER TABLE purchase_orders ADD COLUMN IF NOT EXISTS requisition_id INT; -- FK added later
ALTER TABLE purchase_orders ADD COLUMN IF NOT EXISTS quotation_ref VARCHAR(50);
ALTER TABLE purchase_orders ADD COLUMN IF NOT EXISTS contract_ref VARCHAR(50);

-- 💡 Payment tracking
-- Reason: Track payment status without going to AP module
ALTER TABLE purchase_orders ADD COLUMN IF NOT EXISTS payment_status VARCHAR(20) DEFAULT 'PENDING'; -- PENDING, PARTIAL, PAID
ALTER TABLE purchase_orders ADD COLUMN IF NOT EXISTS amount_paid DECIMAL(12,2) DEFAULT 0.00;

-- 💡 Cancellation info
-- Reason: Track why POs were cancelled for analysis
ALTER TABLE purchase_orders ADD COLUMN IF NOT EXISTS cancelled_by INT REFERENCES employees(id);
ALTER TABLE purchase_orders ADD COLUMN IF NOT EXISTS cancelled_date DATE;
ALTER TABLE purchase_orders ADD COLUMN IF NOT EXISTS cancellation_reason TEXT;

-- 💡 SUGGESTED: PO Line enhancements
-- Reason: Better tracking of line-level details

-- 💡 Tax per line (some items may have different tax rates)
ALTER TABLE purchase_order_lines ADD COLUMN IF NOT EXISTS tax_rate DECIMAL(5,2) DEFAULT 7.00;
ALTER TABLE purchase_order_lines ADD COLUMN IF NOT EXISTS tax_amount DECIMAL(12,2) DEFAULT 0.00;

-- 💡 Requested delivery date per line
-- Reason: Different items may need different delivery dates
ALTER TABLE purchase_order_lines ADD COLUMN IF NOT EXISTS requested_date DATE;

-- 💡 Notes per line
-- Reason: Line-specific instructions or notes
ALTER TABLE purchase_order_lines ADD COLUMN IF NOT EXISTS notes TEXT;

-- =====================================================
-- 1.5 PRODUCT TABLE UPDATES
-- =====================================================

-- ✅ REQUIRED FIELDS FROM DocText.md
ALTER TABLE products ADD COLUMN IF NOT EXISTS brand TEXT;
ALTER TABLE products ADD COLUMN IF NOT EXISTS product_group TEXT;
ALTER TABLE products ADD COLUMN IF NOT EXISTS price_1 DECIMAL(12,4) DEFAULT 0.00;
ALTER TABLE products ADD COLUMN IF NOT EXISTS price_2 DECIMAL(12,4) DEFAULT 0.00;
ALTER TABLE products ADD COLUMN IF NOT EXISTS vat_percent DECIMAL(5,2) DEFAULT 7.00;

-- 💡 SUGGESTED: Enhanced product management
-- Reason: Comprehensive product data for food distribution business

-- 💡 Additional price levels (common in wholesale)
-- Reason: Different prices for different customer levels
ALTER TABLE products ADD COLUMN IF NOT EXISTS price_3 DECIMAL(12,4) DEFAULT 0.00;
ALTER TABLE products ADD COLUMN IF NOT EXISTS wholesale_price DECIMAL(12,4) DEFAULT 0.00;
ALTER TABLE products ADD COLUMN IF NOT EXISTS retail_price DECIMAL(12,4) DEFAULT 0.00;

-- 💡 Cost tracking
-- Reason: Margin analysis and profitability
ALTER TABLE products ADD COLUMN IF NOT EXISTS standard_cost DECIMAL(12,4) DEFAULT 0.00;
ALTER TABLE products ADD COLUMN IF NOT EXISTS last_purchase_cost DECIMAL(12,4) DEFAULT 0.00;
ALTER TABLE products ADD COLUMN IF NOT EXISTS average_cost DECIMAL(12,4) DEFAULT 0.00;

-- 💡 Inventory planning
-- Reason: Automatic reorder suggestions
ALTER TABLE products ADD COLUMN IF NOT EXISTS reorder_point DECIMAL(12,3) DEFAULT 0.00;
ALTER TABLE products ADD COLUMN IF NOT EXISTS reorder_quantity DECIMAL(12,3) DEFAULT 0.00;
ALTER TABLE products ADD COLUMN IF NOT EXISTS safety_stock DECIMAL(12,3) DEFAULT 0.00;
ALTER TABLE products ADD COLUMN IF NOT EXISTS lead_time_days INT DEFAULT 7;

-- 💡 Sales analytics
-- Reason: Demand planning and inventory optimization
ALTER TABLE products ADD COLUMN IF NOT EXISTS avg_daily_sales DECIMAL(12,3) DEFAULT 0.00;
ALTER TABLE products ADD COLUMN IF NOT EXISTS last_sale_date DATE;
ALTER TABLE products ADD COLUMN IF NOT EXISTS total_quantity_sold DECIMAL(15,3) DEFAULT 0.00;

-- 💡 Seasonal product flag
-- Reason: Important for demand planning
ALTER TABLE products ADD COLUMN IF NOT EXISTS is_seasonal BOOLEAN DEFAULT false;
ALTER TABLE products ADD COLUMN IF NOT EXISTS season_start_month INT; -- 1-12
ALTER TABLE products ADD COLUMN IF NOT EXISTS season_end_month INT; -- 1-12

-- 💡 Allergen information (critical for food)
-- Reason: Food safety compliance
ALTER TABLE products ADD COLUMN IF NOT EXISTS allergen_info TEXT; -- JSON or comma-separated

-- 💡 Nutritional info
-- Reason: Required for food products in many jurisdictions
ALTER TABLE products ADD COLUMN IF NOT EXISTS nutritional_info TEXT; -- JSON

-- 💡 Storage requirements (critical for cold storage)
-- Reason: Ensure proper storage conditions
ALTER TABLE products ADD COLUMN IF NOT EXISTS storage_temp_min DECIMAL(5,2);
ALTER TABLE products ADD COLUMN IF NOT EXISTS storage_temp_max DECIMAL(5,2);
ALTER TABLE products ADD COLUMN IF NOT EXISTS storage_instructions TEXT;

-- 💡 Dimensions and weight for logistics
-- Reason: Shipping and storage planning
ALTER TABLE products ADD COLUMN IF NOT EXISTS weight_kg DECIMAL(10,4);
ALTER TABLE products ADD COLUMN IF NOT EXISTS length_cm DECIMAL(10,2);
ALTER TABLE products ADD COLUMN IF NOT EXISTS width_cm DECIMAL(10,2);
ALTER TABLE products ADD COLUMN IF NOT EXISTS height_cm DECIMAL(10,2);
ALTER TABLE products ADD COLUMN IF NOT EXISTS volume_cm3 DECIMAL(12,2);

-- 💡 Image management
-- Reason: Product display in modern interfaces
ALTER TABLE products ADD COLUMN IF NOT EXISTS image_url TEXT;
ALTER TABLE products ADD COLUMN IF NOT EXISTS thumbnail_url TEXT;

-- 💡 Supplier relationship
-- Reason: Quick reference to primary supplier
ALTER TABLE products ADD COLUMN IF NOT EXISTS primary_vendor_id INT REFERENCES vendors(id);

-- 💡 Product notes
-- Reason: Important product-specific information
ALTER TABLE products ADD COLUMN IF NOT EXISTS internal_notes TEXT;

-- Create product_groups table ✅
CREATE TABLE IF NOT EXISTS product_groups (
    id SERIAL PRIMARY KEY,
    code VARCHAR(20) UNIQUE NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    parent_id INT REFERENCES product_groups(id),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create brands table ✅
CREATE TABLE IF NOT EXISTS brands (
    id SERIAL PRIMARY KEY,
    code VARCHAR(20) UNIQUE NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    logo_url TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    -- 💡 SUGGESTED: Additional brand info
    website VARCHAR(255),
    country_of_origin VARCHAR(100),
    manufacturer TEXT
);

-- Add foreign keys to products
ALTER TABLE products ADD COLUMN IF NOT EXISTS brand_id INT REFERENCES brands(id);
ALTER TABLE products ADD COLUMN IF NOT EXISTS product_group_id INT REFERENCES product_groups(id);

-- =====================================================
-- 1.6 NEW MODULE: PURCHASE REQUISITION ✅
-- =====================================================

CREATE TABLE IF NOT EXISTS purchase_requisitions (
    id SERIAL PRIMARY KEY,
    requisition_number VARCHAR(20) UNIQUE NOT NULL,
    document_number VARCHAR(50),
    
    -- Requester Information ✅
    requester_id INT REFERENCES employees(id),
    department_id INT REFERENCES departments(id),
    reason TEXT,
    
    -- Supplier Information ✅
    supplier_id INT REFERENCES vendors(id),
    
    -- Dates ✅
    document_date DATE DEFAULT CURRENT_DATE,
    request_date DATE DEFAULT CURRENT_DATE,
    
    -- Amounts ✅
    total_amount DECIMAL(12,2) DEFAULT 0.00,
    
    -- Notes & Remarks ✅
    remark TEXT,
    
    -- Approval Chain ✅
    checked_by INT REFERENCES employees(id),
    checked_date DATE,
    authorized_by INT REFERENCES employees(id),
    authorized_date DATE,
    purchasing_employee_id INT REFERENCES employees(id),
    purchasing_date DATE,
    
    -- Linked PO (if converted) ✅
    converted_po_id INT REFERENCES purchase_orders(id),
    
    -- Status & Audit ✅
    status pr_status DEFAULT 'DRAFT',
    created_by INT REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    -- 💡 SUGGESTED: Enhanced PR tracking
    
    -- 💡 Priority for urgent purchases
    priority priority_level DEFAULT 'NORMAL',
    
    -- 💡 Required by date
    required_date DATE,
    
    -- 💡 Budget/cost center tracking
    budget_code VARCHAR(50),
    cost_center VARCHAR(50),
    project_id INT,
    
    -- 💡 Alternative suppliers considered
    alternative_suppliers TEXT, -- JSON array of vendor IDs
    
    -- 💡 Justification for single source
    single_source_justification TEXT,
    
    -- 💡 Rejection tracking
    rejected_by INT REFERENCES employees(id),
    rejected_date DATE,
    rejection_reason TEXT,
    
    -- 💡 Revision tracking
    revision_number INT DEFAULT 1,
    previous_version_id INT REFERENCES purchase_requisitions(id)
);

CREATE TABLE IF NOT EXISTS purchase_requisition_lines (
    id SERIAL PRIMARY KEY,
    requisition_id INT REFERENCES purchase_requisitions(id) ON DELETE CASCADE,
    line_number INT NOT NULL,
    
    -- Product Information ✅
    product_id INT REFERENCES products(id),
    product_code VARCHAR(50),
    description TEXT,
    
    -- Stock & Quantity ✅
    stock_balance DECIMAL(12,3) DEFAULT 0.00,
    quantity_requested DECIMAL(10,3) NOT NULL,
    
    -- Pricing ✅
    unit_of_measure VARCHAR(10),
    unit_price DECIMAL(10,4) DEFAULT 0.00,
    amount DECIMAL(12,2) DEFAULT 0.00,
    
    -- 💡 SUGGESTED: Additional line tracking
    
    -- 💡 Quantity approved (may differ from requested)
    quantity_approved DECIMAL(10,3),
    
    -- 💡 Preferred vendor for this item
    preferred_vendor_id INT REFERENCES vendors(id),
    
    -- 💡 Notes per line
    notes TEXT,
    
    UNIQUE(requisition_id, line_number)
);

CREATE INDEX idx_purchase_requisitions_status ON purchase_requisitions(status);
CREATE INDEX idx_purchase_requisitions_requester ON purchase_requisitions(requester_id);
CREATE INDEX idx_purchase_requisitions_date ON purchase_requisitions(document_date);
-- 💡 Index for priority
CREATE INDEX idx_purchase_requisitions_priority ON purchase_requisitions(priority);

-- =====================================================
-- 1.7 NEW MODULE: ADVANCE REQUEST ✅
-- =====================================================

CREATE TABLE IF NOT EXISTS advance_requests (
    id SERIAL PRIMARY KEY,
    request_number VARCHAR(20) UNIQUE NOT NULL,
    
    -- Employee Information ✅
    employee_id INT REFERENCES employees(id) NOT NULL,
    employee_name TEXT, -- Denormalized for display
    position TEXT,
    
    -- Reference ✅
    po_number VARCHAR(50),
    
    -- Dates ✅
    request_date DATE DEFAULT CURRENT_DATE,
    
    -- Currency & Amount ✅
    currency currency_code DEFAULT 'LAK',
    total_amount DECIMAL(12,2) DEFAULT 0.00,
    exchange_rate DECIMAL(12,6) DEFAULT 1.00,
    
    -- Main Description ✅
    description TEXT,
    
    -- Approval Chain (3 signatures) ✅
    approved_by_1 INT REFERENCES employees(id),
    approved_date_1 DATE,
    approved_by_2 INT REFERENCES employees(id),
    approved_date_2 DATE,
    approved_by_3 INT REFERENCES employees(id),
    approved_date_3 DATE,
    
    -- Status & Audit ✅
    status advance_request_status DEFAULT 'DRAFT',
    created_by INT REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    -- 💡 SUGGESTED: Enhanced advance request tracking
    
    -- 💡 Purpose categorization
    purpose_category VARCHAR(50), -- TRAVEL, PURCHASE, EVENT, MAINTENANCE, OTHER
    
    -- 💡 Project/cost center linkage
    project_id INT,
    cost_center VARCHAR(50),
    budget_code VARCHAR(50),
    
    -- 💡 Bank transfer details for payment
    bank_account_number VARCHAR(50),
    bank_name VARCHAR(100),
    
    -- 💡 Expected settlement date
    expected_settlement_date DATE,
    
    -- 💡 Linked voucher (when cleared)
    voucher_id INT, -- FK added after voucher table created
    
    -- 💡 Rejection tracking
    rejected_by INT REFERENCES employees(id),
    rejected_date DATE,
    rejection_reason TEXT,
    
    -- 💡 Base currency equivalent
    total_in_base_currency DECIMAL(12,2)
);

CREATE TABLE IF NOT EXISTS advance_request_lines (
    id SERIAL PRIMARY KEY,
    request_id INT REFERENCES advance_requests(id) ON DELETE CASCADE,
    line_number INT NOT NULL,
    
    -- Hierarchy support (for sub-items like A:, B:) ✅
    parent_line_id INT REFERENCES advance_request_lines(id),
    line_type VARCHAR(20) DEFAULT 'MAIN', -- MAIN, SUB_A, SUB_B, etc.
    
    -- Content ✅
    description TEXT NOT NULL,
    amount DECIMAL(12,2) DEFAULT 0.00,
    
    -- 💡 SUGGESTED: Line categorization
    expense_category VARCHAR(50), -- GOODS, SHIPPING, SERVICE, OTHER
    
    UNIQUE(request_id, line_number)
);

CREATE INDEX idx_advance_requests_status ON advance_requests(status);
CREATE INDEX idx_advance_requests_employee ON advance_requests(employee_id);
CREATE INDEX idx_advance_requests_date ON advance_requests(request_date);

-- =====================================================
-- 1.8 NEW MODULE: ADVANCE VOUCHER ✅
-- =====================================================

CREATE TABLE IF NOT EXISTS advance_vouchers (
    id SERIAL PRIMARY KEY,
    voucher_number VARCHAR(20) UNIQUE NOT NULL,
    
    -- Link to Advance Request ✅
    advance_request_id INT REFERENCES advance_requests(id),
    
    -- Employee Information ✅
    employee_id INT REFERENCES employees(id) NOT NULL,
    employee_name TEXT,
    position TEXT,
    
    -- Reference ✅
    po_number VARCHAR(50),
    
    -- Dates ✅
    voucher_date DATE DEFAULT CURRENT_DATE,
    
    -- Currency ✅
    currency currency_code DEFAULT 'LAK',
    
    -- Amounts ✅
    advance_amount DECIMAL(12,2) DEFAULT 0.00,
    expenditure_amount DECIMAL(12,2) DEFAULT 0.00,
    balance_amount DECIMAL(12,2) GENERATED ALWAYS AS (advance_amount - expenditure_amount) STORED,
    
    -- Approval ✅
    accountant_id INT REFERENCES employees(id),
    accountant_date DATE,
    returned_by_id INT REFERENCES employees(id),
    returned_date DATE,
    
    -- Status & Audit ✅
    status advance_voucher_status DEFAULT 'DRAFT',
    created_by INT REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    -- 💡 SUGGESTED: Enhanced voucher tracking
    
    -- 💡 Exchange rate at settlement
    exchange_rate DECIMAL(12,6) DEFAULT 1.00,
    
    -- 💡 Balance handling
    balance_returned BOOLEAN DEFAULT false,
    balance_return_date DATE,
    balance_received_by INT REFERENCES employees(id),
    
    -- 💡 Additional funds requested (if expenditure > advance)
    additional_amount_requested DECIMAL(12,2) DEFAULT 0.00,
    additional_approved_by INT REFERENCES employees(id),
    additional_approved_date DATE,
    
    -- 💡 Receipt count for verification
    total_receipt_count INT DEFAULT 0,
    
    -- 💡 Verification status
    verified_by INT REFERENCES employees(id),
    verified_date DATE,
    verification_notes TEXT,
    
    -- 💡 Rejection tracking
    rejected_by INT REFERENCES employees(id),
    rejected_date DATE,
    rejection_reason TEXT
);

-- Add FK from advance_requests to advance_vouchers
ALTER TABLE advance_requests ADD CONSTRAINT fk_advance_request_voucher 
    FOREIGN KEY (voucher_id) REFERENCES advance_vouchers(id);

CREATE TABLE IF NOT EXISTS advance_voucher_lines (
    id SERIAL PRIMARY KEY,
    voucher_id INT REFERENCES advance_vouchers(id) ON DELETE CASCADE,
    line_number INT NOT NULL,
    
    -- Line Type: ADVANCE, EXPEND, SUB_EXPEND ✅
    line_type VARCHAR(20) NOT NULL,
    
    -- Hierarchy for sub-items under EXPEND ✅
    parent_line_id INT REFERENCES advance_voucher_lines(id),
    
    -- Content ✅
    description TEXT NOT NULL,
    amount DECIMAL(12,2) DEFAULT 0.00,
    
    -- 💡 SUGGESTED: Line enhancements
    
    -- 💡 Receipt reference
    receipt_number VARCHAR(50),
    receipt_date DATE,
    
    -- 💡 Vendor for this expense
    vendor_name TEXT,
    
    -- 💡 Tax info
    tax_amount DECIMAL(12,2) DEFAULT 0.00,
    
    UNIQUE(voucher_id, line_number)
);

CREATE TABLE IF NOT EXISTS advance_voucher_documents (
    id SERIAL PRIMARY KEY,
    voucher_id INT REFERENCES advance_vouchers(id) ON DELETE CASCADE,
    
    -- Document Type ✅
    document_type voucher_document_type NOT NULL,
    
    -- Status ✅
    has_document BOOLEAN DEFAULT false,
    
    -- Optional file attachment ✅
    document_path TEXT,
    document_name TEXT,
    
    created_at TIMESTAMP DEFAULT NOW(),
    
    -- 💡 SUGGESTED: Document management
    
    -- 💡 File metadata
    file_size INT,
    mime_type VARCHAR(100),
    
    -- 💡 Verification
    verified BOOLEAN DEFAULT false,
    verified_by INT REFERENCES employees(id),
    verified_date DATE
);

CREATE INDEX idx_advance_vouchers_status ON advance_vouchers(status);
CREATE INDEX idx_advance_vouchers_employee ON advance_vouchers(employee_id);
CREATE INDEX idx_advance_vouchers_date ON advance_vouchers(voucher_date);

-- =====================================================
-- 1.9 💡 SUGGESTED: AUDIT TRAIL TABLE
-- Reason: Track all changes for compliance and debugging
-- =====================================================

CREATE TABLE IF NOT EXISTS audit_log (
    id SERIAL PRIMARY KEY,
    table_name VARCHAR(100) NOT NULL,
    record_id INT NOT NULL,
    action VARCHAR(20) NOT NULL, -- INSERT, UPDATE, DELETE
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
-- 1.10 💡 SUGGESTED: DOCUMENT ATTACHMENTS TABLE
-- Reason: Central storage for all document attachments
-- =====================================================

CREATE TABLE IF NOT EXISTS document_attachments (
    id SERIAL PRIMARY KEY,
    entity_type VARCHAR(50) NOT NULL, -- customer, vendor, po, pr, etc.
    entity_id INT NOT NULL,
    
    -- File info
    file_name TEXT NOT NULL,
    file_path TEXT NOT NULL,
    file_size INT,
    mime_type VARCHAR(100),
    
    -- Description
    document_type VARCHAR(50), -- CONTRACT, CERTIFICATE, PHOTO, etc.
    description TEXT,
    
    -- Metadata
    uploaded_by INT REFERENCES employees(id),
    uploaded_at TIMESTAMP DEFAULT NOW(),
    
    is_active BOOLEAN DEFAULT true
);

CREATE INDEX idx_document_attachments_entity ON document_attachments(entity_type, entity_id);

-- =====================================================
-- 1.11 💡 SUGGESTED: NOTIFICATIONS TABLE
-- Reason: System notifications for approvals, alerts, etc.
-- =====================================================

CREATE TABLE IF NOT EXISTS notifications (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES employees(id) NOT NULL,
    
    -- Content
    title TEXT NOT NULL,
    message TEXT NOT NULL,
    notification_type VARCHAR(50), -- APPROVAL, ALERT, INFO, WARNING
    
    -- Link to related entity
    entity_type VARCHAR(50),
    entity_id INT,
    link_url TEXT,
    
    -- Status
    is_read BOOLEAN DEFAULT false,
    read_at TIMESTAMP,
    
    -- Priority
    priority priority_level DEFAULT 'NORMAL',
    
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP
);

CREATE INDEX idx_notifications_user ON notifications(user_id);
CREATE INDEX idx_notifications_unread ON notifications(user_id, is_read) WHERE is_read = false;

-- =====================================================
-- 1.12 💡 SUGGESTED: APPROVAL WORKFLOW CONFIGURATION
-- Reason: Flexible approval workflows
-- =====================================================

CREATE TABLE IF NOT EXISTS approval_workflows (
    id SERIAL PRIMARY KEY,
    workflow_name VARCHAR(100) NOT NULL,
    entity_type VARCHAR(50) NOT NULL, -- pr, advance_request, po, etc.
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS approval_workflow_steps (
    id SERIAL PRIMARY KEY,
    workflow_id INT REFERENCES approval_workflows(id) ON DELETE CASCADE,
    step_number INT NOT NULL,
    step_name VARCHAR(100) NOT NULL,
    
    -- Who can approve
    approver_role_id INT REFERENCES roles(id),
    approver_department_id INT REFERENCES departments(id),
    specific_approver_id INT REFERENCES employees(id),
    
    -- Conditions
    min_amount DECIMAL(12,2),
    max_amount DECIMAL(12,2),
    
    -- Options
    can_skip BOOLEAN DEFAULT false,
    auto_approve_after_days INT,
    
    UNIQUE(workflow_id, step_number)
);

-- =====================================================
-- 1.13 💡 SUGGESTED: CUSTOMER GROUPS TABLE
-- Reason: Group-based pricing and promotions
-- =====================================================

CREATE TABLE IF NOT EXISTS customer_groups (
    id SERIAL PRIMARY KEY,
    code VARCHAR(20) UNIQUE NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    
    -- Pricing
    default_price_type price_type DEFAULT 'PRICE_1',
    discount_percent DECIMAL(5,2) DEFAULT 0.00,
    
    -- Terms
    default_payment_days INT DEFAULT 30,
    credit_limit_default DECIMAL(12,2) DEFAULT 0.00,
    
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW()
);

ALTER TABLE customers ADD CONSTRAINT fk_customer_group 
    FOREIGN KEY (customer_group_id) REFERENCES customer_groups(id);

-- =====================================================
-- 1.14 ADD PAGES FOR NEW MODULES ✅
-- =====================================================

INSERT INTO pages (page_name, route_name, icon) VALUES 
    ('Purchase Requisitions', '/purchase-requisitions', 'file-text'),
    ('Advance Requests', '/advance-requests', 'wallet'),
    ('Advance Vouchers', '/advance-vouchers', 'receipt'),
    ('Brands', '/brands', 'tag'),
    ('Product Groups', '/product-groups', 'layers'),
    -- 💡 SUGGESTED: Additional pages
    ('Customer Groups', '/customer-groups', 'users'),
    ('Audit Log', '/audit-log', 'history'),
    ('Notifications', '/notifications', 'bell'),
    ('Approval Workflows', '/approval-workflows', 'git-branch'),
    ('Documents', '/documents', 'folder')
ON CONFLICT (route_name) DO NOTHING;

-- =====================================================
-- 1.15 VIEWS FOR ENHANCED INVENTORY DISPLAY ✅
-- =====================================================

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
    p.vat_percent,
    i.expiry_date,
    i.lot_number,
    i.warehouse_id,
    w.name as warehouse_name,
    pc.name as category_name,
    pg.name as product_group_name,
    b.name as brand_name,
    -- 💡 Additional computed fields
    p.price_3,
    p.wholesale_price,
    p.retail_price,
    i.average_cost,
    (i.quantity_on_hand * i.average_cost) as inventory_value,
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

-- 💡 SUGGESTED: Customer summary view
CREATE OR REPLACE VIEW customer_summary AS
SELECT 
    c.*,
    e.english_name as sales_rep_name,
    w.name as warehouse_name,
    cg.name as customer_group_name,
    COALESCE(c.total_orders, 0) as order_count,
    COALESCE(c.total_spent, 0) as lifetime_value,
    c.credit_limit - c.current_balance as available_credit
FROM customers c
LEFT JOIN employees e ON c.sales_rep_id = e.id
LEFT JOIN warehouses w ON c.default_warehouse_id = w.id
LEFT JOIN customer_groups cg ON c.customer_group_id = cg.id;

-- 💡 SUGGESTED: Vendor summary view
CREATE OR REPLACE VIEW vendor_summary AS
SELECT 
    v.*,
    e.english_name as buyer_name,
    COALESCE(v.total_po_count, 0) as order_count,
    COALESCE(v.total_po_value, 0) as total_purchases
FROM vendors v
LEFT JOIN employees e ON v.buyer_id = e.id;
```

---

## Phase 2: Backend Models (Go)

### 2.1 Enhanced Customer Model

**File: `registration/src/v1/models/customer.go`**

```go
package models

// ============================================
// Customer Models - Enhanced for Complete System
// ============================================

// PriceType enum ✅
type PriceType string

const (
	PriceType1     PriceType = "PRICE_1"
	PriceType2     PriceType = "PRICE_2"
	PriceType3     PriceType = "PRICE_3"
	PriceWholesale PriceType = "WHOLESALE"
	PriceRetail    PriceType = "RETAIL"
)

// CustomerLevel enum ✅
type CustomerLevel string

const (
	LevelGeneral  CustomerLevel = "GENERAL"
	LevelSilver   CustomerLevel = "SILVER"
	LevelGold     CustomerLevel = "GOLD"
	LevelPlatinum CustomerLevel = "PLATINUM"
	LevelVIP      CustomerLevel = "VIP"
)

// 💡 SUGGESTED: Contact preference enum
type ContactPreference string

const (
	ContactPhone    ContactPreference = "PHONE"
	ContactEmail    ContactPreference = "EMAIL"
	ContactFax      ContactPreference = "FAX"
	ContactMobile   ContactPreference = "MOBILE"
	ContactWhatsApp ContactPreference = "WHATSAPP"
	ContactLine     ContactPreference = "LINE"
)

// 💡 SUGGESTED: Customer type for B2B/B2C
type CustomerType string

const (
	CustomerTypeRetail      CustomerType = "RETAIL"
	CustomerTypeWholesale   CustomerType = "WHOLESALE"
	CustomerTypeDistributor CustomerType = "DISTRIBUTOR"
	CustomerTypeChain       CustomerType = "CHAIN"
)

type Customer struct {
	ID                 int        `json:"id"`
	CustomerCode       string     `json:"customer_code"`
	Name               string     `json:"name"`
	
	// ✅ REQUIRED: Contact Information from DocText.md
	DefaultContactName string     `json:"default_contact_name,omitempty"`
	FullAddress        string     `json:"full_address,omitempty"`
	DeliveryName       string     `json:"delivery_name,omitempty"`
	Tel                string     `json:"tel,omitempty"`
	Fax                string     `json:"fax,omitempty"`
	Mobile             string     `json:"mobile,omitempty"`
	Email              string     `json:"email,omitempty"`
	WebEmail           string     `json:"web_email,omitempty"`
	
	// ✅ REQUIRED: Tax & Identity from DocText.md
	TaxCode            string     `json:"tax_code,omitempty"`
	PassportNo         string     `json:"passport_no,omitempty"`
	CountryCode        string     `json:"country_code,omitempty"`
	Gender             string     `json:"gender,omitempty"`
	Birthday           CustomDate `json:"birthday,omitempty"`
	
	// ✅ REQUIRED: Account Settings from DocText.md
	PriceType          PriceType     `json:"price_type"`
	CustomerLevel      CustomerLevel `json:"customer_level"`
	BillDiscount       float64    `json:"bill_discount"`
	CreditLimit        float64    `json:"credit_limit"`
	CreditDays         int        `json:"credit_days"`
	CurrentBalance     float64    `json:"current_balance"`
	PaymentTermsDays   int        `json:"payment_terms_days"`
	Currency           string     `json:"currency"`
	
	// ✅ REQUIRED: Membership from DocText.md
	MemberExpire       CustomDate `json:"member_expire,omitempty"`
	UsePromotion       bool       `json:"use_promotion"`
	UsePMTBill         bool       `json:"use_pmt_bill"`
	CollectPoint       bool       `json:"collect_point"`
	
	// ✅ REQUIRED: Other from DocText.md
	Payee              string     `json:"payee,omitempty"`
	RDBranch           string     `json:"rd_branch,omitempty"`
	Barcode            string     `json:"barcode,omitempty"`
	Notes              string     `json:"notes,omitempty"`
	
	// 💡 SUGGESTED: Analytics & Engagement
	// Reason: Customer relationship management and sales analysis
	LastOrderDate      CustomDate `json:"last_order_date,omitempty"`
	TotalOrders        int        `json:"total_orders"`
	TotalSpent         float64    `json:"total_spent"`
	LoyaltyPoints      float64    `json:"loyalty_points"`
	CreditScore        int        `json:"credit_score"`
	
	// 💡 SUGGESTED: Delivery preferences
	// Reason: Improves delivery efficiency and customer satisfaction
	PreferredDeliveryTime string `json:"preferred_delivery_time,omitempty"`
	DeliveryInstructions  string `json:"delivery_instructions,omitempty"`
	Latitude              *float64 `json:"latitude,omitempty"`
	Longitude             *float64 `json:"longitude,omitempty"`
	
	// 💡 SUGGESTED: Referral tracking
	// Reason: Supports referral programs
	ReferredByID       *int   `json:"referred_by_id,omitempty"`
	ReferralCode       string `json:"referral_code,omitempty"`
	AcquisitionSource  string `json:"acquisition_source,omitempty"`
	
	// 💡 SUGGESTED: Communication preferences
	// Reason: Modern multi-channel communication
	PreferredContact       ContactPreference `json:"preferred_contact,omitempty"`
	SecondaryContactName   string `json:"secondary_contact_name,omitempty"`
	SecondaryContactPhone  string `json:"secondary_contact_phone,omitempty"`
	LineID                 string `json:"line_id,omitempty"`
	WhatsApp               string `json:"whatsapp,omitempty"`
	
	// 💡 SUGGESTED: B2B classification
	// Reason: Different handling for business customers
	CustomerType           CustomerType `json:"customer_type,omitempty"`
	CompanyRegistrationNo  string `json:"company_registration_no,omitempty"`
	VATRegistrationNo      string `json:"vat_registration_no,omitempty"`
	
	// 💡 SUGGESTED: Account management
	// Reason: Separates order-taking from relationship management
	AccountManagerID   *int `json:"account_manager_id,omitempty"`
	CustomerGroupID    *int `json:"customer_group_id,omitempty"`
	
	// 💡 SUGGESTED: Block status
	// Reason: Prevent orders for problematic customers
	IsBlocked          bool       `json:"is_blocked"`
	BlockReason        string     `json:"block_reason,omitempty"`
	BlockedDate        CustomDate `json:"blocked_date,omitempty"`
	BlockedBy          *int       `json:"blocked_by,omitempty"`
	
	// References
	BillingAddressID   *int       `json:"billing_address_id,omitempty"`
	SalesRepID         *int       `json:"sales_rep_id,omitempty"`
	DefaultRouteID     *int       `json:"default_route_id,omitempty"`
	DefaultWarehouseID *int       `json:"default_warehouse_id,omitempty"`
	
	// Flags
	TaxExempt          bool       `json:"tax_exempt"`
	IsActive           bool       `json:"is_active"`
	
	// Audit
	CreatedBy          int        `json:"created_by"`
	CreatedAt          CustomDate `json:"created_at"`
	UpdatedAt          CustomDate `json:"updated_at"`
}

// CreateCustomerRequest with all fields
type CreateCustomerRequest struct {
	// ✅ Required fields
	CustomerCode       string  `json:"customer_code"`
	Name               string  `json:"name"`
	
	// ✅ DocText.md fields
	DefaultContactName string  `json:"default_contact_name,omitempty"`
	FullAddress        string  `json:"full_address,omitempty"`
	DeliveryName       string  `json:"delivery_name,omitempty"`
	Tel                string  `json:"tel,omitempty"`
	Fax                string  `json:"fax,omitempty"`
	Mobile             string  `json:"mobile,omitempty"`
	Email              string  `json:"email,omitempty"`
	WebEmail           string  `json:"web_email,omitempty"`
	TaxCode            string  `json:"tax_code,omitempty"`
	PassportNo         string  `json:"passport_no,omitempty"`
	CountryCode        string  `json:"country_code,omitempty"`
	Gender             string  `json:"gender,omitempty"`
	Birthday           string  `json:"birthday,omitempty"`
	PriceType          PriceType     `json:"price_type"`
	CustomerLevel      CustomerLevel `json:"customer_level"`
	BillDiscount       float64 `json:"bill_discount"`
	CreditLimit        float64 `json:"credit_limit"`
	CreditDays         int     `json:"credit_days"`
	PaymentTermsDays   int     `json:"payment_terms_days"`
	Currency           string  `json:"currency"`
	MemberExpire       string  `json:"member_expire,omitempty"`
	UsePromotion       bool    `json:"use_promotion"`
	UsePMTBill         bool    `json:"use_pmt_bill"`
	CollectPoint       bool    `json:"collect_point"`
	Payee              string  `json:"payee,omitempty"`
	RDBranch           string  `json:"rd_branch,omitempty"`
	Barcode            string  `json:"barcode,omitempty"`
	Notes              string  `json:"notes,omitempty"`
	SalesRepID         *int    `json:"sales_rep_id,omitempty"`
	DefaultWarehouseID *int    `json:"default_warehouse_id,omitempty"`
	TaxExempt          bool    `json:"tax_exempt"`
	
	// 💡 Suggested fields
	PreferredDeliveryTime string `json:"preferred_delivery_time,omitempty"`
	DeliveryInstructions  string `json:"delivery_instructions,omitempty"`
	Latitude              *float64 `json:"latitude,omitempty"`
	Longitude             *float64 `json:"longitude,omitempty"`
	PreferredContact      ContactPreference `json:"preferred_contact,omitempty"`
	SecondaryContactName  string `json:"secondary_contact_name,omitempty"`
	SecondaryContactPhone string `json:"secondary_contact_phone,omitempty"`
	LineID                string `json:"line_id,omitempty"`
	WhatsApp              string `json:"whatsapp,omitempty"`
	CustomerType          CustomerType `json:"customer_type,omitempty"`
	CompanyRegistrationNo string `json:"company_registration_no,omitempty"`
	VATRegistrationNo     string `json:"vat_registration_no,omitempty"`
	AccountManagerID      *int   `json:"account_manager_id,omitempty"`
	CustomerGroupID       *int   `json:"customer_group_id,omitempty"`
}
```

### 2.2 Enhanced Vendor Model

**File: `registration/src/v1/models/vendor.go` - Additional fields**

```go
// Add to existing Vendor struct:

type Vendor struct {
	// ... existing fields ...
	
	// ✅ REQUIRED FROM DocText.md
	ContactName     string  `json:"contact_name,omitempty"`
	RDBranch        string  `json:"rd_branch,omitempty"`
	PictureName     string  `json:"picture_name,omitempty"`
	Mobile          string  `json:"mobile,omitempty"`
	Fax             string  `json:"fax,omitempty"`
	TaxCode         string  `json:"tax_code,omitempty"`
	CreditDays      int     `json:"credit_days"`
	DiscountPercent float64 `json:"discount_percent"`
	POCreditMoney   float64 `json:"po_credit_money"`
	
	// 💡 SUGGESTED: Vendor performance metrics
	// Reason: Enables vendor evaluation for strategic sourcing
	VendorRating         float64 `json:"vendor_rating"`          // 0-5 scale
	QualityScore         float64 `json:"quality_score"`          // percentage
	OnTimeDeliveryRate   float64 `json:"on_time_delivery_rate"`  // percentage
	
	// 💡 SUGGESTED: Order tracking
	// Reason: Quick reference for vendor activity
	LastPODate           CustomDate `json:"last_po_date,omitempty"`
	TotalPOValue         float64    `json:"total_po_value"`
	TotalPOCount         int        `json:"total_po_count"`
	
	// 💡 SUGGESTED: Web presence
	// Reason: Quick access to vendor information
	Website              string `json:"website,omitempty"`
	
	// 💡 SUGGESTED: Banking info for payments
	// Reason: Required for AP processing
	BankName             string `json:"bank_name,omitempty"`
	BankAccountNumber    string `json:"bank_account_number,omitempty"`
	BankAccountName      string `json:"bank_account_name,omitempty"`
	BankBranch           string `json:"bank_branch,omitempty"`
	SwiftCode            string `json:"swift_code,omitempty"`
	
	// 💡 SUGGESTED: Payment preferences
	// Reason: Streamlines AP processing
	PreferredPayment     string `json:"preferred_payment,omitempty"` // BANK_TRANSFER, CHECK, etc.
	
	// 💡 SUGGESTED: Certifications (important for food industry)
	// Reason: Track vendor compliance (HACCP, ISO, etc.)
	Certifications       string     `json:"certifications,omitempty"` // JSON
	CertificationExpiry  CustomDate `json:"certification_expiry,omitempty"`
	
	// 💡 SUGGESTED: Notes
	// Reason: Store important vendor-specific information
	Notes                string `json:"notes,omitempty"`
	
	// 💡 SUGGESTED: Communication preferences
	PreferredContact     ContactPreference `json:"preferred_contact,omitempty"`
	SecondaryContactName  string `json:"secondary_contact_name,omitempty"`
	SecondaryContactPhone string `json:"secondary_contact_phone,omitempty"`
	
	// 💡 SUGGESTED: Vendor classification
	// Reason: Different vendor types may have different handling
	VendorType           string `json:"vendor_type,omitempty"` // SUPPLIER, MANUFACTURER, DISTRIBUTOR, IMPORTER
	
	// 💡 SUGGESTED: Return policy
	// Reason: Important for handling defective/expired goods
	ReturnPolicy         string `json:"return_policy,omitempty"`
	ReturnWindowDays     int    `json:"return_window_days"`
	
	// 💡 SUGGESTED: Block status
	// Reason: Prevent POs to problematic vendors
	IsBlocked            bool   `json:"is_blocked"`
	BlockReason          string `json:"block_reason,omitempty"`
}
```

### 2.3 Enhanced Purchase Order Model

**File: `registration/src/v1/models/purchase_order.go` - Additional fields**

```go
// Add to existing PurchaseOrder struct:

type PurchaseOrder struct {
	// ... existing fields ...
	
	// ✅ REQUIRED FROM DocText.md
	ContactPerson       string     `json:"contact_person,omitempty"`
	CustomerName        string     `json:"customer_name,omitempty"`
	Address             string     `json:"address,omitempty"`
	Tel                 string     `json:"tel,omitempty"`
	Email               string     `json:"email,omitempty"`
	DocumentRef         string     `json:"document_ref,omitempty"`
	TermsOfPayment      string     `json:"terms_of_payment,omitempty"`
	Currency            string     `json:"currency"`
	TradeDiscount       float64    `json:"trade_discount"`
	AmountAfterDiscount float64    `json:"amount_after_discount"`
	VATRate             float64    `json:"vat_rate"`
	PreparedBy          *int       `json:"prepared_by,omitempty"`
	AuthorizedBy        *int       `json:"authorized_by,omitempty"`
	CheckedBy           *int       `json:"checked_by,omitempty"`
	PreparedDate        CustomDate `json:"prepared_date,omitempty"`
	AuthorizedDate      CustomDate `json:"authorized_date,omitempty"`
	CheckedDate         CustomDate `json:"checked_date,omitempty"`
	
	// 💡 SUGGESTED: Shipping & logistics
	// Reason: Track delivery and improve receiving process
	ShippingMethod       string  `json:"shipping_method,omitempty"`
	ShippingCarrier      string  `json:"shipping_carrier,omitempty"`
	TrackingNumber       string  `json:"tracking_number,omitempty"`
	InsuranceAmount      float64 `json:"insurance_amount"`
	DeliveryInstructions string  `json:"delivery_instructions,omitempty"`
	QualityRequirements  string  `json:"quality_requirements,omitempty"`
	
	// 💡 SUGGESTED: Multi-currency support
	// Reason: Critical for international procurement
	ExchangeRate          float64 `json:"exchange_rate"`
	BaseCurrency          string  `json:"base_currency"`
	TotalInBaseCurrency   float64 `json:"total_in_base_currency"`
	
	// 💡 SUGGESTED: Priority
	// Reason: Helps prioritize receiving
	Priority              string `json:"priority"` // LOW, NORMAL, HIGH, URGENT, CRITICAL
	
	// 💡 SUGGESTED: Document linkage
	// Reason: Connect PO to related documents
	RequisitionID         *int   `json:"requisition_id,omitempty"`
	QuotationRef          string `json:"quotation_ref,omitempty"`
	ContractRef           string `json:"contract_ref,omitempty"`
	
	// 💡 SUGGESTED: Payment tracking
	// Reason: Quick payment status without going to AP
	PaymentStatus         string  `json:"payment_status"` // PENDING, PARTIAL, PAID
	AmountPaid            float64 `json:"amount_paid"`
	
	// 💡 SUGGESTED: Cancellation tracking
	// Reason: Track why POs were cancelled
	CancelledBy           *int       `json:"cancelled_by,omitempty"`
	CancelledDate         CustomDate `json:"cancelled_date,omitempty"`
	CancellationReason    string     `json:"cancellation_reason,omitempty"`
}

// Enhanced PurchaseOrderLine
type PurchaseOrderLine struct {
	// ... existing fields ...
	
	// ✅ REQUIRED: Line discount
	DiscountPercent  float64 `json:"discount_percent"`
	DiscountAmount   float64 `json:"discount_amount"`
	
	// 💡 SUGGESTED: Line-level tax
	// Reason: Some items may have different tax rates
	TaxRate          float64 `json:"tax_rate"`
	TaxAmount        float64 `json:"tax_amount"`
	
	// 💡 SUGGESTED: Line-specific delivery date
	// Reason: Different items may need different delivery dates
	RequestedDate    CustomDate `json:"requested_date,omitempty"`
	
	// 💡 SUGGESTED: Line notes
	// Reason: Line-specific instructions
	Notes            string `json:"notes,omitempty"`
}
```

---

## Phase 3: Additional Backend Models

### 3.1 Notification Model (💡 Suggested)

**File: `registration/src/v1/models/notification.go`**

```go
package models

// 💡 SUGGESTED: Notification system
// Reason: Enables in-app notifications for approvals, alerts, etc.

type PriorityLevel string

const (
	PriorityLow      PriorityLevel = "LOW"
	PriorityNormal   PriorityLevel = "NORMAL"
	PriorityHigh     PriorityLevel = "HIGH"
	PriorityUrgent   PriorityLevel = "URGENT"
	PriorityCritical PriorityLevel = "CRITICAL"
)

type Notification struct {
	ID               int           `json:"id"`
	UserID           int           `json:"user_id"`
	Title            string        `json:"title"`
	Message          string        `json:"message"`
	NotificationType string        `json:"notification_type"` // APPROVAL, ALERT, INFO, WARNING
	EntityType       string        `json:"entity_type,omitempty"`
	EntityID         *int          `json:"entity_id,omitempty"`
	LinkURL          string        `json:"link_url,omitempty"`
	IsRead           bool          `json:"is_read"`
	ReadAt           CustomDate    `json:"read_at,omitempty"`
	Priority         PriorityLevel `json:"priority"`
	CreatedAt        CustomDate    `json:"created_at"`
	ExpiresAt        CustomDate    `json:"expires_at,omitempty"`
}
```

### 3.2 Audit Log Model (💡 Suggested)

**File: `registration/src/v1/models/audit_log.go`**

```go
package models

// 💡 SUGGESTED: Audit logging
// Reason: Compliance, debugging, and security tracking

type AuditLog struct {
	ID         int        `json:"id"`
	TableName  string     `json:"table_name"`
	RecordID   int        `json:"record_id"`
	Action     string     `json:"action"` // INSERT, UPDATE, DELETE
	OldValues  string     `json:"old_values,omitempty"` // JSON
	NewValues  string     `json:"new_values,omitempty"` // JSON
	ChangedBy  int        `json:"changed_by"`
	ChangedAt  CustomDate `json:"changed_at"`
	IPAddress  string     `json:"ip_address,omitempty"`
	UserAgent  string     `json:"user_agent,omitempty"`
}

type AuditLogFilters struct {
	TableName  string `json:"table_name,omitempty"`
	RecordID   *int   `json:"record_id,omitempty"`
	Action     string `json:"action,omitempty"`
	ChangedBy  *int   `json:"changed_by,omitempty"`
	DateFrom   string `json:"date_from,omitempty"`
	DateTo     string `json:"date_to,omitempty"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
}
```

### 3.3 Document Attachment Model (💡 Suggested)

**File: `registration/src/v1/models/document_attachment.go`**

```go
package models

// 💡 SUGGESTED: Document attachment system
// Reason: Centralized document storage for all entities

type DocumentAttachment struct {
	ID           int        `json:"id"`
	EntityType   string     `json:"entity_type"` // customer, vendor, po, pr, etc.
	EntityID     int        `json:"entity_id"`
	FileName     string     `json:"file_name"`
	FilePath     string     `json:"file_path"`
	FileSize     int        `json:"file_size,omitempty"`
	MimeType     string     `json:"mime_type,omitempty"`
	DocumentType string     `json:"document_type,omitempty"` // CONTRACT, CERTIFICATE, PHOTO, etc.
	Description  string     `json:"description,omitempty"`
	UploadedBy   int        `json:"uploaded_by"`
	UploadedAt   CustomDate `json:"uploaded_at"`
	IsActive     bool       `json:"is_active"`
}
```

---

## Phase 4: Frontend Type Definitions

### 4.1 Enhanced Customer Interface

**File: `foodhive-erp-frontend/client/src/types/customer.ts`**

```typescript
// ✅ Required + 💡 Suggested fields

export type PriceType = 'PRICE_1' | 'PRICE_2' | 'PRICE_3' | 'WHOLESALE' | 'RETAIL';
export type CustomerLevel = 'GENERAL' | 'SILVER' | 'GOLD' | 'PLATINUM' | 'VIP';
export type ContactPreference = 'PHONE' | 'EMAIL' | 'FAX' | 'MOBILE' | 'WHATSAPP' | 'LINE';
export type CustomerType = 'RETAIL' | 'WHOLESALE' | 'DISTRIBUTOR' | 'CHAIN';

export interface Customer {
  id: number;
  customer_code: string;
  name: string;
  
  // ✅ Contact Information (DocText.md)
  default_contact_name?: string;
  full_address?: string;
  delivery_name?: string;
  tel?: string;
  fax?: string;
  mobile?: string;
  email?: string;
  web_email?: string;
  
  // ✅ Tax & Identity (DocText.md)
  tax_code?: string;
  passport_no?: string;
  country_code?: string;
  gender?: 'MALE' | 'FEMALE';
  birthday?: string;
  
  // ✅ Account Settings (DocText.md)
  price_type?: PriceType;
  customer_level?: CustomerLevel;
  bill_discount?: number;
  credit_limit: number;
  credit_days?: number;
  current_balance?: number;
  payment_terms_days: number;
  currency?: string;
  
  // ✅ Membership (DocText.md)
  member_expire?: string;
  use_promotion?: boolean;
  use_pmt_bill?: boolean;
  collect_point?: boolean;
  
  // ✅ Other (DocText.md)
  payee?: string;
  rd_branch?: string;
  barcode?: string;
  notes?: string;
  
  // 💡 SUGGESTED: Analytics (for CRM features)
  last_order_date?: string;
  total_orders?: number;
  total_spent?: number;
  loyalty_points?: number;
  credit_score?: number;
  
  // 💡 SUGGESTED: Delivery preferences
  preferred_delivery_time?: string;
  delivery_instructions?: string;
  latitude?: number;
  longitude?: number;
  
  // 💡 SUGGESTED: Communication
  preferred_contact?: ContactPreference;
  secondary_contact_name?: string;
  secondary_contact_phone?: string;
  line_id?: string;
  whatsapp?: string;
  
  // 💡 SUGGESTED: B2B
  customer_type?: CustomerType;
  company_registration_no?: string;
  vat_registration_no?: string;
  
  // 💡 SUGGESTED: Management
  account_manager_id?: number;
  customer_group_id?: number;
  
  // 💡 SUGGESTED: Block status
  is_blocked?: boolean;
  block_reason?: string;
  
  // Standard fields
  sales_rep_id?: number;
  default_warehouse_id?: number;
  tax_exempt?: boolean;
  is_active: boolean;
  created_at?: string;
  updated_at?: string;
}
```

---

## Phase 5: Additional Frontend Services

### 5.1 Notification Service (💡 Suggested)

**File: `foodhive-erp-frontend/client/src/services/notificationService.ts`**

```typescript
import api from '@/lib/api';

// 💡 SUGGESTED: Real-time notification system
// Reason: Enables in-app notifications for approvals, alerts, etc.

export interface Notification {
  id: number;
  user_id: number;
  title: string;
  message: string;
  notification_type: 'APPROVAL' | 'ALERT' | 'INFO' | 'WARNING';
  entity_type?: string;
  entity_id?: number;
  link_url?: string;
  is_read: boolean;
  read_at?: string;
  priority: 'LOW' | 'NORMAL' | 'HIGH' | 'URGENT' | 'CRITICAL';
  created_at: string;
  expires_at?: string;
}

export const notificationService = {
  // Get all notifications for current user
  getNotifications: async (params?: { unread_only?: boolean; limit?: number }) => {
    const response = await api.get('/notifications/list', { params });
    return response.data?.data || response.data || [];
  },

  // Get unread count
  getUnreadCount: async () => {
    const response = await api.get('/notifications/unread-count');
    return response.data?.data?.count || 0;
  },

  // Mark as read
  markAsRead: async (id: string) => {
    const response = await api.post(`/notifications/read/${id}`);
    return response.data;
  },

  // Mark all as read
  markAllAsRead: async () => {
    const response = await api.post('/notifications/read-all');
    return response.data;
  },

  // Delete notification
  delete: async (id: string) => {
    const response = await api.delete(`/notifications/delete/${id}`);
    return response.data;
  }
};
```

### 5.2 Audit Log Service (💡 Suggested)

**File: `foodhive-erp-frontend/client/src/services/auditService.ts`**

```typescript
import api from '@/lib/api';

// 💡 SUGGESTED: Audit log viewing
// Reason: Compliance and debugging

export interface AuditLog {
  id: number;
  table_name: string;
  record_id: number;
  action: 'INSERT' | 'UPDATE' | 'DELETE';
  old_values?: object;
  new_values?: object;
  changed_by: number;
  changed_by_name?: string;
  changed_at: string;
  ip_address?: string;
}

export const auditService = {
  // Get audit logs with filters
  getLogs: async (params?: {
    table_name?: string;
    record_id?: number;
    action?: string;
    changed_by?: number;
    date_from?: string;
    date_to?: string;
    page?: number;
    page_size?: number;
  }) => {
    const response = await api.get('/audit-log/list', { params });
    return response.data?.data || response.data || [];
  },

  // Get logs for specific entity
  getEntityLogs: async (entityType: string, entityId: string) => {
    const response = await api.get(`/audit-log/entity/${entityType}/${entityId}`);
    return response.data?.data || response.data || [];
  }
};
```

---

## Phase 6: System Enhancements (💡 Suggested Summary)

### 6.1 Why These Suggestions?

| Feature | Business Value | Technical Benefit |
|---------|----------------|-------------------|
| **Customer Analytics** | Better customer insights, targeted marketing | Performance optimization, reporting |
| **Vendor Performance** | Strategic sourcing, quality control | Automated vendor selection |
| **GPS Coordinates** | Route optimization, delivery efficiency | Maps integration |
| **Loyalty Points** | Customer retention, repeat business | Marketing automation |
| **Block Status** | Risk management, AR control | Prevent bad orders |
| **Audit Logging** | Compliance, debugging | Security, troubleshooting |
| **Notifications** | User engagement, workflow efficiency | Real-time updates |
| **Document Attachments** | Paperless operations, compliance | File management |
| **Multi-currency Exchange** | International trade support | Accurate financials |
| **Approval Workflows** | Flexible business rules | Configurable processes |

### 6.2 Implementation Priority

| Priority | Feature | Effort | Business Impact |
|----------|---------|--------|-----------------|
| **P1 - Critical** | All DocText.md requirements | High | Must have |
| **P2 - High** | Audit logging, Block status | Medium | Compliance |
| **P3 - Medium** | Notifications, Analytics | Medium | Efficiency |
| **P4 - Low** | GPS, Loyalty, Workflows | Low | Nice to have |

---

## Implementation Checklist

### Database (Phase 1)
- [ ] Run SQL migration `006_requirements_gap_updates.sql`
- [ ] Verify all new columns and tables created
- [ ] Test constraints and indexes

### Backend (Phase 2-3)
- [ ] Update Customer model with all fields
- [ ] Update Vendor model with all fields
- [ ] Update PurchaseOrder model with all fields
- [ ] Create PurchaseRequisition module
- [ ] Create AdvanceRequest module
- [ ] Create AdvanceVoucher module
- [ ] 💡 Create Notification module (suggested)
- [ ] 💡 Create AuditLog module (suggested)
- [ ] 💡 Create DocumentAttachment module (suggested)

### Frontend (Phase 4-5)
- [ ] Update Customer form with all fields
- [ ] Update Vendor form with all fields
- [ ] Update PurchaseOrder form with all fields
- [ ] Create PurchaseRequisition pages
- [ ] Create AdvanceRequest pages
- [ ] Create AdvanceVoucher pages
- [ ] 💡 Create Notification component (suggested)
- [ ] 💡 Create AuditLog viewer (suggested)

### Testing
- [ ] Unit tests for new models
- [ ] Integration tests for new services
- [ ] Frontend component tests
- [ ] End-to-end workflow tests

---

*Document Version: 2.0 Enhanced*
*Last Updated: January 2026*
*Author: FoodHive Development Team*

---

## Notes on Suggested Fields

All suggested fields are marked with 💡 and include:
1. **Comment with reason** - Why this field is valuable
2. **No deletion** - All existing fields preserved
3. **Optional** - Can be implemented in later phases
4. **Backward compatible** - New fields have defaults

The suggested enhancements transform this from a basic ERP to a complete enterprise system with:
- Customer relationship management (CRM)
- Vendor performance management
- Compliance and audit capabilities
- Modern communication channels
- Advanced analytics foundation
