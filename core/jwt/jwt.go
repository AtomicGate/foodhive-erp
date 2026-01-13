package jwt

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	GenerateToken(ID int, email, role string, pages []map[string]interface{}) (string, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
	ParseToken(tokenString string) (map[string]interface{}, error)
	ParseTokenFromRequest(r *http.Request) (map[string]interface{}, error)
}

type JWTServiceImpl struct {
	secretKey string
}

type PageClaim struct {
	RouteName   string      `json:"route_name"`
	Permissions Permissions `json:"permissions"`
}

type Permissions struct {
	CanCreate bool `json:"can_create"`
	CanDelete bool `json:"can_delete"`
	CanUpdate bool `json:"can_update"`
	CanView   bool `json:"can_view"`
}

// New creates a new JWTService with the provided secret key.
func New(secretKey string) JWTService {
	return &JWTServiceImpl{
		secretKey: secretKey,
	}
}

// GenerateToken creates a JWT token with user details and permissions.
func (s *JWTServiceImpl) GenerateToken(ID int, email, role string, pages []map[string]interface{}) (string, error) {
	// Format pages to structured claims
	formattedPages := make([]PageClaim, len(pages))
	for i, page := range pages {
		formattedPages[i] = PageClaim{
			RouteName: page["route_name"].(string),
			Permissions: Permissions{
				CanCreate: page["can_create"].(bool),
				CanDelete: page["can_delete"].(bool),
				CanUpdate: page["can_update"].(bool),
				CanView:   page["can_view"].(bool),
			},
		}
	}

	// Create claims
	claims := jwt.MapClaims{
		"id":    ID,
		"email": email,
		"role":  role,
		"pages": formattedPages,
		"exp":   time.Now().Add(24 * time.Hour).Unix(), // Expiration time
	}

	// Generate token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key
	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates the provided token string.
func (s *JWTServiceImpl) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.secretKey), nil
	})
	return token, err
}

// ParseToken extracts claims from a valid token.
func (s *JWTServiceImpl) ParseToken(tokenString string) (map[string]interface{}, error) {
	// Validate token
	token, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Check for required claims
	requiredClaims := []string{"id", "email", "role", "pages"}
	for _, claim := range requiredClaims {
		if _, exists := claims[claim]; !exists {
			return nil, errors.New("missing required claims")
		}
	}

	// Return parsed claims as a map
	return map[string]interface{}{
		"id":    claims["id"],
		"email": claims["email"],
		"role":  claims["role"],
		"pages": claims["pages"],
	}, nil
}

// ParseTokenFromRequest extracts and parses a token from the Authorization header.
func (s *JWTServiceImpl) ParseTokenFromRequest(r *http.Request) (map[string]interface{}, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("authorization header missing")
	}

	// Extract Bearer token
	parts := strings.Fields(authHeader) // Splits by whitespace
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return nil, errors.New("invalid authorization header format")
	}

	// Parse token
	tokenString := parts[1]
	return s.ParseToken(tokenString)
}
