package models

// ============================================
// General Ledger (GL) Enums
// ============================================

type GLAccountType string

const (
	GLAccountTypeAsset     GLAccountType = "ASSET"
	GLAccountTypeLiability GLAccountType = "LIABILITY"
	GLAccountTypeEquity    GLAccountType = "EQUITY"
	GLAccountTypeRevenue   GLAccountType = "REVENUE"
	GLAccountTypeExpense   GLAccountType = "EXPENSE"
)

type GLAccountSubType string

const (
	// Assets
	GLSubTypeCash        GLAccountSubType = "CASH"
	GLSubTypeBank        GLAccountSubType = "BANK"
	GLSubTypeReceivables GLAccountSubType = "RECEIVABLES"
	GLSubTypeInventory   GLAccountSubType = "INVENTORY"
	GLSubTypeFixedAssets GLAccountSubType = "FIXED_ASSETS"
	GLSubTypeOtherAssets GLAccountSubType = "OTHER_ASSETS"
	// Liabilities
	GLSubTypePayables     GLAccountSubType = "PAYABLES"
	GLSubTypeAccruedLiab  GLAccountSubType = "ACCRUED_LIABILITIES"
	GLSubTypeLongTermDebt GLAccountSubType = "LONG_TERM_DEBT"
	GLSubTypeOtherLiab    GLAccountSubType = "OTHER_LIABILITIES"
	// Equity
	GLSubTypeCapital      GLAccountSubType = "CAPITAL"
	GLSubTypeRetainedEarn GLAccountSubType = "RETAINED_EARNINGS"
	// Revenue
	GLSubTypeSales       GLAccountSubType = "SALES"
	GLSubTypeOtherIncome GLAccountSubType = "OTHER_INCOME"
	// Expense
	GLSubTypeCOGS         GLAccountSubType = "COGS"
	GLSubTypeOperatingExp GLAccountSubType = "OPERATING_EXPENSES"
	GLSubTypePayrollExp   GLAccountSubType = "PAYROLL_EXPENSES"
	GLSubTypeOtherExpense GLAccountSubType = "OTHER_EXPENSES"
)

type JournalEntryStatus string

const (
	JournalStatusDraft    JournalEntryStatus = "DRAFT"
	JournalStatusPending  JournalEntryStatus = "PENDING"
	JournalStatusPosted   JournalEntryStatus = "POSTED"
	JournalStatusReversed JournalEntryStatus = "REVERSED"
	JournalStatusVoided   JournalEntryStatus = "VOIDED"
)

type JournalEntryType string

const (
	JournalTypeManual     JournalEntryType = "MANUAL"
	JournalTypeAR         JournalEntryType = "AR"
	JournalTypeAP         JournalEntryType = "AP"
	JournalTypeInventory  JournalEntryType = "INVENTORY"
	JournalTypePayroll    JournalEntryType = "PAYROLL"
	JournalTypeBank       JournalEntryType = "BANK"
	JournalTypeRecurring  JournalEntryType = "RECURRING"
	JournalTypeAdjustment JournalEntryType = "ADJUSTMENT"
	JournalTypeClosing    JournalEntryType = "CLOSING"
)

type PeriodStatus string

const (
	PeriodStatusOpen   PeriodStatus = "OPEN"
	PeriodStatusClosed PeriodStatus = "CLOSED"
	PeriodStatusLocked PeriodStatus = "LOCKED"
)

// ============================================
// GL Account (Chart of Accounts)
// ============================================

type GLAccount struct {
	ID             int              `json:"id"`
	AccountCode    string           `json:"account_code"`
	AccountName    string           `json:"account_name"`
	AccountType    GLAccountType    `json:"account_type"`
	AccountSubType GLAccountSubType `json:"account_sub_type,omitempty"`
	ParentID       *int             `json:"parent_id,omitempty"`
	Description    string           `json:"description,omitempty"`
	Currency       string           `json:"currency"`
	IsActive       bool             `json:"is_active"`
	IsPostable     bool             `json:"is_postable"` // Can post transactions to this account
	IsBankAccount  bool             `json:"is_bank_account"`
	BankAccountID  *int             `json:"bank_account_id,omitempty"`
	NormalBalance  string           `json:"normal_balance"` // DEBIT or CREDIT
	OpeningBalance float64          `json:"opening_balance"`
	CurrentBalance float64          `json:"current_balance"`
	BudgetAmount   float64          `json:"budget_amount,omitempty"`
	DepartmentID   *int             `json:"department_id,omitempty"`
	CreatedAt      CustomDateTime   `json:"created_at"`
	UpdatedAt      CustomDateTime   `json:"updated_at"`
}

