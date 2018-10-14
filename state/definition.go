package state

import (
	"encoding/json"
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

var types = map[string]struct{}{
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

	var stateDef *BaseStateDefinition
	if err := json.Unmarshal(rawState, &stateDef); err != nil {
		return nil, errors.Wrap(err, "error unmarshaling state json")
	}

	switch stateDef.StateType {
	case TaskStateType:
		var taskStateDef *TaskStateDefinition
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

// Definition defines the definition interface which all state definitions must implement
type Definition interface {
	Typer
	Validator
}

// BaseStateDefinition represents an AWS states language state. It contains fields that can appear in all state types.
type BaseStateDefinition struct {
	StateType    string `json:"Type"`
	StateComment string `json:"Comment"`
}

func (s BaseStateDefinition) Validate() error {
	validationErrs := ValidationErrors{}

	if s.StateType == "" {
		validationErrs = append(validationErrs, NewValidationError(MissingRequiredFieldErrType, "Type", ""))
	} else {
		if _, ok := types[s.StateType]; !ok {
			validationErrs = append(validationErrs, NewValidationError(InvalidValueErrType, "Type", s.StateType))
		}
	}

	if len(validationErrs) > 0 {
		return validationErrs
	}
	return nil
}

type TransitionStateDefinition struct {
	NextState string `json:"Next"`
	EndState  bool   `json:"End"`
}

func (t TransitionStateDefinition) Validate() error {
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

func (t TransitionStateDefinition) Next() string {
	return t.NextState
}

func (t TransitionStateDefinition) End() bool {
	return t.EndState
}

// TaskStateDefinition represents an AWS states language task state.
type TaskStateDefinition struct {
	BaseStateDefinition
	TransitionStateDefinition
	Resource string `json:"Resource"`
	/*
		InputPath        string `json:"InputPath"`
		OutputPath       string `json:"OutputPath"`
		ResultPath       string `json:"ResultPath"`
		TimeoutSeconds   int    `json:"TimeoutSeconds"`
		HeartbeatSeconds int    `json:"HeartbeatSeconds"`
	*/
}

func (t TaskStateDefinition) Type() string {
	return TaskStateType
}

func (t TaskStateDefinition) Validate() error {
	validationErrs := ValidationErrors{}

	if err := t.BaseStateDefinition.Validate(); err != nil {
		validationErrs = append(validationErrs, err.(ValidationErrors)...)
	}

	if err := t.TransitionStateDefinition.Validate(); err != nil {
		validationErrs = append(validationErrs, err.(ValidationErrors)...)
	}

	if t.Resource == "" {
		validationErrs = append(validationErrs, NewValidationError(
			MissingRequiredFieldErrType,
			"Resource", "",
		))
	}

	if len(validationErrs) > 0 {
		return validationErrs
	}
	return nil
}

// RetryDefinition represents an AWS states language retry block
type RetryDefinition struct {
	ErrorEquals     []string `json:"ErrorEquals"`
	IntervalSeconds uint     `json:"IntervalSeconds"`
	MaxAttempts     uint     `json:"MaxAttempts"`
	BackoffRate     float64  `json:"BackoffRate"`
}

type CatchDefinition struct {
	ErrorEquals []string `json:"ErrorEquals"`
	ResultPath  string   `json:"ResultPath"`
	Next        string   `json:"Next"`
}
