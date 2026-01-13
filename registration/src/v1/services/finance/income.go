package finance

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
	ErrIncomeNotFound = errors.New("income not found")
)

type IncomeService interface {
	Create(ctx context.Context, req models.CreateIncomeRequest, createdBy int) (int, error)
	GetByID(ctx context.Context, id int) (*models.Income, error)
	Update(ctx context.Context, id int, req models.UpdateIncomeRequest) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, filters models.IncomeListFilters) ([]models.Income, int64, error)
}

type incomeServiceImpl struct {
	db postgres.Connection
}

func NewIncomeService(db postgres.Connection) IncomeService {
	return &incomeServiceImpl{db: db}
}

func (s *incomeServiceImpl) Create(ctx context.Context, req models.CreateIncomeRequest, createdBy int) (int, error) {
	query := `
		INSERT INTO incomes (
			income_type, type_id, date, receipt_date, amount,
			note, cash_box_id, customer_id, invoice_id, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id`

	var id int
	err := s.db.QueryRow(ctx, query,
		req.IncomeType,
		req.TypeID,
		time.Time(req.Date),
		time.Time(req.ReceiptDate),
		req.Amount,
		req.Note,
		req.CashBoxID,
		req.CustomerID,
		req.InvoiceID,
		createdBy,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create income: %w", err)
	}

	// Update cash box balance
	_, err = s.db.Exec(ctx, `
		UPDATE cash_boxes SET current_balance = current_balance + $1 WHERE id = $2`,
		req.Amount, req.CashBoxID)
	if err != nil {
		return 0, fmt.Errorf("failed to update cash box: %w", err)
	}

	return id, nil
}

func (s *incomeServiceImpl) GetByID(ctx context.Context, id int) (*models.Income, error) {
	query := `
		SELECT id, income_type, type_id, date, receipt_date, amount,
			note, cash_box_id, customer_id, invoice_id, created_by, created_at
		FROM incomes WHERE id = $1`

	var income models.Income
	err := s.db.QueryRow(ctx, query, id).Scan(
		&income.ID,
		&income.IncomeType,
		&income.TypeID,
		&income.Date,
		&income.ReceiptDate,
		&income.Amount,
		&income.Note,
		&income.CashBoxID,
		&income.CustomerID,
		&income.InvoiceID,
		&income.CreatedBy,
		&income.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrIncomeNotFound
		}
		return nil, fmt.Errorf("failed to get income: %w", err)
	}

	return &income, nil
}

func (s *incomeServiceImpl) Update(ctx context.Context, id int, req models.UpdateIncomeRequest) error {
	// Get current amount for balance adjustment
	var currentAmount float64
	var currentCashBoxID int
	err := s.db.QueryRow(ctx, `SELECT amount, cash_box_id FROM incomes WHERE id = $1`, id).Scan(&currentAmount, &currentCashBoxID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return ErrIncomeNotFound
		}
		return fmt.Errorf("failed to get current income: %w", err)
	}

	query := `
		UPDATE incomes SET
			income_type = COALESCE($2, income_type),
			type_id = COALESCE($3, type_id),
			date = COALESCE($4, date),
			receipt_date = COALESCE($5, receipt_date),
			amount = COALESCE($6, amount),
			note = COALESCE($7, note),
			cash_box_id = COALESCE($8, cash_box_id)
		WHERE id = $1`

	result, err := s.db.Exec(ctx, query,
		id,
		req.IncomeType,
		req.TypeID,
		req.Date,
		req.ReceiptDate,
		req.Amount,
		req.Note,
		req.CashBoxID,
	)

	if err != nil {
		return fmt.Errorf("failed to update income: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrIncomeNotFound
	}

	// Adjust cash box balance if amount changed
	if req.Amount != nil && *req.Amount != currentAmount {
		diff := *req.Amount - currentAmount
		_, err = s.db.Exec(ctx, `
			UPDATE cash_boxes SET current_balance = current_balance + $1 WHERE id = $2`,
			diff, currentCashBoxID)
		if err != nil {
			return fmt.Errorf("failed to update cash box: %w", err)
		}
	}

	return nil
}