type GLAccountWithChildren struct {
	Account  GLAccount               `json:"account"`
	Children []GLAccountWithChildren `json:"children,omitempty"`
	Level    int                     `json:"level"`
}

// ============================================
// GL Fiscal Period
// ============================================

type GLFiscalYear struct {
	ID        int            `json:"id"`
	YearCode  string         `json:"year_code"`
	StartDate CustomDate     `json:"start_date"`
	EndDate   CustomDate     `json:"end_date"`
	IsCurrent bool           `json:"is_current"`
	IsClosed  bool           `json:"is_closed"`
	ClosedBy  *int           `json:"closed_by,omitempty"`
	ClosedAt  CustomDateTime `json:"closed_at,omitempty"`
	CreatedAt CustomDateTime `json:"created_at"`
}

type GLPeriod struct {
	ID           int            `json:"id"`
	FiscalYearID int            `json:"fiscal_year_id"`
	PeriodNumber int            `json:"period_number"` // 1-12 or 1-13 (13th for adjustments)
	PeriodName   string         `json:"period_name"`
	StartDate    CustomDate     `json:"start_date"`
	EndDate      CustomDate     `json:"end_date"`
	Status       PeriodStatus   `json:"status"`
	IsAdjustment bool           `json:"is_adjustment"` // True for period 13
	ClosedBy     *int           `json:"closed_by,omitempty"`
	ClosedAt     CustomDateTime `json:"closed_at,omitempty"`
}

// ============================================
// GL Journal Entry
// ============================================

type GLJournalEntry struct {
	ID              int                `json:"id"`
	JournalNumber   string             `json:"journal_number"`
	EntryDate       CustomDate         `json:"entry_date"`
	PostingDate     CustomDate         `json:"posting_date"`
	PeriodID        int                `json:"period_id"`
	EntryType       JournalEntryType   `json:"entry_type"`
	Status          JournalEntryStatus `json:"status"`
	Description     string             `json:"description"`
	Reference       string             `json:"reference,omitempty"`
	SourceDocument  string             `json:"source_document,omitempty"` // e.g., "AR-INV-001"
	SourceModule    string             `json:"source_module,omitempty"`   // e.g., "AR", "AP"
	SourceID        *int               `json:"source_id,omitempty"`       // ID of source document
	TotalDebit      float64            `json:"total_debit"`
	TotalCredit     float64            `json:"total_credit"`
	Currency        string             `json:"currency"`
	ExchangeRate    float64            `json:"exchange_rate"`
	IsRecurring     bool               `json:"is_recurring"`
	RecurringID     *int               `json:"recurring_id,omitempty"`
	ReversedEntryID *int               `json:"reversed_entry_id,omitempty"`
	AutoReverse     bool               `json:"auto_reverse"`
	AutoReverseDate CustomDate         `json:"auto_reverse_date,omitempty"`
	CreatedBy       int                `json:"created_by"`
	PostedBy        *int               `json:"posted_by,omitempty"`
	PostedAt        CustomDateTime     `json:"posted_at,omitempty"`
	CreatedAt       CustomDateTime     `json:"created_at"`
}

type GLJournalLine struct {
	ID           int            `json:"id"`
	JournalID    int            `json:"journal_id"`
	LineNumber   int            `json:"line_number"`
	AccountID    int            `json:"account_id"`
	Description  string         `json:"description,omitempty"`
	DebitAmount  float64        `json:"debit_amount"`
	CreditAmount float64        `json:"credit_amount"`
	DepartmentID *int           `json:"department_id,omitempty"`
	ProjectID    *int           `json:"project_id,omitempty"`
	Reference    string         `json:"reference,omitempty"`
	CreatedAt    CustomDateTime `json:"created_at"`
}

type GLJournalEntryWithLines struct {
	Entry       GLJournalEntry  `json:"entry"`
	Lines       []GLJournalLine `json:"lines"`
	PeriodName  string          `json:"period_name"`
	CreatorName string          `json:"creator_name"`
}

// ============================================
// GL Recurring Entry Template
// ============================================

