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
	ErrExpenseNotFound = errors.New("expense not found")
)

type ExpenseService interface {
	Create(ctx context.Context, req models.CreateExpenseRequest, createdBy int) (int, error)
	GetByID(ctx context.Context, id int) (*models.Expense, error)
	Update(ctx context.Context, id int, req models.UpdateExpenseRequest) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, filters models.ExpenseListFilters) ([]models.Expense, int64, error)
}

type expenseServiceImpl struct {
	db postgres.Connection
}

func NewExpenseService(db postgres.Connection) ExpenseService {
	return &expenseServiceImpl{db: db}
}

func (s *expenseServiceImpl) Create(ctx context.Context, req models.CreateExpenseRequest, createdBy int) (int, error) {
	query := `
		INSERT INTO expenses (
			type_id, date, expense_date, amount,
			note, cash_box_id, vendor_id, bill_id, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	var id int
	err := s.db.QueryRow(ctx, query,
		req.TypeID,
		time.Time(req.Date),
		time.Time(req.ExpenseDate),
		req.Amount,
		req.Note,
		req.CashBoxID,
		req.VendorID,
		req.BillID,
		createdBy,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create expense: %w", err)
	}

	// Update cash box balance (subtract)
	_, err = s.db.Exec(ctx, `
		UPDATE cash_boxes SET current_balance = current_balance - $1 WHERE id = $2`,
		req.Amount, req.CashBoxID)
	if err != nil {
		return 0, fmt.Errorf("failed to update cash box: %w", err)
	}

	return id, nil
}

func (s *expenseServiceImpl) GetByID(ctx context.Context, id int) (*models.Expense, error) {
	query := `
		SELECT id, type_id, date, expense_date, amount,
			note, cash_box_id, vendor_id, bill_id, is_marked, created_by, created_at
		FROM expenses WHERE id = $1`

	var expense models.Expense
	err := s.db.QueryRow(ctx, query, id).Scan(
		&expense.ID,
		&expense.TypeID,
		&expense.Date,
		&expense.ExpenseDate,
		&expense.Amount,
		&expense.Note,
		&expense.CashBoxID,
		&expense.VendorID,
		&expense.BillID,
		&expense.IsMarked,
		&expense.CreatedBy,
		&expense.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrExpenseNotFound
		}
		return nil, fmt.Errorf("failed to get expense: %w", err)
	}

	return &expense, nil
}

func (s *expenseServiceImpl) Update(ctx context.Context, id int, req models.UpdateExpenseRequest) error {
	// Get current amount for balance adjustment
	var currentAmount float64
	var currentCashBoxID int
	err := s.db.QueryRow(ctx, `SELECT amount, cash_box_id FROM expenses WHERE id = $1`, id).Scan(&currentAmount, &currentCashBoxID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return ErrExpenseNotFound
		}
		return fmt.Errorf("failed to get current expense: %w", err)
	}

	query := `
		UPDATE expenses SET
			type_id = COALESCE($2, type_id),
			date = COALESCE($3, date),
			expense_date = COALESCE($4, expense_date),
			amount = COALESCE($5, amount),
			note = COALESCE($6, note),
			cash_box_id = COALESCE($7, cash_box_id),
			is_marked = COALESCE($8, is_marked)
		WHERE id = $1`

	result, err := s.db.Exec(ctx, query,
		id,
		req.TypeID,
		req.Date,
		req.ExpenseDate,
		req.Amount,
		req.Note,
		req.CashBoxID,
		req.IsMarked,
	)

	if err != nil {
		return fmt.Errorf("failed to update expense: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrExpenseNotFound
	}

	// Adjust cash box balance if amount changed
	if req.Amount != nil && *req.Amount != currentAmount {
		diff := currentAmount - *req.Amount // Reverse the difference for expenses
		_, err = s.db.Exec(ctx, `
			UPDATE cash_boxes SET current_balance = current_balance + $1 WHERE id = $2`,
			diff, currentCashBoxID)
		if err != nil {
			return fmt.Errorf("failed to update cash box: %w", err)
		}
	}

	return nil
}

func (s *expenseServiceImpl) Delete(ctx context.Context, id int) error {
	// Get amount for balance adjustment
	var amount float64
	var cashBoxID int
	err := s.db.QueryRow(ctx, `SELECT amount, cash_box_id FROM expenses WHERE id = $1`, id).Scan(&amount, &cashBoxID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return ErrExpenseNotFound
		}
		return fmt.Errorf("failed to get expense: %w", err)
	}

	result, err := s.db.Exec(ctx, `DELETE FROM expenses WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete expense: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrExpenseNotFound
	}

	// Add back to cash box balance
	_, err = s.db.Exec(ctx, `
		UPDATE cash_boxes SET current_balance = current_balance + $1 WHERE id = $2`,
		amount, cashBoxID)
	if err != nil {
		return fmt.Errorf("failed to update cash box: %w", err)
	}

	return nil
}

func (s *expenseServiceImpl) List(ctx context.Context, filters models.ExpenseListFilters) ([]models.Expense, int64, error) {
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PageSize < 1 {
		filters.PageSize = 20
	}

	// Count query
	countQuery := `SELECT COUNT(*) FROM expenses WHERE 1=1`
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
	if filters.VendorID != nil {
		countQuery += fmt.Sprintf(" AND vendor_id = $%d", argIndex)
		args = append(args, *filters.VendorID)
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
	if filters.IsMarked != nil {
		countQuery += fmt.Sprintf(" AND is_marked = $%d", argIndex)
		args = append(args, *filters.IsMarked)
		argIndex++
	}

	var total int64
	err := s.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count: %w", err)
	}

	// Data query
	query := `
		SELECT id, type_id, date, expense_date, amount,
			note, cash_box_id, vendor_id, bill_id, is_marked, created_by, created_at
		FROM expenses WHERE 1=1`

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
	if filters.VendorID != nil {
		query += fmt.Sprintf(" AND vendor_id = $%d", argIndex)
		args = append(args, *filters.VendorID)
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
	if filters.IsMarked != nil {
		query += fmt.Sprintf(" AND is_marked = $%d", argIndex)
		args = append(args, *filters.IsMarked)
		argIndex++
	}

	query += " ORDER BY date DESC, id DESC"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, filters.PageSize, (filters.Page-1)*filters.PageSize)

	rows := s.db.Query(ctx, query, args...)
	defer rows.Close()

	var results []models.Expense
	for rows.Next() {
		var expense models.Expense
		err := rows.Scan(
			&expense.ID,
			&expense.TypeID,
			&expense.Date,
			&expense.ExpenseDate,
			&expense.Amount,
			&expense.Note,
			&expense.CashBoxID,
			&expense.VendorID,
			&expense.BillID,
			&expense.IsMarked,
			&expense.CreatedBy,
			&expense.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan: %w", err)
		}
		results = append(results, expense)
	}

	return results, total, nil
}
