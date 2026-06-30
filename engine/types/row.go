package types

// Row represents a single row of data, consisting of a slice of Values.
type Row struct {
	values []Value
}

func NewRow(len int) Row {
	return Row{
		values: make([]Value, len),
	}
}

// Get the value at the specified index in the row. Returns the value and a boolean indicating if the index is valid.
func (r Row) Get(idx ColumnIndex) (Value, bool) {
	if idx < 0 || int(idx) >= len(r.values) {
		return Value{}, false
	}

	return r.values[idx], true
}

// Set the value at the specified index in the row. If the index is out of bounds, the function does nothing.
func (r Row) Set(idx int, value Value) {
	if idx < 0 || idx >= len(r.values) {
		return
	}

	r.values[idx] = value
}

// Len returns the number of values in the row.
func (r Row) Len() int {
	return len(r.values)
}

// Values returns a copy of the slice of values in the row.
func (r Row) Values() []Value {
	return r.values
}

// GetValueByColumn retrieves the value from the row corresponding to the specified column name in the given schema.
func (r Row) GetValueByColumn(schema Schema, columnName string) (Value, bool) {
	_, idx := schema.GetColumn(columnName)
	if idx == -1 {
		return Value{}, false
	}

	return r.Get(idx)
}
