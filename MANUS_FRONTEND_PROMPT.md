# ERP System Frontend Development - Complete Project Brief

## Project Overview

Build a modern, responsive web-based frontend for a complete ERP (Enterprise Resource Planning) system. The backend is built with Go and provides RESTful APIs for all operations.

**Project Name:** FoodHive ERP
**Industry:** Food Distribution / Wholesale
**Backend:** Go with Chi Router, PostgreSQL
**Authentication:** JWT Token-based

---

## Tech Stack Requirements

### Recommended Frontend Stack:
- **Framework:** Next.js 14+ (App Router) OR React 18+ with Vite
- **UI Library:** shadcn/ui OR Ant Design (for ERP-style components)
- **State Management:** Zustand or TanStack Query (React Query)
- **Forms:** React Hook Form + Zod validation
- **Charts:** Recharts or Chart.js (for dashboards/reports)
- **Tables:** TanStack Table (for data grids with sorting/filtering)
- **Icons:** Lucide React
- **Styling:** Tailwind CSS

### Design Requirements:
- **Theme:** Professional, clean, enterprise-grade
- **Color Scheme:** Blue primary (trust/corporate) with accent colors for status indicators
- **Responsive:** Desktop-first (most users on desktop), but mobile-friendly
- **Dark Mode:** Support dark/light theme toggle

---

## API Configuration

### Base URL:
```
Development: http://localhost:8080/api/v1
Production: https://api.foodhive.com/api/v1
```

### Authentication:
All API calls (except `/login`) require JWT token in header:
```
Authorization: Bearer <jwt_token>
```

### Login Endpoint:
```
POST /login
Body: { "email": "string", "password": "string" }
Response: { "token": "jwt_token", "employee": { ... } }
```

### Standard Response Format:
```json
{
  "data": { ... },
  "message": "Success message",
  "success": true
}
```

### Error Response:
```json
{
  "error": "Error message",
  "success": false
}
```

### Pagination:
```
GET /endpoint?page=1&per_page=20
Response includes: { "data": [...], "total": 100, "page": 1, "per_page": 20 }
```

---

## Complete API Reference

### 1. AUTHENTICATION
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/login` | User login, returns JWT token |

---

### 2. EMPLOYEES (User Management)
**Base:** `/employees`

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/create` | Create employee |
| GET | `/get/{id}` | Get by ID |
| PUT | `/update/{id}` | Update employee |
| DELETE | `/delete/{id}` | Delete employee |
| GET | `/list` | List employees (with filters) |

**Create Employee Request:**
```json
{
  "email": "john@example.com",
  "password": "securepassword",
  "full_name": "John Doe",
  "phone": "+1234567890",
  "department_id": 1,
  "role_id": 1,
  "is_active": true
}
```

---

### 3. DEPARTMENTS
**Base:** `/departments`

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/create` | Create department |
| GET | `/get/{id}` | Get by ID |
| PUT | `/update/{id}` | Update department |
| DELETE | `/delete/{id}` | Delete department |
| GET | `/list` | List departments |

---

### 4. ROLES & PERMISSIONS
**Base:** `/roles`

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/create` | Create role |
| GET | `/get/{id}` | Get role with permissions |
| PUT | `/update/{id}` | Update role |
| DELETE | `/delete/{id}` | Delete role |
| GET | `/list` | List roles |
| POST | `/{roleId}/permissions` | Assign permissions to role |
| GET | `/pages` | List all pages/modules |
| POST | `/pages` | Create new page/module |

**Role with Permissions:**
```json
{
  "id": 1,
  "name": "Sales Manager",
  "description": "Manages sales team",
  "permissions": [
    { "page_id": 1, "page_name": "Sales Orders", "can_create": true, "can_view": true, "can_update": true, "can_delete": false }
  ]
}
```

---

