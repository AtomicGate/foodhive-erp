package ap

import (
	"context"
	"net/http"

	"github.com/anas-dev-92/FoodHive/core/postgres"
	apService "github.com/anas-dev-92/FoodHive/registration/src/v1/services/ap"
)

type contextKey string

const apKey = contextKey("ap_service")

// New creates a middleware that injects the AP service into the request context
func New(db postgres.Executor) func(http.Handler) http.Handler {
	svc := apService.New(db)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), apKey, svc)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Instance retrieves the AP service from the context
func Instance(ctx context.Context) (apService.APService, bool) {
	svc, ok := ctx.Value(apKey).(apService.APService)
	return svc, ok
}
