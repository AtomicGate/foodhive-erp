package helper

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// ============================================
// Custom Date Type
// ============================================

type CustomDate time.Time

const dateFormat = "2006-01-02"

func (cd CustomDate) Format(layout string) string {
	return time.Time(cd).Format(layout)
}

func (cd CustomDate) IsZero() bool {
	return time.Time(cd).IsZero()
}

func (cd *CustomDate) Scan(value interface{}) error {
	if value == nil {
		*cd = CustomDate{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		*cd = CustomDate(v)
		return nil
	case string:
		parsedTime, err := time.Parse(dateFormat, v)
		if err != nil {
			return fmt.Errorf("failed to parse date string: %w", err)
		}
		*cd = CustomDate(parsedTime)
		return nil
	default:
		return fmt.Errorf("unsupported type for CustomDate: %T", value)
	}
}

func (cd CustomDate) Value() (driver.Value, error) {
	t := time.Time(cd)
	if t.IsZero() {
		return nil, nil
	}
	return t.Format(dateFormat), nil
}

func (cd *CustomDate) UnmarshalJSON(b []byte) error {
	str := strings.Trim(string(b), `"`)
	if str == "" || str == "null" {
		return nil
	}

	// Try parsing in "2006-01-02" format
	t, err := time.Parse(dateFormat, str)
	if err != nil {
		// Fallback to RFC3339 format
		t, err = time.Parse(time.RFC3339, str)
		if err != nil {
			return fmt.Errorf("failed to parse date: %w", err)
		}
	}

	*cd = CustomDate(t)
	return nil
}

func (cd CustomDate) MarshalJSON() ([]byte, error) {
	t := time.Time(cd)
	if t.IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + t.Format(dateFormat) + `"`), nil
}

// ============================================
// JSON Helpers
// ============================================

type Envelope map[string]any

// ReadJSON securely reads JSON from request body with size limits and validation
func ReadJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	maxBytes := 1_048_576 // 1MB limit
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	// Allow unknown fields - frontend may send extra fields that backend ignores
	// dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains invalid JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return fmt.Errorf("body contains incomplete JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return fmt.Errorf("body must not be empty")
		case strings.HasPrefix(err.Error(), "json: unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field")
			return fmt.Errorf("body contains unknown key %s", fieldName)
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}

	// Ensure only single JSON value
	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("body must contain only a single JSON value")
	}
	return nil
}

// WriteJSON writes a JSON response with proper headers
func WriteJSON(w http.ResponseWriter, status int, data Envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

// ============================================
// Response Helpers
// ============================================

func ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := Envelope{"error": message}
	err := WriteJSON(w, status, env, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func SuccessResponse(w http.ResponseWriter, r *http.Request, status int, data any) {
	env := Envelope{"data": data}
	err := WriteJSON(w, status, env, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	ErrorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Bad request: %v", err)
	ErrorResponse(w, r, http.StatusBadRequest, err.Error())
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, r, http.StatusNotFound, "the requested resource could not be found")
}

func ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Server error: %v", err)
	ErrorResponse(w, r, http.StatusInternalServerError, "the server encountered a problem and could not process your request")
}

func UnauthorizedResponse(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, r, http.StatusUnauthorized, "unauthorized")
}

func ForbiddenResponse(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, r, http.StatusForbidden, "you don't have permission to access this resource")
}

func ConflictResponse(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, r, http.StatusConflict, "unable to update the record due to an edit conflict, please try again")
}

func MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, r, http.StatusMethodNotAllowed, fmt.Sprintf("the %s method is not supported for this resource", r.Method))
}

func RateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, r, http.StatusTooManyRequests, "rate limit exceeded")
}

// CreatedResponse sends a 201 response with the created resource ID
func CreatedResponse(w http.ResponseWriter, r *http.Request, id int, message string) {
	env := Envelope{
		"id":      id,
		"message": message,
	}
	WriteJSON(w, http.StatusCreated, env, nil)
}
