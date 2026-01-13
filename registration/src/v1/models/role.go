package models

// ============================================
// Role Models
// ============================================

type Role struct {
	ID          int    `json:"id"`
	RoleName    string `json:"role_name"`
	Description string `json:"description,omitempty"`
	IsActive    bool   `json:"is_active"`
}

// ============================================
// Permission Models (Route-based access control)
// ============================================

type Page struct {
	ID          int    `json:"id"`
	RouteName   string `json:"route_name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description,omitempty"`
	Module      string `json:"module,omitempty"` // e.g., "sales", "inventory", "finance"
	IsActive    bool   `json:"is_active"`
}

type EmployeePermission struct {
	ID         int  `json:"id"`
	EmployeeID int  `json:"employee_id"`
	PageID     int  `json:"page_id"`
	CanCreate  bool `json:"can_create"`
	CanView    bool `json:"can_view"`
	CanUpdate  bool `json:"can_update"`
	CanDelete  bool `json:"can_delete"`
}

type EmployeePermissionWithPage struct {
	EmployeePermission
	Page Page `json:"page"`
}

// ============================================
// Role Request/Response DTOs
// ============================================

type CreateRoleRequest struct {
	RoleName    string `json:"role_name"`
	Description string `json:"description,omitempty"`
}

type UpdateRoleRequest struct {
	RoleName    *string `json:"role_name,omitempty"`
	Description *string `json:"description,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

type AssignPermissionRequest struct {
	EmployeeID int  `json:"employee_id"`
	PageID     int  `json:"page_id"`
	CanCreate  bool `json:"can_create"`
	CanView    bool `json:"can_view"`
	CanUpdate  bool `json:"can_update"`
	CanDelete  bool `json:"can_delete"`
}

type BulkAssignPermissionsRequest struct {
	EmployeeID  int                    `json:"employee_id"`
	Permissions []PermissionAssignment `json:"permissions"`
}

type PermissionAssignment struct {
	PageID    int  `json:"page_id"`
	CanCreate bool `json:"can_create"`
	CanView   bool `json:"can_view"`
	CanUpdate bool `json:"can_update"`
	CanDelete bool `json:"can_delete"`
}

// ============================================
// Page Management DTOs
// ============================================

type CreatePageRequest struct {
	RouteName   string `json:"route_name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description,omitempty"`
	Module      string `json:"module"`
	Icon        string `json:"icon,omitempty"`
	ParentID    *int   `json:"parent_id,omitempty"`
	SortOrder   int    `json:"sort_order,omitempty"`
}

type UpdatePageRequest struct {
	RouteName   *string `json:"route_name,omitempty"`
	DisplayName *string `json:"display_name,omitempty"`
	Description *string `json:"description,omitempty"`
	Module      *string `json:"module,omitempty"`
	Icon        *string `json:"icon,omitempty"`
	ParentID    *int    `json:"parent_id,omitempty"`
	SortOrder   *int    `json:"sort_order,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

type PageWithPermissions struct {
	Page      Page `json:"page"`
	CanCreate bool `json:"can_create"`
	CanView   bool `json:"can_view"`
	CanUpdate bool `json:"can_update"`
	CanDelete bool `json:"can_delete"`
}

// ============================================
// Role Template (for assigning standard permissions to roles)
// ============================================

type RoleTemplate struct {
	ID        int  `json:"id"`
	RoleID    int  `json:"role_id"`
	PageID    int  `json:"page_id"`
	CanCreate bool `json:"can_create"`
	CanView   bool `json:"can_view"`
	CanUpdate bool `json:"can_update"`
	CanDelete bool `json:"can_delete"`
}

type RoleWithPermissions struct {
	Role        Role           `json:"role"`
	Permissions []RoleTemplate `json:"permissions"`
}

type AssignRoleTemplateRequest struct {
	RoleID      int                    `json:"role_id"`
	Permissions []PermissionAssignment `json:"permissions"`
}

type ApplyRoleToEmployeeRequest struct {
	EmployeeID int `json:"employee_id"`
	RoleID     int `json:"role_id"`
}

// ============================================
// Security Level
// ============================================

type SecurityLevel int

const (
	SecurityLevelViewer     SecurityLevel = 1  // View only
	SecurityLevelOperator   SecurityLevel = 2  // Basic operations
	SecurityLevelUser       SecurityLevel = 3  // Standard user
	SecurityLevelSuperUser  SecurityLevel = 5  // Advanced user
	SecurityLevelManager    SecurityLevel = 7  // Department manager
	SecurityLevelAdmin      SecurityLevel = 9  // Administrator
	SecurityLevelSuperAdmin SecurityLevel = 10 // Super administrator
)

type SecurityLevelInfo struct {
	Level       SecurityLevel `json:"level"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
}

// ============================================
// User Access Summary
// ============================================

type UserAccessSummary struct {
	EmployeeID    int                   `json:"employee_id"`
	EmployeeName  string                `json:"employee_name"`
	Email         string                `json:"email"`
	RoleID        int                   `json:"role_id"`
	RoleName      string                `json:"role_name"`
	SecurityLevel SecurityLevel         `json:"security_level"`
	Permissions   []PageWithPermissions `json:"permissions"`
	ModuleAccess  map[string]bool       `json:"module_access"`
}
