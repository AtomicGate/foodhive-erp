package login

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/anas-dev-92/FoodHive/core/jwt"
	"github.com/anas-dev-92/FoodHive/core/postgres"
	helper "github.com/anas-dev-92/FoodHive/registration/src/v1/utils/helper"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  UserDetails `json:"user"`
}

type UserDetails struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

// Handler handles user login requests.
// @Summary User login
// @Description Authenticates a user and returns a JWT token with permissions
// @Tags Authentication
// @Accept json
// @Produce json
// @Param login body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse "Login successful"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 401 {string} string "Invalid email or password"
// @Failure 500 {string} string "Internal server error"
// @Router /login [post]
func Handler(service jwt.JWTService, db postgres.Executor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := helper.ReadJSON(w, r, &req); err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		// Validate request
		v := helper.New()
		v.Check(req.Email != "", "email", "must be provided")
		v.Check(req.Password != "", "password", "must be provided")
		if !v.Valid() {
			helper.FailedValidationResponse(w, r, v.Errors)
			return
		}

		var storedPassword, email, englishName, roleName string
		var id, roleID int

		// Retrieve the stored password, email, and role
		err := db.QueryRow(r.Context(), `
			SELECT e.id, e.email, e.password, COALESCE(e.english_name, ''), COALESCE(e.role_id, 0), COALESCE(r.role_name, '')
			FROM employees e
			LEFT JOIN roles r ON e.role_id = r.id
			WHERE e.email = $1 AND e.status = 'CONTINUED'
		`, req.Email).Scan(&id, &email, &storedPassword, &englishName, &roleID, &roleName)
		if err != nil {
			log.Printf("Error fetching user: %v", err)
			helper.UnauthorizedResponse(w, r)
			return
		}

		// Compare passwords (use bcrypt in production!)
		// TODO: Replace with proper password hashing
		if storedPassword != req.Password {
			log.Printf("Invalid password for user: %s", email)
			helper.UnauthorizedResponse(w, r)
			return
		}

		// Fetch routes with permissions
		rows := db.Query(r.Context(), `
			SELECT p.route_name, ep.can_create, ep.can_update, ep.can_delete, ep.can_view
			FROM pages p
			JOIN emp_page ep ON p.id = ep.page_id
			WHERE ep.user_id = $1
		`, id)
		defer rows.Close()

		var routes []map[string]interface{}
		for rows.Next() {
			var routeName string
			var canCreate, canUpdate, canDelete, canView bool
			if err := rows.Scan(&routeName, &canCreate, &canUpdate, &canDelete, &canView); err != nil {
				log.Printf("Error scanning route: %v", err)
				helper.ServerErrorResponse(w, r, err)
				return
			}
			routes = append(routes, map[string]interface{}{
				"route_name": routeName,
				"can_create": canCreate,
				"can_update": canUpdate,
				"can_delete": canDelete,
				"can_view":   canView,
			})
		}

		if err := rows.Err(); err != nil {
			log.Printf("Error reading routes: %v", err)
			helper.ServerErrorResponse(w, r, err)
			return
		}

		// Generate token with ID, email, role, and routes with permissions
		token, err := service.GenerateToken(id, email, roleName, routes)
		if err != nil {
			log.Printf("Error generating token: %v", err)
			helper.ServerErrorResponse(w, r, err)
			return
		}

		// Send the response
		response := LoginResponse{
			Token: token,
			User: UserDetails{
				ID:    id,
				Email: email,
				Name:  englishName,
				Role:  roleName,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
