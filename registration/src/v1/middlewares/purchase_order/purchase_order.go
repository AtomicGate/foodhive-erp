package purchase_order

import (
	"context"
	"net/http"

	"github.com/anas-dev-92/FoodHive/core/postgres"
	poService "github.com/anas-dev-92/FoodHive/registration/src/v1/services/purchase_order"
)

type contextKey string

const purchaseOrderKey = contextKey("purchase_order_service")

// New creates a middleware that injects the purchase order service into the request context
func New(db postgres.Executor) func(http.Handler) http.Handler {
	svc := poService.New(db)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), purchaseOrderKey, svc)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Instance retrieves the purchase order service from the context
func Instance(ctx context.Context) (poService.PurchaseOrderService, bool) {
	svc, ok := ctx.Value(purchaseOrderKey).(poService.PurchaseOrderService)
	return svc, ok
}
