// +build unit

package state_test

import (
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/eggsbenjamin/stepFnLocal/state"
	"github.com/stretchr/testify/require"
)

func TestChoiceRuleDefinitions(t *testing.T) {
	t.Run("ChoiceRuleDefinition", func(t *testing.T) {
		t.Run("Validate", func(t *testing.T) {
			tests := []struct {
				title         string
				def           state.ChoiceRuleDefinition
				expectedError *state.ValidationError
			}{
				{
					"Variable with And",
					state.ChoiceRuleDefinition{
						VariableExp: "$",
						And:         []state.ChoiceRuleDefinition{},
					},
					state.NewValidationError(
						state.InvalidCombinationErrType,
						state.OnlyOneMustExistErrMsg,
						"Variable/And/Or/Not",
					),
				},
				{
					"Variable with Or",
					state.ChoiceRuleDefinition{
						VariableExp: "$",
						Or:          []state.ChoiceRuleDefinition{},
					},
					state.NewValidationError(
						state.InvalidCombinationErrType,
						state.OnlyOneMustExistErrMsg,
						"Variable/And/Or/Not",
					),
				},
				{
					"Variable with Not",
					state.ChoiceRuleDefinition{
						VariableExp: "$",
						Not:         &state.ChoiceRuleDefinition{},
					},
					state.NewValidationError(
						state.InvalidCombinationErrType,
						state.OnlyOneMustExistErrMsg,
						"Variable/And/Or/Not",
					),
				},
				{
					"And with Or",
					state.ChoiceRuleDefinition{
						And: []state.ChoiceRuleDefinition{},
						Or:  []state.ChoiceRuleDefinition{},
					},
					state.NewValidationError(
						state.InvalidCombinationErrType,
						state.OnlyOneMustExistErrMsg,
						"Variable/And/Or/Not",
					),
				},
				{
					"And with Not",
					state.ChoiceRuleDefinition{
						And: []state.ChoiceRuleDefinition{},
						Not: &state.ChoiceRuleDefinition{},
					},
					state.NewValidationError(
						state.InvalidCombinationErrType,
						state.OnlyOneMustExistErrMsg,
						"Variable/And/Or/Not",
					),
				},
				{
					"Or with Not",
					state.ChoiceRuleDefinition{
						Or:  []state.ChoiceRuleDefinition{},
						Not: &state.ChoiceRuleDefinition{},
					},
					state.NewValidationError(
						state.InvalidCombinationErrType,
						state.OnlyOneMustExistErrMsg,
						"Variable/And/Or/Not",
					),
				},
				{
					"And with operator",
					state.ChoiceRuleDefinition{
						And:          []state.ChoiceRuleDefinition{},
						StringEquals: aws.String("test"),
					},
					state.NewValidationError(
						state.InvalidCombinationErrType,
						state.OnlyOneMustExistErrMsg,
						"And || "+strings.Join(state.VariableOperators, "/"),
					),
				},
				{
					"Or with operator",
					state.ChoiceRuleDefinition{
						Or:           []state.ChoiceRuleDefinition{},
						StringEquals: aws.String("test"),
					},
					state.NewValidationError(
						state.InvalidCombinationErrType,
						state.OnlyOneMustExistErrMsg,
						"Or || "+strings.Join(state.VariableOperators, "/"),
					),
				},
				{
					"Not with operator",
					state.ChoiceRuleDefinition{
						Not:          &state.ChoiceRuleDefinition{},
						StringEquals: aws.String("test"),
					},
					state.NewValidationError(
						state.InvalidCombinationErrType,
						state.OnlyOneMustExistErrMsg,
						"Not || "+strings.Join(state.VariableOperators, "/"),
					),
				},
				{
					"missing Variable/And/Or/Not",
					state.ChoiceRuleDefinition{},
					state.NewValidationError(
						state.MissingRequiredFieldErrType,
						"Variable/And/Or/Not", "",
					),
				},
				{
					"variable with no operator",
					state.ChoiceRuleDefinition{
						VariableExp: "$",
					},
					state.NewValidationError(
						state.MissingRequiredFieldErrType,
						strings.Join(state.VariableOperators, "/"), "",
					),
				},
				{
					"variable with more than one operator",
					state.ChoiceRuleDefinition{
						VariableExp:          "$",
						StringEquals:         aws.String("test"),
						StringLessThanEquals: aws.String("test"),
					},
					state.NewValidationError(
						state.InvalidCombinationErrType,
						strings.Join(state.VariableOperators, "/"),
						state.OnlyOneMustExistErrMsg,
					),
				},
				{
					"invalid Variable json path",
					state.ChoiceRuleDefinition{
						VariableExp: "invalid json path",
					},
					state.NewValidationError(
						state.InvalidJSONPathErrType,
						"Variable", "invalid json path",
					),
				},
				{
					"missing top level next",
					state.ChoiceRuleDefinition{},
					state.NewValidationError(
						state.MissingRequiredFieldErrType,
						"Next", "",
					),
				},
				{
					"And Next populated > top level",
					state.ChoiceRuleDefinition{
						And: []state.ChoiceRuleDefinition{
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
					state.ChoiceRuleDefinition{
						Or: []state.ChoiceRuleDefinition{
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
					state.ChoiceRuleDefinition{
						Not: &state.ChoiceRuleDefinition{
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
					state.ChoiceRuleDefinition{
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
