package customer

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/anas-dev-92/FoodHive/core/auth"
	"github.com/anas-dev-92/FoodHive/core/jwt"
	"github.com/anas-dev-92/FoodHive/core/postgres"
	authMiddleware "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/auth"
	customerMiddleware "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/customer"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/models"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/utils/helper"
	"github.com/go-chi/chi/v5"
)

func Router(db postgres.Executor, jwtService jwt.JWTService, authService auth.AuthService) chi.Router {
	r := chi.NewRouter()

	// Inject customer service via middleware
	r.Use(customerMiddleware.New(db))

	// Apply authentication middleware globally
	r.Use(authMiddleware.Authenticate(jwtService))

	// Routes with authorization
	r.With(authMiddleware.Authorize(jwtService)).Post("/create", handleCreate())
	r.With(authMiddleware.Authorize(jwtService)).Get("/get/{id}", handleGetByID())
	r.With(authMiddleware.Authorize(jwtService)).Get("/code/{code}", handleGetByCode())
	r.With(authMiddleware.Authorize(jwtService)).Put("/update/{id}", handleUpdate())
	r.With(authMiddleware.Authorize(jwtService)).Delete("/delete/{id}", handleDelete())
	r.With(authMiddleware.Authorize(jwtService)).Get("/list", handleList())

	// Order Guide routes
	r.With(authMiddleware.Authorize(jwtService)).Get("/{customerId}/order-guide", handleGetOrderGuide())

	// Ship-to routes
	r.With(authMiddleware.Authorize(jwtService)).Post("/{customerId}/ship-to", handleAddShipTo())

	return r
}

func handleCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := customerMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		var req models.CreateCustomerRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		// Validate
		v := helper.New()
		v.Check(req.CustomerCode != "", "customer_code", "Customer code is required")
		v.Check(req.Name != "", "name", "Name is required")
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		// Get user ID from context
		userID, _ := authMiddleware.GetUserID(r.Context())

		id, err := svc.Create(r.Context(), req, userID)
		if err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		helper.CreatedResponse(w, r, id, "Customer created successfully")
	}
}

func handleGetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := customerMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid customer ID"))
			return
		}

		customer, err := svc.GetByID(r.Context(), id)
		if err != nil {
			helper.NotFoundResponse(w, r)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, customer)
	}
}

func handleGetByCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := customerMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		code := chi.URLParam(r, "code")
		if code == "" {
			helper.BadRequestResponse(w, r, errors.New("customer code is required"))
			return
		}

		customer, err := svc.GetByCode(r.Context(), code)
		if err != nil {
			helper.NotFoundResponse(w, r)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, customer)
	}
}

func handleUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := customerMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid customer ID"))
			return
		}

		var req models.UpdateCustomerRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		err = svc.Update(r.Context(), id, req)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Customer updated successfully"})
	}
}

func handleDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := customerMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid customer ID"))
			return
		}

		err = svc.Delete(r.Context(), id)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Customer deleted successfully"})
	}
}

func handleList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := customerMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		// Parse query params
		filters := models.CustomerListFilters{
			Search:   r.URL.Query().Get("search"),
			Page:     1,
			PageSize: 20,
		}

		if page, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && page > 0 {
			filters.Page = page
		}
		if pageSize, err := strconv.Atoi(r.URL.Query().Get("per_page")); err == nil && pageSize > 0 {
			filters.PageSize = pageSize
		}
		if salesRepID, err := strconv.Atoi(r.URL.Query().Get("sales_rep_id")); err == nil {
			filters.SalesRepID = &salesRepID
		}
		if isActive := r.URL.Query().Get("is_active"); isActive != "" {
			active := isActive == "true"
			filters.IsActive = &active
		}

		customers, total, err := svc.List(r.Context(), filters)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		totalPages := int(total) / filters.PageSize
		if int(total)%filters.PageSize > 0 {
			totalPages++
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{
			"data": customers,
			"pagination": helper.Envelope{
				"page":        filters.Page,
				"page_size":   filters.PageSize,
				"total_items": total,
				"total_pages": totalPages,
			},
		})
	}
}

func handleGetOrderGuide() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := customerMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		customerID, err := strconv.Atoi(chi.URLParam(r, "customerId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid customer ID"))
			return
		}

		guides, err := svc.GetOrderGuide(r.Context(), customerID)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, guides)
	}
}

func handleAddShipTo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := customerMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		customerID, err := strconv.Atoi(chi.URLParam(r, "customerId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid customer ID"))
			return
		}

		var shipTo models.CustomerShipTo
		if err := helper.ReadJSON(w, r, &shipTo); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		// Validate
		v := helper.New()
		v.Check(shipTo.ShipToCode != "", "ship_to_code", "Ship-to code is required")
		v.Check(shipTo.Name != "", "name", "Name is required")
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		id, err := svc.AddShipTo(r.Context(), customerID, shipTo)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.CreatedResponse(w, r, id, "Ship-to address added successfully")
	}
}
