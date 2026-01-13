package models

// ============================================
// Income/Expense Types (For GL categorization)
// ============================================

type TransactionCategory string

const (
	TransactionCategoryInternal TransactionCategory = "INTERNAL"
	TransactionCategoryExternal TransactionCategory = "EXTERNAL"
)

// ============================================
// Cash Box Models (Treasury Management)
// ============================================

type CashBox struct {
	ID             int            `json:"id"`
	Name           string         `json:"name"`
	Code           string         `json:"code,omitempty"`
	Currency       string         `json:"currency"`
	CurrentBalance float64        `json:"current_balance"`
	IsActive       bool           `json:"is_active"`
	WarehouseID    *int           `json:"warehouse_id,omitempty"`
	CreatedAt      CustomDateTime `json:"created_at"`
	UpdatedAt      CustomDateTime `json:"updated_at"`
}

// ============================================
// Income Models (AR-related)
// ============================================

type IncomeType struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
	GLAccountID *int   `json:"gl_account_id,omitempty"`
	IsActive    bool   `json:"is_active"`
}

type Income struct {
	ID          int                 `json:"id"`
	IncomeType  TransactionCategory `json:"income_type"` // INTERNAL or EXTERNAL
	TypeID      int                 `json:"type_id"`     // References IncomeType
	Date        CustomDate          `json:"date"`
	ReceiptDate CustomDate          `json:"receipt_date,omitempty"`
	Amount      float64             `json:"amount"`
	Note        string              `json:"note,omitempty"`
	CashBoxID   int                 `json:"cash_box_id"`
	CustomerID  *int                `json:"customer_id,omitempty"`
	InvoiceID   *int                `json:"invoice_id,omitempty"`
	CreatedBy   int                 `json:"created_by"`
	CreatedAt   CustomDateTime      `json:"created_at"`
}

// ============================================
// Expense Models (AP-related)
// ============================================

type ExpenseType struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
	GLAccountID *int   `json:"gl_account_id,omitempty"`
	IsActive    bool   `json:"is_active"`
}

type Expense struct {
	ID          int            `json:"id"`
	TypeID      int            `json:"type_id"` // References ExpenseType
	Date        CustomDate     `json:"date"`
	ExpenseDate CustomDate     `json:"expense_date,omitempty"`
	Amount      float64        `json:"amount"`
	Note        string         `json:"note,omitempty"`
	CashBoxID   int            `json:"cash_box_id"`
	VendorID    *int           `json:"vendor_id,omitempty"`
	BillID      *int           `json:"bill_id,omitempty"`
	IsMarked    bool           `json:"is_marked"`
	CreatedBy   int            `json:"created_by"`
	CreatedAt   CustomDateTime `json:"created_at"`
}

// ============================================
// Payment Type Models
// ============================================

type PaymentType struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Code          string `json:"code,omitempty"`
	IsInstallment bool   `json:"is_installment"`
	IsActive      bool   `json:"is_active"`
}

// ============================================
// Finance Request/Response DTOs
// ============================================

type CreateCashBoxRequest struct {
	Name        string `json:"name"`
	Code        string `json:"code,omitempty"`
	Currency    string `json:"currency"`
	WarehouseID *int   `json:"warehouse_id,omitempty"`
}

type UpdateCashBoxRequest struct {
	Name        *string `json:"name,omitempty"`
	Currency    *string `json:"currency,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
	WarehouseID *int    `json:"warehouse_id,omitempty"`
}

type CreateIncomeRequest struct {
	IncomeType  TransactionCategory `json:"income_type"`
	TypeID      int                 `json:"type_id"`
	Date        CustomDate          `json:"date"`
	ReceiptDate CustomDate          `json:"receipt_date,omitempty"`
	Amount      float64             `json:"amount"`
	Note        string              `json:"note,omitempty"`
	CashBoxID   int                 `json:"cash_box_id"`
	CustomerID  *int                `json:"customer_id,omitempty"`
	InvoiceID   *int                `json:"invoice_id,omitempty"`
}

type UpdateIncomeRequest struct {
	IncomeType  *TransactionCategory `json:"income_type,omitempty"`
	TypeID      *int                 `json:"type_id,omitempty"`
	Date        *CustomDate          `json:"date,omitempty"`
	ReceiptDate *CustomDate          `json:"receipt_date,omitempty"`
	Amount      *float64             `json:"amount,omitempty"`
	Note        *string              `json:"note,omitempty"`
	CashBoxID   *int                 `json:"cash_box_id,omitempty"`
}

type CreateExpenseRequest struct {
	TypeID      int        `json:"type_id"`
	Date        CustomDate `json:"date"`
	ExpenseDate CustomDate `json:"expense_date,omitempty"`
	Amount      float64    `json:"amount"`
	Note        string     `json:"note,omitempty"`
	CashBoxID   int        `json:"cash_box_id"`
	VendorID    *int       `json:"vendor_id,omitempty"`
	BillID      *int       `json:"bill_id,omitempty"`
}

type UpdateExpenseRequest struct {
	TypeID      *int        `json:"type_id,omitempty"`
	Date        *CustomDate `json:"date,omitempty"`
	ExpenseDate *CustomDate `json:"expense_date,omitempty"`
	Amount      *float64    `json:"amount,omitempty"`
	Note        *string     `json:"note,omitempty"`
	CashBoxID   *int        `json:"cash_box_id,omitempty"`
	IsMarked    *bool       `json:"is_marked,omitempty"`
}

type CreatePaymentTypeRequest struct {
	Name          string `json:"name"`
	Code          string `json:"code,omitempty"`
	IsInstallment bool   `json:"is_installment"`
}

type CreateIncomeTypeRequest struct {
	Name        string `json:"name"`
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
	GLAccountID *int   `json:"gl_account_id,omitempty"`
}

type CreateExpenseTypeRequest struct {
	Name        string `json:"name"`
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
	GLAccountID *int   `json:"gl_account_id,omitempty"`
}

// ============================================
// Finance List Filters
// ============================================

type IncomeListFilters struct {
	TypeID     *int        `json:"type_id,omitempty"`
	CashBoxID  *int        `json:"cash_box_id,omitempty"`
	CustomerID *int        `json:"customer_id,omitempty"`
	DateFrom   *CustomDate `json:"date_from,omitempty"`
	DateTo     *CustomDate `json:"date_to,omitempty"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
}

type ExpenseListFilters struct {
	TypeID    *int        `json:"type_id,omitempty"`
	CashBoxID *int        `json:"cash_box_id,omitempty"`
	VendorID  *int        `json:"vendor_id,omitempty"`
	DateFrom  *CustomDate `json:"date_from,omitempty"`
	DateTo    *CustomDate `json:"date_to,omitempty"`
	IsMarked  *bool       `json:"is_marked,omitempty"`
	Page      int         `json:"page"`
	PageSize  int         `json:"page_size"`
}
