package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/anas-dev-92/FoodHive/core/auth"
	"github.com/anas-dev-92/FoodHive/core/jwt"
	"github.com/anas-dev-92/FoodHive/core/postgres"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/utils/helper"
	"github.com/go-chi/chi/v5"
)

// Context keys for user data
type contextKey string

const (
	EmailKey  contextKey = "email"
	UserIDKey contextKey = "userID"
	RoleKey   contextKey = "role"
	PagesKey  contextKey = "pages"
)

// New creates a middleware that injects the auth service into the request context
func New(db postgres.Executor) func(http.Handler) http.Handler {
	authService := auth.New(db)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "auth", authService)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Instance retrieves the auth service from context
func Instance(ctx context.Context) auth.AuthService {
	return ctx.Value("auth").(auth.AuthService)
}

// GetUserID retrieves the user ID from context
func GetUserID(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(UserIDKey).(int)
	return id, ok
}

// GetEmail retrieves the email from context
func GetEmail(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(EmailKey).(string)
	return email, ok
}

// GetRole retrieves the role from context
func GetRole(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(RoleKey).(string)
	return role, ok
}

// Authenticate validates the JWT token and extracts user information
func Authenticate(jwtService jwt.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			parts := strings.Fields(authHeader)

			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				helper.UnauthorizedResponse(w, r)
				return
			}

			tokenString := strings.Trim(parts[1], "\"")

			claims, err := jwtService.ParseToken(tokenString)
			if err != nil {
				helper.UnauthorizedResponse(w, r)
				return
			}

			// Extract email
			email, ok := claims["email"].(string)
			if !ok {
				helper.UnauthorizedResponse(w, r)
				return
			}

			// Extract user ID
			userIDFloat, ok := claims["id"].(float64)
			if !ok {
				helper.UnauthorizedResponse(w, r)
				return
			}
			userID := int(userIDFloat)

			// Extract role
			role, ok := claims["role"].(string)
			if !ok {
				helper.UnauthorizedResponse(w, r)
				return
			}

			// Extract pages/permissions (optional)
			pages, _ := claims["pages"].([]interface{})

			// Add values to context
			ctx := r.Context()
			ctx = context.WithValue(ctx, EmailKey, email)
			ctx = context.WithValue(ctx, UserIDKey, userID)
			ctx = context.WithValue(ctx, RoleKey, role)
			ctx = context.WithValue(ctx, PagesKey, pages)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Authorize checks if the user has permission to access the route
func Authorize(jwtService jwt.JWTService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Parse JWT token to extract user permissions
			tokenData, err := jwtService.ParseTokenFromRequest(r)
			if err != nil {
				helper.UnauthorizedResponse(w, r)
				return
			}

			// Extract user permissions from token
			pagesInterface, ok := tokenData["pages"].([]interface{})
			if !ok {
				helper.ForbiddenResponse(w, r)
				return
			}

			route := chi.RouteContext(r.Context()).RoutePattern()

			// Match route with permissions
			for _, page := range pagesInterface {
				pageData, ok := page.(map[string]interface{})
				if !ok {
					continue
				}

				routeName, ok := pageData["route_name"].(string)
				if !ok {
					continue
				}

				permissions, ok := pageData["permissions"].(map[string]interface{})
				if !ok {
					continue
				}

				// Check if the current route matches
				// Route format: /v1/sales-orders/list, routeName format: /sales-orders
				// We need to check if route contains the routeName
				if strings.Contains(route, routeName) || strings.HasPrefix(route, routeName) {
					// Authorization based on HTTP method and route action
					switch {
					case r.Method == http.MethodPost && getBool(permissions, "can_create"):
						next.ServeHTTP(w, r)
						return
					case r.Method == http.MethodGet && getBool(permissions, "can_view"):
						next.ServeHTTP(w, r)
						return
					case r.Method == http.MethodPut && getBool(permissions, "can_update"):
						next.ServeHTTP(w, r)
						return
					case r.Method == http.MethodPatch && getBool(permissions, "can_update"):
						next.ServeHTTP(w, r)
						return
					case r.Method == http.MethodDelete && getBool(permissions, "can_delete"):
						next.ServeHTTP(w, r)
						return
					}
				}
			}

			helper.ForbiddenResponse(w, r)
		})
	}
}

// AuthorizeRoles checks if the user has one of the specified roles
func AuthorizeRoles(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := GetRole(r.Context())
			if !ok {
				helper.UnauthorizedResponse(w, r)
				return
			}

			for _, allowedRole := range allowedRoles {
				if strings.EqualFold(role, allowedRole) {
					next.ServeHTTP(w, r)
					return
				}
			}

			helper.ForbiddenResponse(w, r)
		})
	}
}

// Helper functions

func containsAction(route, action string) bool {
	return strings.Contains(strings.ToLower(route), action)
}

func getBool(m map[string]interface{}, key string) bool {
	if val, ok := m[key].(bool); ok {
		return val
	}
	return false
}
