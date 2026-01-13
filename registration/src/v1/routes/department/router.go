package department

import (
	"net/http"
	"strconv"

	"github.com/anas-dev-92/FoodHive/core/auth"
	"github.com/anas-dev-92/FoodHive/core/jwt"
	"github.com/anas-dev-92/FoodHive/core/postgres"
	authMiddleware "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/auth"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/models"
	departmentService "github.com/anas-dev-92/FoodHive/registration/src/v1/services/department"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/utils/helper"
	"github.com/go-chi/chi/v5"
)

func Router(db postgres.Executor, jwtService jwt.JWTService, authService auth.AuthService) chi.Router {
	r := chi.NewRouter()

	service := departmentService.New(db.(postgres.Connection))

	r.Use(authMiddleware.Authenticate(jwtService))

	r.With(authMiddleware.Authorize(jwtService)).Post("/create", handleCreate(service))
	r.With(authMiddleware.Authorize(jwtService)).Get("/get/{id}", handleGetByID(service))
	r.With(authMiddleware.Authorize(jwtService)).Put("/update/{id}", handleUpdate(service))
	r.With(authMiddleware.Authorize(jwtService)).Delete("/delete/{id}", handleDelete(service))
	r.With(authMiddleware.Authorize(jwtService)).Get("/list", handleList(service))

	return r
}

func handleCreate(service departmentService.DepartmentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.CreateDepartmentRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		v := helper.New()
		v.Check(req.Name != "", "name", "must be provided")
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		id, err := service.Create(r.Context(), req)
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.CreatedResponse(w, r, id, "department created successfully")
	}
}

func handleGetByID(service departmentService.DepartmentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || id < 1 {
			helper.NotFoundResponse(w, r)
			return
		}

		dept, err := service.GetByID(r.Context(), id)
		if err != nil {
			if err == departmentService.ErrNotFound {
				helper.NotFoundResponse(w, r)
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, dept)
	}
}

func handleUpdate(service departmentService.DepartmentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || id < 1 {
			helper.NotFoundResponse(w, r)
			return
		}

		var req models.UpdateDepartmentRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		err = service.Update(r.Context(), id, req)
		if err != nil {
			if err == departmentService.ErrNotFound {
				helper.NotFoundResponse(w, r)
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "department updated successfully"})
	}
}

func handleDelete(service departmentService.DepartmentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || id < 1 {
			helper.NotFoundResponse(w, r)
			return
		}

		err = service.Delete(r.Context(), id)
		if err != nil {
			if err == departmentService.ErrNotFound {
				helper.NotFoundResponse(w, r)
			} else {
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, helper.Envelope{"message": "department deleted successfully"})
	}
}

func handleList(service departmentService.DepartmentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		departments, err := service.List(r.Context())
		if err != nil {
			helper.ServerErrorResponse(w, r, err)
			return
		}

		helper.SuccessResponse(w, r, http.StatusOK, departments)
	}
}
