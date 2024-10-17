package datatable

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// CommonTransactionDataTable defines the structure of common transaction data table
type CommonTransactionDataTable struct {
	innerDataTable       CommonDataTable
	supportedDataColumns map[TransactionDataTableColumn]bool
	rowParser            CommonTransactionDataRowParser
}

// CommonTransactionDataRow defines the structure of common transaction data row
type CommonTransactionDataRow struct {
	transactionDataTable *CommonTransactionDataTable
	rowData              map[TransactionDataTableColumn]string
	rowDataValid         bool
}

// CommonTransactionDataRowIterator defines the structure of common transaction data row iterator
type CommonTransactionDataRowIterator struct {
	transactionDataTable *CommonTransactionDataTable
	innerIterator        CommonDataRowIterator
}

// CommonTransactionDataRowParser defines the structure of common transaction data row parser
type CommonTransactionDataRowParser interface {
	// Parse returns the converted transaction data row
	Parse(ctx core.Context, user *models.User, dataTable *CommonTransactionDataTable, dataRow CommonDataRow, rowId string) (rowData map[TransactionDataTableColumn]string, rowDataValid bool, err error)
}

// HasColumn returns whether the data table has specified column
func (t *CommonTransactionDataTable) HasColumn(column TransactionDataTableColumn) bool {
	_, exists := t.supportedDataColumns[column]
	return exists
}

// HasOriginalColumn returns whether the original data table has specified column name
func (t *CommonTransactionDataTable) HasOriginalColumn(columnName string) bool {
	return columnName != "" && t.innerDataTable.HasColumn(columnName)
}

// TransactionRowCount returns the total count of transaction data row
func (t *CommonTransactionDataTable) TransactionRowCount() int {
	return t.innerDataTable.DataRowCount()
}

// TransactionRowIterator returns the iterator of transaction data row
func (t *CommonTransactionDataTable) TransactionRowIterator() TransactionDataRowIterator {
	return &CommonTransactionDataRowIterator{
		transactionDataTable: t,
		innerIterator:        t.innerDataTable.DataRowIterator(),
	}
}

// IsValid returns whether this row is valid data for importing
func (r *CommonTransactionDataRow) IsValid() bool {
	return r.rowDataValid
}

// GetData returns the data in the specified column type
func (r *CommonTransactionDataRow) GetData(column TransactionDataTableColumn) string {
	if !r.rowDataValid {
		return ""
	}

	_, exists := r.transactionDataTable.supportedDataColumns[column]

	if !exists {
		return ""
	}

	return r.rowData[column]
}

// HasNext returns whether the iterator does not reach the end
func (t *CommonTransactionDataRowIterator) HasNext() bool {
	return t.innerIterator.HasNext()
}

// Next returns the next transaction data row
func (t *CommonTransactionDataRowIterator) Next(ctx core.Context, user *models.User) (daraRow TransactionDataRow, err error) {
	commonRow := t.innerIterator.Next()

	if commonRow == nil {
		return nil, nil
	}

	rowId := t.innerIterator.CurrentRowId()
	rowData, rowDataValid, err := t.transactionDataTable.rowParser.Parse(ctx, user, t.transactionDataTable, commonRow, rowId)

	if err != nil {
		log.Errorf(ctx, "[common_transaction_data_table.Next] cannot parse data row, because %s", err.Error())
		return nil, err
	}

	return &CommonTransactionDataRow{
		transactionDataTable: t.transactionDataTable,
		rowData:              rowData,
		rowDataValid:         rowDataValid,
	}, nil
}

// CreateNewCommonTransactionDataTable returns transaction data table from Common data table
func CreateNewCommonTransactionDataTable(dataTable CommonDataTable, supportedDataColumns map[TransactionDataTableColumn]bool, rowParser CommonTransactionDataRowParser) *CommonTransactionDataTable {
	return &CommonTransactionDataTable{
		innerDataTable:       dataTable,
		supportedDataColumns: supportedDataColumns,
		rowParser:            rowParser,
	}
}
