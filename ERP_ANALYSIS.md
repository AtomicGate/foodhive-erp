# 📋 ERP System Requirements Analysis

## UMS to ERP Transformation Guide

**Document Version:** 1.0  
**Date:** January 9, 2026  
**Based on:** CamScanner 01-09-2026 18.34.pdf (Client ERP Requirements)

---

## 📑 Table of Contents

1. [Executive Summary](#executive-summary)
2. [Current UMS Capabilities](#current-ums-capabilities)
3. [ERP Requirements Overview](#erp-requirements-overview)
4. [Reusability Analysis](#reusability-analysis)
5. [Modules to Build from Scratch](#modules-to-build-from-scratch)
6. [Database Schema Design](#database-schema-design)
7. [Implementation Roadmap](#implementation-roadmap)
8. [Technical Recommendations](#technical-recommendations)

---

## Executive Summary

This document analyzes the client's ERP requirements against our existing University Management System (UMS) to determine:

- **What can be reused** with modifications
- **What must be built** from scratch
- **Estimated effort** for each module
- **Recommended implementation** approach

### Quick Stats

| Category | Percentage |
|----------|------------|
| Fully Reusable | ~25% |
| Partially Reusable | ~20% |
| Build from Scratch | ~55% |

---

## Current UMS Capabilities

### Technology Stack

| Component | Technology |
|-----------|------------|
| Language | Go 1.23 |
| Web Framework | chi/v5 |
| Database | PostgreSQL (pgx/v5) |
| Object Storage | MinIO |
| Authentication | JWT (golang-jwt/v5) |
| API Docs | Swagger (swaggo) |
| Config | godotenv + Netflix/go-env |

### Existing Modules

#### ✅ Core Infrastructure
- PostgreSQL database layer with transactions
- JWT authentication service
- Role-based access control (RBAC)
- MinIO file storage integration
- Middleware-based service injection
- Chi router with versioned API

#### ✅ User Management
- Employee CRUD operations
- Employee details (gender, job title, passport, etc.)
- Employee finances (salary, allowances, bonuses)
- Contract management (YEARLY/FULL)
- Role and permission management
- Address management

#### ✅ Financial Tracking
- Income tracking (INTERNAL/EXTERNAL)
- Expense management with types
- Cash boxes with detailed transactions
- Student payments with installments

#### ✅ Document Management
- Correspondence system
- Document upload/download (MinIO)
- Recipient tracking with status
- Read tracking per recipient

#### ✅ Academic
- Student registration (temp_students)
- Department management
- Study types and year management
- Admin orders with file attachments
- Archive management

#### ✅ Payroll (Basic)
- Payroll types: LECTURER, PROFESSOR, EMPLOYEE
- Department-based payroll
- Insurance and allowance flags

---

## ERP Requirements Overview

Based on the client document, the following modules are required:

### 1. Sales/Order Entry
- Rapid entry mode, order guide mode, traditional entry
- Multiple order types (standard, advance, pre-paid, on hold, quotes, credit memo, pick up)
- Multiple warehouse sales
- Auto-routing based on customer's route
- Customer order guides with buying stats
- Inventory on hand, allocated, on order display
- Selling below cost notification
- Integrated margin management
- Auto pricing based on customer setup
- Catch weight integration
- Lot and Best Before/Expiry Date lookup
- Lost sales reporting
- Special order requests to buyer
- Customer account drill-down
- EDI, e-commerce integration
- Commission reporting

### 2. Picking & Routing
- Route setup and maintenance
- Zone picking or order picking by date/route
- Master pick reporting
- Mobile units, pick tickets, labels
- Stop re-numbering and order shuffling
- Batch release and printing
- Lot tracking, catch weight, barcode capture
- Suggested picking based on rotation
- Re-packing/processing integration
- Shipping labels
- Pallet functions
- EAN UCC barcode standards

### 3. Customer & AR Management
- Customer data with security permissions
- Multiple currency AR
- Multiple ship-to locations
- Flexible statements
- Default warehouse/shipping assignment
- Customer product code links
- Rep integration for commissions
- Bill To / Ship To grouping
- Customer order guides from history
- Credit management profiling
- Extensive AR reporting

### 4. Pricing & Cost Management
- 5-level price hierarchy (contract, price level, custom, promotional, order)
- Future effective dates for price/cost
- Price lists by effective date
- Mass price maintenance
- Price override with security
- Rebate tracking integration
- Multiple costing methods (average, last, landed, all-in, market, default vendor)
- Cost error correction with backflush

### 5. Inventory Control
- Product group codes (posting, category, sub-category)
- UCC EAN Barcodes
- Catch weight
- Multiple units of measure
- Broken case handling
- Portion costing
- COOL (Country of Origin Labeling)
- Real-time inventory inquiry
- Grade/HACCP/QA comments
- Lb/Kg conversion
- Age of inventory reporting
- Multiple warehouse tracking
- Lot and date tracking
- One up & one down traceability
- Product reservations

### 6. Purchasing & Receiving
- Buyer/supplier/stocking level integration
- Reorder based on min/max or average sales
- Quick PO entry
- Multiple warehouse purchasing
- ETA scheduling with lead times
- Inventory and supplier order guide lookup
- LB/KG conversion
- Special order requests view
- Short shipment reporting
- Buyer's tools - auto generate suggested buying
- WMS integration

### 7. Accounts Payable
- Vendor account history
- Vendor discount tables
- Automatic check runs by date/vendor type
- Inventory adjustments from AP invoice
- PO and vendor settings integration
- Vendor templates for costing matrix
- Vendor lists and labels
- Invoice history and activity reporting
- Snapshot aged payable reports

### 8. Bank & Reconciliation
- Check and receipt activity management
- Automatic payable runs by date
- Manual check entry
- Check reprint, void, view
- Open invoices by customer
- Receipt entry and posting
- Bank reconciliation by account and date

### 9. General Ledger
- Multi-currency at transaction level
- Multi branch/division/department
- Integrated financial reporting
- Comparison reporting (actual vs budget, prior period, prior year)
- Recurring entry management
- Auto reversal tool
- Full audit trail
- Drill-down capability
- Snapshot "as at any date"
- Flexible fiscal year structure

### 10. Warehouse Management System (WMS)
- Receiving
- Put-Away
- Replenishment
- Skid Tracking
- Picking
- Warehouse Transfers
- Cycle Counts
- Labeling
- Skid Maintenance
- Physical and Cycle counts
- Re-Work
- Disposals

### 11. Payroll
- GL integration
- Bank Checks integration
- AR integration

---

## Reusability Analysis

### ✅ 100% Reusable - Core Infrastructure

| Component | Location | Notes |
|-----------|----------|-------|
| Database Layer | `core/postgres/db.go` | Ready to use, supports transactions |
| JWT Service | `core/jwt/jwt.go` | Minor adjustments for ERP roles |
| Auth Service | `core/auth/auth.go` | Add new ERP pages/permissions |
| Storage Service | `core/storage/storage.go` | Add new buckets for ERP docs |
| Middleware Pattern | `registration/src/v1/middlewares/` | Replicate for new modules |
| Router | `registration/src/v1/v1.go` | Extend with new routes |
| Swagger | `registration/docs/` | Add new endpoints |
| Helper Functions | `core/utils/helpers/` | JSON responses, pagination |
| Environment Config | `core/utils/env/` | Add new config options |

### 🔄 70-80% Reusable - Payroll Module

**Current Structure:**
```go
type Payroll struct {
    ID           int
    CreatedBy    int
    DepartmentID int
    CreatedAt    CustomDate
    Date         string
    Note         string
    Insurance    bool
    ActiveOne    bool
    ActiveTwo    bool
    ActiveThree  bool
    ActiveFour   bool
    Mode         bool
    Type         Payroll_type // LECTURER, PROFESSOR, EMPLOYEE
}
```

**Modifications Required:**
- Add GL account integration
- Add bank check integration
- Add AR integration
- Add pay period fields
- Add gross/net pay calculations
- Add deduction tracking

**New Fields to Add:**
```sql
ALTER TABLE payroll ADD COLUMN gl_account_id INT;
ALTER TABLE payroll ADD COLUMN check_id INT;
ALTER TABLE payroll ADD COLUMN pay_period_start DATE;
ALTER TABLE payroll ADD COLUMN pay_period_end DATE;
ALTER TABLE payroll ADD COLUMN gross_pay DECIMAL(12,2);
ALTER TABLE payroll ADD COLUMN deductions DECIMAL(12,2);
ALTER TABLE payroll ADD COLUMN net_pay DECIMAL(12,2);
```

### 🔄 80% Reusable - Employee/HR Module

**Current Capabilities:**
- ✅ Employee CRUD with details
- ✅ Contracts (YEARLY/FULL)
- ✅ Finances (salary, allowances, bonuses)
- ✅ Departments & roles
- ✅ Addresses

**Additions for ERP:**
- Commission tracking for sales reps
- Buyer assignment for purchasing
- Security levels per user
- Driver assignments for routing

**New Fields:**
```sql
ALTER TABLE employees ADD COLUMN is_sales_rep BOOLEAN DEFAULT false;
ALTER TABLE employees ADD COLUMN is_buyer BOOLEAN DEFAULT false;
ALTER TABLE employees ADD COLUMN is_driver BOOLEAN DEFAULT false;
ALTER TABLE employees ADD COLUMN commission_rate DECIMAL(5,2);
ALTER TABLE employees ADD COLUMN security_level INT DEFAULT 1;
```

### 🔄 60% Reusable - Financial Tracking

**Current Capabilities:**
- Income tracking (INTERNAL/EXTERNAL)
- Expense management with types
- Cash boxes with detailed transactions

**Can Adapt For:**
- Accounts Payable base
- Accounts Receivable base
- Bank reconciliation foundation

### 🔄 50% Reusable - Document Management

**Current Capabilities:**
- Document upload/download (MinIO)
- Recipient tracking with status
- Read tracking
- Status workflow (SENT → READ → ACCEPTED)

**Can Adapt For:**
- Purchase Order documents
- Invoice attachments
- Shipping labels storage
- Vendor contracts
- Customer agreements

---

## Modules to Build from Scratch

### 🔴 1. Sales/Order Entry Module

**Effort Level:** HIGH

**New Tables Required:**
- `customers`
- `customer_ship_to`
- `customer_order_guides`
- `sales_orders`
- `sales_order_lines`
- `order_types`
- `lost_sales`
- `special_order_requests`

**New Services Required:**
- `CustomerService`
- `SalesOrderService`
- `OrderGuideService`
- `PricingService`

**New Routes:**
- `POST /v1/customers`
- `GET /v1/customers/{id}`
- `PUT /v1/customers/{id}`
- `DELETE /v1/customers/{id}`
- `GET /v1/customers/list`
- `POST /v1/orders`
- `GET /v1/orders/{id}`
- `PUT /v1/orders/{id}`
- `POST /v1/orders/{id}/lines`
- `GET /v1/orders/by-route/{routeId}`
- `GET /v1/customers/{id}/order-guide`

### 🔴 2. Inventory Control Module

**Effort Level:** HIGH

**New Tables Required:**
- `warehouses`
- `products`
- `product_units`
- `product_categories`
- `inventory`
- `inventory_transactions`
- `lot_tracking`
- `product_reservations`

**New Services Required:**
- `WarehouseService`
- `ProductService`
- `InventoryService`
- `LotTrackingService`

**New Routes:**
- `POST /v1/warehouses`
- `GET /v1/warehouses/{id}`
- `GET /v1/warehouses/{id}/inventory`
- `POST /v1/products`
- `GET /v1/products/{id}`
- `GET /v1/products/{id}/inventory`
- `POST /v1/inventory/adjust`
- `POST /v1/inventory/transfer`
- `GET /v1/inventory/by-lot/{lotNumber}`

### 🔴 3. Purchasing & Receiving Module

**Effort Level:** HIGH

**New Tables Required:**
- `vendors`
- `vendor_products`
- `purchase_orders`
- `purchase_order_lines`
- `receiving`
- `receiving_lines`
- `short_shipments`

**New Services Required:**
- `VendorService`
- `PurchaseOrderService`
- `ReceivingService`
- `BuyerToolsService`

**New Routes:**
- `POST /v1/vendors`
- `GET /v1/vendors/{id}`
- `GET /v1/vendors/{id}/products`
- `POST /v1/purchase-orders`
- `GET /v1/purchase-orders/{id}`
- `POST /v1/purchase-orders/{id}/receive`
- `GET /v1/buyers/{id}/suggested-orders`

### 🔴 4. Accounts Receivable (AR) Module

**Effort Level:** HIGH

**New Tables Required:**
- `ar_invoices`
- `ar_invoice_lines`
- `ar_payments`
- `ar_payment_allocations`
- `ar_credit_memos`
- `customer_statements`
- `credit_limits`

**New Services Required:**
- `ARInvoiceService`
- `ARPaymentService`
- `CreditManagementService`
- `StatementService`

**New Routes:**
- `POST /v1/ar/invoices`
- `GET /v1/ar/invoices/{id}`
- `GET /v1/ar/customers/{id}/invoices`
- `POST /v1/ar/payments`
- `GET /v1/ar/customers/{id}/aging`
- `GET /v1/ar/customers/{id}/statement`

### 🔴 5. Accounts Payable (AP) Module

**Effort Level:** MEDIUM

**New Tables Required:**
- `ap_invoices`
- `ap_invoice_lines`
- `ap_payments`
- `ap_payment_allocations`
- `vendor_discounts`
- `check_runs`

**New Services Required:**
- `APInvoiceService`
- `APPaymentService`
- `CheckRunService`

**New Routes:**
- `POST /v1/ap/invoices`
- `GET /v1/ap/invoices/{id}`
- `GET /v1/ap/vendors/{id}/invoices`
- `POST /v1/ap/payments`
- `POST /v1/ap/check-runs`
- `GET /v1/ap/vendors/{id}/aging`

### 🔴 6. General Ledger Module

**Effort Level:** HIGH

**New Tables Required:**
- `gl_accounts`
- `gl_transactions`
- `gl_journal_entries`
- `gl_recurring_entries`
- `gl_budgets`
- `gl_periods`
- `gl_fiscal_years`

**New Services Required:**
- `GLAccountService`
- `GLTransactionService`
- `GLReportingService`
- `BudgetService`

**New Routes:**
- `POST /v1/gl/accounts`
- `GET /v1/gl/accounts/{id}`
- `GET /v1/gl/accounts/tree`
- `POST /v1/gl/journal-entries`
- `GET /v1/gl/trial-balance`
- `GET /v1/gl/income-statement`
- `GET /v1/gl/balance-sheet`
- `GET /v1/gl/comparison-report`

### 🔴 7. Pricing & Cost Management Module

**Effort Level:** MEDIUM

**New Tables Required:**
- `price_levels`
- `product_prices`
- `customer_prices`
- `contract_prices`
- `promotional_prices`
- `cost_history`
- `rebates`

**New Services Required:**
- `PricingService`
- `CostingService`
- `RebateService`

**New Routes:**
- `POST /v1/pricing/levels`
- `GET /v1/pricing/products/{id}/prices`
- `PUT /v1/pricing/products/{id}/prices`
- `POST /v1/pricing/promotions`
- `GET /v1/pricing/customers/{id}/prices`
- `POST /v1/pricing/mass-update`

### 🔴 8. Picking & Routing Module

**Effort Level:** MEDIUM

**New Tables Required:**
- `routes`
- `route_stops`
- `pick_lists`
- `pick_list_lines`
- `pick_zones`
- `shipping_labels`

**New Services Required:**
- `RouteService`
- `PickListService`
- `ShippingService`

**New Routes:**
- `POST /v1/routes`
- `GET /v1/routes/{id}`
- `GET /v1/routes/{id}/stops`
- `POST /v1/pick-lists/generate`
- `GET /v1/pick-lists/{id}`
- `PUT /v1/pick-lists/{id}/lines/{lineId}`
- `POST /v1/shipping/labels`

### 🔴 9. Bank & Reconciliation Module

**Effort Level:** MEDIUM

**New Tables Required:**
- `bank_accounts`
- `bank_transactions`
- `bank_reconciliations`
- `checks`
- `receipts`

**New Services Required:**
- `BankAccountService`
- `BankReconciliationService`
- `CheckService`

**New Routes:**
- `POST /v1/bank/accounts`
- `GET /v1/bank/accounts/{id}`
- `GET /v1/bank/accounts/{id}/transactions`
- `POST /v1/bank/reconciliation`
- `POST /v1/bank/checks`
- `PUT /v1/bank/checks/{id}/void`

### 🔴 10. Warehouse Management System (WMS)

**Effort Level:** HIGH

**New Tables Required:**
- `warehouse_zones`
- `warehouse_locations`
- `put_away_tasks`
- `replenishment_tasks`
- `pick_tasks`
- `cycle_counts`
- `skid_tracking`
- `warehouse_transfers`
- `disposals`

**New Services Required:**
- `WMSLocationService`
- `PutAwayService`
- `ReplenishmentService`
- `CycleCountService`
- `SkidService`
- `TransferService`

**New Routes:**
- `POST /v1/wms/locations`
- `GET /v1/wms/locations/{id}`
- `GET /v1/wms/locations/{id}/inventory`
- `POST /v1/wms/put-away`
- `POST /v1/wms/replenish`
- `POST /v1/wms/cycle-count`
- `POST /v1/wms/transfer`
- `POST /v1/wms/dispose`

---

## Database Schema Design

### Core ERP Tables

```sql
-- =====================================================
-- CUSTOMERS & SALES
-- =====================================================

CREATE TYPE order_type AS ENUM (
    'STANDARD',
    'ADVANCE',
    'PRE_PAID',
    'ON_HOLD',
    'QUOTE',
    'CREDIT_MEMO',
    'PICK_UP'
);

CREATE TYPE order_status AS ENUM (
    'DRAFT',
    'CONFIRMED',
    'PICKING',
    'SHIPPED',
    'DELIVERED',
    'INVOICED',
    'CANCELLED'
);

CREATE TABLE IF NOT EXISTS customers (
    id SERIAL PRIMARY KEY,
    customer_code VARCHAR(20) UNIQUE NOT NULL,
    name TEXT NOT NULL,
    billing_address_id INT,
    credit_limit DECIMAL(12,2) DEFAULT 0,
    current_balance DECIMAL(12,2) DEFAULT 0,
    payment_terms_days INT DEFAULT 30,
    currency VARCHAR(3) DEFAULT 'USD',
    sales_rep_id INT REFERENCES employees(id),
    default_route_id INT,
    default_warehouse_id INT,
    tax_exempt BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    created_by INT REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS customer_ship_to (
    id SERIAL PRIMARY KEY,
    customer_id INT REFERENCES customers(id) ON DELETE CASCADE,
    ship_to_code VARCHAR(20),
    name TEXT,
    address_line1 TEXT,
    address_line2 TEXT,
    city VARCHAR(100),
    state VARCHAR(100),
    postal_code VARCHAR(20),
    country VARCHAR(100),
    phone VARCHAR(50),
    is_default BOOLEAN DEFAULT false,
    warehouse_id INT,
    route_id INT,
    UNIQUE(customer_id, ship_to_code)
);

CREATE TABLE IF NOT EXISTS customer_order_guides (
    id SERIAL PRIMARY KEY,
    customer_id INT REFERENCES customers(id) ON DELETE CASCADE,
    product_id INT REFERENCES products(id) ON DELETE CASCADE,
    default_quantity DECIMAL(10,3),
    last_ordered_date DATE,
    last_ordered_quantity DECIMAL(10,3),
    avg_weekly_quantity DECIMAL(10,3),
    times_ordered INT DEFAULT 0,
    is_push_item BOOLEAN DEFAULT false,
    custom_price DECIMAL(10,4),
    UNIQUE(customer_id, product_id)
);

CREATE TABLE IF NOT EXISTS sales_orders (
    id SERIAL PRIMARY KEY,
    order_number VARCHAR(20) UNIQUE NOT NULL,
    customer_id INT REFERENCES customers(id),
    ship_to_id INT REFERENCES customer_ship_to(id),
    order_type order_type DEFAULT 'STANDARD',
    order_date DATE DEFAULT CURRENT_DATE,
    requested_ship_date DATE,
    actual_ship_date DATE,
    warehouse_id INT REFERENCES warehouses(id),
    route_id INT REFERENCES routes(id),
    status order_status DEFAULT 'DRAFT',
    subtotal DECIMAL(12,2) DEFAULT 0,
    tax_amount DECIMAL(12,2) DEFAULT 0,
    freight_amount DECIMAL(12,2) DEFAULT 0,
    discount_amount DECIMAL(12,2) DEFAULT 0,
    total_amount DECIMAL(12,2) DEFAULT 0,
    notes TEXT,
    po_number VARCHAR(50),
    sales_rep_id INT REFERENCES employees(id),
    created_by INT REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS sales_order_lines (
    id SERIAL PRIMARY KEY,
    order_id INT REFERENCES sales_orders(id) ON DELETE CASCADE,
    line_number INT,
    product_id INT REFERENCES products(id),
    description TEXT,
    quantity_ordered DECIMAL(10,3),
    quantity_shipped DECIMAL(10,3) DEFAULT 0,
    unit_of_measure VARCHAR(10),
    unit_price DECIMAL(10,4),
    discount_percent DECIMAL(5,2) DEFAULT 0,
    line_total DECIMAL(12,2),
    lot_number VARCHAR(50),
    expiry_date DATE,
    catch_weight DECIMAL(10,3),
    cost DECIMAL(10,4),
    UNIQUE(order_id, line_number)
);

CREATE TABLE IF NOT EXISTS lost_sales (
    id SERIAL PRIMARY KEY,
    order_id INT REFERENCES sales_orders(id),
    product_id INT REFERENCES products(id),
    quantity_requested DECIMAL(10,3),
    quantity_available DECIMAL(10,3),
    reason TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- PRODUCTS & INVENTORY
-- =====================================================

CREATE TABLE IF NOT EXISTS warehouses (
    id SERIAL PRIMARY KEY,
    warehouse_code VARCHAR(10) UNIQUE NOT NULL,
    name TEXT NOT NULL,
    address_line1 TEXT,
    address_line2 TEXT,
    city VARCHAR(100),
    state VARCHAR(100),
    postal_code VARCHAR(20),
    country VARCHAR(100),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS product_categories (
    id SERIAL PRIMARY KEY,
    code VARCHAR(20) UNIQUE,
    name TEXT NOT NULL,
    parent_id INT REFERENCES product_categories(id),
    gl_sales_account_id INT,
    gl_cogs_account_id INT,
    gl_inventory_account_id INT
);

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    sku VARCHAR(50) UNIQUE NOT NULL,
    barcode VARCHAR(50),
    upc VARCHAR(50),
    name TEXT NOT NULL,
    description TEXT,
    category_id INT REFERENCES product_categories(id),
    base_unit VARCHAR(10) DEFAULT 'EA',
    is_catch_weight BOOLEAN DEFAULT false,
    catch_weight_unit VARCHAR(10),
    country_of_origin VARCHAR(3),
    shelf_life_days INT,
    min_shelf_life_days INT,
    is_lot_tracked BOOLEAN DEFAULT false,
    is_serialized BOOLEAN DEFAULT false,
    haccp_category VARCHAR(50),
    qc_required BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS product_units (
    id SERIAL PRIMARY KEY,
    product_id INT REFERENCES products(id) ON DELETE CASCADE,
    unit_name VARCHAR(10) NOT NULL,
    description TEXT,
    conversion_factor DECIMAL(10,6) NOT NULL,
    barcode VARCHAR(50),
    weight DECIMAL(10,4),
    is_purchase_unit BOOLEAN DEFAULT false,
    is_sales_unit BOOLEAN DEFAULT true,
    UNIQUE(product_id, unit_name)
);

CREATE TABLE IF NOT EXISTS inventory (
    id SERIAL PRIMARY KEY,
    product_id INT REFERENCES products(id),
    warehouse_id INT REFERENCES warehouses(id),
    location_code VARCHAR(50),
    lot_number VARCHAR(50),
    production_date DATE,
    expiry_date DATE,
    quantity_on_hand DECIMAL(12,3) DEFAULT 0,
    quantity_allocated DECIMAL(12,3) DEFAULT 0,
    quantity_on_order DECIMAL(12,3) DEFAULT 0,
    quantity_available DECIMAL(12,3) GENERATED ALWAYS AS (quantity_on_hand - quantity_allocated) STORED,
    last_cost DECIMAL(10,4),
    average_cost DECIMAL(10,4),
    last_counted_date DATE,
    last_movement_date DATE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(product_id, warehouse_id, location_code, lot_number)
);

CREATE TYPE inventory_transaction_type AS ENUM (
    'RECEIVE',
    'SHIP',
    'ADJUST_IN',
    'ADJUST_OUT',
    'TRANSFER_IN',
    'TRANSFER_OUT',
    'RETURN',
    'DISPOSE',
    'CYCLE_COUNT'
);

CREATE TABLE IF NOT EXISTS inventory_transactions (
    id SERIAL PRIMARY KEY,
    product_id INT REFERENCES products(id),
    warehouse_id INT REFERENCES warehouses(id),
    location_code VARCHAR(50),
    transaction_type inventory_transaction_type,
    quantity DECIMAL(12,3),
    lot_number VARCHAR(50),
    unit_cost DECIMAL(10,4),
    reference_type VARCHAR(20),
    reference_id INT,
    reference_number VARCHAR(50),
    notes TEXT,
    created_by INT REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS product_reservations (
    id SERIAL PRIMARY KEY,
    product_id INT REFERENCES products(id),
    warehouse_id INT REFERENCES warehouses(id),
    reserved_for_type VARCHAR(20),
    reserved_for_id INT,
    quantity DECIMAL(10,3),
    expiry_date DATE,
    created_by INT REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- VENDORS & PURCHASING
-- =====================================================

CREATE TABLE IF NOT EXISTS vendors (
    id SERIAL PRIMARY KEY,
    vendor_code VARCHAR(20) UNIQUE NOT NULL,
    name TEXT NOT NULL,
    address_line1 TEXT,
    address_line2 TEXT,
    city VARCHAR(100),
    state VARCHAR(100),
    postal_code VARCHAR(20),
    country VARCHAR(100),
    phone VARCHAR(50),
    email VARCHAR(255),
    payment_terms_days INT DEFAULT 30,
    currency VARCHAR(3) DEFAULT 'USD',
    lead_time_days INT DEFAULT 7,
    minimum_order DECIMAL(12,2),
    buyer_id INT REFERENCES employees(id),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS vendor_products (
    id SERIAL PRIMARY KEY,
    vendor_id INT REFERENCES vendors(id) ON DELETE CASCADE,
    product_id INT REFERENCES products(id) ON DELETE CASCADE,
    vendor_sku VARCHAR(50),
    vendor_description TEXT,
    unit_of_measure VARCHAR(10),
    unit_cost DECIMAL(10,4),
    minimum_order_qty DECIMAL(10,3),
    lead_time_days INT,
    is_preferred BOOLEAN DEFAULT false,
    UNIQUE(vendor_id, product_id)
);

CREATE TYPE po_status AS ENUM (
    'DRAFT',
    'SUBMITTED',
    'CONFIRMED',
    'PARTIAL',
    'RECEIVED',
    'CLOSED',
    'CANCELLED'
);

CREATE TABLE IF NOT EXISTS purchase_orders (
    id SERIAL PRIMARY KEY,
    po_number VARCHAR(20) UNIQUE NOT NULL,
    vendor_id INT REFERENCES vendors(id),
    warehouse_id INT REFERENCES warehouses(id),
    order_date DATE DEFAULT CURRENT_DATE,
    expected_date DATE,
    received_date DATE,
    status po_status DEFAULT 'DRAFT',
    subtotal DECIMAL(12,2) DEFAULT 0,
    tax_amount DECIMAL(12,2) DEFAULT 0,
    freight_amount DECIMAL(12,2) DEFAULT 0,
    total_amount DECIMAL(12,2) DEFAULT 0,
    notes TEXT,
    buyer_id INT REFERENCES employees(id),
    created_by INT REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS purchase_order_lines (
    id SERIAL PRIMARY KEY,
    po_id INT REFERENCES purchase_orders(id) ON DELETE CASCADE,
    line_number INT,
    product_id INT REFERENCES products(id),
    description TEXT,
    quantity_ordered DECIMAL(10,3),
    quantity_received DECIMAL(10,3) DEFAULT 0,
    unit_of_measure VARCHAR(10),
    unit_cost DECIMAL(10,4),
    line_total DECIMAL(12,2),
    expected_date DATE,
    UNIQUE(po_id, line_number)
);

CREATE TABLE IF NOT EXISTS receiving (
    id SERIAL PRIMARY KEY,
    receiving_number VARCHAR(20) UNIQUE NOT NULL,
    po_id INT REFERENCES purchase_orders(id),
    warehouse_id INT REFERENCES warehouses(id),
    vendor_id INT REFERENCES vendors(id),
    received_date DATE DEFAULT CURRENT_DATE,
    notes TEXT,
    received_by INT REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS receiving_lines (
    id SERIAL PRIMARY KEY,
    receiving_id INT REFERENCES receiving(id) ON DELETE CASCADE,
    po_line_id INT REFERENCES purchase_order_lines(id),
    product_id INT REFERENCES products(id),
    quantity_received DECIMAL(10,3),
    unit_of_measure VARCHAR(10),
    lot_number VARCHAR(50),
    production_date DATE,
    expiry_date DATE,
    location_code VARCHAR(50),
    unit_cost DECIMAL(10,4),
    is_short_shipment BOOLEAN DEFAULT false,
    notes TEXT
);

-- =====================================================
-- ACCOUNTS RECEIVABLE
-- =====================================================

CREATE TYPE ar_invoice_status AS ENUM (
    'DRAFT',
    'POSTED',
    'PARTIAL',
    'PAID',
    'OVERDUE',
    'VOID'
);

CREATE TABLE IF NOT EXISTS ar_invoices (
    id SERIAL PRIMARY KEY,
    invoice_number VARCHAR(20) UNIQUE NOT NULL,
    customer_id INT REFERENCES customers(id),
    ship_to_id INT REFERENCES customer_ship_to(id),
    sales_order_id INT REFERENCES sales_orders(id),
    invoice_date DATE DEFAULT CURRENT_DATE,
    due_date DATE,
    subtotal DECIMAL(12,2) DEFAULT 0,
    tax_amount DECIMAL(12,2) DEFAULT 0,
    freight_amount DECIMAL(12,2) DEFAULT 0,
    discount_amount DECIMAL(12,2) DEFAULT 0,
    total_amount DECIMAL(12,2) DEFAULT 0,
    paid_amount DECIMAL(12,2) DEFAULT 0,
    balance DECIMAL(12,2) DEFAULT 0,
    currency VARCHAR(3) DEFAULT 'USD',
    status ar_invoice_status DEFAULT 'DRAFT',
    notes TEXT,
    gl_posted BOOLEAN DEFAULT false,
    created_by INT REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ar_invoice_lines (
    id SERIAL PRIMARY KEY,
    invoice_id INT REFERENCES ar_invoices(id) ON DELETE CASCADE,
    line_number INT,
    product_id INT REFERENCES products(id),
    description TEXT,
    quantity DECIMAL(10,3),
    unit_price DECIMAL(10,4),
    discount_percent DECIMAL(5,2) DEFAULT 0,
    tax_amount DECIMAL(12,2) DEFAULT 0,
    line_total DECIMAL(12,2),
    gl_account_id INT,
    UNIQUE(invoice_id, line_number)
);

CREATE TABLE IF NOT EXISTS ar_payments (
    id SERIAL PRIMARY KEY,
    payment_number VARCHAR(20) UNIQUE NOT NULL,
    customer_id INT REFERENCES customers(id),
    payment_date DATE DEFAULT CURRENT_DATE,
    amount DECIMAL(12,2),
    currency VARCHAR(3) DEFAULT 'USD',
    payment_method VARCHAR(20),
    reference_number VARCHAR(50),
    check_number VARCHAR(50),
    bank_account_id INT,
    notes TEXT,
    gl_posted BOOLEAN DEFAULT false,
    created_by INT REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ar_payment_allocations (
    id SERIAL PRIMARY KEY,
    payment_id INT REFERENCES ar_payments(id) ON DELETE CASCADE,
    invoice_id INT REFERENCES ar_invoices(id),
    amount_applied DECIMAL(12,2),
    discount_taken DECIMAL(12,2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- ACCOUNTS PAYABLE
-- =====================================================

CREATE TYPE ap_invoice_status AS ENUM (
    'PENDING',
    'APPROVED',
    'PARTIAL',
    'PAID',
    'VOID'
);

CREATE TABLE IF NOT EXISTS ap_invoices (
    id SERIAL PRIMARY KEY,
    invoice_number VARCHAR(50) NOT NULL,
    vendor_id INT REFERENCES vendors(id),
    po_id INT REFERENCES purchase_orders(id),
    invoice_date DATE,
    received_date DATE DEFAULT CURRENT_DATE,
    due_date DATE,
    subtotal DECIMAL(12,2) DEFAULT 0,
    tax_amount DECIMAL(12,2) DEFAULT 0,
    freight_amount DECIMAL(12,2) DEFAULT 0,
    total_amount DECIMAL(12,2) DEFAULT 0,
    paid_amount DECIMAL(12,2) DEFAULT 0,
    balance DECIMAL(12,2) DEFAULT 0,
    currency VARCHAR(3) DEFAULT 'USD',
    status ap_invoice_status DEFAULT 'PENDING',
    notes TEXT,
    gl_posted BOOLEAN DEFAULT false,
    approved_by INT REFERENCES employees(id),
    approved_at TIMESTAMP,
    created_by INT REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(vendor_id, invoice_number)
);

CREATE TABLE IF NOT EXISTS ap_invoice_lines (
    id SERIAL PRIMARY KEY,
    invoice_id INT REFERENCES ap_invoices(id) ON DELETE CASCADE,
    line_number INT,
    description TEXT,
    quantity DECIMAL(10,3),
    unit_cost DECIMAL(10,4),
    line_total DECIMAL(12,2),
    gl_account_id INT,
    po_line_id INT REFERENCES purchase_order_lines(id),
    UNIQUE(invoice_id, line_number)
);

CREATE TABLE IF NOT EXISTS vendor_discounts (
    id SERIAL PRIMARY KEY,
    vendor_id INT REFERENCES vendors(id) ON DELETE CASCADE,
    discount_days INT,
    discount_percent DECIMAL(5,2),
    UNIQUE(vendor_id, discount_days)
);

CREATE TABLE IF NOT EXISTS ap_payments (
    id SERIAL PRIMARY KEY,
    payment_number VARCHAR(20) UNIQUE NOT NULL,
    vendor_id INT REFERENCES vendors(id),
    payment_date DATE DEFAULT CURRENT_DATE,
    amount DECIMAL(12,2),
    currency VARCHAR(3) DEFAULT 'USD',
    payment_method VARCHAR(20),
    check_number VARCHAR(50),
    bank_account_id INT,
    notes TEXT,
    gl_posted BOOLEAN DEFAULT false,
    created_by INT REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ap_payment_allocations (
    id SERIAL PRIMARY KEY,
    payment_id INT REFERENCES ap_payments(id) ON DELETE CASCADE,
    invoice_id INT REFERENCES ap_invoices(id),
    amount_applied DECIMAL(12,2),
    discount_taken DECIMAL(12,2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS check_runs (
    id SERIAL PRIMARY KEY,
    run_date DATE DEFAULT CURRENT_DATE,
    bank_account_id INT,
    cut_off_date DATE,
    vendor_type VARCHAR(20),
    total_amount DECIMAL(12,2),
    check_count INT,
    status VARCHAR(20) DEFAULT 'PENDING',
    created_by INT REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- GENERAL LEDGER
-- =====================================================

CREATE TYPE gl_account_type AS ENUM (
    'ASSET',
    'LIABILITY',
    'EQUITY',
    'REVENUE',
    'EXPENSE'
);

CREATE TABLE IF NOT EXISTS gl_fiscal_years (
    id SERIAL PRIMARY KEY,
    year_name VARCHAR(20) UNIQUE NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_closed BOOLEAN DEFAULT false,
    closed_by INT REFERENCES employees(id),
    closed_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS gl_periods (
    id SERIAL PRIMARY KEY,
    fiscal_year_id INT REFERENCES gl_fiscal_years(id),
    period_number INT,
    period_name VARCHAR(20),
    start_date DATE,
    end_date DATE,
    is_closed BOOLEAN DEFAULT false,
    is_adjustment_period BOOLEAN DEFAULT false,
    UNIQUE(fiscal_year_id, period_number)
);

CREATE TABLE IF NOT EXISTS gl_accounts (
    id SERIAL PRIMARY KEY,
    account_number VARCHAR(20) UNIQUE NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    account_type gl_account_type NOT NULL,
    parent_id INT REFERENCES gl_accounts(id),
    is_header BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    normal_balance VARCHAR(10) DEFAULT 'DEBIT',
    currency VARCHAR(3) DEFAULT 'USD',
    department_id INT REFERENCES departments(id),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS gl_transactions (
    id SERIAL PRIMARY KEY,
    transaction_date DATE NOT NULL,
    period_id INT REFERENCES gl_periods(id),
    account_id INT REFERENCES gl_accounts(id),
    debit DECIMAL(14,2) DEFAULT 0,
    credit DECIMAL(14,2) DEFAULT 0,
    description TEXT,
    reference_type VARCHAR(30),
    reference_id INT,
    reference_number VARCHAR(50),
    currency VARCHAR(3) DEFAULT 'USD',
    exchange_rate DECIMAL(10,6) DEFAULT 1,
    is_posted BOOLEAN DEFAULT false,
    posted_by INT REFERENCES employees(id),
    posted_at TIMESTAMP,
    created_by INT REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS gl_journal_entries (
    id SERIAL PRIMARY KEY,
    entry_number VARCHAR(20) UNIQUE NOT NULL,
    entry_date DATE NOT NULL,
    period_id INT REFERENCES gl_periods(id),
    description TEXT,
    is_reversing BOOLEAN DEFAULT false,
    reverse_date DATE,
    source VARCHAR(20) DEFAULT 'MANUAL',
    is_posted BOOLEAN DEFAULT false,
    posted_by INT REFERENCES employees(id),
    posted_at TIMESTAMP,
    created_by INT REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS gl_journal_entry_lines (
    id SERIAL PRIMARY KEY,
    journal_entry_id INT REFERENCES gl_journal_entries(id) ON DELETE CASCADE,
    line_number INT,
    account_id INT REFERENCES gl_accounts(id),
    debit DECIMAL(14,2) DEFAULT 0,
    credit DECIMAL(14,2) DEFAULT 0,
    description TEXT,
    department_id INT REFERENCES departments(id),
    UNIQUE(journal_entry_id, line_number)
);

CREATE TABLE IF NOT EXISTS gl_recurring_entries (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    description TEXT,
    frequency VARCHAR(20),
    next_run_date DATE,
    last_run_date DATE,
    is_active BOOLEAN DEFAULT true,
    created_by INT REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS gl_recurring_entry_lines (
    id SERIAL PRIMARY KEY,
    recurring_entry_id INT REFERENCES gl_recurring_entries(id) ON DELETE CASCADE,
    line_number INT,
    account_id INT REFERENCES gl_accounts(id),
    debit DECIMAL(14,2) DEFAULT 0,
    credit DECIMAL(14,2) DEFAULT 0,
    description TEXT,
    UNIQUE(recurring_entry_id, line_number)
);

CREATE TABLE IF NOT EXISTS gl_budgets (
    id SERIAL PRIMARY KEY,
    account_id INT REFERENCES gl_accounts(id),
    period_id INT REFERENCES gl_periods(id),
    budget_amount DECIMAL(14,2) DEFAULT 0,
    notes TEXT,
    created_by INT REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(account_id, period_id)
);

-- =====================================================
-- PRICING
-- =====================================================

CREATE TABLE IF NOT EXISTS price_levels (
    id SERIAL PRIMARY KEY,
    level_code VARCHAR(10) UNIQUE NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    discount_percent DECIMAL(5,2) DEFAULT 0,
    is_active BOOLEAN DEFAULT true
);

CREATE TABLE IF NOT EXISTS product_prices (
    id SERIAL PRIMARY KEY,
    product_id INT REFERENCES products(id) ON DELETE CASCADE,
    price_level_id INT REFERENCES price_levels(id),
    unit_of_measure VARCHAR(10),
    unit_price DECIMAL(10,4),
    effective_date DATE DEFAULT CURRENT_DATE,
    expiry_date DATE,
    min_quantity DECIMAL(10,3) DEFAULT 0,
    UNIQUE(product_id, price_level_id, unit_of_measure, effective_date)
);

CREATE TABLE IF NOT EXISTS customer_prices (
    id SERIAL PRIMARY KEY,
    customer_id INT REFERENCES customers(id) ON DELETE CASCADE,
    product_id INT REFERENCES products(id) ON DELETE CASCADE,
    unit_of_measure VARCHAR(10),
    unit_price DECIMAL(10,4),
    effective_date DATE DEFAULT CURRENT_DATE,
    expiry_date DATE,
    contract_number VARCHAR(50),
    UNIQUE(customer_id, product_id, unit_of_measure, effective_date)
);

CREATE TABLE IF NOT EXISTS promotional_prices (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    product_id INT REFERENCES products(id),
    category_id INT REFERENCES product_categories(id),
    promo_price DECIMAL(10,4),
    discount_percent DECIMAL(5,2),
    start_date DATE,
    end_date DATE,
    min_quantity DECIMAL(10,3) DEFAULT 0,
    is_active BOOLEAN DEFAULT true
);

CREATE TABLE IF NOT EXISTS product_costs (
    id SERIAL PRIMARY KEY,
    product_id INT REFERENCES products(id) ON DELETE CASCADE,
    warehouse_id INT REFERENCES warehouses(id),
    cost_type VARCHAR(20),
    unit_cost DECIMAL(10,4),
    effective_date DATE DEFAULT CURRENT_DATE,
    UNIQUE(product_id, warehouse_id, cost_type, effective_date)
);

-- =====================================================
-- ROUTING & PICKING
-- =====================================================

CREATE TABLE IF NOT EXISTS routes (
    id SERIAL PRIMARY KEY,
    route_code VARCHAR(10) UNIQUE NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    warehouse_id INT REFERENCES warehouses(id),
    driver_id INT REFERENCES employees(id),
    vehicle_id INT,
    departure_time TIME,
    is_active BOOLEAN DEFAULT true
);

CREATE TABLE IF NOT EXISTS route_stops (
    id SERIAL PRIMARY KEY,
    route_id INT REFERENCES routes(id) ON DELETE CASCADE,
    customer_id INT REFERENCES customers(id),
    ship_to_id INT REFERENCES customer_ship_to(id),
    stop_sequence INT,
    estimated_arrival TIME,
    notes TEXT,
    UNIQUE(route_id, stop_sequence)
);

CREATE TYPE pick_list_status AS ENUM (
    'PENDING',
    'IN_PROGRESS',
    'COMPLETE',
    'CANCELLED'
);

CREATE TABLE IF NOT EXISTS pick_lists (
    id SERIAL PRIMARY KEY,
    pick_number VARCHAR(20) UNIQUE NOT NULL,
    warehouse_id INT REFERENCES warehouses(id),
    route_id INT REFERENCES routes(id),
    pick_date DATE DEFAULT CURRENT_DATE,
    status pick_list_status DEFAULT 'PENDING',
    picker_id INT REFERENCES employees(id),
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_by INT REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS pick_list_lines (
    id SERIAL PRIMARY KEY,
    pick_list_id INT REFERENCES pick_lists(id) ON DELETE CASCADE,
    order_id INT REFERENCES sales_orders(id),
    order_line_id INT REFERENCES sales_order_lines(id),
    product_id INT REFERENCES products(id),
    location_code VARCHAR(50),
    lot_number VARCHAR(50),
    quantity_to_pick DECIMAL(10,3),
    quantity_picked DECIMAL(10,3) DEFAULT 0,
    picked_by INT REFERENCES employees(id),
    picked_at TIMESTAMP,
    notes TEXT
);

CREATE TABLE IF NOT EXISTS pick_zones (
    id SERIAL PRIMARY KEY,
    warehouse_id INT REFERENCES warehouses(id),
    zone_code VARCHAR(10) NOT NULL,
    name TEXT,
    pick_sequence INT,
    UNIQUE(warehouse_id, zone_code)
);

-- =====================================================
-- BANK & RECONCILIATION
-- =====================================================

CREATE TABLE IF NOT EXISTS bank_accounts (
    id SERIAL PRIMARY KEY,
    account_code VARCHAR(20) UNIQUE NOT NULL,
    name TEXT NOT NULL,
    bank_name TEXT,
    account_number VARCHAR(50),
    routing_number VARCHAR(50),
    currency VARCHAR(3) DEFAULT 'USD',
    current_balance DECIMAL(14,2) DEFAULT 0,
    gl_account_id INT REFERENCES gl_accounts(id),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TYPE bank_transaction_type AS ENUM (
    'DEPOSIT',
    'WITHDRAWAL',
    'CHECK',
    'TRANSFER',
    'FEE',
    'INTEREST',
    'ADJUSTMENT'
);

CREATE TABLE IF NOT EXISTS bank_transactions (
    id SERIAL PRIMARY KEY,
    bank_account_id INT REFERENCES bank_accounts(id),
    transaction_date DATE,
    transaction_type bank_transaction_type,
    amount DECIMAL(12,2),
    reference_number VARCHAR(50),
    description TEXT,
    is_reconciled BOOLEAN DEFAULT false,
    reconciled_date DATE,
    source_type VARCHAR(20),
    source_id INT,
    created_by INT REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS bank_reconciliations (
    id SERIAL PRIMARY KEY,
    bank_account_id INT REFERENCES bank_accounts(id),
    statement_date DATE,
    statement_balance DECIMAL(14,2),
    reconciled_balance DECIMAL(14,2),
    difference DECIMAL(14,2),
    is_complete BOOLEAN DEFAULT false,
    completed_by INT REFERENCES employees(id),
    completed_at TIMESTAMP,
    created_by INT REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS checks (
    id SERIAL PRIMARY KEY,
    bank_account_id INT REFERENCES bank_accounts(id),
    check_number VARCHAR(20) NOT NULL,
    check_date DATE,
    payee_type VARCHAR(20),
    payee_id INT,
    payee_name TEXT,
    amount DECIMAL(12,2),
    memo TEXT,
    is_printed BOOLEAN DEFAULT false,
    is_voided BOOLEAN DEFAULT false,
    voided_by INT REFERENCES employees(id),
    voided_at TIMESTAMP,
    void_reason TEXT,
    cleared_date DATE,
    created_by INT REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(bank_account_id, check_number)
);

-- =====================================================
-- WAREHOUSE MANAGEMENT SYSTEM (WMS)
-- =====================================================

CREATE TABLE IF NOT EXISTS warehouse_zones (
    id SERIAL PRIMARY KEY,
    warehouse_id INT REFERENCES warehouses(id),
    zone_code VARCHAR(10) NOT NULL,
    name TEXT,
    zone_type VARCHAR(20),
    temperature_controlled BOOLEAN DEFAULT false,
    min_temperature DECIMAL(5,2),
    max_temperature DECIMAL(5,2),
    UNIQUE(warehouse_id, zone_code)
);

CREATE TABLE IF NOT EXISTS warehouse_locations (
    id SERIAL PRIMARY KEY,
    warehouse_id INT REFERENCES warehouses(id),
    zone_id INT REFERENCES warehouse_zones(id),
    location_code VARCHAR(50) UNIQUE NOT NULL,
    aisle VARCHAR(10),
    rack VARCHAR(10),
    shelf VARCHAR(10),
    bin VARCHAR(10),
    location_type VARCHAR(20),
    max_weight DECIMAL(10,2),
    max_volume DECIMAL(10,2),
    is_active BOOLEAN DEFAULT true,
    pick_sequence INT
);

CREATE TYPE wms_task_status AS ENUM (
    'PENDING',
    'ASSIGNED',
    'IN_PROGRESS',
    'COMPLETE',
    'CANCELLED'
);

CREATE TABLE IF NOT EXISTS put_away_tasks (
    id SERIAL PRIMARY KEY,
    receiving_id INT REFERENCES receiving(id),
    receiving_line_id INT REFERENCES receiving_lines(id),
    product_id INT REFERENCES products(id),
    from_location VARCHAR(50),
    to_location VARCHAR(50),
    quantity DECIMAL(10,3),
    lot_number VARCHAR(50),
    status wms_task_status DEFAULT 'PENDING',
    assigned_to INT REFERENCES employees(id),
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS replenishment_tasks (
    id SERIAL PRIMARY KEY,
    warehouse_id INT REFERENCES warehouses(id),
    product_id INT REFERENCES products(id),
    from_location VARCHAR(50),
    to_location VARCHAR(50),
    quantity DECIMAL(10,3),
    lot_number VARCHAR(50),
    priority INT DEFAULT 5,
    status wms_task_status DEFAULT 'PENDING',
    assigned_to INT REFERENCES employees(id),
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS cycle_counts (
    id SERIAL PRIMARY KEY,
    warehouse_id INT REFERENCES warehouses(id),
    count_date DATE DEFAULT CURRENT_DATE,
    location_code VARCHAR(50),
    product_id INT REFERENCES products(id),
    lot_number VARCHAR(50),
    system_quantity DECIMAL(10,3),
    counted_quantity DECIMAL(10,3),
    variance DECIMAL(10,3),
    variance_value DECIMAL(12,2),
    is_approved BOOLEAN DEFAULT false,
    approved_by INT REFERENCES employees(id),
    approved_at TIMESTAMP,
    counted_by INT REFERENCES employees(id),
    counted_at TIMESTAMP,
    notes TEXT
);

CREATE TABLE IF NOT EXISTS skid_tracking (
    id SERIAL PRIMARY KEY,
    skid_number VARCHAR(50) UNIQUE NOT NULL,
    warehouse_id INT REFERENCES warehouses(id),
    location_code VARCHAR(50),
    product_id INT REFERENCES products(id),
    lot_number VARCHAR(50),
    quantity DECIMAL(10,3),
    status VARCHAR(20) DEFAULT 'ACTIVE',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS warehouse_transfers (
    id SERIAL PRIMARY KEY,
    transfer_number VARCHAR(20) UNIQUE NOT NULL,
    from_warehouse_id INT REFERENCES warehouses(id),
    to_warehouse_id INT REFERENCES warehouses(id),
    transfer_date DATE DEFAULT CURRENT_DATE,
    status VARCHAR(20) DEFAULT 'PENDING',
    notes TEXT,
    created_by INT REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS warehouse_transfer_lines (
    id SERIAL PRIMARY KEY,
    transfer_id INT REFERENCES warehouse_transfers(id) ON DELETE CASCADE,
    product_id INT REFERENCES products(id),
    lot_number VARCHAR(50),
    from_location VARCHAR(50),
    to_location VARCHAR(50),
    quantity DECIMAL(10,3),
    received_quantity DECIMAL(10,3) DEFAULT 0
);

CREATE TABLE IF NOT EXISTS disposals (
    id SERIAL PRIMARY KEY,
    warehouse_id INT REFERENCES warehouses(id),
    disposal_date DATE DEFAULT CURRENT_DATE,
    product_id INT REFERENCES products(id),
    location_code VARCHAR(50),
    lot_number VARCHAR(50),
    quantity DECIMAL(10,3),
    reason VARCHAR(100),
    cost DECIMAL(12,2),
    approved_by INT REFERENCES employees(id),
    approved_at TIMESTAMP,
    created_by INT REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- INDEXES FOR PERFORMANCE
-- =====================================================

CREATE INDEX idx_customers_code ON customers(customer_code);
CREATE INDEX idx_customers_sales_rep ON customers(sales_rep_id);
CREATE INDEX idx_sales_orders_customer ON sales_orders(customer_id);
CREATE INDEX idx_sales_orders_date ON sales_orders(order_date);
CREATE INDEX idx_sales_orders_status ON sales_orders(status);
CREATE INDEX idx_sales_order_lines_product ON sales_order_lines(product_id);
CREATE INDEX idx_products_sku ON products(sku);
CREATE INDEX idx_products_barcode ON products(barcode);
CREATE INDEX idx_inventory_product ON inventory(product_id);
CREATE INDEX idx_inventory_warehouse ON inventory(warehouse_id);
CREATE INDEX idx_inventory_lot ON inventory(lot_number);
CREATE INDEX idx_inventory_expiry ON inventory(expiry_date);
CREATE INDEX idx_inventory_transactions_product ON inventory_transactions(product_id);
CREATE INDEX idx_inventory_transactions_date ON inventory_transactions(created_at);
CREATE INDEX idx_vendors_code ON vendors(vendor_code);
CREATE INDEX idx_purchase_orders_vendor ON purchase_orders(vendor_id);
CREATE INDEX idx_purchase_orders_status ON purchase_orders(status);
CREATE INDEX idx_ar_invoices_customer ON ar_invoices(customer_id);
CREATE INDEX idx_ar_invoices_status ON ar_invoices(status);
CREATE INDEX idx_ar_invoices_due_date ON ar_invoices(due_date);
CREATE INDEX idx_ap_invoices_vendor ON ap_invoices(vendor_id);
CREATE INDEX idx_ap_invoices_status ON ap_invoices(status);
CREATE INDEX idx_ap_invoices_due_date ON ap_invoices(due_date);
CREATE INDEX idx_gl_transactions_date ON gl_transactions(transaction_date);
CREATE INDEX idx_gl_transactions_account ON gl_transactions(account_id);
CREATE INDEX idx_gl_transactions_period ON gl_transactions(period_id);
CREATE INDEX idx_pick_lists_date ON pick_lists(pick_date);
CREATE INDEX idx_pick_lists_status ON pick_lists(status);
CREATE INDEX idx_warehouse_locations_code ON warehouse_locations(location_code);
```

---

## Implementation Roadmap

### Phase 1: Foundation (Weeks 1-4)

**Week 1-2: Core Setup**
- [ ] Fork UMS codebase for ERP project
- [ ] Set up new database schema
- [ ] Extend employee module (sales rep, buyer, driver flags)
- [ ] Create base middleware structure

**Week 3-4: Master Data**
- [ ] Implement Warehouses module
- [ ] Implement Products module with units of measure
- [ ] Implement Product Categories
- [ ] Implement Customers module
- [ ] Implement Vendors module

**Deliverables:**
- Working API for master data CRUD
- Database migrations
- API documentation

### Phase 2: Core ERP (Weeks 5-12)

**Week 5-6: Inventory Foundation**
- [ ] Implement Inventory module
- [ ] Implement Inventory Transactions
- [ ] Implement Lot Tracking
- [ ] Real-time quantity calculations

**Week 7-8: Purchasing**
- [ ] Implement Purchase Orders
- [ ] Implement Receiving
- [ ] Implement Short Shipment handling
- [ ] PO to Inventory integration

**Week 9-10: Sales Orders**
- [ ] Implement Sales Orders
- [ ] Implement Order Lines
- [ ] Implement Customer Order Guides
- [ ] Lost Sales tracking
- [ ] Inventory allocation

**Week 11-12: AR/AP Foundation**
- [ ] Implement AR Invoices
- [ ] Implement AR Payments
- [ ] Implement AP Invoices
- [ ] Implement AP Payments

**Deliverables:**
- Complete purchase-to-pay cycle
- Complete order-to-cash cycle
- Inventory management
- API documentation

### Phase 3: Advanced Features (Weeks 13-20)

**Week 13-14: General Ledger**
- [ ] Implement GL Accounts (Chart of Accounts)
- [ ] Implement GL Transactions
- [ ] Implement Journal Entries
- [ ] Implement GL Periods & Fiscal Years

**Week 15-16: Financial Integration**
- [ ] AR to GL posting
- [ ] AP to GL posting
- [ ] Inventory to GL posting
- [ ] Payroll to GL posting

**Week 17-18: Pricing & Costing**
- [ ] Implement Price Levels
- [ ] Implement Customer Pricing
- [ ] Implement Promotional Pricing
- [ ] Implement Cost Tracking

**Week 19-20: Bank & Reconciliation**
- [ ] Implement Bank Accounts
- [ ] Implement Bank Transactions
- [ ] Implement Check Management
- [ ] Implement Bank Reconciliation

**Deliverables:**
- Full financial system
- GL integration
- Financial reporting
- Bank management

### Phase 4: WMS & Advanced (Weeks 21-28)

**Week 21-22: Routing & Picking**
- [ ] Implement Routes
- [ ] Implement Route Stops
- [ ] Implement Pick Lists
- [ ] Implement Pick Zones

**Week 23-24: WMS Foundation**
- [ ] Implement Warehouse Locations
- [ ] Implement Zone Management
- [ ] Implement Put-Away Tasks
- [ ] Implement Replenishment

**Week 25-26: WMS Advanced**
- [ ] Implement Cycle Counts
- [ ] Implement Skid Tracking
- [ ] Implement Warehouse Transfers
- [ ] Implement Disposals

**Week 27-28: Integration & Polish**
- [ ] Mobile API endpoints
- [ ] Barcode scanning integration
- [ ] Label printing
- [ ] Performance optimization

**Deliverables:**
- Complete WMS
- Mobile-ready APIs
- Barcode/label support

### Phase 5: Reporting & Polish (Weeks 29-32)

**Week 29-30: Reporting**
- [ ] Financial statements (Income Statement, Balance Sheet)
- [ ] AR Aging reports
- [ ] AP Aging reports
- [ ] Inventory reports
- [ ] Sales reports

**Week 31-32: Final Integration**
- [ ] End-to-end testing
- [ ] Performance testing
- [ ] Security audit
- [ ] Documentation completion

---

## Technical Recommendations

### 1. Security Improvements

```go
// Use bcrypt for password hashing (NOT plain text)
import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

### 2. Database Connection Pooling

```go
// Use pgxpool instead of single connection
import "github.com/jackc/pgx/v5/pgxpool"

func NewPool(connString string) (*pgxpool.Pool, error) {
    config, err := pgxpool.ParseConfig(connString)
    if err != nil {
        return nil, err
    }
    
    config.MaxConns = 25
    config.MinConns = 5
    config.MaxConnLifetime = time.Hour
    config.MaxConnIdleTime = 30 * time.Minute
    
    return pgxpool.NewWithConfig(context.Background(), config)
}
```

### 3. Input Validation

```go
// Add validation to all services
type Validator struct {
    Errors map[string]string
}

func (v *Validator) Check(ok bool, key, message string) {
    if !ok {
        v.AddError(key, message)
    }
}

func (v *Validator) Valid() bool {
    return len(v.Errors) == 0
}
```

### 4. Transaction Management

```go
// Use context-based transactions
func (s *Service) CreateWithTransaction(ctx context.Context, data Data) error {
    tx, err := s.db.BeginTx(ctx)
    if err != nil {
        return err
    }
    defer tx.Rollback(ctx)
    
    // Multiple operations...
    
    return tx.Commit(ctx)
}
```

### 5. API Versioning

```go
// Already implemented - extend for v2
app.Mount("/v1", v1.Router(...))
app.Mount("/v2", v2.Router(...)) // Future ERP version
```

### 6. Audit Logging

```sql
-- Add audit table
CREATE TABLE audit_log (
    id SERIAL PRIMARY KEY,
    table_name VARCHAR(100),
    record_id INT,
    action VARCHAR(20),
    old_values JSONB,
    new_values JSONB,
    user_id INT,
    ip_address VARCHAR(50),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create audit trigger function
CREATE OR REPLACE FUNCTION audit_trigger_func()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'DELETE' THEN
        INSERT INTO audit_log (table_name, record_id, action, old_values, user_id)
        VALUES (TG_TABLE_NAME, OLD.id, 'DELETE', row_to_json(OLD), current_setting('app.user_id')::int);
        RETURN OLD;
    ELSIF TG_OP = 'UPDATE' THEN
        INSERT INTO audit_log (table_name, record_id, action, old_values, new_values, user_id)
        VALUES (TG_TABLE_NAME, NEW.id, 'UPDATE', row_to_json(OLD), row_to_json(NEW), current_setting('app.user_id')::int);
        RETURN NEW;
    ELSIF TG_OP = 'INSERT' THEN
        INSERT INTO audit_log (table_name, record_id, action, new_values, user_id)
        VALUES (TG_TABLE_NAME, NEW.id, 'INSERT', row_to_json(NEW), current_setting('app.user_id')::int);
        RETURN NEW;
    END IF;
END;
$$ LANGUAGE plpgsql;
```

---

## Summary Matrix

| Module | Reusability | Effort | Priority |
|--------|-------------|--------|----------|
| Core Infrastructure | 100% ✅ | Low | P0 |
| Employee/HR | 80% 🔄 | Low | P0 |
| Payroll | 70% 🔄 | Medium | P2 |
| Financial Base | 60% 🔄 | Medium | P1 |
| Documents | 50% 🔄 | Medium | P2 |
| Products & Inventory | 0% 🔴 | **High** | P0 |
| Customers | 0% 🔴 | **High** | P0 |
| Vendors | 0% 🔴 | **High** | P0 |
| Sales Orders | 0% 🔴 | **High** | P1 |
| Purchasing | 0% 🔴 | **High** | P1 |
| AR Management | 0% 🔴 | **High** | P1 |
| AP Management | 0% 🔴 | Medium | P1 |
| General Ledger | 0% 🔴 | **High** | P2 |
| Pricing | 0% 🔴 | Medium | P2 |
| Picking & Routing | 0% 🔴 | Medium | P2 |
| Bank Reconciliation | 0% 🔴 | Medium | P2 |
| WMS | 0% 🔴 | **High** | P3 |

---

## Appendix: File Structure for New Modules

```
registration/
├── src/
│   └── v1/
│       ├── middlewares/
│       │   ├── customer/
│       │   │   └── customer.go
│       │   ├── vendor/
│       │   │   └── vendor.go
│       │   ├── product/
│       │   │   └── product.go
│       │   ├── inventory/
│       │   │   └── inventory.go
│       │   ├── sales_order/
│       │   │   └── sales_order.go
│       │   ├── purchase_order/
│       │   │   └── purchase_order.go
│       │   ├── ar/
│       │   │   └── ar.go
│       │   ├── ap/
│       │   │   └── ap.go
│       │   ├── gl/
│       │   │   └── gl.go
│       │   └── wms/
│       │       └── wms.go
│       ├── routes/
│       │   ├── customer/
│       │   │   ├── router.go
│       │   │   ├── create.go
│       │   │   ├── get.go
│       │   │   ├── update.go
│       │   │   └── delete.go
│       │   ├── vendor/
│       │   ├── product/
│       │   ├── inventory/
│       │   ├── sales_order/
│       │   ├── purchase_order/
│       │   ├── ar/
│       │   ├── ap/
│       │   ├── gl/
│       │   └── wms/
│       ├── services/
│       │   ├── customer/
│       │   │   └── customer.go
│       │   ├── vendor/
│       │   │   └── vendor.go
│       │   ├── product/
│       │   │   └── product.go
│       │   ├── inventory/
│       │   │   └── inventory.go
│       │   ├── sales_order/
│       │   │   └── sales_order.go
│       │   ├── purchase_order/
│       │   │   └── purchase_order.go
│       │   ├── ar/
│       │   │   ├── invoice.go
│       │   │   └── payment.go
│       │   ├── ap/
│       │   │   ├── invoice.go
│       │   │   └── payment.go
│       │   ├── gl/
│       │   │   ├── account.go
│       │   │   ├── transaction.go
│       │   │   └── journal.go
│       │   └── wms/
│       │       ├── location.go
│       │       ├── put_away.go
│       │       ├── replenishment.go
│       │       └── cycle_count.go
│       └── models/
│           ├── customer.go
│           ├── vendor.go
│           ├── product.go
│           ├── inventory.go
│           ├── sales_order.go
│           ├── purchase_order.go
│           ├── ar.go
│           ├── ap.go
│           ├── gl.go
│           └── wms.go
└── sql/
    ├── erp_schema.sql
    ├── erp_indexes.sql
    ├── erp_triggers.sql
    └── erp_seed_data.sql
```

---

**Document prepared for:** ERP System Development Project  
**Based on:** UMS Codebase Analysis & Client Requirements (CamScanner 01-09-2026 18.34.pdf)

