package sfn

import (
	"encoding/json"
	"sync"

	"github.com/eggsbenjamin/stepFnLocal/state"
)

type ParallelState struct {
	def           state.ParallelDefinition
	stateMachines []StepFunction
}

func NewParallelState(def state.ParallelDefinition, stateMachines ...StepFunction) ParallelState {
	return ParallelState{
		def:           def,
		stateMachines: stateMachines,
	}
}

func (p ParallelState) Run(input []byte) ([]byte, error) {
	type stateMachineResult struct {
		Index  int // order of output is important
		Output json.RawMessage
	}

	var wg sync.WaitGroup
	stateMachineResults := make(chan stateMachineResult)
	errs := make(chan error)

	wg.Add(len(p.stateMachines))
	go func() {
		wg.Wait()
		close(stateMachineResults)
	}()

	for i, stateMachine := range p.stateMachines {
		go func(index int, stateMachine StepFunction) {
			defer wg.Done()

			result, err := stateMachine.StartExecution(input)
			if err != nil {
				errs <- err
				return
			}
			stateMachineResults <- stateMachineResult{
				Index:  index,
				Output: result.Output,
			}
		}(i, stateMachine)
	}

	results := make([]json.RawMessage, len(p.stateMachines))
	running := true
	for running {
		select {
		case result, ok := <-stateMachineResults:
			if !ok {
				running = false
				break
			}
			results[result.Index] = result.Output
		case err, ok := <-errs:
			if ok {
				return []byte{}, err
			}
			errs = nil
		}
	}

	return json.Marshal(results)
}

func (p ParallelState) Next() string {
	return p.def.NextState
}

func (p ParallelState) IsEnd() bool {
	return p.def.EndState
}
