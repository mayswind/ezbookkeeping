package datatable

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// CommonTransactionDataRowParser defines the structure of common transaction data row parser
type CommonTransactionDataRowParser interface {
	// Parse returns the converted transaction data row
	Parse(ctx core.Context, user *models.User, dataRow CommonDataTableRow, rowId string) (rowData map[TransactionDataTableColumn]string, rowDataValid bool, err error)
}

// commonDataTableToTransactionDataTableWrapper defines the structure of common data table to transaction data table wrapper
type commonDataTableToTransactionDataTableWrapper struct {
	innerDataTable       CommonDataTable
	supportedDataColumns map[TransactionDataTableColumn]bool
	rowParser            CommonTransactionDataRowParser
}

// commonDataTableToTransactionDataTableWrapperRow defines the data row structure of common data table to transaction data table wrapper
type commonDataTableToTransactionDataTableWrapperRow struct {
	transactionDataTable *commonDataTableToTransactionDataTableWrapper
	rowData              map[TransactionDataTableColumn]string
	rowDataValid         bool
}

// commonDataTableToTransactionDataTableWrapperRowIterator defines the data row iterator structure of common data table to transaction data table wrapper
type commonDataTableToTransactionDataTableWrapperRowIterator struct {
	transactionDataTable *commonDataTableToTransactionDataTableWrapper
	innerIterator        CommonDataTableRowIterator
}

// HasColumn returns whether the data table has specified column
func (t *commonDataTableToTransactionDataTableWrapper) HasColumn(column TransactionDataTableColumn) bool {
	_, exists := t.supportedDataColumns[column]
	return exists
}

// TransactionRowCount returns the total count of transaction data row
func (t *commonDataTableToTransactionDataTableWrapper) TransactionRowCount() int {
	return t.innerDataTable.DataRowCount()
}

// TransactionRowIterator returns the iterator of transaction data row
func (t *commonDataTableToTransactionDataTableWrapper) TransactionRowIterator() TransactionDataRowIterator {
	return &commonDataTableToTransactionDataTableWrapperRowIterator{
		transactionDataTable: t,
		innerIterator:        t.innerDataTable.DataRowIterator(),
	}
}

// IsValid returns whether this row is valid data for importing
func (r *commonDataTableToTransactionDataTableWrapperRow) IsValid() bool {
	return r.rowDataValid
}

// GetData returns the data in the specified column type
func (r *commonDataTableToTransactionDataTableWrapperRow) GetData(column TransactionDataTableColumn) string {
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
func (t *commonDataTableToTransactionDataTableWrapperRowIterator) HasNext() bool {
	return t.innerIterator.HasNext()
}

// Next returns the next transaction data row
func (t *commonDataTableToTransactionDataTableWrapperRowIterator) Next(ctx core.Context, user *models.User) (daraRow TransactionDataRow, err error) {
	commonDataRow := t.innerIterator.Next()

	if commonDataRow == nil {
		return nil, nil
	}

	rowId := t.innerIterator.CurrentRowId()
	rowData, rowDataValid, err := t.transactionDataTable.rowParser.Parse(ctx, user, commonDataRow, rowId)

	if err != nil {
		log.Errorf(ctx, "[common_data_table_to_transaction_data_table_wrapper.Next] cannot parse data row, because %s", err.Error())
		return nil, err
	}

	return &commonDataTableToTransactionDataTableWrapperRow{
		transactionDataTable: t.transactionDataTable,
		rowData:              rowData,
		rowDataValid:         rowDataValid,
	}, nil
}

// CreateNewTransactionDataTableFromCommonDataTable returns transaction data table from Common data table
func CreateNewTransactionDataTableFromCommonDataTable(dataTable CommonDataTable, supportedDataColumns map[TransactionDataTableColumn]bool, rowParser CommonTransactionDataRowParser) TransactionDataTable {
	return &commonDataTableToTransactionDataTableWrapper{
		innerDataTable:       dataTable,
		supportedDataColumns: supportedDataColumns,
		rowParser:            rowParser,
	}
}
