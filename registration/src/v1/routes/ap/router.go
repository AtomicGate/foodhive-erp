package ap

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/anas-dev-92/FoodHive/core/auth"
	"github.com/anas-dev-92/FoodHive/core/jwt"
	"github.com/anas-dev-92/FoodHive/core/postgres"
	apMiddleware "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/ap"
	authMiddleware "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/auth"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/models"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/utils/helper"
	"github.com/go-chi/chi/v5"
)

func Router(db postgres.Executor, jwtService jwt.JWTService, authService auth.AuthService) chi.Router {
	app := chi.NewRouter()

	// Inject AP service
	app.Use(apMiddleware.New(db))

	// Apply authentication
	app.Use(authMiddleware.Authenticate(jwtService))

	// ===========================================
	// Invoices
	// ===========================================
	app.With(authMiddleware.Authorize(jwtService)).Post("/invoices/create", handleCreateInvoice())
	app.With(authMiddleware.Authorize(jwtService)).Get("/invoices/get/{id}", handleGetInvoice())
	app.With(authMiddleware.Authorize(jwtService)).Get("/invoices/list", handleListInvoices())
	app.With(authMiddleware.Authorize(jwtService)).Post("/invoices/{id}/approve", handleApproveInvoice())
	app.With(authMiddleware.Authorize(jwtService)).Post("/invoices/{id}/void", handleVoidInvoice())
	app.With(authMiddleware.Authorize(jwtService)).Post("/invoices/from-receiving/{receivingId}", handleCreateFromReceiving())

	// ===========================================
	// Payments
	// ===========================================
	app.With(authMiddleware.Authorize(jwtService)).Post("/payments/create", handleCreatePayment())
	app.With(authMiddleware.Authorize(jwtService)).Get("/payments/get/{id}", handleGetPayment())
	app.With(authMiddleware.Authorize(jwtService)).Get("/payments/list", handleListPayments())
	app.With(authMiddleware.Authorize(jwtService)).Post("/payments/{id}/void", handleVoidPayment())

	// ===========================================
	// Vendor Balance & Aging
	// ===========================================
	app.With(authMiddleware.Authorize(jwtService)).Get("/balance/{vendorId}", handleGetVendorBalance())
	app.With(authMiddleware.Authorize(jwtService)).Get("/aging/{vendorId}", handleGetVendorAging())
	app.With(authMiddleware.Authorize(jwtService)).Get("/aging/report", handleGetAgingReport())

	// ===========================================
	// Due & Overdue
	// ===========================================
	app.With(authMiddleware.Authorize(jwtService)).Get("/due", handleGetDueInvoices())
	app.With(authMiddleware.Authorize(jwtService)).Get("/overdue", handleGetOverdueInvoices())

	return app
}

// ===========================================
// Invoice Handlers
// ===========================================

func handleCreateInvoice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := apMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		var req models.CreateAPInvoiceRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		v := models.NewValidator()
		models.ValidateAPInvoice(v, &req)
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

		helper.CreatedResponse(w, r, id, "AP Invoice created successfully")
	}
}

func handleGetInvoice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := apMiddleware.Instance(r.Context())
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

func handleListInvoices() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := apMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		filters := models.APInvoiceListFilters{
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
		if vendorID, err := strconv.Atoi(r.URL.Query().Get("vendor_id")); err == nil {
			filters.VendorID = &vendorID
		}
		if status := r.URL.Query().Get("status"); status != "" {
			s := models.APInvoiceStatus(status)
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

func handleApproveInvoice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := apMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid invoice ID"))
			return
		}

		approvedBy := 1 // TODO: Get from auth context

		err = svc.ApproveInvoice(r.Context(), id, approvedBy)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Invoice approved successfully"})
	}
}

func handleVoidInvoice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := apMiddleware.Instance(r.Context())
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

func handleCreateFromReceiving() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := apMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		receivingID, err := strconv.Atoi(chi.URLParam(r, "receivingId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid receiving ID"))
			return
		}

		createdBy := 1 // TODO: Get from auth context

		id, err := svc.CreateFromReceiving(r.Context(), receivingID, createdBy)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.CreatedResponse(w, r, id, "AP Invoice created from receiving successfully")
	}
}

// ===========================================
// Payment Handlers
// ===========================================

func handleCreatePayment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := apMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		var req models.CreateAPPaymentRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		v := models.NewValidator()
		models.ValidateAPPayment(v, &req)
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		preparedBy := 1 // TODO: Get from auth context

		id, err := svc.CreatePayment(r.Context(), &req, preparedBy)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.CreatedResponse(w, r, id, "Payment created successfully")
	}
}

func handleGetPayment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := apMiddleware.Instance(r.Context())
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
		svc, ok := apMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		var vendorID *int
		if vid, err := strconv.Atoi(r.URL.Query().Get("vendor_id")); err == nil {
			vendorID = &vid
		}

		limit := 100
		if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && l > 0 {
			limit = l
		}

		payments, err := svc.ListPayments(r.Context(), vendorID, limit)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, payments)
	}
}

func handleVoidPayment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := apMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid payment ID"))
			return
		}

		err = svc.VoidPayment(r.Context(), id)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Payment voided successfully"})
	}
}

// ===========================================
// Balance & Aging Handlers
// ===========================================

func handleGetVendorBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := apMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		vendorID, err := strconv.Atoi(chi.URLParam(r, "vendorId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid vendor ID"))
			return
		}

		balance, err := svc.GetVendorBalance(r.Context(), vendorID)
		if err != nil {
			helper.NotFoundResponse(w, r)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, balance)
	}
}

func handleGetVendorAging() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := apMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		vendorID, err := strconv.Atoi(chi.URLParam(r, "vendorId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid vendor ID"))
			return
		}

		aging, err := svc.GetVendorAging(r.Context(), vendorID)
		if err != nil {
			helper.NotFoundResponse(w, r)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, aging)
	}
}

func handleGetAgingReport() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := apMiddleware.Instance(r.Context())
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

// ===========================================
// Due & Overdue Handlers
// ===========================================

func handleGetDueInvoices() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := apMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		withinDays := 7
		if d, err := strconv.Atoi(r.URL.Query().Get("days")); err == nil && d > 0 {
			withinDays = d
		}

		invoices, err := svc.GetDueInvoices(r.Context(), withinDays)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, invoices)
	}
}

func handleGetOverdueInvoices() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := apMiddleware.Instance(r.Context())
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
