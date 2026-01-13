package vendor

import (
	"context"
	"net/http"

	"github.com/anas-dev-92/FoodHive/core/postgres"
	vendorService "github.com/anas-dev-92/FoodHive/registration/src/v1/services/vendor"
)

type contextKey string

const vendorKey = contextKey("vendor_service")

// New creates a middleware that injects the vendor service into the request context
func New(db postgres.Executor) func(http.Handler) http.Handler {
	svc := vendorService.New(db)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), vendorKey, svc)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Instance retrieves the vendor service from the context
func Instance(ctx context.Context) (vendorService.VendorService, bool) {
	svc, ok := ctx.Value(vendorKey).(vendorService.VendorService)
	return svc, ok
}
