package gl

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/anas-dev-92/FoodHive/core/auth"
	"github.com/anas-dev-92/FoodHive/core/jwt"
	"github.com/anas-dev-92/FoodHive/core/postgres"
	authMiddleware "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/auth"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/models"
	glService "github.com/anas-dev-92/FoodHive/registration/src/v1/services/gl"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/utils/helper"
	"github.com/go-chi/chi/v5"
)

func Router(db postgres.Executor, jwtService jwt.JWTService, authService auth.AuthService) chi.Router {
	r := chi.NewRouter()

	service := glService.New(db.(postgres.Connection))

	r.Use(authMiddleware.Authenticate(jwtService))

	// ===== Chart of Accounts =====
	r.With(authMiddleware.Authorize(jwtService)).Post("/accounts", handleCreateAccount(service))
	r.With(authMiddleware.Authorize(jwtService)).Get("/accounts/{id}", handleGetAccountByID(service))
	r.With(authMiddleware.Authorize(jwtService)).Get("/accounts/code/{code}", handleGetAccountByCode(service))
	r.With(authMiddleware.Authorize(jwtService)).Put("/accounts/{id}", handleUpdateAccount(service))
	r.With(authMiddleware.Authorize(jwtService)).Delete("/accounts/{id}", handleDeleteAccount(service))
	r.With(authMiddleware.Authorize(jwtService)).Get("/accounts", handleListAccounts(service))
	r.With(authMiddleware.Authorize(jwtService)).Get("/chart-of-accounts", handleGetChartOfAccounts(service))

	// ===== Fiscal Years =====
	r.With(authMiddleware.Authorize(jwtService)).Post("/fiscal-years", handleCreateFiscalYear(service))
	r.With(authMiddleware.Authorize(jwtService)).Get("/fiscal-years/{id}", handleGetFiscalYearByID(service))
	r.With(authMiddleware.Authorize(jwtService)).Get("/fiscal-years/current", handleGetCurrentFiscalYear(service))
	r.With(authMiddleware.Authorize(jwtService)).Get("/fiscal-years", handleListFiscalYears(service))
	r.With(authMiddleware.Authorize(jwtService)).Post("/fiscal-years/{id}/close", handleCloseFiscalYear(service, jwtService))

	// ===== Periods =====
	r.With(authMiddleware.Authorize(jwtService)).Get("/periods/{id}", handleGetPeriodByID(service))
	r.With(authMiddleware.Authorize(jwtService)).Get("/periods/current", handleGetCurrentPeriod(service))
	r.With(authMiddleware.Authorize(jwtService)).Get("/fiscal-years/{fiscalYearId}/periods", handleListPeriods(service))
	r.With(authMiddleware.Authorize(jwtService)).Post("/periods/{id}/close", handleClosePeriod(service, jwtService))
	r.With(authMiddleware.Authorize(jwtService)).Post("/periods/{id}/reopen", handleReopenPeriod(service))

	// ===== Journal Entries =====
	r.With(authMiddleware.Authorize(jwtService)).Post("/journal-entries", handleCreateJournalEntry(service, jwtService))
	r.With(authMiddleware.Authorize(jwtService)).Get("/journal-entries/{id}", handleGetJournalEntryByID(service))
	r.With(authMiddleware.Authorize(jwtService)).Put("/journal-entries/{id}", handleUpdateJournalEntry(service))
	r.With(authMiddleware.Authorize(jwtService)).Delete("/journal-entries/{id}", handleDeleteJournalEntry(service))
	r.With(authMiddleware.Authorize(jwtService)).Get("/journal-entries", handleListJournalEntries(service))
	r.With(authMiddleware.Authorize(jwtService)).Post("/journal-entries/{id}/post", handlePostJournalEntry(service, jwtService))
	r.With(authMiddleware.Authorize(jwtService)).Post("/journal-entries/{id}/reverse", handleReverseJournalEntry(service, jwtService))
	r.With(authMiddleware.Authorize(jwtService)).Post("/journal-entries/{id}/void", handleVoidJournalEntry(service))

	// ===== Reports =====
	r.With(authMiddleware.Authorize(jwtService)).Get("/reports/trial-balance", handleGetTrialBalance(service))
	r.With(authMiddleware.Authorize(jwtService)).Get("/reports/income-statement", handleGetIncomeStatement(service))
	r.With(authMiddleware.Authorize(jwtService)).Get("/reports/balance-sheet", handleGetBalanceSheet(service))
	r.With(authMiddleware.Authorize(jwtService)).Get("/reports/account-activity/{accountId}", handleGetAccountActivity(service))

	return r
}

