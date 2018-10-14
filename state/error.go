package state

import "errors"

const (
	MissingRequiredFieldErrType = "Missing required field"
	InvalidValueErrType         = "Invalid Value"
)

var (
	ErrStateNotFound = errors.New("state not found")
	ErrUnknownState  = errors.New("unknown state")
)

// ValidationError represents a single AWS states language validation error.
type ValidationError struct {
	Type  string
	Field string
	Value string
}

// Error implements the error interface for ValidationError
func (v ValidationError) Error() string {
	str := v.Type + " '" + v.Field + "'"
	if v.Value != "" {
		str += ": '" + v.Value + "'"
	}

	return str
}

// NewValidationError instantiates a ValidationError
func NewValidationError(typ, field, value string) *ValidationError {
	return &ValidationError{
		Type:  typ,
		Field: field,
		Value: value,
	}
}

// ValidationErrors represents zero or more AWS states language validation errors.
type ValidationErrors []error

// Error implements the error interface for ValidationErrors
func (v ValidationErrors) Error() string {
	var str string
	for i, err := range []error(v) {
		str += err.Error()
		if i != len(v)-1 {
			str += " : "
		}
	}

	return str
}
