# Requirements Gap Analysis
## Based on DocText.md vs Current Implementation

This document analyzes the gaps between the requirements documented in `DocText.md` (from the provided forms/images) and the current codebase implementation.

---

## 1. Customer Form - Missing Fields

### Current Implementation
The `customers` table currently has:
- `id`, `customer_code`, `name`
- `billing_address_id`, `credit_limit`, `current_balance`
- `payment_terms_days`, `currency`
- `sales_rep_id`, `default_route_id`, `default_warehouse_id`
- `tax_exempt`, `is_active`, `created_by`, `created_at`, `updated_at`

### Required Fields from DocText.md (Section 3)

#### Basic Information
- ✅ `ID` - Auto-generated (exists)
- ✅ `Name` (ຊື່) - exists
- ✅ `Customer Code` (ລະຫັດລູກຄ້າ) - exists
- ❌ **`Default Contact Name`** (ຊື່ຜູ້ຕິດຕໍ່ໃຊ້ເປັນຄ່າເລີ່ມຕົ້ນ) - MISSING
- ❌ **`Full Address`** (ທີ່ຢູ່ເຕັມ) - MISSING (only has `billing_address_id`)

#### Delivery & Contact Information
- ❌ **`Delivery Name`** (ຊື່ສະຖານທີ່ຈັດສົ່ງ) - MISSING
- ❌ **`Tel`** (ໂທລະສັບ) - MISSING
- ❌ **`Fax`** (ແຟັກ) - MISSING
- ❌ **`Mobile`** (ມືຖື) - MISSING
- ❌ **`Tax`** (ອາກອນ) - MISSING (only has `tax_exempt` boolean)

#### Account Information
- ❌ **`Birthday`** (ວັນເກີດ) - MISSING
- ❌ **`Web Email`** (ອີເມວເວັບ) - MISSING
- ❌ **`Price Type`** (ປະເພດລາຄາ) - MISSING (dropdown: ລາຄາຂາຍ 1)
- ❌ **`Level`** (ລະດັບ) - MISSING (dropdown: membership level)
- ❌ **`Bill Discount`** (ສ່ວນຫຼຸດໃບບິນ) - MISSING (percentage, default 0.00%)
- ✅ `Limit Credit Money` (ຈໍາກັດເງິນເຊື່ອ) - exists as `credit_limit`
- ✅ `Comment` (ໝາຍເຫດ) - could use `notes` field if exists

#### Membership & Identity
- ❌ **`Member Expire`** (ໝົດອາຍຸສະມາຊິກ) - MISSING (date)
- ❌ **`Credit`** (ເງິນເຊື່ອ) - MISSING (number, default 0.00) - different from credit_limit
- ❌ **`Passport No.`** (ເລກໜັງສືຜ່ານແດນ) - MISSING
- ❌ **`Country Code`** (ລະຫັດປະເທດ) - MISSING
- ❌ **`Gender`** (ເພດ) - MISSING (dropdown)

#### Settings & Preferences
- ❌ **`Use Promotion`** (ໃຊ້ໂປຣໂມຊັ່ນ) - MISSING (checkbox)
- ❌ **`Use PMT Bill`** (ໃຊ້ໃບບິນ PMT) - MISSING (checkbox)
- ❌ **`Collect Point`** (ຈຸດສະສົມແຕ້ມ) - MISSING (checkbox)
- ❌ **`Payee`** (ຜູ້ຮັບເງິນ) - MISSING
- ❌ **`RD Branch`** (ສາຂາ RD) - MISSING
- ❌ **`Barcode`** (ບາໂຄດ) - MISSING
- ❌ **`EMail`** (ອີເມວ) - MISSING

#### Additional Features
- ❌ **`Smart Card Reader`** - Feature not implemented

---

## 2. Supplier/Vendor Form - Missing Fields

### Current Implementation
The `vendors` table currently has:
- `id`, `vendor_code`, `name`
- `address_line1`, `address_line2`, `city`, `state`, `postal_code`, `country`
- `phone`, `email`
- `payment_terms_days`, `currency`, `lead_time_days`, `minimum_order`
- `buyer_id`, `is_active`, `created_at`, `updated_at`

