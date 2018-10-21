package sfn

import (
	"encoding/json"
	"time"

	"github.com/eggsbenjamin/stepFnLocal/state"
	"github.com/pkg/errors"
)

// define a ChoiceRule interface
// define a type that implements the interface for each choice rule type
// define a recursive choice rule factory to create the 'choice tree'

// 	TimestampEquals
// 	TimestampLessThan
// 	TimestampGreaterThan
// 	TimestampLessThanEquals
// 	TimestampGreaterThanEquals
// 	And
// 	Or
// 	Not

type ChoiceRule interface {
	Run([]byte) (bool, error)
}

type StringEqualsChoiceRule struct {
	def state.ChoiceRuleDefinition
}

func NewStringEqualsChoiceRule(def state.ChoiceRuleDefinition) StringEqualsChoiceRule {
	return StringEqualsChoiceRule{
		def: def,
	}
}

func (s StringEqualsChoiceRule) Run(input []byte) (bool, error) {
	if s.def.StringEquals == nil {
		return false, errors.Errorf("StringEquals is nil")
	}

	jsonOperand, err := s.def.VariableExp.Search(input)
	if err != nil {
		return false, errors.Wrap(err, "error searching json input")
	}

	var operand string
	if err = json.Unmarshal(jsonOperand, &operand); err != nil {
		return false, errors.Wrap(err, "error unmarshaling operand")
	}

	return operand == *s.def.StringEquals, nil
}

type StringLessThanChoiceRule struct {
	def state.ChoiceRuleDefinition
}

func NewStringLessThanChoiceRule(def state.ChoiceRuleDefinition) StringLessThanChoiceRule {
	return StringLessThanChoiceRule{
		def: def,
	}
}

func (s StringLessThanChoiceRule) Run(input []byte) (bool, error) {
	if s.def.StringLessThan == nil {
		return false, errors.Errorf("StringLessThan is nil")
	}

	jsonOperand, err := s.def.VariableExp.Search(input)
	if err != nil {
		return false, errors.Wrap(err, "error searching json input")
	}

	var operand string
	if err = json.Unmarshal(jsonOperand, &operand); err != nil {
		return false, errors.Wrap(err, "error unmarshaling operand")
	}

	return operand < *s.def.StringLessThan, nil
}

type StringGreaterThanChoiceRule struct {
	def state.ChoiceRuleDefinition
}

func NewStringGreaterThanChoiceRule(def state.ChoiceRuleDefinition) StringGreaterThanChoiceRule {
	return StringGreaterThanChoiceRule{
		def: def,
	}
}

func (s StringGreaterThanChoiceRule) Run(input []byte) (bool, error) {
	if s.def.StringGreaterThan == nil {
		return false, errors.Errorf("StringGreaterThan is nil")
	}

	jsonOperand, err := s.def.VariableExp.Search(input)
	if err != nil {
		return false, errors.Wrap(err, "error searching json input")
	}

	var operand string
	if err = json.Unmarshal(jsonOperand, &operand); err != nil {
		return false, errors.Wrap(err, "error unmarshaling operand")
	}

	return operand > *s.def.StringGreaterThan, nil
}

type StringLessThanEqualsChoiceRule struct {
	def state.ChoiceRuleDefinition
}

func NewStringLessThanEqualsChoiceRule(def state.ChoiceRuleDefinition) StringLessThanEqualsChoiceRule {
	return StringLessThanEqualsChoiceRule{
		def: def,
	}
}

func (s StringLessThanEqualsChoiceRule) Run(input []byte) (bool, error) {
	if s.def.StringLessThanEquals == nil {
		return false, errors.Errorf("StringLessThanEquals is nil")
	}

	jsonOperand, err := s.def.VariableExp.Search(input)
	if err != nil {
		return false, errors.Wrap(err, "error searching json input")
	}

	var operand string
	if err = json.Unmarshal(jsonOperand, &operand); err != nil {
		return false, errors.Wrap(err, "error unmarshaling operand")
	}

	return operand <= *s.def.StringLessThanEquals, nil
}

type StringGreaterThanEqualsChoiceRule struct {
	def state.ChoiceRuleDefinition
}

func NewStringGreaterThanEqualsChoiceRule(def state.ChoiceRuleDefinition) StringGreaterThanEqualsChoiceRule {
	return StringGreaterThanEqualsChoiceRule{
		def: def,
	}
}

func (s StringGreaterThanEqualsChoiceRule) Run(input []byte) (bool, error) {
	if s.def.StringGreaterThanEquals == nil {
		return false, errors.Errorf("StringGreaterThanEquals is nil")
	}

	jsonOperand, err := s.def.VariableExp.Search(input)
	if err != nil {
		return false, errors.Wrap(err, "error searching json input")
	}

	var operand string
	if err = json.Unmarshal(jsonOperand, &operand); err != nil {
		return false, errors.Wrap(err, "error unmarshaling operand")
	}

	return operand >= *s.def.StringGreaterThanEquals, nil
}

type NumericEqualsChoiceRule struct {
	def state.ChoiceRuleDefinition
}

