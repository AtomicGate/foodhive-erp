package warehouse

import (
	"context"
	"net/http"

	"github.com/anas-dev-92/FoodHive/core/postgres"
	warehouseService "github.com/anas-dev-92/FoodHive/registration/src/v1/services/warehouse"
)

type contextKey string

const warehouseKey = contextKey("warehouse_service")

// New creates a middleware that injects the warehouse service into the request context
func New(db postgres.Executor) func(http.Handler) http.Handler {
	svc := warehouseService.New(db)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), warehouseKey, svc)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Instance retrieves the warehouse service from the context
func Instance(ctx context.Context) (warehouseService.WarehouseService, bool) {
	svc, ok := ctx.Value(warehouseKey).(warehouseService.WarehouseService)
	return svc, ok
}
