package executor

import (
	"context"

	"github.com/azzz/strata/engine/types"
)

type MockOperator struct {
	source []struct {
		row types.Row
		err error
	}
	cursor int

	row types.Row
	err error
}

func (m *MockOperator) addRow(row types.Row) *MockOperator {
	m.source = append(m.source,
		struct {
			row types.Row
			err error
		}{row: row},
	)

	return m
}

func (m *MockOperator) addErr(err error) *MockOperator {
	m.source = append(m.source,
		struct {
			row types.Row
			err error
		}{err: err},
	)

	return m
}

func (m *MockOperator) Row() types.Row { return m.row }
func (m *MockOperator) Err() error     { return m.err }
func (m *MockOperator) Close() error   { return nil }

func (m *MockOperator) Next(ctx context.Context) bool {
	for m.cursor < len(m.source) {
		fetcher := m.source[m.cursor]

		row, err := fetcher.row, fetcher.err
		if err != nil {
			m.err = err
			return false
		} else {
			m.row = row
			m.cursor++

			return true
		}
	}

	return false
}

func newRow(values ...types.Value) types.Row {
	row := types.NewRow(len(values))
	for i, v := range values {
		row.Set(i, v)
	}

	return row
}
