package jsonl

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/azzz/strata/engine/types"
)

// Decoder is responsible for decoding JSON lines into Rows based on a given schema.
type Decoder struct {
	schema types.Schema
}

// NewDecoder creates a new Decoder with the provided schema.
func NewDecoder(schema types.Schema) *Decoder {
	return &Decoder{
		schema: schema,
	}
}

// Decode takes a JSON line as input and decodes it into a Row according to the schema.
func (d *Decoder) Decode(data []byte) (types.Row, error) {
	var fields map[string]any

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()

	if err := decoder.Decode(&fields); err != nil {
		return types.Row{}, fmt.Errorf("failed to decode json: %w", err)
	}

	row := types.NewRow(d.schema.Len())

	for idx, field := range d.schema.Columns() {
		value, exists := fields[field.Name]
		if !exists {
			row.Set(idx, types.NewNullValue())
			continue
		}

		decodedValue, err := d.decodeValue(field.Name, field.Type, value)
		if err != nil {
			return types.Row{}, err
		}

		row.Set(idx, decodedValue)
	}

	return row, nil
}

func (d *Decoder) decodeValue(name string, kind types.Kind, raw any) (types.Value, error) {
	switch kind {
	case types.KindString:
		if val, ok := raw.(string); !ok {
			return types.NewNullValue(), &types.TypeError{
				Field:    name,
				Expected: types.KindString,
			}
		} else {
			return types.NewStringValue(val), nil
		}
	case types.KindInt64:
		if n, ok := raw.(json.Number); !ok { // JSON numbers are float64
			return types.NewNullValue(), &types.TypeError{
				Field:    name,
				Expected: types.KindInt64,
			}
		} else {
			val, err := n.Int64()
			if err != nil {
				return types.NewNullValue(), &types.TypeError{
					Field:    name,
					Expected: types.KindInt64,
				}
			}

			return types.NewInt64Value(val), nil
		}
	case types.KindFloat64:
		if val, ok := raw.(json.Number); !ok {
			return types.NewNullValue(), &types.TypeError{
				Field:    name,
				Expected: types.KindFloat64,
			}
		} else {
			val, err := val.Float64()
			if err != nil {
				return types.NewNullValue(), &types.TypeError{
					Field:    name,
					Expected: types.KindFloat64,
				}
			}

			return types.NewFloat64Value(val), nil
		}
	case types.KindBool:
		if val, ok := raw.(bool); !ok {
			return types.NewNullValue(), &types.TypeError{
				Field:    name,
				Expected: types.KindBool,
			}
		} else {
			return types.NewBoolValue(val), nil
		}
	case types.KindNull:
		if raw != nil {
			return types.NewNullValue(), &types.TypeError{
				Field:    name,
				Expected: types.KindNull,
			}
		} else {
			return types.NewNullValue(), nil
		}
	case types.KindUInt64:
		if n, ok := raw.(json.Number); !ok {
			return types.NewNullValue(), &types.TypeError{
				Field:    name,
				Expected: types.KindUInt64,
			}
		} else {
			val, err := n.Int64()
			if err != nil || val < 0 {
				return types.NewNullValue(), &types.TypeError{
					Field:    name,
					Expected: types.KindUInt64,
				}
			}

			return types.NewUInt64Value(uint64(val)), nil
		}

	case types.KindTimestamp:
		return types.NewNullValue(), &types.TypeNotImplementedError{
			Field: name,
			Kind:  kind,
		}

	default:
		return types.NewNullValue(), &types.TypeError{
			Field:    name,
			Expected: kind,
		}
	}
}
