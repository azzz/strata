package jsonl

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/azzz/strata/engine/types"
)

func TestJSONDecoder_decodeValue(t *testing.T) {
	t.Parallel()

	decoder := &Decoder{}

	tests := []struct {
		name        string
		field       string
		kind        types.Kind
		raw         any
		want        types.Value
		wantTypeErr *types.TypeError
	}{
		{
			name:  "string",
			field: "name",
			kind:  types.KindString,
			raw:   "alice",
			want:  types.NewStringValue("alice"),
		},
		{
			name:  "int64",
			field: "age",
			kind:  types.KindInt64,
			raw:   json.Number("42"),
			want:  types.NewInt64Value(42),
		},
		{
			name:  "int64 invalid number",
			field: "age",
			kind:  types.KindInt64,
			raw:   json.Number("42.5"),
			want:  types.NewNullValue(),
			wantTypeErr: &types.TypeError{
				Field:    "age",
				Expected: types.KindInt64,
			},
		},
		{
			name:  "float64",
			field: "score",
			kind:  types.KindFloat64,
			raw:   json.Number("3.14"),
			want:  types.NewFloat64Value(3.14),
		},
		{
			name:  "float64 invalid number",
			field: "score",
			kind:  types.KindFloat64,
			raw:   json.Number("not-a-number"),
			want:  types.NewNullValue(),
			wantTypeErr: &types.TypeError{
				Field:    "score",
				Expected: types.KindFloat64,
			},
		},
		{
			name:  "bool",
			field: "active",
			kind:  types.KindBool,
			raw:   true,
			want:  types.NewBoolValue(true),
		},
		{
			name:  "null",
			field: "deleted_at",
			kind:  types.KindNull,
			raw:   nil,
			want:  types.NewNullValue(),
		},
		{
			name:  "string type mismatch",
			field: "name",
			kind:  types.KindString,
			raw:   123,
			want:  types.NewNullValue(),
			wantTypeErr: &types.TypeError{
				Field:    "name",
				Expected: types.KindString,
			},
		},
		{
			name:  "int64 type mismatch",
			field: "age",
			kind:  types.KindInt64,
			raw:   float64(42),
			want:  types.NewNullValue(),
			wantTypeErr: &types.TypeError{
				Field:    "age",
				Expected: types.KindInt64,
			},
		},
		{
			name:  "float64 type mismatch",
			field: "score",
			kind:  types.KindFloat64,
			raw:   float64(3.14),
			want:  types.NewNullValue(),
			wantTypeErr: &types.TypeError{
				Field:    "score",
				Expected: types.KindFloat64,
			},
		},
		{
			name:  "bool type mismatch",
			field: "active",
			kind:  types.KindBool,
			raw:   "true",
			want:  types.NewNullValue(),
			wantTypeErr: &types.TypeError{
				Field:    "active",
				Expected: types.KindBool,
			},
		},
		{
			name:  "null type mismatch",
			field: "deleted_at",
			kind:  types.KindNull,
			raw:   "2024-01-01",
			want:  types.NewNullValue(),
			wantTypeErr: &types.TypeError{
				Field:    "deleted_at",
				Expected: types.KindNull,
			},
		},
		{
			name:  "unsupported kind",
			field: "created_at",
			kind:  types.KindTimestamp,
			raw:   "2024-01-01T00:00:00Z",
			want:  types.NewNullValue(),
			wantTypeErr: &types.TypeError{
				Field: "created_at",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := decoder.decodeValue(tt.field, tt.kind, tt.raw)
			if got != tt.want {
				t.Fatalf("decodeValue() value = %#v, want %#v", got, tt.want)
			}

			if tt.wantTypeErr == nil {
				if err != nil {
					t.Fatalf("decodeValue() error = %v, want nil", err)
				}

				return
			}

			if err == nil {
				t.Fatal("decodeValue() error = nil, want type error")
			}

			var typeErr *types.TypeError
			if !errors.As(err, &typeErr) {
				t.Fatalf("decodeValue() error = %T, want *types.TypeError", err)
			}

			if *typeErr != *tt.wantTypeErr {
				t.Fatalf("decodeValue() type error = %#v, want %#v", *typeErr, *tt.wantTypeErr)
			}
		})
	}
}

func TestJSONDecoder_Decode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		schema      types.Schema
		data        []byte
		want        []types.Value
		wantTypeErr *types.TypeError
		wantErr     bool
	}{
		{
			name: "decode row and fill missing field with null",
			schema: types.NewSchema(
				types.Column{Name: "name", Type: types.KindString},
				types.Column{Name: "active", Type: types.KindBool},
				types.Column{Name: "deleted_at", Type: types.KindNull},
				types.Column{Name: "missing", Type: types.KindString},
			),
			data: []byte(`{"name":"alice","active":true,"deleted_at":null}`),
			want: []types.Value{
				types.NewStringValue("alice"),
				types.NewBoolValue(true),
				types.NewNullValue(),
				types.NewNullValue(),
			},
		},
		{
			name: "invalid json",
			schema: types.NewSchema(
				types.Column{Name: "name", Type: types.KindString},
			),
			data:    []byte(`{"name":`),
			wantErr: true,
		},
		{
			name: "numeric fields decode from json numbers",
			schema: types.NewSchema(
				types.Column{Name: "age", Type: types.KindInt64},
				types.Column{Name: "load", Type: types.KindFloat64},
			),
			data: []byte(`{"age":42,"load":0.5}`),
			want: []types.Value{
				types.NewInt64Value(42),
				types.NewFloat64Value(0.5),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			decoder := NewDecoder(tt.schema)
			got, err := decoder.Decode(tt.data)

			if tt.wantErr {
				if err == nil {
					t.Fatal("Decode() error = nil, want non-nil")
				}

				return
			}

			if tt.wantTypeErr != nil {
				if err == nil {
					t.Fatal("Decode() error = nil, want type error")
				}

				var typeErr *types.TypeError
				if !errors.As(err, &typeErr) {
					t.Fatalf("Decode() error = %T, want *types.TypeError", err)
				}

				if *typeErr != *tt.wantTypeErr {
					t.Fatalf("Decode() type error = %#v, want %#v", *typeErr, *tt.wantTypeErr)
				}

				return
			}

			if err != nil {
				t.Fatalf("Decode() error = %v, want nil", err)
			}

			for i, want := range tt.want {
				value, ok := got.Get(i)
				if !ok {
					t.Fatalf("Decode() row missing value at index %d", i)
				}

				if value != want {
					t.Fatalf("Decode() row[%d] = %#v, want %#v", i, value, want)
				}
			}
		})
	}
}
