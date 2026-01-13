package finance

import (
	"context"
	"errors"
	"fmt"

	"github.com/anas-dev-92/FoodHive/core/postgres"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/models"
	"github.com/jackc/pgx/v5"
)

var (
	ErrPaymentTypeNotFound = errors.New("payment type not found")
)

type PaymentTypeService interface {
	Create(ctx context.Context, req models.CreatePaymentTypeRequest) (int, error)
	GetByID(ctx context.Context, id int) (*models.PaymentType, error)
	Update(ctx context.Context, id int, name *string, isInstallment *bool) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context) ([]models.PaymentType, error)
}

type paymentTypeServiceImpl struct {
	db postgres.Connection
}

func NewPaymentTypeService(db postgres.Connection) PaymentTypeService {
	return &paymentTypeServiceImpl{db: db}
}

func (s *paymentTypeServiceImpl) Create(ctx context.Context, req models.CreatePaymentTypeRequest) (int, error) {
	query := `
		INSERT INTO payment_types (name, code, is_installment)
		VALUES ($1, $2, $3)
		RETURNING id`

	var id int
	err := s.db.QueryRow(ctx, query,
		req.Name,
		req.Code,
		req.IsInstallment,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create payment type: %w", err)
	}

	return id, nil
}

func (s *paymentTypeServiceImpl) GetByID(ctx context.Context, id int) (*models.PaymentType, error) {
	query := `
		SELECT id, name, code, is_installment, is_active
		FROM payment_types WHERE id = $1`

	var pt models.PaymentType
	err := s.db.QueryRow(ctx, query, id).Scan(
		&pt.ID,
		&pt.Name,
		&pt.Code,
		&pt.IsInstallment,
		&pt.IsActive,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrPaymentTypeNotFound
		}
		return nil, fmt.Errorf("failed to get payment type: %w", err)
	}

	return &pt, nil
}

func (s *paymentTypeServiceImpl) Update(ctx context.Context, id int, name *string, isInstallment *bool) error {
	query := `
		UPDATE payment_types SET
			name = COALESCE($2, name),
			is_installment = COALESCE($3, is_installment)
		WHERE id = $1`

	result, err := s.db.Exec(ctx, query, id, name, isInstallment)
	if err != nil {
		return fmt.Errorf("failed to update payment type: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrPaymentTypeNotFound
	}

	return nil
}

func (s *paymentTypeServiceImpl) Delete(ctx context.Context, id int) error {
	// Soft delete
	query := `UPDATE payment_types SET is_active = false WHERE id = $1`
	result, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete payment type: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrPaymentTypeNotFound
	}

	return nil
}

func (s *paymentTypeServiceImpl) List(ctx context.Context) ([]models.PaymentType, error) {
	query := `
		SELECT id, name, code, is_installment, is_active
		FROM payment_types WHERE is_active = true
		ORDER BY name`

	rows := s.db.Query(ctx, query)
	defer rows.Close()

	var results []models.PaymentType
	for rows.Next() {
		var pt models.PaymentType
		err := rows.Scan(
			&pt.ID,
			&pt.Name,
			&pt.Code,
			&pt.IsInstallment,
			&pt.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}
		results = append(results, pt)
	}

	return results, nil
}

// ============================================
// Income Type Service
// ============================================

type IncomeTypeService interface {
	Create(ctx context.Context, req models.CreateIncomeTypeRequest) (int, error)
	GetByID(ctx context.Context, id int) (*models.IncomeType, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context) ([]models.IncomeType, error)
}

type incomeTypeServiceImpl struct {
	db postgres.Connection
}

func NewIncomeTypeService(db postgres.Connection) IncomeTypeService {
	return &incomeTypeServiceImpl{db: db}
}

