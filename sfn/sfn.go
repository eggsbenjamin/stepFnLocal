package sfn

import (
	"encoding/json"
	"time"

	"github.com/eggsbenjamin/stepFnLocal/state"
	"github.com/pkg/errors"
)

const (
	ExecutionStatusRunning   = "RUNNING"
	ExecutionStatusSucceeded = "SUCCEEDED"
	ExecutionStatusFailed    = "FAILED"
	ExecutionStatusTimedOut  = "TIMED_OUT"
	ExecutionStatusAborted   = "ABORTED"
)

type StepFunction interface {
	StartExecution([]byte) (state.ExecutionResult, error)
}

type stepFunction struct {
	stateMachineDef state.MachineDefinition
	stateFactory    StateFactory
}

func New(stateMachineDef state.MachineDefinition, stateFactory StateFactory) (StepFunction, error) {
	if err := stateMachineDef.Validate(); err != nil {
		return stepFunction{}, err
	}

	return stepFunction{
		stateMachineDef: stateMachineDef,
		stateFactory:    stateFactory,
	}, nil
}

func (s stepFunction) StartExecution(input []byte) (state.ExecutionResult, error) {
	start := time.Now()

	status := ExecutionStatusSucceeded
	output, err := s.run(s.stateMachineDef.StartAt, input)
	if err != nil {
		status = ExecutionStatusFailed
	}

	return state.ExecutionResult{
		Input:  input,
		Output: output,
		Status: status,
		Start:  start,
		End:    time.Now(),
	}, err
}

func (r stepFunction) run(stateTitle string, input json.RawMessage) ([]byte, error) {
	def, err := r.stateMachineDef.States.GetDefinition(stateTitle)
	if err != nil {
		return []byte{}, err
	}

	var _state state.State
	switch def.Type() {
	case state.TaskStateType:
		task, ok := r.stateFactory.Create(def)
		if !ok {
			return []byte{}, errors.New("unknown state")
		}
		_state = task
	}

	output, err := _state.Run(input)
	if err != nil {
		return []byte{}, err
	}

	if _state.IsEnd() {
		return output, nil
	}

	return r.run(_state.Next(), output)
}
