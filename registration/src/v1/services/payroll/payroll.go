package payroll

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/anas-dev-92/FoodHive/core/postgres"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/models"
	"github.com/jackc/pgx/v5"
)

var (
	ErrNotFound      = errors.New("payroll not found")
	ErrEditConflict  = errors.New("edit conflict")
	ErrInvalidStatus = errors.New("invalid payroll status for this operation")
)

type PayrollService interface {
	Create(ctx context.Context, req models.CreatePayrollRequest, createdBy int) (int, error)
	GetByID(ctx context.Context, id int) (*models.PayrollWithLines, error)
	Update(ctx context.Context, id int, req models.UpdatePayrollRequest) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, filters models.PayrollListFilters) ([]models.Payroll, int64, error)
	AddLine(ctx context.Context, payrollID int, req models.AddPayrollLineRequest) (int, error)
	RemoveLine(ctx context.Context, lineID int) error
	CalculatePayroll(ctx context.Context, id int) error
	ApprovePayroll(ctx context.Context, id int, approvedBy int) error
}

type payrollServiceImpl struct {
	db postgres.Connection
}

func New(db postgres.Connection) PayrollService {
	return &payrollServiceImpl{db: db}
}

func (s *payrollServiceImpl) Create(ctx context.Context, req models.CreatePayrollRequest, createdBy int) (int, error) {
	query := `
		INSERT INTO payrolls (
			created_by, department_id, payroll_period, pay_date, status, type,
			note, include_insurance, include_bonus, include_overtime, include_tax
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id`

	var id int
	err := s.db.QueryRow(ctx, query,
		createdBy,
		req.DepartmentID,
		req.PayrollPeriod,
		time.Time(req.PayDate),
		models.PayrollStatusDraft,
		req.Type,
		req.Note,
		req.IncludeInsurance,
		req.IncludeBonus,
		req.IncludeOvertime,
		req.IncludeTax,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create payroll: %w", err)
	}

	return id, nil
}

func (s *payrollServiceImpl) GetByID(ctx context.Context, id int) (*models.PayrollWithLines, error) {
	query := `
		SELECT id, created_by, department_id, payroll_period, pay_date, status, type,
			note, total_gross_pay, total_deductions, total_net_pay,
			include_insurance, include_bonus, include_overtime, include_tax,
			created_at, updated_at, approved_by, approved_at
		FROM payrolls WHERE id = $1`

	var result models.PayrollWithLines
	var p models.Payroll

	err := s.db.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.CreatedBy, &p.DepartmentID, &p.PayrollPeriod, &p.PayDate, &p.Status, &p.Type,
		&p.Note, &p.TotalGrossPay, &p.TotalDeductions, &p.TotalNetPay,
		&p.IncludeInsurance, &p.IncludeBonus, &p.IncludeOvertime, &p.IncludeTax,
		&p.CreatedAt, &p.UpdatedAt, &p.ApprovedBy, &p.ApprovedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get payroll: %w", err)
	}
	result.Payroll = p

	// Get payroll lines
	linesQuery := `
		SELECT id, payroll_id, employee_id, base_salary, allowances, bonuses,
			overtime_pay, overtime_hours, gross_pay, tax_deduction,
			insurance_deduction, other_deductions, net_pay, notes
		FROM payroll_lines WHERE payroll_id = $1
		ORDER BY id`

	rows := s.db.Query(ctx, linesQuery, id)
	defer rows.Close()

	for rows.Next() {
		var line models.PayrollLine
		err := rows.Scan(
			&line.ID, &line.PayrollID, &line.EmployeeID, &line.BaseSalary,
			&line.Allowances, &line.Bonuses, &line.OvertimePay, &line.OvertimeHours,
			&line.GrossPay, &line.TaxDeduction, &line.InsuranceDeduction,
			&line.OtherDeductions, &line.NetPay, &line.Notes,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payroll line: %w", err)
		}
		result.Lines = append(result.Lines, line)
	}

	return &result, nil
}

