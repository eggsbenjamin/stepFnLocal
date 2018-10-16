// +build unit

package sfn_test

import (
	"testing"

	"github.com/eggsbenjamin/stepFnLocal/sfn"
	"github.com/eggsbenjamin/stepFnLocal/state"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestStepFunction(t *testing.T) {
	t.Run("StartExecution", func(t *testing.T) {
		t.Run("tasks", func(t *testing.T) {
			def := state.MachineDefinition{
				StartAt: "test1",
				States: state.MachineStates{
					"test1": []byte(`{"Type":"Task","Next":"test2","Resource":"test"}`),
					"test2": []byte(`{"Type":"Task","Next":"test3","Resource":"test"}`),
					"test3": []byte(`{"Type":"Task","End":true,"Resource":"test"}`),
				},
			}

			input := []byte("input")
			output := []byte("output")

			ctrl := gomock.NewController(t)
			mockState := state.NewMockState(ctrl)
			mockStateFactory := sfn.NewMockStateFactory(ctrl)

			gomock.InOrder(
				mockStateFactory.EXPECT().Create(gomock.Any()).Return(mockState, nil),
				mockState.EXPECT().Run(input).Return(output, nil),
				mockState.EXPECT().IsEnd().Return(false),
				mockState.EXPECT().Next().Return("test2"),
				mockStateFactory.EXPECT().Create(gomock.Any()).Return(mockState, nil),
				mockState.EXPECT().Run(output).Return(output, nil),
				mockState.EXPECT().IsEnd().Return(false),
				mockState.EXPECT().Next().Return("test3"),
				mockStateFactory.EXPECT().Create(gomock.Any()).Return(mockState, nil),
				mockState.EXPECT().Run(output).Return(output, nil),
				mockState.EXPECT().IsEnd().Return(true),
			)

			expectedResult := state.ExecutionResult{
				Input:  input,
				Output: output,
				Status: sfn.ExecutionStatusSucceeded,
			}

			fn, err := sfn.New(def, mockStateFactory)
			require.NoError(t, err)

			result, err := fn.StartExecution(input)
			require.NoError(t, err)
			require.Equal(t, expectedResult.Input, result.Input)
			require.Equal(t, expectedResult.Output, result.Output)
			require.Equal(t, expectedResult.Status, result.Status)
			ctrl.Finish()
		})
	})
}
