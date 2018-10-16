//go:generate mockgen -package state -source=definition.go -destination definition_mock.go

package state

import (
	"encoding/json"

	"github.com/eggsbenjamin/stepFnLocal/jsonpath"
	"github.com/pkg/errors"
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
	PassStateType:    {},
	TaskStateType:    {},
	ChoiceStateType:  {},
	WaitStateType:    {},
	SucceedStateType: {},
	FailStateType:    {},
}

// MachineStates represents an AWS states language 'States' object.
type MachineStates map[string]json.RawMessage

func (m MachineStates) GetDefinition(name string) (Definition, error) {
	rawState, ok := m[name]
	if !ok {
		return nil, ErrStateNotFound
	}

	var stateDef *BaseDefinition
	if err := json.Unmarshal(rawState, &stateDef); err != nil {
		return nil, errors.Wrap(err, "error unmarshaling state json")
	}

	switch stateDef.StateType {
	case TaskStateType:
		var taskStateDef *TaskDefinition
		if err := json.Unmarshal(rawState, &taskStateDef); err != nil {
			return nil, errors.Wrap(err, "error unmarshaling task state json")
		}
		return *taskStateDef, nil
	}

	return nil, ErrUnknownState
}

// MachineDefinition represents an AWS states language state machine
type MachineDefinition struct {
	Comment        string        `json:"Comment"`
	StartAt        string        `json:"StartAt"`
	Version        string        `json:"Version"`
	TimeoutSeconds int           `json:"TimeoutSeconds"`
	States         MachineStates `json:"States"`
}

func (m MachineDefinition) Validate() error {
	validationErrs := ValidationErrors{}

	if m.StartAt == "" {
		validationErrs = append(validationErrs, NewValidationError(MissingRequiredFieldErrType, "StartAt", ""))
	} else {
		if _, ok := m.States[m.StartAt]; !ok {
			validationErrs = append(validationErrs, NewValidationError(InvalidValueErrType, "StartAt", m.StartAt))
		}
	}

	if m.States == nil {
		validationErrs = append(validationErrs, NewValidationError(MissingRequiredFieldErrType, "States", ""))
	}

	for title := range m.States {
		def, err := m.States.GetDefinition(title)
		if err != nil {
			return errors.Wrapf(err, "error getting state definiton for %s", title)
		}

		if err := def.Validate(); err != nil {
			validationErrs = append(validationErrs, err.(ValidationErrors)...)
		}

		transitioner, ok := def.(Transitioner)
		if !ok {
			continue
		}

		if transitioner.Next() != "" && !transitioner.End() {
			if _, ok = m.States[transitioner.Next()]; !ok {
				validationErrs = append(validationErrs, NewValidationError(
					InvalidValueErrType, "Next", transitioner.Next(),
				))
			}
		}
	}

	if len(validationErrs) > 0 {
		return validationErrs
	}
	return nil
}

// Typer defines the typer interface which all state definitions must implement
type Typer interface {
	Type() string
}

// Validator defines the validator interface which all state definitions must implement
type Validator interface {
	Validate() error
}

type Transitioner interface {
	Next() string
	End() bool
}

type InputPather interface {
	InputPath() JSONPathExp
}

type OutputPather interface {
	OutputPath() JSONPathExp
}

type IOPather interface {
	InputPather
	OutputPather
}

type ResultPather interface {
	ResultPath() JSONPathExp
}

// Definition defines the definition interface which all state definitions must implement
type Definition interface {
	Typer
	Validator
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

// BaseDefinition represents an AWS states language state. It contains fields that can appear in all state types.
type BaseDefinition struct {
	StateType    string `json:"Type"`
	StateComment string `json:"Comment"`
}

func (s BaseDefinition) Validate() error {
	validationErrs := ValidationErrors{}

	if s.StateType == "" {
		validationErrs = append(validationErrs, NewValidationError(MissingRequiredFieldErrType, "Type", ""))
	} else {
		if _, ok := validTypes[s.StateType]; !ok {
			validationErrs = append(validationErrs, NewValidationError(InvalidValueErrType, "Type", s.StateType))
		}
	}

	if len(validationErrs) > 0 {
		return validationErrs
	}
	return nil
}

type TransitionDefinition struct {
	NextState string `json:"Next"`
	EndState  bool   `json:"End"`
}

func (t TransitionDefinition) Validate() error {
	validationErrs := ValidationErrors{}

	if t.NextState == "" && t.EndState != true {
		validationErrs = append(validationErrs, NewValidationError(
			MissingRequiredFieldErrType,
			"Next/End:true", "",
		))
	}

	if len(validationErrs) > 0 {
		return validationErrs
	}
	return nil
}

func (t TransitionDefinition) Next() string {
	return t.NextState
}

func (t TransitionDefinition) End() bool {
	return t.EndState
}

type IODefinition struct {
	InputPathExp  JSONPathExp `json:"InputPath"`
	OutputPathExp JSONPathExp `json:"OutputPath"`
}

func (i IODefinition) Validate() error {
	validationErrs := ValidationErrors{}

	if i.InputPathExp != "" {
		if err := i.InputPathExp.Validate(); err != nil {
			validationErrs = append(validationErrs, NewValidationError(
				InvalidJSONPathErrType,
				"InputPath", string(i.InputPathExp),
			))
		}
	}

	if i.OutputPathExp != "" {
		if err := i.OutputPathExp.Validate(); err != nil {
			validationErrs = append(validationErrs, NewValidationError(
				InvalidJSONPathErrType,
				"OutputPath", string(i.OutputPathExp),
			))
		}
	}

	if len(validationErrs) > 0 {
		return validationErrs
	}
	return nil
}

func (i IODefinition) InputPath() JSONPathExp {
	return i.InputPathExp
}

func (i IODefinition) OutputPath() JSONPathExp {
	return i.OutputPathExp
}

type ResultDefinition struct {
	ResultPathExp JSONPathExp `json:"ResultPath"`
}

func (r ResultDefinition) Validate() error {
	validationErrs := ValidationErrors{}

	if r.ResultPathExp != "" {
		if err := r.ResultPathExp.Validate(); err != nil {
			validationErrs = append(validationErrs, NewValidationError(
				InvalidJSONPathErrType,
				"ResultPath", string(r.ResultPathExp),
			))
		}
	}

	if len(validationErrs) > 0 {
		return validationErrs
	}
	return nil
}

func (r ResultDefinition) ResultPath() JSONPathExp {
	return r.ResultPathExp
}

// TaskDefinition represents an AWS states language task state.
type TaskDefinition struct {
	BaseDefinition
	TransitionDefinition
	IODefinition
	ResultDefinition
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

	if err := t.IODefinition.Validate(); err != nil {
		validationErrs = append(validationErrs, err.(ValidationErrors)...)
	}

	if err := t.ResultDefinition.Validate(); err != nil {
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
