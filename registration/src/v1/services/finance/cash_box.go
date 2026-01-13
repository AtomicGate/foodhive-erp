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
	ErrCashBoxNotFound = errors.New("cash box not found")
)

type CashBoxService interface {
	Create(ctx context.Context, req models.CreateCashBoxRequest) (int, error)
	GetByID(ctx context.Context, id int) (*models.CashBox, error)
	Update(ctx context.Context, id int, req models.UpdateCashBoxRequest) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context) ([]models.CashBox, error)
}

type cashBoxServiceImpl struct {
	db postgres.Connection
}

func NewCashBoxService(db postgres.Connection) CashBoxService {
	return &cashBoxServiceImpl{db: db}
}

func (s *cashBoxServiceImpl) Create(ctx context.Context, req models.CreateCashBoxRequest) (int, error) {
	if req.Currency == "" {
		req.Currency = "USD"
	}

	query := `
		INSERT INTO cash_boxes (name, code, currency, warehouse_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	var id int
	err := s.db.QueryRow(ctx, query,
		req.Name,
		req.Code,
		req.Currency,
		req.WarehouseID,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create cash box: %w", err)
	}

	return id, nil
}

func (s *cashBoxServiceImpl) GetByID(ctx context.Context, id int) (*models.CashBox, error) {
	query := `
		SELECT id, name, code, currency, current_balance, is_active, warehouse_id, created_at, updated_at
		FROM cash_boxes WHERE id = $1`

	var cashBox models.CashBox
	err := s.db.QueryRow(ctx, query, id).Scan(
		&cashBox.ID,
		&cashBox.Name,
		&cashBox.Code,
		&cashBox.Currency,
		&cashBox.CurrentBalance,
		&cashBox.IsActive,
		&cashBox.WarehouseID,
		&cashBox.CreatedAt,
		&cashBox.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrCashBoxNotFound
		}
		return nil, fmt.Errorf("failed to get cash box: %w", err)
	}

	return &cashBox, nil
}

func (s *cashBoxServiceImpl) Update(ctx context.Context, id int, req models.UpdateCashBoxRequest) error {
	query := `
		UPDATE cash_boxes SET
			name = COALESCE($2, name),
			currency = COALESCE($3, currency),
			is_active = COALESCE($4, is_active),
			warehouse_id = COALESCE($5, warehouse_id),
			updated_at = NOW()
		WHERE id = $1`

	result, err := s.db.Exec(ctx, query,
		id,
		req.Name,
		req.Currency,
		req.IsActive,
		req.WarehouseID,
	)

	if err != nil {
		return fmt.Errorf("failed to update cash box: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrCashBoxNotFound
	}

	return nil
}

func (s *cashBoxServiceImpl) Delete(ctx context.Context, id int) error {
	// Soft delete - mark as inactive
	query := `UPDATE cash_boxes SET is_active = false, updated_at = NOW() WHERE id = $1`
	result, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete cash box: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrCashBoxNotFound
	}

	return nil
}

func (s *cashBoxServiceImpl) List(ctx context.Context) ([]models.CashBox, error) {
	query := `
		SELECT id, name, code, currency, current_balance, is_active, warehouse_id, created_at, updated_at
		FROM cash_boxes
		ORDER BY name`

	rows := s.db.Query(ctx, query)
	defer rows.Close()

	var results []models.CashBox
	for rows.Next() {
		var cashBox models.CashBox
		err := rows.Scan(
			&cashBox.ID,
			&cashBox.Name,
			&cashBox.Code,
			&cashBox.Currency,
			&cashBox.CurrentBalance,
			&cashBox.IsActive,
			&cashBox.WarehouseID,
			&cashBox.CreatedAt,
			&cashBox.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}
		results = append(results, cashBox)
	}

	return results, nil
}
