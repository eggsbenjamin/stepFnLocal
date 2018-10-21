// +build unit

package sfn

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/eggsbenjamin/stepFnLocal/state"
	"github.com/stretchr/testify/require"
)

func TestChoiceRules(t *testing.T) {
	t.Run("StringEqualsChoiceRule", func(t *testing.T) {
		def := state.ChoiceRuleDefinition{
			VariableExp:  state.JSONPathExp("$"),
			StringEquals: aws.String("test"),
		}
		rule := NewStringEqualsChoiceRule(def)

		tests := []struct {
			title          string
			input          []byte
			expectedResult bool
		}{
			{
				"not equal",
				[]byte(`"not equal"`),
				false,
			},
			{
				"equal",
				[]byte(`"test"`),
				true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.title, func(t *testing.T) {
				result, err := rule.Run(tt.input)
				require.NoError(t, err)
				require.Equal(t, tt.expectedResult, result)
			})
		}
	})

	t.Run("StringLessThanChoiceRule", func(t *testing.T) {
		def := state.ChoiceRuleDefinition{
			VariableExp:    state.JSONPathExp("$"),
			StringLessThan: aws.String("test"),
		}
		rule := NewStringLessThanChoiceRule(def)

		tests := []struct {
			title          string
			input          []byte
			expectedResult bool
		}{
			{
				"not less than",
				[]byte(`"test1"`),
				false,
			},
			{
				"less than",
				[]byte(`"te"`),
				true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.title, func(t *testing.T) {
				result, err := rule.Run(tt.input)
				require.NoError(t, err)
				require.Equal(t, tt.expectedResult, result)
			})
		}
	})

	t.Run("StringGreaterThanChoiceRule", func(t *testing.T) {
		def := state.ChoiceRuleDefinition{
			VariableExp:       state.JSONPathExp("$"),
			StringGreaterThan: aws.String("test"),
		}
		rule := NewStringGreaterThanChoiceRule(def)

		tests := []struct {
			title          string
			input          []byte
			expectedResult bool
		}{
			{
				"not greater than",
				[]byte(`"te"`),
				false,
			},
			{
				"greater than",
				[]byte(`"test1"`),
				true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.title, func(t *testing.T) {
				result, err := rule.Run(tt.input)
				require.NoError(t, err)
				require.Equal(t, tt.expectedResult, result)
			})
		}
	})

	t.Run("StringLessThanEqualsChoiceRule", func(t *testing.T) {
		def := state.ChoiceRuleDefinition{
			VariableExp:          state.JSONPathExp("$"),
			StringLessThanEquals: aws.String("test"),
		}
		rule := NewStringLessThanEqualsChoiceRule(def)

		tests := []struct {
			title          string
			input          []byte
			expectedResult bool
		}{
			{
				"not less than",
				[]byte(`"test1"`),
				false,
			},
			{
				"less than",
				[]byte(`"te"`),
				true,
			},
			{
				"equal",
				[]byte(`"test"`),
				true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.title, func(t *testing.T) {
				result, err := rule.Run(tt.input)
				require.NoError(t, err)
				require.Equal(t, tt.expectedResult, result)
			})
		}
	})

	t.Run("StringGreaterThanEqualsChoiceRule", func(t *testing.T) {
		def := state.ChoiceRuleDefinition{
			VariableExp:             state.JSONPathExp("$"),
			StringGreaterThanEquals: aws.String("test"),
		}
		rule := NewStringGreaterThanEqualsChoiceRule(def)

		tests := []struct {
			title          string
			input          []byte
			expectedResult bool
		}{
			{
				"not greater than",
				[]byte(`"te"`),
				false,
			},
			{
				"greater than",
				[]byte(`"test1"`),
				true,
			},
			{
				"equal",
				[]byte(`"test"`),
				true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.title, func(t *testing.T) {
				result, err := rule.Run(tt.input)
				require.NoError(t, err)
				require.Equal(t, tt.expectedResult, result)
			})
		}
	})

	t.Run("NumericEqualsChoiceRule", func(t *testing.T) {
		def := state.ChoiceRuleDefinition{
			VariableExp:   state.JSONPathExp("$"),
			NumericEquals: aws.Float64(1986),
		}
		rule := NewNumericEqualsChoiceRule(def)

		tests := []struct {
			title          string
			input          []byte
			expectedResult bool
		}{
			{
				"not equal",
				[]byte(`7`),
				false,
			},
			{
				"equal",
				[]byte(`1986`),
				true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.title, func(t *testing.T) {
				result, err := rule.Run(tt.input)
				require.NoError(t, err)
				require.Equal(t, tt.expectedResult, result)
			})
		}
	})

	t.Run("NumericLessThanChoiceRule", func(t *testing.T) {
		def := state.ChoiceRuleDefinition{
			VariableExp:     state.JSONPathExp("$"),
			NumericLessThan: aws.Float64(1986),
		}
		rule := NewNumericLessThanChoiceRule(def)

		tests := []struct {
			title          string
			input          []byte
			expectedResult bool
		}{
			{
				"not less than",
				[]byte(`1986.1`),
				false,
			},
			{
				"less than",
				[]byte(`1985.9`),
				true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.title, func(t *testing.T) {
				result, err := rule.Run(tt.input)
				require.NoError(t, err)
				require.Equal(t, tt.expectedResult, result)
			})
		}
	})

	t.Run("NumericGreaterThanChoiceRule", func(t *testing.T) {
		def := state.ChoiceRuleDefinition{
			VariableExp:        state.JSONPathExp("$"),
			NumericGreaterThan: aws.Float64(1986),
		}
		rule := NewNumericGreaterThanChoiceRule(def)

		tests := []struct {
			title          string
			input          []byte
			expectedResult bool
		}{
			{
				"not greater than",
				[]byte(`1985.9`),
				false,
			},
			{
				"greater than",
				[]byte(`1986.1`),
				true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.title, func(t *testing.T) {
				result, err := rule.Run(tt.input)
				require.NoError(t, err)
				require.Equal(t, tt.expectedResult, result)
			})
		}
	})

	t.Run("NumericLessThanEqualsChoiceRule", func(t *testing.T) {
		def := state.ChoiceRuleDefinition{
			VariableExp:           state.JSONPathExp("$"),
			NumericLessThanEquals: aws.Float64(1986),
		}
		rule := NewNumericLessThanEqualsChoiceRule(def)

		tests := []struct {
			title          string
			input          []byte
			expectedResult bool
		}{
			{
				"not less than",
				[]byte(`1986.1`),
				false,
			},
			{
				"less than",
				[]byte(`1985.9`),
				true,
			},
			{
				"equal",
				[]byte(`1986`),
				true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.title, func(t *testing.T) {
				result, err := rule.Run(tt.input)
				require.NoError(t, err)
				require.Equal(t, tt.expectedResult, result)
			})
		}
	})

	t.Run("NumericGreaterThanEqualsChoiceRule", func(t *testing.T) {
		def := state.ChoiceRuleDefinition{
			VariableExp:              state.JSONPathExp("$"),
			NumericGreaterThanEquals: aws.Float64(1986),
		}
		rule := NewNumericGreaterThanEqualsChoiceRule(def)

		tests := []struct {
			title          string
			input          []byte
			expectedResult bool
		}{
			{
				"not greater than",
				[]byte(`1985.9`),
				false,
			},
			{
				"greater than",
				[]byte(`1986.1`),
				true,
			},
			{
				"equal",
				[]byte(`1986`),
				true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.title, func(t *testing.T) {
				result, err := rule.Run(tt.input)
				require.NoError(t, err)
				require.Equal(t, tt.expectedResult, result)
			})
		}
	})

	t.Run("BooleanEqualsChoiceRule", func(t *testing.T) {
		def := state.ChoiceRuleDefinition{
			VariableExp:   state.JSONPathExp("$"),
			BooleanEquals: aws.Bool(true),
		}
		rule := NewBooleanEqualsChoiceRule(def)

		tests := []struct {
			title          string
			input          []byte
			expectedResult bool
		}{
			{
				"not equal",
				[]byte(`false`),
				false,
			},
			{
				"equal",
				[]byte(`true`),
				true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.title, func(t *testing.T) {
				result, err := rule.Run(tt.input)
				require.NoError(t, err)
				require.Equal(t, tt.expectedResult, result)
			})
		}
	})

	t.Run("TimestampEqualsChoiceRule", func(t *testing.T) {
		dummyTime, _ := time.Parse(time.RFC3339, "2018-10-21T19:53:03Z")
		def := state.ChoiceRuleDefinition{
			VariableExp:     state.JSONPathExp("$"),
			TimestampEquals: &dummyTime,
		}
		rule := NewTimestampEqualsChoiceRule(def)

		tests := []struct {
			title          string
			input          []byte
			expectedResult bool
		}{
			{
				"not equal",
				[]byte(`"3030-10-21T19:53:03Z"`),
				false,
			},
			{
				"equal",
				[]byte(`"2018-10-21T19:53:03Z"`),
				true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.title, func(t *testing.T) {
				result, err := rule.Run(tt.input)
				require.NoError(t, err)
				require.Equal(t, tt.expectedResult, result)
			})
		}
	})

	t.Run("TimestampLessThanChoiceRule", func(t *testing.T) {
		dummyTime, _ := time.Parse(time.RFC3339, "2018-10-21T19:53:03Z")
		def := state.ChoiceRuleDefinition{
			VariableExp:       state.JSONPathExp("$"),
			TimestampLessThan: &dummyTime,
		}
		rule := NewTimestampLessThanChoiceRule(def)

		tests := []struct {
			title          string
			input          []byte
			expectedResult bool
		}{
			{
				"not less than",
				[]byte(`"2018-10-21T19:53:04Z"`),
				false,
			},
			{
				"less than",
				[]byte(`"2018-10-21T19:53:02Z"`),
				true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.title, func(t *testing.T) {
				result, err := rule.Run(tt.input)
				require.NoError(t, err)
				require.Equal(t, tt.expectedResult, result)
			})
		}
	})

	t.Run("TimestampGreaterThanChoiceRule", func(t *testing.T) {
		dummyTime, _ := time.Parse(time.RFC3339, "2018-10-21T19:53:03Z")
		def := state.ChoiceRuleDefinition{
			VariableExp:          state.JSONPathExp("$"),
			TimestampGreaterThan: &dummyTime,
		}
		rule := NewTimestampGreaterThanChoiceRule(def)

		tests := []struct {
			title          string
			input          []byte
			expectedResult bool
		}{
			{
				"not greater than",
				[]byte(`"2018-10-21T19:53:02Z"`),
				false,
			},
			{
				"greater than",
				[]byte(`"2018-10-21T19:53:04Z"`),
				true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.title, func(t *testing.T) {
				result, err := rule.Run(tt.input)
				require.NoError(t, err)
				require.Equal(t, tt.expectedResult, result)
			})
		}
	})

	t.Run("TimestampLessThanEqualsChoiceRule", func(t *testing.T) {
		dummyTime, _ := time.Parse(time.RFC3339, "2018-10-21T19:53:03Z")
		def := state.ChoiceRuleDefinition{
			VariableExp:             state.JSONPathExp("$"),
			TimestampLessThanEquals: &dummyTime,
		}
		rule := NewTimestampLessThanEqualsChoiceRule(def)

		tests := []struct {
			title          string
			input          []byte
			expectedResult bool
		}{
			{
				"not less than",
				[]byte(`"2018-10-21T19:53:04Z"`),
				false,
			},
			{
				"less than",
				[]byte(`"2018-10-21T19:53:02Z"`),
				true,
			},
			{
				"equal",
				[]byte(`"2018-10-21T19:53:03Z"`),
				true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.title, func(t *testing.T) {
				result, err := rule.Run(tt.input)
				require.NoError(t, err)
				require.Equal(t, tt.expectedResult, result)
			})
		}
	})

	t.Run("TimestampGreaterThanEqualsChoiceRule", func(t *testing.T) {
		dummyTime, _ := time.Parse(time.RFC3339, "2018-10-21T19:53:03Z")
		def := state.ChoiceRuleDefinition{
			VariableExp:                state.JSONPathExp("$"),
			TimestampGreaterThanEquals: &dummyTime,
		}
		rule := NewTimestampGreaterThanEqualsChoiceRule(def)

		tests := []struct {
			title          string
			input          []byte
			expectedResult bool
		}{
			{
				"not greater than",
				[]byte(`"2018-10-21T19:53:02Z"`),
				false,
			},
			{
				"greater than",
				[]byte(`"2018-10-21T19:53:04Z"`),
				true,
			},
			{
				"equal",
				[]byte(`"2018-10-21T19:53:03Z"`),
				true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.title, func(t *testing.T) {
				result, err := rule.Run(tt.input)
				require.NoError(t, err)
				require.Equal(t, tt.expectedResult, result)
			})
		}
	})
}