func NewNumericEqualsChoiceRule(def state.ChoiceRuleDefinition) NumericEqualsChoiceRule {
	return NumericEqualsChoiceRule{
		def: def,
	}
}

func (s NumericEqualsChoiceRule) Run(input []byte) (bool, error) {
	if s.def.NumericEquals == nil {
		return false, errors.Errorf("NumericEquals is nil")
	}

	jsonOperand, err := s.def.VariableExp.Search(input)
	if err != nil {
		return false, errors.Wrap(err, "error searching json input")
	}

	var operand float64
	if err = json.Unmarshal(jsonOperand, &operand); err != nil {
		return false, errors.Wrap(err, "error unmarshaling operand")
	}

	return operand == *s.def.NumericEquals, nil
}

type NumericLessThanChoiceRule struct {
	def state.ChoiceRuleDefinition
}

func NewNumericLessThanChoiceRule(def state.ChoiceRuleDefinition) NumericLessThanChoiceRule {
	return NumericLessThanChoiceRule{
		def: def,
	}
}

func (s NumericLessThanChoiceRule) Run(input []byte) (bool, error) {
	if s.def.NumericLessThan == nil {
		return false, errors.Errorf("NumericLessThan is nil")
	}

	jsonOperand, err := s.def.VariableExp.Search(input)
	if err != nil {
		return false, errors.Wrap(err, "error searching json input")
	}

	var operand float64
	if err = json.Unmarshal(jsonOperand, &operand); err != nil {
		return false, errors.Wrap(err, "error unmarshaling operand")
	}

	return operand < *s.def.NumericLessThan, nil
}

type NumericGreaterThanChoiceRule struct {
	def state.ChoiceRuleDefinition
}

func NewNumericGreaterThanChoiceRule(def state.ChoiceRuleDefinition) NumericGreaterThanChoiceRule {
	return NumericGreaterThanChoiceRule{
		def: def,
	}
}

func (s NumericGreaterThanChoiceRule) Run(input []byte) (bool, error) {
	if s.def.NumericGreaterThan == nil {
		return false, errors.Errorf("NumericGreaterThan is nil")
	}

	jsonOperand, err := s.def.VariableExp.Search(input)
	if err != nil {
		return false, errors.Wrap(err, "error searching json input")
	}

	var operand float64
	if err = json.Unmarshal(jsonOperand, &operand); err != nil {
		return false, errors.Wrap(err, "error unmarshaling operand")
	}

	return operand > *s.def.NumericGreaterThan, nil
}

type NumericLessThanEqualsChoiceRule struct {
	def state.ChoiceRuleDefinition
}

func NewNumericLessThanEqualsChoiceRule(def state.ChoiceRuleDefinition) NumericLessThanEqualsChoiceRule {
	return NumericLessThanEqualsChoiceRule{
		def: def,
	}
}

func (s NumericLessThanEqualsChoiceRule) Run(input []byte) (bool, error) {
	if s.def.NumericLessThanEquals == nil {
		return false, errors.Errorf("NumericLessThanEquals is nil")
	}

	jsonOperand, err := s.def.VariableExp.Search(input)
	if err != nil {
		return false, errors.Wrap(err, "error searching json input")
	}

	var operand float64
	if err = json.Unmarshal(jsonOperand, &operand); err != nil {
		return false, errors.Wrap(err, "error unmarshaling operand")
	}

	return operand <= *s.def.NumericLessThanEquals, nil
}

type NumericGreaterThanEqualsChoiceRule struct {
	def state.ChoiceRuleDefinition
}

func NewNumericGreaterThanEqualsChoiceRule(def state.ChoiceRuleDefinition) NumericGreaterThanEqualsChoiceRule {
	return NumericGreaterThanEqualsChoiceRule{
		def: def,
	}
}

func (s NumericGreaterThanEqualsChoiceRule) Run(input []byte) (bool, error) {
	if s.def.NumericGreaterThanEquals == nil {
		return false, errors.Errorf("NumericGreaterThanEquals is nil")
	}

	jsonOperand, err := s.def.VariableExp.Search(input)
	if err != nil {
		return false, errors.Wrap(err, "error searching json input")
	}

	var operand float64
	if err = json.Unmarshal(jsonOperand, &operand); err != nil {
		return false, errors.Wrap(err, "error unmarshaling operand")
	}

	return operand >= *s.def.NumericGreaterThanEquals, nil
}

type BooleanEqualsChoiceRule struct {
	def state.ChoiceRuleDefinition
}

func NewBooleanEqualsChoiceRule(def state.ChoiceRuleDefinition) BooleanEqualsChoiceRule {
	return BooleanEqualsChoiceRule{
		def: def,
	}
}

func (s BooleanEqualsChoiceRule) Run(input []byte) (bool, error) {
	if s.def.BooleanEquals == nil {
		return false, errors.Errorf("BooleanEquals is nil")
	}

	jsonOperand, err := s.def.VariableExp.Search(input)
	if err != nil {
		return false, errors.Wrap(err, "error searching json input")
	}

	var operand bool
	if err = json.Unmarshal(jsonOperand, &operand); err != nil {
		return false, errors.Wrap(err, "error unmarshaling operand")
	}

	return operand == *s.def.BooleanEquals, nil
}