### 5. CUSTOMERS
**Base:** `/customers`

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/create` | Create customer |
| GET | `/get/{id}` | Get customer with details |
| GET | `/code/{code}` | Get by customer code |
| PUT | `/update/{id}` | Update customer |
| DELETE | `/delete/{id}` | Delete customer |
| GET | `/list` | List customers |
| GET | `/{id}/order-guide` | Get customer order guide |
| POST | `/{id}/ship-to` | Add ship-to address |

**Customer Model:**
```json
{
  "id": 1,
  "customer_code": "CUST001",
  "name": "ABC Restaurant",
  "credit_limit": 10000.00,
  "payment_terms_days": 30,
  "currency": "USD",
  "sales_rep_id": 5,
  "default_warehouse_id": 1,
  "tax_exempt": false,
  "ship_to_addresses": [...]
}
```

---

### 6. VENDORS
**Base:** `/vendors`

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/create` | Create vendor |
| GET | `/get/{id}` | Get vendor |
| PUT | `/update/{id}` | Update vendor |
| DELETE | `/delete/{id}` | Delete vendor |
| GET | `/list` | List vendors |

---

### 7. WAREHOUSES
**Base:** `/warehouses`

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/create` | Create warehouse |
| GET | `/get/{id}` | Get warehouse |
| PUT | `/update/{id}` | Update warehouse |
| DELETE | `/delete/{id}` | Delete warehouse |
| GET | `/list` | List warehouses |

---

### 8. PRODUCTS
**Base:** `/products`

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/create` | Create product |
| GET | `/get/{id}` | Get product with details |
| GET | `/sku/{sku}` | Get by SKU |
| GET | `/barcode/{barcode}` | Get by barcode |
| PUT | `/update/{id}` | Update product |
| DELETE | `/delete/{id}` | Delete product |
| GET | `/list` | List products |

**Categories:**
| POST | `/categories/create` | Create category |
| GET | `/categories/get/{id}` | Get category |
| PUT | `/categories/update/{id}` | Update category |
| DELETE | `/categories/delete/{id}` | Delete category |
| GET | `/categories/list` | List categories |

**Units of Measure:**
| POST | `/units/add` | Add unit conversion |
| GET | `/units/product/{productId}` | Get product units |
| PUT | `/units/update/{id}` | Update unit |
| DELETE | `/units/delete/{id}` | Delete unit |

**Product Model:**
```json
{
  "id": 1,
  "sku": "BEEF001",
  "barcode": "123456789012",
  "upc": "012345678901",
  "name": "Ground Beef",
  "description": "Fresh ground beef",
  "category_id": 5,
  "base_unit": "LB",
  "is_catch_weight": true,
  "catch_weight_unit": "LB",
  "country_of_origin": "USA",
  "shelf_life_days": 14,
  "is_lot_tracked": true,
  "haccp_category": "MEAT",
  "qc_required": true
}
```

---

### 9. INVENTORY
**Base:** `/inventory`

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/get/{id}` | Get inventory record |
| GET | `/product/{productId}` | Get inventory by product |
| GET | `/warehouse/{warehouseId}` | Get inventory by warehouse |
| GET | `/lot/{lotNumber}` | Get inventory by lot |
| GET | `/list` | List inventory with filters |
| GET | `/summary/product/{productId}` | Get product summary |
| GET | `/expiring` | Get expiring inventory |
| POST | `/receive` | Receive inventory |
| POST | `/adjust` | Adjust inventory |
| POST | `/transfer` | Transfer inventory |
| GET | `/transactions` | Get transaction history |

**List Filters:**
```
?product_id=1&warehouse_id=1&lot_number=LOT001&expiring_within_days=30&page=1&per_page=20
```

**Inventory Summary Response:**
```json
{
  "product_id": 1,
  "product_name": "Ground Beef",
  "total_on_hand": 500.00,
  "total_allocated": 50.00,
  "total_on_order": 100.00,
  "total_available": 450.00,
  "warehouses": [
    { "warehouse_id": 1, "warehouse_name": "Main", "quantity": 300.00 }
  ]
}
```

---

### 10. PURCHASE ORDERS
**Base:** `/purchase-orders`

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/create` | Create PO |
| GET | `/get/{id}` | Get PO with lines |
| GET | `/number/{poNumber}` | Get by PO number |
| PUT | `/update/{id}` | Update PO |
| DELETE | `/delete/{id}` | Delete PO |
| GET | `/list` | List POs |
| POST | `/submit/{id}` | Submit PO |
| POST | `/cancel/{id}` | Cancel PO |
| POST | `/{poId}/lines` | Add line |
| PUT | `/lines/{lineId}` | Update line |
| DELETE | `/lines/{lineId}` | Delete line |

