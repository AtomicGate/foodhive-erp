package catch_weight

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/anas-dev-92/FoodHive/core/auth"
	"github.com/anas-dev-92/FoodHive/core/jwt"
	"github.com/anas-dev-92/FoodHive/core/postgres"
	authMiddleware "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/auth"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/models"
	cwService "github.com/anas-dev-92/FoodHive/registration/src/v1/services/catch_weight"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/utils/helper"
	"github.com/go-chi/chi/v5"
)

// Router creates the catch weight routes
func Router(db postgres.Executor, jwtService jwt.JWTService, authService auth.AuthService) chi.Router {
	r := chi.NewRouter()

	service := cwService.New(db.(postgres.Connection))

	r.Use(authMiddleware.Authenticate(jwtService))

	// Weight Capture Routes
	r.With(authMiddleware.Authorize(jwtService)).Post("/capture", handleCaptureCatchWeight(service, jwtService))
	r.With(authMiddleware.Authorize(jwtService)).Post("/capture/quick", handleQuickCaptureCatchWeight(service, jwtService))
	r.With(authMiddleware.Authorize(jwtService)).Post("/entries/{entryId}/pieces", handleAddPieceWeight(service))
	r.With(authMiddleware.Authorize(jwtService)).Put("/pieces/{pieceId}", handleUpdatePieceWeight(service))
	r.With(authMiddleware.Authorize(jwtService)).Delete("/pieces/{pieceId}", handleDeletePieceWeight(service))

	// Retrieval Routes
	r.With(authMiddleware.Authorize(jwtService)).Get("/entries/{id}", handleGetEntryByID(service))
	r.With(authMiddleware.Authorize(jwtService)).Get("/entries", handleListEntries(service))
	r.With(authMiddleware.Authorize(jwtService)).Get("/entries/{entryId}/pieces", handleGetPiecesByEntry(service))
	r.With(authMiddleware.Authorize(jwtService)).Get("/reference/{refType}/{refId}/product/{productId}", handleGetEntryByReference(service))

	// Product Configuration Routes
	r.With(authMiddleware.Authorize(jwtService)).Get("/products/{productId}/config", handleGetProductConfig(service))
	r.With(authMiddleware.Authorize(jwtService)).Put("/products/{productId}/config", handleUpdateProductConfig(service))

	// Report Routes
	r.With(authMiddleware.Authorize(jwtService)).Get("/reports/variance", handleGetVarianceReport(service))
	r.With(authMiddleware.Authorize(jwtService)).Get("/reports/lot/{productId}/{lotNumber}", handleGetLotSummary(service))

	// Billing Routes
	r.With(authMiddleware.Authorize(jwtService)).Post("/billing/adjustment", handleCalculateBillingAdjustment(service))
	r.With(authMiddleware.Authorize(jwtService)).Post("/entries/{entryId}/mark-billed", handleMarkAsBilled(service))

	// Validation Routes
	r.With(authMiddleware.Authorize(jwtService)).Post("/validate-weight", handleValidateWeight(service))

	return r
}

// ============================================
// Weight Capture Handlers
// ============================================

