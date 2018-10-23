// +build unit

package sfn_test

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/eggsbenjamin/stepFnLocal/sfn"
	"github.com/eggsbenjamin/stepFnLocal/state"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestChoiceRules(t *testing.T) {
	t.Run("StringEqualsChoiceRule", func(t *testing.T) {
		def := state.ChoiceRuleDefinition{
			VariableExp:  state.JSONPathExp("$"),
			StringEquals: aws.String("test"),
		}
		rule := sfn.NewStringEqualsChoiceRule(def)

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
		rule := sfn.NewStringLessThanChoiceRule(def)

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
		rule := sfn.NewStringGreaterThanChoiceRule(def)

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
		rule := sfn.NewStringLessThanEqualsChoiceRule(def)

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
		rule := sfn.NewStringGreaterThanEqualsChoiceRule(def)

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
		rule := sfn.NewNumericEqualsChoiceRule(def)

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
		rule := sfn.NewNumericLessThanChoiceRule(def)

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
		rule := sfn.NewNumericGreaterThanChoiceRule(def)

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
		rule := sfn.NewNumericLessThanEqualsChoiceRule(def)

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
		rule := sfn.NewNumericGreaterThanEqualsChoiceRule(def)

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
		rule := sfn.NewBooleanEqualsChoiceRule(def)

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
		rule := sfn.NewTimestampEqualsChoiceRule(def)

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
		rule := sfn.NewTimestampLessThanChoiceRule(def)

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
		rule := sfn.NewTimestampGreaterThanChoiceRule(def)

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
		rule := sfn.NewTimestampLessThanEqualsChoiceRule(def)

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
		rule := sfn.NewTimestampGreaterThanEqualsChoiceRule(def)

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

	// TODO; reconsider testing approach, are mocks suitable for this? (hint: no, use stubs instead...)
	t.Run("And", func(t *testing.T) {
		tests := []struct {
			title            string
			setupChoiceRules func(*gomock.Controller) []sfn.ChoiceRule
			expectedResult   bool
		}{
			{
				"false",
				func(ctrl *gomock.Controller) []sfn.ChoiceRule {
					choiceRules := []sfn.ChoiceRule{}
					mockChoiceRule := sfn.NewMockChoiceRule(ctrl)
					mockChoiceRule.EXPECT().Run([]byte{}).Return(false, nil)

					return append(choiceRules, mockChoiceRule)
				},
				false,
			},
			{
				"true",
				func(ctrl *gomock.Controller) []sfn.ChoiceRule {
					choiceRules := []sfn.ChoiceRule{}
					mockChoiceRule := sfn.NewMockChoiceRule(ctrl)
					mockChoiceRule.EXPECT().Run([]byte{}).Return(true, nil)

					return append(choiceRules, mockChoiceRule)
				},
				true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.title, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				rule := sfn.NewAndChoiceRule(state.ChoiceRuleDefinition{}, tt.setupChoiceRules(ctrl)...)

				result, err := rule.Run([]byte{})
				require.NoError(t, err)
				require.Equal(t, tt.expectedResult, result)
				ctrl.Finish()
			})
		}
	})

	t.Run("Or", func(t *testing.T) {
		tests := []struct {
			title            string
			setupChoiceRules func(*gomock.Controller) []sfn.ChoiceRule
			expectedResult   bool
		}{
			{
				"false",
				func(ctrl *gomock.Controller) []sfn.ChoiceRule {
					choiceRules := []sfn.ChoiceRule{}
					mockChoiceRule := sfn.NewMockChoiceRule(ctrl)
					mockChoiceRule.EXPECT().Run([]byte{}).Return(false, nil)

					return append(choiceRules, mockChoiceRule)
				},
				false,
			},
			{
				"true",
				func(ctrl *gomock.Controller) []sfn.ChoiceRule {
					choiceRules := []sfn.ChoiceRule{}
					mockChoiceRule := sfn.NewMockChoiceRule(ctrl)
					mockChoiceRule.EXPECT().Run([]byte{}).Return(true, nil)

					return append(choiceRules, mockChoiceRule)
				},
				true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.title, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				rule := sfn.NewOrChoiceRule(state.ChoiceRuleDefinition{}, tt.setupChoiceRules(ctrl)...)

				result, err := rule.Run([]byte{})
				require.NoError(t, err)
				require.Equal(t, tt.expectedResult, result)
				ctrl.Finish()
			})
		}
	})

	t.Run("Not", func(t *testing.T) {
		tests := []struct {
			title           string
			setupChoiceRule func(*gomock.Controller) sfn.ChoiceRule
			expectedResult  bool
		}{
			{
				"false",
				func(ctrl *gomock.Controller) sfn.ChoiceRule {
					mockChoiceRule := sfn.NewMockChoiceRule(ctrl)
					mockChoiceRule.EXPECT().Run([]byte{}).Return(false, nil)

					return mockChoiceRule
				},
				true,
			},
			{
				"true",
				func(ctrl *gomock.Controller) sfn.ChoiceRule {
					mockChoiceRule := sfn.NewMockChoiceRule(ctrl)
					mockChoiceRule.EXPECT().Run([]byte{}).Return(true, nil)

					return mockChoiceRule
				},
				false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.title, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				rule := sfn.NewNotChoiceRule(state.ChoiceRuleDefinition{}, tt.setupChoiceRule(ctrl))

				result, err := rule.Run([]byte{})
				require.NoError(t, err)
				require.Equal(t, tt.expectedResult, result)
				ctrl.Finish()
			})
		}
	})
}

