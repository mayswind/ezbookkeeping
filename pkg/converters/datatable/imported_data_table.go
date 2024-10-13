package datatable

// ImportedDataTable defines the structure of imported data table
type ImportedDataTable interface {
	// DataRowCount returns the total count of data row
	DataRowCount() int

	// HeaderColumnNames returns the header column name list
	HeaderColumnNames() []string

	// DataRowIterator returns the iterator of data row
	DataRowIterator() ImportedDataRowIterator
}

// ImportedDataRow defines the structure of imported data row
type ImportedDataRow interface {
	// ColumnCount returns the total count of column in this data row
	ColumnCount() int

	// GetData returns the data in the specified column index
	GetData(columnIndex int) string
}

// ImportedDataRowIterator defines the structure of imported data row iterator
type ImportedDataRowIterator interface {
	// HasNext returns whether the iterator does not reach the end
	HasNext() bool

	// Next returns the next imported data row
	Next() ImportedDataRow
}
