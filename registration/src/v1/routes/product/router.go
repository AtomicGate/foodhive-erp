package product

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/anas-dev-92/FoodHive/core/auth"
	"github.com/anas-dev-92/FoodHive/core/jwt"
	"github.com/anas-dev-92/FoodHive/core/postgres"
	authMiddleware "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/auth"
	productMiddleware "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/product"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/models"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/utils/helper"
	"github.com/go-chi/chi/v5"
)

func Router(db postgres.Executor, jwtService jwt.JWTService, authService auth.AuthService) chi.Router {
	app := chi.NewRouter()

	// Inject product service
	app.Use(productMiddleware.New(db))

	// Apply authentication
	app.Use(authMiddleware.Authenticate(jwtService))

	// ===========================================
	// Product Routes
	// ===========================================
	app.With(authMiddleware.Authorize(jwtService)).Post("/create", handleCreate())
	app.With(authMiddleware.Authorize(jwtService)).Get("/get/{id}", handleGetByID())
	app.With(authMiddleware.Authorize(jwtService)).Get("/sku/{sku}", handleGetBySKU())
	app.With(authMiddleware.Authorize(jwtService)).Get("/barcode/{barcode}", handleGetByBarcode())
	app.With(authMiddleware.Authorize(jwtService)).Put("/update/{id}", handleUpdate())
	app.With(authMiddleware.Authorize(jwtService)).Delete("/delete/{id}", handleDelete())
	app.With(authMiddleware.Authorize(jwtService)).Get("/list", handleList())

	// ===========================================
	// Category Routes
	// ===========================================
	app.Route("/categories", func(r chi.Router) {
		r.With(authMiddleware.Authorize(jwtService)).Post("/create", handleCreateCategory())
		r.With(authMiddleware.Authorize(jwtService)).Get("/get/{id}", handleGetCategory())
		r.With(authMiddleware.Authorize(jwtService)).Put("/update/{id}", handleUpdateCategory())
		r.With(authMiddleware.Authorize(jwtService)).Delete("/delete/{id}", handleDeleteCategory())
		r.With(authMiddleware.Authorize(jwtService)).Get("/list", handleListCategories())
	})

	// ===========================================
	// Unit Routes
	// ===========================================
	app.Route("/units", func(r chi.Router) {
		r.With(authMiddleware.Authorize(jwtService)).Post("/add", handleAddUnit())
		r.With(authMiddleware.Authorize(jwtService)).Get("/product/{productId}", handleGetUnits())
		r.With(authMiddleware.Authorize(jwtService)).Put("/update/{id}", handleUpdateUnit())
		r.With(authMiddleware.Authorize(jwtService)).Delete("/delete/{id}", handleDeleteUnit())
	})

	return app
}

// ===========================================
// Product Handlers
// ===========================================

func handleCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := productMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		var req models.CreateProductRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		// Validate
		v := models.NewValidator()
		models.ValidateProduct(v, &req)
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		id, err := svc.Create(r.Context(), &req)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.CreatedResponse(w, r, id, "Product created successfully")
	}
}

func handleGetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := productMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid product ID"))
			return
		}

		product, err := svc.GetByID(r.Context(), id)
		if err != nil {
			helper.NotFoundResponse(w, r)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, product)
	}
}

func handleGetBySKU() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := productMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		sku := chi.URLParam(r, "sku")
		product, err := svc.GetBySKU(r.Context(), sku)
		if err != nil {
			helper.NotFoundResponse(w, r)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, product)
	}
}

func handleGetByBarcode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := productMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		barcode := chi.URLParam(r, "barcode")
		product, err := svc.GetByBarcode(r.Context(), barcode)
		if err != nil {
			helper.NotFoundResponse(w, r)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, product)
	}
}

func handleUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := productMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid product ID"))
			return
		}

		var req models.UpdateProductRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		log.Printf("Product update request for ID %d: Name=%v, Description=%v, QCRequired=%v, IsActive=%v",
			id, req.Name, req.Description, req.QCRequired, req.IsActive)

		err = svc.Update(r.Context(), id, &req)
		if err != nil {
			log.Printf("Product update failed: %v", err)
			helper.ServerErrorResponse(w, r, err)
			return
		}

		log.Printf("Product %d updated successfully", id)
		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Product updated successfully"})
	}
}

func handleDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := productMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid product ID"))
			return
		}

		err = svc.Delete(r.Context(), id)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Product deleted successfully"})
	}
}

func handleList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := productMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		// Parse query params
		filters := models.ProductListFilters{
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
		if categoryID, err := strconv.Atoi(r.URL.Query().Get("category_id")); err == nil {
			filters.CategoryID = &categoryID
		}
		if isActive := r.URL.Query().Get("is_active"); isActive != "" {
			active := isActive == "true"
			filters.IsActive = &active
		}

		products, total, err := svc.List(r.Context(), &filters)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		totalPages := int(total) / filters.PageSize
		if int(total)%filters.PageSize > 0 {
			totalPages++
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{
			"data": products,
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
// Category Handlers
// ===========================================

func handleCreateCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := productMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		var req models.ProductCategory
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		// Validate
		v := models.NewValidator()
		v.Check(req.Name != "", "name", "Category name is required")
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		id, err := svc.CreateCategory(r.Context(), &req)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.CreatedResponse(w, r, id, "Category created successfully")
	}
}

func handleGetCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := productMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid category ID"))
			return
		}

		category, err := svc.GetCategoryByID(r.Context(), id)
		if err != nil {
			helper.NotFoundResponse(w, r)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, category)
	}
}

func handleUpdateCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := productMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid category ID"))
			return
		}

		var req models.ProductCategory
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		err = svc.UpdateCategory(r.Context(), id, &req)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Category updated successfully"})
	}
}

func handleDeleteCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := productMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid category ID"))
			return
		}

		err = svc.DeleteCategory(r.Context(), id)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Category deleted successfully"})
	}
}

func handleListCategories() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := productMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		categories, err := svc.ListCategories(r.Context())
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, categories)
	}
}

// ===========================================
// Unit Handlers
// ===========================================

func handleAddUnit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := productMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		var req models.ProductUnit
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		// Validate
		v := models.NewValidator()
		v.Check(req.ProductID > 0, "product_id", "Product ID is required")
		v.Check(req.UnitName != "", "unit_name", "Unit name is required")
		v.Check(req.ConversionFactor > 0, "conversion_factor", "Conversion factor must be greater than 0")
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		id, err := svc.AddUnit(r.Context(), &req)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.CreatedResponse(w, r, id, "Unit added successfully")
	}
}

func handleGetUnits() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := productMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		productID, err := strconv.Atoi(chi.URLParam(r, "productId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid product ID"))
			return
		}

		units, err := svc.GetUnits(r.Context(), productID)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, units)
	}
}

func handleUpdateUnit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := productMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid unit ID"))
			return
		}

		var req models.ProductUnit
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		err = svc.UpdateUnit(r.Context(), id, &req)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Unit updated successfully"})
	}
}

func handleDeleteUnit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := productMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid unit ID"))
			return
		}

		err = svc.DeleteUnit(r.Context(), id)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Unit deleted successfully"})
	}
}