type GLRecurringEntry struct {
	ID            int               `json:"id"`
	TemplateName  string            `json:"template_name"`
	Description   string            `json:"description"`
	Frequency     string            `json:"frequency"` // MONTHLY, QUARTERLY, YEARLY
	NextRunDate   CustomDate        `json:"next_run_date"`
	EndDate       CustomDate        `json:"end_date,omitempty"`
	LastRunDate   CustomDate        `json:"last_run_date,omitempty"`
	TotalRuns     int               `json:"total_runs"`
	IsActive      bool              `json:"is_active"`
	AutoReverse   bool              `json:"auto_reverse"`
	DaysToReverse int               `json:"days_to_reverse,omitempty"`
	Lines         []GLRecurringLine `json:"lines,omitempty"`
	CreatedBy     int               `json:"created_by"`
	CreatedAt     CustomDateTime    `json:"created_at"`
}

type GLRecurringLine struct {
	ID           int     `json:"id"`
	RecurringID  int     `json:"recurring_id"`
	LineNumber   int     `json:"line_number"`
	AccountID    int     `json:"account_id"`
	Description  string  `json:"description,omitempty"`
	DebitAmount  float64 `json:"debit_amount"`
	CreditAmount float64 `json:"credit_amount"`
}

// ============================================
// GL Budget
// ============================================

type GLBudget struct {
	ID           int            `json:"id"`
	FiscalYearID int            `json:"fiscal_year_id"`
	BudgetName   string         `json:"budget_name"`
	Description  string         `json:"description,omitempty"`
	IsApproved   bool           `json:"is_approved"`
	ApprovedBy   *int           `json:"approved_by,omitempty"`
	ApprovedAt   CustomDateTime `json:"approved_at,omitempty"`
	CreatedBy    int            `json:"created_by"`
	CreatedAt    CustomDateTime `json:"created_at"`
}

type GLBudgetLine struct {
	ID           int     `json:"id"`
	BudgetID     int     `json:"budget_id"`
	AccountID    int     `json:"account_id"`
	PeriodID     int     `json:"period_id"`
	BudgetAmount float64 `json:"budget_amount"`
	Notes        string  `json:"notes,omitempty"`
}

// ============================================
// GL Reports
// ============================================

type TrialBalanceRow struct {
	AccountID     int           `json:"account_id"`
	AccountCode   string        `json:"account_code"`
	AccountName   string        `json:"account_name"`
	AccountType   GLAccountType `json:"account_type"`
	OpeningDebit  float64       `json:"opening_debit"`
	OpeningCredit float64       `json:"opening_credit"`
	PeriodDebit   float64       `json:"period_debit"`
	PeriodCredit  float64       `json:"period_credit"`
	ClosingDebit  float64       `json:"closing_debit"`
	ClosingCredit float64       `json:"closing_credit"`
	Level         int           `json:"level"`
}

type TrialBalanceReport struct {
	AsOfDate    CustomDate        `json:"as_of_date"`
	FiscalYear  string            `json:"fiscal_year"`
	Period      string            `json:"period"`
	Rows        []TrialBalanceRow `json:"rows"`
	TotalDebit  float64           `json:"total_debit"`
	TotalCredit float64           `json:"total_credit"`
}

type IncomeStatementRow struct {
	AccountID   int           `json:"account_id"`
	AccountCode string        `json:"account_code"`
	AccountName string        `json:"account_name"`
	AccountType GLAccountType `json:"account_type"`
	Amount      float64       `json:"amount"`
	Budget      float64       `json:"budget,omitempty"`
	Variance    float64       `json:"variance,omitempty"`
	PriorPeriod float64       `json:"prior_period,omitempty"`
	PriorYear   float64       `json:"prior_year,omitempty"`
	Level       int           `json:"level"`
	IsSubtotal  bool          `json:"is_subtotal"`
}

type IncomeStatementReport struct {
	PeriodFrom    CustomDate           `json:"period_from"`
	PeriodTo      CustomDate           `json:"period_to"`
	Revenue       []IncomeStatementRow `json:"revenue"`
	TotalRevenue  float64              `json:"total_revenue"`
	Expenses      []IncomeStatementRow `json:"expenses"`
	TotalExpenses float64              `json:"total_expenses"`
	NetIncome     float64              `json:"net_income"`
}

