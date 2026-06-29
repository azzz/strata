package executor

import (
	"context"
	"errors"
	"testing"

	"github.com/azzz/strata/engine/types"
	"github.com/azzz/strata/engine/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOffset(t *testing.T) {
	newInput := func() *MockOperator {
		return new(MockOperator).
			addRow(newRow(types.NewInt64Value(4))).
			addRow(newRow(types.NewInt64Value(8))).
			addRow(newRow(types.NewInt64Value(15)))
	}

	t.Run("no offset", func(t *testing.T) {
		input := newInput()

		op := NewOffset(input, 0)
		retrieved, err := utils.Collect(context.Background(), op)

		require.NoError(t, err)
		assert.Equal(t,
			[]types.Row{
				newRow(types.NewInt64Value(4)),
				newRow(types.NewInt64Value(8)),
				newRow(types.NewInt64Value(15)),
			}, retrieved)
	})

	t.Run("offset is higher than number of rows in the source", func(t *testing.T) {
		input := newInput()

		op := NewOffset(input, 5)
		retrieved, err := utils.Collect(context.Background(), op)

		require.NoError(t, err)
		assert.Empty(t, retrieved)
		assert.Equal(t, 3, input.cursor)
	})

	t.Run("offset is equal to number of rows in the source", func(t *testing.T) {
		input := newInput()

		op := NewOffset(input, 3)
		retrieved, err := utils.Collect(context.Background(), op)

		require.NoError(t, err)
		assert.Empty(t, retrieved)
		assert.Equal(t, 3, input.cursor)
	})

	t.Run("offset is lower than number of rows in the source", func(t *testing.T) {
		input := newInput()

		op := NewOffset(input, 2)
		retrieved, err := utils.Collect(context.Background(), op)

		require.NoError(t, err)
		assert.Equal(t,
			[]types.Row{
				newRow(types.NewInt64Value(15)),
			}, retrieved)
		assert.Equal(t, 3, input.cursor)
	})

	t.Run("next after source is exhausted does not change state", func(t *testing.T) {
		input := newInput()

		op := NewOffset(input, 3)

		require.False(t, op.Next(context.Background()))
		require.False(t, op.Next(context.Background()))

		assert.NoError(t, op.Err())
		assert.Equal(t, 3, input.cursor, "offset should not read past the end of the source")
	})

	t.Run("upstream error while skipping rows is returned", func(t *testing.T) {
		expectedErr := errors.New("boom")
		input := new(MockOperator).
			addRow(newRow(types.NewInt64Value(4))).
			addErr(expectedErr)

		op := NewOffset(input, 2)
		retrieved, err := utils.Collect(context.Background(), op)

		require.ErrorIs(t, err, expectedErr)
		assert.Empty(t, retrieved)
	})

	t.Run("upstream error after offset is returned with collected rows", func(t *testing.T) {
		expectedErr := errors.New("boom")
		input := new(MockOperator).
			addRow(newRow(types.NewInt64Value(4))).
			addRow(newRow(types.NewInt64Value(8))).
			addErr(expectedErr)

		op := NewOffset(input, 1)
		retrieved, err := utils.Collect(context.Background(), op)

		require.ErrorIs(t, err, expectedErr)
		assert.Equal(t,
			[]types.Row{
				newRow(types.NewInt64Value(8)),
			}, retrieved)
	})
}
