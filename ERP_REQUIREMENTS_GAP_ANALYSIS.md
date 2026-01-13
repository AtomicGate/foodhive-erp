## ERP Requirements – Coverage & Gap Analysis

This document maps the features listed in `CamScanner 01-09-2026 18.34.pdf` to the current FoodHive ERP implementation (backend + frontend) and identifies what is **implemented**, **partially implemented**, and **missing / TODO**.

Legend:
- **✅ Implemented**: Core logic and API exist; basic UI exists or is straightforward to add.
- **🟡 Partial**: Some logic/APIs exist, but not all sub‑features or screens are done.
- **❌ Missing**: No meaningful implementation yet.

---

## 1. Sales / Order Entry

**Status: 🟡 Partial**

We have a Sales Orders module (APIs + frontend pages) with basic CRUD, listing, and integration to inventory and AR/GL hooks, but several convenience and analytical features are not yet implemented.

**Document features vs status**

- **Rapid entry mode / order-guide / traditional**  
  - ✅ Standard order entry (sales order creation APIs + UI).  
  - 🟡 Order-guide style entry: backend has `customer` order‑guide endpoint and model; frontend still basic and not integrated into the order entry screen.  
  - ❌ Dedicated “rapid entry” UX (keyboard-driven, minimal clicks) is not yet designed.

- **Multiple order types (standard, advance, pre-paid, on hold, quotes, credit memo, pickup memo)**  
  - 🟡 Order type field exists conceptually in sales order models, but only a subset is wired end‑to‑end; business rules per type (e.g., prepaid vs quote vs credit memo) are not fully enforced in services or UI.

- **Multiple warehouse sales capability**  
  - 🟡 Warehouse dimension exists in schema and services; some endpoints accept/return warehouse IDs.  
  - ❌ End‑to‑end multi‑warehouse flow from pricing → order → picking → invoicing is not fully validated.

- **Auto-routing based on customer's next available route**  
  - ❌ Route master data and auto-assignment logic not implemented.

- **Customer order guide (previous products, stats, default prices, push items, etc.)**  
  - ✅ Backend has `customer` order‑guide endpoint and model using sales history.  
  - ❌ Frontend integration into the Sales Order screen (order guide mode, quick add from guide).

- **Inventory visibility on screen (on hand, allocated, on order)**  
  - ✅ Inventory service exposes on hand / allocated concepts.  
  - 🟡 Sales order UI shows basic availability but not full three-way (on hand / allocated / on order) view per line as described.

- **Selling below cost notification & margin management**  
  - ❌ No explicit below-cost check in pricing/sales services.  
  - ❌ Margin management UI & approval workflows not implemented.

- **Auto price from customer setup & integrated pricing**  
  - ✅ Pricing module exists with contracts / price levels, and APIs for price determination.  
  - 🟡 Not all price hierarchy levels (contract, price level, custom, promo, order price) are fully enforced at order-entry time; frontend uses simplified logic.

- **Catch weight integration with picking**  
  - ✅ Catch weight core models & APIs exist.  
  - 🟡 Integration of actual picked weights into sales invoices and margin reporting is not fully wired.

- **Lot and best-before/expiry lookup**  
  - 🟡 Lot & expiry metadata supported in inventory schema.  
  - ❌ Dedicated lot/expiry lookup UI from order entry is missing.

- **Lost sales reporting due to lack of inventory**  
  - ❌ No explicit lost-sales logging/report/reporting implemented.

- **Special order request to buyer from order screen**  
  - ❌ No automated special‑order request generation from Sales → Purchasing.

- **View customer account and drill down past invoices from order screen**  
  - 🟡 AR/customer statement APIs exist; partial reporting.  
  - ❌ Deep drill‑down from the Sales Order screen (UI workflow) not implemented.

- **EDI / e‑commerce / remote sales integration**  
  - ❌ No EDI or web‑order integration layer implemented.

- **Commission reporting**  
  - 🟡 Sales rep fields exist on customers/orders; commission logic and reports are not complete.

**Key TODOs for Sales / Order Entry**

