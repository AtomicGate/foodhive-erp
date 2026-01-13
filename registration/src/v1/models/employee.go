package models

// ============================================
// Employee Status Types
// ============================================

type EmployeeStatus string

const (
	EmployeeStatusResign    EmployeeStatus = "RESIGN"
	EmployeeStatusContinued EmployeeStatus = "CONTINUED"
	EmployeeStatusOnLeave   EmployeeStatus = "ON_LEAVE"
	EmployeeStatusSuspended EmployeeStatus = "SUSPENDED"
)

type Gender string

const (
	GenderMale   Gender = "MALE"
	GenderFemale Gender = "FEMALE"
)

// ============================================
// Contract Types
// ============================================

type ContractType string

const (
	ContractTypeYearly    ContractType = "YEARLY"
	ContractTypeFull      ContractType = "FULL"
	ContractTypePartTime  ContractType = "PART_TIME"
	ContractTypeTemporary ContractType = "TEMPORARY"
)

// ============================================
// Employee Models
// ============================================

type Employee struct {
	ID            int            `json:"id"`
	Email         string         `json:"email"`
	Password      string         `json:"password,omitempty"`
	AccountStatus string         `json:"account_status"`
	EnglishName   string         `json:"english_name,omitempty"`
	ArabicName    string         `json:"arabic_name,omitempty"`
	Nationality   string         `json:"nationality,omitempty"`
	Phone         string         `json:"phone,omitempty"`
	DateOfBirth   CustomDate     `json:"date_of_birth,omitempty"`
	Status        EmployeeStatus `json:"status,omitempty"`
	ContractID    *int           `json:"contract_id,omitempty"`
	RoleID        int            `json:"role_id,omitempty"`
	DepartmentID  *int           `json:"department_id,omitempty"`
	WarehouseID   *int           `json:"warehouse_id,omitempty"`
	CreatedAt     CustomDateTime `json:"created_at,omitempty"`
	UpdatedAt     CustomDateTime `json:"updated_at,omitempty"`
}

type EmployeeDetails struct {
	ID               int    `json:"id"`
	EmployeeID       int    `json:"employee_id"`
	Gender           Gender `json:"gender,omitempty"`
	JobTitle         string `json:"job_title,omitempty"`
	MajorStudy       string `json:"major_study,omitempty"`
	Notes            string `json:"notes,omitempty"`
	PassportNumber   string `json:"passport_number,omitempty"`
	NationalID       string `json:"national_id,omitempty"`
	IsRetired        bool   `json:"is_retired,omitempty"`
	IsMarried        bool   `json:"is_married,omitempty"`
	NumberOfChildren int16  `json:"number_of_children,omitempty"`
}

type EmployeeFinances struct {
	ID                   int     `json:"id"`
	EmployeeID           int     `json:"employee_id"`
	BaseSalary           float64 `json:"base_salary"`
	YearsOfService       float64 `json:"years_of_service"`
	AcademicAllowance    float64 `json:"academic_allowance,omitempty"`
	DegreeAllowance      float64 `json:"degree_allowance,omitempty"`
	PositionAllowance    float64 `json:"position_allowance,omitempty"`
	ProfessionAllowance  float64 `json:"profession_allowance,omitempty"`
	IncentiveBonus       float64 `json:"incentive_bonus,omitempty"`
	TransportAllowance   float64 `json:"transport_allowance,omitempty"`
	HousingAllowance     float64 `json:"housing_allowance,omitempty"`
	OvertimeRate         float64 `json:"overtime_rate,omitempty"`
	TaxDeduction         float64 `json:"tax_deduction,omitempty"`
	InsuranceDeduction   float64 `json:"insurance_deduction,omitempty"`
	BankAccountNumber    string  `json:"bank_account_number,omitempty"`
	BankName             string  `json:"bank_name,omitempty"`
}

type Contract struct {
	ID           int          `json:"id"`
	StartDate    CustomDate   `json:"start_date,omitempty"`
	EndDate      CustomDate   `json:"end_date,omitempty"`
	ContractType ContractType `json:"contract_type,omitempty"`
	Notes        string       `json:"notes,omitempty"`
}

// ============================================
// Address Model (Used for Employees, Customers, Vendors)
// ============================================

type Address struct {
	ID                   int    `json:"id"`
	EntityType           string `json:"entity_type,omitempty"` // "employee", "customer", "vendor"
	EntityID             int    `json:"entity_id,omitempty"`
	AddressLine1         string `json:"address_line1,omitempty"`
	AddressLine2         string `json:"address_line2,omitempty"`
	City                 string `json:"city,omitempty"`
	State                string `json:"state,omitempty"`
	Country              string `json:"country,omitempty"`
	PostalCode           string `json:"postal_code,omitempty"`
	House                string `json:"house,omitempty"`
	Avenue               string `json:"avenue,omitempty"`
	Neighborhood         string `json:"neighborhood,omitempty"`
	EmergencyPhoneNumber string `json:"emergency_phone_number,omitempty"`
	IsDefault            bool   `json:"is_default,omitempty"`
}

