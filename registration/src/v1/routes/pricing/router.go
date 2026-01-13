package pricing

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/anas-dev-92/FoodHive/core/auth"
	"github.com/anas-dev-92/FoodHive/core/jwt"
	"github.com/anas-dev-92/FoodHive/core/postgres"
	authMiddleware "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/auth"
	pricingMiddleware "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/pricing"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/models"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/utils/helper"
	"github.com/go-chi/chi/v5"
)

func Router(db postgres.Executor, jwtService jwt.JWTService, authService auth.AuthService) chi.Router {
	app := chi.NewRouter()

	// Inject pricing service
	app.Use(pricingMiddleware.New(db))

	// Apply authentication
	app.Use(authMiddleware.Authenticate(jwtService))

	// ===========================================
	// Price Lookup
	// ===========================================
	app.With(authMiddleware.Authorize(jwtService)).Post("/lookup", handlePriceLookup())
	app.With(authMiddleware.Authorize(jwtService)).Post("/lookup/batch", handleBatchPriceLookup())
	app.With(authMiddleware.Authorize(jwtService)).Get("/check-margin", handleCheckMargin())

	// ===========================================
	// Product Prices (Base & Level)
	// ===========================================
	app.With(authMiddleware.Authorize(jwtService)).Post("/product/set", handleSetProductPrice())
	app.With(authMiddleware.Authorize(jwtService)).Get("/product/{productId}", handleGetProductPrices())
	app.With(authMiddleware.Authorize(jwtService)).Delete("/product/price/{id}", handleDeleteProductPrice())

	// ===========================================
	// Customer Prices
	// ===========================================
	app.With(authMiddleware.Authorize(jwtService)).Post("/customer/set", handleSetCustomerPrice())
	app.With(authMiddleware.Authorize(jwtService)).Get("/customer/{customerId}", handleGetCustomerPrices())
	app.With(authMiddleware.Authorize(jwtService)).Delete("/customer/price/{id}", handleDeleteCustomerPrice())

	// ===========================================
	// Contracts
	// ===========================================
	app.With(authMiddleware.Authorize(jwtService)).Post("/contracts/create", handleCreateContract())
	app.With(authMiddleware.Authorize(jwtService)).Get("/contracts/get/{id}", handleGetContract())
	app.With(authMiddleware.Authorize(jwtService)).Get("/contracts/list", handleListContracts())
	app.With(authMiddleware.Authorize(jwtService)).Post("/contracts/{id}/deactivate", handleDeactivateContract())

	// ===========================================
	// Promotions
	// ===========================================
	app.With(authMiddleware.Authorize(jwtService)).Post("/promotions/create", handleCreatePromotion())
	app.With(authMiddleware.Authorize(jwtService)).Get("/promotions/get/{id}", handleGetPromotion())
	app.With(authMiddleware.Authorize(jwtService)).Get("/promotions/list", handleListPromotions())
	app.With(authMiddleware.Authorize(jwtService)).Post("/promotions/{id}/deactivate", handleDeactivatePromotion())

	// ===========================================
	// Product Costs
	// ===========================================
	app.With(authMiddleware.Authorize(jwtService)).Post("/costs/update", handleUpdateProductCost())
	app.With(authMiddleware.Authorize(jwtService)).Get("/costs/{productId}", handleGetProductCost())

	// ===========================================
	// Mass Price Maintenance
	// ===========================================
	app.With(authMiddleware.Authorize(jwtService)).Post("/mass-update", handleMassPriceUpdate())

	// ===========================================
	// Price Lists
	// ===========================================
	app.With(authMiddleware.Authorize(jwtService)).Get("/list", handleGetPriceList())

	return app
}

// ===========================================
// Price Lookup Handlers
// ===========================================

func handlePriceLookup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := pricingMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		var req models.PriceLookupRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		if req.ProductID <= 0 {
			helper.BadRequestResponse(w, r, errors.New("product_id is required"))
			return
		}
		if req.Quantity <= 0 {
			req.Quantity = 1
		}

		result, err := svc.GetPrice(r.Context(), &req)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, result)
	}
}

func handleBatchPriceLookup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := pricingMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		var req struct {
			ProductIDs []int `json:"product_ids"`
			CustomerID *int  `json:"customer_id,omitempty"`
		}
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		if len(req.ProductIDs) == 0 {
			helper.BadRequestResponse(w, r, errors.New("product_ids is required"))
			return
		}

		results, err := svc.GetPricesForProducts(r.Context(), req.ProductIDs, req.CustomerID)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, results)
	}
}

func handleCheckMargin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := pricingMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		productID, err := strconv.Atoi(r.URL.Query().Get("product_id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("product_id is required"))
			return
		}

		price, err := strconv.ParseFloat(r.URL.Query().Get("price"), 64)
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("price is required"))
			return
		}

		isBelowCost, cost, err := svc.CheckBelowCost(r.Context(), productID, price)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{
			"is_below_cost": isBelowCost,
			"cost":          cost,
			"price":         price,
			"margin":        price - cost,
		})
	}
}

// ===========================================
// Product Price Handlers
// ===========================================

func handleSetProductPrice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := pricingMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		var req models.SetProductPriceRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		v := models.NewValidator()
		models.ValidateProductPrice(v, &req)
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		createdBy := 1 // TODO: Get from auth context

		id, err := svc.SetProductPrice(r.Context(), &req, createdBy)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.CreatedResponse(w, r, id, "Product price set successfully")
	}
}

func handleGetProductPrices() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := pricingMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		productID, err := strconv.Atoi(chi.URLParam(r, "productId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid product ID"))
			return
		}

		prices, err := svc.GetProductPrices(r.Context(), productID)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, prices)
	}
}