### Required Fields from DocText.md (Section 7)

#### Basic Information
- ✅ `Supplier ID` - Auto-generated (exists)
- ✅ `Name` (ຊື່) - exists
- ✅ `Supplier Code` (ລະຫັດຜູ້ສະໜອງ) - exists
- ❌ **`Contact Name`** (ຊື່ຜູ້ຕິດຕໍ່) - MISSING
- ❌ **`RD Branch`** (ສາຂາ RD) - MISSING
- ✅ `Address` (ທີ່ຢູ່) - exists (address_line1, etc.)

#### Additional Information
- ❌ **`Picture Name`** (ຊື່ຮູບພາບ) - MISSING
- ✅ `Tel` (ໂທລະສັບ) - exists as `phone`
- ❌ **`Mobile`** (ມືຖື) - MISSING
- ❌ **`Fax`** (ແຟັກ) - MISSING
- ❌ **`TaxCode`** (ລະຫັດອາກອນ) - MISSING

#### Financial Settings
- ❌ **`Credit`** (ເງິນເຊື່ອ) - MISSING (number, default 0, unit: days)
- ❌ **`Discount`** (ສ່ວນຫຼຸດ) - MISSING (percentage, default 0.00%)
- ❌ **`PO Credit Money`** (ວົງເງິນເຊື່ອ PO) - MISSING (number, default 0.00)
- ✅ `Currency` (ສະກຸນເງິນ) - exists (default: KIP in form, USD in code)

---

## 3. Purchase Order Form - Missing Fields

### Current Implementation
The `purchase_orders` table currently has:
- `id`, `po_number`, `vendor_id`, `warehouse_id`
- `order_date`, `expected_date`, `received_date`
- `status`, `subtotal`, `tax_amount`, `freight_amount`, `total_amount`
- `notes`, `buyer_id`, `created_by`, `created_at`, `updated_at`

The `purchase_order_lines` table has:
- `id`, `po_id`, `line_number`, `product_id`, `description`
- `quantity_ordered`, `quantity_received`, `unit_of_measure`, `unit_cost`, `line_total`
- `expected_date`

### Required Fields from DocText.md (Section 1)

#### Left Section (Supplier/Customer Info)
- ❌ **`Contact Person`** (ຜູ້ຕິດຕໍ່) - MISSING
- ✅ `Supplier Code` (ລະຫັດຜູ້ສະໜອງ) - exists via `vendor_id`
- ❌ **`Customer Name`** (ຊື່ບໍລິສັດ) - MISSING (on PO form)
- ❌ **`Address`** (ທີ່ຢູ່) - MISSING (on PO form)
- ❌ **`Tel`** (ໂທລະສັບ) - MISSING (on PO form)
- ❌ **`E-mail`** (ອີເມວ) - MISSING (on PO form)

#### Right Section (Document Details)
- ✅ `Document No` (ເລກທີເອກະສານ) - exists as `po_number`
- ✅ `Document Date` (ວັນທີເອກະສານ) - exists as `order_date`
- ❌ **`Document Ref`** (ເອກະສານອ້າງອີງ) - MISSING
- ❌ **`Terms of Payment`** (ເງື່ອນໄຂການຊຳລະເງິນ) - MISSING
- ✅ `Currency` (ສະກຸນເງິນ) - exists

#### Line Items Table
- ✅ `Item` (ລໍາດັບ) - exists as `line_number`
- ✅ `Product Code` (ລະຫັດສິນຄ້າ) - exists via `product_id`
- ✅ `Description` (ລາຍການ) - exists
- ✅ `Qty` (ຈໍານວນ) - exists as `quantity_ordered`
- ✅ `Unit` (ຫົວໜ່ວຍ) - exists as `unit_of_measure`
- ✅ `Unit Price` (ລາຄາ) - exists as `unit_cost`
- ❌ **`Discount`** (ສ່ວນຫຼຸດ) - MISSING (per line item)
- ✅ `Amount` (ຈໍານວນເງິນ) - exists as `line_total`

