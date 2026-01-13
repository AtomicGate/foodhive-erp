package inventory

import (
	"context"
	"net/http"

	"github.com/anas-dev-92/FoodHive/core/postgres"
	inventoryService "github.com/anas-dev-92/FoodHive/registration/src/v1/services/inventory"
)

type contextKey string

const inventoryKey = contextKey("inventory_service")

// New creates a middleware that injects the inventory service into the request context
func New(db postgres.Executor) func(http.Handler) http.Handler {
	svc := inventoryService.New(db)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), inventoryKey, svc)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Instance retrieves the inventory service from the context
func Instance(ctx context.Context) (inventoryService.InventoryService, bool) {
	svc, ok := ctx.Value(inventoryKey).(inventoryService.InventoryService)
	return svc, ok
}
