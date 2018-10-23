package sfn

import (
	"encoding/base64"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	awslambda "github.com/aws/aws-sdk-go/service/lambda"
	"github.com/eggsbenjamin/stepFnLocal/lambda"
	"github.com/eggsbenjamin/stepFnLocal/state"
	"github.com/pkg/errors"
)

type LambdaTask struct {
	definition   state.TaskDefinition
	arn          arn.ARN
	lambdaClient lambda.Client
}

func NewLambdaTask(def state.TaskDefinition, arn arn.ARN, lambdaClient lambda.Client) State {
	return LambdaTask{
		definition:   def,
		arn:          arn,
		lambdaClient: lambdaClient,
	}
}

func (l LambdaTask) Run(input []byte) ([]byte, error) {
	// TODO: call the lambda async using InvocationType: "Event"... how to get response...??? (poll cloudwatch?)
	invokeOutput, err := l.lambdaClient.Invoke(&awslambda.InvokeInput{
		FunctionName: aws.String(l.arn.String()),
		LogType:      aws.String("Tail"),
		Payload:      input,
	})
	if invokeOutput.LogResult != nil {
		logResult, err := base64.StdEncoding.DecodeString(*invokeOutput.LogResult)
		if err != nil {
			return nil, err
		}
		log.Printf("%s\n%s", l.arn.Resource, logResult)
	}
	if err != nil {
		return nil, errors.Wrap(err, "lambda client error")
	}
	if invokeOutput.FunctionError != nil {
		return nil, state.NewError(
			state.ErrTaskFailedCode,
			string(invokeOutput.Payload),
		)
	}

	return invokeOutput.Payload, nil
}

func (l LambdaTask) Next() string {
	return l.definition.Next()
}

func (l LambdaTask) IsEnd() bool {
	return l.definition.End()
}

type OverrideFn func(input []byte) ([]byte, error)

type OverrideTask struct {
	definition state.TaskDefinition
	fn         OverrideFn
}

func NewOverrideTask(def state.TaskDefinition, fn OverrideFn) State {
	return OverrideTask{
		definition: def,
		fn:         fn,
	}
}

func (o OverrideTask) Run(input []byte) ([]byte, error) {
	return o.fn(input)
}

func (o OverrideTask) Next() string {
	return o.definition.Next()
}

func (o OverrideTask) IsEnd() bool {
	return o.definition.End()
}
