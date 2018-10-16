package jsonpath

import (
	"encoding/json"

	"github.com/oliveagle/jsonpath"
)

type Expression interface {
	Search(json []byte) ([]byte, error)
}

func NewExpression(input string) (Expression, error) {
	exp := expression(input)
	return exp, exp.compile()
}

type expression string

func (e expression) Search(inputJSON []byte) ([]byte, error) {
	var input interface{}
	if err := json.Unmarshal(inputJSON, &input); err != nil {
		return []byte{}, err
	}

	res, _ := jsonpath.JsonPathLookup(input, string(e))
	outputJSON, err := json.Marshal(res)
	if err != nil {
		return []byte{}, err
	}

	if string(outputJSON) == "null" {
		return nil, nil
	}
	return outputJSON, nil
}

func (e expression) compile() error {
	_, err := jsonpath.Compile(string(e))
	return err
}
