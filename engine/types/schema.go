package types

// Schema represents the schema of a dataset, consisting of a slice of Columns.
type Schema struct {
	columns []Column
}

func NewSchema(columns ...Column) Schema {
	return Schema{
		columns: columns,
	}
}

// GetColumn retrieves a Column by its name from the Schema. It returns the Column and its index, or an empty Column and -1 if not found.
func (s Schema) GetColumn(name string) (Column, ColumnIndex) {
	for i, field := range s.columns {
		if field.Name == name {
			return field, ColumnIndex(i)
		}
	}

	return Column{}, -1
}

// GetColumnByIndex retrieves a Column by its index from the Schema. It returns the Column and a boolean indicating if the index is valid.
func (s Schema) GetColumnByIndex(idx ColumnIndex) (Column, bool) {
	if idx < 0 || int(idx) >= len(s.columns) {
		return Column{}, false
	}

	return s.columns[idx], true
}

// Columns returns a copy of the slice of Columns in the Schema.
func (s Schema) Columns() []Column {
	return s.columns
}

// Len returns the number of columns in the Schema.
func (s Schema) Len() int {
	return len(s.columns)
}

// Column represents a single column in a schema, consisting of a name and a data type (Kind).
type Column struct {
	Name string // The name of the column
	Type Kind   // The data type of the column, represented by the Kind type
}

// ColumnIndex is an index of a column is Schema
type ColumnIndex int
