package datatable

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// TransactionDataTable defines the structure of transaction data table
type TransactionDataTable interface {
	// HasColumn returns whether the transaction data table has specified column
	HasColumn(column TransactionDataTableColumn) bool

	// TransactionRowCount returns the total count of transaction data row
	TransactionRowCount() int

	// TransactionRowIterator returns the iterator of transaction data row
	TransactionRowIterator() TransactionDataRowIterator
}

// TransactionDataRow defines the structure of transaction data row
type TransactionDataRow interface {
	// IsValid returns whether this row is valid data for importing
	IsValid() bool

	// GetData returns the data in the specified column type
	GetData(column TransactionDataTableColumn) string
}

// TransactionDataRowIterator defines the structure of transaction data row iterator
type TransactionDataRowIterator interface {
	// HasNext returns whether the iterator does not reach the end
	HasNext() bool

	// Next returns the next transaction data row
	Next(ctx core.Context, user *models.User) (daraRow TransactionDataRow, err error)
}

// TransactionDataRowParser defines the structure of transaction data row parser
type TransactionDataRowParser interface {
	// GetAddedColumns returns the added columns after converting the data row
	GetAddedColumns() []TransactionDataTableColumn

	// Parse returns the converted transaction data row
	Parse(data map[TransactionDataTableColumn]string) (rowData map[TransactionDataTableColumn]string, rowDataValid bool, err error)
}

// TransactionDataTableBuilder defines the structure of data table builder
type TransactionDataTableBuilder interface {
	// AppendTransaction appends the specified transaction to data builder
	AppendTransaction(data map[TransactionDataTableColumn]string)

	// ReplaceDelimiters returns the text after removing the delimiters
	ReplaceDelimiters(text string) string
}

// TransactionDataTableColumn represents the data column type of data table
type TransactionDataTableColumn byte

// Transaction data table columns
const (
	TRANSACTION_DATA_TABLE_TRANSACTION_TIME         TransactionDataTableColumn = 1
	TRANSACTION_DATA_TABLE_TRANSACTION_TIMEZONE     TransactionDataTableColumn = 2
	TRANSACTION_DATA_TABLE_TRANSACTION_TYPE         TransactionDataTableColumn = 3
	TRANSACTION_DATA_TABLE_CATEGORY                 TransactionDataTableColumn = 4
	TRANSACTION_DATA_TABLE_SUB_CATEGORY             TransactionDataTableColumn = 5
	TRANSACTION_DATA_TABLE_ACCOUNT_NAME             TransactionDataTableColumn = 6
	TRANSACTION_DATA_TABLE_ACCOUNT_CURRENCY         TransactionDataTableColumn = 7
	TRANSACTION_DATA_TABLE_AMOUNT                   TransactionDataTableColumn = 8
	TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_NAME     TransactionDataTableColumn = 9
	TRANSACTION_DATA_TABLE_RELATED_ACCOUNT_CURRENCY TransactionDataTableColumn = 10
	TRANSACTION_DATA_TABLE_RELATED_AMOUNT           TransactionDataTableColumn = 11
	TRANSACTION_DATA_TABLE_GEOGRAPHIC_LOCATION      TransactionDataTableColumn = 12
	TRANSACTION_DATA_TABLE_TAGS                     TransactionDataTableColumn = 13
	TRANSACTION_DATA_TABLE_DESCRIPTION              TransactionDataTableColumn = 14
)