- Implement **order-type behaviors** (prepaid, quotes, credit memo, on‑hold) with validations & lifecycle rules.
- Finish **order-guide mode** in frontend and wire to `/customers/{id}/order-guide`.
- Add **below-cost checks** and **margin calculation** per line and per order (with configurable thresholds).
- Implement **lost-sales tracking** and reports.
- Add **special order request** flow that creates purchase suggestions for buyers.
- Create **customer AR drill‑down** panel on the order screen.

---

## 2. Picking & Routing

**Status: 🟡 Partial**

Pick list and routing concepts exist, integrated with inventory and catch weight, but the WMS-level flow is not fully implemented.

**Document features vs status**

- **Route setup & maintenance** – ❌ No dedicated route master with calendars and zones.  
- **Zone picking / picking by date, route or batch** – 🟡 Basic picking by order/route exists; zone/batch logic is minimal.  
- **Master pick reporting** – ❌ Not yet implemented as dedicated reports.  
- **Mobile / labels output** – 🟡 Label printing and barcode support exist at model level; no actual mobile integration.  
- **Re‑number stops, shuffle orders between routes** – ❌ Not implemented.  
- **Batch release & batch printing** – ❌ No batch release UI/workflow.  
- **Integrated with lot tracking, catch weight, bar code capture** – 🟡 Integration points exist; needs end‑to‑end testing and UI.  
- **Suggested picking based on rotation (FIFO/FEFO, etc.)** – ❌ Not implemented.  
- **Integration with re‑packing/processing** – ❌ Not implemented.  
- **Shipping labels, pallet functions, UCC/EAN** – 🟡 Barcodes supported conceptually; pallet & shipping label workflows missing.

**Key TODOs for Picking & Routing**

- Design & implement **route master** (routes, stops, calendars).
- Implement **picking strategies** (by route/date/batch, zone picking).
- Add **pick waves / master pick reports** and **shipping label** generation.
- Integrate **catch-weight capture** directly into the pick confirmation flow.

---

## 3. Customer & AR Management

**Status: 🟡 Partial**

Core customer master, AR aging, and statements are present; some advanced profiling and multi‑currency aspects are not.

**Document features vs status**

- Customer data access from order screen – 🟡 Data available via APIs; UX still basic.  
- Multiple currency AR – ❌ System largely assumes single base currency.  
- Multiple ship‑to linked to bill‑to – ✅ Customer ship‑to endpoints exist; frontend list/maintenance exists.  
- Flexible statements – 🟡 Statements and aging exist, but configuration options are limited.  
- Default warehouse & shipping method per customer – 🟡 Fields exist or can be added; logic not fully applied.  
- Customer product code links – ❌ Not implemented.  
- Rep integration & commission – 🟡 Rep fields present; commission engine and reports incomplete.  
- Grouping by bill‑to/ship‑to – 🟡 Some grouping in reporting; not fully generalized.  
- Order guide based on history – ✅ Implemented in backend; not fully used in UI.  
- Credit management profiling – 🟡 Credit limits & basic checks exist; detailed profiling rules & scoring not implemented.  
- Extensive AR reports – 🟡 Core reports exist; additional formats and filters are still to be added.

**Key TODOs for Customer & AR**

- Add **multi-currency support** at AR transaction level.  
- Implement **customer product codes** (customer‑specific SKU mapping).  
- Build **credit profile dashboard** (limits, history, risk flags).  
- Flesh out **AR reporting screens** (aging buckets, statements, export).

---

## 4. Pricing & Cost Management

**Status: 🟡 Partial**

Pricing module exists with contracts, levels, effective dates; backflush cost correction is not done.

**Document features vs status**

- 5‑level price default hierarchy – 🟡 Model supports contracts, lists, promotions; enforcement still simplified.  
- Price & cost by future effective dates – ✅ Effective‑date fields exist in pricing models.  
- Print price lists by effective date – ❌ No dedicated report/export yet.  
- Mass price maintenance – ❌ No bulk price‑update tools in UI.  
- Overwrite price with security levels – 🟡 Manual override possible; security controls granularly not finished.  
- Rebate tracking integration – ❌ Rebate engine and accruals missing.  
- Multiple costing methods (avg, last, landed, etc.) – 🟡 Core cost fields exist; per-product method selection and reporting not complete.  
- Automatic cost error backflush in GL – ❌ Not implemented.

**Key TODOs for Pricing**

