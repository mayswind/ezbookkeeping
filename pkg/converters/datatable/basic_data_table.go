package datatable

// BasicDataTable defines the structure of basic data table
type BasicDataTable interface {
	// DataRowCount returns the total count of data row
	DataRowCount() int

	// HeaderColumnNames returns the header column name list
	HeaderColumnNames() []string

	// DataRowIterator returns the iterator of data row
	DataRowIterator() BasicDataTableRowIterator
}

// BasicDataTableRow defines the structure of basic data row
type BasicDataTableRow interface {
	// ColumnCount returns the total count of column in this data row
	ColumnCount() int

	// GetData returns the data in the specified column index
	GetData(columnIndex int) string
}

// BasicDataTableRowIterator defines the structure of basic data row iterator
type BasicDataTableRowIterator interface {
	// HasNext returns whether the iterator does not reach the end
	HasNext() bool

	// CurrentRowId returns current row id
	CurrentRowId() string

	// Next returns the next basic data row
	Next() BasicDataTableRow
}
