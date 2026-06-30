package storage

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/azzz/strata/engine/types"
)

func TestJSONLScanner_Integration(t *testing.T) {
	t.Parallel()

	scanner := NewJSONLScanner(
		filepath.Join("testdata", "part_1.jsonl"),
		types.NewSchema(
			types.Column{Name: "node_id", Type: types.KindInt64},
			types.Column{Name: "active", Type: types.KindBool},
			types.Column{Name: "load", Type: types.KindFloat64},
		),
	)

	if err := scanner.Open(); err != nil {
		t.Fatalf("Open() error = %v", err)
	}

	defer func() {
		if err := scanner.Close(); err != nil {
			t.Fatalf("Close() error = %v", err)
		}
	}()

	ctx := context.Background()
	wantRows := [][]types.Value{
		{
			types.NewInt64Value(17),
			types.NewBoolValue(true),
			types.NewFloat64Value(0.5),
		},
		{
			types.NewInt64Value(17),
			types.NewBoolValue(true),
			types.NewFloat64Value(0.7),
		},
		{
			types.NewInt64Value(23),
			types.NewBoolValue(false),
			types.NewFloat64Value(1.0),
		},
	}

	for rowIndex, wantRow := range wantRows {
		if !scanner.Next(ctx) {
			t.Fatalf("Next() = false at row %d, err = %v", rowIndex, scanner.Err())
		}

		row := scanner.Row()
		for columnIndex, wantValue := range wantRow {
			gotValue, ok := row.Get(types.ColumnIndex(columnIndex))
			if !ok {
				t.Fatalf("Row().Get(%d) = !ok at row %d", columnIndex, rowIndex)
			}

			if gotValue != wantValue {
				t.Fatalf("Row().Get(%d) = %#v at row %d, want %#v", columnIndex, gotValue, rowIndex, wantValue)
			}
		}
	}

	if scanner.Next(ctx) {
		t.Fatal("Next() = true after EOF, want false")
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("Err() = %v after EOF, want nil", err)
	}
}
