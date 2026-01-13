package pricing

import (
	"context"
	"net/http"

	"github.com/anas-dev-92/FoodHive/core/postgres"
	pricingService "github.com/anas-dev-92/FoodHive/registration/src/v1/services/pricing"
)

type contextKey string

const pricingKey = contextKey("pricing_service")

// New creates a middleware that injects the pricing service into the request context
func New(db postgres.Executor) func(http.Handler) http.Handler {
	svc := pricingService.New(db)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), pricingKey, svc)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Instance retrieves the pricing service from the context
func Instance(ctx context.Context) (pricingService.PricingService, bool) {
	svc, ok := ctx.Value(pricingKey).(pricingService.PricingService)
	return svc, ok
}
