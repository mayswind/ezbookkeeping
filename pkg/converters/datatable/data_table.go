package datatable

import (
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// ImportedDataTable defines the structure of imported data table
type ImportedDataTable interface {
	// DataRowCount returns the total count of data row
	DataRowCount() int

	// HeaderLineColumnNames returns the header column name list
	HeaderLineColumnNames() []string

	// DataRowIterator returns the iterator of data row
	DataRowIterator() ImportedDataRowIterator
}

// ImportedDataRow defines the structure of imported data row
type ImportedDataRow interface {
	// IsValid returns whether this row contains valid data for importing
	IsValid() bool

	// ColumnCount returns the total count of column in this data row
	ColumnCount() int

	// GetData returns the data in the specified column index
	GetData(columnIndex int) string

	// GetTime returns the time in the specified column index
	GetTime(columnIndex int, timezoneOffset int16) (time.Time, error)

	// GetTimezoneOffset returns the time zone offset in the specified column index
	GetTimezoneOffset(columnIndex int) (*time.Location, error)
}

// ImportedDataRowIterator defines the structure of imported data row iterator
type ImportedDataRowIterator interface {
	// HasNext returns whether the iterator does not reach the end
	HasNext() bool

	// Next returns the next imported data row
	Next(ctx core.Context, user *models.User) ImportedDataRow
}

// DataTableBuilder defines the structure of data table builder
type DataTableBuilder interface {
	// AppendTransaction appends the specified transaction to data builder
	AppendTransaction(data map[DataTableColumn]string)

	// ReplaceDelimiters returns the text after removing the delimiters
	ReplaceDelimiters(text string) string
}
