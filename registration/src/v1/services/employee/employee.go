package employee

import (
	"context"
	"errors"
	"fmt"

	"github.com/anas-dev-92/FoodHive/core/postgres"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/models"
	"github.com/jackc/pgx/v5"
)

var (
	ErrNotFound       = errors.New("employee not found")
	ErrDuplicateEmail = errors.New("email already exists")
	ErrRoleNotFound   = errors.New("role not found")
)

type EmployeeService interface {
	Create(ctx context.Context, req models.CreateEmployeeRequest, createdBy int) (int, error)
	GetByID(ctx context.Context, id int) (*models.EmployeeInfo, error)
	Update(ctx context.Context, id int, req models.UpdateEmployeeRequest) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, filters models.EmployeeListFilters) ([]models.EmployeeInfo, int64, error)
}

type employeeServiceImpl struct {
	db postgres.Connection
}

func New(db postgres.Connection) EmployeeService {
	return &employeeServiceImpl{db: db}
}

func (s *employeeServiceImpl) Create(ctx context.Context, req models.CreateEmployeeRequest, createdBy int) (int, error) {
	if req.Email == "" {
		return 0, errors.New("email is required")
	}

	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return 0, fmt.Errorf("could not start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Check if role exists
	if req.RoleID > 0 {
		var existingRoleID int
		err = tx.QueryRow(ctx, `SELECT id FROM roles WHERE id = $1`, req.RoleID).Scan(&existingRoleID)
		if err != nil {
			if err == pgx.ErrNoRows {
				return 0, ErrRoleNotFound
			}
			return 0, fmt.Errorf("failed to check role: %w", err)
		}
	}

	// Insert employee - simplified schema
	status := req.Status
	if status == "" {
		status = "CONTINUED"
	}

	employeeQuery := `
		INSERT INTO employees (
			email, password, english_name, arabic_name, 
			nationality, phone, role_id, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`

	var employeeID int
	err = tx.QueryRow(ctx, employeeQuery,
		req.Email,
		req.Password,
		req.EnglishName,
		req.ArabicName,
		req.Nationality,
		req.Phone,
		req.RoleID,
		status,
	).Scan(&employeeID)
	if err != nil {
		return 0, fmt.Errorf("failed to create employee: %w", err)
	}

	// Note: Contract, Details, Finances, and Address tables were removed in simplified schema
	// These can be added back later if needed

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return employeeID, nil
}

func (s *employeeServiceImpl) GetByID(ctx context.Context, id int) (*models.EmployeeInfo, error) {
	query := `
		SELECT 
			e.id, e.email, e.account_status, e.english_name, e.arabic_name,
			e.nationality, e.phone, e.date_of_birth, e.status, e.contract_id, e.role_id,
			e.department_id, e.warehouse_id, e.created_at, e.updated_at,
			COALESCE(c.id, 0), c.start_date, c.end_date, c.contract_type,
			COALESCE(ed.id, 0), ed.gender, ed.job_title, ed.major_study, ed.notes,
			ed.passport_number, ed.national_id, ed.is_retired, ed.is_married, ed.number_of_children,
			COALESCE(ef.id, 0), ef.base_salary, ef.years_of_service, ef.academic_allowance,
			ef.degree_allowance, ef.position_allowance, ef.profession_allowance,
			ef.transport_allowance, ef.housing_allowance, ef.bank_account_number, ef.bank_name
		FROM employees e
		LEFT JOIN contracts c ON e.contract_id = c.id
		LEFT JOIN employee_details ed ON ed.employee_id = e.id
		LEFT JOIN employee_finances ef ON ef.employee_id = e.id
		WHERE e.id = $1`

	var info models.EmployeeInfo
	var contractID, detailsID, financesID int

	err := s.db.QueryRow(ctx, query, id).Scan(
		&info.Employee.ID,
		&info.Employee.Email,
		&info.Employee.AccountStatus,
		&info.Employee.EnglishName,
		&info.Employee.ArabicName,
		&info.Employee.Nationality,
		&info.Employee.Phone,
		&info.Employee.DateOfBirth,
		&info.Employee.Status,
		&info.Employee.ContractID,
		&info.Employee.RoleID,
		&info.Employee.DepartmentID,
		&info.Employee.WarehouseID,
		&info.Employee.CreatedAt,
		&info.Employee.UpdatedAt,
		&contractID,
		&info.Contract.StartDate,
		&info.Contract.EndDate,
		&info.Contract.ContractType,
		&detailsID,
		&info.EmployeeDetails.Gender,
		&info.EmployeeDetails.JobTitle,
		&info.EmployeeDetails.MajorStudy,
		&info.EmployeeDetails.Notes,
		&info.EmployeeDetails.PassportNumber,
		&info.EmployeeDetails.NationalID,
		&info.EmployeeDetails.IsRetired,
		&info.EmployeeDetails.IsMarried,
		&info.EmployeeDetails.NumberOfChildren,
		&financesID,
		&info.EmployeeFinances.BaseSalary,
		&info.EmployeeFinances.YearsOfService,
		&info.EmployeeFinances.AcademicAllowance,
		&info.EmployeeFinances.DegreeAllowance,
		&info.EmployeeFinances.PositionAllowance,
		&info.EmployeeFinances.ProfessionAllowance,
		&info.EmployeeFinances.TransportAllowance,
		&info.EmployeeFinances.HousingAllowance,
		&info.EmployeeFinances.BankAccountNumber,
		&info.EmployeeFinances.BankName,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get employee: %w", err)
	}

	info.Contract.ID = contractID
	info.EmployeeDetails.ID = detailsID
	info.EmployeeDetails.EmployeeID = id
	info.EmployeeFinances.ID = financesID
	info.EmployeeFinances.EmployeeID = id

	// Fetch address
	addressQuery := `
		SELECT id, address_line1, address_line2, city, state, country,
			postal_code, house, avenue, neighborhood, emergency_phone_number
		FROM addresses
		WHERE entity_type = 'employee' AND entity_id = $1`

	err = s.db.QueryRow(ctx, addressQuery, id).Scan(
		&info.Address.ID,
		&info.Address.AddressLine1,
		&info.Address.AddressLine2,
		&info.Address.City,
		&info.Address.State,
		&info.Address.Country,
		&info.Address.PostalCode,
		&info.Address.House,
		&info.Address.Avenue,
		&info.Address.Neighborhood,
		&info.Address.EmergencyPhoneNumber,
	)
	if err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("failed to get address: %w", err)
	}
	info.Address.EntityType = "employee"
	info.Address.EntityID = id

	return &info, nil
}