- Implement **full price hierarchy resolver** used consistently in Sales & AR.  
- Add **mass maintenance** and **price list** reporting screens.  
- Design and implement **rebate tracking** (from invoice lines to GL).  
- Add **cost method configuration** per product and proper cost rollups.

---

## 5. Inventory Control

**Status: 🟡 Partial**

Inventory tables, units of measure, catch weight, and lot/expiry support exist. Some advanced traceability and reservation features are missing.

**Document features vs status**

- Product group codes / categories / posting – ✅ Product master supports category & posting accounts.  
- Barcodes (EAN/UCC) – 🟡 Barcode fields available; generation/printing not fully done.  
- Catch weight – ✅ Catch weight module implemented; integrated partially with inventory, sales, purchasing.  
- Multiple units of measure – ✅ Units table and conversions exist.  
- Broken case / portion costing – 🟡 Supported conceptually via UOM & pricing; not fully enforced in workflows.  
- COOL (country of origin labeling) – ❌ Not implemented.  
- Real‑time inventory inquiry – ✅ Inventory list and detail endpoints; UI exists.  
- Grade/HACCP/QA comments – ❌ Not fully represented in schema.  
- Lb/Kg conversion – ✅ Conversion logic supported via UOM factors.  
- Age of inventory reporting – 🟡 Partially supported; dedicated reports missing.  
- Multi‑warehouse tracking – ✅ Warehouses and stock by warehouse exist.  
- Lot/date tracking; one‑up/one‑down traceability – 🟡 Lot & expiry exist; full trace (source → destination) requires more work.  
- Product reservations by customer/group/rep – ❌ Not implemented.

**Key TODOs for Inventory**

- Implement **country of origin, grade, HACCP** attributes on products/batches.  
- Build **age of inventory** and **traceability** reports.  
- Add **reservation logic** for customers / reps and integrate into order allocation.

---

## 6. Purchasing & Receiving

**Status: 🟡 Partial**

Purchase Orders and basic receiving are present; advanced buyer tools are limited.

**Document features vs status**

- Routine buying by buyer/product/supplier – 🟡 Basic PO creation and vendor link exist; buyer assignment is minimal.  
- Reorder by min/max or average sales – ❌ Automatic suggestion engine not implemented.  
- Quick PO entry – ✅ APIs + simple UI exist.  
- Multi‑warehouse purchasing – 🟡 Warehouse field exists; business rules not exhaustive.  
- ETA scheduling (lead times, expected dates) – 🟡 Fields exist; scheduling logic and alerts minimal.  
- Inventory & supplier order guides from PO – 🟡 Inventory lookup exists; supplier order guide not done.  
- LB/KG conversion for supplier units – ✅ Supported through UOM; needs UI polish.  
- Short shipment reporting – ❌ Not implemented as a dedicated report/workflow.  
- Buyer’s reporting & tools (auto‑suggested buying) – ❌ Not implemented.

**Key TODOs for Purchasing**

- Build **replenishment engine** (min/max & forecast/avg‑sales based).  
- Implement **short‑shipment logging** and reports.  
- Create **buyer workbench** screen with suggested POs and exceptions.

---

## 7. Accounts Payable (AP)

**Status: 🟡 Partial**

AP module and integration with inventory/GL exist; scheduling and discount logic are basic.

**Document features vs status**

- Vendor account history (screen & reports) – 🟡 Basic history screens exist; reports to expand.  
- Vendor discount tables for payables – ❌ Not implemented.  
- Automatic check runs by date & vendor type – ❌ Payment proposal engine missing.  
- Inventory adjustments from AP invoices – 🟡 Partial integration: cost/qty adjustments impact inventory, but edge cases need testing.  
- AP → PO → vendor settings integration – 🟡 Partially wired.  
- Vendor templates for other‑in costing – ❌ Not implemented.  
- Vendor lists & labels – ❌ Reporting not yet created.  
- Aged payables “as of any date” – 🟡 Core aging logic exists in GL/AP; reporting UI limited.

**Key TODOs for AP**

- Implement **payment proposal** and batch payment runs.  
- Add **discount calendar/tables** and show “pay by” suggestions.  
- Build **AP aging** and vendor reports as per spec.

---

## 8. Bank & Reconciliation

**Status: ❌ Missing (Core design only in GL)**

No dedicated bank module yet; only GL cash accounts.

**Document features vs status**