// ============================================
// Account Handlers
// ============================================

func handleCreateAccount(service glService.GLService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.CreateGLAccountRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		v := models.NewValidator()
		models.ValidateGLAccount(v, &req)
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		id, err := service.CreateAccount(r.Context(), req)
		if err != nil {
			if errors.Is(err, glService.ErrDuplicateCode) {
				helper.BadRequestResponse(w, r, errors.New("account code already exists"))
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.CreatedResponse(w, r, id, "GL account created")
	}
}

func handleGetAccountByID(service glService.GLService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid account ID"))
			return
		}

		account, err := service.GetAccountByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, glService.ErrAccountNotFound) {
				helper.NotFoundResponse(w, r)
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.WriteJSON(w, http.StatusOK, helper.Envelope{"account": account}, nil)
	}
}

func handleGetAccountByCode(service glService.GLService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := chi.URLParam(r, "code")
		if code == "" {
			helper.BadRequestResponse(w, r, errors.New("account code is required"))
			return
		}

		account, err := service.GetAccountByCode(r.Context(), code)
		if err != nil {
			if errors.Is(err, glService.ErrAccountNotFound) {
				helper.NotFoundResponse(w, r)
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.WriteJSON(w, http.StatusOK, helper.Envelope{"account": account}, nil)
	}
}

func handleUpdateAccount(service glService.GLService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid account ID"))
			return
		}

		var req models.UpdateGLAccountRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		if err := service.UpdateAccount(r.Context(), id, req); err != nil {
			if errors.Is(err, glService.ErrAccountNotFound) {
				helper.NotFoundResponse(w, r)
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.WriteJSON(w, http.StatusOK, helper.Envelope{"message": "account updated"}, nil)
	}
}

func handleDeleteAccount(service glService.GLService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid account ID"))
			return
		}

		if err := service.DeleteAccount(r.Context(), id); err != nil {
			if errors.Is(err, glService.ErrAccountNotFound) {
				helper.NotFoundResponse(w, r)
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.WriteJSON(w, http.StatusOK, helper.Envelope{"message": "account deleted"}, nil)
	}
}

func handleListAccounts(service glService.GLService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filters := models.GLAccountListFilters{
			Search:   r.URL.Query().Get("search"),
			Page:     1,
			PageSize: 50,
		}

		if accountType := r.URL.Query().Get("account_type"); accountType != "" {
			t := models.GLAccountType(accountType)
			filters.AccountType = &t
		}

		if isActive := r.URL.Query().Get("is_active"); isActive != "" {
			b := isActive == "true"
			filters.IsActive = &b
		}

		if isPostable := r.URL.Query().Get("is_postable"); isPostable != "" {
			b := isPostable == "true"
			filters.IsPostable = &b
		}

		if page := r.URL.Query().Get("page"); page != "" {
			if p, err := strconv.Atoi(page); err == nil && p > 0 {
				filters.Page = p
			}
		}

		if pageSize := r.URL.Query().Get("page_size"); pageSize != "" {
			if ps, err := strconv.Atoi(pageSize); err == nil && ps > 0 && ps <= 100 {
				filters.PageSize = ps
			}
		}

		accounts, total, err := service.ListAccounts(r.Context(), filters)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.WriteJSON(w, http.StatusOK, helper.Envelope{
			"accounts":  accounts,
			"total":     total,
			"page":      filters.Page,
			"page_size": filters.PageSize,
		}, nil)
	}
}

