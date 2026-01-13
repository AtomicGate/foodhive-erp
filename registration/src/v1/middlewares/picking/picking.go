package picking

import (
	"context"
	"net/http"

	"github.com/anas-dev-92/FoodHive/core/postgres"
	pickingService "github.com/anas-dev-92/FoodHive/registration/src/v1/services/picking"
)

type contextKey string

const pickingKey = contextKey("picking_service")

// New creates a middleware that injects the picking service into the request context
func New(db postgres.Executor) func(http.Handler) http.Handler {
	svc := pickingService.New(db)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), pickingKey, svc)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Instance retrieves the picking service from the context
func Instance(ctx context.Context) (pickingService.PickingService, bool) {
	svc, ok := ctx.Value(pickingKey).(pickingService.PickingService)
	return svc, ok
}
