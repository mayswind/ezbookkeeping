package datatable

// CommonDataTable defines the structure of common data table
type CommonDataTable interface {
	// HeaderColumnCount returns the total count of column in header row
	HeaderColumnCount() int

	// HasColumn returns whether the common data table has specified column name
	HasColumn(columnName string) bool

	// DataRowCount returns the total count of common data row
	DataRowCount() int

	// DataRowIterator returns the iterator of common data row
	DataRowIterator() CommonDataRowIterator
}

// CommonDataRow defines the structure of common data row
type CommonDataRow interface {
	// ColumnCount returns the total count of column in this data row
	ColumnCount() int

	// HasData returns whether the common data row has specified column data
	HasData(columnName string) bool

	// GetData returns the data in the specified column name
	GetData(columnName string) string
}

// CommonDataRowIterator defines the structure of common data row iterator
type CommonDataRowIterator interface {
	// HasNext returns whether the iterator does not reach the end
	HasNext() bool

	// CurrentRowId returns current row id
	CurrentRowId() string

	// Next returns the next common data row
	Next() CommonDataRow
}
