package inventory

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/anas-dev-92/FoodHive/core/auth"
	"github.com/anas-dev-92/FoodHive/core/jwt"
	"github.com/anas-dev-92/FoodHive/core/postgres"
	authMiddleware "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/auth"
	inventoryMiddleware "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/inventory"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/models"
	inventoryService "github.com/anas-dev-92/FoodHive/registration/src/v1/services/inventory"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/utils/helper"
	"github.com/go-chi/chi/v5"
)

func Router(db postgres.Executor, jwtService jwt.JWTService, authService auth.AuthService) chi.Router {
	app := chi.NewRouter()

	// Inject inventory service
	app.Use(inventoryMiddleware.New(db))

	// Apply authentication
	app.Use(authMiddleware.Authenticate(jwtService))

	// ===========================================
	// Inventory Query Routes
	// ===========================================
	app.With(authMiddleware.Authorize(jwtService)).Get("/get/{id}", handleGetByID())
	app.With(authMiddleware.Authorize(jwtService)).Get("/product/{productId}", handleGetByProduct())
	app.With(authMiddleware.Authorize(jwtService)).Get("/warehouse/{warehouseId}", handleGetByWarehouse())
	app.With(authMiddleware.Authorize(jwtService)).Get("/lot/{lotNumber}", handleGetByLot())
	app.With(authMiddleware.Authorize(jwtService)).Get("/list", handleList())

	// ===========================================
	// Summary Routes
	// ===========================================
	app.With(authMiddleware.Authorize(jwtService)).Get("/summary/product/{productId}", handleGetProductSummary())
	app.With(authMiddleware.Authorize(jwtService)).Get("/expiring", handleGetExpiring())

	// ===========================================
	// Inventory Operations
	// ===========================================
	app.With(authMiddleware.Authorize(jwtService)).Post("/receive", handleReceive())
	app.With(authMiddleware.Authorize(jwtService)).Post("/adjust", handleAdjust())
	app.With(authMiddleware.Authorize(jwtService)).Post("/transfer", handleTransfer())

	// ===========================================
	// Transaction History
	// ===========================================
	app.With(authMiddleware.Authorize(jwtService)).Get("/transactions", handleGetTransactions())

	return app
}

// ===========================================
// Query Handlers
// ===========================================

func handleGetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := inventoryMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid inventory ID"))
			return
		}

		inv, err := svc.GetByID(r.Context(), id)
		if err != nil {
			helper.NotFoundResponse(w, r)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, inv)
	}
}

func handleGetByProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := inventoryMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		productID, err := strconv.Atoi(chi.URLParam(r, "productId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid product ID"))
			return
		}

		inventory, err := svc.GetByProduct(r.Context(), productID)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, inventory)
	}
}

func handleGetByWarehouse() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := inventoryMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		warehouseID, err := strconv.Atoi(chi.URLParam(r, "warehouseId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid warehouse ID"))
			return
		}

		inventory, err := svc.GetByWarehouse(r.Context(), warehouseID)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, inventory)
	}
}

func handleGetByLot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := inventoryMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		lotNumber := chi.URLParam(r, "lotNumber")
		inventory, err := svc.GetByLot(r.Context(), lotNumber)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, inventory)
	}
}

func handleList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := inventoryMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		filters := models.InventoryInquiryRequest{
			LotNumber: r.URL.Query().Get("lot_number"),
			Page:      1,
			PageSize:  20,
		}

		if page, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil {
			filters.Page = page
		}
		if pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size")); err == nil {
			filters.PageSize = pageSize
		}
		if productID, err := strconv.Atoi(r.URL.Query().Get("product_id")); err == nil {
			filters.ProductID = &productID
		}
		if warehouseID, err := strconv.Atoi(r.URL.Query().Get("warehouse_id")); err == nil {
			filters.WarehouseID = &warehouseID
		}
		if r.URL.Query().Get("show_expiring") == "true" {
			filters.ShowExpiring = true
			if days, err := strconv.Atoi(r.URL.Query().Get("days_to_expiry")); err == nil {
				filters.DaysToExpiry = days
			} else {
				filters.DaysToExpiry = 30
			}
		}

		inventory, total, err := svc.List(r.Context(), &filters)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		totalPages := int(total) / filters.PageSize
		if int(total)%filters.PageSize > 0 {
			totalPages++
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{
			"data": inventory,
			"pagination": helper.Envelope{
				"page":        filters.Page,
				"page_size":   filters.PageSize,
				"total_items": total,
				"total_pages": totalPages,
			},
		})
	}
}

// ===========================================
// Summary Handlers
// ===========================================

func handleGetProductSummary() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := inventoryMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		productID, err := strconv.Atoi(chi.URLParam(r, "productId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid product ID"))
			return
		}

		summary, err := svc.GetProductSummary(r.Context(), productID)
		if err != nil {
			helper.NotFoundResponse(w, r)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, summary)
	}
}

func handleGetExpiring() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := inventoryMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		daysToExpiry := 30
		if days, err := strconv.Atoi(r.URL.Query().Get("days")); err == nil && days > 0 {
			daysToExpiry = days
		}

		var warehouseID *int
		if wid, err := strconv.Atoi(r.URL.Query().Get("warehouse_id")); err == nil {
			warehouseID = &wid
		}

		inventory, err := svc.GetExpiringInventory(r.Context(), daysToExpiry, warehouseID)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, inventory)
	}
}

// ===========================================
// Operation Handlers
// ===========================================

func handleReceive() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := inventoryMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		var req inventoryService.ReceiveRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		// Validate
		v := models.NewValidator()
		v.Check(req.ProductID > 0, "product_id", "Product ID is required")
		v.Check(req.WarehouseID > 0, "warehouse_id", "Warehouse ID is required")
		v.Check(req.Quantity > 0, "quantity", "Quantity must be positive")
		v.Check(req.UnitCost >= 0, "unit_cost", "Unit cost must be 0 or greater")
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		// Get user ID from context (assuming it's set by auth middleware)
		createdBy := 1 // TODO: Get from auth context

		id, err := svc.Receive(r.Context(), &req, createdBy)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.CreatedResponse(w, r, id, "Inventory received successfully")
	}
}

func handleAdjust() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := inventoryMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		var req models.AdjustInventoryRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		v := models.NewValidator()
		models.ValidateAdjustInventory(v, &req)
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		createdBy := 1 // TODO: Get from auth context

		err := svc.Adjust(r.Context(), &req, createdBy)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Inventory adjusted successfully"})
	}
}

func handleTransfer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := inventoryMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		var req models.TransferInventoryRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		v := models.NewValidator()
		models.ValidateTransferInventory(v, &req)
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		createdBy := 1 // TODO: Get from auth context

		err := svc.Transfer(r.Context(), &req, createdBy)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Inventory transferred successfully"})
	}
}

// ===========================================
// Transaction History Handler
// ===========================================

func handleGetTransactions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := inventoryMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		var productID, warehouseID *int
		if pid, err := strconv.Atoi(r.URL.Query().Get("product_id")); err == nil {
			productID = &pid
		}
		if wid, err := strconv.Atoi(r.URL.Query().Get("warehouse_id")); err == nil {
			warehouseID = &wid
		}

		limit := 100
		if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && l > 0 {
			limit = l
		}

		transactions, err := svc.GetTransactions(r.Context(), productID, warehouseID, limit)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, transactions)
	}
}
