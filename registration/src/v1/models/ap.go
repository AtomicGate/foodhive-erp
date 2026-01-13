package models

// ============================================
// AP Enums
// ============================================

type APInvoiceStatus string

const (
	APInvoiceStatusPending  APInvoiceStatus = "PENDING"
	APInvoiceStatusApproved APInvoiceStatus = "APPROVED"
	APInvoiceStatusPartial  APInvoiceStatus = "PARTIAL"
	APInvoiceStatusPaid     APInvoiceStatus = "PAID"
	APInvoiceStatusVoid     APInvoiceStatus = "VOID"
)

type APPaymentMethod string

const (
	APPaymentCheck    APPaymentMethod = "CHECK"
	APPaymentTransfer APPaymentMethod = "TRANSFER"
	APPaymentACH      APPaymentMethod = "ACH"
	APPaymentCash     APPaymentMethod = "CASH"
)

// ============================================
// AP Invoice Models
// ============================================

type APInvoice struct {
	ID             int             `json:"id"`
	InvoiceNumber  string          `json:"invoice_number"`
	VendorID       int             `json:"vendor_id"`
	POID           *int            `json:"po_id,omitempty"`
	ReceivingID    *int            `json:"receiving_id,omitempty"`
	InvoiceDate    CustomDate      `json:"invoice_date"`
	DueDate        CustomDate      `json:"due_date"`
	Status         APInvoiceStatus `json:"status"`
	Subtotal       float64         `json:"subtotal"`
	TaxAmount      float64         `json:"tax_amount"`
	FreightAmount  float64         `json:"freight_amount"`
	DiscountAmount float64         `json:"discount_amount"`
	TotalAmount    float64         `json:"total_amount"`
	AmountPaid     float64         `json:"amount_paid"`
	BalanceDue     float64         `json:"balance_due"`
	Currency       string          `json:"currency"`
	Notes          string          `json:"notes,omitempty"`
	ApprovedBy     *int            `json:"approved_by,omitempty"`
	ApprovedAt     CustomDateTime  `json:"approved_at,omitempty"`
	CreatedBy      int             `json:"created_by"`
	CreatedAt      CustomDateTime  `json:"created_at"`
	UpdatedAt      CustomDateTime  `json:"updated_at"`
}

type APInvoiceLine struct {
	ID              int     `json:"id"`
	InvoiceID       int     `json:"invoice_id"`
	LineNumber      int     `json:"line_number"`
	ProductID       *int    `json:"product_id,omitempty"`
	Description     string  `json:"description"`
	Quantity        float64 `json:"quantity"`
	UnitCost        float64 `json:"unit_cost"`
	TaxPercent      float64 `json:"tax_percent"`
	LineTotal       float64 `json:"line_total"`
	POLineID        *int    `json:"po_line_id,omitempty"`
	ReceivingLineID *int    `json:"receiving_line_id,omitempty"`
	GLAccountID     *int    `json:"gl_account_id,omitempty"`
}

type APInvoiceWithDetails struct {
	Invoice     APInvoice       `json:"invoice"`
	Lines       []APInvoiceLine `json:"lines"`
	VendorName  string          `json:"vendor_name"`
	VendorCode  string          `json:"vendor_code"`
	PONumber    string          `json:"po_number,omitempty"`
	DaysOverdue int             `json:"days_overdue,omitempty"`
}

// ============================================
// AP Payment (Disbursement) Models
// ============================================

type APPayment struct {
	ID            int             `json:"id"`
	PaymentNumber string          `json:"payment_number"`
	VendorID      int             `json:"vendor_id"`
	PaymentDate   CustomDate      `json:"payment_date"`
	PaymentMethod APPaymentMethod `json:"payment_method"`
	Amount        float64         `json:"amount"`
	Currency      string          `json:"currency"`
	CheckNumber   string          `json:"check_number,omitempty"`
	BankAccountID *int            `json:"bank_account_id,omitempty"`
	ReferenceNo   string          `json:"reference_no,omitempty"`
	Notes         string          `json:"notes,omitempty"`
	PreparedBy    int             `json:"prepared_by"`
	ApprovedBy    *int            `json:"approved_by,omitempty"`
	IsVoided      bool            `json:"is_voided"`
	CreatedAt     CustomDateTime  `json:"created_at"`
}

type APPaymentApplication struct {
	ID        int     `json:"id"`
	PaymentID int     `json:"payment_id"`
	InvoiceID int     `json:"invoice_id"`
	Amount    float64 `json:"amount"`
}

