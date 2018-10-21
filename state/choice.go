package state

import (
	"strings"
	"time"
)

/*
	TODO

	Add ChoiceRule constants
	Add ChoiceRuleDefinition
	Add ChoiceDefinition
*/

const (
	StringEquals               = "StringEquals"
	StringLessThan             = "StringLessThan"
	StringGreaterThan          = "StringGreaterThan"
	StringLessThanEquals       = "StringLessThanEquals"
	StringGreaterThanEquals    = "StringGreaterThanEquals"
	NumericEquals              = "NumericEquals"
	NumericLessThan            = "NumericLessThan"
	NumericGreaterThan         = "NumericGreaterThan"
	NumericLessThanEquals      = "NumericLessThanEquals"
	NumericGreaterThanEquals   = "NumericGreaterThanEquals"
	BooleanEquals              = "BooleanEquals"
	TimestampEquals            = "TimestampEquals"
	TimestampLessThan          = "TimestampLessThan"
	TimestampGreaterThan       = "TimestampGreaterThan"
	TimestampLessThanEquals    = "TimestampLessThanEquals"
	TimestampGreaterThanEquals = "TimestampGreaterThanEquals"
	And                        = "And"
	Or                         = "Or"
	Not                        = "Not"
)

var (
	VariableOperators = []string{
		StringEquals,
		StringLessThan,
		StringGreaterThan,
		StringLessThanEquals,
		StringGreaterThanEquals,
		NumericEquals,
		NumericLessThan,
		NumericGreaterThan,
		NumericLessThanEquals,
		NumericGreaterThanEquals,
		BooleanEquals,
		TimestampEquals,
		TimestampLessThan,
		TimestampGreaterThan,
		TimestampLessThanEquals,
		TimestampGreaterThanEquals,
	}
	LogicalOperators = []string{
		And,
		Or,
		Not,
	}
)

var operatorLookup = map[string]struct{}{
	StringEquals:               {},
	StringLessThan:             {},
	StringGreaterThan:          {},
	StringLessThanEquals:       {},
	StringGreaterThanEquals:    {},
	NumericEquals:              {},
	NumericLessThan:            {},
	NumericGreaterThan:         {},
	NumericLessThanEquals:      {},
	NumericGreaterThanEquals:   {},
	BooleanEquals:              {},
	TimestampEquals:            {},
	TimestampLessThan:          {},
	TimestampGreaterThan:       {},
	TimestampLessThanEquals:    {},
	TimestampGreaterThanEquals: {},
	And: {},
	Or:  {},
	Not: {},
}

type VariableOperator interface {
	Operand() interface{}
	Variable() JSONPathExp
}

type RuleDefinition interface {
	Validate(depth int) error
	Type() string
}

type ChoiceRuleDefinition struct {
	VariableExp                JSONPathExp            `json:"Variable"`
	NextState                  string                 `json:"Next"`
	StringEquals               *string                `json:"StringEquals"`
	StringLessThan             *string                `json:"StringLessThan"`
	StringGreaterThan          *string                `json:"StringGreaterThan"`
	StringLessThanEquals       *string                `json:"StringLessThanEquals"`
	StringGreaterThanEquals    *string                `json:"StringGreaterThanEquals"`
	NumericEquals              *float64               `json:"NumericEquals"`
	NumericLessThan            *float64               `json:"NumericLessThan"`
	NumericGreaterThan         *float64               `json:"NumericGreaterThan"`
	NumericLessThanEquals      *float64               `json:"NumericLessThanEquals"`
	NumericGreaterThanEquals   *float64               `json:"NumericGreaterThanEquals"`
	BooleanEquals              *bool                  `json:"BooleanEquals"`
	TimestampEquals            *time.Time             `json:"TimestampLessThan"`
	TimestampLessThan          *time.Time             `json:"TimestampLessThan"`
	TimestampGreaterThan       *time.Time             `json:"TimestampGreaterThan"`
	TimestampLessThanEquals    *time.Time             `json:"TimestampLessThanEquals"`
	TimestampGreaterThanEquals *time.Time             `json:"TimestampGreaterThanEquals"`
	And                        []ChoiceRuleDefinition `json:"And"`
	Or                         []ChoiceRuleDefinition `json:"Or"`
	Not                        *ChoiceRuleDefinition  `json:"Not"`
}