type BalanceSheetRow struct {
	AccountID   int           `json:"account_id"`
	AccountCode string        `json:"account_code"`
	AccountName string        `json:"account_name"`
	AccountType GLAccountType `json:"account_type"`
	Balance     float64       `json:"balance"`
	PriorYear   float64       `json:"prior_year,omitempty"`
	Level       int           `json:"level"`
	IsSubtotal  bool          `json:"is_subtotal"`
}

type BalanceSheetReport struct {
	AsOfDate         CustomDate        `json:"as_of_date"`
	Assets           []BalanceSheetRow `json:"assets"`
	TotalAssets      float64           `json:"total_assets"`
	Liabilities      []BalanceSheetRow `json:"liabilities"`
	TotalLiabilities float64           `json:"total_liabilities"`
	Equity           []BalanceSheetRow `json:"equity"`
	TotalEquity      float64           `json:"total_equity"`
}

type AccountActivityRow struct {
	Date         CustomDate `json:"date"`
	JournalNum   string     `json:"journal_number"`
	Description  string     `json:"description"`
	Reference    string     `json:"reference,omitempty"`
	Debit        float64    `json:"debit"`
	Credit       float64    `json:"credit"`
	Balance      float64    `json:"balance"`
	SourceModule string     `json:"source_module,omitempty"`
}

type AccountActivityReport struct {
	AccountID      int                  `json:"account_id"`
	AccountCode    string               `json:"account_code"`
	AccountName    string               `json:"account_name"`
	DateFrom       CustomDate           `json:"date_from"`
	DateTo         CustomDate           `json:"date_to"`
	OpeningBalance float64              `json:"opening_balance"`
	Activity       []AccountActivityRow `json:"activity"`
	ClosingBalance float64              `json:"closing_balance"`
}

// ============================================
// Request/Response DTOs
// ============================================

type CreateGLAccountRequest struct {
	AccountCode    string           `json:"account_code"`
	AccountName    string           `json:"account_name"`
	AccountType    GLAccountType    `json:"account_type"`
	AccountSubType GLAccountSubType `json:"account_sub_type,omitempty"`
	ParentID       *int             `json:"parent_id,omitempty"`
	Description    string           `json:"description,omitempty"`
	Currency       string           `json:"currency"`
	IsPostable     bool             `json:"is_postable"`
	IsBankAccount  bool             `json:"is_bank_account"`
	BankAccountID  *int             `json:"bank_account_id,omitempty"`
	NormalBalance  string           `json:"normal_balance"` // DEBIT or CREDIT
	OpeningBalance float64          `json:"opening_balance"`
	DepartmentID   *int             `json:"department_id,omitempty"`
}

type UpdateGLAccountRequest struct {
	AccountName    *string           `json:"account_name,omitempty"`
	AccountSubType *GLAccountSubType `json:"account_sub_type,omitempty"`
	ParentID       *int              `json:"parent_id,omitempty"`
	Description    *string           `json:"description,omitempty"`
	IsPostable     *bool             `json:"is_postable,omitempty"`
	IsActive       *bool             `json:"is_active,omitempty"`
	BudgetAmount   *float64          `json:"budget_amount,omitempty"`
	DepartmentID   *int              `json:"department_id,omitempty"`
}

type CreateJournalEntryRequest struct {
	EntryDate       string                     `json:"entry_date"`
	EntryType       JournalEntryType           `json:"entry_type"`
	Description     string                     `json:"description"`
	Reference       string                     `json:"reference,omitempty"`
	Currency        string                     `json:"currency"`
	ExchangeRate    float64                    `json:"exchange_rate"`
	AutoReverse     bool                       `json:"auto_reverse"`
	AutoReverseDate string                     `json:"auto_reverse_date,omitempty"`
	Lines           []CreateJournalLineRequest `json:"lines"`
}

type CreateJournalLineRequest struct {
	AccountID    int     `json:"account_id"`
	Description  string  `json:"description,omitempty"`
	DebitAmount  float64 `json:"debit_amount"`
	CreditAmount float64 `json:"credit_amount"`
	DepartmentID *int    `json:"department_id,omitempty"`
	ProjectID    *int    `json:"project_id,omitempty"`
	Reference    string  `json:"reference,omitempty"`
}

