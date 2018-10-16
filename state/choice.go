package state

import (
	"encoding/json"
	"strconv"
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

var validComparisonOperators = map[string]struct{}{
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

type ChoiceComparisonFn func(v string, ctrl string) (bool, error)

func StringEqualsFn(v string, ctrl string) (bool, error) {
	return v == ctrl, nil
}

func StringLessThanFn(v string, ctrl string) (bool, error) {
	return v < ctrl, nil
}

func StringGreaterThanFn(v string, ctrl string) (bool, error) {
	return v > ctrl, nil
}

func StringLessThanEqualsFn(v string, ctrl string) (bool, error) {
	return v <= ctrl, nil
}

func StringGreaterThanEqualsFn(v string, ctrl string) (bool, error) {
	return v >= ctrl, nil
}

func NumericEqualsFn(v string, ctrl string) (bool, error) {
	n1, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return false, err
	}

	n2, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return false, err
	}

	return n1 == n2, nil
}

func NumericLessThanFn(v string, ctrl string) (bool, error) {
	n1, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return false, err
	}

	n2, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return false, err
	}

	return n1 < n2, nil
}

func NumericGreaterThanFn(v string, ctrl string) (bool, error) {
	n1, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return false, err
	}

	n2, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return false, err
	}

	return n1 > n2, nil
}

func NumericLessThanEqualsFn(v string, ctrl string) (bool, error) {
	n1, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return false, err
	}

	n2, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return false, err
	}

	return n1 <= n2, nil
}

func NumericGreaterThanEqualsFn(v string, ctrl string) (bool, error) {
	n1, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return false, err
	}

	n2, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return false, err
	}

	return n1 >= n2, nil
}

func BooleanEqualsFn(v string, ctrl string) (bool, error) {
	b1, err := strconv.ParseBool(v)
	if err != nil {
		return false, err
	}

	b2, err := strconv.ParseBool(ctrl)
	if err != nil {
		return false, err
	}

	return b1 == b2, nil
}

func TimestampEqualsFn(v string, ctrl string) (bool, error) {
	t1, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return false, err
	}

	t2, err := time.Parse(time.RFC3339, ctrl)
	if err != nil {
		return false, err
	}

	return t1.Equal(t2), nil
}

func TimestampLessThanFn(v string, ctrl string) (bool, error) {
	t1, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return false, err
	}

	t2, err := time.Parse(time.RFC3339, ctrl)
	if err != nil {
		return false, err
	}

	return t1.Before(t2), nil
}

func TimestampGreaterThanFn(v string, ctrl string) (bool, error) {
	t1, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return false, err
	}

	t2, err := time.Parse(time.RFC3339, ctrl)
	if err != nil {
		return false, err
	}

	return t1.After(t2), nil
}

func TimestampLessThanEqualsFn(v string, ctrl string) (bool, error) {
	t1, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return false, err
	}

	t2, err := time.Parse(time.RFC3339, ctrl)
	if err != nil {
		return false, err
	}

	return t1.Before(t2) || t1.Equal(t2), nil
}

func TimestampGreaterThanEqualsFn(v string, ctrl string) (bool, error) {
	t1, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return false, err
	}

	t2, err := time.Parse(time.RFC3339, ctrl)
	if err != nil {
		return false, err
	}

	return t1.After(t2) || t1.Equal(t2), nil
}

// TODO: think if these should actually be defined like this?
func AndFn(v string, ctrl string) (bool, error) { return false, nil }
func OrFn(v string, ctrl string) (bool, error)  { return false, nil }
func NotFn(v string, ctrl string) (bool, error) { return false, nil }

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
	TimestampLessThan          *time.Time             `json:"TimestampLessThan"`
	TimestampGreaterThan       *time.Time             `json:"TimestampGreaterThan"`
	TimestampLessThanEquals    *time.Time             `json:"TimestampLessThanEquals"`
	TimestampGreaterThanEquals *time.Time             `json:"TimestampGreaterThanEquals"`
	And                        []ChoiceRuleDefinition `json:"And"`
	Or                         []ChoiceRuleDefinition `json:"Or"`
	Not                        *ChoiceRuleDefinition  `json:"Not"`
}

type ChoiceRules []json.RawMessage

func (c ChoiceRules) GetDefinitionAtIndex(i int) (ChoiceRuleDefinition, error) {
	out := map[string]json.RawMessage{}
	if err := json.Unmarshal(c[i], &out); err != nil {
		return nil, nil
	}

	for k, v := range out {
		if k == "Variable" || k == "Next" {
			continue
		}

		operator, ok := validComparisonOperators[k]
		if !ok {
			return nil, errors.Errof("unknown choice state operator: %s", k)
		}

		var ruleDef ChoiceRuleDefinition
		switch operator {
		case StringEquals:
		}
	}

	return nil, nil
}
