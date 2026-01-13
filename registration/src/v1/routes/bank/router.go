package bank

import (
	"github.com/anas-dev-92/FoodHive/core/auth"
	"github.com/anas-dev-92/FoodHive/core/jwt"
	"github.com/anas-dev-92/FoodHive/core/postgres"
	"github.com/go-chi/chi/v5"
)

// TODO: Implement Bank & Reconciliation routes
func Router(db postgres.Executor, jwtService jwt.JWTService, authService auth.AuthService) chi.Router {
	r := chi.NewRouter()
	// Implement Bank routes here
	return r
}
