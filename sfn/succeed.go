package sfn

import "github.com/eggsbenjamin/stepFnLocal/state"

type SucceedState struct {
	def state.SucceedDefinition
}

func (p SucceedState) Run(input []byte) ([]byte, error) {
	return input, nil
}

func (p SucceedState) Next() string {
	return ""
}

func (p SucceedState) IsEnd() bool {
	return true
}
