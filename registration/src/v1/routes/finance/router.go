package finance

import (
	"net/http"
	"strconv"
	"time"

	"github.com/anas-dev-92/FoodHive/core/auth"
	"github.com/anas-dev-92/FoodHive/core/jwt"
	"github.com/anas-dev-92/FoodHive/core/postgres"
	authMiddleware "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/auth"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/models"
	financeService "github.com/anas-dev-92/FoodHive/registration/src/v1/services/finance"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/utils/helper"
	"github.com/go-chi/chi/v5"
)

func Router(db postgres.Executor, jwtService jwt.JWTService, authService auth.AuthService) chi.Router {
	r := chi.NewRouter()

	incomeService := financeService.NewIncomeService(db.(postgres.Connection))
	expenseService := financeService.NewExpenseService(db.(postgres.Connection))
	cashBoxService := financeService.NewCashBoxService(db.(postgres.Connection))
	paymentTypeService := financeService.NewPaymentTypeService(db.(postgres.Connection))

	r.Use(authMiddleware.Authenticate(jwtService))

	// Income routes
	r.Route("/incomes", func(r chi.Router) {
		r.With(authMiddleware.Authorize(jwtService)).Post("/create", handleCreateIncome(incomeService))
		r.With(authMiddleware.Authorize(jwtService)).Get("/get/{id}", handleGetIncome(incomeService))
		r.With(authMiddleware.Authorize(jwtService)).Put("/update/{id}", handleUpdateIncome(incomeService))
		r.With(authMiddleware.Authorize(jwtService)).Delete("/delete/{id}", handleDeleteIncome(incomeService))
		r.With(authMiddleware.Authorize(jwtService)).Get("/list", handleListIncomes(incomeService))
	})

	// Expense routes
	r.Route("/expenses", func(r chi.Router) {
		r.With(authMiddleware.Authorize(jwtService)).Post("/create", handleCreateExpense(expenseService))
		r.With(authMiddleware.Authorize(jwtService)).Get("/get/{id}", handleGetExpense(expenseService))
		r.With(authMiddleware.Authorize(jwtService)).Put("/update/{id}", handleUpdateExpense(expenseService))
		r.With(authMiddleware.Authorize(jwtService)).Delete("/delete/{id}", handleDeleteExpense(expenseService))
		r.With(authMiddleware.Authorize(jwtService)).Get("/list", handleListExpenses(expenseService))
	})

	// Cash Box routes
	r.Route("/cash-boxes", func(r chi.Router) {
		r.With(authMiddleware.Authorize(jwtService)).Post("/create", handleCreateCashBox(cashBoxService))
		r.With(authMiddleware.Authorize(jwtService)).Get("/get/{id}", handleGetCashBox(cashBoxService))
		r.With(authMiddleware.Authorize(jwtService)).Put("/update/{id}", handleUpdateCashBox(cashBoxService))
		r.With(authMiddleware.Authorize(jwtService)).Delete("/delete/{id}", handleDeleteCashBox(cashBoxService))
		r.With(authMiddleware.Authorize(jwtService)).Get("/list", handleListCashBoxes(cashBoxService))
	})

	// Payment Type routes
	r.Route("/payment-types", func(r chi.Router) {
		r.With(authMiddleware.Authorize(jwtService)).Post("/create", handleCreatePaymentType(paymentTypeService))
		r.With(authMiddleware.Authorize(jwtService)).Get("/get/{id}", handleGetPaymentType(paymentTypeService))
		r.With(authMiddleware.Authorize(jwtService)).Delete("/delete/{id}", handleDeletePaymentType(paymentTypeService))
		r.With(authMiddleware.Authorize(jwtService)).Get("/list", handleListPaymentTypes(paymentTypeService))
	})

	return r
}

// ============================================
// Income Handlers
// ============================================

func handleCreateIncome(service financeService.IncomeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.CreateIncomeRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		v := helper.New()
		v.Check(req.TypeID > 0, "type_id", "must be provided")
		v.Check(req.Amount > 0, "amount", "must be positive")
		v.Check(req.CashBoxID > 0, "cash_box_id", "must be provided")
		v.Check(!req.Date.IsZero(), "date", "must be provided")
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

		helper.CreatedResponse(w, r, id, "income created successfully")
	}
}

