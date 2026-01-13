package purchase_order

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/anas-dev-92/FoodHive/core/auth"
	"github.com/anas-dev-92/FoodHive/core/jwt"
	"github.com/anas-dev-92/FoodHive/core/postgres"
	authMiddleware "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/auth"
	poMiddleware "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/purchase_order"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/models"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/utils/helper"
	"github.com/go-chi/chi/v5"
)

func Router(db postgres.Executor, jwtService jwt.JWTService, authService auth.AuthService) chi.Router {
	app := chi.NewRouter()

	// Inject purchase order service
	app.Use(poMiddleware.New(db))

	// Apply authentication
	app.Use(authMiddleware.Authenticate(jwtService))

	// ===========================================
	// Purchase Order Routes
	// ===========================================
	app.With(authMiddleware.Authorize(jwtService)).Post("/create", handleCreate())
	app.With(authMiddleware.Authorize(jwtService)).Get("/get/{id}", handleGetByID())
	app.With(authMiddleware.Authorize(jwtService)).Get("/number/{poNumber}", handleGetByPONumber())
	app.With(authMiddleware.Authorize(jwtService)).Put("/update/{id}", handleUpdate())
	app.With(authMiddleware.Authorize(jwtService)).Delete("/delete/{id}", handleDelete())
	app.With(authMiddleware.Authorize(jwtService)).Get("/list", handleList())
	app.With(authMiddleware.Authorize(jwtService)).Post("/submit/{id}", handleSubmit())
	app.With(authMiddleware.Authorize(jwtService)).Post("/cancel/{id}", handleCancel())

	// ===========================================
	// Purchase Order Lines Routes
	// ===========================================
	app.With(authMiddleware.Authorize(jwtService)).Post("/{poId}/lines", handleAddLine())
	app.With(authMiddleware.Authorize(jwtService)).Put("/lines/{lineId}", handleUpdateLine())
	app.With(authMiddleware.Authorize(jwtService)).Delete("/lines/{lineId}", handleDeleteLine())

	// ===========================================
	// Receiving Routes
	// ===========================================
	app.With(authMiddleware.Authorize(jwtService)).Post("/receive", handleCreateReceiving())
	app.With(authMiddleware.Authorize(jwtService)).Get("/receiving/{id}", handleGetReceiving())
	app.With(authMiddleware.Authorize(jwtService)).Get("/receivings", handleListReceivings())

	return app
}

// ===========================================
// Purchase Order Handlers
// ===========================================

func handleCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := poMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		var req models.CreatePurchaseOrderRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		v := models.NewValidator()
		models.ValidatePurchaseOrder(v, &req)
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

		helper.CreatedResponse(w, r, id, "Purchase order created successfully")
	}
}

func handleGetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := poMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid purchase order ID"))
			return
		}

		po, err := svc.GetByID(r.Context(), id)
		if err != nil {
			helper.NotFoundResponse(w, r)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, po)
	}
}

func handleGetByPONumber() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := poMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		poNumber := chi.URLParam(r, "poNumber")
		po, err := svc.GetByPONumber(r.Context(), poNumber)
		if err != nil {
			helper.NotFoundResponse(w, r)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, po)
	}
}

func handleUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := poMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid purchase order ID"))
			return
		}

		var req models.UpdatePurchaseOrderRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		err = svc.Update(r.Context(), id, &req)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Purchase order updated successfully"})
	}
}

func handleDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := poMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid purchase order ID"))
			return
		}

		err = svc.Delete(r.Context(), id)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Purchase order deleted successfully"})
	}
}

func handleList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := poMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		filters := models.PurchaseOrderListFilters{
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
		if vendorID, err := strconv.Atoi(r.URL.Query().Get("vendor_id")); err == nil {
			filters.VendorID = &vendorID
		}
		if warehouseID, err := strconv.Atoi(r.URL.Query().Get("warehouse_id")); err == nil {
			filters.WarehouseID = &warehouseID
		}
		if status := r.URL.Query().Get("status"); status != "" {
			s := models.POStatus(status)
			filters.Status = &s
		}
		if buyerID, err := strconv.Atoi(r.URL.Query().Get("buyer_id")); err == nil {
			filters.BuyerID = &buyerID
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

func handleSubmit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := poMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid purchase order ID"))
			return
		}

		err = svc.Submit(r.Context(), id)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Purchase order submitted successfully"})
	}
}

func handleCancel() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := poMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid purchase order ID"))
			return
		}

		err = svc.Cancel(r.Context(), id)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Purchase order cancelled successfully"})
	}
}

// ===========================================
// Line Handlers
// ===========================================

func handleAddLine() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := poMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		poID, err := strconv.Atoi(chi.URLParam(r, "poId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid purchase order ID"))
			return
		}

		var req models.CreatePurchaseOrderLineRequest
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

		id, err := svc.AddLine(r.Context(), poID, &req)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.CreatedResponse(w, r, id, "Line added successfully")
	}
}

func handleUpdateLine() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := poMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		lineID, err := strconv.Atoi(chi.URLParam(r, "lineId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid line ID"))
			return
		}

		var req models.CreatePurchaseOrderLineRequest
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
		svc, ok := poMiddleware.Instance(r.Context())
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
// Receiving Handlers
// ===========================================

func handleCreateReceiving() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := poMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		var req models.CreateReceivingRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		v := models.NewValidator()
		models.ValidateReceiving(v, &req)
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		receivedBy := 1 // TODO: Get from auth context

		id, err := svc.CreateReceiving(r.Context(), &req, receivedBy)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.CreatedResponse(w, r, id, "Receiving created successfully")
	}
}

func handleGetReceiving() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := poMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid receiving ID"))
			return
		}

		receiving, err := svc.GetReceiving(r.Context(), id)
		if err != nil {
			helper.NotFoundResponse(w, r)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, receiving)
	}
}

func handleListReceivings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := poMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		var poID, warehouseID *int
		if pid, err := strconv.Atoi(r.URL.Query().Get("po_id")); err == nil {
			poID = &pid
		}
		if wid, err := strconv.Atoi(r.URL.Query().Get("warehouse_id")); err == nil {
			warehouseID = &wid
		}

		limit := 100
		if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && l > 0 {
			limit = l
		}

		receivings, err := svc.ListReceivings(r.Context(), poID, warehouseID, limit)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, receivings)
	}
}
