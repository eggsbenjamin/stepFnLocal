//go:generate mockgen -package state -source=state.go -destination state_mock.go

package state

import (
	"time"
)

// Machine defines the standard machine API for state machine implementations
type Machine interface {
	StartExecution([]byte) (ExecutionResult, error)
}

// ExecutionResult represents the result of a state machine execution
type ExecutionResult struct {
	Input  []byte
	Output []byte
	Status string
	Start  time.Time
	End    time.Time
}

// State defines the standard state API for state machine implementations
type State interface {
	Run([]byte) ([]byte, error)
	Next() string
	IsEnd() bool
}
