package state

import (
	"time"
)

type Machine interface {
	StartExecution([]byte) (ExecutionResult, error)
}

type ExecutionResult struct {
	Input  json.RawMessage
	Output string
	Status string
	Start  time.Time
	End    time.Time
}

type State interface {
	Run(input []byte) ([]byte, error)
	Next() string
	End() bool
}

type TaskState struct {
	def TaskDefinition
}

func (t TaskState) Run(input []byte) ([]byte, error) {
	return nil, nil
}

func (t TaskState) Next() string {
	return t.def.Next()
}
