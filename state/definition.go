package state

import (
	"github.com/eggsbenjamin/stepFnLocal/jsonpath"
)

const (
	PassStateType     = "Pass"
	TaskStateType     = "Task"
	ChoiceStateType   = "Choice"
	WaitStateType     = "Wait"
	SucceedStateType  = "Succeed"
	FailStateType     = "Fail"
	ParallelStateType = "Parallel"
)

var validTypes = map[string]struct{}{
	PassStateType:     {},
	TaskStateType:     {},
	ChoiceStateType:   {},
	WaitStateType:     {},
	SucceedStateType:  {},
	FailStateType:     {},
	ParallelStateType: {},
}

type JSONPathExp string

func (j JSONPathExp) Validate() error {
	_, err := jsonpath.NewExpression(string(j))
	return err
}

func (j JSONPathExp) Search(input []byte) ([]byte, error) {
	if string(j) == "" {
		return input, nil
	}

	exp, err := jsonpath.NewExpression(string(j))
	if err != nil {
		return []byte{}, err
	}

	return exp.Search(input)
}

// TaskDefinition represents an AWS states language task state.
type TaskDefinition struct {
	BaseDefinition
	TransitionDefinition
	IOPathDefinition
	ResultPathDefinition
	Resource string `json:"Resource"`
	/*
		TimeoutSeconds   int    `json:"TimeoutSeconds"`
		HeartbeatSeconds int    `json:"HeartbeatSeconds"`
	*/
}

func (t TaskDefinition) Type() string {
	return TaskStateType
}

func (t TaskDefinition) Validate() error {
	validationErrs := ValidationErrors{}

	if err := t.BaseDefinition.Validate(); err != nil {
		validationErrs = append(validationErrs, err.(ValidationErrors)...)
	}

	if err := t.TransitionDefinition.Validate(); err != nil {
		validationErrs = append(validationErrs, err.(ValidationErrors)...)
	}

	if err := t.IOPathDefinition.Validate(); err != nil {
		validationErrs = append(validationErrs, err.(ValidationErrors)...)
	}

	if err := t.ResultPathDefinition.Validate(); err != nil {
		validationErrs = append(validationErrs, err.(ValidationErrors)...)
	}

	if t.Resource == "" {
		validationErrs = append(validationErrs, NewValidationError(
			MissingRequiredFieldErrType,
			"Resource", "",
		))
	}

	if t.InputPathExp != "" {
		if err := t.InputPathExp.Validate(); err != nil {
			validationErrs = append(validationErrs, NewValidationError(
				InvalidJSONPathErrType,
				"InputPath", string(t.InputPathExp),
			))
		}
	}

	if t.OutputPathExp != "" {
		if err := t.OutputPathExp.Validate(); err != nil {
			validationErrs = append(validationErrs, NewValidationError(
				InvalidJSONPathErrType,
				"OutputPath", string(t.OutputPathExp),
			))
		}
	}

	if t.ResultPathExp != "" {
		if err := t.ResultPathExp.Validate(); err != nil {
			validationErrs = append(validationErrs, NewValidationError(
				InvalidJSONPathErrType,
				"ResultPath", string(t.ResultPathExp),
			))
		}
	}

	if len(validationErrs) > 0 {
		return validationErrs
	}
	return nil
}

// RetryDefinition represents an AWS states language retry block
type RetryDefinition struct {
	ErrorEquals     []string `json:"ErrorEquals"`
	IntervalSeconds int      `json:"IntervalSeconds"`
	MaxAttempts     int      `json:"MaxAttempts"`
	BackoffRate     float64  `json:"BackoffRate"`
}

type CatchDefinition struct {
	ErrorEquals []string `json:"ErrorEquals"`
	ResultPath  string   `json:"ResultPath"`
	Next        string   `json:"Next"`
}