**Receiving:**
| POST | `/receive` | Create receiving |
| GET | `/receiving/{id}` | Get receiving |
| GET | `/receivings` | List receivings |

**PO Status Flow:** `DRAFT` → `SUBMITTED` → `APPROVED` → `PARTIAL` → `RECEIVED` / `CANCELLED`

**Create PO Request:**
```json
{
  "vendor_id": 1,
  "warehouse_id": 1,
  "expected_date": "2026-01-20",
  "buyer_id": 5,
  "notes": "Weekly order",
  "lines": [
    { "product_id": 1, "quantity": 100, "unit_cost": 5.99, "unit": "LB" }
  ]
}
```

---

### 11. SALES ORDERS
**Base:** `/sales-orders`

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/create` | Create order |
| GET | `/get/{id}` | Get order with lines |
| GET | `/number/{orderNumber}` | Get by order number |
| PUT | `/update/{id}` | Update order |
| DELETE | `/delete/{id}` | Delete order |
| GET | `/list` | List orders |
| POST | `/confirm/{id}` | Confirm order |
| POST | `/cancel/{id}` | Cancel order |
| POST | `/ship/{id}` | Mark as shipped |
| POST | `/{orderId}/lines` | Add line |
| PUT | `/lines/{lineId}` | Update line |
| DELETE | `/lines/{lineId}` | Delete line |
| GET | `/order-guide/{customerId}` | Get customer order guide |
| POST | `/lost-sale` | Record lost sale |
| GET | `/lost-sales` | Get lost sales |

**Order Types:** `STANDARD`, `QUOTE`, `CREDIT_MEMO`, `PICKUP`, `ADVANCE`, `PREPAID`, `ON_HOLD`

**Order Status Flow:** `DRAFT` → `CONFIRMED` → `PICKING` → `SHIPPED` → `INVOICED` / `CANCELLED`

**Create Order Request:**
```json
{
  "customer_id": 1,
  "ship_to_id": 1,
  "order_type": "STANDARD",
  "requested_ship_date": "2026-01-15",
  "warehouse_id": 1,
  "route_id": 3,
  "po_number": "CUST-PO-123",
  "notes": "Leave at back door",
  "lines": [
    { "product_id": 1, "quantity": 50, "unit": "LB", "unit_price": 7.99 }
  ]
}
```

---

### 12. PICKING & ROUTING
**Base:** `/picking`

**Routes:**
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/routes/create` | Create route |
| GET | `/routes/get/{id}` | Get route with stops |
| PUT | `/routes/update/{id}` | Update route |
| DELETE | `/routes/delete/{id}` | Delete route |
| GET | `/routes/list` | List routes |
| POST | `/routes/{routeId}/stops` | Add stop |
| PUT | `/routes/stops/{stopId}` | Update stop |
| DELETE | `/routes/stops/{stopId}` | Delete stop |
| POST | `/routes/{routeId}/reorder` | Reorder stops |

**Pick Lists:**
| POST | `/create` | Create pick list |
| POST | `/generate` | Generate from route/date |
| GET | `/get/{id}` | Get pick list |
| GET | `/list` | List pick lists |
| POST | `/{id}/start` | Start picking |
| POST | `/{id}/complete` | Complete picking |
| POST | `/{id}/cancel` | Cancel pick list |
| GET | `/{id}/lines` | Get pick lines |
| POST | `/lines/{lineId}/confirm` | Confirm pick line |

**Reports:**
| GET | `/master-pick` | Master pick report |
| GET | `/suggested-picking` | Suggested picking (FEFO) |

**Pick Status Flow:** `PENDING` → `IN_PROGRESS` → `COMPLETED` / `CANCELLED`

---

