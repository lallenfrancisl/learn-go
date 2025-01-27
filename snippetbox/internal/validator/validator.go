package validator

import (
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FormErrors  []string
	FieldErrors map[string]string
}

// Valid returns true if the FieldErrors doesn't contain
// any entries
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.FormErrors) == 0
}

func (v *Validator) AddFieldError(key string, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) AddFormError(message string) {
	v.FormErrors = append(v.FormErrors, message)
}

func (v *Validator) CheckField(ok bool, key string, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func MinChars(value string, length int) bool {
	return utf8.RuneCountInString(value) >= length
}

func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

var EmailRX = regexp.MustCompile(
	"^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$",
)

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
