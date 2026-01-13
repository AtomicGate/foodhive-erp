package catch_weight

import (
	"context"
	"net/http"

	"github.com/anas-dev-92/FoodHive/core/postgres"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/services/catch_weight"
)

type contextKey string

const catchWeightServiceKey contextKey = "catchWeightService"

// New creates a middleware that injects the catch weight service into the request context.
func New(db postgres.Executor) func(http.Handler) http.Handler {
	service := catch_weight.New(db.(postgres.Connection))
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), catchWeightServiceKey, service)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Instance retrieves the catch weight service from the request context.
func Instance(ctx context.Context) catch_weight.CatchWeightService {
	return ctx.Value(catchWeightServiceKey).(catch_weight.CatchWeightService)
}