### 13. PRICING
**Base:** `/pricing`

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/lookup` | Get price for product/customer |
| POST | `/lookup/batch` | Batch price lookup |
| GET | `/check-margin` | Check if price below cost |
| POST | `/product/set` | Set product price |
| GET | `/product/{productId}` | Get product prices |
| DELETE | `/product/price/{id}` | Delete product price |
| POST | `/customer/set` | Set customer price |
| GET | `/customer/{customerId}` | Get customer prices |
| DELETE | `/customer/price/{id}` | Delete customer price |
| POST | `/contracts/create` | Create contract |
| GET | `/contracts/get/{id}` | Get contract |
| GET | `/contracts/list` | List contracts |
| POST | `/contracts/{id}/deactivate` | Deactivate contract |
| POST | `/promotions/create` | Create promotion |
| GET | `/promotions/get/{id}` | Get promotion |
| GET | `/promotions/list` | List promotions |
| POST | `/promotions/{id}/deactivate` | Deactivate promotion |
| POST | `/costs/update` | Update product cost |
| GET | `/costs/{productId}` | Get product cost |
| POST | `/mass-update` | Mass price update |
| GET | `/list` | Get price list |

**5-Level Price Hierarchy (highest priority first):**
1. Contract Price (customer-specific contracts)
2. Price Level (customer price level)
3. Customer Price (special pricing)
4. Promotional Price (date-based promotions)
5. Order Price (base product price)

**Price Lookup Request:**
```json
{
  "product_id": 1,
  "customer_id": 5,
  "quantity": 100,
  "as_of_date": "2026-01-15"
}
```

**Price Lookup Response:**
```json
{
  "product_id": 1,
  "price": 6.99,
  "price_source": "CONTRACT",
  "contract_id": 12,
  "list_price": 7.99,
  "discount_percent": 12.5,
  "cost": 4.50,
  "margin_percent": 35.6
}
```

---

### 14. ACCOUNTS RECEIVABLE (AR)
**Base:** `/ar`

**Invoices:**
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/invoices/create` | Create invoice |
| GET | `/invoices/get/{id}` | Get invoice |
| GET | `/invoices/number/{number}` | Get by invoice number |
| GET | `/invoices/list` | List invoices |
| POST | `/invoices/{id}/post` | Post invoice |
| POST | `/invoices/{id}/void` | Void invoice |
| POST | `/invoices/from-order/{orderId}` | Create from sales order |

**Payments:**
| POST | `/payments/create` | Create payment |
| GET | `/payments/get/{id}` | Get payment |
| GET | `/payments/list` | List payments |

**Credit Management:**
| GET | `/credit/{customerId}` | Get customer credit |
| GET | `/credit/{customerId}/check` | Check credit availability |
| PUT | `/credit/{customerId}/limit` | Update credit limit |

**Aging & Reports:**
| GET | `/aging/{customerId}` | Get customer aging |
| GET | `/aging/report` | Get aging report (all customers) |
| GET | `/statement/{customerId}` | Get customer statement |
| GET | `/overdue` | Get overdue invoices |

**Aging Report Response:**
```json
{
  "customer_id": 1,
  "customer_name": "ABC Restaurant",
  "current": 1000.00,
  "days_1_30": 500.00,
  "days_31_60": 250.00,
  "days_61_90": 0.00,
  "over_90": 0.00,
  "total": 1750.00
}
```

---

### 15. ACCOUNTS PAYABLE (AP)
**Base:** `/ap`

**Invoices:**
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/invoices/create` | Create invoice |
| GET | `/invoices/get/{id}` | Get invoice |
| GET | `/invoices/list` | List invoices |
| POST | `/invoices/{id}/approve` | Approve invoice |
| POST | `/invoices/{id}/void` | Void invoice |
| POST | `/invoices/from-receiving/{receivingId}` | Create from receiving |

**Payments:**
| POST | `/payments/create` | Create payment |
| GET | `/payments/get/{id}` | Get payment |
| GET | `/payments/list` | List payments |
| POST | `/payments/{id}/void` | Void payment |

**Balance & Aging:**
| GET | `/balance/{vendorId}` | Get vendor balance |
| GET | `/aging/{vendorId}` | Get vendor aging |
| GET | `/aging/report` | Get aging report (all vendors) |
| GET | `/due` | Get due invoices |
| GET | `/overdue` | Get overdue invoices |

---

### 16. GENERAL LEDGER (GL)
**Base:** `/gl`

**Chart of Accounts:**
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/accounts` | Create account |
| GET | `/accounts/{id}` | Get account |
| GET | `/accounts/code/{code}` | Get by account code |
| PUT | `/accounts/{id}` | Update account |
| DELETE | `/accounts/{id}` | Delete account |
| GET | `/accounts` | List accounts |
| GET | `/chart-of-accounts` | Get hierarchical COA |

