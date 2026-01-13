package employee

import (
	"github.com/anas-dev-92/FoodHive/core/auth"
	"github.com/anas-dev-92/FoodHive/core/jwt"
	"github.com/anas-dev-92/FoodHive/core/postgres"
	authMiddleware "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/auth"
	employeeService "github.com/anas-dev-92/FoodHive/registration/src/v1/services/employee"
	"github.com/go-chi/chi/v5"
)

func Router(db postgres.Executor, jwtService jwt.JWTService, authService auth.AuthService) chi.Router {
	r := chi.NewRouter()

	// Initialize service
	service := employeeService.New(db.(postgres.Connection))

	// Apply authentication middleware globally
	r.Use(authMiddleware.Authenticate(jwtService))

	// Routes with authorization
	r.With(authMiddleware.Authorize(jwtService)).Post("/create", HandlerCreate(service))
	r.With(authMiddleware.Authorize(jwtService)).Get("/get/{id}", HandlerGetByID(service))
	r.With(authMiddleware.Authorize(jwtService)).Put("/update/{id}", HandlerUpdate(service))
	r.With(authMiddleware.Authorize(jwtService)).Delete("/delete/{id}", HandlerDelete(service))
	r.With(authMiddleware.Authorize(jwtService)).Get("/list", HandlerList(service))

	return r
}
