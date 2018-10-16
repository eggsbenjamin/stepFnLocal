//go:generate mockgen -package sfn -source=factory.go -destination factory_mock.go

package sfn

import (
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/eggsbenjamin/stepFnLocal/lambda"
	state "github.com/eggsbenjamin/stepFnLocal/state"
	"github.com/pkg/errors"
)

type StateFactory interface {
	Create(state.Definition) (state.State, error)
}

type stateFactory struct {
	overrides    map[string]OverrideFn
	lambdaClient lambda.Client
}

func NewStateFactory(overrides map[string]OverrideFn, lambdaClient lambda.Client) StateFactory {
	return stateFactory{
		overrides:    overrides,
		lambdaClient: lambdaClient,
	}
}

func (s stateFactory) Create(def state.Definition) (state.State, error) {
	switch def.Type() {
	case state.TaskStateType:
		taskDef, ok := def.(state.TaskDefinition)
		if !ok {
			return nil, errors.New("invalid task state definition")
		}
		return s.createTaskState(taskDef)
	}

	return nil, state.ErrUnknownState
}

func (s stateFactory) createTaskState(def state.TaskDefinition) (state.State, error) {
	if overrideFn, ok := s.overrides[def.Resource]; ok {
		return NewOverrideTask(def, overrideFn), nil
	}

	arn, err := arn.Parse(def.Resource)
	if err != nil {
		return nil, errors.Wrap(err, "invalid arn")
	}

	// TODO: add support for aws resources other than lambda
	return NewLambdaTask(def, arn, s.lambdaClient), nil
}
