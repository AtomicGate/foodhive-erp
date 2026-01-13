package login

import (
	"github.com/anas-dev-92/FoodHive/core/jwt"
	"github.com/anas-dev-92/FoodHive/core/postgres"
	"github.com/go-chi/chi/v5"
)

func Router(secretKey string, db postgres.Executor) chi.Router {

	app := chi.NewRouter()

	// Initialize your custom JWT service with the secret key
	jwtService := jwt.New(secretKey)

	// Pass jwtService and db to the Handler
	app.Post("/login", Handler(jwtService, db))

	return app
}
