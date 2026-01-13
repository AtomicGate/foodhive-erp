package ar

import (
	"context"
	"net/http"

	"github.com/anas-dev-92/FoodHive/core/postgres"
	arService "github.com/anas-dev-92/FoodHive/registration/src/v1/services/ar"
)

type contextKey string

const arKey = contextKey("ar_service")

// New creates a middleware that injects the AR service into the request context
func New(db postgres.Executor) func(http.Handler) http.Handler {
	svc := arService.New(db)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), arKey, svc)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Instance retrieves the AR service from the context
func Instance(ctx context.Context) (arService.ARService, bool) {
	svc, ok := ctx.Value(arKey).(arService.ARService)
	return svc, ok
}
