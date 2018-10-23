package sfn_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	awslambda "github.com/aws/aws-sdk-go/service/lambda"
	"github.com/eggsbenjamin/stepFnLocal/lambda"
	"github.com/eggsbenjamin/stepFnLocal/sfn"
	"github.com/eggsbenjamin/stepFnLocal/state"
	"github.com/stretchr/testify/require"
)

func TestLambdaTask(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		dummyErr := errors.New("error")

		t.Run("client", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockClient := lambda.NewMockClient(ctrl)
			mockClient.EXPECT().Invoke(gomock.Any()).Return(&awslambda.InvokeOutput{}, dummyErr)

			task := sfn.NewLambdaTask(
				state.TaskDefinition{},
				arn.ARN{},
				mockClient,
			)

			_, err := task.Run([]byte{})
			require.Equal(t, dummyErr, errors.Cause(err))
			ctrl.Finish()
		})

		t.Run("invocation", func(t *testing.T) {
			errorPayload := []byte(`{"errorMessage" : "error"}`)
			expectedError := state.NewError(
				state.ErrTaskFailedCode,
				string(errorPayload),
			)

			ctrl := gomock.NewController(t)
			mockClient := lambda.NewMockClient(ctrl)
			mockClient.EXPECT().Invoke(gomock.Any()).Return(&awslambda.InvokeOutput{
				FunctionError: aws.String("Handled"),
				Payload:       errorPayload,
			}, nil)

			task := sfn.NewLambdaTask(
				state.TaskDefinition{},
				arn.ARN{},
				mockClient,
			)

			_, err := task.Run([]byte{})
			require.Equal(t, expectedError, err)
			ctrl.Finish()
		})
	})

	t.Run("success", func(t *testing.T) {
		arn, _ := arn.Parse("arn:aws:lambda:eu-west-1:1234567890:function:test")
		arnStr := arn.String()
		input := []byte(`{"test-input":"success"}`)
		output := []byte(`{"test-output":"success"}`)

		ctrl := gomock.NewController(t)
		mockClient := lambda.NewMockClient(ctrl)
		mockClient.EXPECT().Invoke(&awslambda.InvokeInput{
			FunctionName: &arnStr,
			LogType:      aws.String("Tail"),
			Payload:      input,
		}).Return(&awslambda.InvokeOutput{
			Payload: output,
		}, nil)

		task := sfn.NewLambdaTask(
			state.TaskDefinition{},
			arn,
			mockClient,
		)

		result, err := task.Run(input)
		require.NoError(t, err)
		require.Equal(t, output, result)
		ctrl.Finish()
	})
}
