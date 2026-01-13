package models

// ============================================
// AR (Accounts Receivable) Enums
// ============================================

type ARInvoiceStatus string

const (
	ARInvoiceStatusDraft   ARInvoiceStatus = "DRAFT"
	ARInvoiceStatusPosted  ARInvoiceStatus = "POSTED"
	ARInvoiceStatusPartial ARInvoiceStatus = "PARTIAL"
	ARInvoiceStatusPaid    ARInvoiceStatus = "PAID"
	ARInvoiceStatusOverdue ARInvoiceStatus = "OVERDUE"
	ARInvoiceStatusVoid    ARInvoiceStatus = "VOID"
)

type PaymentMethod string

const (
	PaymentMethodCash     PaymentMethod = "CASH"
	PaymentMethodCheck    PaymentMethod = "CHECK"
	PaymentMethodCard     PaymentMethod = "CARD"
	PaymentMethodTransfer PaymentMethod = "TRANSFER"
	PaymentMethodCredit   PaymentMethod = "CREDIT"
)

// ============================================
// AR Invoice Models
// ============================================

type ARInvoice struct {
	ID             int             `json:"id"`
	InvoiceNumber  string          `json:"invoice_number"`
	CustomerID     int             `json:"customer_id"`
	OrderID        *int            `json:"order_id,omitempty"`
	InvoiceDate    CustomDate      `json:"invoice_date"`
	DueDate        CustomDate      `json:"due_date"`
	Status         ARInvoiceStatus `json:"status"`
	Subtotal       float64         `json:"subtotal"`
	TaxAmount      float64         `json:"tax_amount"`
	FreightAmount  float64         `json:"freight_amount"`
	DiscountAmount float64         `json:"discount_amount"`
	TotalAmount    float64         `json:"total_amount"`
	AmountPaid     float64         `json:"amount_paid"`
	BalanceDue     float64         `json:"balance_due"`
	Currency       string          `json:"currency"`
	Notes          string          `json:"notes,omitempty"`
	PostedBy       *int            `json:"posted_by,omitempty"`
	PostedAt       CustomDateTime  `json:"posted_at,omitempty"`
	CreatedBy      int             `json:"created_by"`
	CreatedAt      CustomDateTime  `json:"created_at"`
	UpdatedAt      CustomDateTime  `json:"updated_at"`
}

