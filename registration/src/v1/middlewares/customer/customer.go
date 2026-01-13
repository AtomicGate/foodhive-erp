package customer

import (
	"context"
	"net/http"

	"github.com/anas-dev-92/FoodHive/core/postgres"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/services/customer"
)

type contextKey string

const customerKey = contextKey("customer_service")

func New(db postgres.Executor) func(http.Handler) http.Handler {
	customerService := customer.New(db.(postgres.Connection))

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), customerKey, customerService)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Instance(ctx context.Context) (customer.CustomerService, bool) {
	svc, ok := ctx.Value(customerKey).(customer.CustomerService)
	return svc, ok
}
