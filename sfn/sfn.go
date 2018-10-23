package sfn

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/eggsbenjamin/stepFnLocal/state"
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
	fmt.Printf("running state: %s\n", stateTitle)
	def, err := r.stateMachineDef.States.GetDefinition(stateTitle)
	if err != nil {
		return []byte{}, err
	}

	_state, err := r.stateFactory.Create(def)
	if err != nil {
		return []byte{}, err
	}

	if v, ok := def.(state.InputPather); ok {
		input, err = v.InputPath().Search(input)
		if err != nil {
			return []byte{}, err
		}
	}

	output, err := _state.Run(input)
	if err != nil {
		return []byte{}, err
	}

	/*
		TODO: add support for ResultPath

		ResultPath modifies the output of a state using jsonpath.

		e.g.

		output = {"hello":"world"}
		ResultPath = "$.test"
		modified = {"test":{"hello":"world"}}
	*/

	if v, ok := def.(state.OutputPather); ok {
		output, err = v.OutputPath().Search(output)
		if err != nil {
			return []byte{}, err
		}
	}

	if _state.IsEnd() {
		return output, nil
	}

	return r.run(_state.Next(), output)
}
