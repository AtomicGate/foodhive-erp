package sales_order

import (
	"context"
	"net/http"

	"github.com/anas-dev-92/FoodHive/core/postgres"
	soService "github.com/anas-dev-92/FoodHive/registration/src/v1/services/sales_order"
)

type contextKey string

const salesOrderKey = contextKey("sales_order_service")

// New creates a middleware that injects the sales order service into the request context
func New(db postgres.Executor) func(http.Handler) http.Handler {
	svc := soService.New(db)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), salesOrderKey, svc)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Instance retrieves the sales order service from the context
func Instance(ctx context.Context) (soService.SalesOrderService, bool) {
	svc, ok := ctx.Value(salesOrderKey).(soService.SalesOrderService)
	return svc, ok
}
