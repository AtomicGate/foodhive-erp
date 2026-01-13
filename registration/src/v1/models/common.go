package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// ============================================
// Custom Types
// ============================================

// CustomDate handles date serialization
type CustomDate time.Time

func (cd *CustomDate) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = s[1 : len(s)-1] // Remove quotes
	if s == "" || s == "null" {
		return nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*cd = CustomDate(t)
	return nil
}

func (cd CustomDate) MarshalJSON() ([]byte, error) {
	t := time.Time(cd)
	if t.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(t.Format("2006-01-02"))
}

func (cd *CustomDate) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		*cd = CustomDate(v)
		return nil
	default:
		return errors.New("invalid type for CustomDate")
	}
}

func (cd CustomDate) Value() (driver.Value, error) {
	t := time.Time(cd)
	if t.IsZero() {
		return nil, nil
	}
	return t, nil
}

func (cd CustomDate) IsZero() bool {
	return time.Time(cd).IsZero()
}

// CustomDateTime handles datetime serialization
type CustomDateTime time.Time

func (cdt *CustomDateTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = s[1 : len(s)-1]
	if s == "" || s == "null" {
		return nil
	}
	t, err := time.Parse("2006-01-02T15:04:05Z07:00", s)
	if err != nil {
		t, err = time.Parse("2006-01-02 15:04:05", s)
		if err != nil {
			return err
		}
	}
	*cdt = CustomDateTime(t)
	return nil
}

func (cdt CustomDateTime) MarshalJSON() ([]byte, error) {
	t := time.Time(cdt)
	if t.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(t.Format("2006-01-02T15:04:05Z07:00"))
}

func (cdt *CustomDateTime) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		*cdt = CustomDateTime(v)
		return nil
	default:
		return errors.New("invalid type for CustomDateTime")
	}
}

func (cdt CustomDateTime) Value() (driver.Value, error) {
	t := time.Time(cdt)
	if t.IsZero() {
		return nil, nil
	}
	return t, nil
}

func (cdt CustomDateTime) IsZero() bool {
	return time.Time(cdt).IsZero()
}

// ============================================
// Pagination
// ============================================

type Pagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

// ============================================
// Common Response Types
// ============================================

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Error   string            `json:"error"`
	Details map[string]string `json:"details,omitempty"`
}

type IDResponse struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
}

// ============================================
// Validation
// ============================================

type Validator struct {
	Errors map[string]string
}

func NewValidator() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}