**Fiscal Years:**
| POST | `/fiscal-years` | Create fiscal year (auto-creates 13 periods) |
| GET | `/fiscal-years/{id}` | Get fiscal year |
| GET | `/fiscal-years/current` | Get current fiscal year |
| GET | `/fiscal-years` | List fiscal years |
| POST | `/fiscal-years/{id}/close` | Close fiscal year |

**Periods:**
| GET | `/periods/{id}` | Get period |
| GET | `/periods/current` | Get current period |
| GET | `/fiscal-years/{fiscalYearId}/periods` | List periods |
| POST | `/periods/{id}/close` | Close period |
| POST | `/periods/{id}/reopen` | Reopen period |

**Journal Entries:**
| POST | `/journal-entries` | Create journal entry |
| GET | `/journal-entries/{id}` | Get with lines |
| PUT | `/journal-entries/{id}` | Update (draft only) |
| DELETE | `/journal-entries/{id}` | Delete (draft only) |
| GET | `/journal-entries` | List entries |
| POST | `/journal-entries/{id}/post` | Post to GL |
| POST | `/journal-entries/{id}/reverse` | Create reversal |
| POST | `/journal-entries/{id}/void` | Void entry |

**Financial Reports:**
| GET | `/reports/trial-balance` | Trial Balance |
| GET | `/reports/income-statement` | Income Statement (P&L) |
| GET | `/reports/balance-sheet` | Balance Sheet |
| GET | `/reports/account-activity/{accountId}` | Account Ledger |

**Account Types:** `ASSET`, `LIABILITY`, `EQUITY`, `REVENUE`, `EXPENSE`

**Journal Entry Status:** `DRAFT` → `POSTED` / `VOIDED`

**Create Journal Entry:**
```json
{
  "entry_date": "2026-01-15",
  "entry_type": "STANDARD",
  "description": "Monthly rent",
  "lines": [
    { "account_id": 50, "description": "Rent expense", "debit_amount": 5000.00, "credit_amount": 0 },
    { "account_id": 10, "description": "Cash payment", "debit_amount": 0, "credit_amount": 5000.00 }
  ]
}
```

---

### 17. CATCH WEIGHT
**Base:** `/catch-weight`

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/capture` | Capture catch weight (with pieces) |
| POST | `/capture/quick` | Quick capture (total only) |
| POST | `/entries/{entryId}/pieces` | Add piece weight |
| PUT | `/pieces/{pieceId}` | Update piece weight |
| DELETE | `/pieces/{pieceId}` | Delete piece |
| GET | `/entries/{id}` | Get entry by ID |
| GET | `/entries` | List entries |
| GET | `/entries/{entryId}/pieces` | Get pieces for entry |
| GET | `/reference/{refType}/{refId}/product/{productId}` | Get by reference |
| GET | `/products/{productId}/config` | Get product config |
| PUT | `/products/{productId}/config` | Update product config |
| GET | `/reports/variance` | Variance report |
| GET | `/reports/lot/{productId}/{lotNumber}` | Lot summary |
| POST | `/billing/adjustment` | Calculate billing adjustment |
| POST | `/entries/{entryId}/mark-billed` | Mark as billed |
| POST | `/validate-weight` | Validate weight |

**Catch Weight Entry:**
```json
{
  "product_id": 1,
  "reference_type": "SALES_ORDER",
  "reference_id": 123,
  "lot_number": "LOT2026-001",
  "expected_weight": 50.00,
  "actual_weight": 52.35,
  "pieces": [
    { "piece_number": 1, "weight": 17.45 },
    { "piece_number": 2, "weight": 17.55 },
    { "piece_number": 3, "weight": 17.35 }
  ]
}
```

---

### 18. PAYROLL
**Base:** `/payroll`

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/create` | Create payroll batch |
| GET | `/get/{id}` | Get payroll with lines |
| PUT | `/update/{id}` | Update payroll |
| DELETE | `/delete/{id}` | Delete payroll |
| GET | `/list` | List payrolls |
| POST | `/{id}/lines` | Add employee line |
| DELETE | `/lines/{lineId}` | Remove line |
| POST | `/{id}/calculate` | Calculate payroll |
| POST | `/{id}/approve` | Approve payroll |

