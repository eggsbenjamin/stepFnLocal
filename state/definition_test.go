package state_test

import (
	"encoding/json"
	"testing"

	"github.com/eggsbenjamin/stepFnLocal/state"
	"github.com/stretchr/testify/require"
)

func TestState(t *testing.T) {
	t.Run("MachineStates", func(t *testing.T) {
		t.Run("Get", func(t *testing.T) {
			t.Run("not found", func(t *testing.T) {
				input := `{"Test":{"Type":"Task"}}`
				var machineStates state.MachineStates
				require.NoError(t, json.Unmarshal([]byte(input), &machineStates))

				_, err := machineStates.GetDefinition("Non-existent")
				require.Equal(t, state.ErrStateNotFound, err)
			})

			t.Run("unknown state", func(t *testing.T) {
				input := `{"Test":{"Type":"Unknown"}}`
				var machineStates state.MachineStates
				require.NoError(t, json.Unmarshal([]byte(input), &machineStates))

				_, err := machineStates.GetDefinition("Test")
				require.Equal(t, state.ErrUnknownState, err)
			})

			t.Run("existing", func(t *testing.T) {
				tests := []struct {
					stateType     string
					isCorrectType func(state.Definition) bool
				}{
					{
						state.TaskStateType,
						func(def state.Definition) bool {
							_, ok := def.(state.TaskStateDefinition)
							return ok
						},
					},
					// TODO
					// -  PassStateType
					// - ChoiceStateType
					// - WaitStateType
					// - SucceedStateType
					// - FailStateType
					// - ParallelStateType
				}

				for _, tt := range tests {
					t.Run(tt.stateType, func(t *testing.T) {
						input := `{"Test":{"Type":"` + tt.stateType + `"}}`
						var machineStates state.MachineStates
						require.NoError(t, json.Unmarshal([]byte(input), &machineStates))

						def, err := machineStates.GetDefinition("Test")
						require.NoError(t, err)
						require.True(t, tt.isCorrectType(def))
					})
				}
			})
		})
	})

	t.Run("MachineDefinition", func(t *testing.T) {
		t.Run("Validate", func(t *testing.T) {
			tests := []struct {
				title         string
				stateMachine  state.MachineDefinition
				expectedError *state.ValidationError
			}{
				{
					"missing StartAt",
					state.MachineDefinition{},
					state.NewValidationError(
						state.MissingRequiredFieldErrType,
						"StartAt", "",
					),
				},
				{
					"invalid StartAt",
					state.MachineDefinition{
						StartAt: "unknown",
					},
					state.NewValidationError(
						state.InvalidValueErrType,
						"StartAt", "unknown",
					),
				},
				{
					"missing States",
					state.MachineDefinition{
						StartAt: "test",
					},
					state.NewValidationError(
						state.MissingRequiredFieldErrType,
						"States", "",
					),
				},
				{
					"valid",
					state.MachineDefinition{
						StartAt: "test1",
						States: map[string]json.RawMessage{
							"test1": []byte(`{"Type":"Task", "Resource":"Test", "Next":"test2"}`),
							"test2": []byte(`{"Type":"Task", "Resource":"Test", "End": true}`),
						},
					},
					nil,
				},
			}

			for _, tt := range tests {
				t.Run(tt.title, func(t *testing.T) {
					err := tt.stateMachine.Validate()
					if tt.expectedError == nil {
						require.NoError(t, err)
						return
					}

					require.Error(t, err)
					vErr, ok := err.(state.ValidationErrors)
					require.True(t, ok)
					require.Contains(t, vErr, tt.expectedError)
				})
			}
		})
	})

	t.Run("BaseStateDefinition", func(t *testing.T) {
		t.Run("Validate", func(t *testing.T) {
			tests := []struct {
				title         string
				_state        state.BaseStateDefinition
				expectedError *state.ValidationError
			}{
				{
					"missing type",
					state.BaseStateDefinition{},
					state.NewValidationError(
						state.MissingRequiredFieldErrType,
						"Type", "",
					),
				},
				{
					"invalid type",
					state.BaseStateDefinition{
						StateType: "INVALID",
					},
					state.NewValidationError(
						state.InvalidValueErrType,
						"Type", "INVALID",
					),
				},
				{
					"valid",
					state.BaseStateDefinition{
						StateType: state.TaskStateType,
					},
					nil,
				},
			}

			for _, tt := range tests {
				err := tt._state.Validate()
				if tt.expectedError == nil {
					require.NoError(t, err)
					continue
				}

				require.Error(t, err)
				vErr, ok := err.(state.ValidationErrors)
				require.True(t, ok)
				require.Contains(t, vErr, tt.expectedError)
			}
		})
	})

	t.Run("TransitionStateDefinition", func(t *testing.T) {
		t.Run("Validate", func(t *testing.T) {
			tests := []struct {
				title         string
				_state        state.TransitionStateDefinition
				expectedError *state.ValidationError
			}{
				{
					"missing Next and End:true",
					state.TransitionStateDefinition{},
					state.NewValidationError(
						state.MissingRequiredFieldErrType,
						"Next/End:true", "",
					),
				},
			}

			for _, tt := range tests {
				err := tt._state.Validate()
				if tt.expectedError == nil {
					require.NoError(t, err)
					continue
				}

				require.Error(t, err)
				vErr, ok := err.(state.ValidationErrors)
				require.True(t, ok)
				require.Contains(t, vErr, tt.expectedError)
			}
		})
	})

	t.Run("TaskStateDefinition", func(t *testing.T) {
		t.Run("Validate", func(t *testing.T) {
			tests := []struct {
				title         string
				task          state.TaskStateDefinition
				expectedError *state.ValidationError
			}{
				{
					"missing resource",
					state.TaskStateDefinition{},
					state.NewValidationError(
						state.MissingRequiredFieldErrType,
						"Resource", "",
					),
				},
				{
					"valid",
					state.TaskStateDefinition{
						BaseStateDefinition: state.BaseStateDefinition{
							StateType: state.TaskStateType,
						},
						TransitionStateDefinition: state.TransitionStateDefinition{
							EndState: true,
						},
						Resource: "test",
					},
					nil,
				},
			}

			for _, tt := range tests {
				t.Run(tt.title, func(t *testing.T) {
					err := tt.task.Validate()
					if tt.expectedError == nil {
						require.NoError(t, err)
						return
					}

					require.Error(t, err)
					vErr, ok := err.(state.ValidationErrors)
					require.True(t, ok)
					require.Contains(t, vErr, tt.expectedError)
				})
			}
		})
	})
}
