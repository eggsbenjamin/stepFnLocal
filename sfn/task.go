//go:generate mockgen -package sfn -source=task.go -destination task_mock.go

package sfn

import (
	"github.com/aws/aws-sdk-go/aws/arn"
	awslambda "github.com/aws/aws-sdk-go/service/lambda"
	"github.com/eggsbenjamin/stepFnLocal/lambda"
	"github.com/eggsbenjamin/stepFnLocal/state"
	"github.com/pkg/errors"
)

type StateFactory interface {
	Create(state.Definition) (state.State, bool)
}

type stateFactory struct{}

func (stateFactory) Create(def state.Definition) (state.State, bool) {
	return nil, false
}

type LambdaTask struct {
	definition   state.TaskDefinition
	arn          arn.ARN
	lambdaClient lambda.Client
}

func NewLambdaTask(def state.TaskDefinition, arn arn.ARN, lambdaClient lambda.Client) state.State {
	return LambdaTask{
		definition:   def,
		arn:          arn,
		lambdaClient: lambdaClient,
	}
}

func (l LambdaTask) Run(input []byte) ([]byte, error) {
	arnStr := l.arn.String()

	// TODO: call the lambda async using InvocationType: "Event"... how to get response...??? (poll cloudwatch?)
	invokeOutput, err := l.lambdaClient.Invoke(&awslambda.InvokeInput{
		FunctionName: &arnStr,
		Payload:      input,
	})
	if err != nil {
		return nil, errors.Wrap(err, "lambda client error")
	}
	if invokeOutput.FunctionError != nil {
		return nil, state.NewError(invokeOutput.Payload)
	}

	return invokeOutput.Payload, nil
}

func (l LambdaTask) Next() string {
	return l.definition.Next()
}

func (l LambdaTask) IsEnd() bool {
	return l.definition.End()
}