func (s *payrollServiceImpl) Update(ctx context.Context, id int, req models.UpdatePayrollRequest) error {
	query := `
		UPDATE payrolls SET
			department_id = COALESCE($2, department_id),
			pay_date = COALESCE($3, pay_date),
			note = COALESCE($4, note),
			include_insurance = COALESCE($5, include_insurance),
			include_bonus = COALESCE($6, include_bonus),
			include_overtime = COALESCE($7, include_overtime),
			include_tax = COALESCE($8, include_tax),
			updated_at = NOW()
		WHERE id = $1 AND status = 'DRAFT'`

	result, err := s.db.Exec(ctx, query,
		id,
		req.DepartmentID,
		req.PayDate,
		req.Note,
		req.IncludeInsurance,
		req.IncludeBonus,
		req.IncludeOvertime,
		req.IncludeTax,
	)
	if err != nil {
		return fmt.Errorf("failed to update payroll: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *payrollServiceImpl) Delete(ctx context.Context, id int) error {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Delete lines first
	_, err = tx.Exec(ctx, `DELETE FROM payroll_lines WHERE payroll_id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete payroll lines: %w", err)
	}

	// Delete payroll (only if DRAFT)
	result, err := tx.Exec(ctx, `DELETE FROM payrolls WHERE id = $1 AND status = 'DRAFT'`, id)
	if err != nil {
		return fmt.Errorf("failed to delete payroll: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}

func (s *payrollServiceImpl) List(ctx context.Context, filters models.PayrollListFilters) ([]models.Payroll, int64, error) {
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PageSize < 1 {
		filters.PageSize = 20
	}

	// Count query
	countQuery := `SELECT COUNT(*) FROM payrolls WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if filters.DepartmentID != nil {
		countQuery += fmt.Sprintf(" AND department_id = $%d", argIndex)
		args = append(args, *filters.DepartmentID)
		argIndex++
	}
	if filters.PayrollPeriod != nil {
		countQuery += fmt.Sprintf(" AND payroll_period = $%d", argIndex)
		args = append(args, *filters.PayrollPeriod)
		argIndex++
	}
	if filters.Status != nil {
		countQuery += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, *filters.Status)
		argIndex++
	}
	if filters.Type != nil {
		countQuery += fmt.Sprintf(" AND type = $%d", argIndex)
		args = append(args, *filters.Type)
		argIndex++
	}

	var total int64
	err := s.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count: %w", err)
	}

	// Data query
	query := `
		SELECT id, created_by, department_id, payroll_period, pay_date, status, type,
			note, total_gross_pay, total_deductions, total_net_pay,
			include_insurance, include_bonus, include_overtime, include_tax,
			created_at, updated_at, approved_by, approved_at
		FROM payrolls WHERE 1=1`

	args = []interface{}{}
	argIndex = 1

	if filters.DepartmentID != nil {
		query += fmt.Sprintf(" AND department_id = $%d", argIndex)
		args = append(args, *filters.DepartmentID)
		argIndex++
	}
	if filters.PayrollPeriod != nil {
		query += fmt.Sprintf(" AND payroll_period = $%d", argIndex)
		args = append(args, *filters.PayrollPeriod)
		argIndex++
	}
	if filters.Status != nil {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, *filters.Status)
		argIndex++
	}
	if filters.Type != nil {
		query += fmt.Sprintf(" AND type = $%d", argIndex)
		args = append(args, *filters.Type)
		argIndex++
	}

	query += " ORDER BY created_at DESC"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, filters.PageSize, (filters.Page-1)*filters.PageSize)

	rows := s.db.Query(ctx, query, args...)
	defer rows.Close()

	var results []models.Payroll
	for rows.Next() {
		var p models.Payroll
		err := rows.Scan(
			&p.ID, &p.CreatedBy, &p.DepartmentID, &p.PayrollPeriod, &p.PayDate, &p.Status, &p.Type,
			&p.Note, &p.TotalGrossPay, &p.TotalDeductions, &p.TotalNetPay,
			&p.IncludeInsurance, &p.IncludeBonus, &p.IncludeOvertime, &p.IncludeTax,
			&p.CreatedAt, &p.UpdatedAt, &p.ApprovedBy, &p.ApprovedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan: %w", err)
		}
		results = append(results, p)
	}

	return results, total, nil
}

