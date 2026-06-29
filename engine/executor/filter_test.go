package executor

import (
	"context"
	"errors"
	"testing"

	"github.com/azzz/strata/engine/expr"
	"github.com/azzz/strata/engine/types"
	"github.com/azzz/strata/engine/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockFilter struct {
	match func(types.Row) (bool, error)
	calls int
}

func (m *mockFilter) Match(row types.Row) (bool, error) {
	m.calls++
	return m.match(row)
}

func TestFilter(t *testing.T) {
	newInput := func() *MockOperator {
		return new(MockOperator).
			addRow(newRow(types.NewInt64Value(4))).
			addRow(newRow(types.NewInt64Value(8))).
			addRow(newRow(types.NewInt64Value(15)))
	}

	t.Run("all rows match", func(t *testing.T) {
		input := newInput()

		op := NewFilter(input, &expr.GreaterThan{Col: 0, Value: types.NewInt64Value(0)})
		retrieved, err := utils.Collect(context.Background(), op)

		require.NoError(t, err)
		assert.Equal(t,
			[]types.Row{
				newRow(types.NewInt64Value(4)),
				newRow(types.NewInt64Value(8)),
				newRow(types.NewInt64Value(15)),
			}, retrieved)
	})

	t.Run("non matching rows are skipped", func(t *testing.T) {
		input := newInput()

		op := NewFilter(input, &expr.GreaterThan{Col: 0, Value: types.NewInt64Value(8)})
		retrieved, err := utils.Collect(context.Background(), op)

		require.NoError(t, err)
		assert.Equal(t,
			[]types.Row{
				newRow(types.NewInt64Value(15)),
			}, retrieved)
	})

	t.Run("no rows match", func(t *testing.T) {
		input := newInput()

		op := NewFilter(input, &expr.GreaterThan{Col: 0, Value: types.NewInt64Value(20)})
		retrieved, err := utils.Collect(context.Background(), op)

		require.NoError(t, err)
		assert.Empty(t, retrieved)
	})

	t.Run("upstream error is returned with collected rows", func(t *testing.T) {
		expectedErr := errors.New("boom")
		input := new(MockOperator).
			addRow(newRow(types.NewInt64Value(4))).
			addRow(newRow(types.NewInt64Value(8))).
			addErr(expectedErr)

		op := NewFilter(input, &expr.GreaterThan{Col: 0, Value: types.NewInt64Value(0)})
		retrieved, err := utils.Collect(context.Background(), op)

		require.ErrorIs(t, err, expectedErr)
		assert.Equal(t,
			[]types.Row{
				newRow(types.NewInt64Value(4)),
				newRow(types.NewInt64Value(8)),
			}, retrieved)
	})

	t.Run("filter error stops scan immediately", func(t *testing.T) {
		expectedErr := errors.New("boom")
		input := newInput()
		filter := &mockFilter{
			match: func(row types.Row) (bool, error) {
				value, ok := row.Get(0)
				if !ok {
					return false, errors.New("missing value")
				}

				if value == types.NewInt64Value(8) {
					return false, expectedErr
				}

				return true, nil
			},
		}

		op := NewFilter(input, filter)
		retrieved, err := utils.Collect(context.Background(), op)

		require.ErrorIs(t, err, expectedErr)
		assert.Equal(t,
			[]types.Row{
				newRow(types.NewInt64Value(4)),
			}, retrieved)
		assert.Equal(t, 2, input.cursor, "filter should stop reading upstream after evaluation error")
		assert.Equal(t, 2, filter.calls)
	})
}
