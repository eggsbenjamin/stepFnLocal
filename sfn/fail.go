package sfn

import "github.com/eggsbenjamin/stepFnLocal/state"

type FailState struct {
	def state.FailDefinition
}

func NewFailState(def state.FailDefinition) FailState {
	return FailState{
		def: def,
	}
}

func (p FailState) Run(input []byte) ([]byte, error) {
	return input, nil
}

func (p FailState) Next() string {
	return ""
}

func (p FailState) IsEnd() bool {
	return true
}
