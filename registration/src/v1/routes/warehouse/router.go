package warehouse

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/anas-dev-92/FoodHive/core/auth"
	"github.com/anas-dev-92/FoodHive/core/jwt"
	"github.com/anas-dev-92/FoodHive/core/postgres"
	authMiddleware "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/auth"
	warehouseMiddleware "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/warehouse"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/models"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/utils/helper"
	"github.com/go-chi/chi/v5"
)

func Router(db postgres.Executor, jwtService jwt.JWTService, authService auth.AuthService) chi.Router {
	app := chi.NewRouter()

	// Inject warehouse service
	app.Use(warehouseMiddleware.New(db))

	// Apply authentication
	app.Use(authMiddleware.Authenticate(jwtService))

	// ===========================================
	// Warehouse Routes
	// ===========================================
	app.With(authMiddleware.Authorize(jwtService)).Post("/create", handleCreate())
	app.With(authMiddleware.Authorize(jwtService)).Get("/get/{id}", handleGetByID())
	app.With(authMiddleware.Authorize(jwtService)).Get("/code/{code}", handleGetByCode())
	app.With(authMiddleware.Authorize(jwtService)).Put("/update/{id}", handleUpdate())
	app.With(authMiddleware.Authorize(jwtService)).Delete("/delete/{id}", handleDelete())
	app.With(authMiddleware.Authorize(jwtService)).Get("/list", handleList())

	// ===========================================
	// Zone Routes
	// ===========================================
	app.Route("/zones", func(r chi.Router) {
		r.With(authMiddleware.Authorize(jwtService)).Post("/create", handleCreateZone())
		r.With(authMiddleware.Authorize(jwtService)).Get("/get/{id}", handleGetZone())
		r.With(authMiddleware.Authorize(jwtService)).Get("/warehouse/{warehouseId}", handleGetZonesByWarehouse())
		r.With(authMiddleware.Authorize(jwtService)).Put("/update/{id}", handleUpdateZone())
		r.With(authMiddleware.Authorize(jwtService)).Delete("/delete/{id}", handleDeleteZone())
	})

	// ===========================================
	// Location Routes
	// ===========================================
	app.Route("/locations", func(r chi.Router) {
		r.With(authMiddleware.Authorize(jwtService)).Post("/create", handleCreateLocation())
		r.With(authMiddleware.Authorize(jwtService)).Get("/get/{id}", handleGetLocation())
		r.With(authMiddleware.Authorize(jwtService)).Get("/code/{code}", handleGetLocationByCode())
		r.With(authMiddleware.Authorize(jwtService)).Get("/warehouse/{warehouseId}", handleGetLocationsByWarehouse())
		r.With(authMiddleware.Authorize(jwtService)).Get("/zone/{zoneId}", handleGetLocationsByZone())
		r.With(authMiddleware.Authorize(jwtService)).Put("/update/{id}", handleUpdateLocation())
		r.With(authMiddleware.Authorize(jwtService)).Delete("/delete/{id}", handleDeleteLocation())
	})

	return app
}

// ===========================================
// Warehouse Handlers
// ===========================================

func handleCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := warehouseMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		var req models.CreateWarehouseRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		v := models.NewValidator()
		models.ValidateWarehouse(v, &req)
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		id, err := svc.Create(r.Context(), &req)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.CreatedResponse(w, r, id, "Warehouse created successfully")
	}
}

func handleGetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := warehouseMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid warehouse ID"))
			return
		}

		warehouse, err := svc.GetByID(r.Context(), id)
		if err != nil {
			helper.NotFoundResponse(w, r)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, warehouse)
	}
}

func handleGetByCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := warehouseMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		code := chi.URLParam(r, "code")
		warehouse, err := svc.GetByCode(r.Context(), code)
		if err != nil {
			helper.NotFoundResponse(w, r)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, warehouse)
	}
}

func handleUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := warehouseMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid warehouse ID"))
			return
		}

		var req models.UpdateWarehouseRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		err = svc.Update(r.Context(), id, &req)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Warehouse updated successfully"})
	}
}

func handleDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := warehouseMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid warehouse ID"))
			return
		}

		err = svc.Delete(r.Context(), id)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Warehouse deleted successfully"})
	}
}

func handleList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := warehouseMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		filters := models.WarehouseListFilters{
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
		if isActive := r.URL.Query().Get("is_active"); isActive != "" {
			active := isActive == "true"
			filters.IsActive = &active
		}

		warehouses, total, err := svc.List(r.Context(), &filters)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		totalPages := int(total) / filters.PageSize
		if int(total)%filters.PageSize > 0 {
			totalPages++
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{
			"data": warehouses,
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
// Zone Handlers
// ===========================================

func handleCreateZone() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := warehouseMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		var req models.CreateZoneRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		v := models.NewValidator()
		models.ValidateZone(v, &req)
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		id, err := svc.CreateZone(r.Context(), &req)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.CreatedResponse(w, r, id, "Zone created successfully")
	}
}

func handleGetZone() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := warehouseMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid zone ID"))
			return
		}

		zone, err := svc.GetZoneByID(r.Context(), id)
		if err != nil {
			helper.NotFoundResponse(w, r)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, zone)
	}
}

func handleGetZonesByWarehouse() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := warehouseMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		warehouseID, err := strconv.Atoi(chi.URLParam(r, "warehouseId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid warehouse ID"))
			return
		}

		zones, err := svc.GetZonesByWarehouse(r.Context(), warehouseID)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, zones)
	}
}

func handleUpdateZone() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := warehouseMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid zone ID"))
			return
		}

		var req models.CreateZoneRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		err = svc.UpdateZone(r.Context(), id, &req)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Zone updated successfully"})
	}
}

func handleDeleteZone() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := warehouseMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid zone ID"))
			return
		}

		err = svc.DeleteZone(r.Context(), id)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Zone deleted successfully"})
	}
}

// ===========================================
// Location Handlers
// ===========================================

func handleCreateLocation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := warehouseMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		var req models.CreateLocationRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		v := models.NewValidator()
		models.ValidateLocation(v, &req)
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		id, err := svc.CreateLocation(r.Context(), &req)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.CreatedResponse(w, r, id, "Location created successfully")
	}
}

func handleGetLocation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := warehouseMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid location ID"))
			return
		}

		location, err := svc.GetLocationByID(r.Context(), id)
		if err != nil {
			helper.NotFoundResponse(w, r)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, location)
	}
}

func handleGetLocationByCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := warehouseMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		code := chi.URLParam(r, "code")
		location, err := svc.GetLocationByCode(r.Context(), code)
		if err != nil {
			helper.NotFoundResponse(w, r)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, location)
	}
}

func handleGetLocationsByWarehouse() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := warehouseMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		warehouseID, err := strconv.Atoi(chi.URLParam(r, "warehouseId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid warehouse ID"))
			return
		}

		locations, err := svc.GetLocationsByWarehouse(r.Context(), warehouseID)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, locations)
	}
}

func handleGetLocationsByZone() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := warehouseMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		zoneID, err := strconv.Atoi(chi.URLParam(r, "zoneId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid zone ID"))
			return
		}

		locations, err := svc.GetLocationsByZone(r.Context(), zoneID)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, locations)
	}
}

func handleUpdateLocation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := warehouseMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid location ID"))
			return
		}

		var req models.CreateLocationRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		err = svc.UpdateLocation(r.Context(), id, &req)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Location updated successfully"})
	}
}

func handleDeleteLocation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := warehouseMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid location ID"))
			return
		}

		err = svc.DeleteLocation(r.Context(), id)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Location deleted successfully"})
	}
}
