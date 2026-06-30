package types

import "time"

// Kind represents the data type of a Value. It can be one of several predefined types, such as string, int64, uint64, float64 etc.
type Kind string

const (
	KindString    Kind = "string"
	KindInt64     Kind = "int64"
	KindUInt64    Kind = "uint64"
	KindFloat64   Kind = "float64"
	KindBool      Kind = "bool"
	KindTimestamp Kind = "timestamp"
	KindNull      Kind = "null"
)

// Value represents a single value in the dataset, consisting of a type and a value. It's a union type that can hold different types of data
type Value struct {
	Kind Kind

	I64 int64
	U64 uint64
	F64 float64
	S   string
	B   bool
	TS  time.Time
}

func (v Value) IsNull() bool {
	return v.Kind == KindNull
}

func (v Value) IsString() bool {
	return v.Kind == KindString
}

func (v Value) IsInt64() bool {
	return v.Kind == KindInt64
}

func (v Value) IsUInt64() bool {
	return v.Kind == KindUInt64
}

func (v Value) IsFloat64() bool {
	return v.Kind == KindFloat64
}

func (v Value) IsBool() bool {
	return v.Kind == KindBool
}

func NewStringValue(s string) Value {
	return Value{
		Kind: KindString,
		S:    s,
	}
}

func NewInt64Value(i int64) Value {
	return Value{
		Kind: KindInt64,
		I64:  i,
	}
}

func NewUInt64Value(u uint64) Value {
	return Value{
		Kind: KindUInt64,
		U64:  u,
	}
}

func NewFloat64Value(f float64) Value {
	return Value{
		Kind: KindFloat64,
		F64:  f,
	}
}

func NewBoolValue(b bool) Value {
	return Value{
		Kind: KindBool,
		B:    b,
	}
}

func NewNullValue() Value {
	return Value{
		Kind: KindNull,
	}
}
