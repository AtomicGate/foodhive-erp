package ar

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/anas-dev-92/FoodHive/core/auth"
	"github.com/anas-dev-92/FoodHive/core/jwt"
	"github.com/anas-dev-92/FoodHive/core/postgres"
	arMiddleware "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/ar"
	authMiddleware "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/auth"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/models"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/utils/helper"
	"github.com/go-chi/chi/v5"
)

func Router(db postgres.Executor, jwtService jwt.JWTService, authService auth.AuthService) chi.Router {
	app := chi.NewRouter()

	// Inject AR service
	app.Use(arMiddleware.New(db))

	// Apply authentication
	app.Use(authMiddleware.Authenticate(jwtService))

	// ===========================================
	// Invoices
	// ===========================================
	app.With(authMiddleware.Authorize(jwtService)).Post("/invoices/create", handleCreateInvoice())
	app.With(authMiddleware.Authorize(jwtService)).Get("/invoices/get/{id}", handleGetInvoice())
	app.With(authMiddleware.Authorize(jwtService)).Get("/invoices/number/{number}", handleGetInvoiceByNumber())
	app.With(authMiddleware.Authorize(jwtService)).Get("/invoices/list", handleListInvoices())
	app.With(authMiddleware.Authorize(jwtService)).Post("/invoices/{id}/post", handlePostInvoice())
	app.With(authMiddleware.Authorize(jwtService)).Post("/invoices/{id}/void", handleVoidInvoice())
	app.With(authMiddleware.Authorize(jwtService)).Post("/invoices/from-order/{orderId}", handleCreateFromOrder())

	// ===========================================
	// Payments
	// ===========================================
	app.With(authMiddleware.Authorize(jwtService)).Post("/payments/create", handleCreatePayment())
	app.With(authMiddleware.Authorize(jwtService)).Get("/payments/get/{id}", handleGetPayment())
	app.With(authMiddleware.Authorize(jwtService)).Get("/payments/list", handleListPayments())

	// ===========================================
	// Credit Management
	// ===========================================
	app.With(authMiddleware.Authorize(jwtService)).Get("/credit/{customerId}", handleGetCustomerCredit())
	app.With(authMiddleware.Authorize(jwtService)).Get("/credit/{customerId}/check", handleCheckCredit())
	app.With(authMiddleware.Authorize(jwtService)).Put("/credit/{customerId}/limit", handleUpdateCreditLimit())

	// ===========================================
	// Aging & Statements
	// ===========================================
	app.With(authMiddleware.Authorize(jwtService)).Get("/aging/{customerId}", handleGetCustomerAging())
	app.With(authMiddleware.Authorize(jwtService)).Get("/aging/report", handleGetAgingReport())
	app.With(authMiddleware.Authorize(jwtService)).Get("/statement/{customerId}", handleGetStatement())
	app.With(authMiddleware.Authorize(jwtService)).Get("/overdue", handleGetOverdueInvoices())

	return app
}

// ===========================================
// Invoice Handlers
// ===========================================

func handleCreateInvoice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := arMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		var req models.CreateARInvoiceRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		v := models.NewValidator()
		models.ValidateARInvoice(v, &req)
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		createdBy := 1 // TODO: Get from auth context

		id, err := svc.CreateInvoice(r.Context(), &req, createdBy)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.CreatedResponse(w, r, id, "Invoice created successfully")
	}
}

func handleGetInvoice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := arMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid invoice ID"))
			return
		}

		invoice, err := svc.GetInvoice(r.Context(), id)
		if err != nil {
			helper.NotFoundResponse(w, r)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, invoice)
	}
}

func handleGetInvoiceByNumber() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := arMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		number := chi.URLParam(r, "number")
		invoice, err := svc.GetInvoiceByNumber(r.Context(), number)
		if err != nil {
			helper.NotFoundResponse(w, r)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, invoice)
	}
}

func handleListInvoices() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := arMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		filters := models.ARInvoiceListFilters{
			DateFrom: r.URL.Query().Get("date_from"),
			DateTo:   r.URL.Query().Get("date_to"),
			Overdue:  r.URL.Query().Get("overdue") == "true",
			Page:     1,
			PageSize: 20,
		}

		if page, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil {
			filters.Page = page
		}
		if pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size")); err == nil {
			filters.PageSize = pageSize
		}
		if customerID, err := strconv.Atoi(r.URL.Query().Get("customer_id")); err == nil {
			filters.CustomerID = &customerID
		}
		if status := r.URL.Query().Get("status"); status != "" {
			s := models.ARInvoiceStatus(status)
			filters.Status = &s
		}

		invoices, total, err := svc.ListInvoices(r.Context(), &filters)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		totalPages := int(total) / filters.PageSize
		if int(total)%filters.PageSize > 0 {
			totalPages++
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{
			"data": invoices,
			"pagination": helper.Envelope{
				"page":        filters.Page,
				"page_size":   filters.PageSize,
				"total_items": total,
				"total_pages": totalPages,
			},
		})
	}
}

