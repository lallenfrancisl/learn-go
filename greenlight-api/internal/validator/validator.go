package validator

import "regexp"

var (
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

type Validator struct {
	Errors ValidationErrors
}

type ValidationErrors map[string]string

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// Valid returns true if the errors map doesn't contain any entries
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// Add an error message to the map (as long as no entry already exists for the given key)
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Check adds an error message to the map only if a validation check is not 'ok'
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// In returns true if a specific value is in a list of strings
func In(value string, list ...string) bool {
	for i := range list {
		if value == list[i] {
			return true
		}
	}

	return false
}

// Matches returns true if a string value matches a specific regexp pattern
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// Unique returns true if all string values in a slice are unique
func Unique(values []string) bool {
	uniqueValues := make(map[string]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(values) == len(uniqueValues)
}

// Returns true if value is not blank
func NotBlank(value string) bool {
	return value != ""
}

// Checks if the value is nil
func NotNil(value interface{}) bool {
	return value != nil
}

// Returns true when `value` is less than or equal to `limit`
func Max(value, limit int) bool {
	return value <= limit
}

// Returns true when `value` is more than or equal to `limit`
func Min(value, limit int) bool {
	return value >= limit
}

// Returns true when `value` is more than `limit`
func GreaterThan(value, limit int) bool {
	return value > limit
}

// Returns true when `value` is less than `limit`
func LessThan(value, limit int) bool {
	return value < limit
}

// Returns true when length of `value` is less than or equal to limit
func MaxLen(value string, limit int) bool {
	return len(value) <= limit
}

// Returns true when length of `value` is more than or equal to limit
func MinLen(value string, limit int) bool {
	return len(value) >= limit
}

// Returns true when `value` is equal to `expected`
func Equal(value, expected string) bool {
	return value == expected
}