func handleGetChartOfAccounts(service glService.GLService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chart, err := service.GetChartOfAccounts(r.Context())
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.WriteJSON(w, http.StatusOK, helper.Envelope{"chart_of_accounts": chart}, nil)
	}
}

// ============================================
// Fiscal Year Handlers
// ============================================

func handleCreateFiscalYear(service glService.GLService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.CreateFiscalYearRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		v := models.NewValidator()
		models.ValidateFiscalYear(v, &req)
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		id, err := service.CreateFiscalYear(r.Context(), req)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.CreatedResponse(w, r, id, "fiscal year created with 13 periods")
	}
}

func handleGetFiscalYearByID(service glService.GLService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid fiscal year ID"))
			return
		}

		fy, err := service.GetFiscalYearByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, glService.ErrFiscalYearNotFound) {
				helper.NotFoundResponse(w, r)
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.WriteJSON(w, http.StatusOK, helper.Envelope{"fiscal_year": fy}, nil)
	}
}

func handleGetCurrentFiscalYear(service glService.GLService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fy, err := service.GetCurrentFiscalYear(r.Context())
		if err != nil {
			if errors.Is(err, glService.ErrFiscalYearNotFound) {
				helper.NotFoundResponse(w, r)
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.WriteJSON(w, http.StatusOK, helper.Envelope{"fiscal_year": fy}, nil)
	}
}

func handleListFiscalYears(service glService.GLService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		years, err := service.ListFiscalYears(r.Context())
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.WriteJSON(w, http.StatusOK, helper.Envelope{"fiscal_years": years}, nil)
	}
}

func handleCloseFiscalYear(service glService.GLService, jwtService jwt.JWTService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid fiscal year ID"))
			return
		}

		tokenData, err := jwtService.ParseTokenFromRequest(r)
		if err != nil {
			helper.UnauthorizedResponse(w, r)
			return
		}
		userID := int(tokenData["user_id"].(float64))

		if err := service.CloseFiscalYear(r.Context(), id, userID); err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.WriteJSON(w, http.StatusOK, helper.Envelope{"message": "fiscal year closed"}, nil)
	}
}

// ============================================
// Period Handlers
// ============================================

func handleGetPeriodByID(service glService.GLService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid period ID"))
			return
		}

		period, err := service.GetPeriodByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, glService.ErrPeriodNotFound) {
				helper.NotFoundResponse(w, r)
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.WriteJSON(w, http.StatusOK, helper.Envelope{"period": period}, nil)
	}
}

func handleGetCurrentPeriod(service glService.GLService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		period, err := service.GetCurrentPeriod(r.Context())
		if err != nil {
			if errors.Is(err, glService.ErrPeriodNotFound) {
				helper.NotFoundResponse(w, r)
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.WriteJSON(w, http.StatusOK, helper.Envelope{"period": period}, nil)
	}
}

func handleListPeriods(service glService.GLService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fiscalYearID, err := strconv.Atoi(chi.URLParam(r, "fiscalYearId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid fiscal year ID"))
			return
		}

		periods, err := service.ListPeriods(r.Context(), fiscalYearID)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.WriteJSON(w, http.StatusOK, helper.Envelope{"periods": periods}, nil)
	}
}

func handleClosePeriod(service glService.GLService, jwtService jwt.JWTService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid period ID"))
			return
		}

		tokenData, err := jwtService.ParseTokenFromRequest(r)
		if err != nil {
			helper.UnauthorizedResponse(w, r)
			return
		}
		userID := int(tokenData["user_id"].(float64))

		if err := service.ClosePeriod(r.Context(), id, userID); err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.WriteJSON(w, http.StatusOK, helper.Envelope{"message": "period closed"}, nil)
	}
}

func handleReopenPeriod(service glService.GLService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid period ID"))
			return
		}

		if err := service.ReopenPeriod(r.Context(), id); err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.WriteJSON(w, http.StatusOK, helper.Envelope{"message": "period reopened"}, nil)
	}
}