func (s *incomeTypeServiceImpl) Create(ctx context.Context, req models.CreateIncomeTypeRequest) (int, error) {
	query := `
		INSERT INTO income_types (name, code, description, gl_account_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	var id int
	err := s.db.QueryRow(ctx, query,
		req.Name,
		req.Code,
		req.Description,
		req.GLAccountID,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create income type: %w", err)
	}

	return id, nil
}

func (s *incomeTypeServiceImpl) GetByID(ctx context.Context, id int) (*models.IncomeType, error) {
	query := `
		SELECT id, name, code, description, gl_account_id, is_active
		FROM income_types WHERE id = $1`

	var it models.IncomeType
	err := s.db.QueryRow(ctx, query, id).Scan(
		&it.ID,
		&it.Name,
		&it.Code,
		&it.Description,
		&it.GLAccountID,
		&it.IsActive,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("income type not found")
		}
		return nil, fmt.Errorf("failed to get income type: %w", err)
	}

	return &it, nil
}

func (s *incomeTypeServiceImpl) Delete(ctx context.Context, id int) error {
	query := `UPDATE income_types SET is_active = false WHERE id = $1`
	result, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete: %w", err)
	}
	if result.RowsAffected() == 0 {
		return errors.New("income type not found")
	}
	return nil
}

func (s *incomeTypeServiceImpl) List(ctx context.Context) ([]models.IncomeType, error) {
	query := `
		SELECT id, name, code, description, gl_account_id, is_active
		FROM income_types WHERE is_active = true
		ORDER BY name`

	rows := s.db.Query(ctx, query)
	defer rows.Close()

	var results []models.IncomeType
	for rows.Next() {
		var it models.IncomeType
		err := rows.Scan(
			&it.ID,
			&it.Name,
			&it.Code,
			&it.Description,
			&it.GLAccountID,
			&it.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}
		results = append(results, it)
	}

	return results, nil
}

// ============================================
// Expense Type Service
// ============================================

type ExpenseTypeService interface {
	Create(ctx context.Context, req models.CreateExpenseTypeRequest) (int, error)
	GetByID(ctx context.Context, id int) (*models.ExpenseType, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context) ([]models.ExpenseType, error)
}

type expenseTypeServiceImpl struct {
	db postgres.Connection
}

func NewExpenseTypeService(db postgres.Connection) ExpenseTypeService {
	return &expenseTypeServiceImpl{db: db}
}

func (s *expenseTypeServiceImpl) Create(ctx context.Context, req models.CreateExpenseTypeRequest) (int, error) {
	query := `
		INSERT INTO expense_types (name, code, description, gl_account_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	var id int
	err := s.db.QueryRow(ctx, query,
		req.Name,
		req.Code,
		req.Description,
		req.GLAccountID,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create expense type: %w", err)
	}

	return id, nil
}

func (s *expenseTypeServiceImpl) GetByID(ctx context.Context, id int) (*models.ExpenseType, error) {
	query := `
		SELECT id, name, code, description, gl_account_id, is_active
		FROM expense_types WHERE id = $1`

	var et models.ExpenseType
	err := s.db.QueryRow(ctx, query, id).Scan(
		&et.ID,
		&et.Name,
		&et.Code,
		&et.Description,
		&et.GLAccountID,
		&et.IsActive,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("expense type not found")
		}
		return nil, fmt.Errorf("failed to get expense type: %w", err)
	}

	return &et, nil
}

func (s *expenseTypeServiceImpl) Delete(ctx context.Context, id int) error {
	query := `UPDATE expense_types SET is_active = false WHERE id = $1`
	result, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete: %w", err)
	}
	if result.RowsAffected() == 0 {
		return errors.New("expense type not found")
	}
	return nil
}

func (s *expenseTypeServiceImpl) List(ctx context.Context) ([]models.ExpenseType, error) {
	query := `
		SELECT id, name, code, description, gl_account_id, is_active
		FROM expense_types WHERE is_active = true
		ORDER BY name`

	rows := s.db.Query(ctx, query)
	defer rows.Close()

	var results []models.ExpenseType
	for rows.Next() {
		var et models.ExpenseType
		err := rows.Scan(
			&et.ID,
			&et.Name,
			&et.Code,
			&et.Description,
			&et.GLAccountID,
			&et.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}
		results = append(results, et)
	}

	return results, nil
}
