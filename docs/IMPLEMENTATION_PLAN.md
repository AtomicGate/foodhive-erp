# FoodHive ERP - Full Implementation Plan
## Bridging the Gap: From Current State to DocText.md Requirements

This document provides a comprehensive, step-by-step implementation plan to align the current FoodHive ERP system with the requirements specified in DocText.md.

---

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [Phase 1: Database Schema Changes](#phase-1-database-schema-changes)
3. [Phase 2: Backend Models (Go)](#phase-2-backend-models-go)
4. [Phase 3: Backend Services & Routes](#phase-3-backend-services--routes)
5. [Phase 4: Frontend Services](#phase-4-frontend-services)
6. [Phase 5: Frontend Components](#phase-5-frontend-components)
7. [Implementation Priority & Timeline](#implementation-priority--timeline)

---

## Executive Summary

### Current State Analysis

| Module | Current Fields | Required Fields | Gap |
|--------|---------------|-----------------|-----|
| Customers | 17 fields | 37+ fields | **20+ fields missing** |
| Vendors | 14 fields | 22+ fields | **8+ fields missing** |
| Purchase Orders | 13 fields | 23+ fields | **10+ fields missing** |
| Purchase Requisition | ❌ Not exists | Full module | **NEW MODULE** |
| Advance Request | ❌ Not exists | Full module | **NEW MODULE** |
| Advance Voucher | ❌ Not exists | Full module | **NEW MODULE** |
| Products | 17 fields | 20+ fields | **3+ fields missing** |
| Inventory Display | Basic | Full with pricing | **Enhanced views needed** |

---

## Phase 1: Database Schema Changes

### File: `sql/006_requirements_gap_updates.sql`

```sql
-- =====================================================
-- PHASE 1: Database Schema Updates for DocText.md Requirements
-- Version: 2.0
-- =====================================================

-- =====================================================
-- 1.1 ENUM TYPES FOR NEW MODULES
-- =====================================================

-- Price Type Enum
CREATE TYPE price_type AS ENUM ('PRICE_1', 'PRICE_2', 'PRICE_3', 'WHOLESALE', 'RETAIL');

-- Customer Level Enum
CREATE TYPE customer_level AS ENUM ('GENERAL', 'SILVER', 'GOLD', 'PLATINUM', 'VIP');

-- Currency Enum (extend for multi-currency)
CREATE TYPE currency_code AS ENUM ('LAK', 'THB', 'USD', 'KIP');

-- Purchase Requisition Status
CREATE TYPE pr_status AS ENUM ('DRAFT', 'SUBMITTED', 'CHECKED', 'AUTHORIZED', 'CONVERTED', 'CANCELLED');

-- Advance Request Status
CREATE TYPE advance_request_status AS ENUM ('DRAFT', 'SUBMITTED', 'APPROVED', 'REJECTED', 'CLEARED');

-- Advance Voucher Status
CREATE TYPE advance_voucher_status AS ENUM ('DRAFT', 'SUBMITTED', 'APPROVED', 'REJECTED', 'COMPLETED');

-- Document Type for Advance Voucher
CREATE TYPE voucher_document_type AS ENUM ('RECEIPT', 'INVOICE', 'BANK_TRANSFER', 'PO_PR', 'TRANSPORTATION');

-- =====================================================
-- 1.2 CUSTOMER TABLE UPDATES
-- =====================================================

-- Add missing fields to customers table
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

-- =====================================================
-- 1.3 VENDOR TABLE UPDATES
-- =====================================================

-- Add missing fields to vendors table
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS contact_name TEXT;
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS rd_branch TEXT;
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS picture_name TEXT;
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS mobile VARCHAR(50);
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS fax VARCHAR(50);
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS tax_code VARCHAR(50);
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS credit_days INT DEFAULT 0;
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS discount_percent DECIMAL(5,2) DEFAULT 0.00;
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS po_credit_money DECIMAL(12,2) DEFAULT 0.00;

-- =====================================================
-- 1.4 PURCHASE ORDER TABLE UPDATES
-- =====================================================

-- Add missing fields to purchase_orders table
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

-- Add discount to purchase_order_lines
ALTER TABLE purchase_order_lines ADD COLUMN IF NOT EXISTS discount_percent DECIMAL(5,2) DEFAULT 0.00;
ALTER TABLE purchase_order_lines ADD COLUMN IF NOT EXISTS discount_amount DECIMAL(12,2) DEFAULT 0.00;

-- =====================================================
-- 1.5 PRODUCT TABLE UPDATES
-- =====================================================

-- Add brand and product group to products
ALTER TABLE products ADD COLUMN IF NOT EXISTS brand TEXT;
ALTER TABLE products ADD COLUMN IF NOT EXISTS product_group TEXT;
ALTER TABLE products ADD COLUMN IF NOT EXISTS price_1 DECIMAL(12,4) DEFAULT 0.00;
ALTER TABLE products ADD COLUMN IF NOT EXISTS price_2 DECIMAL(12,4) DEFAULT 0.00;
ALTER TABLE products ADD COLUMN IF NOT EXISTS vat_percent DECIMAL(5,2) DEFAULT 7.00;

-- Create product_groups table
CREATE TABLE IF NOT EXISTS product_groups (
    id SERIAL PRIMARY KEY,
    code VARCHAR(20) UNIQUE NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    parent_id INT REFERENCES product_groups(id),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create brands table
CREATE TABLE IF NOT EXISTS brands (
    id SERIAL PRIMARY KEY,
    code VARCHAR(20) UNIQUE NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    logo_url TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Add foreign keys to products
ALTER TABLE products ADD COLUMN IF NOT EXISTS brand_id INT REFERENCES brands(id);
ALTER TABLE products ADD COLUMN IF NOT EXISTS product_group_id INT REFERENCES product_groups(id);

-- =====================================================
-- 1.6 NEW MODULE: PURCHASE REQUISITION
-- =====================================================

CREATE TABLE IF NOT EXISTS purchase_requisitions (
    id SERIAL PRIMARY KEY,
    requisition_number VARCHAR(20) UNIQUE NOT NULL,
    document_number VARCHAR(50),
    
    -- Requester Information
    requester_id INT REFERENCES employees(id),
    department_id INT REFERENCES departments(id),
    reason TEXT,
    
    -- Supplier Information
    supplier_id INT REFERENCES vendors(id),
    
    -- Dates
    document_date DATE DEFAULT CURRENT_DATE,
    request_date DATE DEFAULT CURRENT_DATE,
    
    -- Amounts
    total_amount DECIMAL(12,2) DEFAULT 0.00,
    
    -- Notes & Remarks
    remark TEXT,
    
    -- Approval Chain
    checked_by INT REFERENCES employees(id),
    checked_date DATE,
    authorized_by INT REFERENCES employees(id),
    authorized_date DATE,
    purchasing_employee_id INT REFERENCES employees(id),
    purchasing_date DATE,
    
    -- Linked PO (if converted)
    converted_po_id INT REFERENCES purchase_orders(id),
    
    -- Status & Audit
    status pr_status DEFAULT 'DRAFT',
    created_by INT REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS purchase_requisition_lines (
    id SERIAL PRIMARY KEY,
    requisition_id INT REFERENCES purchase_requisitions(id) ON DELETE CASCADE,
    line_number INT NOT NULL,
    
    -- Product Information
    product_id INT REFERENCES products(id),
    product_code VARCHAR(50),
    description TEXT,
    
    -- Stock & Quantity
    stock_balance DECIMAL(12,3) DEFAULT 0.00,
    quantity_requested DECIMAL(10,3) NOT NULL,
    
    -- Pricing
    unit_of_measure VARCHAR(10),
    unit_price DECIMAL(10,4) DEFAULT 0.00,
    amount DECIMAL(12,2) DEFAULT 0.00,
    
    UNIQUE(requisition_id, line_number)
);

CREATE INDEX idx_purchase_requisitions_status ON purchase_requisitions(status);
CREATE INDEX idx_purchase_requisitions_requester ON purchase_requisitions(requester_id);
CREATE INDEX idx_purchase_requisitions_date ON purchase_requisitions(document_date);

-- =====================================================
-- 1.7 NEW MODULE: ADVANCE REQUEST
-- =====================================================

CREATE TABLE IF NOT EXISTS advance_requests (
    id SERIAL PRIMARY KEY,
    request_number VARCHAR(20) UNIQUE NOT NULL,
    
    -- Employee Information
    employee_id INT REFERENCES employees(id) NOT NULL,
    employee_name TEXT, -- Denormalized for display
    position TEXT,
    
    -- Reference
    po_number VARCHAR(50),
    
    -- Dates
    request_date DATE DEFAULT CURRENT_DATE,
    
    -- Currency & Amount
    currency currency_code DEFAULT 'LAK',
    total_amount DECIMAL(12,2) DEFAULT 0.00,
    exchange_rate DECIMAL(12,6) DEFAULT 1.00,
    
    -- Main Description
    description TEXT,
    
    -- Approval Chain (3 signatures)
    approved_by_1 INT REFERENCES employees(id),
    approved_date_1 DATE,
    approved_by_2 INT REFERENCES employees(id),
    approved_date_2 DATE,
    approved_by_3 INT REFERENCES employees(id),
    approved_date_3 DATE,
    
    -- Status & Audit
    status advance_request_status DEFAULT 'DRAFT',
    created_by INT REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS advance_request_lines (
    id SERIAL PRIMARY KEY,
    request_id INT REFERENCES advance_requests(id) ON DELETE CASCADE,
    line_number INT NOT NULL,
    
    -- Hierarchy support (for sub-items like A:, B:)
    parent_line_id INT REFERENCES advance_request_lines(id),
    line_type VARCHAR(20) DEFAULT 'MAIN', -- MAIN, SUB_A, SUB_B, etc.
    
    -- Content
    description TEXT NOT NULL,
    amount DECIMAL(12,2) DEFAULT 0.00,
    
    UNIQUE(request_id, line_number)
);

CREATE INDEX idx_advance_requests_status ON advance_requests(status);
CREATE INDEX idx_advance_requests_employee ON advance_requests(employee_id);
CREATE INDEX idx_advance_requests_date ON advance_requests(request_date);

-- =====================================================
-- 1.8 NEW MODULE: ADVANCE VOUCHER
-- =====================================================

CREATE TABLE IF NOT EXISTS advance_vouchers (
    id SERIAL PRIMARY KEY,
    voucher_number VARCHAR(20) UNIQUE NOT NULL,
    
    -- Link to Advance Request
    advance_request_id INT REFERENCES advance_requests(id),
    
    -- Employee Information
    employee_id INT REFERENCES employees(id) NOT NULL,
    employee_name TEXT,
    position TEXT,
    
    -- Reference
    po_number VARCHAR(50),
    
    -- Dates
    voucher_date DATE DEFAULT CURRENT_DATE,
    
    -- Currency
    currency currency_code DEFAULT 'LAK',
    
    -- Amounts
    advance_amount DECIMAL(12,2) DEFAULT 0.00,
    expenditure_amount DECIMAL(12,2) DEFAULT 0.00,
    balance_amount DECIMAL(12,2) GENERATED ALWAYS AS (advance_amount - expenditure_amount) STORED,
    
    -- Approval
    accountant_id INT REFERENCES employees(id),
    accountant_date DATE,
    returned_by_id INT REFERENCES employees(id),
    returned_date DATE,
    
    -- Status & Audit
    status advance_voucher_status DEFAULT 'DRAFT',
    created_by INT REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS advance_voucher_lines (
    id SERIAL PRIMARY KEY,
    voucher_id INT REFERENCES advance_vouchers(id) ON DELETE CASCADE,
    line_number INT NOT NULL,
    
    -- Line Type: ADVANCE, EXPEND, SUB_EXPEND
    line_type VARCHAR(20) NOT NULL,
    
    -- Hierarchy for sub-items under EXPEND
    parent_line_id INT REFERENCES advance_voucher_lines(id),
    
    -- Content
    description TEXT NOT NULL,
    amount DECIMAL(12,2) DEFAULT 0.00,
    
    UNIQUE(voucher_id, line_number)
);

CREATE TABLE IF NOT EXISTS advance_voucher_documents (
    id SERIAL PRIMARY KEY,
    voucher_id INT REFERENCES advance_vouchers(id) ON DELETE CASCADE,
    
    -- Document Type
    document_type voucher_document_type NOT NULL,
    
    -- Status
    has_document BOOLEAN DEFAULT false,
    
    -- Optional file attachment
    document_path TEXT,
    document_name TEXT,
    
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_advance_vouchers_status ON advance_vouchers(status);
CREATE INDEX idx_advance_vouchers_employee ON advance_vouchers(employee_id);
CREATE INDEX idx_advance_vouchers_date ON advance_vouchers(voucher_date);

-- =====================================================
-- 1.9 ADD PAGES FOR NEW MODULES
-- =====================================================

INSERT INTO pages (page_name, route_name, icon) VALUES 
    ('Purchase Requisitions', '/purchase-requisitions', 'file-text'),
    ('Advance Requests', '/advance-requests', 'wallet'),
    ('Advance Vouchers', '/advance-vouchers', 'receipt'),
    ('Brands', '/brands', 'tag'),
    ('Product Groups', '/product-groups', 'layers')
ON CONFLICT (route_name) DO NOTHING;

-- =====================================================
-- 1.10 VIEWS FOR ENHANCED INVENTORY DISPLAY
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
    b.name as brand_name
FROM inventory i
JOIN products p ON i.product_id = p.id
LEFT JOIN warehouses w ON i.warehouse_id = w.id
LEFT JOIN product_categories pc ON p.category_id = pc.id
LEFT JOIN product_groups pg ON p.product_group_id = pg.id
LEFT JOIN brands b ON p.brand_id = b.id;
```

---

## Phase 2: Backend Models (Go)

### 2.1 Update Customer Model

**File: `registration/src/v1/models/customer.go`**

```go
package models

// ============================================
// Customer Models - Updated for DocText.md
// ============================================

type PriceType string

const (
	PriceType1     PriceType = "PRICE_1"
	PriceType2     PriceType = "PRICE_2"
	PriceType3     PriceType = "PRICE_3"
	PriceWholesale PriceType = "WHOLESALE"
	PriceRetail    PriceType = "RETAIL"
)

type CustomerLevel string

const (
	LevelGeneral  CustomerLevel = "GENERAL"
	LevelSilver   CustomerLevel = "SILVER"
	LevelGold     CustomerLevel = "GOLD"
	LevelPlatinum CustomerLevel = "PLATINUM"
	LevelVIP      CustomerLevel = "VIP"
)

type Customer struct {
	ID                 int        `json:"id"`
	CustomerCode       string     `json:"customer_code"`
	Name               string     `json:"name"`
	
	// Contact Information (NEW)
	DefaultContactName string     `json:"default_contact_name,omitempty"`
	FullAddress        string     `json:"full_address,omitempty"`
	DeliveryName       string     `json:"delivery_name,omitempty"`
	Tel                string     `json:"tel,omitempty"`
	Fax                string     `json:"fax,omitempty"`
	Mobile             string     `json:"mobile,omitempty"`
	Email              string     `json:"email,omitempty"`
	WebEmail           string     `json:"web_email,omitempty"`
	
	// Tax & Identity (NEW)
	TaxCode            string     `json:"tax_code,omitempty"`
	PassportNo         string     `json:"passport_no,omitempty"`
	CountryCode        string     `json:"country_code,omitempty"`
	Gender             string     `json:"gender,omitempty"`
	Birthday           CustomDate `json:"birthday,omitempty"`
	
	// Account Settings
	PriceType          PriceType     `json:"price_type"`
	CustomerLevel      CustomerLevel `json:"customer_level"`
	BillDiscount       float64    `json:"bill_discount"`
	CreditLimit        float64    `json:"credit_limit"`
	CreditDays         int        `json:"credit_days"`
	CurrentBalance     float64    `json:"current_balance"`
	PaymentTermsDays   int        `json:"payment_terms_days"`
	Currency           string     `json:"currency"`
	
	// Membership (NEW)
	MemberExpire       CustomDate `json:"member_expire,omitempty"`
	UsePromotion       bool       `json:"use_promotion"`
	UsePMTBill         bool       `json:"use_pmt_bill"`
	CollectPoint       bool       `json:"collect_point"`
	
	// Other (NEW)
	Payee              string     `json:"payee,omitempty"`
	RDBranch           string     `json:"rd_branch,omitempty"`
	Barcode            string     `json:"barcode,omitempty"`
	Notes              string     `json:"notes,omitempty"`
	
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

type CreateCustomerRequest struct {
	CustomerCode       string     `json:"customer_code"`
	Name               string     `json:"name"`
	DefaultContactName string     `json:"default_contact_name,omitempty"`
	FullAddress        string     `json:"full_address,omitempty"`
	DeliveryName       string     `json:"delivery_name,omitempty"`
	Tel                string     `json:"tel,omitempty"`
	Fax                string     `json:"fax,omitempty"`
	Mobile             string     `json:"mobile,omitempty"`
	Email              string     `json:"email,omitempty"`
	WebEmail           string     `json:"web_email,omitempty"`
	TaxCode            string     `json:"tax_code,omitempty"`
	PassportNo         string     `json:"passport_no,omitempty"`
	CountryCode        string     `json:"country_code,omitempty"`
	Gender             string     `json:"gender,omitempty"`
	Birthday           string     `json:"birthday,omitempty"`
	PriceType          PriceType     `json:"price_type"`
	CustomerLevel      CustomerLevel `json:"customer_level"`
	BillDiscount       float64    `json:"bill_discount"`
	CreditLimit        float64    `json:"credit_limit"`
	CreditDays         int        `json:"credit_days"`
	PaymentTermsDays   int        `json:"payment_terms_days"`
	Currency           string     `json:"currency"`
	MemberExpire       string     `json:"member_expire,omitempty"`
	UsePromotion       bool       `json:"use_promotion"`
	UsePMTBill         bool       `json:"use_pmt_bill"`
	CollectPoint       bool       `json:"collect_point"`
	Payee              string     `json:"payee,omitempty"`
	RDBranch           string     `json:"rd_branch,omitempty"`
	Barcode            string     `json:"barcode,omitempty"`
	Notes              string     `json:"notes,omitempty"`
	SalesRepID         *int       `json:"sales_rep_id,omitempty"`
	DefaultWarehouseID *int       `json:"default_warehouse_id,omitempty"`
	TaxExempt          bool       `json:"tax_exempt"`
}
```

### 2.2 Update Vendor Model

**File: `registration/src/v1/models/vendor.go`**

```go
// Add these fields to Vendor struct:
type Vendor struct {
	// ... existing fields ...
	
	// NEW FIELDS
	ContactName     string  `json:"contact_name,omitempty"`
	RDBranch        string  `json:"rd_branch,omitempty"`
	PictureName     string  `json:"picture_name,omitempty"`
	Mobile          string  `json:"mobile,omitempty"`
	Fax             string  `json:"fax,omitempty"`
	TaxCode         string  `json:"tax_code,omitempty"`
	CreditDays      int     `json:"credit_days"`
	DiscountPercent float64 `json:"discount_percent"`
	POCreditMoney   float64 `json:"po_credit_money"`
}
```

### 2.3 Update Purchase Order Model

**File: `registration/src/v1/models/purchase_order.go`**

```go
// Add these fields to PurchaseOrder struct:
type PurchaseOrder struct {
	// ... existing fields ...
	
	// NEW FIELDS - Contact Info on PO
	ContactPerson       string     `json:"contact_person,omitempty"`
	CustomerName        string     `json:"customer_name,omitempty"`
	Address             string     `json:"address,omitempty"`
	Tel                 string     `json:"tel,omitempty"`
	Email               string     `json:"email,omitempty"`
	
	// Document Reference
	DocumentRef         string     `json:"document_ref,omitempty"`
	TermsOfPayment      string     `json:"terms_of_payment,omitempty"`
	Currency            string     `json:"currency"`
	
	// Discount & VAT
	TradeDiscount       float64    `json:"trade_discount"`
	AmountAfterDiscount float64    `json:"amount_after_discount"`
	VATRate             float64    `json:"vat_rate"`
	
	// Signature Chain
	PreparedBy          *int       `json:"prepared_by,omitempty"`
	AuthorizedBy        *int       `json:"authorized_by,omitempty"`
	CheckedBy           *int       `json:"checked_by,omitempty"`
	PreparedDate        CustomDate `json:"prepared_date,omitempty"`
	AuthorizedDate      CustomDate `json:"authorized_date,omitempty"`
	CheckedDate         CustomDate `json:"checked_date,omitempty"`
}

// Update PurchaseOrderLine:
type PurchaseOrderLine struct {
	// ... existing fields ...
	
	// NEW FIELDS
	DiscountPercent  float64 `json:"discount_percent"`
	DiscountAmount   float64 `json:"discount_amount"`
}
```

### 2.4 New Model: Purchase Requisition

**File: `registration/src/v1/models/purchase_requisition.go`**

```go
package models

// ============================================
// Purchase Requisition Status
// ============================================

type PRStatus string

const (
	PRStatusDraft      PRStatus = "DRAFT"
	PRStatusSubmitted  PRStatus = "SUBMITTED"
	PRStatusChecked    PRStatus = "CHECKED"
	PRStatusAuthorized PRStatus = "AUTHORIZED"
	PRStatusConverted  PRStatus = "CONVERTED"
	PRStatusCancelled  PRStatus = "CANCELLED"
)

// ============================================
// Purchase Requisition Models
// ============================================

type PurchaseRequisition struct {
	ID                    int        `json:"id"`
	RequisitionNumber     string     `json:"requisition_number"`
	DocumentNumber        string     `json:"document_number,omitempty"`
	
	// Requester Information
	RequesterID           int        `json:"requester_id"`
	DepartmentID          *int       `json:"department_id,omitempty"`
	Reason                string     `json:"reason,omitempty"`
	
	// Supplier Information
	SupplierID            *int       `json:"supplier_id,omitempty"`
	
	// Dates
	DocumentDate          CustomDate `json:"document_date"`
	RequestDate           CustomDate `json:"request_date"`
	
	// Amounts
	TotalAmount           float64    `json:"total_amount"`
	
	// Remarks
	Remark                string     `json:"remark,omitempty"`
	
	// Approval Chain
	CheckedBy             *int       `json:"checked_by,omitempty"`
	CheckedDate           CustomDate `json:"checked_date,omitempty"`
	AuthorizedBy          *int       `json:"authorized_by,omitempty"`
	AuthorizedDate        CustomDate `json:"authorized_date,omitempty"`
	PurchasingEmployeeID  *int       `json:"purchasing_employee_id,omitempty"`
	PurchasingDate        CustomDate `json:"purchasing_date,omitempty"`
	
	// Linked PO
	ConvertedPOID         *int       `json:"converted_po_id,omitempty"`
	
	// Status & Audit
	Status                PRStatus   `json:"status"`
	CreatedBy             int        `json:"created_by"`
	CreatedAt             CustomDate `json:"created_at"`
	UpdatedAt             CustomDate `json:"updated_at"`
}

type PurchaseRequisitionLine struct {
	ID                int     `json:"id"`
	RequisitionID     int     `json:"requisition_id"`
	LineNumber        int     `json:"line_number"`
	ProductID         *int    `json:"product_id,omitempty"`
	ProductCode       string  `json:"product_code,omitempty"`
	Description       string  `json:"description"`
	StockBalance      float64 `json:"stock_balance"`
	QuantityRequested float64 `json:"quantity_requested"`
	UnitOfMeasure     string  `json:"unit_of_measure"`
	UnitPrice         float64 `json:"unit_price"`
	Amount            float64 `json:"amount"`
}

type PurchaseRequisitionWithDetails struct {
	Requisition        PurchaseRequisition       `json:"requisition"`
	Lines              []PurchaseRequisitionLine `json:"lines"`
	RequesterName      string                    `json:"requester_name"`
	DepartmentName     string                    `json:"department_name,omitempty"`
	SupplierName       string                    `json:"supplier_name,omitempty"`
	CheckedByName      string                    `json:"checked_by_name,omitempty"`
	AuthorizedByName   string                    `json:"authorized_by_name,omitempty"`
	PurchasingEmpName  string                    `json:"purchasing_emp_name,omitempty"`
}

// ============================================
// Request/Response Types
// ============================================

type CreatePurchaseRequisitionRequest struct {
	DepartmentID   *int                              `json:"department_id,omitempty"`
	Reason         string                            `json:"reason,omitempty"`
	SupplierID     *int                              `json:"supplier_id,omitempty"`
	RequestDate    string                            `json:"request_date,omitempty"`
	Remark         string                            `json:"remark,omitempty"`
	Lines          []CreatePRLineRequest             `json:"lines"`
}

type CreatePRLineRequest struct {
	ProductID         *int    `json:"product_id,omitempty"`
	ProductCode       string  `json:"product_code,omitempty"`
	Description       string  `json:"description"`
	StockBalance      float64 `json:"stock_balance"`
	QuantityRequested float64 `json:"quantity_requested"`
	UnitOfMeasure     string  `json:"unit_of_measure"`
	UnitPrice         float64 `json:"unit_price"`
}

type SubmitPRRequest struct {
	CheckedBy    *int   `json:"checked_by,omitempty"`
	AuthorizedBy *int   `json:"authorized_by,omitempty"`
}

type ConvertPRToPORequest struct {
	WarehouseID  int    `json:"warehouse_id"`
	ExpectedDate string `json:"expected_date,omitempty"`
	BuyerID      *int   `json:"buyer_id,omitempty"`
}

// ============================================
// Validation
// ============================================

func ValidatePurchaseRequisition(v *Validator, req *CreatePurchaseRequisitionRequest) {
	v.Check(len(req.Lines) > 0, "lines", "At least one line is required")
	for _, line := range req.Lines {
		v.Check(line.Description != "" || line.ProductID != nil, "lines", "Description or Product ID is required")
		v.Check(line.QuantityRequested > 0, "lines", "Quantity must be positive")
	}
}
```

### 2.5 New Model: Advance Request

**File: `registration/src/v1/models/advance_request.go`**

```go
package models

// ============================================
// Advance Request Status
// ============================================

type AdvanceRequestStatus string

const (
	AdvReqStatusDraft     AdvanceRequestStatus = "DRAFT"
	AdvReqStatusSubmitted AdvanceRequestStatus = "SUBMITTED"
	AdvReqStatusApproved  AdvanceRequestStatus = "APPROVED"
	AdvReqStatusRejected  AdvanceRequestStatus = "REJECTED"
	AdvReqStatusCleared   AdvanceRequestStatus = "CLEARED"
)

type CurrencyCode string

const (
	CurrencyLAK CurrencyCode = "LAK"
	CurrencyTHB CurrencyCode = "THB"
	CurrencyUSD CurrencyCode = "USD"
)

// ============================================
// Advance Request Models
// ============================================

type AdvanceRequest struct {
	ID            int                  `json:"id"`
	RequestNumber string               `json:"request_number"`
	
	// Employee Info
	EmployeeID    int                  `json:"employee_id"`
	EmployeeName  string               `json:"employee_name,omitempty"`
	Position      string               `json:"position,omitempty"`
	
	// Reference
	PONumber      string               `json:"po_number,omitempty"`
	
	// Date
	RequestDate   CustomDate           `json:"request_date"`
	
	// Currency & Amount
	Currency      CurrencyCode         `json:"currency"`
	TotalAmount   float64              `json:"total_amount"`
	ExchangeRate  float64              `json:"exchange_rate"`
	
	// Description
	Description   string               `json:"description,omitempty"`
	
	// Approval Chain (3 signatures)
	ApprovedBy1   *int                 `json:"approved_by_1,omitempty"`
	ApprovedDate1 CustomDate           `json:"approved_date_1,omitempty"`
	ApprovedBy2   *int                 `json:"approved_by_2,omitempty"`
	ApprovedDate2 CustomDate           `json:"approved_date_2,omitempty"`
	ApprovedBy3   *int                 `json:"approved_by_3,omitempty"`
	ApprovedDate3 CustomDate           `json:"approved_date_3,omitempty"`
	
	// Status & Audit
	Status        AdvanceRequestStatus `json:"status"`
	CreatedBy     int                  `json:"created_by"`
	CreatedAt     CustomDate           `json:"created_at"`
	UpdatedAt     CustomDate           `json:"updated_at"`
}

type AdvanceRequestLine struct {
	ID           int     `json:"id"`
	RequestID    int     `json:"request_id"`
	LineNumber   int     `json:"line_number"`
	ParentLineID *int    `json:"parent_line_id,omitempty"`
	LineType     string  `json:"line_type"` // MAIN, SUB_A, SUB_B
	Description  string  `json:"description"`
	Amount       float64 `json:"amount"`
}

type AdvanceRequestWithDetails struct {
	Request        AdvanceRequest       `json:"request"`
	Lines          []AdvanceRequestLine `json:"lines"`
	ApprovedBy1Name string              `json:"approved_by_1_name,omitempty"`
	ApprovedBy2Name string              `json:"approved_by_2_name,omitempty"`
	ApprovedBy3Name string              `json:"approved_by_3_name,omitempty"`
}

// Request/Response Types
type CreateAdvanceRequestRequest struct {
	EmployeeID   int                       `json:"employee_id"`
	Position     string                    `json:"position,omitempty"`
	PONumber     string                    `json:"po_number,omitempty"`
	RequestDate  string                    `json:"request_date,omitempty"`
	Currency     CurrencyCode              `json:"currency"`
	ExchangeRate float64                   `json:"exchange_rate"`
	Description  string                    `json:"description,omitempty"`
	Lines        []CreateAdvanceLineRequest `json:"lines"`
}

type CreateAdvanceLineRequest struct {
	ParentLineID *int    `json:"parent_line_id,omitempty"`
	LineType     string  `json:"line_type"`
	Description  string  `json:"description"`
	Amount       float64 `json:"amount"`
}
```

### 2.6 New Model: Advance Voucher

**File: `registration/src/v1/models/advance_voucher.go`**

```go
package models

// ============================================
// Advance Voucher Status
// ============================================

type AdvanceVoucherStatus string

const (
	AdvVoucherStatusDraft     AdvanceVoucherStatus = "DRAFT"
	AdvVoucherStatusSubmitted AdvanceVoucherStatus = "SUBMITTED"
	AdvVoucherStatusApproved  AdvanceVoucherStatus = "APPROVED"
	AdvVoucherStatusRejected  AdvanceVoucherStatus = "REJECTED"
	AdvVoucherStatusCompleted AdvanceVoucherStatus = "COMPLETED"
)

type VoucherDocumentType string

const (
	DocTypeReceipt        VoucherDocumentType = "RECEIPT"
	DocTypeInvoice        VoucherDocumentType = "INVOICE"
	DocTypeBankTransfer   VoucherDocumentType = "BANK_TRANSFER"
	DocTypePOPR           VoucherDocumentType = "PO_PR"
	DocTypeTransportation VoucherDocumentType = "TRANSPORTATION"
)

// ============================================
// Advance Voucher Models
// ============================================

type AdvanceVoucher struct {
	ID                 int                  `json:"id"`
	VoucherNumber      string               `json:"voucher_number"`
	
	// Link to Request
	AdvanceRequestID   *int                 `json:"advance_request_id,omitempty"`
	
	// Employee Info
	EmployeeID         int                  `json:"employee_id"`
	EmployeeName       string               `json:"employee_name,omitempty"`
	Position           string               `json:"position,omitempty"`
	
	// Reference
	PONumber           string               `json:"po_number,omitempty"`
	
	// Date
	VoucherDate        CustomDate           `json:"voucher_date"`
	
	// Currency & Amounts
	Currency           CurrencyCode         `json:"currency"`
	AdvanceAmount      float64              `json:"advance_amount"`
	ExpenditureAmount  float64              `json:"expenditure_amount"`
	BalanceAmount      float64              `json:"balance_amount"` // Computed
	
	// Signatures
	AccountantID       *int                 `json:"accountant_id,omitempty"`
	AccountantDate     CustomDate           `json:"accountant_date,omitempty"`
	ReturnedByID       *int                 `json:"returned_by_id,omitempty"`
	ReturnedDate       CustomDate           `json:"returned_date,omitempty"`
	
	// Status & Audit
	Status             AdvanceVoucherStatus `json:"status"`
	CreatedBy          int                  `json:"created_by"`
	CreatedAt          CustomDate           `json:"created_at"`
	UpdatedAt          CustomDate           `json:"updated_at"`
}

type AdvanceVoucherLine struct {
	ID           int     `json:"id"`
	VoucherID    int     `json:"voucher_id"`
	LineNumber   int     `json:"line_number"`
	LineType     string  `json:"line_type"` // ADVANCE, EXPEND, SUB_EXPEND
	ParentLineID *int    `json:"parent_line_id,omitempty"`
	Description  string  `json:"description"`
	Amount       float64 `json:"amount"`
}

type AdvanceVoucherDocument struct {
	ID           int                 `json:"id"`
	VoucherID    int                 `json:"voucher_id"`
	DocumentType VoucherDocumentType `json:"document_type"`
	HasDocument  bool                `json:"has_document"`
	DocumentPath string              `json:"document_path,omitempty"`
	DocumentName string              `json:"document_name,omitempty"`
}

type AdvanceVoucherWithDetails struct {
	Voucher          AdvanceVoucher           `json:"voucher"`
	Lines            []AdvanceVoucherLine     `json:"lines"`
	Documents        []AdvanceVoucherDocument `json:"documents"`
	AccountantName   string                   `json:"accountant_name,omitempty"`
	ReturnedByName   string                   `json:"returned_by_name,omitempty"`
}
```

---

## Phase 3: Backend Services & Routes

### 3.1 Services to Create

| Service File | Module | Key Functions |
|--------------|--------|---------------|
| `purchase_requisition/purchase_requisition.go` | Purchase Requisition | CRUD, Submit, Approve, ConvertToPO |
| `advance_request/advance_request.go` | Advance Request | CRUD, Submit, Approve (3-level) |
| `advance_voucher/advance_voucher.go` | Advance Voucher | CRUD, Submit, Approve, Document Upload |
| `brand/brand.go` | Brands | CRUD |
| `product_group/product_group.go` | Product Groups | CRUD |

### 3.2 Service Function Signatures

**Purchase Requisition Service:**
```go
func (s *PurchaseRequisitionService) List(filters PRListFilters) ([]PurchaseRequisitionWithDetails, *models.Pagination, error)
func (s *PurchaseRequisitionService) GetByID(id int) (*PurchaseRequisitionWithDetails, error)
func (s *PurchaseRequisitionService) Create(req CreatePurchaseRequisitionRequest, userID int) (*PurchaseRequisition, error)
func (s *PurchaseRequisitionService) Update(id int, req UpdatePRRequest) (*PurchaseRequisition, error)
func (s *PurchaseRequisitionService) Delete(id int) error
func (s *PurchaseRequisitionService) Submit(id int, req SubmitPRRequest) error
func (s *PurchaseRequisitionService) Approve(id int, approverID int, action string) error
func (s *PurchaseRequisitionService) ConvertToPO(id int, req ConvertPRToPORequest, userID int) (*PurchaseOrder, error)
```

### 3.3 Routes to Add

**File: `registration/src/v1/routes/routes.go`**

```go
// Purchase Requisitions
prGroup := v1.Group("/purchase-requisitions")
{
    prGroup.GET("/list", prHandler.List)
    prGroup.GET("/get/:id", prHandler.GetByID)
    prGroup.POST("/create", prHandler.Create)
    prGroup.PUT("/update/:id", prHandler.Update)
    prGroup.DELETE("/delete/:id", prHandler.Delete)
    prGroup.POST("/submit/:id", prHandler.Submit)
    prGroup.POST("/approve/:id", prHandler.Approve)
    prGroup.POST("/convert-to-po/:id", prHandler.ConvertToPO)
}

// Advance Requests
arGroup := v1.Group("/advance-requests")
{
    arGroup.GET("/list", arHandler.List)
    arGroup.GET("/get/:id", arHandler.GetByID)
    arGroup.POST("/create", arHandler.Create)
    arGroup.PUT("/update/:id", arHandler.Update)
    arGroup.DELETE("/delete/:id", arHandler.Delete)
    arGroup.POST("/submit/:id", arHandler.Submit)
    arGroup.POST("/approve/:id", arHandler.Approve)
}

// Advance Vouchers
avGroup := v1.Group("/advance-vouchers")
{
    avGroup.GET("/list", avHandler.List)
    avGroup.GET("/get/:id", avHandler.GetByID)
    avGroup.POST("/create", avHandler.Create)
    avGroup.PUT("/update/:id", avHandler.Update)
    avGroup.DELETE("/delete/:id", avHandler.Delete)
    avGroup.POST("/submit/:id", avHandler.Submit)
    avGroup.POST("/approve/:id", avHandler.Approve)
    avGroup.POST("/:id/documents", avHandler.UploadDocument)
    avGroup.PUT("/:id/documents/:docId", avHandler.UpdateDocument)
}

// Brands
brandsGroup := v1.Group("/brands")
{
    brandsGroup.GET("/list", brandHandler.List)
    brandsGroup.GET("/get/:id", brandHandler.GetByID)
    brandsGroup.POST("/create", brandHandler.Create)
    brandsGroup.PUT("/update/:id", brandHandler.Update)
    brandsGroup.DELETE("/delete/:id", brandHandler.Delete)
}

// Product Groups
pgGroup := v1.Group("/product-groups")
{
    pgGroup.GET("/list", pgHandler.List)
    pgGroup.GET("/get/:id", pgHandler.GetByID)
    pgGroup.POST("/create", pgHandler.Create)
    pgGroup.PUT("/update/:id", pgHandler.Update)
    pgGroup.DELETE("/delete/:id", pgHandler.Delete)
}
```

---

## Phase 4: Frontend Services

### 4.1 Update masterDataService.ts

**File: `foodhive-erp-frontend/client/src/services/masterDataService.ts`**

Add extended Customer interface:

```typescript
export interface Customer {
  id: number;
  customer_code: string;
  name: string;
  
  // Contact Information
  default_contact_name?: string;
  full_address?: string;
  delivery_name?: string;
  tel?: string;
  fax?: string;
  mobile?: string;
  email?: string;
  web_email?: string;
  
  // Tax & Identity
  tax_code?: string;
  passport_no?: string;
  country_code?: string;
  gender?: 'MALE' | 'FEMALE';
  birthday?: string;
  
  // Account Settings
  price_type?: 'PRICE_1' | 'PRICE_2' | 'PRICE_3' | 'WHOLESALE' | 'RETAIL';
  customer_level?: 'GENERAL' | 'SILVER' | 'GOLD' | 'PLATINUM' | 'VIP';
  bill_discount?: number;
  credit_limit: number;
  credit_days?: number;
  payment_terms_days: number;
  
  // Membership
  member_expire?: string;
  use_promotion?: boolean;
  use_pmt_bill?: boolean;
  collect_point?: boolean;
  
  // Other
  payee?: string;
  rd_branch?: string;
  barcode?: string;
  notes?: string;
  
  sales_rep_id?: number;
  is_active: boolean;
}

export interface Vendor {
  id: number;
  vendor_code: string;
  name: string;
  
  // Contact
  contact_name?: string;
  phone?: string;
  mobile?: string;
  fax?: string;
  email?: string;
  
  // Address
  address_line1?: string;
  address_line2?: string;
  city?: string;
  state?: string;
  postal_code?: string;
  country?: string;
  
  // Tax & Financial
  tax_code?: string;
  rd_branch?: string;
  credit_days?: number;
  discount_percent?: number;
  po_credit_money?: number;
  payment_terms_days: number;
  currency?: string;
  
  // Other
  picture_name?: string;
  
  is_active: boolean;
}
```

### 4.2 New Service: purchaseRequisitionService.ts

**File: `foodhive-erp-frontend/client/src/services/purchaseRequisitionService.ts`**

```typescript
import api from '@/lib/api';

export interface PurchaseRequisition {
  id: number;
  requisition_number: string;
  document_number?: string;
  requester_id: number;
  requester_name?: string;
  department_id?: number;
  department_name?: string;
  reason?: string;
  supplier_id?: number;
  supplier_name?: string;
  document_date: string;
  request_date: string;
  total_amount: number;
  remark?: string;
  status: 'DRAFT' | 'SUBMITTED' | 'CHECKED' | 'AUTHORIZED' | 'CONVERTED' | 'CANCELLED';
  lines?: PurchaseRequisitionLine[];
}

export interface PurchaseRequisitionLine {
  id: number;
  line_number: number;
  product_id?: number;
  product_code?: string;
  description: string;
  stock_balance: number;
  quantity_requested: number;
  unit_of_measure: string;
  unit_price: number;
  amount: number;
}

export const purchaseRequisitionService = {
  getList: async (params?: any) => {
    const response = await api.get('/purchase-requisitions/list', { params });
    return response.data?.data || response.data || [];
  },

  getById: async (id: string) => {
    const response = await api.get(`/purchase-requisitions/get/${id}`);
    return response.data?.data || response.data;
  },

  create: async (data: Partial<PurchaseRequisition>) => {
    const response = await api.post('/purchase-requisitions/create', data);
    return response.data;
  },

  update: async (id: string, data: Partial<PurchaseRequisition>) => {
    const response = await api.put(`/purchase-requisitions/update/${id}`, data);
    return response.data;
  },

  delete: async (id: string) => {
    const response = await api.delete(`/purchase-requisitions/delete/${id}`);
    return response.data;
  },

  submit: async (id: string) => {
    const response = await api.post(`/purchase-requisitions/submit/${id}`);
    return response.data;
  },

  approve: async (id: string, action: 'check' | 'authorize' | 'purchase') => {
    const response = await api.post(`/purchase-requisitions/approve/${id}`, { action });
    return response.data;
  },

  convertToPO: async (id: string, data: { warehouse_id: number; expected_date?: string; buyer_id?: number }) => {
    const response = await api.post(`/purchase-requisitions/convert-to-po/${id}`, data);
    return response.data;
  }
};
```

### 4.3 New Service: advanceService.ts

**File: `foodhive-erp-frontend/client/src/services/advanceService.ts`**

```typescript
import api from '@/lib/api';

export type Currency = 'LAK' | 'THB' | 'USD';

export interface AdvanceRequest {
  id: number;
  request_number: string;
  employee_id: number;
  employee_name?: string;
  position?: string;
  po_number?: string;
  request_date: string;
  currency: Currency;
  total_amount: number;
  exchange_rate: number;
  description?: string;
  status: 'DRAFT' | 'SUBMITTED' | 'APPROVED' | 'REJECTED' | 'CLEARED';
  lines?: AdvanceRequestLine[];
}

export interface AdvanceRequestLine {
  id: number;
  line_number: number;
  parent_line_id?: number;
  line_type: string;
  description: string;
  amount: number;
}

export interface AdvanceVoucher {
  id: number;
  voucher_number: string;
  advance_request_id?: number;
  employee_id: number;
  employee_name?: string;
  position?: string;
  po_number?: string;
  voucher_date: string;
  currency: Currency;
  advance_amount: number;
  expenditure_amount: number;
  balance_amount: number;
  status: 'DRAFT' | 'SUBMITTED' | 'APPROVED' | 'REJECTED' | 'COMPLETED';
  lines?: AdvanceVoucherLine[];
  documents?: AdvanceVoucherDocument[];
}

export interface AdvanceVoucherLine {
  id: number;
  line_number: number;
  line_type: 'ADVANCE' | 'EXPEND' | 'SUB_EXPEND';
  parent_line_id?: number;
  description: string;
  amount: number;
}

export interface AdvanceVoucherDocument {
  id: number;
  document_type: 'RECEIPT' | 'INVOICE' | 'BANK_TRANSFER' | 'PO_PR' | 'TRANSPORTATION';
  has_document: boolean;
  document_path?: string;
  document_name?: string;
}

export const advanceService = {
  // Advance Requests
  getRequests: async (params?: any) => {
    const response = await api.get('/advance-requests/list', { params });
    return response.data?.data || response.data || [];
  },

  getRequest: async (id: string) => {
    const response = await api.get(`/advance-requests/get/${id}`);
    return response.data?.data || response.data;
  },

  createRequest: async (data: Partial<AdvanceRequest>) => {
    const response = await api.post('/advance-requests/create', data);
    return response.data;
  },

  updateRequest: async (id: string, data: Partial<AdvanceRequest>) => {
    const response = await api.put(`/advance-requests/update/${id}`, data);
    return response.data;
  },

  deleteRequest: async (id: string) => {
    const response = await api.delete(`/advance-requests/delete/${id}`);
    return response.data;
  },

  submitRequest: async (id: string) => {
    const response = await api.post(`/advance-requests/submit/${id}`);
    return response.data;
  },

  approveRequest: async (id: string, level: 1 | 2 | 3) => {
    const response = await api.post(`/advance-requests/approve/${id}`, { level });
    return response.data;
  },

  // Advance Vouchers
  getVouchers: async (params?: any) => {
    const response = await api.get('/advance-vouchers/list', { params });
    return response.data?.data || response.data || [];
  },

  getVoucher: async (id: string) => {
    const response = await api.get(`/advance-vouchers/get/${id}`);
    return response.data?.data || response.data;
  },

  createVoucher: async (data: Partial<AdvanceVoucher>) => {
    const response = await api.post('/advance-vouchers/create', data);
    return response.data;
  },

  updateVoucher: async (id: string, data: Partial<AdvanceVoucher>) => {
    const response = await api.put(`/advance-vouchers/update/${id}`, data);
    return response.data;
  },

  deleteVoucher: async (id: string) => {
    const response = await api.delete(`/advance-vouchers/delete/${id}`);
    return response.data;
  },

  submitVoucher: async (id: string) => {
    const response = await api.post(`/advance-vouchers/submit/${id}`);
    return response.data;
  },

  approveVoucher: async (id: string, role: 'accountant' | 'returned_by') => {
    const response = await api.post(`/advance-vouchers/approve/${id}`, { role });
    return response.data;
  },

  updateDocument: async (voucherId: string, docType: string, hasDocument: boolean, file?: File) => {
    const formData = new FormData();
    formData.append('document_type', docType);
    formData.append('has_document', String(hasDocument));
    if (file) formData.append('file', file);
    
    const response = await api.post(`/advance-vouchers/${voucherId}/documents`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    });
    return response.data;
  }
};
```

---

## Phase 5: Frontend Components

### 5.1 Updated Customer Form

**File: `foodhive-erp-frontend/client/src/pages/customers/CustomerForm.tsx`**

The form needs to be expanded to include all new fields organized into tabs:

1. **Basic Information Tab**
   - ID, Customer Code, Name
   - Default Contact Name, Full Address
   - Delivery Name

2. **Contact Information Tab**
   - Tel, Fax, Mobile
   - Email, Web Email
   - Birthday

3. **Identity & Tax Tab**
   - Tax Code, Passport No
   - Country Code, Gender

4. **Account Settings Tab**
   - Price Type (dropdown)
   - Customer Level (dropdown)
   - Bill Discount, Credit Limit, Credit Days
   - Payment Terms, Currency

5. **Membership Tab**
   - Member Expire Date
   - Use Promotion (checkbox)
   - Use PMT Bill (checkbox)
   - Collect Point (checkbox)

6. **Other Tab**
   - Payee, RD Branch
   - Barcode, Notes

### 5.2 New Pages to Create

| Page | Path | Description |
|------|------|-------------|
| `PurchaseRequisitionList.tsx` | `/purchase-requisitions` | List all PRs with filters |
| `PurchaseRequisitionForm.tsx` | `/purchase-requisitions/new` | Create/Edit PR |
| `PurchaseRequisitionDetail.tsx` | `/purchase-requisitions/:id` | View PR with approval workflow |
| `AdvanceRequestList.tsx` | `/advance-requests` | List all advance requests |
| `AdvanceRequestForm.tsx` | `/advance-requests/new` | Create/Edit advance request |
| `AdvanceRequestDetail.tsx` | `/advance-requests/:id` | View with 3-level approval |
| `AdvanceVoucherList.tsx` | `/advance-vouchers` | List all vouchers |
| `AdvanceVoucherForm.tsx` | `/advance-vouchers/new` | Create/Edit voucher |
| `AdvanceVoucherDetail.tsx` | `/advance-vouchers/:id` | View with document checklist |
| `BrandList.tsx` | `/brands` | Manage brands |
| `ProductGroupList.tsx` | `/product-groups` | Manage product groups |

### 5.3 Enhanced Inventory Page

**Updates to `foodhive-erp-frontend/client/src/pages/inventory/InventoryList.tsx`**

Add new features:
- **Display Columns**: Price 1, Price 2, VAT%
- **Filter Options**: Product Group, Brand
- **Action Buttons**: Edit Sell Price, Print Barcode, Group View
- **Options**: Set Buy Price toggle, Show Num Order, Expand Group

---

## Implementation Priority & Timeline

### Phase 1: High Priority (Week 1-2)
1. ✅ Database schema updates (1 day)
2. ✅ Customer model & service updates (2 days)
3. ✅ Vendor model & service updates (1 day)
4. ✅ Purchase Order updates (1 day)
5. ✅ Frontend Customer form update (2 days)
6. ✅ Frontend Vendor form update (1 day)

### Phase 2: New Modules (Week 3-4)
1. Purchase Requisition - Full module (3 days)
2. Advance Request - Full module (2 days)
3. Advance Voucher - Full module (2 days)

### Phase 3: Enhancements (Week 5)
1. Brand & Product Group modules (1 day)
2. Product pricing fields (1 day)
3. Enhanced Inventory display (2 days)

### Phase 4: Testing & Polish (Week 6)
1. Integration testing
2. UI/UX refinements
3. Documentation updates

---

## File Changes Summary

### Database
- [ ] `sql/006_requirements_gap_updates.sql` - New file

### Backend Models
- [ ] `models/customer.go` - Update
- [ ] `models/vendor.go` - Update
- [ ] `models/purchase_order.go` - Update
- [ ] `models/product.go` - Update
- [ ] `models/purchase_requisition.go` - New file
- [ ] `models/advance_request.go` - New file
- [ ] `models/advance_voucher.go` - New file

### Backend Services
- [ ] `services/customer/customer.go` - Update
- [ ] `services/vendor/vendor.go` - Update
- [ ] `services/purchase_order/purchase_order.go` - Update
- [ ] `services/purchase_requisition/purchase_requisition.go` - New file
- [ ] `services/advance_request/advance_request.go` - New file
- [ ] `services/advance_voucher/advance_voucher.go` - New file
- [ ] `services/brand/brand.go` - New file
- [ ] `services/product_group/product_group.go` - New file

### Backend Handlers & Routes
- [ ] `handlers/purchase_requisition.go` - New file
- [ ] `handlers/advance_request.go` - New file
- [ ] `handlers/advance_voucher.go` - New file
- [ ] `routes/routes.go` - Update

### Frontend Services
- [ ] `services/masterDataService.ts` - Update interfaces
- [ ] `services/purchasingService.ts` - Update
- [ ] `services/purchaseRequisitionService.ts` - New file
- [ ] `services/advanceService.ts` - New file

### Frontend Pages
- [ ] `pages/customers/CustomerList.tsx` - Update form
- [ ] `pages/vendors/VendorList.tsx` - Update form
- [ ] `pages/purchasing/PurchaseOrderList.tsx` - Update
- [ ] `pages/purchasing/PurchaseRequisitionList.tsx` - New file
- [ ] `pages/purchasing/PurchaseRequisitionForm.tsx` - New file
- [ ] `pages/finance/AdvanceRequestList.tsx` - New file
- [ ] `pages/finance/AdvanceRequestForm.tsx` - New file
- [ ] `pages/finance/AdvanceVoucherList.tsx` - New file
- [ ] `pages/finance/AdvanceVoucherForm.tsx` - New file
- [ ] `pages/inventory/InventoryList.tsx` - Update

### Navigation & Routing
- [ ] `App.tsx` - Add new routes
- [ ] `components/layout/Sidebar.tsx` - Add new menu items

---

*Document Version: 1.0*
*Last Updated: January 2026*
*Author: FoodHive Development Team*