func handlePostInvoice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := arMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid invoice ID"))
			return
		}

		postedBy := 1 // TODO: Get from auth context

		err = svc.PostInvoice(r.Context(), id, postedBy)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Invoice posted successfully"})
	}
}

func handleVoidInvoice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := arMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid invoice ID"))
			return
		}

		err = svc.VoidInvoice(r.Context(), id)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Invoice voided successfully"})
	}
}

func handleCreateFromOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := arMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		orderID, err := strconv.Atoi(chi.URLParam(r, "orderId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid order ID"))
			return
		}

		createdBy := 1 // TODO: Get from auth context

		id, err := svc.CreateFromOrder(r.Context(), orderID, createdBy)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.CreatedResponse(w, r, id, "Invoice created from order successfully")
	}
}

// ===========================================
// Payment Handlers
// ===========================================

func handleCreatePayment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := arMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		var req models.CreateARPaymentRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		v := models.NewValidator()
		models.ValidateARPayment(v, &req)
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		receivedBy := 1 // TODO: Get from auth context

		id, err := svc.CreatePayment(r.Context(), &req, receivedBy)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.CreatedResponse(w, r, id, "Payment recorded successfully")
	}
}

func handleGetPayment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := arMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid payment ID"))
			return
		}

		payment, err := svc.GetPayment(r.Context(), id)
		if err != nil {
			helper.NotFoundResponse(w, r)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, payment)
	}
}

func handleListPayments() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := arMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		var customerID *int
		if cid, err := strconv.Atoi(r.URL.Query().Get("customer_id")); err == nil {
			customerID = &cid
		}

		limit := 100
		if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && l > 0 {
			limit = l
		}

		payments, err := svc.ListPayments(r.Context(), customerID, limit)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, payments)
	}
}

// ===========================================
// Credit Management Handlers
// ===========================================

func handleGetCustomerCredit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := arMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		customerID, err := strconv.Atoi(chi.URLParam(r, "customerId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid customer ID"))
			return
		}

		credit, err := svc.GetCustomerCredit(r.Context(), customerID)
		if err != nil {
			helper.NotFoundResponse(w, r)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, credit)
	}
}

func handleCheckCredit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := arMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		customerID, err := strconv.Atoi(chi.URLParam(r, "customerId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid customer ID"))
			return
		}

		amount, err := strconv.ParseFloat(r.URL.Query().Get("amount"), 64)
		if err != nil || amount <= 0 {
			helper.BadRequestResponse(w, r, errors.New("valid amount is required"))
			return
		}

		available, creditAvailable, err := svc.CheckCreditAvailable(r.Context(), customerID, amount)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{
			"is_available":     available,
			"available_credit": creditAvailable,
			"requested_amount": amount,
		})
	}
}

func handleUpdateCreditLimit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := arMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		customerID, err := strconv.Atoi(chi.URLParam(r, "customerId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid customer ID"))
			return
		}

		var req struct {
			CreditLimit float64 `json:"credit_limit"`
		}
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		err = svc.UpdateCreditLimit(r.Context(), customerID, req.CreditLimit)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Credit limit updated"})
	}
}

// ===========================================
// Aging & Statement Handlers
// ===========================================

func handleGetCustomerAging() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := arMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		customerID, err := strconv.Atoi(chi.URLParam(r, "customerId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid customer ID"))
			return
		}

		aging, err := svc.GetCustomerAging(r.Context(), customerID)
		if err != nil {
			helper.NotFoundResponse(w, r)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, aging)
	}
}

func handleGetAgingReport() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := arMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		report, err := svc.GetAgingReport(r.Context())
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, report)
	}
}

func handleGetStatement() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := arMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		customerID, err := strconv.Atoi(chi.URLParam(r, "customerId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid customer ID"))
			return
		}

		fromDate := r.URL.Query().Get("from_date")
		toDate := r.URL.Query().Get("to_date")
		if fromDate == "" || toDate == "" {
			helper.BadRequestResponse(w, r, errors.New("from_date and to_date are required"))
			return
		}

		statement, err := svc.GetStatement(r.Context(), customerID, fromDate, toDate)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, statement)
	}
}

func handleGetOverdueInvoices() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := arMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		daysOverdue := 1
		if d, err := strconv.Atoi(r.URL.Query().Get("days")); err == nil && d > 0 {
			daysOverdue = d
		}

		invoices, err := svc.GetOverdueInvoices(r.Context(), daysOverdue)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, invoices)
	}
}
