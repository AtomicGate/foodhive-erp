package employee

import (
	"net/http"
	"strconv"

	authMiddleware "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/auth"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/models"
	employeeService "github.com/anas-dev-92/FoodHive/registration/src/v1/services/employee"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/utils/helper"
	"github.com/go-chi/chi/v5"
)

func HandlerCreate(service employeeService.EmployeeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.CreateEmployeeRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		// Validate
		v := helper.New()
		v.Check(req.Email != "", "email", "must be provided")
		v.Check(helper.ValidEmail(req.Email), "email", "must be a valid email address")
		v.Check(req.Password != "", "password", "must be provided")
		v.Check(len(req.Password) >= 8, "password", "must be at least 8 characters")
		v.Check(req.EnglishName != "", "english_name", "must be provided")
		v.Check(req.RoleID > 0, "role_id", "must be provided")

		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		// Get current user ID from context
		userID, ok := authMiddleware.GetUserID(r.Context())
		if !ok {
			helper.UnauthorizedResponse(w, r)
			return
		}

		id, err := service.Create(r.Context(), req, userID)
		if err != nil {
			switch err {
			case employeeService.ErrDuplicateEmail:
				helper.ErrorResponse(w, r, http.StatusConflict, "email already exists")
			case employeeService.ErrRoleNotFound:
				helper.ErrorResponse(w, r, http.StatusBadRequest, "role not found")
			default:
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.CreatedResponse(w, r, id, "employee created successfully")
	}
}

func HandlerGetByID(service employeeService.EmployeeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id < 1 {
			helper.NotFoundResponse(w, r)
			return
		}

		employee, err := service.GetByID(r.Context(), id)
		if err != nil {
			if err == employeeService.ErrNotFound {
				helper.NotFoundResponse(w, r)
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		// Don't return password
		employee.Employee.Password = ""

		helper.SuccessResponse(w, r, http.StatusOK, employee)
	}
}

func HandlerUpdate(service employeeService.EmployeeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id < 1 {
			helper.NotFoundResponse(w, r)
			return
		}

		var req models.UpdateEmployeeRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		err = service.Update(r.Context(), id, req)
		if err != nil {
			if err == employeeService.ErrNotFound {
				helper.NotFoundResponse(w, r)
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "employee updated successfully"})
	}
}

func HandlerDelete(service employeeService.EmployeeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id < 1 {
			helper.NotFoundResponse(w, r)
			return
		}

		err = service.Delete(r.Context(), id)
		if err != nil {
			if err == employeeService.ErrNotFound {
				helper.NotFoundResponse(w, r)
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "employee deleted successfully"})
	}
}

func HandlerList(service employeeService.EmployeeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters
		filters := models.EmployeeListFilters{
			Search:   r.URL.Query().Get("search"),
			Page:     1,
			PageSize: 20,
		}

		if page, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && page > 0 {
			filters.Page = page
		}
		if pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size")); err == nil && pageSize > 0 {
			filters.PageSize = pageSize
		}
		if deptID, err := strconv.Atoi(r.URL.Query().Get("department_id")); err == nil && deptID > 0 {
			filters.DepartmentID = &deptID
		}
		if roleID, err := strconv.Atoi(r.URL.Query().Get("role_id")); err == nil && roleID > 0 {
			filters.RoleID = &roleID
		}
		if status := r.URL.Query().Get("status"); status != "" {
			filters.Status = &status
		}

		employees, total, err := service.List(r.Context(), filters)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		// Calculate total pages
		totalPages := int(total) / filters.PageSize
		if int(total)%filters.PageSize != 0 {
			totalPages++
		}

		response := models.PaginatedResponse{
			Data: employees,
			Pagination: models.Pagination{
				Page:       filters.Page,
				PageSize:   filters.PageSize,
				TotalItems: total,
				TotalPages: totalPages,
			},
		}

		helper.SuccessResponse(w, r, http.StatusOK, response)
	}
}
