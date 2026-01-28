package login

import (
	"github.com/anas-dev-92/FoodHive/core/jwt"
	"github.com/anas-dev-92/FoodHive/core/postgres"
	"github.com/go-chi/chi/v5"
)

// Router creates the login routes (public - no auth middleware)
func Router(db postgres.Executor, jwtService jwt.JWTService) chi.Router {
	r := chi.NewRouter()

	// Login endpoint - no authentication required (public)
	r.Post("/login", Handler(jwtService, db))

	return r
}
