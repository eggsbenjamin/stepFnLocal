package sfn

import "github.com/eggsbenjamin/stepFnLocal/state"

type PassState struct {
	def state.PassDefinition
}

func NewPassState(def state.PassDefinition) PassState {
	return PassState{
		def: def,
	}
}

func (p PassState) Run(input []byte) ([]byte, error) {
	if p.def.Result != nil {
		return p.def.Result, nil
	}

	return input, nil
}

func (p PassState) Next() string {
	return p.def.Next()
}

func (p PassState) IsEnd() bool {
	return p.def.End()
}
