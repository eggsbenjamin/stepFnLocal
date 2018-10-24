//go:generate mockgen -package sfn -source=factory.go -destination factory_mock.go

package sfn

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/eggsbenjamin/stepFnLocal/lambda"
	state "github.com/eggsbenjamin/stepFnLocal/state"
	"github.com/pkg/errors"
)

type StateFactory interface {
	Create(state.Definition) (State, error)
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

func (s stateFactory) Create(def state.Definition) (State, error) {
	switch def.Type() {
	case state.TaskStateType:
		taskDef, ok := def.(state.TaskDefinition)
		if !ok {
			return nil, errors.New("invalid task state definition")
		}
		return s.createTaskState(taskDef)
	case state.ChoiceStateType:
		choiceDef, ok := def.(state.ChoiceDefinition)
		if !ok {
			return nil, errors.New("invalid choice state definition")
		}
		return s.createChoiceState(choiceDef)
	case state.PassStateType:
		passDef, ok := def.(state.PassDefinition)
		if !ok {
			return nil, errors.New("invalid pass state definition")
		}
		return NewPassState(passDef), nil
	case state.SucceedStateType:
		succeedDef, ok := def.(state.SucceedDefinition)
		if !ok {
			return nil, errors.New("invalid succeed state definition")
		}
		return NewSucceedState(succeedDef), nil
	case state.FailStateType:
		failDef, ok := def.(state.FailDefinition)
		if !ok {
			return nil, errors.New("invalid fail state definition")
		}
		return NewFailState(failDef), nil
	case state.ParallelStateType:
		parallelDef, ok := def.(state.ParallelDefinition)
		if !ok {
			return nil, errors.New("invalid parallel state definition")
		}
		return s.createParallelState(parallelDef)
	}

	return nil, state.ErrUnknownState
}

// TODO: these shouldn't be methods exposed by the factory, instad they should be dependencies (see abstract factory pattern)
func (s stateFactory) createTaskState(def state.TaskDefinition) (State, error) {
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

func (s stateFactory) createChoiceState(def state.ChoiceDefinition) (State, error) {
	choiceRuleFactory := NewChoiceRuleFactory()
	choiceRules := []ChoiceRule{}
	for _, choiceRuleDef := range def.Choices {
		choiceRule, err := choiceRuleFactory.Create(choiceRuleDef)
		if err != nil {
			return nil, errors.Wrap(err, "error creating choice state")
		}

		choiceRules = append(choiceRules, choiceRule)
	}

	return NewChoiceState(def, choiceRules...), nil
}

func (s stateFactory) createParallelState(def state.ParallelDefinition) (State, error) {
	stateMachines := []StepFunction{}

	for _, branchDef := range def.Branches {
		stateMachine, err := NewWithAWSConfig(branchDef, s.overrides, &aws.Config{}) // TODO: use same config
		if err != nil {
			return nil, errors.Wrap(err, "error creating parallel state branch")
		}
		stateMachines = append(stateMachines, stateMachine)
	}

	return NewParallelState(def, stateMachines...), nil
}
