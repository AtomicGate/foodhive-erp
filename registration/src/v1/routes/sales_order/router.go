package sales_order

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/anas-dev-92/FoodHive/core/auth"
	"github.com/anas-dev-92/FoodHive/core/jwt"
	"github.com/anas-dev-92/FoodHive/core/postgres"
	authMiddleware "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/auth"
	soMiddleware "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/sales_order"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/models"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/utils/helper"
	"github.com/go-chi/chi/v5"
)

func Router(db postgres.Executor, jwtService jwt.JWTService, authService auth.AuthService) chi.Router {
	app := chi.NewRouter()

	// Inject sales order service
	app.Use(soMiddleware.New(db))

	// Apply authentication
	app.Use(authMiddleware.Authenticate(jwtService))

	// ===========================================
	// Sales Order Routes
	// ===========================================
	app.With(authMiddleware.Authorize(jwtService)).Post("/create", handleCreate())
	app.With(authMiddleware.Authorize(jwtService)).Get("/get/{id}", handleGetByID())
	app.With(authMiddleware.Authorize(jwtService)).Get("/number/{orderNumber}", handleGetByOrderNumber())
	app.With(authMiddleware.Authorize(jwtService)).Put("/update/{id}", handleUpdate())
	app.With(authMiddleware.Authorize(jwtService)).Delete("/delete/{id}", handleDelete())
	app.With(authMiddleware.Authorize(jwtService)).Get("/list", handleList())
	app.With(authMiddleware.Authorize(jwtService)).Post("/confirm/{id}", handleConfirm())
	app.With(authMiddleware.Authorize(jwtService)).Post("/cancel/{id}", handleCancel())
	app.With(authMiddleware.Authorize(jwtService)).Post("/ship/{id}", handleShip())

	// ===========================================
	// Sales Order Lines Routes
	// ===========================================
	app.With(authMiddleware.Authorize(jwtService)).Post("/{orderId}/lines", handleAddLine())
	app.With(authMiddleware.Authorize(jwtService)).Put("/lines/{lineId}", handleUpdateLine())
	app.With(authMiddleware.Authorize(jwtService)).Delete("/lines/{lineId}", handleDeleteLine())

	// ===========================================
	// Order Guide Routes
	// ===========================================
	app.With(authMiddleware.Authorize(jwtService)).Get("/order-guide/{customerId}", handleGetOrderGuide())

	// ===========================================
	// Lost Sales Routes
	// ===========================================
	app.With(authMiddleware.Authorize(jwtService)).Post("/lost-sale", handleRecordLostSale())
	app.With(authMiddleware.Authorize(jwtService)).Get("/lost-sales", handleGetLostSales())

	return app
}

// ===========================================
// Sales Order Handlers
// ===========================================

func handleCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := soMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		var req models.CreateSalesOrderRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		v := models.NewValidator()
		models.ValidateSalesOrder(v, &req)
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		createdBy := 1 // TODO: Get from auth context

		id, err := svc.Create(r.Context(), &req, createdBy)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.CreatedResponse(w, r, id, "Sales order created successfully")
	}
}

func handleGetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := soMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid sales order ID"))
			return
		}

		order, err := svc.GetByID(r.Context(), id)
		if err != nil {
			helper.NotFoundResponse(w, r)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, order)
	}
}

func handleGetByOrderNumber() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := soMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		orderNumber := chi.URLParam(r, "orderNumber")
		order, err := svc.GetByOrderNumber(r.Context(), orderNumber)
		if err != nil {
			helper.NotFoundResponse(w, r)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, order)
	}
}

func handleUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := soMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid sales order ID"))
			return
		}

		var req models.UpdateSalesOrderRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		err = svc.Update(r.Context(), id, &req)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Sales order updated successfully"})
	}
}

func handleDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := soMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid sales order ID"))
			return
		}

		err = svc.Delete(r.Context(), id)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Sales order deleted successfully"})
	}
}

