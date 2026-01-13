package vendor

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/anas-dev-92/FoodHive/core/auth"
	"github.com/anas-dev-92/FoodHive/core/jwt"
	"github.com/anas-dev-92/FoodHive/core/postgres"
	authMiddleware "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/auth"
	vendorMiddleware "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/vendor"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/models"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/utils/helper"
	"github.com/go-chi/chi/v5"
)

func Router(db postgres.Executor, jwtService jwt.JWTService, authService auth.AuthService) chi.Router {
	app := chi.NewRouter()

	// Inject vendor service
	app.Use(vendorMiddleware.New(db))

	// Apply authentication
	app.Use(authMiddleware.Authenticate(jwtService))

	// ===========================================
	// Vendor Routes
	// ===========================================
	app.With(authMiddleware.Authorize(jwtService)).Post("/create", handleCreate())
	app.With(authMiddleware.Authorize(jwtService)).Get("/get/{id}", handleGetByID())
	app.With(authMiddleware.Authorize(jwtService)).Get("/code/{code}", handleGetByCode())
	app.With(authMiddleware.Authorize(jwtService)).Put("/update/{id}", handleUpdate())
	app.With(authMiddleware.Authorize(jwtService)).Delete("/delete/{id}", handleDelete())
	app.With(authMiddleware.Authorize(jwtService)).Get("/list", handleList())

	// ===========================================
	// Vendor Product Routes
	// ===========================================
	app.Route("/products", func(r chi.Router) {
		r.With(authMiddleware.Authorize(jwtService)).Post("/add", handleAddProduct())
		r.With(authMiddleware.Authorize(jwtService)).Get("/vendor/{vendorId}", handleGetProducts())
		r.With(authMiddleware.Authorize(jwtService)).Get("/product/{productId}", handleGetProductsByProduct())
		r.With(authMiddleware.Authorize(jwtService)).Put("/update/{id}", handleUpdateProduct())
		r.With(authMiddleware.Authorize(jwtService)).Delete("/delete/{id}", handleDeleteProduct())
	})

	// ===========================================
	// Vendor Discount Routes
	// ===========================================
	app.Route("/discounts", func(r chi.Router) {
		r.With(authMiddleware.Authorize(jwtService)).Post("/add", handleAddDiscount())
		r.With(authMiddleware.Authorize(jwtService)).Get("/vendor/{vendorId}", handleGetDiscounts())
		r.With(authMiddleware.Authorize(jwtService)).Put("/update/{id}", handleUpdateDiscount())
		r.With(authMiddleware.Authorize(jwtService)).Delete("/delete/{id}", handleDeleteDiscount())
	})

	return app
}

// ===========================================
// Vendor Handlers
// ===========================================

func handleCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := vendorMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		var req models.CreateVendorRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		v := models.NewValidator()
		models.ValidateVendor(v, &req)
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		id, err := svc.Create(r.Context(), &req)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.CreatedResponse(w, r, id, "Vendor created successfully")
	}
}

func handleGetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := vendorMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid vendor ID"))
			return
		}

		vendor, err := svc.GetByID(r.Context(), id)
		if err != nil {
			helper.NotFoundResponse(w, r)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, vendor)
	}
}

func handleGetByCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := vendorMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		code := chi.URLParam(r, "code")
		vendor, err := svc.GetByCode(r.Context(), code)
		if err != nil {
			helper.NotFoundResponse(w, r)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, vendor)
	}
}

func handleUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := vendorMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid vendor ID"))
			return
		}

		var req models.UpdateVendorRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		err = svc.Update(r.Context(), id, &req)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Vendor updated successfully"})
	}
}

func handleDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := vendorMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid vendor ID"))
			return
		}

		err = svc.Delete(r.Context(), id)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Vendor deleted successfully"})
	}
}

func handleList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := vendorMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		filters := models.VendorListFilters{
			Search:   r.URL.Query().Get("search"),
			Page:     1,
			PageSize: 20,
		}

		if page, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil {
			filters.Page = page
		}
		if pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size")); err == nil {
			filters.PageSize = pageSize
		}
		if buyerID, err := strconv.Atoi(r.URL.Query().Get("buyer_id")); err == nil {
			filters.BuyerID = &buyerID
		}
		if isActive := r.URL.Query().Get("is_active"); isActive != "" {
			active := isActive == "true"
			filters.IsActive = &active
		}

		vendors, total, err := svc.List(r.Context(), &filters)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		totalPages := int(total) / filters.PageSize
		if int(total)%filters.PageSize > 0 {
			totalPages++
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{
			"data": vendors,
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
// Vendor Product Handlers
// ===========================================

func handleAddProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := vendorMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		var req models.VendorProduct
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		v := models.NewValidator()
		models.ValidateVendorProduct(v, &req)
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		id, err := svc.AddProduct(r.Context(), &req)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.CreatedResponse(w, r, id, "Vendor product added successfully")
	}
}

func handleGetProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := vendorMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		vendorID, err := strconv.Atoi(chi.URLParam(r, "vendorId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid vendor ID"))
			return
		}

		products, err := svc.GetProducts(r.Context(), vendorID)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, products)
	}
}

func handleGetProductsByProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := vendorMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		productID, err := strconv.Atoi(chi.URLParam(r, "productId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid product ID"))
			return
		}

		products, err := svc.GetProductsByProductID(r.Context(), productID)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, products)
	}
}

func handleUpdateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := vendorMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid vendor product ID"))
			return
		}

		var req models.VendorProduct
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		err = svc.UpdateProduct(r.Context(), id, &req)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Vendor product updated successfully"})
	}
}

func handleDeleteProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := vendorMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid vendor product ID"))
			return
		}

		err = svc.DeleteProduct(r.Context(), id)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Vendor product deleted successfully"})
	}
}

// ===========================================
// Vendor Discount Handlers
// ===========================================

func handleAddDiscount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := vendorMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		var req models.VendorDiscount
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		v := models.NewValidator()
		models.ValidateVendorDiscount(v, &req)
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		id, err := svc.AddDiscount(r.Context(), &req)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.CreatedResponse(w, r, id, "Vendor discount added successfully")
	}
}

func handleGetDiscounts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := vendorMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		vendorID, err := strconv.Atoi(chi.URLParam(r, "vendorId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid vendor ID"))
			return
		}

		discounts, err := svc.GetDiscounts(r.Context(), vendorID)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, discounts)
	}
}

func handleUpdateDiscount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := vendorMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid discount ID"))
			return
		}

		var req models.VendorDiscount
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		err = svc.UpdateDiscount(r.Context(), id, &req)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Vendor discount updated successfully"})
	}
}

func handleDeleteDiscount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := vendorMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid discount ID"))
			return
		}

		err = svc.DeleteDiscount(r.Context(), id)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Vendor discount deleted successfully"})
	}
}
