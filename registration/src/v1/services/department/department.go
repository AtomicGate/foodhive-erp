package department

import (
	"context"
	"errors"
	"fmt"

	"github.com/anas-dev-92/FoodHive/core/postgres"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/models"
	"github.com/jackc/pgx/v5"
)

var (
	ErrNotFound = errors.New("department not found")
)

type DepartmentService interface {
	Create(ctx context.Context, req models.CreateDepartmentRequest) (int, error)
	GetByID(ctx context.Context, id int) (*models.DepartmentWithManager, error)
	Update(ctx context.Context, id int, req models.UpdateDepartmentRequest) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context) ([]models.DepartmentWithManager, error)
}

type departmentServiceImpl struct {
	db postgres.Connection
}

func New(db postgres.Connection) DepartmentService {
	return &departmentServiceImpl{db: db}
}

func (s *departmentServiceImpl) Create(ctx context.Context, req models.CreateDepartmentRequest) (int, error) {
	query := `
		INSERT INTO departments (name, code, description, manager_id, parent_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	var id int
	err := s.db.QueryRow(ctx, query,
		req.Name,
		req.Code,
		req.Description,
		req.ManagerID,
		req.ParentID,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create department: %w", err)
	}

	return id, nil
}

func (s *departmentServiceImpl) GetByID(ctx context.Context, id int) (*models.DepartmentWithManager, error) {
	query := `
		SELECT d.id, d.name, d.code, d.description, d.manager_id, d.parent_id, d.is_active,
			d.created_at, d.updated_at, COALESCE(e.english_name, '') as manager_name
		FROM departments d
		LEFT JOIN employees e ON d.manager_id = e.id
		WHERE d.id = $1`

	var dept models.DepartmentWithManager
	err := s.db.QueryRow(ctx, query, id).Scan(
		&dept.ID,
		&dept.Name,
		&dept.Code,
		&dept.Description,
		&dept.ManagerID,
		&dept.ParentID,
		&dept.IsActive,
		&dept.CreatedAt,
		&dept.UpdatedAt,
		&dept.ManagerName,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get department: %w", err)
	}

	return &dept, nil
}

func (s *departmentServiceImpl) Update(ctx context.Context, id int, req models.UpdateDepartmentRequest) error {
	query := `
		UPDATE departments SET
			name = COALESCE($2, name),
			code = COALESCE($3, code),
			description = COALESCE($4, description),
			manager_id = COALESCE($5, manager_id),
			parent_id = COALESCE($6, parent_id),
			is_active = COALESCE($7, is_active),
			updated_at = NOW()
		WHERE id = $1`

	result, err := s.db.Exec(ctx, query,
		id,
		req.Name,
		req.Code,
		req.Description,
		req.ManagerID,
		req.ParentID,
		req.IsActive,
	)
	if err != nil {
		return fmt.Errorf("failed to update department: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *departmentServiceImpl) Delete(ctx context.Context, id int) error {
	// Soft delete
	query := `UPDATE departments SET is_active = false, updated_at = NOW() WHERE id = $1`
	result, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete department: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *departmentServiceImpl) List(ctx context.Context) ([]models.DepartmentWithManager, error) {
	query := `
		SELECT d.id, d.name, d.code, d.description, d.manager_id, d.parent_id, d.is_active,
			d.created_at, d.updated_at, COALESCE(e.english_name, '') as manager_name
		FROM departments d
		LEFT JOIN employees e ON d.manager_id = e.id
		WHERE d.is_active = true
		ORDER BY d.name`

	rows := s.db.Query(ctx, query)
	defer rows.Close()

	var results []models.DepartmentWithManager
	for rows.Next() {
		var dept models.DepartmentWithManager
		err := rows.Scan(
			&dept.ID,
			&dept.Name,
			&dept.Code,
			&dept.Description,
			&dept.ManagerID,
			&dept.ParentID,
			&dept.IsActive,
			&dept.CreatedAt,
			&dept.UpdatedAt,
			&dept.ManagerName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}
		results = append(results, dept)
	}

	return results, nil
}
