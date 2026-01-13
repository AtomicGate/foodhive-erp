package role

import (
	"context"
	"errors"
	"fmt"

	"github.com/anas-dev-92/FoodHive/core/postgres"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/models"
	"github.com/jackc/pgx/v5"
)

var (
	ErrNotFound = errors.New("role not found")
)

type RoleService interface {
	Create(ctx context.Context, req models.CreateRoleRequest) (int, error)
	GetByID(ctx context.Context, id int) (*models.Role, error)
	Update(ctx context.Context, id int, req models.UpdateRoleRequest) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context) ([]models.Role, error)
}

type roleServiceImpl struct {
	db postgres.Connection
}

func New(db postgres.Connection) RoleService {
	return &roleServiceImpl{db: db}
}

func (s *roleServiceImpl) Create(ctx context.Context, req models.CreateRoleRequest) (int, error) {
	query := `INSERT INTO roles (role_name, role_desc) VALUES ($1, $2) RETURNING id`

	var id int
	err := s.db.QueryRow(ctx, query, req.RoleName, req.Description).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create role: %w", err)
	}

	return id, nil
}

func (s *roleServiceImpl) GetByID(ctx context.Context, id int) (*models.Role, error) {
	query := `SELECT id, role_name, role_desc, is_active FROM roles WHERE id = $1`

	var role models.Role
	err := s.db.QueryRow(ctx, query, id).Scan(
		&role.ID,
		&role.RoleName,
		&role.Description,
		&role.IsActive,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	return &role, nil
}

func (s *roleServiceImpl) Update(ctx context.Context, id int, req models.UpdateRoleRequest) error {
	query := `
		UPDATE roles SET
			role_name = COALESCE($2, role_name),
			description = COALESCE($3, description),
			is_active = COALESCE($4, is_active)
		WHERE id = $1`

	result, err := s.db.Exec(ctx, query, id, req.RoleName, req.Description, req.IsActive)
	if err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *roleServiceImpl) Delete(ctx context.Context, id int) error {
	// Soft delete
	query := `UPDATE roles SET is_active = false WHERE id = $1`
	result, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *roleServiceImpl) List(ctx context.Context) ([]models.Role, error) {
	query := `SELECT id, role_name, role_desc, is_active FROM roles WHERE is_active = true ORDER BY role_name`

	rows := s.db.Query(ctx, query)
	defer rows.Close()

	var results []models.Role
	for rows.Next() {
		var role models.Role
		err := rows.Scan(
			&role.ID,
			&role.RoleName,
			&role.Description,
			&role.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}
		results = append(results, role)
	}

	return results, nil
}

// ============================================
// Permission Service
// ============================================

type PermissionService interface {
	AssignPermission(ctx context.Context, req models.AssignPermissionRequest) error
	BulkAssignPermissions(ctx context.Context, req models.BulkAssignPermissionsRequest) error
	GetEmployeePermissions(ctx context.Context, employeeID int) ([]models.EmployeePermissionWithPage, error)
	RemovePermission(ctx context.Context, employeeID int, pageID int) error
}

type permissionServiceImpl struct {
	db postgres.Connection
}

func NewPermissionService(db postgres.Connection) PermissionService {
	return &permissionServiceImpl{db: db}
}

func (s *permissionServiceImpl) AssignPermission(ctx context.Context, req models.AssignPermissionRequest) error {
	query := `
		INSERT INTO employee_permissions (employee_id, page_id, can_create, can_view, can_update, can_delete)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (employee_id, page_id) DO UPDATE SET
			can_create = EXCLUDED.can_create,
			can_view = EXCLUDED.can_view,
			can_update = EXCLUDED.can_update,
			can_delete = EXCLUDED.can_delete`

	_, err := s.db.Exec(ctx, query,
		req.EmployeeID,
		req.PageID,
		req.CanCreate,
		req.CanView,
		req.CanUpdate,
		req.CanDelete,
	)
	if err != nil {
		return fmt.Errorf("failed to assign permission: %w", err)
	}

	return nil
}

func (s *permissionServiceImpl) BulkAssignPermissions(ctx context.Context, req models.BulkAssignPermissionsRequest) error {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Clear existing permissions
	_, err = tx.Exec(ctx, `DELETE FROM employee_permissions WHERE employee_id = $1`, req.EmployeeID)
	if err != nil {
		return fmt.Errorf("failed to clear permissions: %w", err)
	}

	// Insert new permissions
	for _, perm := range req.Permissions {
		query := `
			INSERT INTO employee_permissions (employee_id, page_id, can_create, can_view, can_update, can_delete)
			VALUES ($1, $2, $3, $4, $5, $6)`
		_, err = tx.Exec(ctx, query,
			req.EmployeeID,
			perm.PageID,
			perm.CanCreate,
			perm.CanView,
			perm.CanUpdate,
			perm.CanDelete,
		)
		if err != nil {
			return fmt.Errorf("failed to insert permission: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}

func (s *permissionServiceImpl) GetEmployeePermissions(ctx context.Context, employeeID int) ([]models.EmployeePermissionWithPage, error) {
	query := `
		SELECT 
			ep.id, ep.employee_id, ep.page_id, ep.can_create, ep.can_view, ep.can_update, ep.can_delete,
			p.id, p.route_name, p.display_name, p.description, p.module, p.is_active
		FROM employee_permissions ep
		JOIN pages p ON ep.page_id = p.id
		WHERE ep.employee_id = $1
		ORDER BY p.module, p.display_name`

	rows := s.db.Query(ctx, query, employeeID)
	defer rows.Close()

	var results []models.EmployeePermissionWithPage
	for rows.Next() {
		var ep models.EmployeePermissionWithPage
		err := rows.Scan(
			&ep.ID, &ep.EmployeeID, &ep.PageID, &ep.CanCreate, &ep.CanView, &ep.CanUpdate, &ep.CanDelete,
			&ep.Page.ID, &ep.Page.RouteName, &ep.Page.DisplayName, &ep.Page.Description, &ep.Page.Module, &ep.Page.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}
		results = append(results, ep)
	}

	return results, nil
}

func (s *permissionServiceImpl) RemovePermission(ctx context.Context, employeeID int, pageID int) error {
	query := `DELETE FROM employee_permissions WHERE employee_id = $1 AND page_id = $2`
	result, err := s.db.Exec(ctx, query, employeeID, pageID)
	if err != nil {
		return fmt.Errorf("failed to remove permission: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("permission not found")
	}

	return nil
}

// ============================================
// Page Service (System Modules Management)
// ============================================

type PageService interface {
	CreatePage(ctx context.Context, req models.CreatePageRequest) (int, error)
	GetPage(ctx context.Context, id int) (*models.Page, error)
	UpdatePage(ctx context.Context, id int, req models.UpdatePageRequest) error
	DeletePage(ctx context.Context, id int) error
	ListPages(ctx context.Context, module string) ([]models.Page, error)
	ListModules(ctx context.Context) ([]string, error)
}

type pageServiceImpl struct {
	db postgres.Connection
}

func NewPageService(db postgres.Connection) PageService {
	return &pageServiceImpl{db: db}
}

func (s *pageServiceImpl) CreatePage(ctx context.Context, req models.CreatePageRequest) (int, error) {
	query := `
		INSERT INTO pages (route_name, display_name, description, module, icon, parent_id, sort_order, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, true)
		RETURNING id`

	var id int
	err := s.db.QueryRow(ctx, query,
		req.RouteName, req.DisplayName, req.Description, req.Module,
		req.Icon, req.ParentID, req.SortOrder,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create page: %w", err)
	}

	return id, nil
}

func (s *pageServiceImpl) GetPage(ctx context.Context, id int) (*models.Page, error) {
	query := `SELECT id, route_name, display_name, description, module, is_active FROM pages WHERE id = $1`

	var page models.Page
	err := s.db.QueryRow(ctx, query, id).Scan(
		&page.ID, &page.RouteName, &page.DisplayName, &page.Description, &page.Module, &page.IsActive,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("page not found")
		}
		return nil, fmt.Errorf("failed to get page: %w", err)
	}

	return &page, nil
}

func (s *pageServiceImpl) UpdatePage(ctx context.Context, id int, req models.UpdatePageRequest) error {
	query := `
		UPDATE pages SET
			route_name = COALESCE($2, route_name),
			display_name = COALESCE($3, display_name),
			description = COALESCE($4, description),
			module = COALESCE($5, module),
			is_active = COALESCE($6, is_active)
		WHERE id = $1`

	result, err := s.db.Exec(ctx, query, id, req.RouteName, req.DisplayName, req.Description, req.Module, req.IsActive)
	if err != nil {
		return fmt.Errorf("failed to update page: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("page not found")
	}

	return nil
}

func (s *pageServiceImpl) DeletePage(ctx context.Context, id int) error {
	result, err := s.db.Exec(ctx, `UPDATE pages SET is_active = false WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete page: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("page not found")
	}

	return nil
}

func (s *pageServiceImpl) ListPages(ctx context.Context, module string) ([]models.Page, error) {
	query := `SELECT id, route_name, display_name, description, module, is_active FROM pages WHERE is_active = true`
	args := []interface{}{}

	if module != "" {
		query += ` AND module = $1`
		args = append(args, module)
	}

	query += ` ORDER BY module, sort_order, display_name`

	rows := s.db.Query(ctx, query, args...)
	defer rows.Close()

	var results []models.Page
	for rows.Next() {
		var page models.Page
		err := rows.Scan(&page.ID, &page.RouteName, &page.DisplayName, &page.Description, &page.Module, &page.IsActive)
		if err != nil {
			return nil, fmt.Errorf("failed to scan page: %w", err)
		}
		results = append(results, page)
	}

	return results, nil
}

func (s *pageServiceImpl) ListModules(ctx context.Context) ([]string, error) {
	query := `SELECT DISTINCT module FROM pages WHERE is_active = true AND module IS NOT NULL ORDER BY module`

	rows := s.db.Query(ctx, query)
	defer rows.Close()

	var modules []string
	for rows.Next() {
		var module string
		if err := rows.Scan(&module); err != nil {
			return nil, fmt.Errorf("failed to scan module: %w", err)
		}
		modules = append(modules, module)
	}

	return modules, nil
}

// ============================================
// Role Template Service
// ============================================

type RoleTemplateService interface {
	SetRolePermissions(ctx context.Context, req models.AssignRoleTemplateRequest) error
	GetRolePermissions(ctx context.Context, roleID int) (*models.RoleWithPermissions, error)
	ApplyRoleToEmployee(ctx context.Context, employeeID int, roleID int) error
	GetUserAccessSummary(ctx context.Context, employeeID int) (*models.UserAccessSummary, error)
}

type roleTemplateServiceImpl struct {
	db postgres.Connection
}

func NewRoleTemplateService(db postgres.Connection) RoleTemplateService {
	return &roleTemplateServiceImpl{db: db}
}

func (s *roleTemplateServiceImpl) SetRolePermissions(ctx context.Context, req models.AssignRoleTemplateRequest) error {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Clear existing role permissions
	_, err = tx.Exec(ctx, `DELETE FROM role_templates WHERE role_id = $1`, req.RoleID)
	if err != nil {
		return fmt.Errorf("failed to clear role permissions: %w", err)
	}

	// Insert new permissions
	for _, perm := range req.Permissions {
		query := `
			INSERT INTO role_templates (role_id, page_id, can_create, can_view, can_update, can_delete)
			VALUES ($1, $2, $3, $4, $5, $6)`
		_, err = tx.Exec(ctx, query, req.RoleID, perm.PageID, perm.CanCreate, perm.CanView, perm.CanUpdate, perm.CanDelete)
		if err != nil {
			return fmt.Errorf("failed to insert role permission: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}

func (s *roleTemplateServiceImpl) GetRolePermissions(ctx context.Context, roleID int) (*models.RoleWithPermissions, error) {
	// Get role
	var role models.Role
	err := s.db.QueryRow(ctx, `SELECT id, role_name, role_desc, is_active FROM roles WHERE id = $1`, roleID).Scan(
		&role.ID, &role.RoleName, &role.Description, &role.IsActive,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("role not found")
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	// Get permissions
	query := `
		SELECT rt.id, rt.role_id, rt.page_id, rt.can_create, rt.can_view, rt.can_update, rt.can_delete
		FROM role_templates rt
		WHERE rt.role_id = $1`

	rows := s.db.Query(ctx, query, roleID)
	defer rows.Close()

	var permissions []models.RoleTemplate
	for rows.Next() {
		var perm models.RoleTemplate
		err := rows.Scan(&perm.ID, &perm.RoleID, &perm.PageID, &perm.CanCreate, &perm.CanView, &perm.CanUpdate, &perm.CanDelete)
		if err != nil {
			return nil, fmt.Errorf("failed to scan permission: %w", err)
		}
		permissions = append(permissions, perm)
	}

	return &models.RoleWithPermissions{
		Role:        role,
		Permissions: permissions,
	}, nil
}

func (s *roleTemplateServiceImpl) ApplyRoleToEmployee(ctx context.Context, employeeID int, roleID int) error {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Update employee's role
	_, err = tx.Exec(ctx, `UPDATE employees SET role_id = $1 WHERE id = $2`, roleID, employeeID)
	if err != nil {
		return fmt.Errorf("failed to update employee role: %w", err)
	}

	// Clear existing employee permissions
	_, err = tx.Exec(ctx, `DELETE FROM employee_permissions WHERE employee_id = $1`, employeeID)
	if err != nil {
		return fmt.Errorf("failed to clear employee permissions: %w", err)
	}

	// Copy role template permissions to employee
	query := `
		INSERT INTO employee_permissions (employee_id, page_id, can_create, can_view, can_update, can_delete)
		SELECT $1, page_id, can_create, can_view, can_update, can_delete
		FROM role_templates WHERE role_id = $2`
	_, err = tx.Exec(ctx, query, employeeID, roleID)
	if err != nil {
		return fmt.Errorf("failed to copy permissions: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}

func (s *roleTemplateServiceImpl) GetUserAccessSummary(ctx context.Context, employeeID int) (*models.UserAccessSummary, error) {
	// Get employee info
	var summary models.UserAccessSummary
	err := s.db.QueryRow(ctx, `
		SELECT e.id, COALESCE(e.english_name, e.email), e.email, COALESCE(e.role_id, 0), 
		       COALESCE(r.role_name, ''), COALESCE(e.security_level, 1)
		FROM employees e
		LEFT JOIN roles r ON e.role_id = r.id
		WHERE e.id = $1`, employeeID).Scan(
		&summary.EmployeeID, &summary.EmployeeName, &summary.Email,
		&summary.RoleID, &summary.RoleName, &summary.SecurityLevel,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("employee not found")
		}
		return nil, fmt.Errorf("failed to get employee: %w", err)
	}

	// Get permissions with page details
	query := `
		SELECT p.id, p.route_name, p.display_name, p.description, p.module, p.is_active,
		       ep.can_create, ep.can_view, ep.can_update, ep.can_delete
		FROM employee_permissions ep
		JOIN pages p ON ep.page_id = p.id
		WHERE ep.employee_id = $1 AND p.is_active = true
		ORDER BY p.module, p.display_name`

	rows := s.db.Query(ctx, query, employeeID)
	defer rows.Close()

	summary.Permissions = []models.PageWithPermissions{}
	summary.ModuleAccess = make(map[string]bool)

	for rows.Next() {
		var perm models.PageWithPermissions
		err := rows.Scan(
			&perm.Page.ID, &perm.Page.RouteName, &perm.Page.DisplayName,
			&perm.Page.Description, &perm.Page.Module, &perm.Page.IsActive,
			&perm.CanCreate, &perm.CanView, &perm.CanUpdate, &perm.CanDelete,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan permission: %w", err)
		}
		summary.Permissions = append(summary.Permissions, perm)

		// Track module access
		if perm.CanView || perm.CanCreate || perm.CanUpdate || perm.CanDelete {
			summary.ModuleAccess[perm.Page.Module] = true
		}
	}

	return &summary, nil
}

// ============================================
// Security Level Helpers
// ============================================

func GetSecurityLevels() []models.SecurityLevelInfo {
	return []models.SecurityLevelInfo{
		{Level: models.SecurityLevelViewer, Name: "Viewer", Description: "View-only access"},
		{Level: models.SecurityLevelOperator, Name: "Operator", Description: "Basic operations"},
		{Level: models.SecurityLevelUser, Name: "User", Description: "Standard user access"},
		{Level: models.SecurityLevelSuperUser, Name: "Super User", Description: "Advanced user access"},
		{Level: models.SecurityLevelManager, Name: "Manager", Description: "Department manager access"},
		{Level: models.SecurityLevelAdmin, Name: "Admin", Description: "Administrator access"},
		{Level: models.SecurityLevelSuperAdmin, Name: "Super Admin", Description: "Full system access"},
	}
}