// TODO: test logical operator creation
func TestChoiceRuleFactory(t *testing.T) {
	tests := []struct {
		title          string
		input          state.ChoiceRuleDefinition
		expectedResult sfn.ChoiceRule
	}{}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			result, err := sfn.NewChoiceRuleFactory().Create(tt.input)
			require.NoError(t, err)
			require.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestChoiceState(t *testing.T) {
	input := []byte(`"input"`)
	tests := []struct {
		title        string
		setup        func(*gomock.Controller) sfn.State
		expectedNext string
		expectedErr  error
	}{
		{
			"first rule true",
			func(ctrl *gomock.Controller) sfn.State {
				choiceRules := []sfn.ChoiceRule{}
				mockChoiceRule1 := sfn.NewMockChoiceRule(ctrl)
				mockChoiceRule1.EXPECT().Run(input).Return(true, nil)
				mockChoiceRule1.EXPECT().Next().Return("test")

				mockChoiceRule2 := sfn.NewMockChoiceRule(ctrl)

				choiceRules = append(choiceRules, mockChoiceRule1, mockChoiceRule2)
				return sfn.NewChoiceState(state.ChoiceDefinition{}, choiceRules...)
			},
			"test",
			nil,
		},
		{
			"non-first rule true",
			func(ctrl *gomock.Controller) sfn.State {
				choiceRules := []sfn.ChoiceRule{}
				mockChoiceRule1 := sfn.NewMockChoiceRule(ctrl)
				mockChoiceRule1.EXPECT().Run(input).Return(false, nil)

				mockChoiceRule2 := sfn.NewMockChoiceRule(ctrl)
				mockChoiceRule2.EXPECT().Run(input).Return(true, nil)
				mockChoiceRule2.EXPECT().Next().Return("test")

				choiceRules = append(choiceRules, mockChoiceRule1, mockChoiceRule2)
				return sfn.NewChoiceState(state.ChoiceDefinition{}, choiceRules...)
			},
			"test",
			nil,
		},
		{
			"no match with default",
			func(ctrl *gomock.Controller) sfn.State {
				choiceRules := []sfn.ChoiceRule{}
				mockChoiceRule1 := sfn.NewMockChoiceRule(ctrl)
				mockChoiceRule1.EXPECT().Run(input).Return(false, nil)

				mockChoiceRule2 := sfn.NewMockChoiceRule(ctrl)
				mockChoiceRule2.EXPECT().Run(input).Return(false, nil)

				def := state.ChoiceDefinition{
					DefaultState: "test",
				}
				choiceRules = append(choiceRules, mockChoiceRule1, mockChoiceRule2)
				return sfn.NewChoiceState(def, choiceRules...)
			},
			"test",
			nil,
		},
		{
			"no match without default",
			func(ctrl *gomock.Controller) sfn.State {
				choiceRules := []sfn.ChoiceRule{}
				mockChoiceRule1 := sfn.NewMockChoiceRule(ctrl)
				mockChoiceRule1.EXPECT().Run(input).Return(false, nil)

				mockChoiceRule2 := sfn.NewMockChoiceRule(ctrl)
				mockChoiceRule2.EXPECT().Run(input).Return(false, nil)

				def := state.ChoiceDefinition{}
				choiceRules = append(choiceRules, mockChoiceRule1, mockChoiceRule2)
				return sfn.NewChoiceState(def, choiceRules...)
			},
			"",
			state.NewError(
				state.ErrNoChoiceMatchedCode,
				"",
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			choiceState := tt.setup(ctrl)

			result, err := choiceState.Run(input)
			require.Equal(t, input, result) // should never modify it's input

			if tt.expectedErr != nil {
				require.Equal(t, tt.expectedErr, errors.Cause(err))
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expectedNext, choiceState.Next())
			ctrl.Finish()
		})
	}
}
