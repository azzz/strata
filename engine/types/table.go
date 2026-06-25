package types

// Table represents a table in the dataset, consisting of a schema, partitioning columns, and sorting columns.
type Table struct {
	Schema      Schema   // The schema of the table, defining the structure and types of the data
	PartitionBy []Column // The columns used for partitioning the data in the table
	SortBy      []Column // The columns used for sorting the data in the table
}