func handleDeleteProductPrice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := pricingMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid price ID"))
			return
		}

		err = svc.DeleteProductPrice(r.Context(), id)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Price deleted successfully"})
	}
}

// ===========================================
// Customer Price Handlers
// ===========================================

func handleSetCustomerPrice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := pricingMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		var req models.SetCustomerPriceRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		v := models.NewValidator()
		models.ValidateCustomerPrice(v, &req)
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		createdBy := 1 // TODO: Get from auth context

		id, err := svc.SetCustomerPrice(r.Context(), &req, createdBy)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.CreatedResponse(w, r, id, "Customer price set successfully")
	}
}

func handleGetCustomerPrices() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := pricingMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		customerID, err := strconv.Atoi(chi.URLParam(r, "customerId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid customer ID"))
			return
		}

		prices, err := svc.GetCustomerPrices(r.Context(), customerID)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, prices)
	}
}

func handleDeleteCustomerPrice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := pricingMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid price ID"))
			return
		}

		err = svc.DeleteCustomerPrice(r.Context(), id)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Customer price deleted successfully"})
	}
}

// ===========================================
// Contract Handlers
// ===========================================

func handleCreateContract() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := pricingMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		var req models.CreatePriceContractRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		v := models.NewValidator()
		models.ValidatePriceContract(v, &req)
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		createdBy := 1 // TODO: Get from auth context

		id, err := svc.CreateContract(r.Context(), &req, createdBy)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.CreatedResponse(w, r, id, "Contract created successfully")
	}
}

func handleGetContract() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := pricingMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid contract ID"))
			return
		}

		contract, err := svc.GetContract(r.Context(), id)
		if err != nil {
			helper.NotFoundResponse(w, r)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, contract)
	}
}

func handleListContracts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := pricingMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		var customerID *int
		if cid, err := strconv.Atoi(r.URL.Query().Get("customer_id")); err == nil {
			customerID = &cid
		}

		activeOnly := r.URL.Query().Get("active_only") == "true"

		contracts, err := svc.ListContracts(r.Context(), customerID, activeOnly)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, contracts)
	}
}

func handleDeactivateContract() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := pricingMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid contract ID"))
			return
		}

		err = svc.DeactivateContract(r.Context(), id)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Contract deactivated"})
	}
}

// ===========================================
// Promotion Handlers
// ===========================================

func handleCreatePromotion() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := pricingMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		var req models.CreatePromotionRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		v := models.NewValidator()
		models.ValidatePromotion(v, &req)
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		createdBy := 1 // TODO: Get from auth context

		id, err := svc.CreatePromotion(r.Context(), &req, createdBy)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.CreatedResponse(w, r, id, "Promotion created successfully")
	}
}

func handleGetPromotion() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := pricingMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid promotion ID"))
			return
		}

		promotion, err := svc.GetPromotion(r.Context(), id)
		if err != nil {
			helper.NotFoundResponse(w, r)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, promotion)
	}
}

func handleListPromotions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := pricingMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		activeOnly := r.URL.Query().Get("active_only") == "true"

		promotions, err := svc.ListPromotions(r.Context(), activeOnly)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, promotions)
	}
}

func handleDeactivatePromotion() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := pricingMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid promotion ID"))
			return
		}

		err = svc.DeactivatePromotion(r.Context(), id)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Promotion deactivated"})
	}
}

// ===========================================
// Cost Handlers
// ===========================================

func handleUpdateProductCost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := pricingMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		var req models.UpdateProductCostRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		v := models.NewValidator()
		v.Check(req.ProductID > 0, "product_id", "Product is required")
		v.Check(req.Cost >= 0, "cost", "Cost must be non-negative")
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		updatedBy := 1 // TODO: Get from auth context

		err := svc.UpdateProductCost(r.Context(), &req, updatedBy)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Product cost updated successfully"})
	}
}

func handleGetProductCost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := pricingMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		productID, err := strconv.Atoi(chi.URLParam(r, "productId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid product ID"))
			return
		}

		cost, err := svc.GetProductCost(r.Context(), productID)
		if err != nil {
			helper.NotFoundResponse(w, r)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, cost)
	}
}

// ===========================================
// Mass Update Handler
// ===========================================

func handleMassPriceUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := pricingMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		var req models.MassPriceUpdateRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		v := models.NewValidator()
		models.ValidateMassPriceUpdate(v, &req)
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		createdBy := 1 // TODO: Get from auth context

		count, err := svc.MassPriceUpdate(r.Context(), &req, createdBy)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{
			"message":          "Mass price update completed",
			"products_updated": count,
		})
	}
}

// ===========================================
// Price List Handler
// ===========================================

func handleGetPriceList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := pricingMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		filters := models.PriceListFilters{
			EffectiveDate: r.URL.Query().Get("effective_date"),
			IncludeCost:   r.URL.Query().Get("include_cost") == "true",
			Page:          1,
			PageSize:      50,
		}

		if page, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil {
			filters.Page = page
		}
		if pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size")); err == nil {
			filters.PageSize = pageSize
		}
		if categoryID, err := strconv.Atoi(r.URL.Query().Get("category_id")); err == nil {
			filters.CategoryID = &categoryID
		}

		items, total, err := svc.GetPriceList(r.Context(), &filters)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		totalPages := int(total) / filters.PageSize
		if int(total)%filters.PageSize > 0 {
			totalPages++
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{
			"data": items,
			"pagination": helper.Envelope{
				"page":        filters.Page,
				"page_size":   filters.PageSize,
				"total_items": total,
				"total_pages": totalPages,
			},
		})
	}
}
