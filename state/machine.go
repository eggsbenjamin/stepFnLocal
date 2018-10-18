package state

import (
	"encoding/json"
	"github.com/pkg/errors"
)

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
