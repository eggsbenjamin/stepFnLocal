package sfn_test

import (
	"testing"

	"github.com/eggsbenjamin/stepFnLocal/sfn"
	"github.com/eggsbenjamin/stepFnLocal/state"
	"github.com/stretchr/testify/require"
)

func TestPassState(t *testing.T) {
	t.Run("no result specified", func(t *testing.T) {
		input := []byte("test")
		pass := sfn.NewPassState(state.PassDefinition{})

		result, err := pass.Run(input)
		require.NoError(t, err)
		require.Equal(t, input, result)
	})

	t.Run("result specified", func(t *testing.T) {
		input := []byte("test")
		def := state.PassDefinition{
			Result: []byte("result"),
		}
		pass := sfn.NewPassState(def)

		result, err := pass.Run(input)
		require.NoError(t, err)
		require.Equal(t, []byte(def.Result), result)
	})
}