#### Summary Section
- ✅ `Amount before Trade Discount` - exists as `subtotal`
- ❌ **`Trade Discount`** (ສ່ວນຫຼຸດການຄ້າ) - MISSING
- ❌ **`Amount after Trade Discount`** (ລາຄາຫຼັງສ່ວນຫຼຸດ) - MISSING
- ✅ `VAT 7%` (ອາກອນມູນຄ່າເພີ່ມ) - exists as `tax_amount` but not specifically 7%
- ✅ `Grand Total (VAT included)` - exists as `total_amount`

#### Signature Section
- ❌ **`Prepared By`** (ຜູ້ກະກຽມເອກະສານ) - MISSING (signature field)
- ❌ **`Authorized By`** (ຜູ້ອະນຸມັດ) - MISSING (signature field)
- ❌ **`Check By`** (ຜູ້ກວດເອກະສານ) - MISSING (signature field)

---

## 4. Purchase Requisition Form - COMPLETELY MISSING

### Required from DocText.md (Section 2)

This is a **completely new feature** that doesn't exist in the current system.

#### Required Tables/Models:
- `purchase_requisitions` table with:
  - `id`, `requisition_number`, `document_number`
  - `requester_id` (employee), `department_id`
  - `reason`, `supplier_id` (vendor)
  - `document_date`, `request_date`
  - `total_amount`, `remark`
  - `checked_by`, `authorized_by`, `purchasing_employee_id`
  - `status`, `created_by`, `created_at`, `updated_at`

- `purchase_requisition_lines` table with:
  - `id`, `requisition_id`, `line_number`
  - `product_id`, `product_code`, `description`
  - `stock_balance`, `quantity_requested`
  - `unit_of_measure`, `unit_price`, `amount`

#### Workflow:
1. Internal department requests purchase
2. Gets checked by supervisor
3. Gets authorized by manager
4. Goes to purchasing department
5. Can be converted to Purchase Order

---

## 5. Clear Advance Request Form - COMPLETELY MISSING

### Required from DocText.md (Section 6)

This is a **completely new feature** that doesn't exist in the current system.

#### Required Tables/Models:
- `advance_requests` table with:
  - `id`, `request_number`
  - `employee_id` (name), `position`
  - `po_number` (optional reference)
  - `request_date`
  - `currency` (LAK, THB, USD)
  - `total_amount`
  - `description` (can be multi-line with sub-items)
  - `exchange_rate`
  - `status`, `created_by`, `created_at`, `updated_at`

- `advance_request_lines` table with:
  - `id`, `request_id`, `line_number`
  - `description`, `amount`
  - `parent_line_id` (for sub-items like "A:", "B:")

#### Features:
- Multi-currency support (LAK, THB, USD)
- Hierarchical line items (main items with sub-items)
- Exchange rate tracking
- Approval workflow (signatures)

---

## 6. Clear Advance Voucher Form - COMPLETELY MISSING

### Required from DocText.md (Section 5)

This is a **completely new feature** that doesn't exist in the current system.

#### Required Tables/Models:
- `advance_vouchers` table with:
  - `id`, `voucher_number`
  - `advance_request_id` (links to advance request)
  - `employee_id` (name), `position`
  - `po_number` (optional reference)
  - `voucher_date`
  - `currency` (LAK, THB, USD)
  - `advance_amount`, `expenditure_amount`
  - `balance_amount`
  - `status`, `created_by`, `created_at`, `updated_at`

- `advance_voucher_lines` table with:
  - `id`, `voucher_id`, `line_number`
  - `line_type` (ADVANCE, EXPEND, SUB_EXPEND)
  - `description`, `amount`
  - `parent_line_id` (for sub-items under EXPEND)

- `advance_voucher_documents` table with:
  - `id`, `voucher_id`
  - `document_type` (RECEIPT, INVOICE, BANK_TRANSFER, PO_PR, TRANSPORTATION)
  - `has_document` (boolean)
  - `document_path` (optional file path)

- `advance_voucher_signatures` table with:
  - `id`, `voucher_id`
  - `role` (ACCOUNTANT, RETURNED_BY)
  - `employee_id`, `signature_date`

#### Features:
- Links to Advance Request (same request number)
- Tracks advance amount vs expenditure
- Supporting documents checklist (Receipt, Invoice, Bank Transfer, PO/PR, Transportation)
- Signature tracking (Accountant, Returned by)
- Multi-currency support