type ARInvoiceLine struct {
	ID          int     `json:"id"`
	InvoiceID   int     `json:"invoice_id"`
	LineNumber  int     `json:"line_number"`
	ProductID   *int    `json:"product_id,omitempty"`
	Description string  `json:"description"`
	Quantity    float64 `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	TaxPercent  float64 `json:"tax_percent"`
	LineTotal   float64 `json:"line_total"`
	OrderLineID *int    `json:"order_line_id,omitempty"`
}

type ARInvoiceWithDetails struct {
	Invoice      ARInvoice       `json:"invoice"`
	Lines        []ARInvoiceLine `json:"lines"`
	CustomerName string          `json:"customer_name"`
	CustomerCode string          `json:"customer_code"`
	OrderNumber  string          `json:"order_number,omitempty"`
	DaysOverdue  int             `json:"days_overdue,omitempty"`
}

// ============================================
// AR Payment (Receipt) Models
// ============================================

type ARPayment struct {
	ID            int            `json:"id"`
	ReceiptNumber string         `json:"receipt_number"`
	CustomerID    int            `json:"customer_id"`
	PaymentDate   CustomDate     `json:"payment_date"`
	PaymentMethod PaymentMethod  `json:"payment_method"`
	Amount        float64        `json:"amount"`
	Currency      string         `json:"currency"`
	ReferenceNo   string         `json:"reference_no,omitempty"`
	CheckNumber   string         `json:"check_number,omitempty"`
	BankAccount   string         `json:"bank_account,omitempty"`
	Notes         string         `json:"notes,omitempty"`
	ReceivedBy    int            `json:"received_by"`
	PostedAt      CustomDateTime `json:"posted_at,omitempty"`
	CreatedAt     CustomDateTime `json:"created_at"`
}

type ARPaymentApplication struct {
	ID        int     `json:"id"`
	PaymentID int     `json:"payment_id"`
	InvoiceID int     `json:"invoice_id"`
	Amount    float64 `json:"amount"`
}

type ARPaymentWithDetails struct {
	Payment      ARPayment              `json:"payment"`
	Applications []ARPaymentApplication `json:"applications"`
	CustomerName string                 `json:"customer_name"`
	ReceiverName string                 `json:"receiver_name"`
}

// ============================================
// Credit Management
// ============================================

type CustomerCredit struct {
	CustomerID       int     `json:"customer_id"`
	CustomerName     string  `json:"customer_name"`
	CreditLimit      float64 `json:"credit_limit"`
	CurrentBalance   float64 `json:"current_balance"`
	AvailableCredit  float64 `json:"available_credit"`
	TotalOverdue     float64 `json:"total_overdue"`
	OldestOverdue    int     `json:"oldest_overdue_days"`
	CreditStatus     string  `json:"credit_status"`
	PaymentTermsDays int     `json:"payment_terms_days"`
}

type AgingBucket struct {
	Current   float64 `json:"current"`
	Days1_30  float64 `json:"days_1_30"`
	Days31_60 float64 `json:"days_31_60"`
	Days61_90 float64 `json:"days_61_90"`
	Over90    float64 `json:"over_90"`
	Total     float64 `json:"total"`
}

type CustomerAging struct {
	CustomerID   int         `json:"customer_id"`
	CustomerName string      `json:"customer_name"`
	CustomerCode string      `json:"customer_code"`
	Aging        AgingBucket `json:"aging"`
}

// ============================================
// Statement
// ============================================

type StatementLine struct {
	Date        CustomDate `json:"date"`
	Type        string     `json:"type"`
	Reference   string     `json:"reference"`
	Description string     `json:"description"`
	Debit       float64    `json:"debit"`
	Credit      float64    `json:"credit"`
	Balance     float64    `json:"balance"`
}

type CustomerStatement struct {
	CustomerID     int             `json:"customer_id"`
	CustomerName   string          `json:"customer_name"`
	CustomerCode   string          `json:"customer_code"`
	StatementDate  CustomDate      `json:"statement_date"`
	OpeningBalance float64         `json:"opening_balance"`
	ClosingBalance float64         `json:"closing_balance"`
	Lines          []StatementLine `json:"lines"`
}

// ============================================
// Request/Response Types
// ============================================

type CreateARInvoiceRequest struct {
	CustomerID    int                      `json:"customer_id"`
	OrderID       *int                     `json:"order_id,omitempty"`
	InvoiceDate   string                   `json:"invoice_date"`
	DueDate       string                   `json:"due_date,omitempty"`
	TaxAmount     float64                  `json:"tax_amount,omitempty"`
	FreightAmount float64                  `json:"freight_amount,omitempty"`
	Notes         string                   `json:"notes,omitempty"`
	Lines         []CreateARInvoiceLineReq `json:"lines"`
}

type CreateARInvoiceLineReq struct {
	ProductID   *int    `json:"product_id,omitempty"`
	Description string  `json:"description"`
	Quantity    float64 `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	TaxPercent  float64 `json:"tax_percent,omitempty"`
}

type CreateARPaymentRequest struct {
	CustomerID    int                     `json:"customer_id"`
	PaymentDate   string                  `json:"payment_date"`
	PaymentMethod PaymentMethod           `json:"payment_method"`
	Amount        float64                 `json:"amount"`
	ReferenceNo   string                  `json:"reference_no,omitempty"`
	CheckNumber   string                  `json:"check_number,omitempty"`
	Notes         string                  `json:"notes,omitempty"`
	Applications  []PaymentApplicationReq `json:"applications"`
}

type PaymentApplicationReq struct {
	InvoiceID int     `json:"invoice_id"`
	Amount    float64 `json:"amount"`
}

type ARInvoiceListFilters struct {
	CustomerID *int             `json:"customer_id,omitempty"`
	Status     *ARInvoiceStatus `json:"status,omitempty"`
	DateFrom   string           `json:"date_from,omitempty"`
	DateTo     string           `json:"date_to,omitempty"`
	Overdue    bool             `json:"overdue,omitempty"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
}

// ============================================
// Validation
// ============================================

func ValidateARInvoice(v *Validator, req *CreateARInvoiceRequest) {
	v.Check(req.CustomerID > 0, "customer_id", "Customer is required")
	v.Check(req.InvoiceDate != "", "invoice_date", "Invoice date is required")
	v.Check(len(req.Lines) > 0, "lines", "At least one line is required")

	for _, line := range req.Lines {
		v.Check(line.Description != "", "lines", "Description is required for all lines")
		v.Check(line.Quantity > 0, "lines", "Quantity must be positive")
	}
}

func ValidateARPayment(v *Validator, req *CreateARPaymentRequest) {
	v.Check(req.CustomerID > 0, "customer_id", "Customer is required")
	v.Check(req.PaymentDate != "", "payment_date", "Payment date is required")
	v.Check(req.PaymentMethod != "", "payment_method", "Payment method is required")
	v.Check(req.Amount > 0, "amount", "Amount must be positive")
}
