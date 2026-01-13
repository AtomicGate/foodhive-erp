package models

// ============================================
// Payroll Types
// ============================================

type PayrollType string

const (
	PayrollTypeLecturer  PayrollType = "LECTURER"
	PayrollTypeProfessor PayrollType = "PROFESSOR"
	PayrollTypeEmployee  PayrollType = "EMPLOYEE"
	PayrollTypeHourly    PayrollType = "HOURLY"
	PayrollTypeSalaried  PayrollType = "SALARIED"
)

type PayrollStatus string

const (
	PayrollStatusDraft     PayrollStatus = "DRAFT"
	PayrollStatusPending   PayrollStatus = "PENDING"
	PayrollStatusApproved  PayrollStatus = "APPROVED"
	PayrollStatusProcessed PayrollStatus = "PROCESSED"
	PayrollStatusPaid      PayrollStatus = "PAID"
	PayrollStatusCancelled PayrollStatus = "CANCELLED"
)

// ============================================
// Payroll Models
// ============================================

type Payroll struct {
	ID               int             `json:"id"`
	CreatedBy        int             `json:"created_by"`
	DepartmentID     *int            `json:"department_id,omitempty"`
	PayrollPeriod    string          `json:"payroll_period"` // e.g., "2026-01"
	PayDate          CustomDate      `json:"pay_date"`
	Status           PayrollStatus   `json:"status"`
	Type             PayrollType     `json:"type"`
	Note             string          `json:"note,omitempty"`
	TotalGrossPay    float64         `json:"total_gross_pay"`
	TotalDeductions  float64         `json:"total_deductions"`
	TotalNetPay      float64         `json:"total_net_pay"`
	IncludeInsurance bool            `json:"include_insurance"`
	IncludeBonus     bool            `json:"include_bonus"`
	IncludeOvertime  bool            `json:"include_overtime"`
	IncludeTax       bool            `json:"include_tax"`
	CreatedAt        CustomDateTime  `json:"created_at"`
	UpdatedAt        CustomDateTime  `json:"updated_at"`
	ApprovedBy       *int            `json:"approved_by,omitempty"`
	ApprovedAt       *CustomDateTime `json:"approved_at,omitempty"`
}

type PayrollLine struct {
	ID                 int     `json:"id"`
	PayrollID          int     `json:"payroll_id"`
	EmployeeID         int     `json:"employee_id"`
	BaseSalary         float64 `json:"base_salary"`
	Allowances         float64 `json:"allowances"`
	Bonuses            float64 `json:"bonuses"`
	OvertimePay        float64 `json:"overtime_pay"`
	OvertimeHours      float64 `json:"overtime_hours"`
	GrossPay           float64 `json:"gross_pay"`
	TaxDeduction       float64 `json:"tax_deduction"`
	InsuranceDeduction float64 `json:"insurance_deduction"`
	OtherDeductions    float64 `json:"other_deductions"`
	NetPay             float64 `json:"net_pay"`
	Notes              string  `json:"notes,omitempty"`
}

type PayrollWithLines struct {
	Payroll Payroll       `json:"payroll"`
	Lines   []PayrollLine `json:"lines"`
}

type PayrollLineWithEmployee struct {
	PayrollLine
	EmployeeName   string `json:"employee_name"`
	EmployeeCode   string `json:"employee_code,omitempty"`
	DepartmentName string `json:"department_name,omitempty"`
}

// ============================================
// Payroll Request/Response DTOs
// ============================================

type CreatePayrollRequest struct {
	DepartmentID     *int        `json:"department_id,omitempty"`
	PayrollPeriod    string      `json:"payroll_period"`
	PayDate          CustomDate  `json:"pay_date"`
	Type             PayrollType `json:"type"`
	Note             string      `json:"note,omitempty"`
	IncludeInsurance bool        `json:"include_insurance"`
	IncludeBonus     bool        `json:"include_bonus"`
	IncludeOvertime  bool        `json:"include_overtime"`
	IncludeTax       bool        `json:"include_tax"`
}

type UpdatePayrollRequest struct {
	DepartmentID     *int           `json:"department_id,omitempty"`
	PayDate          *CustomDate    `json:"pay_date,omitempty"`
	Note             *string        `json:"note,omitempty"`
	Status           *PayrollStatus `json:"status,omitempty"`
	IncludeInsurance *bool          `json:"include_insurance,omitempty"`
	IncludeBonus     *bool          `json:"include_bonus,omitempty"`
	IncludeOvertime  *bool          `json:"include_overtime,omitempty"`
	IncludeTax       *bool          `json:"include_tax,omitempty"`
}

type AddPayrollLineRequest struct {
	EmployeeID      int     `json:"employee_id"`
	OvertimeHours   float64 `json:"overtime_hours,omitempty"`
	Bonuses         float64 `json:"bonuses,omitempty"`
	OtherDeductions float64 `json:"other_deductions,omitempty"`
	Notes           string  `json:"notes,omitempty"`
}

type PayrollListFilters struct {
	DepartmentID  *int    `json:"department_id,omitempty"`
	PayrollPeriod *string `json:"payroll_period,omitempty"`
	Status        *string `json:"status,omitempty"`
	Type          *string `json:"type,omitempty"`
	Page          int     `json:"page"`
	PageSize      int     `json:"page_size"`
}