type CreateFiscalYearRequest struct {
	YearCode  string `json:"year_code"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type CreateRecurringEntryRequest struct {
	TemplateName  string                     `json:"template_name"`
	Description   string                     `json:"description"`
	Frequency     string                     `json:"frequency"` // MONTHLY, QUARTERLY, YEARLY
	NextRunDate   string                     `json:"next_run_date"`
	EndDate       string                     `json:"end_date,omitempty"`
	AutoReverse   bool                       `json:"auto_reverse"`
	DaysToReverse int                        `json:"days_to_reverse,omitempty"`
	Lines         []CreateJournalLineRequest `json:"lines"`
}

type GLAccountListFilters struct {
	AccountType *GLAccountType `json:"account_type,omitempty"`
	IsActive    *bool          `json:"is_active,omitempty"`
	IsPostable  *bool          `json:"is_postable,omitempty"`
	ParentID    *int           `json:"parent_id,omitempty"`
	Search      string         `json:"search,omitempty"`
	Page        int            `json:"page"`
	PageSize    int            `json:"page_size"`
}

type JournalEntryListFilters struct {
	PeriodID  *int                `json:"period_id,omitempty"`
	EntryType *JournalEntryType   `json:"entry_type,omitempty"`
	Status    *JournalEntryStatus `json:"status,omitempty"`
	DateFrom  string              `json:"date_from,omitempty"`
	DateTo    string              `json:"date_to,omitempty"`
	AccountID *int                `json:"account_id,omitempty"`
	Search    string              `json:"search,omitempty"`
	Page      int                 `json:"page"`
	PageSize  int                 `json:"page_size"`
}

type ReportFilters struct {
	FiscalYearID  *int   `json:"fiscal_year_id,omitempty"`
	PeriodID      *int   `json:"period_id,omitempty"`
	DateFrom      string `json:"date_from,omitempty"`
	DateTo        string `json:"date_to,omitempty"`
	AccountID     *int   `json:"account_id,omitempty"`
	DepartmentID  *int   `json:"department_id,omitempty"`
	IncludeBudget bool   `json:"include_budget"`
	ComparePrior  bool   `json:"compare_prior"`
}

// ============================================
// Validation
// ============================================

func ValidateGLAccount(v *Validator, req *CreateGLAccountRequest) {
	v.Check(req.AccountCode != "", "account_code", "Account code is required")
	v.Check(len(req.AccountCode) <= 20, "account_code", "Account code must be 20 characters or less")
	v.Check(req.AccountName != "", "account_name", "Account name is required")
	v.Check(req.AccountType != "", "account_type", "Account type is required")
	v.Check(req.NormalBalance == "DEBIT" || req.NormalBalance == "CREDIT", "normal_balance", "Normal balance must be DEBIT or CREDIT")
	if req.Currency == "" {
		req.Currency = "USD"
	}
}

func ValidateJournalEntry(v *Validator, req *CreateJournalEntryRequest) {
	v.Check(req.EntryDate != "", "entry_date", "Entry date is required")
	v.Check(req.Description != "", "description", "Description is required")
	v.Check(len(req.Lines) >= 2, "lines", "At least two lines are required")

	var totalDebit, totalCredit float64
	for i, line := range req.Lines {
		v.Check(line.AccountID > 0, "lines", "Account ID is required for all lines")
		v.Check(line.DebitAmount >= 0, "lines", "Debit amount must be non-negative")
		v.Check(line.CreditAmount >= 0, "lines", "Credit amount must be non-negative")
		v.Check(!(line.DebitAmount > 0 && line.CreditAmount > 0), "lines",
			"Line "+string(rune(i+1))+" cannot have both debit and credit")
		totalDebit += line.DebitAmount
		totalCredit += line.CreditAmount
	}

	// Check that debits equal credits (within small tolerance for floating point)
	diff := totalDebit - totalCredit
	if diff < 0 {
		diff = -diff
	}
	v.Check(diff < 0.01, "lines", "Total debits must equal total credits")
}

func ValidateFiscalYear(v *Validator, req *CreateFiscalYearRequest) {
	v.Check(req.YearCode != "", "year_code", "Year code is required")
	v.Check(req.StartDate != "", "start_date", "Start date is required")
	v.Check(req.EndDate != "", "end_date", "End date is required")
}

func ValidateRecurringEntry(v *Validator, req *CreateRecurringEntryRequest) {
	v.Check(req.TemplateName != "", "template_name", "Template name is required")
	v.Check(req.Frequency != "", "frequency", "Frequency is required")
	v.Check(req.NextRunDate != "", "next_run_date", "Next run date is required")
	v.Check(len(req.Lines) >= 2, "lines", "At least two lines are required")
}