// ============================================
// Journal Entry Handlers
// ============================================

func handleCreateJournalEntry(service glService.GLService, jwtService jwt.JWTService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.CreateJournalEntryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		v := models.NewValidator()
		models.ValidateJournalEntry(v, &req)
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		tokenData, err := jwtService.ParseTokenFromRequest(r)
		if err != nil {
			helper.UnauthorizedResponse(w, r)
			return
		}
		userID := int(tokenData["user_id"].(float64))

		id, err := service.CreateJournalEntry(r.Context(), req, userID)
		if err != nil {
			switch {
			case errors.Is(err, glService.ErrPeriodClosed):
				helper.BadRequestResponse(w, r, errors.New("no open period for this date"))
			case errors.Is(err, glService.ErrAccountNotPostable):
				helper.BadRequestResponse(w, r, errors.New("one or more accounts are not postable"))
			default:
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.CreatedResponse(w, r, id, "journal entry created")
	}
}

func handleGetJournalEntryByID(service glService.GLService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid journal entry ID"))
			return
		}

		entry, err := service.GetJournalEntryByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, glService.ErrJournalNotFound) {
				helper.NotFoundResponse(w, r)
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.WriteJSON(w, http.StatusOK, helper.Envelope{"journal_entry": entry}, nil)
	}
}

func handleUpdateJournalEntry(service glService.GLService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid journal entry ID"))
			return
		}

		var req models.CreateJournalEntryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		if err := service.UpdateJournalEntry(r.Context(), id, req); err != nil {
			switch {
			case errors.Is(err, glService.ErrJournalNotFound):
				helper.NotFoundResponse(w, r)
			case errors.Is(err, glService.ErrJournalPosted):
				helper.BadRequestResponse(w, r, errors.New("cannot modify posted journal entry"))
			default:
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.WriteJSON(w, http.StatusOK, helper.Envelope{"message": "journal entry updated"}, nil)
	}
}

func handleDeleteJournalEntry(service glService.GLService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid journal entry ID"))
			return
		}

		if err := service.DeleteJournalEntry(r.Context(), id); err != nil {
			switch {
			case errors.Is(err, glService.ErrJournalNotFound):
				helper.NotFoundResponse(w, r)
			case errors.Is(err, glService.ErrJournalPosted):
				helper.BadRequestResponse(w, r, errors.New("cannot delete posted journal entry"))
			default:
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.WriteJSON(w, http.StatusOK, helper.Envelope{"message": "journal entry deleted"}, nil)
	}
}

func handleListJournalEntries(service glService.GLService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filters := models.JournalEntryListFilters{
			DateFrom: r.URL.Query().Get("date_from"),
			DateTo:   r.URL.Query().Get("date_to"),
			Search:   r.URL.Query().Get("search"),
			Page:     1,
			PageSize: 50,
		}

		if periodID := r.URL.Query().Get("period_id"); periodID != "" {
			if id, err := strconv.Atoi(periodID); err == nil {
				filters.PeriodID = &id
			}
		}

		if entryType := r.URL.Query().Get("entry_type"); entryType != "" {
			t := models.JournalEntryType(entryType)
			filters.EntryType = &t
		}

		if status := r.URL.Query().Get("status"); status != "" {
			s := models.JournalEntryStatus(status)
			filters.Status = &s
		}

		if page := r.URL.Query().Get("page"); page != "" {
			if p, err := strconv.Atoi(page); err == nil && p > 0 {
				filters.Page = p
			}
		}

		if pageSize := r.URL.Query().Get("page_size"); pageSize != "" {
			if ps, err := strconv.Atoi(pageSize); err == nil && ps > 0 && ps <= 100 {
				filters.PageSize = ps
			}
		}

		entries, total, err := service.ListJournalEntries(r.Context(), filters)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.WriteJSON(w, http.StatusOK, helper.Envelope{
			"journal_entries": entries,
			"total":           total,
			"page":            filters.Page,
			"page_size":       filters.PageSize,
		}, nil)
	}
}

