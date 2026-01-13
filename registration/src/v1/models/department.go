package models

// ============================================
// Department Models
// ============================================

type Department struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
	ManagerID   *int   `json:"manager_id,omitempty"`
	ParentID    *int   `json:"parent_id,omitempty"` // For hierarchical departments
	IsActive    bool   `json:"is_active"`
	CreatedAt   CustomDateTime `json:"created_at,omitempty"`
	UpdatedAt   CustomDateTime `json:"updated_at,omitempty"`
}

type DepartmentWithManager struct {
	Department
	ManagerName string `json:"manager_name,omitempty"`
}

// ============================================
// Department Request/Response DTOs
// ============================================

type CreateDepartmentRequest struct {
	Name        string `json:"name"`
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
	ManagerID   *int   `json:"manager_id,omitempty"`
	ParentID    *int   `json:"parent_id,omitempty"`
}

type UpdateDepartmentRequest struct {
	Name        *string `json:"name,omitempty"`
	Code        *string `json:"code,omitempty"`
	Description *string `json:"description,omitempty"`
	ManagerID   *int    `json:"manager_id,omitempty"`
	ParentID    *int    `json:"parent_id,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}