func handleGetIncome(service financeService.IncomeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || id < 1 {
			helper.NotFoundResponse(w, r)
			return
		}

		income, err := service.GetByID(r.Context(), id)
		if err != nil {
			if err == financeService.ErrIncomeNotFound {
				helper.NotFoundResponse(w, r)
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, income)
	}
}

func handleUpdateIncome(service financeService.IncomeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || id < 1 {
			helper.NotFoundResponse(w, r)
			return
		}

		var req models.UpdateIncomeRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		err = service.Update(r.Context(), id, req)
		if err != nil {
			if err == financeService.ErrIncomeNotFound {
				helper.NotFoundResponse(w, r)
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "income updated successfully"})
	}
}

func handleDeleteIncome(service financeService.IncomeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || id < 1 {
			helper.NotFoundResponse(w, r)
			return
		}

		err = service.Delete(r.Context(), id)
		if err != nil {
			if err == financeService.ErrIncomeNotFound {
				helper.NotFoundResponse(w, r)
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "income deleted successfully"})
	}
}

func handleListIncomes(service financeService.IncomeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filters := models.IncomeListFilters{Page: 1, PageSize: 20}

		if page, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && page > 0 {
			filters.Page = page
		}
		if pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size")); err == nil && pageSize > 0 {
			filters.PageSize = pageSize
		}
		if typeID, err := strconv.Atoi(r.URL.Query().Get("type_id")); err == nil && typeID > 0 {
			filters.TypeID = &typeID
		}
		if cashBoxID, err := strconv.Atoi(r.URL.Query().Get("cash_box_id")); err == nil && cashBoxID > 0 {
			filters.CashBoxID = &cashBoxID
		}
		if dateFrom := r.URL.Query().Get("date_from"); dateFrom != "" {
			if t, err := time.Parse("2006-01-02", dateFrom); err == nil {
				cd := models.CustomDate(t)
				filters.DateFrom = &cd
			}
		}
		if dateTo := r.URL.Query().Get("date_to"); dateTo != "" {
			if t, err := time.Parse("2006-01-02", dateTo); err == nil {
				cd := models.CustomDate(t)
				filters.DateTo = &cd
			}
		}

		incomes, total, err := service.List(r.Context(), filters)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		totalPages := int(total) / filters.PageSize
		if int(total)%filters.PageSize != 0 {
			totalPages++
		}

		response := models.PaginatedResponse{
			Data: incomes,
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

// ============================================
// Expense Handlers
// ============================================

func handleCreateExpense(service financeService.ExpenseService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.CreateExpenseRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		v := helper.New()
		v.Check(req.TypeID > 0, "type_id", "must be provided")
		v.Check(req.Amount > 0, "amount", "must be positive")
		v.Check(req.CashBoxID > 0, "cash_box_id", "must be provided")
		v.Check(!req.Date.IsZero(), "date", "must be provided")
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

		helper.CreatedResponse(w, r, id, "expense created successfully")
	}
}

func handleGetExpense(service financeService.ExpenseService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || id < 1 {
			helper.NotFoundResponse(w, r)
			return
		}

		expense, err := service.GetByID(r.Context(), id)
		if err != nil {
			if err == financeService.ErrExpenseNotFound {
				helper.NotFoundResponse(w, r)
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, expense)
	}
}

func handleUpdateExpense(service financeService.ExpenseService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || id < 1 {
			helper.NotFoundResponse(w, r)
			return
		}

		var req models.UpdateExpenseRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		err = service.Update(r.Context(), id, req)
		if err != nil {
			if err == financeService.ErrExpenseNotFound {
				helper.NotFoundResponse(w, r)
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "expense updated successfully"})
	}
}