func handlePostJournalEntry(service glService.GLService, jwtService jwt.JWTService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid journal entry ID"))
			return
		}

		tokenData, err := jwtService.ParseTokenFromRequest(r)
		if err != nil {
			helper.UnauthorizedResponse(w, r)
			return
		}
		userID := int(tokenData["user_id"].(float64))

		if err := service.PostJournalEntry(r.Context(), id, userID); err != nil {
			switch {
			case errors.Is(err, glService.ErrJournalNotFound):
				helper.NotFoundResponse(w, r)
			case errors.Is(err, glService.ErrJournalPosted):
				helper.BadRequestResponse(w, r, errors.New("journal entry already posted"))
			case errors.Is(err, glService.ErrUnbalancedEntry):
				helper.BadRequestResponse(w, r, errors.New("journal entry is unbalanced"))
			default:
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.WriteJSON(w, http.StatusOK, helper.Envelope{"message": "journal entry posted"}, nil)
	}
}

func handleReverseJournalEntry(service glService.GLService, jwtService jwt.JWTService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid journal entry ID"))
			return
		}

		var req struct {
			ReversalDate string `json:"reversal_date"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		tokenData, err := jwtService.ParseTokenFromRequest(r)
		if err != nil {
			helper.UnauthorizedResponse(w, r)
			return
		}
		userID := int(tokenData["user_id"].(float64))

		reversalID, err := service.ReverseJournalEntry(r.Context(), id, req.ReversalDate, userID)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.WriteJSON(w, http.StatusOK, helper.Envelope{
			"message":           "journal entry reversed",
			"reversal_entry_id": reversalID,
		}, nil)
	}
}

func handleVoidJournalEntry(service glService.GLService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid journal entry ID"))
			return
		}

		if err := service.VoidJournalEntry(r.Context(), id); err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.WriteJSON(w, http.StatusOK, helper.Envelope{"message": "journal entry voided"}, nil)
	}
}

// ============================================
// Report Handlers
// ============================================

func handleGetTrialBalance(service glService.GLService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filters := models.ReportFilters{
			DateFrom: r.URL.Query().Get("date_from"),
			DateTo:   r.URL.Query().Get("date_to"),
		}

		report, err := service.GetTrialBalance(r.Context(), filters)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.WriteJSON(w, http.StatusOK, helper.Envelope{"trial_balance": report}, nil)
	}
}

func handleGetIncomeStatement(service glService.GLService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filters := models.ReportFilters{
			DateFrom: r.URL.Query().Get("date_from"),
			DateTo:   r.URL.Query().Get("date_to"),
		}

		report, err := service.GetIncomeStatement(r.Context(), filters)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.WriteJSON(w, http.StatusOK, helper.Envelope{"income_statement": report}, nil)
	}
}

func handleGetBalanceSheet(service glService.GLService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filters := models.ReportFilters{
			DateFrom: r.URL.Query().Get("date_from"),
			DateTo:   r.URL.Query().Get("date_to"),
		}

		report, err := service.GetBalanceSheet(r.Context(), filters)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.WriteJSON(w, http.StatusOK, helper.Envelope{"balance_sheet": report}, nil)
	}
}

func handleGetAccountActivity(service glService.GLService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountID, err := strconv.Atoi(chi.URLParam(r, "accountId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid account ID"))
			return
		}

		dateFrom := r.URL.Query().Get("date_from")
		dateTo := r.URL.Query().Get("date_to")

		if dateFrom == "" || dateTo == "" {
			helper.BadRequestResponse(w, r, errors.New("date_from and date_to are required"))
			return
		}

		report, err := service.GetAccountActivity(r.Context(), accountID, dateFrom, dateTo)
		if err != nil {
			if errors.Is(err, glService.ErrAccountNotFound) {
				helper.NotFoundResponse(w, r)
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.WriteJSON(w, http.StatusOK, helper.Envelope{"account_activity": report}, nil)
	}
}
