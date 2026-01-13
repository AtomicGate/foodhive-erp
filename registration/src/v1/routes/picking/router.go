package picking

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/anas-dev-92/FoodHive/core/auth"
	"github.com/anas-dev-92/FoodHive/core/jwt"
	"github.com/anas-dev-92/FoodHive/core/postgres"
	authMiddleware "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/auth"
	pickingMiddleware "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/picking"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/models"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/utils/helper"
	"github.com/go-chi/chi/v5"
)

func Router(db postgres.Executor, jwtService jwt.JWTService, authService auth.AuthService) chi.Router {
	app := chi.NewRouter()

	// Inject picking service
	app.Use(pickingMiddleware.New(db))

	// Apply authentication
	app.Use(authMiddleware.Authenticate(jwtService))

	// ===========================================
	// Route Management
	// ===========================================
	app.With(authMiddleware.Authorize(jwtService)).Post("/routes/create", handleCreateRoute())
	app.With(authMiddleware.Authorize(jwtService)).Get("/routes/get/{id}", handleGetRoute())
	app.With(authMiddleware.Authorize(jwtService)).Put("/routes/update/{id}", handleUpdateRoute())
	app.With(authMiddleware.Authorize(jwtService)).Delete("/routes/delete/{id}", handleDeleteRoute())
	app.With(authMiddleware.Authorize(jwtService)).Get("/routes/list", handleListRoutes())

	// ===========================================
	// Route Stops
	// ===========================================
	app.With(authMiddleware.Authorize(jwtService)).Post("/routes/{routeId}/stops", handleAddRouteStop())
	app.With(authMiddleware.Authorize(jwtService)).Put("/routes/stops/{stopId}", handleUpdateRouteStop())
	app.With(authMiddleware.Authorize(jwtService)).Delete("/routes/stops/{stopId}", handleDeleteRouteStop())
	app.With(authMiddleware.Authorize(jwtService)).Post("/routes/{routeId}/reorder", handleReorderStops())

	// ===========================================
	// Pick Lists
	// ===========================================
	app.With(authMiddleware.Authorize(jwtService)).Post("/create", handleCreatePickList())
	app.With(authMiddleware.Authorize(jwtService)).Post("/generate", handleGeneratePickList())
	app.With(authMiddleware.Authorize(jwtService)).Get("/get/{id}", handleGetPickList())
	app.With(authMiddleware.Authorize(jwtService)).Get("/list", handleListPickLists())
	app.With(authMiddleware.Authorize(jwtService)).Post("/{id}/start", handleStartPicking())
	app.With(authMiddleware.Authorize(jwtService)).Post("/{id}/complete", handleCompletePicking())
	app.With(authMiddleware.Authorize(jwtService)).Post("/{id}/cancel", handleCancelPickList())

	// ===========================================
	// Pick Lines
	// ===========================================
	app.With(authMiddleware.Authorize(jwtService)).Get("/{id}/lines", handleGetPickLines())
	app.With(authMiddleware.Authorize(jwtService)).Post("/lines/{lineId}/confirm", handleConfirmPickLine())

	// ===========================================
	// Reports
	// ===========================================
	app.With(authMiddleware.Authorize(jwtService)).Get("/master-pick", handleMasterPickReport())
	app.With(authMiddleware.Authorize(jwtService)).Get("/suggested-picking", handleSuggestedPicking())

	return app
}

// ===========================================
// Route Handlers
// ===========================================

func handleCreateRoute() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := pickingMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		var req models.CreateRouteRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		v := models.NewValidator()
		models.ValidateRoute(v, &req)
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		id, err := svc.CreateRoute(r.Context(), &req)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.CreatedResponse(w, r, id, "Route created successfully")
	}
}

func handleGetRoute() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := pickingMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid route ID"))
			return
		}

		route, err := svc.GetRoute(r.Context(), id)
		if err != nil {
			helper.NotFoundResponse(w, r)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, route)
	}
}

func handleUpdateRoute() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := pickingMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid route ID"))
			return
		}

		var req models.UpdateRouteRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		err = svc.UpdateRoute(r.Context(), id, &req)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Route updated successfully"})
	}
}

func handleDeleteRoute() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := pickingMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid route ID"))
			return
		}

		err = svc.DeleteRoute(r.Context(), id)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Route deleted successfully"})
	}
}

func handleListRoutes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := pickingMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		var warehouseID *int
		if wid, err := strconv.Atoi(r.URL.Query().Get("warehouse_id")); err == nil {
			warehouseID = &wid
		}

		activeOnly := r.URL.Query().Get("active_only") == "true"

		routes, err := svc.ListRoutes(r.Context(), warehouseID, activeOnly)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, routes)
	}
}

// ===========================================
// Route Stop Handlers
// ===========================================

func handleAddRouteStop() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := pickingMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		routeID, err := strconv.Atoi(chi.URLParam(r, "routeId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid route ID"))
			return
		}

		var req models.AddRouteStopRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		v := models.NewValidator()
		models.ValidateRouteStop(v, &req)
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		id, err := svc.AddRouteStop(r.Context(), routeID, &req)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.CreatedResponse(w, r, id, "Route stop added successfully")
	}
}

func handleUpdateRouteStop() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := pickingMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		stopID, err := strconv.Atoi(chi.URLParam(r, "stopId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid stop ID"))
			return
		}

		var req models.AddRouteStopRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		err = svc.UpdateRouteStop(r.Context(), stopID, &req)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Route stop updated successfully"})
	}
}

