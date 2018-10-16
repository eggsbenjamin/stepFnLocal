//go:generate mockgen -package sfn -source=state.go -destination state_mock.go

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
	lambdaClient lambda.Client
}

func NewStateFactory(lambdaClient lambda.Client) StateFactory {
	return stateFactory{
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
	arn, err := arn.Parse(def.Resource)
	if err != nil {
		return nil, errors.Wrap(err, "invalid arn")
	}

	return NewLambdaTask(def, arn, s.lambdaClient), nil
}