---

## Application Screens & Components

### DASHBOARD (Home)
**Priority: HIGH**
- **Today's Summary Cards:**
  - Pending Sales Orders count
  - Orders to Ship today
  - Pending Pick Lists
  - Overdue AR amount
  - Due AP this week
  - Low Stock Alerts count
- **Quick Actions:**
  - New Sales Order
  - New Purchase Order
  - Record Payment
- **Charts:**
  - Sales trend (last 30 days)
  - AR Aging pie chart
  - Top 5 Products sold

### MASTER DATA

#### Customers Screen
- **List View:** Searchable data table with filters
  - Columns: Code, Name, Credit Limit, Balance, Sales Rep, Status
  - Actions: View, Edit, Delete
- **Detail View:**
  - Customer info card
  - Ship-to addresses list
  - Order guide tab
  - Order history tab
  - AR aging tab
  - Credit status indicator

#### Vendors Screen
- **List View:** Same pattern as Customers
- **Detail View:**
  - Vendor info
  - PO history
  - AP aging

#### Products Screen
- **List View:**
  - Columns: SKU, Name, Category, Stock, Cost, Price
  - Filters: Category, Active/Inactive
  - Quick stock view
- **Detail View:**
  - Product info
  - Pricing tab (all price levels)
  - Inventory tab (by warehouse)
  - Units of measure
  - Catch weight config (if applicable)

### TRANSACTIONS

#### Sales Order Screen
**Priority: HIGH**
- **List View:**
  - Status filter tabs: All, Draft, Confirmed, Shipping, Shipped
  - Columns: Order#, Date, Customer, Total, Status
- **Order Entry Form:**
  - Customer selector (with credit check warning)
  - Ship-to dropdown
  - Order type dropdown
  - Line items grid:
    - Product search (by SKU, barcode, name)
    - Order guide quick-add
    - Quantity, Unit, Price, Amount
    - Below-cost warning indicator
  - Totals section
  - Notes field
- **Actions:** Save Draft, Confirm, Cancel, Print, Create Invoice

#### Purchase Order Screen
- **List View:** Similar to Sales Orders
- **PO Entry Form:**
  - Vendor selector
  - Warehouse selector
  - Expected date
  - Line items with vendor product lookup
- **Receiving:**
  - Select PO
  - Enter received quantities
  - Lot number entry
  - Variance display

#### Pick List Screen
- **List View:**
  - Filter by date, route, status
  - Columns: Pick#, Route, Orders, Items, Status
- **Pick Detail:**
  - Order/stop list
  - Pick lines grouped by location (for zone picking)
  - Confirm pick quantities
  - Catch weight entry integration
- **Master Pick Report:** Consolidated view across orders

### FINANCIAL

#### AR Dashboard
- **Summary Cards:** Total Receivable, Current, 30, 60, 90, 120+ days
- **Aging Chart:** Stacked bar by customer
- **Quick Payment Entry**
- **Overdue Alert List**

#### AR Invoices
- **List View:** Filter by customer, status, date range
- **Invoice Detail:** Lines, payments, balance
- **Create from Order:** Select completed orders to invoice

#### AR Payments
- **Payment Entry Form:**
  - Customer selector
  - Payment method
  - Amount
  - Apply to invoices (checkbox grid)
- **Payment List**

#### AP Dashboard (Mirror of AR)
- Same pattern for vendor payables