func handleDeleteRouteStop() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := pickingMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		stopID, err := strconv.Atoi(chi.URLParam(r, "stopId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid stop ID"))
			return
		}

		err = svc.DeleteRouteStop(r.Context(), stopID)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Route stop deleted successfully"})
	}
}

func handleReorderStops() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := pickingMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		routeID, err := strconv.Atoi(chi.URLParam(r, "routeId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid route ID"))
			return
		}

		var req struct {
			StopOrder []int `json:"stop_order"`
		}
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		err = svc.ReorderStops(r.Context(), routeID, req.StopOrder)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Stops reordered successfully"})
	}
}

// ===========================================
// Pick List Handlers
// ===========================================

func handleCreatePickList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := pickingMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		var req models.CreatePickListRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		v := models.NewValidator()
		models.ValidatePickList(v, &req)
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		createdBy := 1 // TODO: Get from auth context

		id, err := svc.CreatePickList(r.Context(), &req, createdBy)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.CreatedResponse(w, r, id, "Pick list created successfully")
	}
}

func handleGeneratePickList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := pickingMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		var req models.GeneratePickListRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		v := models.NewValidator()
		v.Check(req.WarehouseID > 0, "warehouse_id", "Warehouse is required")
		v.Check(req.PickDate != "", "pick_date", "Pick date is required")
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		createdBy := 1 // TODO: Get from auth context

		id, err := svc.GeneratePickListForRoute(r.Context(), &req, createdBy)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.CreatedResponse(w, r, id, "Pick list generated successfully")
	}
}

func handleGetPickList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := pickingMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid pick list ID"))
			return
		}

		pickList, err := svc.GetPickList(r.Context(), id)
		if err != nil {
			helper.NotFoundResponse(w, r)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, pickList)
	}
}

func handleListPickLists() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := pickingMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		filters := models.PickListFilters{
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
		if warehouseID, err := strconv.Atoi(r.URL.Query().Get("warehouse_id")); err == nil {
			filters.WarehouseID = &warehouseID
		}
		if routeID, err := strconv.Atoi(r.URL.Query().Get("route_id")); err == nil {
			filters.RouteID = &routeID
		}
		if pickerID, err := strconv.Atoi(r.URL.Query().Get("picker_id")); err == nil {
			filters.PickerID = &pickerID
		}
		if status := r.URL.Query().Get("status"); status != "" {
			s := models.PickListStatus(status)
			filters.Status = &s
		}

		pickLists, total, err := svc.ListPickLists(r.Context(), &filters)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		totalPages := int(total) / filters.PageSize
		if int(total)%filters.PageSize > 0 {
			totalPages++
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{
			"data": pickLists,
			"pagination": helper.Envelope{
				"page":        filters.Page,
				"page_size":   filters.PageSize,
				"total_items": total,
				"total_pages": totalPages,
			},
		})
	}
}

func handleStartPicking() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := pickingMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid pick list ID"))
			return
		}

		pickerID := 1 // TODO: Get from auth context

		err = svc.StartPicking(r.Context(), id, pickerID)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Picking started"})
	}
}

func handleCompletePicking() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := pickingMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid pick list ID"))
			return
		}

		err = svc.CompletePicking(r.Context(), id)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Picking completed"})
	}
}

func handleCancelPickList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := pickingMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid pick list ID"))
			return
		}

		err = svc.CancelPickList(r.Context(), id)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Pick list cancelled"})
	}
}

// ===========================================
// Pick Line Handlers
// ===========================================

func handleGetPickLines() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := pickingMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid pick list ID"))
			return
		}

		lines, err := svc.GetPickLines(r.Context(), id)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, lines)
	}
}

func handleConfirmPickLine() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := pickingMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		lineID, err := strconv.Atoi(chi.URLParam(r, "lineId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid line ID"))
			return
		}

		var req models.ConfirmPickLineRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		pickerID := 1 // TODO: Get from auth context

		err = svc.ConfirmPickLine(r.Context(), lineID, &req, pickerID)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "Pick line confirmed"})
	}
}

// ===========================================
// Report Handlers
// ===========================================

func handleMasterPickReport() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := pickingMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		warehouseID, err := strconv.Atoi(r.URL.Query().Get("warehouse_id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("warehouse_id is required"))
			return
		}

		pickDate := r.URL.Query().Get("pick_date")
		if pickDate == "" {
			helper.BadRequestResponse(w, r, errors.New("pick_date is required"))
			return
		}

		var routeID *int
		if rid, err := strconv.Atoi(r.URL.Query().Get("route_id")); err == nil {
			routeID = &rid
		}

		items, err := svc.GetMasterPickReport(r.Context(), warehouseID, pickDate, routeID)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, items)
	}
}

func handleSuggestedPicking() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, ok := pickingMiddleware.Instance(r.Context())
		if !ok {
			helper.ServerErrorResponse(w, r, errors.New("service unavailable"))
			return
		}

		productID, err := strconv.Atoi(r.URL.Query().Get("product_id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("product_id is required"))
			return
		}

		warehouseID, err := strconv.Atoi(r.URL.Query().Get("warehouse_id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("warehouse_id is required"))
			return
		}

		quantity, err := strconv.ParseFloat(r.URL.Query().Get("quantity"), 64)
		if err != nil || quantity <= 0 {
			helper.BadRequestResponse(w, r, errors.New("valid quantity is required"))
			return
		}

		suggestions, err := svc.GetSuggestedPicking(r.Context(), productID, warehouseID, quantity)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, suggestions)
	}
}