type TimestampEqualsChoiceRule struct {
	def state.ChoiceRuleDefinition
}

func NewTimestampEqualsChoiceRule(def state.ChoiceRuleDefinition) TimestampEqualsChoiceRule {
	return TimestampEqualsChoiceRule{
		def: def,
	}
}

func (s TimestampEqualsChoiceRule) Run(input []byte) (bool, error) {
	if s.def.TimestampEquals == nil {
		return false, errors.Errorf("TimestampEquals is nil")
	}

	jsonOperand, err := s.def.VariableExp.Search(input)
	if err != nil {
		return false, errors.Wrap(err, "error searching json input")
	}

	var operand time.Time
	if err = json.Unmarshal(jsonOperand, &operand); err != nil {
		return false, errors.Wrap(err, "error unmarshaling operand")
	}

	return operand.Equal(*s.def.TimestampEquals), nil
}

type TimestampLessThanChoiceRule struct {
	def state.ChoiceRuleDefinition
}

func NewTimestampLessThanChoiceRule(def state.ChoiceRuleDefinition) TimestampLessThanChoiceRule {
	return TimestampLessThanChoiceRule{
		def: def,
	}
}

func (s TimestampLessThanChoiceRule) Run(input []byte) (bool, error) {
	if s.def.TimestampLessThan == nil {
		return false, errors.Errorf("TimestampLessThan is nil")
	}

	jsonOperand, err := s.def.VariableExp.Search(input)
	if err != nil {
		return false, errors.Wrap(err, "error searching json input")
	}

	var operand time.Time
	if err = json.Unmarshal(jsonOperand, &operand); err != nil {
		return false, errors.Wrap(err, "error unmarshaling operand")
	}

	return operand.Before(*s.def.TimestampLessThan), nil
}

type TimestampGreaterThanChoiceRule struct {
	def state.ChoiceRuleDefinition
}

func NewTimestampGreaterThanChoiceRule(def state.ChoiceRuleDefinition) TimestampGreaterThanChoiceRule {
	return TimestampGreaterThanChoiceRule{
		def: def,
	}
}

func (s TimestampGreaterThanChoiceRule) Run(input []byte) (bool, error) {
	if s.def.TimestampGreaterThan == nil {
		return false, errors.Errorf("TimestampGreaterThan is nil")
	}

	jsonOperand, err := s.def.VariableExp.Search(input)
	if err != nil {
		return false, errors.Wrap(err, "error searching json input")
	}

	var operand time.Time
	if err = json.Unmarshal(jsonOperand, &operand); err != nil {
		return false, errors.Wrap(err, "error unmarshaling operand")
	}

	return operand.After(*s.def.TimestampGreaterThan), nil
}

type TimestampLessThanEqualsChoiceRule struct {
	def state.ChoiceRuleDefinition
}

func NewTimestampLessThanEqualsChoiceRule(def state.ChoiceRuleDefinition) TimestampLessThanEqualsChoiceRule {
	return TimestampLessThanEqualsChoiceRule{
		def: def,
	}
}

func (s TimestampLessThanEqualsChoiceRule) Run(input []byte) (bool, error) {
	if s.def.TimestampLessThanEquals == nil {
		return false, errors.Errorf("TimestampLessThanEquals is nil")
	}

	jsonOperand, err := s.def.VariableExp.Search(input)
	if err != nil {
		return false, errors.Wrap(err, "error searching json input")
	}

	var operand time.Time
	if err = json.Unmarshal(jsonOperand, &operand); err != nil {
		return false, errors.Wrap(err, "error unmarshaling operand")
	}

	return operand.Before(*s.def.TimestampLessThanEquals) || operand.Equal(*s.def.TimestampLessThanEquals), nil
}

type TimestampGreaterThanEqualsChoiceRule struct {
	def state.ChoiceRuleDefinition
}

func NewTimestampGreaterThanEqualsChoiceRule(def state.ChoiceRuleDefinition) TimestampGreaterThanEqualsChoiceRule {
	return TimestampGreaterThanEqualsChoiceRule{
		def: def,
	}
}

func (s TimestampGreaterThanEqualsChoiceRule) Run(input []byte) (bool, error) {
	if s.def.TimestampGreaterThanEquals == nil {
		return false, errors.Errorf("TimestampGreaterThanEquals is nil")
	}

	jsonOperand, err := s.def.VariableExp.Search(input)
	if err != nil {
		return false, errors.Wrap(err, "error searching json input")
	}

	var operand time.Time
	if err = json.Unmarshal(jsonOperand, &operand); err != nil {
		return false, errors.Wrap(err, "error unmarshaling operand")
	}

	return operand.After(*s.def.TimestampGreaterThanEquals) || operand.Equal(*s.def.TimestampGreaterThanEquals), nil
}
