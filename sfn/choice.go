//go:generate mockgen -package sfn -source=choice.go -destination choice_mock.go

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

type ChoiceState struct {
	def     state.ChoiceDefinition
	choices []ChoiceRule
	next    string
}

func NewChoiceState(def state.ChoiceDefinition, choices ...ChoiceRule) State {
	return &ChoiceState{
		def:     def,
		choices: choices,
	}
}

func (c *ChoiceState) Run(input []byte) ([]byte, error) {
	for _, choice := range c.choices {
		result, err := choice.Run(input)
		if err != nil {
			return []byte{}, errors.Wrap(err, "error running choice state")
		}
		if result {
			c.next = choice.Next()
			return input, nil
		}
	}

	if c.def.DefaultState == "" {
		return input, state.ErrNoChoiceMatched
	}

	c.next = c.def.DefaultState
	return input, nil
}

func (c *ChoiceState) Next() string {
	return c.next
}

func (c ChoiceState) IsEnd() bool {
	return false
}

type ChoiceRule interface {
	Run([]byte) (bool, error)
	Next() string
}

type ChoiceRuleFactory interface {
	Create(state.ChoiceRuleDefinition) (ChoiceRule, error)
}

type choiceRuleFactory struct{}

func NewChoiceRuleFactory() ChoiceRuleFactory {
	return choiceRuleFactory{}
}

func (c choiceRuleFactory) Create(def state.ChoiceRuleDefinition) (ChoiceRule, error) {
	switch def.Type() {
	case state.StringEquals:
		return NewStringEqualsChoiceRule(def), nil
	case state.StringLessThan:
		return NewStringLessThanChoiceRule(def), nil
	case state.StringGreaterThan:
		return NewStringGreaterThanChoiceRule(def), nil
	case state.StringLessThanEquals:
		return NewStringLessThanEqualsChoiceRule(def), nil
	case state.StringGreaterThanEquals:
		return NewStringGreaterThanEqualsChoiceRule(def), nil
	case state.NumericEquals:
		return NewNumericEqualsChoiceRule(def), nil
	case state.NumericLessThan:
		return NewNumericLessThanChoiceRule(def), nil
	case state.NumericGreaterThan:
		return NewNumericGreaterThanChoiceRule(def), nil
	case state.NumericLessThanEquals:
		return NewNumericLessThanEqualsChoiceRule(def), nil
	case state.NumericGreaterThanEquals:
		return NewNumericGreaterThanEqualsChoiceRule(def), nil
	case state.BooleanEquals:
		return NewBooleanEqualsChoiceRule(def), nil
	case state.TimestampEquals:
		return NewTimestampEqualsChoiceRule(def), nil
	case state.TimestampLessThan:
		return NewTimestampLessThanChoiceRule(def), nil
	case state.TimestampGreaterThan:
		return NewTimestampGreaterThanChoiceRule(def), nil
	case state.TimestampLessThanEquals:
		return NewTimestampLessThanEqualsChoiceRule(def), nil
	case state.TimestampGreaterThanEquals:
		return NewTimestampGreaterThanEqualsChoiceRule(def), nil
	case state.And:
		choiceRules := []ChoiceRule{}
		for _, choiceRuleDef := range def.And {
			choiceRule, err := c.Create(choiceRuleDef)
			if err != nil {
				return nil, err
			}
			choiceRules = append(choiceRules, choiceRule)
		}
		return NewAndChoiceRule(def, choiceRules...), nil
	case state.Or:
		choiceRules := []ChoiceRule{}
		for _, choiceRuleDef := range def.Or {
			choiceRule, err := c.Create(choiceRuleDef)
			if err != nil {
				return nil, err
			}
			choiceRules = append(choiceRules, choiceRule)
		}
		return NewOrChoiceRule(def, choiceRules...), nil
	case state.Not:
		choiceRule, err := c.Create(def)
		if err != nil {
			return nil, err
		}
		return NewNotChoiceRule(def, choiceRule), nil
	}

	return nil, nil
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

func (s StringEqualsChoiceRule) Next() string {
	return s.def.NextState
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

func (s StringLessThanChoiceRule) Next() string {
	return s.def.NextState
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

func (s StringGreaterThanChoiceRule) Next() string {
	return s.def.NextState
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

func (s StringLessThanEqualsChoiceRule) Next() string {
	return s.def.NextState
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

func (s StringGreaterThanEqualsChoiceRule) Next() string {
	return s.def.NextState
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

func (s NumericEqualsChoiceRule) Next() string {
	return s.def.NextState
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

func (s NumericLessThanChoiceRule) Next() string {
	return s.def.NextState
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

func (s NumericGreaterThanChoiceRule) Next() string {
	return s.def.NextState
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

func (s NumericLessThanEqualsChoiceRule) Next() string {
	return s.def.NextState
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

func (s NumericGreaterThanEqualsChoiceRule) Next() string {
	return s.def.NextState
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

func (s BooleanEqualsChoiceRule) Next() string {
	return s.def.NextState
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

func (s TimestampEqualsChoiceRule) Next() string {
	return s.def.NextState
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

func (s TimestampLessThanChoiceRule) Next() string {
	return s.def.NextState
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

func (s TimestampGreaterThanChoiceRule) Next() string {
	return s.def.NextState
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

func (s TimestampLessThanEqualsChoiceRule) Next() string {
	return s.def.NextState
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

func (s TimestampGreaterThanEqualsChoiceRule) Next() string {
	return s.def.NextState
}

type AndChoiceRule struct {
	def         state.ChoiceRuleDefinition
	choiceRules []ChoiceRule
}

func NewAndChoiceRule(def state.ChoiceRuleDefinition, choiceRules ...ChoiceRule) AndChoiceRule {
	return AndChoiceRule{
		def:         def,
		choiceRules: choiceRules,
	}
}

func (s AndChoiceRule) Run(input []byte) (bool, error) {
	for _, choiceRule := range s.choiceRules {
		result, err := choiceRule.Run(input)
		if err != nil {
			return false, err
		}
		if !result {
			return false, nil
		}
	}

	return true, nil
}

func (s AndChoiceRule) Next() string {
	return s.def.NextState
}

type OrChoiceRule struct {
	def         state.ChoiceRuleDefinition
	choiceRules []ChoiceRule
}

func NewOrChoiceRule(def state.ChoiceRuleDefinition, choiceRules ...ChoiceRule) OrChoiceRule {
	return OrChoiceRule{
		def:         def,
		choiceRules: choiceRules,
	}
}

func (s OrChoiceRule) Run(input []byte) (bool, error) {
	for _, choiceRule := range s.choiceRules {
		result, err := choiceRule.Run(input)
		if err != nil {
			return false, err
		}
		if result {
			return true, nil
		}
	}

	return false, nil
}

func (s OrChoiceRule) Next() string {
	return s.def.NextState
}

type NotChoiceRule struct {
	def        state.ChoiceRuleDefinition
	choiceRule ChoiceRule
}

func NewNotChoiceRule(def state.ChoiceRuleDefinition, choiceRule ChoiceRule) NotChoiceRule {
	return NotChoiceRule{
		def:        def,
		choiceRule: choiceRule,
	}
}

func (s NotChoiceRule) Run(input []byte) (bool, error) {
	result, err := s.choiceRule.Run(input)
	if err != nil {
		return false, err
	}

	return !result, nil
}

func (s NotChoiceRule) Next() string {
	return s.def.NextState
}