func handleDeleteExpense(service financeService.ExpenseService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || id < 1 {
			helper.NotFoundResponse(w, r)
			return
		}

		err = service.Delete(r.Context(), id)
		if err != nil {
			if err == financeService.ErrExpenseNotFound {
				helper.NotFoundResponse(w, r)
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "expense deleted successfully"})
	}
}

func handleListExpenses(service financeService.ExpenseService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filters := models.ExpenseListFilters{Page: 1, PageSize: 20}

		if page, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && page > 0 {
			filters.Page = page
		}
		if pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size")); err == nil && pageSize > 0 {
			filters.PageSize = pageSize
		}
		if typeID, err := strconv.Atoi(r.URL.Query().Get("type_id")); err == nil && typeID > 0 {
			filters.TypeID = &typeID
		}
		if cashBoxID, err := strconv.Atoi(r.URL.Query().Get("cash_box_id")); err == nil && cashBoxID > 0 {
			filters.CashBoxID = &cashBoxID
		}
		if isMarked := r.URL.Query().Get("is_marked"); isMarked != "" {
			marked := isMarked == "true"
			filters.IsMarked = &marked
		}

		expenses, total, err := service.List(r.Context(), filters)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		totalPages := int(total) / filters.PageSize
		if int(total)%filters.PageSize != 0 {
			totalPages++
		}

		response := models.PaginatedResponse{
			Data: expenses,
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

// ============================================
// Cash Box Handlers
// ============================================

func handleCreateCashBox(service financeService.CashBoxService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.CreateCashBoxRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		v := helper.New()
		v.Check(req.Name != "", "name", "must be provided")
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		id, err := service.Create(r.Context(), req)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.CreatedResponse(w, r, id, "cash box created successfully")
	}
}

func handleGetCashBox(service financeService.CashBoxService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || id < 1 {
			helper.NotFoundResponse(w, r)
			return
		}

		cashBox, err := service.GetByID(r.Context(), id)
		if err != nil {
			if err == financeService.ErrCashBoxNotFound {
				helper.NotFoundResponse(w, r)
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, cashBox)
	}
}

func handleUpdateCashBox(service financeService.CashBoxService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || id < 1 {
			helper.NotFoundResponse(w, r)
			return
		}

		var req models.UpdateCashBoxRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		err = service.Update(r.Context(), id, req)
		if err != nil {
			if err == financeService.ErrCashBoxNotFound {
				helper.NotFoundResponse(w, r)
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "cash box updated successfully"})
	}
}

func handleDeleteCashBox(service financeService.CashBoxService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || id < 1 {
			helper.NotFoundResponse(w, r)
			return
		}

		err = service.Delete(r.Context(), id)
		if err != nil {
			if err == financeService.ErrCashBoxNotFound {
				helper.NotFoundResponse(w, r)
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "cash box deleted successfully"})
	}
}

func handleListCashBoxes(service financeService.CashBoxService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cashBoxes, err := service.List(r.Context())
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, cashBoxes)
	}
}

// ============================================
// Payment Type Handlers
// ============================================

func handleCreatePaymentType(service financeService.PaymentTypeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.CreatePaymentTypeRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		v := helper.New()
		v.Check(req.Name != "", "name", "must be provided")
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		id, err := service.Create(r.Context(), req)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.CreatedResponse(w, r, id, "payment type created successfully")
	}
}

func handleGetPaymentType(service financeService.PaymentTypeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || id < 1 {
			helper.NotFoundResponse(w, r)
			return
		}

		pt, err := service.GetByID(r.Context(), id)
		if err != nil {
			if err == financeService.ErrPaymentTypeNotFound {
				helper.NotFoundResponse(w, r)
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, pt)
	}
}

func handleDeletePaymentType(service financeService.PaymentTypeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || id < 1 {
			helper.NotFoundResponse(w, r)
			return
		}

		err = service.Delete(r.Context(), id)
		if err != nil {
			if err == financeService.ErrPaymentTypeNotFound {
				helper.NotFoundResponse(w, r)
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "payment type deleted successfully"})
	}
}

func handleListPaymentTypes(service financeService.PaymentTypeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		paymentTypes, err := service.List(r.Context())
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, paymentTypes)
	}
}