func handleList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := soMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		filters := models.SalesOrderListFilters{
			DateFrom: r.URL.Query().Get("date_from"),
			DateTo:   r.URL.Query().Get("date_to"),
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
		if warehouseID, err := strconv.Atoi(r.URL.Query().Get("warehouse_id")); err == nil {
			filters.WarehouseID = &warehouseID
		}
		if routeID, err := strconv.Atoi(r.URL.Query().Get("route_id")); err == nil {
			filters.RouteID = &routeID
		}
		if salesRepID, err := strconv.Atoi(r.URL.Query().Get("sales_rep_id")); err == nil {
			filters.SalesRepID = &salesRepID
		}
		if status := r.URL.Query().Get("status"); status != "" {
			s := models.OrderStatus(status)
			filters.Status = &s
		}
		if orderType := r.URL.Query().Get("order_type"); orderType != "" {
			t := models.OrderType(orderType)
			filters.OrderType = &t
		}

		orders, total, err := svc.List(r.Context(), &filters)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		totalPages := int(total) / filters.PageSize
		if int(total)%filters.PageSize > 0 {
			totalPages++
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{
			"data": orders,
			"pagination": helper.Envelope{
				"page":        filters.Page,
				"page_size":   filters.PageSize,
				"total_items": total,
				"total_pages": totalPages,
			},
		})
	}
}

func handleConfirm() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := soMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid sales order ID"))
			return
		}

		err = svc.Confirm(r.Context(), id)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Sales order confirmed successfully"})
	}
}

func handleCancel() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := soMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid sales order ID"))
			return
		}

		err = svc.Cancel(r.Context(), id)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Sales order cancelled successfully"})
	}
}

func handleShip() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := soMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid sales order ID"))
			return
		}

		err = svc.Ship(r.Context(), id)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Sales order shipped successfully"})
	}
}

// ===========================================
// Line Handlers
// ===========================================

func handleAddLine() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := soMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		orderID, err := strconv.Atoi(chi.URLParam(r, "orderId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid sales order ID"))
			return
		}

		var req models.CreateSalesOrderLineRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		v := models.NewValidator()
		v.Check(req.ProductID > 0, "product_id", "Product ID is required")
		v.Check(req.Quantity > 0, "quantity", "Quantity must be positive")
		v.Check(req.UnitOfMeasure != "", "unit_of_measure", "Unit of measure is required")
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		id, err := svc.AddLine(r.Context(), orderID, &req)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.CreatedResponse(w, r, id, "Line added successfully")
	}
}

func handleUpdateLine() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := soMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		lineID, err := strconv.Atoi(chi.URLParam(r, "lineId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid line ID"))
			return
		}

		var req models.CreateSalesOrderLineRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		err = svc.UpdateLine(r.Context(), lineID, &req)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Line updated successfully"})
	}
}

func handleDeleteLine() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := soMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		lineID, err := strconv.Atoi(chi.URLParam(r, "lineId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid line ID"))
			return
		}

		err = svc.DeleteLine(r.Context(), lineID)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Line deleted successfully"})
	}
}

// ===========================================
// Order Guide Handler
// ===========================================

func handleGetOrderGuide() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := soMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		customerID, err := strconv.Atoi(chi.URLParam(r, "customerId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid customer ID"))
			return
		}

		warehouseID, err := strconv.Atoi(r.URL.Query().Get("warehouse_id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("warehouse_id is required"))
			return
		}

		entries, err := svc.GetOrderGuide(r.Context(), customerID, warehouseID)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, entries)
	}
}

// ===========================================
// Lost Sales Handlers
// ===========================================

func handleRecordLostSale() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := soMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		var req struct {
			OrderID           int     `json:"order_id"`
			ProductID         int     `json:"product_id"`
			QuantityRequested float64 `json:"quantity_requested"`
			QuantityAvailable float64 `json:"quantity_available"`
			Reason            string  `json:"reason"`
		}

		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		v := models.NewValidator()
		v.Check(req.OrderID > 0, "order_id", "Order ID is required")
		v.Check(req.ProductID > 0, "product_id", "Product ID is required")
		v.Check(req.QuantityRequested > 0, "quantity_requested", "Quantity requested must be positive")
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		err := svc.RecordLostSale(r.Context(), req.OrderID, req.ProductID, req.QuantityRequested, req.QuantityAvailable, req.Reason)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Lost sale recorded successfully"})
	}
}

func handleGetLostSales() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := soMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		var orderID *int
		if oid, err := strconv.Atoi(r.URL.Query().Get("order_id")); err == nil {
			orderID = &oid
		}

		limit := 100
		if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && l > 0 {
			limit = l
		}

		lostSales, err := svc.GetLostSales(r.Context(), orderID, limit)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, lostSales)
	}
}