#### General Ledger
- **Chart of Accounts:**
  - Tree view (hierarchical)
  - Account detail modal
- **Journal Entry Screen:**
  - Entry form with debit/credit lines
  - Balance check (must equal)
  - Post/Void actions
- **Reports:**
  - Trial Balance (export to Excel/PDF)
  - Income Statement
  - Balance Sheet
  - Account Activity drill-down

### OPERATIONS

#### Catch Weight Screen
**Priority: MEDIUM**
- **Weight Capture Form:**
  - Select reference (SO, PO, Pick List)
  - Scan/enter product
  - Expected vs Actual weight display
  - Piece-by-piece entry option
  - Variance indicator (green/yellow/red)
- **Reports:**
  - Variance report by product/date
  - Lot summary

### SETTINGS/ADMIN

#### User Management
- Employee list
- Create/Edit employee form
- Role assignment
- Activity log

#### Role & Permissions
- Role list
- Permission matrix (checkboxes):
  - Rows: Modules (Sales Orders, Purchase Orders, etc.)
  - Columns: Create, View, Update, Delete

#### Pricing Management
- Product prices
- Customer prices
- Contracts
- Promotions
- Mass price update tool

---

## UI Component Requirements

### Data Tables
- Sortable columns
- Column visibility toggle
- Filter row or sidebar filters
- Pagination (10, 20, 50, 100 per page)
- Row selection (single/multi)
- Export (CSV, Excel)
- Print view

### Forms
- Field validation with error messages
- Required field indicators
- Auto-save drafts (localStorage)
- Confirmation dialogs for destructive actions

### Search/Select Components
- Debounced search input
- Async dropdown with search
- Barcode scanner integration (optional)
- Recently used items

### Status Indicators
- Draft: Gray
- Pending/Confirmed: Blue
- In Progress: Yellow
- Completed/Shipped: Green
- Cancelled/Voided: Red
- Overdue: Orange/Red

### Notifications
- Toast messages for actions
- Bell icon with notification dropdown
- Unread count badge

---

## Business Rules to Implement in Frontend

1. **Credit Check:** Show warning when customer exceeds credit limit on order entry
2. **Below Cost Warning:** Highlight when selling below product cost
3. **Stock Check:** Show available vs ordered quantity
4. **Lot Tracking:** For lot-tracked products, require lot selection
5. **Catch Weight:** For catch weight products, require weight entry
6. **Expiry Check:** Warn when picking products near expiry
7. **Balance Validation:** Journal entries must balance (debits = credits)
8. **Period Check:** Prevent posting to closed periods
9. **Sequential Flow:** Enforce order status flow (can't skip states)
10. **Role Permissions:** Hide/disable features based on user permissions

---

## Responsive Breakpoints

- **Desktop:** 1280px+ (full features)
- **Tablet:** 768px-1279px (simplified layout)
- **Mobile:** <768px (essential features only)

---

## Accessibility

- Keyboard navigation
- ARIA labels
- Screen reader support
- High contrast mode support
- Focus indicators

---

## Performance Requirements

- First contentful paint: <2s
- Time to interactive: <4s
- API response caching (React Query)
- Infinite scroll for large lists
- Lazy loading for routes/components

---

## Deliverables

1. Complete React/Next.js application
2. All screens listed above
3. API integration layer
4. Authentication flow with protected routes
5. Responsive design
6. Dark/light theme toggle
7. Basic error handling and loading states
8. Print-friendly views for reports/invoices

---

## Sample Data

Create mock data seeder or use the following test accounts:

**Admin User:**
- Email: admin@foodhive.com
- Password: Admin123!

**Sales User:**
- Email: sales@foodhive.com
- Password: Sales123!

---

## Questions to Ask During Development

1. What is the primary color scheme preference?
2. Any specific branding/logo to include?
3. Preferred date/number format (locale)?
4. Any specific third-party integrations needed?
5. Offline support requirements?
6. Mobile app needed in future?

---

Good luck building! This system covers 90% of a standard food distribution ERP. Focus on the core flows first: **Order Entry → Picking → Invoice → Payment**, then add the supporting modules.