func (s *payrollServiceImpl) AddLine(ctx context.Context, payrollID int, req models.AddPayrollLineRequest) (int, error) {
	// Get payroll to verify status
	var status string
	err := s.db.QueryRow(ctx, `SELECT status FROM payrolls WHERE id = $1`, payrollID).Scan(&status)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, ErrNotFound
		}
		return 0, fmt.Errorf("failed to check payroll: %w", err)
	}
	if status != string(models.PayrollStatusDraft) {
		return 0, ErrInvalidStatus
	}

	// Get employee finances
	var baseSalary, academicAllowance, degreeAllowance, positionAllowance float64
	var professionAllowance, transportAllowance, housingAllowance float64
	var overtimeRate, taxDeduction, insuranceDeduction float64

	finQuery := `
		SELECT COALESCE(base_salary, 0), COALESCE(academic_allowance, 0),
			COALESCE(degree_allowance, 0), COALESCE(position_allowance, 0),
			COALESCE(profession_allowance, 0), COALESCE(transport_allowance, 0),
			COALESCE(housing_allowance, 0), COALESCE(overtime_rate, 0),
			COALESCE(tax_deduction, 0), COALESCE(insurance_deduction, 0)
		FROM employee_finances WHERE employee_id = $1`

	err = s.db.QueryRow(ctx, finQuery, req.EmployeeID).Scan(
		&baseSalary, &academicAllowance, &degreeAllowance, &positionAllowance,
		&professionAllowance, &transportAllowance, &housingAllowance,
		&overtimeRate, &taxDeduction, &insuranceDeduction,
	)
	if err != nil && err != pgx.ErrNoRows {
		return 0, fmt.Errorf("failed to get employee finances: %w", err)
	}

	// Calculate
	allowances := academicAllowance + degreeAllowance + positionAllowance +
		professionAllowance + transportAllowance + housingAllowance
	overtimePay := req.OvertimeHours * overtimeRate
	grossPay := baseSalary + allowances + req.Bonuses + overtimePay
	totalDeductions := taxDeduction + insuranceDeduction + req.OtherDeductions
	netPay := grossPay - totalDeductions

	query := `
		INSERT INTO payroll_lines (
			payroll_id, employee_id, base_salary, allowances, bonuses,
			overtime_pay, overtime_hours, gross_pay, tax_deduction,
			insurance_deduction, other_deductions, net_pay, notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id`

	var id int
	err = s.db.QueryRow(ctx, query,
		payrollID, req.EmployeeID, baseSalary, allowances, req.Bonuses,
		overtimePay, req.OvertimeHours, grossPay, taxDeduction,
		insuranceDeduction, req.OtherDeductions, netPay, req.Notes,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to add payroll line: %w", err)
	}

	return id, nil
}

func (s *payrollServiceImpl) RemoveLine(ctx context.Context, lineID int) error {
	// Check if payroll is in draft status
	var status string
	err := s.db.QueryRow(ctx, `
		SELECT p.status FROM payrolls p
		JOIN payroll_lines pl ON pl.payroll_id = p.id
		WHERE pl.id = $1`, lineID).Scan(&status)
	if err != nil {
		if err == pgx.ErrNoRows {
			return ErrNotFound
		}
		return fmt.Errorf("failed to check status: %w", err)
	}
	if status != string(models.PayrollStatusDraft) {
		return ErrInvalidStatus
	}

	result, err := s.db.Exec(ctx, `DELETE FROM payroll_lines WHERE id = $1`, lineID)
	if err != nil {
		return fmt.Errorf("failed to remove line: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *payrollServiceImpl) CalculatePayroll(ctx context.Context, id int) error {
	// Sum all lines and update payroll totals
	query := `
		UPDATE payrolls SET
			total_gross_pay = (SELECT COALESCE(SUM(gross_pay), 0) FROM payroll_lines WHERE payroll_id = $1),
			total_deductions = (SELECT COALESCE(SUM(tax_deduction + insurance_deduction + other_deductions), 0) FROM payroll_lines WHERE payroll_id = $1),
			total_net_pay = (SELECT COALESCE(SUM(net_pay), 0) FROM payroll_lines WHERE payroll_id = $1),
			updated_at = NOW()
		WHERE id = $1 AND status = 'DRAFT'`

	result, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to calculate: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *payrollServiceImpl) ApprovePayroll(ctx context.Context, id int, approvedBy int) error {
	// Calculate first
	if err := s.CalculatePayroll(ctx, id); err != nil {
		return err
	}

	query := `
		UPDATE payrolls SET
			status = 'APPROVED',
			approved_by = $2,
			approved_at = NOW(),
			updated_at = NOW()
		WHERE id = $1 AND status = 'DRAFT'`

	result, err := s.db.Exec(ctx, query, id, approvedBy)
	if err != nil {
		return fmt.Errorf("failed to approve: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
