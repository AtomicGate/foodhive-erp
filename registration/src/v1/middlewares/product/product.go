package product

import (
	"context"
	"net/http"

	"github.com/anas-dev-92/FoodHive/core/postgres"
	productService "github.com/anas-dev-92/FoodHive/registration/src/v1/services/product"
)

type contextKey string

const productKey = contextKey("product_service")

// New creates a middleware that injects the product service into the request context
func New(db postgres.Executor) func(http.Handler) http.Handler {
	svc := productService.New(db)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), productKey, svc)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Instance retrieves the product service from the context
func Instance(ctx context.Context) (productService.ProductService, bool) {
	svc, ok := ctx.Value(productKey).(productService.ProductService)
	return svc, ok
}