type APPaymentWithDetails struct {
	Payment        APPayment              `json:"payment"`
	Applications   []APPaymentApplication `json:"applications"`
	VendorName     string                 `json:"vendor_name"`
	PreparedByName string                 `json:"prepared_by_name"`
}

// ============================================
// Vendor Balance & Aging
// ============================================

type VendorBalance struct {
	VendorID       int     `json:"vendor_id"`
	VendorName     string  `json:"vendor_name"`
	CurrentBalance float64 `json:"current_balance"`
	TotalOverdue   float64 `json:"total_overdue"`
	OldestOverdue  int     `json:"oldest_overdue_days"`
	PaymentTerms   int     `json:"payment_terms_days"`
}

type APAgingBucket struct {
	Current   float64 `json:"current"`
	Days1_30  float64 `json:"days_1_30"`
	Days31_60 float64 `json:"days_31_60"`
	Days61_90 float64 `json:"days_61_90"`
	Over90    float64 `json:"over_90"`
	Total     float64 `json:"total"`
}

type VendorAging struct {
	VendorID   int           `json:"vendor_id"`
	VendorName string        `json:"vendor_name"`
	VendorCode string        `json:"vendor_code"`
	Aging      APAgingBucket `json:"aging"`
}

// ============================================
// Request/Response Types
// ============================================

type CreateAPInvoiceRequest struct {
	VendorID      int                      `json:"vendor_id"`
	InvoiceNumber string                   `json:"invoice_number"`
	POID          *int                     `json:"po_id,omitempty"`
	ReceivingID   *int                     `json:"receiving_id,omitempty"`
	InvoiceDate   string                   `json:"invoice_date"`
	DueDate       string                   `json:"due_date,omitempty"`
	TaxAmount     float64                  `json:"tax_amount,omitempty"`
	FreightAmount float64                  `json:"freight_amount,omitempty"`
	Notes         string                   `json:"notes,omitempty"`
	Lines         []CreateAPInvoiceLineReq `json:"lines"`
}

type CreateAPInvoiceLineReq struct {
	ProductID   *int    `json:"product_id,omitempty"`
	Description string  `json:"description"`
	Quantity    float64 `json:"quantity"`
	UnitCost    float64 `json:"unit_cost"`
	TaxPercent  float64 `json:"tax_percent,omitempty"`
	GLAccountID *int    `json:"gl_account_id,omitempty"`
}

type CreateAPPaymentRequest struct {
	VendorID      int                       `json:"vendor_id"`
	PaymentDate   string                    `json:"payment_date"`
	PaymentMethod APPaymentMethod           `json:"payment_method"`
	Amount        float64                   `json:"amount"`
	CheckNumber   string                    `json:"check_number,omitempty"`
	BankAccountID *int                      `json:"bank_account_id,omitempty"`
	ReferenceNo   string                    `json:"reference_no,omitempty"`
	Notes         string                    `json:"notes,omitempty"`
	Applications  []APPaymentApplicationReq `json:"applications"`
}

type APPaymentApplicationReq struct {
	InvoiceID int     `json:"invoice_id"`
	Amount    float64 `json:"amount"`
}

type APInvoiceListFilters struct {
	VendorID *int             `json:"vendor_id,omitempty"`
	Status   *APInvoiceStatus `json:"status,omitempty"`
	DateFrom string           `json:"date_from,omitempty"`
	DateTo   string           `json:"date_to,omitempty"`
	Overdue  bool             `json:"overdue,omitempty"`
	Page     int              `json:"page"`
	PageSize int              `json:"page_size"`
}

// ============================================
// Validation
// ============================================

func ValidateAPInvoice(v *Validator, req *CreateAPInvoiceRequest) {
	v.Check(req.VendorID > 0, "vendor_id", "Vendor is required")
	v.Check(req.InvoiceNumber != "", "invoice_number", "Invoice number is required")
	v.Check(req.InvoiceDate != "", "invoice_date", "Invoice date is required")
	v.Check(len(req.Lines) > 0, "lines", "At least one line is required")

	for _, line := range req.Lines {
		v.Check(line.Description != "", "lines", "Description is required for all lines")
		v.Check(line.Quantity > 0, "lines", "Quantity must be positive")
		v.Check(line.UnitCost >= 0, "lines", "Unit cost must be non-negative")
	}
}

func ValidateAPPayment(v *Validator, req *CreateAPPaymentRequest) {
	v.Check(req.VendorID > 0, "vendor_id", "Vendor is required")
	v.Check(req.PaymentDate != "", "payment_date", "Payment date is required")
	v.Check(req.PaymentMethod != "", "payment_method", "Payment method is required")
	v.Check(req.Amount > 0, "amount", "Amount must be positive")
}
