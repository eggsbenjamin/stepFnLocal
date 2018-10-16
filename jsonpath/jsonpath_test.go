// +build unit

package jsonpath_test

import (
	"testing"

	"github.com/eggsbenjamin/stepFnLocal/jsonpath"
	"github.com/stretchr/testify/require"
)

func TestExpression(t *testing.T) {
	t.Run("invalid JSON path", func(t *testing.T) {
		invalid := "invalid"
		_, err := jsonpath.NewExpression(invalid)
		require.Error(t, err)
	})

	t.Run("Search", func(t *testing.T) {
		t.Run("root element", func(t *testing.T) {
			tests := []struct {
				title string
				input string
			}{
				{
					"string",
					`"hello"`,
				},
				{
					"number",
					`33`,
				},
				{
					"object",
					`{"hello":"world"}`,
				},
				{
					"array",
					`["one",2,{}]`,
				},
			}

			for _, tt := range tests {
				t.Run(tt.title, func(t *testing.T) {
					exp, err := jsonpath.NewExpression("$")
					require.NoError(t, err)

					result, err := exp.Search([]byte(tt.input))
					require.NoError(t, err)
					require.Equal(t, string(tt.input), string(result))
				})
			}
		})

		t.Run("not found", func(t *testing.T) {
			exp, err := jsonpath.NewExpression("$.hello")
			require.NoError(t, err)

			result, err := exp.Search([]byte(`{}`))
			require.NoError(t, err)
			require.Nil(t, result)
		})
	})
}
