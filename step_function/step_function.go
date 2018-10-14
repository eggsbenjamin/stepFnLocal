package stepfunction

import (
	"encoding/json"
	"errors"

	"github.com/eggsbenjamin/stepFnLocal/state"
)

type StepFunction interface {
	StartExecution(json.RawMessage) error
}

type stepFunction struct {
	stateMachine *state.MachineDefinition
	resourceMap  map[string]string
}

func New(stateMachine *state.MachineDefinition, resourceMap map[string]string) StepFunction {
	return &stepFunction{
		stateMachine: stateMachine,
		resourceMap:  resourceMap,
	}
}

func (r *stepFunction) StartExecution(input json.RawMessage) error {
	if r.stateMachine.StartAt == "" {
		return errors.New("missing StartAt")
	}
	return nil
}
