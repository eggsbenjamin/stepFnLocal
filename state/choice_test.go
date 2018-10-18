// +build unit

package state_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/eggsbenjamin/stepFnLocal/state"
	"github.com/stretchr/testify/require"
)

func TestChoiceRuleDefinitions(t *testing.T) {
	t.Run("BaseChoiceRuleDefinition", func(t *testing.T) {
		t.Run("Validate", func(t *testing.T) {
			tests := []struct {
				title         string
				def           state.BaseChoiceRuleDefinition
				expectedError *state.ValidationError
			}{
				{
					"Variable with And",
					state.BaseChoiceRuleDefinition{
						VariableExp: "$",
						And:         []state.BaseChoiceRuleDefinition{},
					},
					state.NewValidationError(
						state.InvalidCombinationErrType,
						state.OnlyOneMustExistErrMsg,
						"Variable/And/Or/Not",
					),
				},
				{
					"Variable with Or",
					state.BaseChoiceRuleDefinition{
						VariableExp: "$",
						Or:          []state.BaseChoiceRuleDefinition{},
					},
					state.NewValidationError(
						state.InvalidCombinationErrType,
						state.OnlyOneMustExistErrMsg,
						"Variable/And/Or/Not",
					),
				},
				{
					"Variable with Not",
					state.BaseChoiceRuleDefinition{
						VariableExp: "$",
						Not:         &state.BaseChoiceRuleDefinition{},
					},
					state.NewValidationError(
						state.InvalidCombinationErrType,
						state.OnlyOneMustExistErrMsg,
						"Variable/And/Or/Not",
					),
				},
				{
					"And with Or",
					state.BaseChoiceRuleDefinition{
						And: []state.BaseChoiceRuleDefinition{},
						Or:  []state.BaseChoiceRuleDefinition{},
					},
					state.NewValidationError(
						state.InvalidCombinationErrType,
						state.OnlyOneMustExistErrMsg,
						"Variable/And/Or/Not",
					),
				},
				{
					"And with Not",
					state.BaseChoiceRuleDefinition{
						And: []state.BaseChoiceRuleDefinition{},
						Not: &state.BaseChoiceRuleDefinition{},
					},
					state.NewValidationError(
						state.InvalidCombinationErrType,
						state.OnlyOneMustExistErrMsg,
						"Variable/And/Or/Not",
					),
				},
				{
					"Or with Not",
					state.BaseChoiceRuleDefinition{
						Or:  []state.BaseChoiceRuleDefinition{},
						Not: &state.BaseChoiceRuleDefinition{},
					},
					state.NewValidationError(
						state.InvalidCombinationErrType,
						state.OnlyOneMustExistErrMsg,
						"Variable/And/Or/Not",
					),
				},
				{
					"And with operator",
					state.BaseChoiceRuleDefinition{
						And:          []state.BaseChoiceRuleDefinition{},
						StringEquals: aws.String("test"),
					},
					state.NewValidationError(
						state.InvalidCombinationErrType,
						state.OnlyOneMustExistErrMsg,
						"And || "+state.VariableOperatorList,
					),
				},
				{
					"Or with operator",
					state.BaseChoiceRuleDefinition{
						Or:           []state.BaseChoiceRuleDefinition{},
						StringEquals: aws.String("test"),
					},
					state.NewValidationError(
						state.InvalidCombinationErrType,
						state.OnlyOneMustExistErrMsg,
						"Or || "+state.VariableOperatorList,
					),
				},
				{
					"Not with operator",
					state.BaseChoiceRuleDefinition{
						Not:          &state.BaseChoiceRuleDefinition{},
						StringEquals: aws.String("test"),
					},
					state.NewValidationError(
						state.InvalidCombinationErrType,
						state.OnlyOneMustExistErrMsg,
						"Not || "+state.VariableOperatorList,
					),
				},
				{
					"missing Variable/And/Or/Not",
					state.BaseChoiceRuleDefinition{},
					state.NewValidationError(
						state.MissingRequiredFieldErrType,
						"Variable/And/Or/Not", "",
					),
				},
				{
					"variable with no operator",
					state.BaseChoiceRuleDefinition{
						VariableExp: "$",
					},
					state.NewValidationError(
						state.MissingRequiredFieldErrType,
						state.VariableOperatorList, "",
					),
				},
				{
					"variable with more than one operator",
					state.BaseChoiceRuleDefinition{
						VariableExp:          "$",
						StringEquals:         aws.String("test"),
						StringLessThanEquals: aws.String("test"),
					},
					state.NewValidationError(
						state.InvalidCombinationErrType,
						state.VariableOperatorList,
						state.OnlyOneMustExistErrMsg,
					),
				},
				{
					"invalid Variable json path",
					state.BaseChoiceRuleDefinition{
						VariableExp: "invalid json path",
					},
					state.NewValidationError(
						state.InvalidJSONPathErrType,
						"Variable", "invalid json path",
					),
				},
				{
					"missing top level next",
					state.BaseChoiceRuleDefinition{},
					state.NewValidationError(
						state.MissingRequiredFieldErrType,
						"Next", "",
					),
				},
				{
					"And Next populated > top level",
					state.BaseChoiceRuleDefinition{
						And: []state.BaseChoiceRuleDefinition{
							{
								NextState: "test",
							},
						},
					},
					state.NewValidationError(
						state.InvalidKeyErrType,
						"Next", "",
					),
				},
				{
					"Or Next populated > top level",
					state.BaseChoiceRuleDefinition{
						Or: []state.BaseChoiceRuleDefinition{
							{
								NextState: "test",
							},
						},
					},
					state.NewValidationError(
						state.InvalidKeyErrType,
						"Next", "",
					),
				},
				{
					"Not Next populated > top level",
					state.BaseChoiceRuleDefinition{
						Not: &state.BaseChoiceRuleDefinition{
							NextState: "test",
						},
					},
					state.NewValidationError(
						state.InvalidKeyErrType,
						"Next", "",
					),
				},
				{
					"valid variable choice rule",
					state.BaseChoiceRuleDefinition{
						VariableExp:  "$",
						StringEquals: aws.String("test"),
						NextState:    "test",
					},
					nil,
				},
			}

			for _, tt := range tests {
				t.Run(tt.title, func(t *testing.T) {
					err := tt.def.Validate(0)
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

func TestChoiceDefinition(t *testing.T) {
	t.Run("Validate", func(t *testing.T) {
		tests := []struct {
			title         string
			def           state.ChoiceDefinition
			expectedError *state.ValidationError
		}{
			{
				"empty Choices array",
				state.ChoiceDefinition{},
				state.NewValidationError(
					state.InvalidValueErrType,
					"Choices", "Is empty",
				),
			},
		}

		for _, tt := range tests {
			t.Run(tt.title, func(t *testing.T) {
				err := tt.def.Validate()
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
}
