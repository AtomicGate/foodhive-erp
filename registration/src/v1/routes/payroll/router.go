package payroll

import (
	"net/http"
	"strconv"

	"github.com/anas-dev-92/FoodHive/core/auth"
	"github.com/anas-dev-92/FoodHive/core/jwt"
	"github.com/anas-dev-92/FoodHive/core/postgres"
	authMiddleware "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/auth"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/models"
	payrollService "github.com/anas-dev-92/FoodHive/registration/src/v1/services/payroll"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/utils/helper"
	"github.com/go-chi/chi/v5"
)

func Router(db postgres.Executor, jwtService jwt.JWTService, authService auth.AuthService) chi.Router {
	r := chi.NewRouter()

	service := payrollService.New(db.(postgres.Connection))

	r.Use(authMiddleware.Authenticate(jwtService))

	// Payroll CRUD
	r.With(authMiddleware.Authorize(jwtService)).Post("/create", handleCreate(service))
	r.With(authMiddleware.Authorize(jwtService)).Get("/get/{id}", handleGetByID(service))
	r.With(authMiddleware.Authorize(jwtService)).Put("/update/{id}", handleUpdate(service))
	r.With(authMiddleware.Authorize(jwtService)).Delete("/delete/{id}", handleDelete(service))
	r.With(authMiddleware.Authorize(jwtService)).Get("/list", handleList(service))

	// Payroll actions
	r.With(authMiddleware.Authorize(jwtService)).Post("/{id}/lines/add", handleAddLine(service))
	r.With(authMiddleware.Authorize(jwtService)).Delete("/lines/{lineId}/delete", handleRemoveLine(service))
	r.With(authMiddleware.Authorize(jwtService)).Post("/{id}/calculate", handleCalculate(service))
	r.With(authMiddleware.Authorize(jwtService)).Post("/{id}/approve", handleApprove(service))

	return r
}

func handleCreate(service payrollService.PayrollService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.CreatePayrollRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		v := helper.New()
		v.Check(req.PayrollPeriod != "", "payroll_period", "must be provided (e.g., 2026-01)")
		v.Check(!req.PayDate.IsZero(), "pay_date", "must be provided")
		v.Check(req.Type != "", "type", "must be provided")
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		userID, ok := authMiddleware.GetUserID(r.Context())
		if !ok {
			helper.UnauthorizedResponse(w, r)
			return
		}

		id, err := service.Create(r.Context(), req, userID)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.CreatedResponse(w, r, id, "payroll created successfully")
	}
}

func handleGetByID(service payrollService.PayrollService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || id < 1 {
			helper.NotFoundResponse(w, r)
			return
		}

		payroll, err := service.GetByID(r.Context(), id)
		if err != nil {
			if err == payrollService.ErrNotFound {
				helper.NotFoundResponse(w, r)
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, payroll)
	}
}

func handleUpdate(service payrollService.PayrollService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || id < 1 {
			helper.NotFoundResponse(w, r)
			return
		}

		var req models.UpdatePayrollRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		err = service.Update(r.Context(), id, req)
		if err != nil {
			if err == payrollService.ErrNotFound {
				helper.NotFoundResponse(w, r)
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "payroll updated successfully"})
	}
}

func handleDelete(service payrollService.PayrollService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || id < 1 {
			helper.NotFoundResponse(w, r)
			return
		}

		err = service.Delete(r.Context(), id)
		if err != nil {
			if err == payrollService.ErrNotFound {
				helper.NotFoundResponse(w, r)
			} else if err == payrollService.ErrInvalidStatus {
				helper.ErrorResponse(w, r, http.StatusBadRequest, "can only delete draft payrolls")
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "payroll deleted successfully"})
	}
}

func handleList(service payrollService.PayrollService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filters := models.PayrollListFilters{
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
		if period := r.URL.Query().Get("payroll_period"); period != "" {
			filters.PayrollPeriod = &period
		}
		if status := r.URL.Query().Get("status"); status != "" {
			filters.Status = &status
		}
		if pType := r.URL.Query().Get("type"); pType != "" {
			filters.Type = &pType
		}

		payrolls, total, err := service.List(r.Context(), filters)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		totalPages := int(total) / filters.PageSize
		if int(total)%filters.PageSize != 0 {
			totalPages++
		}

		response := models.PaginatedResponse{
			Data: payrolls,
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

func handleAddLine(service payrollService.PayrollService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payrollID, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || payrollID < 1 {
			helper.NotFoundResponse(w, r)
			return
		}

		var req models.AddPayrollLineRequest
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

		id, err := service.AddLine(r.Context(), payrollID, req)
		if err != nil {
			if err == payrollService.ErrNotFound {
				helper.NotFoundResponse(w, r)
			} else if err == payrollService.ErrInvalidStatus {
				helper.ErrorResponse(w, r, http.StatusBadRequest, "can only add lines to draft payrolls")
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.CreatedResponse(w, r, id, "payroll line added successfully")
	}
}

func handleRemoveLine(service payrollService.PayrollService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lineID, err := strconv.Atoi(chi.URLParam(r, "lineId"))
		if err != nil || lineID < 1 {
			helper.NotFoundResponse(w, r)
			return
		}

		err = service.RemoveLine(r.Context(), lineID)
		if err != nil {
			if err == payrollService.ErrNotFound {
				helper.NotFoundResponse(w, r)
			} else if err == payrollService.ErrInvalidStatus {
				helper.ErrorResponse(w, r, http.StatusBadRequest, "can only remove lines from draft payrolls")
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "payroll line removed successfully"})
	}
}

func handleCalculate(service payrollService.PayrollService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || id < 1 {
			helper.NotFoundResponse(w, r)
			return
		}

		err = service.CalculatePayroll(r.Context(), id)
		if err != nil {
			if err == payrollService.ErrNotFound {
				helper.NotFoundResponse(w, r)
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "payroll calculated successfully"})
	}
}

func handleApprove(service payrollService.PayrollService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || id < 1 {
			helper.NotFoundResponse(w, r)
			return
		}

		userID, ok := authMiddleware.GetUserID(r.Context())
		if !ok {
			helper.UnauthorizedResponse(w, r)
			return
		}

		err = service.ApprovePayroll(r.Context(), id, userID)
		if err != nil {
			if err == payrollService.ErrNotFound {
				helper.NotFoundResponse(w, r)
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "payroll approved successfully"})
	}
}
