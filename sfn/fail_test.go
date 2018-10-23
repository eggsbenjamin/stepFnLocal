package sfn_test

import (
	"testing"

	"github.com/eggsbenjamin/stepFnLocal/sfn"
	"github.com/eggsbenjamin/stepFnLocal/state"
	"github.com/stretchr/testify/require"
)

func TestFailState(t *testing.T) {
	input := []byte("input")
	def := state.FailDefinition{
		Error: "test error",
		Cause: "test cause",
	}
	failState := sfn.NewFailState(def)

	expectedErr := state.NewError(
		def.Error,
		def.Cause,
	)

	result, err := failState.Run(input)
	require.Equal(t, input, result) // should never modify it's input
	require.Equal(t, "", failState.Next())
	require.Equal(t, expectedErr, err)
}
