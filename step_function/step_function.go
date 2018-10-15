package stepfunction

import (
	"encoding/json"

	"github.com/eggsbenjamin/stepFnLocal/state"
)

type StepFunction interface {
	StartExecution(json.RawMessage) error
}

type stepFunction struct {
	stateMachineDef state.MachineDefinition
}

func New(stateMachineDef state.MachineDefinition) (StepFunction, error) {
	// validate definition in ctr

	return &stepFunction{
		stateMachineDef: stateMachineDef,
	}, nil
}

func (r *stepFunction) run(stateTitle string, input json.RawMessage) error {
	stateDef, err := r.stateMachineDef.States.GetDefinition(stateTitle)
	if err != nil {
		return err
	}

	var state state.State
	switch stateDef.Type() {
	case state.TaskStateType:
		def, ok := stateDef.(state.TaskDefinition)
		state = state.NewTask(def)
	}

	output, err := state.Run(input)
	if err != nil {
		return err
	}

	if state.End() {
		return nil
	}

	return run(state.Next(), output)
}