func (s *incomeServiceImpl) Delete(ctx context.Context, id int) error {
	// Get amount for balance adjustment
	var amount float64
	var cashBoxID int
	err := s.db.QueryRow(ctx, `SELECT amount, cash_box_id FROM incomes WHERE id = $1`, id).Scan(&amount, &cashBoxID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return ErrIncomeNotFound
		}
		return fmt.Errorf("failed to get income: %w", err)
	}

	result, err := s.db.Exec(ctx, `DELETE FROM incomes WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete income: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrIncomeNotFound
	}

	// Subtract from cash box balance
	_, err = s.db.Exec(ctx, `
		UPDATE cash_boxes SET current_balance = current_balance - $1 WHERE id = $2`,
		amount, cashBoxID)
	if err != nil {
		return fmt.Errorf("failed to update cash box: %w", err)
	}

	return nil
}

func (s *incomeServiceImpl) List(ctx context.Context, filters models.IncomeListFilters) ([]models.Income, int64, error) {
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PageSize < 1 {
		filters.PageSize = 20
	}

	// Count query
	countQuery := `SELECT COUNT(*) FROM incomes WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if filters.TypeID != nil {
		countQuery += fmt.Sprintf(" AND type_id = $%d", argIndex)
		args = append(args, *filters.TypeID)
		argIndex++
	}
	if filters.CashBoxID != nil {
		countQuery += fmt.Sprintf(" AND cash_box_id = $%d", argIndex)
		args = append(args, *filters.CashBoxID)
		argIndex++
	}
	if filters.CustomerID != nil {
		countQuery += fmt.Sprintf(" AND customer_id = $%d", argIndex)
		args = append(args, *filters.CustomerID)
		argIndex++
	}
	if filters.DateFrom != nil {
		countQuery += fmt.Sprintf(" AND date >= $%d", argIndex)
		args = append(args, time.Time(*filters.DateFrom))
		argIndex++
	}
	if filters.DateTo != nil {
		countQuery += fmt.Sprintf(" AND date <= $%d", argIndex)
		args = append(args, time.Time(*filters.DateTo))
		argIndex++
	}

	var total int64
	err := s.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count: %w", err)
	}

	// Data query
	query := `
		SELECT id, income_type, type_id, date, receipt_date, amount,
			note, cash_box_id, customer_id, invoice_id, created_by, created_at
		FROM incomes WHERE 1=1`

	args = []interface{}{}
	argIndex = 1

	if filters.TypeID != nil {
		query += fmt.Sprintf(" AND type_id = $%d", argIndex)
		args = append(args, *filters.TypeID)
		argIndex++
	}
	if filters.CashBoxID != nil {
		query += fmt.Sprintf(" AND cash_box_id = $%d", argIndex)
		args = append(args, *filters.CashBoxID)
		argIndex++
	}
	if filters.CustomerID != nil {
		query += fmt.Sprintf(" AND customer_id = $%d", argIndex)
		args = append(args, *filters.CustomerID)
		argIndex++
	}
	if filters.DateFrom != nil {
		query += fmt.Sprintf(" AND date >= $%d", argIndex)
		args = append(args, time.Time(*filters.DateFrom))
		argIndex++
	}
	if filters.DateTo != nil {
		query += fmt.Sprintf(" AND date <= $%d", argIndex)
		args = append(args, time.Time(*filters.DateTo))
		argIndex++
	}

	query += " ORDER BY date DESC, id DESC"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, filters.PageSize, (filters.Page-1)*filters.PageSize)

	rows := s.db.Query(ctx, query, args...)
	defer rows.Close()

	var results []models.Income
	for rows.Next() {
		var income models.Income
		err := rows.Scan(
			&income.ID,
			&income.IncomeType,
			&income.TypeID,
			&income.Date,
			&income.ReceiptDate,
			&income.Amount,
			&income.Note,
			&income.CashBoxID,
			&income.CustomerID,
			&income.InvoiceID,
			&income.CreatedBy,
			&income.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan: %w", err)
		}
		results = append(results, income)
	}

	return results, total, nil
}
