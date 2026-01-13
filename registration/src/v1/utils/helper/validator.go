package helper

import (
	"regexp"
	"slices"
)

var (
	// EmailRX is a regex for validating email addresses
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	// PhoneRX validates phone numbers (basic validation)
	PhoneRX = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)

	// SKURx validates SKU format (alphanumeric with dashes)
	SKURX = regexp.MustCompile(`^[A-Za-z0-9\-]+$`)
)

// ValidationErrors is a map of field names to error messages
type ValidationErrors map[string]string

// Validator provides validation functionality
type Validator struct {
	Errors map[string]string
}

// New creates a new Validator instance
func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// Valid returns true if there are no validation errors
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError adds an error message for a specific field (only if not already present)
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Check adds an error if the condition is false
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// PermittedValue checks if a value is in a list of permitted values
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

// Matches checks if a string matches a regex pattern
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// Unique checks if all values in a slice are unique
func Unique[T comparable](values []T) bool {
	uniqueValues := make(map[T]bool)
	for _, value := range values {
		uniqueValues[value] = true
	}
	return len(values) == len(uniqueValues)
}

// NotBlank checks if a string is not empty after trimming whitespace
func NotBlank(value string) bool {
	return len(value) > 0
}

// MinLength checks if a string has at least n characters
func MinLength(value string, n int) bool {
	return len(value) >= n
}

// MaxLength checks if a string has at most n characters
func MaxLength(value string, n int) bool {
	return len(value) <= n
}

// Between checks if a value is between min and max (inclusive)
func Between[T ~int | ~int64 | ~float64](value, min, max T) bool {
	return value >= min && value <= max
}

// PositiveInt checks if an integer is positive
func PositiveInt(value int) bool {
	return value > 0
}

// PositiveFloat checks if a float is positive
func PositiveFloat(value float64) bool {
	return value > 0
}

// NonNegativeFloat checks if a float is non-negative
func NonNegativeFloat(value float64) bool {
	return value >= 0
}

// ValidEmail checks if a string is a valid email address
func ValidEmail(value string) bool {
	return EmailRX.MatchString(value)
}

// ValidPhone checks if a string is a valid phone number
func ValidPhone(value string) bool {
	return PhoneRX.MatchString(value)
}

// ValidSKU checks if a string is a valid SKU
func ValidSKU(value string) bool {
	return SKURX.MatchString(value)
}