// ============================================
// Composite Employee Info
// ============================================

type EmployeeInfo struct {
	Employee         Employee         `json:"employee"`
	Contract         Contract         `json:"contract,omitempty"`
	EmployeeDetails  EmployeeDetails  `json:"employee_details,omitempty"`
	Address          Address          `json:"address,omitempty"`
	EmployeeFinances EmployeeFinances `json:"employee_finances,omitempty"`
}

// ============================================
// Employee Request/Response DTOs
// ============================================

type CreateEmployeeRequest struct {
	Email        string         `json:"email"`
	Password     string         `json:"password"`
	EnglishName  string         `json:"english_name"`
	ArabicName   string         `json:"arabic_name,omitempty"`
	Nationality  string         `json:"nationality,omitempty"`
	Phone        string         `json:"phone,omitempty"`
	DateOfBirth  CustomDate     `json:"date_of_birth,omitempty"`
	Status       EmployeeStatus `json:"status,omitempty"`
	RoleID       int            `json:"role_id"`
	DepartmentID *int           `json:"department_id,omitempty"`
	WarehouseID  *int           `json:"warehouse_id,omitempty"`

	// Nested details
	Details  *CreateEmployeeDetailsRequest  `json:"details,omitempty"`
	Finances *CreateEmployeeFinancesRequest `json:"finances,omitempty"`
	Contract *CreateContractRequest         `json:"contract,omitempty"`
	Address  *CreateAddressRequest          `json:"address,omitempty"`
}

type CreateEmployeeDetailsRequest struct {
	Gender           Gender `json:"gender,omitempty"`
	JobTitle         string `json:"job_title,omitempty"`
	MajorStudy       string `json:"major_study,omitempty"`
	Notes            string `json:"notes,omitempty"`
	PassportNumber   string `json:"passport_number,omitempty"`
	NationalID       string `json:"national_id,omitempty"`
	IsMarried        bool   `json:"is_married,omitempty"`
	NumberOfChildren int16  `json:"number_of_children,omitempty"`
}

type CreateEmployeeFinancesRequest struct {
	BaseSalary          float64 `json:"base_salary"`
	AcademicAllowance   float64 `json:"academic_allowance,omitempty"`
	DegreeAllowance     float64 `json:"degree_allowance,omitempty"`
	PositionAllowance   float64 `json:"position_allowance,omitempty"`
	ProfessionAllowance float64 `json:"profession_allowance,omitempty"`
	TransportAllowance  float64 `json:"transport_allowance,omitempty"`
	HousingAllowance    float64 `json:"housing_allowance,omitempty"`
	BankAccountNumber   string  `json:"bank_account_number,omitempty"`
	BankName            string  `json:"bank_name,omitempty"`
}

type CreateContractRequest struct {
	StartDate    CustomDate   `json:"start_date"`
	EndDate      CustomDate   `json:"end_date,omitempty"`
	ContractType ContractType `json:"contract_type"`
	Notes        string       `json:"notes,omitempty"`
}

type CreateAddressRequest struct {
	AddressLine1         string `json:"address_line1,omitempty"`
	AddressLine2         string `json:"address_line2,omitempty"`
	City                 string `json:"city,omitempty"`
	State                string `json:"state,omitempty"`
	Country              string `json:"country,omitempty"`
	PostalCode           string `json:"postal_code,omitempty"`
	House                string `json:"house,omitempty"`
	Avenue               string `json:"avenue,omitempty"`
	Neighborhood         string `json:"neighborhood,omitempty"`
	EmergencyPhoneNumber string `json:"emergency_phone_number,omitempty"`
}

type UpdateEmployeeRequest struct {
	EnglishName  *string         `json:"english_name,omitempty"`
	ArabicName   *string         `json:"arabic_name,omitempty"`
	Nationality  *string         `json:"nationality,omitempty"`
	Phone        *string         `json:"phone,omitempty"`
	DateOfBirth  *CustomDate     `json:"date_of_birth,omitempty"`
	Status       *EmployeeStatus `json:"status,omitempty"`
	RoleID       *int            `json:"role_id,omitempty"`
	DepartmentID *int            `json:"department_id,omitempty"`
	WarehouseID  *int            `json:"warehouse_id,omitempty"`
	IsActive     *bool           `json:"is_active,omitempty"`
}

type EmployeeListFilters struct {
	Search       string  `json:"search,omitempty"`
	DepartmentID *int    `json:"department_id,omitempty"`
	RoleID       *int    `json:"role_id,omitempty"`
	Status       *string `json:"status,omitempty"`
	IsActive     *bool   `json:"is_active,omitempty"`
	Page         int     `json:"page"`
	PageSize     int     `json:"page_size"`
}
