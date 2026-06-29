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

func TestLimit(t *testing.T) {
	newInput := func() *MockOperator {
		return new(MockOperator).
			addRow(newRow(types.NewInt64Value(4))).
			addRow(newRow(types.NewInt64Value(8))).
			addRow(newRow(types.NewInt64Value(15)))
	}

	t.Run("no limit", func(t *testing.T) {
		input := newInput()

		op := NewLimit(input, 0)
		retrieved, err := utils.Collect(context.Background(), op)

		require.NoError(t, err)
		assert.Len(t, retrieved, 3)
		assert.Equal(t,
			[]types.Row{
				newRow(types.NewInt64Value(4)),
				newRow(types.NewInt64Value(8)),
				newRow(types.NewInt64Value(15)),
			}, retrieved)
	})

	t.Run("limit is higher than number of rows in the source", func(t *testing.T) {
		input := newInput()

		op := NewLimit(input, 5)
		retrieved, err := utils.Collect(context.Background(), op)

		require.NoError(t, err)
		assert.Len(t, retrieved, 3)
		assert.Equal(t,
			[]types.Row{
				newRow(types.NewInt64Value(4)),
				newRow(types.NewInt64Value(8)),
				newRow(types.NewInt64Value(15)),
			}, retrieved)
	})

	t.Run("limit is equal to number of rows in the source", func(t *testing.T) {
		input := newInput()

		op := NewLimit(input, 3)
		retrieved, err := utils.Collect(context.Background(), op)

		require.NoError(t, err)
		assert.Len(t, retrieved, 3)
		assert.Equal(t,
			[]types.Row{
				newRow(types.NewInt64Value(4)),
				newRow(types.NewInt64Value(8)),
				newRow(types.NewInt64Value(15)),
			}, retrieved)
		assert.Equal(t, 3, input.cursor)
	})

	t.Run("limit is lower than number of rows in the source", func(t *testing.T) {
		input := newInput()

		op := NewLimit(input, 2)
		retrieved, err := utils.Collect(context.Background(), op)

		require.NoError(t, err)
		assert.Len(t, retrieved, 2)
		assert.Equal(t,
			[]types.Row{
				newRow(types.NewInt64Value(4)),
				newRow(types.NewInt64Value(8)),
			}, retrieved)
		assert.Equal(t, 2, input.cursor, "limit should stop reading upstream once enough rows were emitted")
	})

	t.Run("next after reaching the limit does not read upstream", func(t *testing.T) {
		input := newInput()

		op := NewLimit(input, 2)

		require.True(t, op.Next(context.Background()))
		assert.Equal(t, newRow(types.NewInt64Value(4)), op.Row())
		require.True(t, op.Next(context.Background()))
		assert.Equal(t, newRow(types.NewInt64Value(8)), op.Row())
		require.False(t, op.Next(context.Background()))
		require.False(t, op.Next(context.Background()))

		assert.NoError(t, op.Err())
		assert.Equal(t, 2, input.cursor, "limit should not read upstream after reaching the bound")
	})

	t.Run("upstream error before reaching the limit is returned", func(t *testing.T) {
		expectedErr := errors.New("boom")
		input := new(MockOperator).
			addRow(newRow(types.NewInt64Value(4))).
			addErr(expectedErr)

		op := NewLimit(input, 5)
		retrieved, err := utils.Collect(context.Background(), op)

		require.ErrorIs(t, err, expectedErr)
		assert.Equal(t,
			[]types.Row{
				newRow(types.NewInt64Value(4)),
			}, retrieved)
	})

	t.Run("upstream error after reaching the limit is ignored", func(t *testing.T) {
		expectedErr := errors.New("boom")
		input := new(MockOperator).
			addRow(newRow(types.NewInt64Value(4))).
			addRow(newRow(types.NewInt64Value(8))).
			addErr(expectedErr)

		op := NewLimit(input, 2)
		retrieved, err := utils.Collect(context.Background(), op)

		require.NoError(t, err)
		assert.Equal(t,
			[]types.Row{
				newRow(types.NewInt64Value(4)),
				newRow(types.NewInt64Value(8)),
			}, retrieved)
		assert.Equal(t, 2, input.cursor, "limit should not read rows beyond the configured bound")
	})
}