- Manage checks & receipts – ❌  
- Automatic payable runs by date – ❌ (ties into AP payment proposal).  
- Manual checks, reprint, void – ❌  
- View open invoices by customer and assign payments – 🟡 Partially via AR; dedicated cash‑application UI missing.  
- Receipt registers & reports – ❌  
- Bank reconciliation by account/date – ❌

**Key TODOs for Bank & Reconciliation**

- Design & implement **Bank module** (bank accounts, statements, reconciliations).  
- Implement **cash receipt** and **check issuance** flows integrated with AR/AP/GL.

---

## 9. General Ledger (GL)

**Status: 🟡 Partial**

We have:
- Chart of accounts  
- Journal entries with posting and reversal  
- Basic Trial Balance / Income Statement / Balance Sheet  
But some advanced features are not finalized.

**Document features vs status**

- Multi‑currency at transaction level – ❌ Not fully implemented; base‑currency only.  
- Multi‑branch/division/department reporting – 🟡 Dimensions exist (department, warehouse); full multi‑segment reporting is limited.  
- Integrated financial reporting – ✅ Core financial reports exist.  
- Comparisons vs budget / prior period / year – ❌ Budgeting and comparative reporting missing.  
- Recurring journals – 🟡 Simple journal templates partially exist; recurring scheduling not done.  
- Auto‑reversal – ✅ Supported for selected journal types.  
- Full audit trail – ✅ Journal and document IDs with references are logged; UI for viewing audit logs could improve.  
- Drill‑down – 🟡 Some drill‑down from reports to entries exists; end‑to‑end (from GL to source doc) not complete.  
- “As at any date” snapshots – 🟡 Trial balance supports dates; more snapshot tools needed.

**Key TODOs for GL**

- Implement **multi‑currency** for AR/AP/Inventory and GL postings.  
- Add **budget module** and comparative reporting.  
- Improve **drill‑down** from financial statements to underlying documents.

---

## 10. Warehouse Management System (WMS)

**Status: ❌ Mostly Missing (only basic inventory & picking parts)**

The WMS add‑on described (locations, slot/section/skid, re‑work, disposals, cycle counts, etc.) is not fully implemented.

**Document features vs status**

- Receiving, Put‑away, Replenishment – 🟡 Receiving exists; directed put‑away/replenishment not.  
- Skid tracking, location (section/slot/skid) – ❌ Not modeled.  
- Warehouse transfers – 🟡 Simple transfer concepts present; dedicated flows limited.  
- Cycle counts, physical counts – ❌ Counting module not implemented.  
- Labeling, skid maintenance, disposals – ❌ Not implemented.  
- Re‑work – ❌ Not implemented.

**Key TODOs for WMS**

- Design WMS data model (locations, bins, skids) and flows.  
- Implement **cycle count**, **transfer**, and **re‑work** processes and screens.

---

## 11. Payroll

**Status: ❌ Not Implemented**

The document only mentions integration to GL, Bank, AR. No payroll module exists in this codebase.

**Key TODOs for Payroll**

- Decide scope: integrate with external payroll vs create basic internal module.  
- If internal: design employees’ earnings, deductions, runs, and GL postings.

---

## 12. High‑Priority Next Steps

Based on the document and current code, these are the **most impactful next implementations**:

1. **Stabilize Core Master Data & Security**
   - Finalize Employees, Departments, Roles, Permissions (emp_page) with the new clean schema.
   - Finish frontend for managing roles, pages, and permissions.

2. **Complete Sales / Order Entry Experience**
   - Order types, order‑guide mode, margin/below‑cost checks.
   - Customer AR drill‑down from order screen.

3. **Tighten Inventory + Catch Weight + Picking**
   - Make catch‑weight fully drive invoiced quantities and margins.  
   - Improve picking flow, including lot/expiry selection.

4. **GL & Financial Reporting**
   - Harden Trial Balance / IS / BS and add basic “as of date” snapshots.  
   - Start budget/comparison design (even simple version).

5. **AP / AR Aging & Payment Flows**
   - Aging reports “as of date”.  
   - Simple payment proposal for AP and receipt application for AR.

If you’d like, I can now convert this analysis into a concrete **implementation roadmap** (with phases and tasks) or start by implementing one of the missing pieces (for example, complete order‑guide mode or AP aging reports). 

