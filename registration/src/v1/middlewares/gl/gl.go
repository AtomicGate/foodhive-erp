package gl

import (
	"context"
	"net/http"

	"github.com/anas-dev-92/FoodHive/core/postgres"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/services/gl"
)

type contextKey string

const glServiceKey contextKey = "glService"

// New creates a middleware that injects the GL service into the request context.
func New(db postgres.Executor) func(http.Handler) http.Handler {
	service := gl.New(db.(postgres.Connection))
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), glServiceKey, service)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Instance retrieves the GL service from the request context.
func Instance(ctx context.Context) gl.GLService {
	return ctx.Value(glServiceKey).(gl.GLService)
}
