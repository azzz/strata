package executor

import (
	"context"
	"errors"
	"testing"

	"github.com/azzz/strata/engine/storage"
	"github.com/azzz/strata/engine/types"
	"github.com/azzz/strata/engine/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockStorage struct {
	scanners []storage.Scanner
	errs     []error
	calls    int
	requests []storage.ScanRequest
}

func (m *mockStorage) NewScanner(ctx context.Context, req storage.ScanRequest) (storage.Scanner, error) {
	m.requests = append(m.requests, req)

	idx := m.calls
	m.calls++

	if idx < len(m.errs) && m.errs[idx] != nil {
		return nil, m.errs[idx]
	}

	return m.scanners[idx], nil
}

type mockScanner struct {
	rows       []types.Row
	nextErr    error
	openErr    error
	closeErr   error
	cursor     int
	openCalls  int
	closeCalls int
	row        types.Row
	err        error
}

func (m *mockScanner) Open() error {
	m.openCalls++
	return m.openErr
}

func (m *mockScanner) Next(ctx context.Context) bool {
	if m.cursor >= len(m.rows) {
		if m.nextErr != nil {
			m.err = m.nextErr
		}

		return false
	}

	m.row = m.rows[m.cursor]
	m.cursor++

	return true
}

func (m *mockScanner) Err() error     { return m.err }
func (m *mockScanner) Row() types.Row { return m.row }

func (m *mockScanner) Close() error {
	m.closeCalls++
	return m.closeErr
}

func TestSeqScan(t *testing.T) {
	requests := []storage.ScanRequest{
		{URI: "first.jsonl", Format: storage.FormatJSONL},
		{URI: "second.jsonl", Format: storage.FormatJSONL},
	}

	t.Run("reads rows from all requests sequentially", func(t *testing.T) {
		first := &mockScanner{
			rows: []types.Row{
				newRow(types.NewInt64Value(4)),
				newRow(types.NewInt64Value(8)),
			},
		}
		second := &mockScanner{
			rows: []types.Row{
				newRow(types.NewInt64Value(15)),
			},
		}
		input := &mockStorage{scanners: []storage.Scanner{first, second}}

		op := NewSeqScan(input, requests)
		retrieved, err := utils.Collect(context.Background(), op)

		require.NoError(t, err)
		assert.Equal(t,
			[]types.Row{
				newRow(types.NewInt64Value(4)),
				newRow(types.NewInt64Value(8)),
				newRow(types.NewInt64Value(15)),
			}, retrieved)
		assert.Equal(t, requests, input.requests)
		assert.Equal(t, 1, first.openCalls)
		assert.Equal(t, 1, second.openCalls)
		assert.Equal(t, 1, first.closeCalls)
		assert.Equal(t, 1, second.closeCalls)
	})

	t.Run("new scanner error is returned", func(t *testing.T) {
		expectedErr := errors.New("boom")
		input := &mockStorage{
			errs: []error{expectedErr},
		}

		op := NewSeqScan(input, requests[:1])
		retrieved, err := utils.Collect(context.Background(), op)

		require.ErrorIs(t, err, expectedErr)
		assert.Empty(t, retrieved)
		assert.Equal(t, 1, input.calls)
	})

	t.Run("open error is returned", func(t *testing.T) {
		expectedErr := errors.New("boom")
		scanner := &mockScanner{openErr: expectedErr}
		input := &mockStorage{scanners: []storage.Scanner{scanner}}

		op := NewSeqScan(input, requests[:1])
		retrieved, err := utils.Collect(context.Background(), op)

		require.ErrorIs(t, err, expectedErr)
		assert.Empty(t, retrieved)
		assert.Equal(t, 1, scanner.openCalls)
		assert.Equal(t, 0, scanner.closeCalls)
	})

	t.Run("scanner error is returned with collected rows", func(t *testing.T) {
		expectedErr := errors.New("boom")
		scanner := &mockScanner{
			rows: []types.Row{
				newRow(types.NewInt64Value(4)),
			},
			nextErr: expectedErr,
		}
		input := &mockStorage{scanners: []storage.Scanner{scanner}}

		op := NewSeqScan(input, requests[:1])
		retrieved, err := utils.Collect(context.Background(), op)

		require.ErrorIs(t, err, expectedErr)
		assert.Equal(t,
			[]types.Row{
				newRow(types.NewInt64Value(4)),
			}, retrieved)
		assert.Equal(t, 0, scanner.closeCalls)
	})

	t.Run("close error is returned after scanner is exhausted", func(t *testing.T) {
		expectedErr := errors.New("boom")
		scanner := &mockScanner{
			rows: []types.Row{
				newRow(types.NewInt64Value(4)),
			},
			closeErr: expectedErr,
		}
		input := &mockStorage{scanners: []storage.Scanner{scanner}}

		op := NewSeqScan(input, requests[:1])
		retrieved, err := utils.Collect(context.Background(), op)

		require.ErrorIs(t, err, expectedErr)
		assert.Equal(t,
			[]types.Row{
				newRow(types.NewInt64Value(4)),
			}, retrieved)
		assert.Equal(t, 1, scanner.closeCalls)
	})
}
