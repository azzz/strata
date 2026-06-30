package expr

import (
	"errors"
	"testing"

	"github.com/azzz/strata/engine/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComparisonFilters(t *testing.T) {
	t.Parallel()

	row := types.NewRow(4)
	row.Set(0, types.NewInt64Value(10))
	row.Set(1, types.NewStringValue("beta"))
	row.Set(2, types.NewBoolValue(true))
	row.Set(3, types.NewNullValue())

	tests := []struct {
		name   string
		filter Filter
		want   bool
	}{
		{
			name:   "eq matches equal value",
			filter: &Eq{Col: 0, Value: types.NewInt64Value(10)},
			want:   true,
		},
		{
			name:   "eq returns false for different value",
			filter: &Eq{Col: 1, Value: types.NewStringValue("alpha")},
			want:   false,
		},
		{
			name:   "eq matches null value",
			filter: &Eq{Col: 3, Value: types.NewNullValue()},
			want:   true,
		},
		{
			name:   "greater than matches larger value",
			filter: &GreaterThan{Col: 0, Value: types.NewInt64Value(9)},
			want:   true,
		},
		{
			name:   "greater than returns false for equal value",
			filter: &GreaterThan{Col: 0, Value: types.NewInt64Value(10)},
			want:   false,
		},
		{
			name:   "greater than or equal matches equal value",
			filter: &GreaterThanOrEqual{Col: 0, Value: types.NewInt64Value(10)},
			want:   true,
		},
		{
			name:   "greater than or equal returns false for smaller value",
			filter: &GreaterThanOrEqual{Col: 0, Value: types.NewInt64Value(11)},
			want:   false,
		},
		{
			name:   "less than matches smaller value",
			filter: &LessThan{Col: 1, Value: types.NewStringValue("gamma")},
			want:   true,
		},
		{
			name:   "less than returns false for equal value",
			filter: &LessThan{Col: 2, Value: types.NewBoolValue(true)},
			want:   false,
		},
		{
			name:   "less than or equal matches equal value",
			filter: &LessThanOrEqual{Col: 2, Value: types.NewBoolValue(true)},
			want:   true,
		},
		{
			name:   "less than or equal returns false for greater value",
			filter: &LessThanOrEqual{Col: 0, Value: types.NewInt64Value(9)},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.filter.Match(row)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestComparisonFilters_ReturnErrors(t *testing.T) {
	t.Parallel()

	t.Run("missing column", func(t *testing.T) {
		t.Parallel()

		row := types.NewRow(1)
		row.Set(0, types.NewInt64Value(10))

		ok, err := (&Eq{Col: 1, Value: types.NewInt64Value(10)}).Match(row)
		require.Error(t, err)
		assert.False(t, ok)

		var missingColumnErr *MissingColumnError
		require.ErrorAs(t, err, &missingColumnErr)
		assert.Equal(t, types.ColumnIndex(1), missingColumnErr.Col)
	})

	t.Run("incompatible value types", func(t *testing.T) {
		t.Parallel()

		row := types.NewRow(1)
		row.Set(0, types.NewInt64Value(10))

		ok, err := (&Eq{Col: 0, Value: types.NewStringValue("10")}).Match(row)
		require.Error(t, err)
		assert.False(t, ok)

		var compareErr *types.CompareError
		require.ErrorAs(t, err, &compareErr)
		assert.Equal(t, types.KindInt64, compareErr.Left)
		assert.Equal(t, types.KindString, compareErr.Right)
	})
}

func TestLogicalFilters(t *testing.T) {
	t.Parallel()

	row := types.NewRow(1)
	row.Set(0, types.NewInt64Value(10))

	t.Run("and returns true when all filters match", func(t *testing.T) {
		t.Parallel()

		ok, err := (&And{
			Filters: []Filter{
				&GreaterThan{Col: 0, Value: types.NewInt64Value(5)},
				&LessThanOrEqual{Col: 0, Value: types.NewInt64Value(10)},
			},
		}).Match(row)

		require.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("and short circuits on false", func(t *testing.T) {
		t.Parallel()

		first := &stubFilter{result: false}
		second := &stubFilter{result: true}

		ok, err := (&And{Filters: []Filter{first, second}}).Match(row)
		require.NoError(t, err)
		assert.False(t, ok)
		assert.Equal(t, 1, first.calls)
		assert.Equal(t, 0, second.calls)
	})

	t.Run("and returns first error", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New("boom")
		first := &stubFilter{err: expectedErr}
		second := &stubFilter{result: true}

		ok, err := (&And{Filters: []Filter{first, second}}).Match(row)
		require.ErrorIs(t, err, expectedErr)
		assert.False(t, ok)
		assert.Equal(t, 1, first.calls)
		assert.Equal(t, 0, second.calls)
	})

	t.Run("or returns true when any filter matches", func(t *testing.T) {
		t.Parallel()

		ok, err := (&Or{
			Filters: []Filter{
				&LessThan{Col: 0, Value: types.NewInt64Value(5)},
				&Eq{Col: 0, Value: types.NewInt64Value(10)},
			},
		}).Match(row)

		require.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("or short circuits on true", func(t *testing.T) {
		t.Parallel()

		first := &stubFilter{result: true}
		second := &stubFilter{result: false}

		ok, err := (&Or{Filters: []Filter{first, second}}).Match(row)
		require.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, 1, first.calls)
		assert.Equal(t, 0, second.calls)
	})

	t.Run("or returns first error", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New("boom")
		first := &stubFilter{err: expectedErr}
		second := &stubFilter{result: true}

		ok, err := (&Or{Filters: []Filter{first, second}}).Match(row)
		require.ErrorIs(t, err, expectedErr)
		assert.False(t, ok)
		assert.Equal(t, 1, first.calls)
		assert.Equal(t, 0, second.calls)
	})

	t.Run("not negates nested result", func(t *testing.T) {
		t.Parallel()

		ok, err := (&Not{Filter: &Eq{Col: 0, Value: types.NewInt64Value(11)}}).Match(row)
		require.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("not propagates error", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New("boom")
		filter := &stubFilter{err: expectedErr}

		ok, err := (&Not{Filter: filter}).Match(row)
		require.ErrorIs(t, err, expectedErr)
		assert.False(t, ok)
		assert.Equal(t, 1, filter.calls)
	})
}

type stubFilter struct {
	result bool
	err    error
	calls  int
}

func (s *stubFilter) Match(types.Row) (bool, error) {
	s.calls++

	if s.err != nil {
		return false, s.err
	}

	return s.result, nil
}
