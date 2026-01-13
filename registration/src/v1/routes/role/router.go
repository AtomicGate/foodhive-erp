package role

import (
	"net/http"
	"strconv"

	"github.com/anas-dev-92/FoodHive/core/auth"
	"github.com/anas-dev-92/FoodHive/core/jwt"
	"github.com/anas-dev-92/FoodHive/core/postgres"
	authMiddleware "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/auth"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/models"
	roleService "github.com/anas-dev-92/FoodHive/registration/src/v1/services/role"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/utils/helper"
	"github.com/go-chi/chi/v5"
)

func Router(db postgres.Executor, jwtService jwt.JWTService, authService auth.AuthService) chi.Router {
	r := chi.NewRouter()

	service := roleService.New(db.(postgres.Connection))
	permService := roleService.NewPermissionService(db.(postgres.Connection))
	pageService := roleService.NewPageService(db.(postgres.Connection))
	templateService := roleService.NewRoleTemplateService(db.(postgres.Connection))

	r.Use(authMiddleware.Authenticate(jwtService))

	// Role routes
	r.With(authMiddleware.Authorize(jwtService)).Post("/create", handleCreate(service))
	r.With(authMiddleware.Authorize(jwtService)).Get("/get/{id}", handleGetByID(service))
	r.With(authMiddleware.Authorize(jwtService)).Put("/update/{id}", handleUpdate(service))
	r.With(authMiddleware.Authorize(jwtService)).Delete("/delete/{id}", handleDelete(service))
	r.With(authMiddleware.Authorize(jwtService)).Get("/list", handleList(service))

	// Permission routes
	r.With(authMiddleware.Authorize(jwtService)).Post("/permissions/assign", handleAssignPermission(permService))
	r.With(authMiddleware.Authorize(jwtService)).Post("/permissions/bulk-assign", handleBulkAssignPermissions(permService))
	r.With(authMiddleware.Authorize(jwtService)).Get("/permissions/employee/{id}", handleGetEmployeePermissions(permService))

	// Pages (system modules) management
	r.With(authMiddleware.Authorize(jwtService)).Post("/pages/create", handleCreatePage(pageService))
	r.With(authMiddleware.Authorize(jwtService)).Get("/pages/get/{id}", handleGetPage(pageService))
	r.With(authMiddleware.Authorize(jwtService)).Put("/pages/update/{id}", handleUpdatePage(pageService))
	r.With(authMiddleware.Authorize(jwtService)).Delete("/pages/delete/{id}", handleDeletePage(pageService))
	r.With(authMiddleware.Authorize(jwtService)).Get("/pages/list", handleListPages(pageService))
	r.With(authMiddleware.Authorize(jwtService)).Get("/pages/modules", handleListModules(pageService))

	// Role template (standard permissions per role)
	r.With(authMiddleware.Authorize(jwtService)).Post("/templates/set", handleSetRolePermissions(templateService))
	r.With(authMiddleware.Authorize(jwtService)).Get("/templates/get/{roleId}", handleGetRolePermissions(templateService))
	r.With(authMiddleware.Authorize(jwtService)).Post("/templates/apply", handleApplyRoleToEmployee(templateService))

	// User access summary
	r.With(authMiddleware.Authorize(jwtService)).Get("/access/summary/{employeeId}", handleGetUserAccessSummary(templateService))
	r.With(authMiddleware.Authorize(jwtService)).Get("/access/security-levels", handleGetSecurityLevels())

	return r
}

func handleCreate(service roleService.RoleService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.CreateRoleRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		v := helper.New()
		v.Check(req.RoleName != "", "role_name", "must be provided")
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		id, err := service.Create(r.Context(), req)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.CreatedResponse(w, r, id, "role created successfully")
	}
}

func handleGetByID(service roleService.RoleService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || id < 1 {
			helper.NotFoundResponse(w, r)
			return
		}

		role, err := service.GetByID(r.Context(), id)
		if err != nil {
			if err == roleService.ErrNotFound {
				helper.NotFoundResponse(w, r)
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, role)
	}
}

func handleUpdate(service roleService.RoleService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || id < 1 {
			helper.NotFoundResponse(w, r)
			return
		}

		var req models.UpdateRoleRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		err = service.Update(r.Context(), id, req)
		if err != nil {
			if err == roleService.ErrNotFound {
				helper.NotFoundResponse(w, r)
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "role updated successfully"})
	}
}