func (s *employeeServiceImpl) Update(ctx context.Context, id int, req models.UpdateEmployeeRequest) error {
	query := `
		UPDATE employees SET
			english_name = COALESCE($2, english_name),
			arabic_name = COALESCE($3, arabic_name),
			nationality = COALESCE($4, nationality),
			phone = COALESCE($5, phone),
			date_of_birth = COALESCE($6, date_of_birth),
			status = COALESCE($7, status),
			role_id = COALESCE($8, role_id),
			department_id = COALESCE($9, department_id),
			warehouse_id = COALESCE($10, warehouse_id),
			updated_at = NOW()
		WHERE id = $1`

	result, err := s.db.Exec(ctx, query,
		id,
		req.EnglishName,
		req.ArabicName,
		req.Nationality,
		req.Phone,
		req.DateOfBirth,
		req.Status,
		req.RoleID,
		req.DepartmentID,
		req.WarehouseID,
	)
	if err != nil {
		return fmt.Errorf("failed to update employee: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *employeeServiceImpl) Delete(ctx context.Context, id int) error {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Delete related records
	_, err = tx.Exec(ctx, `DELETE FROM employee_finances WHERE employee_id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete finances: %w", err)
	}

	_, err = tx.Exec(ctx, `DELETE FROM addresses WHERE entity_type = 'employee' AND entity_id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete address: %w", err)
	}

	_, err = tx.Exec(ctx, `DELETE FROM employee_details WHERE employee_id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete details: %w", err)
	}

	result, err := tx.Exec(ctx, `DELETE FROM employees WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete employee: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}

func (s *employeeServiceImpl) List(ctx context.Context, filters models.EmployeeListFilters) ([]models.EmployeeInfo, int64, error) {
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PageSize < 1 {
		filters.PageSize = 20
	}

	// Count query - basic columns only
	countQuery := `SELECT COUNT(*) FROM employees WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if filters.Search != "" {
		countQuery += fmt.Sprintf(" AND (english_name ILIKE $%d OR arabic_name ILIKE $%d OR email ILIKE $%d)", argIndex, argIndex, argIndex)
		args = append(args, "%"+filters.Search+"%")
		argIndex++
	}
	if filters.RoleID != nil {
		countQuery += fmt.Sprintf(" AND role_id = $%d", argIndex)
		args = append(args, *filters.RoleID)
		argIndex++
	}
	if filters.Status != nil {
		countQuery += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, *filters.Status)
		argIndex++
	}

	var total int64
	err := s.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count: %w", err)
	}

	// Data query - simple schema
	query := `
		SELECT e.id, e.email, COALESCE(e.english_name, ''), COALESCE(e.arabic_name, ''), 
			COALESCE(e.phone, ''), COALESCE(e.nationality, ''), COALESCE(e.role_id, 0),
			COALESCE(e.status, 'CONTINUED'), e.created_at, e.updated_at
		FROM employees e WHERE 1=1`

	args = []interface{}{}
	argIndex = 1

	if filters.Search != "" {
		query += fmt.Sprintf(" AND (english_name ILIKE $%d OR arabic_name ILIKE $%d OR email ILIKE $%d)", argIndex, argIndex, argIndex)
		args = append(args, "%"+filters.Search+"%")
		argIndex++
	}
	if filters.RoleID != nil {
		query += fmt.Sprintf(" AND role_id = $%d", argIndex)
		args = append(args, *filters.RoleID)
		argIndex++
	}
	if filters.Status != nil {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, *filters.Status)
		argIndex++
	}

	query += " ORDER BY english_name"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, filters.PageSize, (filters.Page-1)*filters.PageSize)

	rows := s.db.Query(ctx, query, args...)
	defer rows.Close()

	var results []models.EmployeeInfo
	for rows.Next() {
		var info models.EmployeeInfo
		var status string
		err := rows.Scan(
			&info.Employee.ID,
			&info.Employee.Email,
			&info.Employee.EnglishName,
			&info.Employee.ArabicName,
			&info.Employee.Phone,
			&info.Employee.Nationality,
			&info.Employee.RoleID,
			&status,
			&info.Employee.CreatedAt,
			&info.Employee.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan: %w", err)
		}
		info.Employee.Status = models.EmployeeStatus(status)
		results = append(results, info)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error reading rows: %w", err)
	}

	return results, total, nil
}
