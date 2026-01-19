# Vientiane Cold Storage - Forms & System Documentation

This document contains all forms, templates, and system screens used by Vientiane Cold Storage for registration and data entry processes.

---

## Table of Contents

1. [Purchase Order Form](#1-purchase-order-form)
2. [Purchase Requisition Form](#2-purchase-requisition-form)
3. [Add New Customer Form](#3-add-new-customer-form)
4. [Stock/Inventory Management Screen](#4-stockinventory-management-screen)
5. [Clear Advance Voucher Form](#5-clear-advance-voucher-form)
6. [Clear Advance Request Form](#6-clear-advance-request-form)
7. [Add New Supplier Form](#7-add-new-supplier-form)

---

## 1. Purchase Order Form

**Company:** Vientiane Food & Cold Storage  
**Document Type:** ใบສັ່ງຊື້ / Purchase Order

### Company Header Information

| Field | Details |
|-------|---------|
| Company Name | Vientiane Food & Cold Storage |
| Address | 025/1 ຖະໜົນດອນນົກຂຸ້ມ, ບ້ານ ດອນນົກຂຸ້ມ, ເມືອງ ສີສັດຕະນາກ |
| P.O. Box | 9279, 025/1 Don Nokkhoum Road, Ban Don Nokkhoum, Sisattanak District, Vientiane Capital, Lao PDR |
| Tel | +856-21-48-6060 |
| Fax | (856) 21 486060 |

### Left Section Fields

| Field (Lao) | Field (English) | Description |
|-------------|-----------------|-------------|
| ຜູ້ຕິດຕໍ່ | Contact Person | Person to contact for this order |
| ລະຫັດຜູ້ສະໜອງ | Supplier Code | Unique code for the supplier |
| ຊື່ບໍລິສັດ | Customer Name | Name of the customer/company |
| ທີ່ຢູ່/Address | Address | Full address (e.g., 0207809702) |
| ໂທລະສັບ | Tel | Telephone number |
| ອີເມວ | E-mail | Email address |

### Right Section Fields

| Field (Lao) | Field (English) | Description |
|-------------|-----------------|-------------|
| ເລກທີເອກະສານ | Document No | Unique document number |
| ວັນທີເອກະສານ | Document Date | Date of document creation |
| ເອກະສານອ້າງອີງ | Document Ref | Reference document number |
| ເງື່ອນໄຂການຊຳລະເງິນ | Terms of Payment | Payment terms |
| ສະກຸນເງິນ | Currency | Currency type (default: LAK) |

### Line Items Table

| Column (Lao) | Column (English) | Description |
|--------------|------------------|-------------|
| ລໍາດັບ | Item | Line item number |
| ລະຫັດສິນຄ້າ | Product Code | Product identification code |
| ລາຍການ | Description | Product description |
| ຈໍານວນ | Qty | Quantity ordered |
| ຫົວໜ່ວຍ | Unit | Unit of measure |
| ລາຄາ | Unit Price | Price per unit |
| ສ່ວນຫຼຸດ | Discount | Discount amount |
| ຈໍານວນເງິນ | Amount | Total amount per line |

### Summary Section

| Field (Lao) | Field (English) | Default |
|-------------|-----------------|---------|
| ຈໍານວນເງິນກ່ອນສ່ວນຫຼຸດການຄ້າ | Amount before Trade Discount | 0.00 |
| ສ່ວນຫຼຸດການຄ້າ | Trade Discount | - |
| ລາຄາຫຼັງສ່ວນຫຼຸດ | Amount after Trade Discount | - |
| ອາກອນມູນຄ່າເພີ່ມ | VAT 7% | - |
| ຈໍານວນເງິນທັງໝົດ | Grand Total (VAT included) | 0.00 |

### Remarks Section

| Field (Lao) | Field (English) |
|-------------|-----------------|
| ໝາຍເຫດ | Remark |

### Signature Section

| Role (Lao) | Role (English) | Date Format |
|------------|----------------|-------------|
| ຜູ້ກະກຽມເອກະສານ | Prepared By | ວັນທີ __/__/__ |
| ຜູ້ອະນຸມັດ | Authorized By | ວັນທີ __/__/__ |
| ຜູ້ກວດເອກະສານ | Check By | ວັນທີ __/__/__ |

---

## 2. Purchase Requisition Form

**Company:** Vientiane Cold  
**Document Type:** ໃບຂໍຊື້ / Purchase Requisition  
**Form Code:** FM-PC-04, Rev:00, 19/07/2012

### Header Information

| Field (Lao) | Field (English) | Description |
|-------------|-----------------|-------------|
| ຜູ້ຮ້ອງຂໍ | Requester | Person requesting the purchase |
| ໜ່ວຍງານ | Department | Department making the request |
| ເຫດຜົນທີ່ຮ້ອງຂໍ | Reason | Reason for the purchase request |
| ບໍລິສັດຂາຍເຄື່ອງ | Supplier | Supplier name |

### Document Details

| Field (Lao) | Field (English) |
|-------------|-----------------|
| ເລກທີເອກະສານ | Document No. |
| ວັນທີເອກະສານ | Document Date |
| ວັນທີເຮັດຄໍາຮ້ອງ | Request Date |

### Line Items Table

| Column (Lao) | Column (English) | Description |
|--------------|------------------|-------------|
| ລໍາດັບ | No | Line item number |
| ລະຫັດ | Code | Product code |
| ລາຍການ | Description | Product description |
| ສິນຄ້າຄົງເຫຼືອ | Stock Balance | Current stock balance |
| ຈໍານວນ | QTY | Quantity requested |
| ຫົວໜ່ວຍ | Unit | Unit of measure |
| ລາຄາ/ໜ່ວຍ | Unit Price | Price per unit |
| ຈໍານວນເງິນ | Amount | Total amount |

### Footer Section

| Field (Lao) | Field (English) |
|-------------|-----------------|
| ໝາຍເຫດ | Remark |
| ລວມຈໍານວນເງິນ | Total |

### Signature Section

| Role (Lao) | Role (English) | Date Format |
|------------|----------------|-------------|
| ຜູ້ກວດສອບ | Checked by | ວັນທີ __/__/__ |
| ຜູ້ອະນຸມັດ | Authorized by | ວັນທີ __/__/__ |
| ພະນັກງານຈັດຊື້ | Purchasing | ວັນທີ __/__/__ |

---

## 3. Add New Customer Form

**Form Type:** Software Interface - Customer Registration  
**Title:** ເພີ່ມລູກຄ້າ (Add new Customer)

### Basic Information (Detail 2)

| Field (Lao) | Field (English) | Type | Example |
|-------------|-----------------|------|---------|
| ID | ID | Auto-generated | 762 |
| ຊື່ | Name | Text | - |
| ລະຫັດລູກຄ້າ | Customer Code | Text | - |
| ຊື່ຜູ້ຕິດຕໍ່ໃຊ້ເປັນຄ່າເລີ່ມຕົ້ນ | Default Contact Name | Text | - |
| ທີ່ຢູ່ເຕັມ | Full Address | Text | - |

### Delivery & Contact Information

| Field (Lao) | Field (English) | Type |
|-------------|-----------------|------|
| ຊື່ສະຖານທີ່ຈັດສົ່ງ | Ture Name (Delivery Name) | Text |
| ໂທລະສັບ | Tel | Text |
| ແຟັກ | Fax | Text |
| ມືຖື | Mobile | Text |
| ອາກອນ | Tax | Text |

### Account Information

| Field (Lao) | Field (English) | Type | Default |
|-------------|-----------------|------|---------|
| ວັນເກີດ | Birthday | Date | dd/mm/yyyy |
| ອີເມວເວັບ | Web EMail | Text | - |
| ປະເພດລາຄາ | Price Type | Dropdown | ລາຄາຂາຍ 1 |
| ລະດັບ | Level | Dropdown | ທົ່ວໄປ |
| ສ່ວນຫຼຸດໃບບິນ | Bill Discount | Percentage | 0.00% |
| ຈໍາກັດເງິນເຊື່ອ | Limit Credit Money | Number | - |
| ໝາຍເຫດ | Comment | Text | - |

### Membership & Identity

| Field (Lao) | Field (English) | Type | Default |
|-------------|-----------------|------|---------|
| ໝົດອາຍຸສະມາຊິກ | Member Expire | Date | dd/mm/yyyy |
| ເງິນເຊື່ອ | Credit | Number | 0.00 |
| ເລກໜັງສືຜ່ານແດນ | Passport No. | Text | - |
| ລະຫັດປະເທດ | Country Code | Text | - |
| ເພດ | Gender | Dropdown | - |

### Settings & Preferences

| Field (Lao) | Field (English) | Type |
|-------------|-----------------|------|
| ໃຊ້ໂປຣໂມຊັ່ນ | Use Promotion | Checkbox |
| ໃຊ້ໃບບິນ PMT | Use PMT Bill | Checkbox |
| ຈຸດສະສົມແຕ້ມ | Collect Point | Checkbox |
| ຜູ້ຮັບເງິນ | Payee | Text |
| ສາຂາ RD | RD Branch | Text |
| ບາໂຄດ | Barcode | Text |
| ອີເມວ | EMail | Text |

### Additional Features

- **Smart Card Reader:** Option to read customer card
- **Buttons:** OK (ຕົກລົງ), Close

---

## 4. Stock/Inventory Management Screen

**Form Type:** Software Interface - Inventory Management  
**Title:** Stock Cold room 1 (Welcome Admin)

### Main Menu Tabs

| Tab (Lao) | Tab (English) |
|-----------|---------------|
| ຂໍ້ມູນ | Data |
| ຊື້ | Buy |
| ສາງ | Stock |
| ຂາຍ | Sell |
| ເງິນ | Money |
| ລາຍງານ | Report |
| ໜ້າຕ່າງ | Windows |
| ຊ່ວຍເຫຼືອ | Help |

### Filter/Condition Section (Condition 2)

| Field (Lao) | Field (English) | Type |
|-------------|-----------------|------|
| ໝວດໝູ່ | Category | Dropdown |
| ກຸ່ມສິນຄ້າ | Product Group | Dropdown |
| ຍີ່ຫໍ້ | Brand | Dropdown |
| ຊື່ສາຂາ | Branch Name | Dropdown |

### Options

| Option (Lao) | Option (English) |
|--------------|------------------|
| ຕັ້ງລາຄາຊື້ | Set Buy Price |
| ສະແດງລໍາດັບ | Show Num Order |
| ຂະຫຍາຍກຸ່ມ | Expand Group |

### Product List Table

| Column | Description |
|--------|-------------|
| ID | Product ID number |
| Product Code | Product identification code (e.g., CAN054, CHB830400) |
| Barcode | Product barcode (e.g., FRS192B, 9780201371550) |
| ລາຍການ (Description) | Product name/description |
| ຈໍານວນເຄື່ອງຄົງເຫຼືອ (Stock Balance) | Current stock amount |
| ໜ່ວຍ | Unit (pack, ctn, Kg, carton) |
| ລາຄາ (Price) | Unit price |
| %VAT | VAT percentage |
| Expire | Expiration date |
| ລາຄາ 1 | Price 1 |
| ລາຄາ 2 | Price 2 |
| Lot Number | Lot/Batch number |

### Sample Products Listed

| Product Code | Description | Unit | Price |
|--------------|-------------|------|-------|
| CAN054 | Smoked Mackerel Fillet 150g | pack | 579.00 |
| - | Mackerel in tomato sauce 50x155g | ctn | 1,722.00 |
| CHB830400 | Beef Striploin steak Grain fed F1Wagyu SB 4-5, 200g | pack | 255.00 |
| CHB830403 | Beef Striploin steak grain fed F1Wagyu SB 4-5, 200g | pack | 194.00 |
| CHB830407 | Beef Cube roll steak grain fed Angus 200g | pack | 375.00 |
| CHB830418 | Beef Striploin steak Grain fed Angus 200g | pack | 291.00 |
| CHB830430 | Beef Cube roll steak Grain fed 200g | pack | 270.00 |
| CHB830433 | Beef Cube roll steak Grain fed 200g | pack | 555.00 |
| CHB830436 | Beef Tenderloin steak Grain fed 200g | pack | 60.00 |
| CHB830551 | Beef Chuck eye roll thin slice grain fed 200g | pack | 260.00 |
| CHB834333 | Beef Navel end brisket W/AC Grain fed | Kg | 183.37 |
| CHB882732 | Beef Chuck eye roll W/VAC F1Wagyu SB 4-5 | Kg | 8.40 |
| CHB999751 | Beef "YP-NEB" Navel end brisket W/AC 12/D | pack | 179.81 |
| CHB9997518 | Beef "YP-NEB" Navel end brisket sliced 150g | pack | 219.00 |
| DHGV707 | Ground Black Pepper 24x30g | carton | 73.00 |
| DHGV708 | Ground White Pepper 24x30g | carton | 80.00 |

### Bottom Section

| Feature (Lao) | Feature (English) |
|---------------|-------------------|
| ແກ້ໄຂລາຄາຂາຍ | Edit Sell Price |
| ລາຍລະອຽດເພີ່ມເຕີມ | More Detail |
| ແກ້ໄຂສະຖານທີ່ | Edit Location |
| ເບິ່ງແບບກຸ່ມ | Group View |
| ພິມ | Print |
| ພິມບາໂຄດ | Print Barcode |

### Summary Display

| Field | Value |
|-------|-------|
| ລວມ VAT | (VAT included total) |
| VAT ທັງໝົດ | 543,475,798.66 |
| ວັນທີເລີ່ມສິ້ນສຸດ | Date range display |

---

## 5. Clear Advance Voucher Form

**Company:** Vientiane Cold  
**Document Type:** ໃບຄືນເງິນເບີກລ່ວງໜ້າ / Clear Advance Voucher Form  
**Form Number:** Nº 003328

### Header Information

| Field (Lao) | Field (English) | Example Value |
|-------------|-----------------|---------------|
| ຊື່ຜູ້ເບີກເງິນ | Name | Sounaly |
| ຕຳແໜ່ງ | Position | - |
| ເລກທີ PO | PO / No | - |
| ວັນທີ | Date | 14/01/2026 |

### Currency Selection

| Option | Symbol |
|--------|--------|
| LAK | ☑ (checked) |
| THB | ☐ |
| USD | ☐ |

### Line Items

| ລໍາດັບ (No) | ລາຍການ (Description) | ຈໍານວນ (Amount) |
|-------------|----------------------|-----------------|
| 1 | Advance (ເງິນເບີກລ່ວງໜ້າ) | 2,700,000 |
| 2 | Expend (ລາຍຈ່າຍ) | 2,700,000 |
| - | - ນ້ຳມັນ (Gas/Fuel) | 2,200,000 |
| - | - ອື່ນໆ (Other) | 500,000 |

### Supporting Documents Checklist

| Document (Lao) | Document (English) | Status |
|----------------|-------------------|--------|
| ໃບຮັບເງິນ | Receipt | ☑ |
| ໃບແຈ້ງໜີ້ | Invoice | ☐ |
| ໃບໂອນເງິນ | Bank Transfer slip | ☐ |
| ໃບສະເໜີຊື້, ໃບສັ່ງຊື້ | PO, PR | ☐ |
| ຄ່າຂົນສົ່ງ | Transportation charge | ☐ |

### Signature Section

| Role (Lao) | Role (English) |
|------------|----------------|
| ບັນຊີ | Accountant |
| ຜູ້ຄືນເງິນ | Returned by |

---

## 6. Clear Advance Request Form

**Company:** Vientiane Cold  
**Document Type:** ໃບຂໍເບີກເງິນຄືນລ່ວງໜ້າ / Clear Advance Request Form  
**Form Number:** Nº 003328

### Header Information

| Field (Lao) | Field (English) | Example Value |
|-------------|-----------------|---------------|
| ຊື່ຜູ້ເບີກ | Name | Sounaly |
| ຕໍາແໜ່ງ | Position | - |
| ເລກທີ PO | PO / No | - |
| ວັນທີ | Date | 06/01/2025 |

### Currency Selection

| Option | Symbol |
|--------|--------|
| LAK | ☑ (checked) |
| THB | ☐ |
| USD | ☐ |

### Line Items Table

| ລໍາດັບ (No) | ລາຍການ (Description) | ຈໍານວນ (Amount) |
|-------------|----------------------|-----------------|
| 1 | ຈ່າຍຄ່າ ສ້ອມແປງ ບໍລິການ PBM+ Segun | - |
| A | ສິນຄ້າ (Goods): 2 ລາຍການ x 1,000,000 | 2,000,000 |
| B | ຄ່າຂົນສົ່ງ (Shipping): 1 ລາຍການ x 200,000 | 200,000 |
| - | ສິນຄ້າ (ວັນຈັນ/ວັນສຸກ/ວັນເສົາ/ວັນອາທິດ) | 500,000 |
| - | ລູກຄ້າ 0% + ສົ່ງຟຣີ 1 ລາຍການ | - |
| - | + ຄ່າຂົນສົ່ງ 3 ລາຍການ | 2,700,000 |

### Footer

| Field (Lao) | Field (English) |
|-------------|-----------------|
| ສິ່ງຄາດແນບເງິນ | (Attachments) |

### Signature Section

Three signature boxes for approval chain.

---

## 7. Add New Supplier Form

**Form Type:** Software Interface - Supplier Registration  
**Title:** ເພີ່ມບໍລິສັດຜູ້ສະໜອງ (Add new Supplier)

### Basic Information

| Field (Lao) | Field (English) | Type | Example |
|-------------|-----------------|------|---------|
| Supplier ID | ID ຜູ້ສະໜອງ | Auto-generated | 196 |
| ຊື່ | Name | Text (Required *) | - |
| ລະຫັດຜູ້ສະໜອງ | Supplier Code | Text | - |
| ຊື່ຜູ້ຕິດຕໍ່ | Contact Name | Text | - |
| ສາຂາ RD | RD Branch | Text | - |
| ທີ່ຢູ່ | Address | Text | - |

### Additional Information

| Field (Lao) | Field (English) | Type |
|-------------|-----------------|------|
| ຊື່ຮູບພາບ | Picture Name | Text/Image |
| ໂທລະສັບ | Tel | Text |
| ມືຖື | Mobile | Text |
| ແຟັກ | Fax | Text |
| ລະຫັດອາກອນ | TaxCode | Text |

### Financial Settings

| Field (Lao) | Field (English) | Type | Default |
|-------------|-----------------|------|---------|
| ເງິນເຊື່ອ | Credit | Number | 0 |
| ສ່ວນຫຼຸດ | Discount | Percentage | 0.00% |
| ວົງເງິນເຊື່ອ PO | PO Credit Money | Number | 0.00 |
| ສະກຸນເງິນ | Currency | Dropdown | KIP |

### Action Buttons

| Button (Lao) | Button (English) |
|--------------|------------------|
| ຕົກລົງ | OK |
| ຍົກເລີກ | Cancel |

---

## Common Field Types Reference

| Type | Description |
|------|-------------|
| Text | Free-form text input |
| Number | Numeric input |
| Date | Date picker (format: dd/mm/yyyy) |
| Dropdown | Selection from predefined options |
| Checkbox | Boolean selection |
| Percentage | Numeric with % symbol |
| Auto-generated | System-generated value |

---

## Currency Codes

| Code | Currency |
|------|----------|
| LAK | Lao Kip (ກີບ) |
| THB | Thai Baht |
| USD | US Dollar |
| KIP | Kip (alternative code) |

---

## Document Workflow

1. **Purchase Requisition** - Internal request for purchase
2. **Purchase Order** - Official order to supplier
3. **Clear Advance Request** - Request for advance payment
4. **Clear Advance Voucher** - Documentation of advance usage

---

*Document generated from company forms and system screenshots*  
*Last updated: January 2026*
