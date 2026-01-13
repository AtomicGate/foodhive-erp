package auth

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/anas-dev-92/FoodHive/core/postgres"
)

type AuthService interface {
	HasUserPermission(ctx context.Context, email, path string) (bool, error)
	GetUserRole(ctx context.Context, email string) (string, error)
}

type AuthServiceImpl struct {
	db postgres.Executor
}

func New(db postgres.Executor) AuthService {
	return &AuthServiceImpl{db: db}
}
func (s *AuthServiceImpl) GetUserPagesAndPermissions(ctx context.Context, userID int) ([]map[string]interface{}, error) {
	query := `
		SELECT p.page_name, p.route_name, ep.can_create, ep.can_update, ep.can_delete, ep.can_view
		FROM emp_page ep
		JOIN pages p ON ep.page_id = p.id
		WHERE ep.user_id = $1
	`
	rows := s.db.Query(ctx, query, userID)
	defer rows.Close()

	var pages []map[string]interface{}
	for rows.Next() {
		var pageName, routeName string
		var canCreate, canUpdate, canDelete, canView bool
		err := rows.Scan(&pageName, &routeName, &canCreate, &canUpdate, &canDelete, &canView)
		if err != nil {
			return nil, err
		}

		page := map[string]interface{}{
			"page_name":  pageName,
			"route_name": routeName,
			"can_create": canCreate,
			"can_update": canUpdate,
			"can_delete": canDelete,
			"can_view":   canView,
		}
		pages = append(pages, page)
	}

	return pages, nil
}
func (s *AuthServiceImpl) HasUserPermission(ctx context.Context, email, path string) (bool, error) {
	role, err := s.GetUserRole(ctx, email)
	if err != nil {
		return false, err
	}
	if role == "admin" {
		return true, nil
	}

	var userID int
	query := `SELECT id FROM employees WHERE email = $1`
	err = s.db.QueryRow(ctx, query, email).Scan(&userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	// Split the path into parts
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		return false, nil // Invalid path format
	}

	// Extract the action (last part of the path)
	action := parts[len(parts)-1]

	// Reconstruct the base path (everything except the action)
	basePath := strings.Join(parts[:len(parts)-1], "/")

	var permissionType string
	switch action {
	case "create":
		permissionType = "can_create"
	case "update":
		permissionType = "can_update"
	case "delete":
		permissionType = "can_delete"
	default:
		permissionType = "can_view"
	}

	query = `SELECT ` + permissionType + ` FROM emp_page WHERE user_id = $1 AND page_id = (SELECT id FROM pages WHERE route_name = $2)`
	var hasPermission bool
	err = s.db.QueryRow(ctx, query, userID, basePath).Scan(&hasPermission)

	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return hasPermission, nil
}

func (s *AuthServiceImpl) GetUserRole(ctx context.Context, email string) (string, error) {
	var role string
	query := `
        SELECT r.role_desc 
        FROM employees u
        JOIN roles r ON u.role_id = r.id
        WHERE u.email = $1`
	err := s.db.QueryRow(ctx, query, email).Scan(&role)
	if err != nil {
		return "", err
	}

	return role, nil
}
