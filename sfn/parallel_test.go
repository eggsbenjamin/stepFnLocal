package sfn_test

import (
	"testing"

	"github.com/eggsbenjamin/stepFnLocal/sfn"
	"github.com/eggsbenjamin/stepFnLocal/state"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestParallelState(t *testing.T) {
	dummyErr := errors.New("error")
	tests := []struct {
		title       string
		setup       func(*gomock.Controller) (state sfn.State, expectedOutput []byte)
		expectedErr error
	}{
		{
			"single branch error",
			func(ctrl *gomock.Controller) (sfn.State, []byte) {
				branch := sfn.NewMockStepFunction(ctrl)
				branch.EXPECT().StartExecution(gomock.Any()).Return(
					sfn.ExecutionResult{},
					dummyErr,
				)
				parallelState := sfn.NewParallelState(
					state.ParallelDefinition{},
					branch,
				)

				return parallelState, []byte{}
			},
			dummyErr,
		},
		{
			"multiple branches one error",
			func(ctrl *gomock.Controller) (sfn.State, []byte) {
				branch1 := sfn.NewMockStepFunction(ctrl)
				branch1.EXPECT().StartExecution(gomock.Any()).Return(
					sfn.ExecutionResult{
						Status: sfn.ExecutionStatusSucceeded,
					},
					nil,
				)
				branch2 := sfn.NewMockStepFunction(ctrl)
				branch2.EXPECT().StartExecution(gomock.Any()).Return(
					sfn.ExecutionResult{},
					dummyErr,
				)
				parallelState := sfn.NewParallelState(
					state.ParallelDefinition{},
					branch1,
					branch2,
				)

				return parallelState, []byte{}
			},
			dummyErr,
		},
		{
			"single branch success",
			func(ctrl *gomock.Controller) (sfn.State, []byte) {
				branch := sfn.NewMockStepFunction(ctrl)
				branch.EXPECT().StartExecution(gomock.Any()).Return(
					sfn.ExecutionResult{
						Status: sfn.ExecutionStatusSucceeded,
						Output: []byte(`{"result":"test"}`),
					},
					nil,
				)
				parallelState := sfn.NewParallelState(
					state.ParallelDefinition{},
					branch,
				)

				return parallelState, []byte(`[{"result":"test"}]`)
			},
			nil,
		},
		{
			"multiple branch success",
			func(ctrl *gomock.Controller) (sfn.State, []byte) {
				branch1 := sfn.NewMockStepFunction(ctrl)
				branch1.EXPECT().StartExecution(gomock.Any()).Return(
					sfn.ExecutionResult{
						Status: sfn.ExecutionStatusSucceeded,
						Output: []byte(`{"result":"test"}`),
					},
					nil,
				)
				branch2 := sfn.NewMockStepFunction(ctrl)
				branch2.EXPECT().StartExecution(gomock.Any()).Return(
					sfn.ExecutionResult{
						Status: sfn.ExecutionStatusSucceeded,
						Output: []byte(`"test"`),
					},
					nil,
				)
				branch3 := sfn.NewMockStepFunction(ctrl)
				branch3.EXPECT().StartExecution(gomock.Any()).Return(
					sfn.ExecutionResult{
						Status: sfn.ExecutionStatusSucceeded,
						Output: []byte(`1351`),
					},
					nil,
				)
				parallelState := sfn.NewParallelState(
					state.ParallelDefinition{},
					branch1,
					branch2,
					branch3,
				)

				return parallelState, []byte(`[{"result":"test"},"test",1351]`)
			},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			parallelState, expectedOutput := tt.setup(ctrl)

			result, err := parallelState.Run([]byte{})
			if tt.expectedErr != nil {
				require.Equal(t, tt.expectedErr, errors.Cause(err))
				return
			}

			require.NoError(t, err)
			require.JSONEq(t, string(expectedOutput), string(result))
			ctrl.Finish()
		})
	}
}