func handleDelete(service roleService.RoleService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || id < 1 {
			helper.NotFoundResponse(w, r)
			return
		}

		err = service.Delete(r.Context(), id)
		if err != nil {
			if err == roleService.ErrNotFound {
				helper.NotFoundResponse(w, r)
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "role deleted successfully"})
	}
}

func handleList(service roleService.RoleService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roles, err := service.List(r.Context())
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, roles)
	}
}

func handleAssignPermission(service roleService.PermissionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.AssignPermissionRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		v := helper.New()
		v.Check(req.EmployeeID > 0, "employee_id", "must be provided")
		v.Check(req.PageID > 0, "page_id", "must be provided")
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		err := service.AssignPermission(r.Context(), req)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "permission assigned successfully"})
	}
}

func handleBulkAssignPermissions(service roleService.PermissionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.BulkAssignPermissionsRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		v := helper.New()
		v.Check(req.EmployeeID > 0, "employee_id", "must be provided")
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		err := service.BulkAssignPermissions(r.Context(), req)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "permissions assigned successfully"})
	}
}

func handleGetEmployeePermissions(service roleService.PermissionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || id < 1 {
			helper.NotFoundResponse(w, r)
			return
		}

		permissions, err := service.GetEmployeePermissions(r.Context(), id)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, permissions)
	}
}

// ============================================
// Page Management Handlers
// ============================================

func handleCreatePage(service roleService.PageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.CreatePageRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		v := helper.New()
		v.Check(req.RouteName != "", "route_name", "must be provided")
		v.Check(req.DisplayName != "", "display_name", "must be provided")
		v.Check(req.Module != "", "module", "must be provided")
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		id, err := service.CreatePage(r.Context(), req)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.CreatedResponse(w, r, id, "page created successfully")
	}
}

func handleGetPage(service roleService.PageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || id < 1 {
			helper.NotFoundResponse(w, r)
			return
		}

		page, err := service.GetPage(r.Context(), id)
		if err != nil {
			helper.NotFoundResponse(w, r)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, page)
	}
}

func handleUpdatePage(service roleService.PageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || id < 1 {
			helper.NotFoundResponse(w, r)
			return
		}

		var req models.UpdatePageRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		err = service.UpdatePage(r.Context(), id, req)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "page updated successfully"})
	}
}

func handleDeletePage(service roleService.PageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || id < 1 {
			helper.NotFoundResponse(w, r)
			return
		}

		err = service.DeletePage(r.Context(), id)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "page deleted successfully"})
	}
}

func handleListPages(service roleService.PageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		module := r.URL.Query().Get("module")

		pages, err := service.ListPages(r.Context(), module)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, pages)
	}
}

func handleListModules(service roleService.PageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		modules, err := service.ListModules(r.Context())
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, modules)
	}
}

// ============================================
// Role Template Handlers
// ============================================

func handleSetRolePermissions(service roleService.RoleTemplateService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.AssignRoleTemplateRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		v := helper.New()
		v.Check(req.RoleID > 0, "role_id", "must be provided")
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		err := service.SetRolePermissions(r.Context(), req)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "role permissions set successfully"})
	}
}

func handleGetRolePermissions(service roleService.RoleTemplateService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roleID, err := strconv.Atoi(chi.URLParam(r, "roleId"))
		if err != nil || roleID < 1 {
			helper.NotFoundResponse(w, r)
			return
		}

		rolePerms, err := service.GetRolePermissions(r.Context(), roleID)
		if err != nil {
			helper.NotFoundResponse(w, r)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, rolePerms)
	}
}

func handleApplyRoleToEmployee(service roleService.RoleTemplateService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.ApplyRoleToEmployeeRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		v := helper.New()
		v.Check(req.EmployeeID > 0, "employee_id", "must be provided")
		v.Check(req.RoleID > 0, "role_id", "must be provided")
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		err := service.ApplyRoleToEmployee(r.Context(), req.EmployeeID, req.RoleID)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "role applied to employee successfully"})
	}
}

func handleGetUserAccessSummary(service roleService.RoleTemplateService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		employeeID, err := strconv.Atoi(chi.URLParam(r, "employeeId"))
		if err != nil || employeeID < 1 {
			helper.NotFoundResponse(w, r)
			return
		}

		summary, err := service.GetUserAccessSummary(r.Context(), employeeID)
		if err != nil {
			helper.NotFoundResponse(w, r)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, summary)
	}
}

func handleGetSecurityLevels() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		levels := roleService.GetSecurityLevels()
		helper.SuccessResponse(w, r, http.StatusOK, levels)
	}
}