---

## 7. Stock/Inventory Management Screen

### Current Implementation
The system has:
- `inventory` table with product, warehouse, location, lot, quantity tracking
- `inventory_transactions` for movement tracking
- `products` table with basic product info

### Required Features from DocText.md (Section 4)

#### Missing Display Fields:
- ❌ **`Price 1`** (ລາຄາ 1) - MISSING in display
- ❌ **`Price 2`** (ລາຄາ 2) - MISSING in display
- ❌ **`%VAT`** - MISSING in display
- ✅ `Expire` - exists as `expiry_date`
- ✅ `Lot Number` - exists

#### Missing Filter Options:
- ✅ `Category` (ໝວດໝູ່) - exists via `category_id`
- ❌ **`Product Group`** (ກຸ່ມສິນຄ້າ) - MISSING
- ❌ **`Brand`** (ຍີ່ຫໍ້) - MISSING
- ✅ `Branch Name` (ຊື່ສາຂາ) - exists via `warehouse_id`

#### Missing Options:
- ❌ **`Set Buy Price`** (ຕັ້ງລາຄາຊື້) - MISSING
- ❌ **`Show Num Order`** (ສະແດງລໍາດັບ) - MISSING
- ❌ **`Expand Group`** (ຂະຫຍາຍກຸ່ມ) - MISSING

#### Missing Actions:
- ❌ **`Edit Sell Price`** (ແກ້ໄຂລາຄາຂາຍ) - MISSING
- ❌ **`More Detail`** (ລາຍລະອຽດເພີ່ມເຕີມ) - MISSING
- ❌ **`Edit Location`** (ແກ້ໄຂສະຖານທີ່) - MISSING
- ❌ **`Group View`** (ເບິ່ງແບບກຸ່ມ) - MISSING
- ✅ `Print` (ພິມ) - likely exists
- ❌ **`Print Barcode`** (ພິມບາໂຄດ) - MISSING

---

## Summary of Missing Features

### High Priority (Core Business Functions)
1. **Purchase Requisition** - Complete new module
2. **Clear Advance Request** - Complete new module
3. **Clear Advance Voucher** - Complete new module
4. **Customer Form** - Missing ~20 fields
5. **Supplier Form** - Missing ~8 fields
6. **Purchase Order** - Missing ~10 fields including signatures

### Medium Priority (Enhanced Features)
1. **Stock Management** - Missing price displays, filters, and actions
2. **Product** - Missing Brand and Product Group fields
3. **Pricing** - Missing Price Type, Price 1, Price 2, VAT% display

### Low Priority (UI/UX Enhancements)
1. Smart Card Reader for customers
2. Enhanced signature tracking
3. Document attachment support

---

## Implementation Notes

### Database Schema Changes Needed:
1. Add columns to `customers` table (~20 new fields)
2. Add columns to `vendors` table (~8 new fields)
3. Add columns to `purchase_orders` table (~10 new fields)
4. Add `discount` column to `purchase_order_lines` table
5. Create `purchase_requisitions` and `purchase_requisition_lines` tables
6. Create `advance_requests` and `advance_request_lines` tables
7. Create `advance_vouchers`, `advance_voucher_lines`, `advance_voucher_documents`, `advance_voucher_signatures` tables
8. Add `brand` and `product_group` to `products` table
9. Add pricing fields to products/inventory

### Model Changes Needed:
- Update all Go models to include new fields
- Create new models for Purchase Requisition, Advance Request, Advance Voucher
- Update request/response DTOs

### Service Changes Needed:
- Update customer service to handle new fields
- Update vendor service to handle new fields
- Update purchase order service to handle new fields
- Create purchase requisition service
- Create advance request service
- Create advance voucher service

### Frontend Changes Needed:
- Update customer form with all new fields
- Update vendor form with all new fields
- Update purchase order form with all new fields
- Create purchase requisition UI
- Create advance request UI
- Create advance voucher UI
- Update stock/inventory screen with new filters and actions

---

*Analysis completed based on DocText.md requirements vs current codebase*
*Date: January 2026*
