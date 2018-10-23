package state

import (
	"encoding/json"
	"errors"
)

const (
	// ValidationErr Types
	MissingRequiredFieldErrType = "Missing required field"
	InvalidKeyErrType           = "Invalid Key"
	InvalidValueErrType         = "Invalid Value"
	InvalidJSONPathErrType      = "Invalid JSON path expression"
	InvalidCombinationErrType   = "Invalid Combination"
	NonRFC3339TimeStampErrType  = "Non RFC3339 timestamp"

	OnlyOneMustExistErrMsg = "Only one must exist"
)

var (
	// states language errors
	ErrAll                    = errors.New("States.ALL")
	ErrTimeout                = errors.New("States.Timeout")
	ErrTaskFailed             = errors.New("States.TaskFailed")
	ErrTaskPermissions        = errors.New("States.TaskPermissions")
	ErrResultPathMatchFailure = errors.New("States.ResultPathMatchFailure")
	ErrBranchFailed           = errors.New("States.BranchFailed")
	ErrNoChoiceMatched        = errors.New("States.NoChoiceMatched")

	// internal errors
	ErrStateNotFound = errors.New("state not found")
	ErrUnknownState  = errors.New("unknown state")
)

type Error struct {
	Message json.RawMessage
}

func (e Error) Error() string {
	return string(e.Message)
}

func NewError(msg json.RawMessage) Error {
	return Error{
		Message: msg,
	}
}

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