func (b ChoiceRuleDefinition) Validate(depth int) error {
	validationErrs := ValidationErrors{}

	if err := b.validateLogicalOperatorCombinations(); err != nil {
		validationErrs = append(validationErrs, err.(ValidationErrors)...)
	}

	variableOperatorCount := b.countVariableOperators()

	if b.VariableExp != "" {
		if err := b.VariableExp.Validate(); err != nil {
			validationErrs = append(validationErrs, NewValidationError(
				InvalidJSONPathErrType,
				"Variable", string(b.VariableExp),
			))
		}

		if variableOperatorCount == 0 {
			validationErrs = append(validationErrs, NewValidationError(
				MissingRequiredFieldErrType,
				strings.Join(VariableOperators, "/"), "",
			))
		}

		if variableOperatorCount > 1 {
			validationErrs = append(validationErrs, NewValidationError(
				InvalidCombinationErrType,
				strings.Join(VariableOperators, "/"),
				OnlyOneMustExistErrMsg,
			))
		}
	}

	if b.And != nil {
		if variableOperatorCount > 0 {
			validationErrs = append(validationErrs, NewValidationError(
				InvalidCombinationErrType,
				OnlyOneMustExistErrMsg,
				"And || "+strings.Join(VariableOperators, "/"),
			))
		}

		for _, choiceRule := range b.And {
			if err := choiceRule.Validate(depth + 1); err != nil {
				validationErrs = append(validationErrs, err.(ValidationErrors)...)
			}
		}
	}

	if b.Or != nil {
		if variableOperatorCount > 0 {
			validationErrs = append(validationErrs, NewValidationError(
				InvalidCombinationErrType,
				OnlyOneMustExistErrMsg,
				"Or || "+strings.Join(VariableOperators, "/"),
			))
		}

		for _, choiceRule := range b.Or {
			if err := choiceRule.Validate(depth + 1); err != nil {
				validationErrs = append(validationErrs, err.(ValidationErrors)...)
			}
		}
	}

	if b.Not != nil {
		if variableOperatorCount > 0 {
			validationErrs = append(validationErrs, NewValidationError(
				InvalidCombinationErrType,
				OnlyOneMustExistErrMsg,
				"Not || "+strings.Join(VariableOperators, "/"),
			))
		}

		if err := b.Not.Validate(depth + 1); err != nil {
			validationErrs = append(validationErrs, err.(ValidationErrors)...)
		}
	}

	if depth == 0 {
		if b.NextState == "" {
			validationErrs = append(validationErrs, NewValidationError(
				MissingRequiredFieldErrType,
				"Next", "",
			))
		}
	}

	if depth > 0 {
		if b.NextState != "" {
			validationErrs = append(validationErrs, NewValidationError(
				InvalidKeyErrType,
				"Next", "",
			))
		}
	}

	if len(validationErrs) > 0 {
		return validationErrs
	}
	return nil
}

func (b ChoiceRuleDefinition) Type() string {
	if b.StringEquals != nil {
		return StringEquals
	}
	if b.StringLessThan != nil {
		return StringLessThan
	}
	if b.StringGreaterThan != nil {
		return StringGreaterThan
	}
	if b.StringLessThanEquals != nil {
		return StringLessThanEquals
	}
	if b.StringGreaterThanEquals != nil {
		return StringGreaterThanEquals
	}
	if b.NumericEquals != nil {
		return NumericEquals
	}
	if b.NumericLessThan != nil {
		return NumericLessThan
	}
	if b.NumericGreaterThan != nil {
		return NumericGreaterThan
	}
	if b.NumericLessThanEquals != nil {
		return NumericLessThanEquals
	}
	if b.NumericGreaterThanEquals != nil {
		return NumericGreaterThanEquals
	}
	if b.BooleanEquals != nil {
		return BooleanEquals
	}
	if b.TimestampLessThan != nil {
		return TimestampLessThan
	}
	if b.TimestampGreaterThan != nil {
		return TimestampGreaterThan
	}
	if b.TimestampLessThanEquals != nil {
		return TimestampLessThanEquals
	}
	if b.TimestampGreaterThanEquals != nil {
		return TimestampGreaterThanEquals
	}
	if b.And != nil {
		return And
	}
	if b.Or != nil {
		return Or
	}
	if b.Not != nil {
		return Not
	}

	return ""
}

func (b ChoiceRuleDefinition) validateLogicalOperatorCombinations() error {
	validationErrs := ValidationErrors{}

	var count int
	if b.VariableExp != "" {
		count++
	}
	if b.And != nil {
		count++
	}
	if b.Or != nil {
		count++
	}
	if b.Not != nil {
		count++
	}

	if count > 1 {
		validationErrs = append(validationErrs, NewValidationError(
			InvalidCombinationErrType,
			OnlyOneMustExistErrMsg,
			"Variable/And/Or/Not",
		))
	}

	if count == 0 {
		validationErrs = append(validationErrs, NewValidationError(
			MissingRequiredFieldErrType,
			"Variable/And/Or/Not", "",
		))
	}

	if len(validationErrs) > 0 {
		return validationErrs
	}
	return nil
}

func (b ChoiceRuleDefinition) countVariableOperators() int {
	var count int

	if b.StringEquals != nil {
		count++
	}
	if b.StringLessThan != nil {
		count++
	}
	if b.StringGreaterThan != nil {
		count++
	}
	if b.StringLessThanEquals != nil {
		count++
	}
	if b.StringGreaterThanEquals != nil {
		count++
	}
	if b.NumericEquals != nil {
		count++
	}
	if b.NumericLessThan != nil {
		count++
	}
	if b.NumericGreaterThan != nil {
		count++
	}
	if b.NumericLessThanEquals != nil {
		count++
	}
	if b.NumericGreaterThanEquals != nil {
		count++
	}
	if b.BooleanEquals != nil {
		count++
	}
	if b.TimestampEquals != nil {
		count++
	}
	if b.TimestampLessThan != nil {
		count++
	}
	if b.TimestampGreaterThan != nil {
		count++
	}
	if b.TimestampLessThanEquals != nil {
		count++
	}
	if b.TimestampGreaterThanEquals != nil {
		count++
	}

	return count
}

type ChoiceDefinition struct {
	Choices      []ChoiceRuleDefinition `json:"Choices"`
	DefaultState string                 `json:"Default"`
	NextState    string                 `json:"-"`
}

func (c ChoiceDefinition) Type() string {
	return ChoiceStateType
}

func (c ChoiceDefinition) Validate() error {
	validationErrs := ValidationErrors{}

	if len(c.Choices) == 0 {
		validationErrs = append(validationErrs, NewValidationError(
			InvalidValueErrType,
			"Choices", "Is empty",
		))
	}

	for _, choiceRule := range c.Choices {
		if err := choiceRule.Validate(0); err != nil {
			validationErrs = append(validationErrs, err.(ValidationErrors)...)
		}
	}

	if len(validationErrs) > 0 {
		return validationErrs
	}
	return nil
}

func (c ChoiceDefinition) Next() string {
	return c.NextState
}