func handleCaptureCatchWeight(service cwService.CatchWeightService, jwtService jwt.JWTService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.CaptureCatchWeightRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		v := models.NewValidator()
		models.ValidateCatchWeightCapture(v, &req)
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		// Get user ID from token
		tokenData, err := jwtService.ParseTokenFromRequest(r)
		if err != nil {
			helper.UnauthorizedResponse(w, r)
			return
		}
		userID := int(tokenData["user_id"].(float64))

		id, err := service.CaptureCatchWeight(r.Context(), req, userID)
		if err != nil {
			switch {
			case errors.Is(err, cwService.ErrProductNotFound):
				helper.NotFoundResponse(w, r)
			case errors.Is(err, cwService.ErrNotCatchWeight):
				helper.BadRequestResponse(w, r, errors.New("product is not configured for catch weight"))
			case errors.Is(err, cwService.ErrDuplicateEntry):
				helper.BadRequestResponse(w, r, errors.New("catch weight already captured for this reference"))
			default:
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.CreatedResponse(w, r, id, "catch weight entry created")
	}
}

func handleQuickCaptureCatchWeight(service cwService.CatchWeightService, jwtService jwt.JWTService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.QuickCatchWeightRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		v := models.NewValidator()
		models.ValidateQuickCatchWeight(v, &req)
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		tokenData, err := jwtService.ParseTokenFromRequest(r)
		if err != nil {
			helper.UnauthorizedResponse(w, r)
			return
		}
		userID := int(tokenData["user_id"].(float64))

		id, err := service.QuickCaptureCatchWeight(r.Context(), req, userID)
		if err != nil {
			switch {
			case errors.Is(err, cwService.ErrProductNotFound):
				helper.NotFoundResponse(w, r)
			case errors.Is(err, cwService.ErrNotCatchWeight):
				helper.BadRequestResponse(w, r, errors.New("product is not configured for catch weight"))
			default:
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.CreatedResponse(w, r, id, "catch weight entry created")
	}
}

func handleAddPieceWeight(service cwService.CatchWeightService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entryID, err := strconv.Atoi(chi.URLParam(r, "entryId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid entry ID"))
			return
		}

		var req models.CapturePieceWeightRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		if req.Weight <= 0 {
			helper.BadRequestResponse(w, r, errors.New("weight must be positive"))
			return
		}

		id, err := service.AddPieceWeight(r.Context(), entryID, req)
		if err != nil {
			if errors.Is(err, cwService.ErrNotFound) {
				helper.NotFoundResponse(w, r)
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.CreatedResponse(w, r, id, "piece weight added")
	}
}

func handleUpdatePieceWeight(service cwService.CatchWeightService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pieceID, err := strconv.Atoi(chi.URLParam(r, "pieceId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid piece ID"))
			return
		}

		var req struct {
			Weight float64 `json:"weight"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		if req.Weight <= 0 {
			helper.BadRequestResponse(w, r, errors.New("weight must be positive"))
			return
		}

		if err := service.UpdatePieceWeight(r.Context(), pieceID, req.Weight); err != nil {
			if errors.Is(err, cwService.ErrNotFound) {
				helper.NotFoundResponse(w, r)
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.WriteJSON(w, http.StatusOK, helper.Envelope{"message": "piece weight updated"}, nil)
	}
}

func handleDeletePieceWeight(service cwService.CatchWeightService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pieceID, err := strconv.Atoi(chi.URLParam(r, "pieceId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid piece ID"))
			return
		}

		if err := service.DeletePieceWeight(r.Context(), pieceID); err != nil {
			if errors.Is(err, cwService.ErrNotFound) {
				helper.NotFoundResponse(w, r)
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.WriteJSON(w, http.StatusOK, helper.Envelope{"message": "piece deleted"}, nil)
	}
}

// ============================================
// Retrieval Handlers
// ============================================

func handleGetEntryByID(service cwService.CatchWeightService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid entry ID"))
			return
		}

		entry, err := service.GetEntryByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, cwService.ErrNotFound) {
				helper.NotFoundResponse(w, r)
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.WriteJSON(w, http.StatusOK, helper.Envelope{"catch_weight_entry": entry}, nil)
	}
}

func handleGetEntryByReference(service cwService.CatchWeightService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		refType := chi.URLParam(r, "refType")
		refID, err := strconv.Atoi(chi.URLParam(r, "refId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid reference ID"))
			return
		}
		productID, err := strconv.Atoi(chi.URLParam(r, "productId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid product ID"))
			return
		}

		entry, err := service.GetEntryByReference(r.Context(), refType, refID, productID)
		if err != nil {
			if errors.Is(err, cwService.ErrNotFound) {
				helper.NotFoundResponse(w, r)
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.WriteJSON(w, http.StatusOK, helper.Envelope{"catch_weight_entry": entry}, nil)
	}
}

func handleListEntries(service cwService.CatchWeightService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filters := models.CatchWeightListFilters{
			ReferenceType: r.URL.Query().Get("reference_type"),
			LotNumber:     r.URL.Query().Get("lot_number"),
			DateFrom:      r.URL.Query().Get("date_from"),
			DateTo:        r.URL.Query().Get("date_to"),
			HasVariance:   r.URL.Query().Get("has_variance") == "true",
			Page:          1,
			PageSize:      50,
		}

		if productID := r.URL.Query().Get("product_id"); productID != "" {
			id, err := strconv.Atoi(productID)
			if err == nil {
				filters.ProductID = &id
			}
		}

		if page := r.URL.Query().Get("page"); page != "" {
			if p, err := strconv.Atoi(page); err == nil && p > 0 {
				filters.Page = p
			}
		}

		if pageSize := r.URL.Query().Get("page_size"); pageSize != "" {
			if ps, err := strconv.Atoi(pageSize); err == nil && ps > 0 && ps <= 100 {
				filters.PageSize = ps
			}
		}

		entries, total, err := service.ListEntries(r.Context(), filters)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.WriteJSON(w, http.StatusOK, helper.Envelope{
			"catch_weight_entries": entries,
			"total":                total,
			"page":                 filters.Page,
			"page_size":            filters.PageSize,
		}, nil)
	}
}

func handleGetPiecesByEntry(service cwService.CatchWeightService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entryID, err := strconv.Atoi(chi.URLParam(r, "entryId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid entry ID"))
			return
		}

		pieces, err := service.GetPiecesByEntry(r.Context(), entryID)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.WriteJSON(w, http.StatusOK, helper.Envelope{"pieces": pieces}, nil)
	}
}

// ============================================
// Product Configuration Handlers
// ============================================

func handleGetProductConfig(service cwService.CatchWeightService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productID, err := strconv.Atoi(chi.URLParam(r, "productId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid product ID"))
			return
		}

		config, err := service.GetProductConfig(r.Context(), productID)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.WriteJSON(w, http.StatusOK, helper.Envelope{"config": config}, nil)
	}
}

func handleUpdateProductConfig(service cwService.CatchWeightService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productID, err := strconv.Atoi(chi.URLParam(r, "productId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid product ID"))
			return
		}

		var req models.UpdateCatchWeightConfigRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		if err := service.UpdateProductConfig(r.Context(), productID, req); err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.WriteJSON(w, http.StatusOK, helper.Envelope{"message": "configuration updated"}, nil)
	}
}

// ============================================
// Report Handlers
// ============================================

func handleGetVarianceReport(service cwService.CatchWeightService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var productID *int
		if pid := r.URL.Query().Get("product_id"); pid != "" {
			id, err := strconv.Atoi(pid)
			if err == nil {
				productID = &id
			}
		}

		dateFrom := r.URL.Query().Get("date_from")
		dateTo := r.URL.Query().Get("date_to")

		reports, err := service.GetVarianceReport(r.Context(), productID, dateFrom, dateTo)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.WriteJSON(w, http.StatusOK, helper.Envelope{"variance_report": reports}, nil)
	}
}

func handleGetLotSummary(service cwService.CatchWeightService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productID, err := strconv.Atoi(chi.URLParam(r, "productId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid product ID"))
			return
		}

		lotNumber := chi.URLParam(r, "lotNumber")
		if lotNumber == "" {
			helper.BadRequestResponse(w, r, errors.New("lot number is required"))
			return
		}

		summary, err := service.GetLotSummary(r.Context(), productID, lotNumber)
		if err != nil {
			if errors.Is(err, cwService.ErrNotFound) {
				helper.NotFoundResponse(w, r)
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.WriteJSON(w, http.StatusOK, helper.Envelope{"lot_summary": summary}, nil)
	}
}

// ============================================
// Billing Handlers
// ============================================

func handleCalculateBillingAdjustment(service cwService.CatchWeightService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			InvoiceID      int     `json:"invoice_id"`
			InvoiceLineID  int     `json:"invoice_line_id"`
			ProductID      int     `json:"product_id"`
			StandardWeight float64 `json:"standard_weight"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		adjustment, err := service.CalculateBillingAdjustment(r.Context(), req.InvoiceID, req.InvoiceLineID, req.ProductID, req.StandardWeight)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.WriteJSON(w, http.StatusOK, helper.Envelope{"billing_adjustment": adjustment}, nil)
	}
}

func handleMarkAsBilled(service cwService.CatchWeightService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entryID, err := strconv.Atoi(chi.URLParam(r, "entryId"))
		if err != nil {
			helper.BadRequestResponse(w, r, errors.New("invalid entry ID"))
			return
		}

		if err := service.MarkAsBilled(r.Context(), entryID); err != nil {
			if errors.Is(err, cwService.ErrNotFound) {
				helper.NotFoundResponse(w, r)
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.WriteJSON(w, http.StatusOK, helper.Envelope{"message": "entry marked as billed"}, nil)
	}
}

// ============================================
// Validation Handlers
// ============================================

func handleValidateWeight(service cwService.CatchWeightService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ProductID int     `json:"product_id"`
			Weight    float64 `json:"weight"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		if err := service.ValidateWeight(r.Context(), req.ProductID, req.Weight); err != nil {
			if errors.Is(err, cwService.ErrWeightOutOfRange) {
				helper.WriteJSON(w, http.StatusOK, helper.Envelope{
					"valid":   false,
					"message": err.Error(),
				}, nil)
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.WriteJSON(w, http.StatusOK, helper.Envelope{
			"valid":   true,
			"message": "weight is within acceptable range",
		}, nil)
	}
}
