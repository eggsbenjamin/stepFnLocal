// +build acceptance

package sfn_test

import (
	"encoding/json"
	"testing"

	"github.com/eggsbenjamin/stepFnLocal/sfn"
	"github.com/eggsbenjamin/stepFnLocal/state"
	"github.com/stretchr/testify/require"
)

func TestParallelStateAcceptance(t *testing.T) {
	def := `
	{
  "Comment": "test",
  "StartAt": "Parallel_Test",
  "States": {
    "Parallel_Test":{
      "Type": "Parallel",
      "Branches": [
        {
          "StartAt": "TEST_BRANCH_1",
          "States": {
            "TEST_BRANCH_1": {
              "Type": "Task",
              "Resource": "TEST_BRANCH_1_ARN",
              "End": true
            }
          }
        },
        {
          "StartAt": "TEST_BRANCH_2",
          "States": {
            "TEST_BRANCH_2": {
              "Type": "Task",
              "Resource": "TEST_BRANCH_2_ARN",
              "End": true
            }
          }
        },
        {
          "StartAt": "TEST_BRANCH_3",
          "States": {
            "TEST_BRANCH_3": {
              "Type": "Task",
              "Resource": "TEST_BRANCH_3_ARN",
              "End": true
            }
          }
        }
      ],
			"End": true
		}
  }
}
	`

	var machineDef state.MachineDefinition
	require.NoError(t, json.Unmarshal([]byte(def), &machineDef))

	overrides := map[string]sfn.OverrideFn{
		"TEST_BRANCH_1_ARN": func([]byte) ([]byte, error) {
			return []byte(`"TEST_BRANCH_1_ARN_RESULT"`), nil
		},
		"TEST_BRANCH_2_ARN": func([]byte) ([]byte, error) {
			return []byte(`"TEST_BRANCH_2_ARN_RESULT"`), nil
		},
		"TEST_BRANCH_3_ARN": func([]byte) ([]byte, error) {
			return []byte(`"TEST_BRANCH_3_ARN_RESULT"`), nil
		},
	}
	fn, err := sfn.New(machineDef, overrides)
	require.NoError(t, err)

	input := []byte(`"input"`)
	expectedResult := sfn.ExecutionResult{
		Input:  input,
		Output: []byte(`["TEST_BRANCH_1_ARN_RESULT","TEST_BRANCH_2_ARN_RESULT","TEST_BRANCH_3_ARN_RESULT"]`),
		Status: sfn.ExecutionStatusSucceeded,
	}

	result, err := fn.StartExecution(input)
	require.NoError(t, err)
	require.Equal(t, expectedResult.Input, result.Input)
	require.JSONEq(t, string(expectedResult.Output), string(result.Output))
	require.Equal(t, expectedResult.Status, result.Status)
}
