package stepfunction_test

import (
	"encoding/json"
	"testing"

	"github.com/eggsbenjamin/stepFnLocal/state"
	"github.com/eggsbenjamin/stepFnLocal/step_function"
	"github.com/stretchr/testify/require"
)

func TestStepFunction(t *testing.T) {
	t.Run("StartExecution", func(t *testing.T) {
		t.Run("invalid state machine definition", func(t *testing.T) {
			tests := []struct {
				title        string
				stateMachine *state.MachineDefinition
			}{
				{
					"nil definition",
					nil,
				},
				{
					"missing StartAt",
					&state.MachineDefinition{},
				},
			}

			for _, tt := range tests {
				sfn := stepfunction.New(tt.stateMachine, nil)
				require.Error(t, sfn.StartExecution(json.RawMessage(`{}`)))
			}
		})
	})
}
