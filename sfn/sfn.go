//go:generate mockgen -package sfn -source=sfn.go -destination sfn_mock.go

package sfn

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/eggsbenjamin/stepFnLocal/state"
)

const (
	ExecutionStatusRunning   = "RUNNING"
	ExecutionStatusSucceeded = "SUCCEEDED"
	ExecutionStatusFailed    = "FAILED"
	ExecutionStatusTimedOut  = "TIMED_OUT"
	ExecutionStatusAborted   = "ABORTED"
)

// ExecutionResult represents the result of a state machine execution
type ExecutionResult struct {
	Input  []byte
	Output []byte
	Status string
	Start  time.Time
	End    time.Time
}

// State defines the standard state API for state machine implementations
type State interface {
	Run([]byte) ([]byte, error)
	Next() string
	IsEnd() bool
}

type StepFunction interface {
	StartExecution([]byte) (ExecutionResult, error)
	SetStateFactory(StateFactory)
}

type stepFunction struct {
	stateMachineDef state.MachineDefinition
	stateFactory    StateFactory
}

func New(def state.MachineDefinition, overrides map[string]OverrideFn) (StepFunction, error) {
	if err := def.Validate(); err != nil {
		return &stepFunction{}, err
	}

	lambdaClient := lambda.New(session.Must(session.NewSession(&aws.Config{})))
	stateFactory := NewStateFactory(overrides, lambdaClient)

	return &stepFunction{
		stateMachineDef: def,
		stateFactory:    stateFactory,
	}, nil
}

func NewWithAWSConfig(def state.MachineDefinition, overrides map[string]OverrideFn, awsCfg *aws.Config) (StepFunction, error) {
	if err := def.Validate(); err != nil {
		return &stepFunction{}, err
	}

	lambdaClient := lambda.New(session.Must(session.NewSession(awsCfg)))
	stateFactory := NewStateFactory(overrides, lambdaClient)

	return &stepFunction{
		stateMachineDef: def,
		stateFactory:    stateFactory,
	}, nil
}

func (s *stepFunction) StartExecution(input []byte) (ExecutionResult, error) {
	start := time.Now()

	status := ExecutionStatusSucceeded
	output, err := s.run(s.stateMachineDef.StartAt, input)
	if err != nil {
		status = ExecutionStatusFailed
	}

	return ExecutionResult{
		Input:  input,
		Output: output,
		Status: status,
		Start:  start,
		End:    time.Now(),
	}, err
}

func (s *stepFunction) SetStateFactory(stateFactory StateFactory) {
	s.stateFactory = stateFactory
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
