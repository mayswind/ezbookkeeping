package datatable

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// basicDataTableToTransactionDataTableWrapper defines the structure of basic data table to transaction data table wrapper
type basicDataTableToTransactionDataTableWrapper struct {
	innerDataTable    BasicDataTable
	dataColumnMapping map[TransactionDataTableColumn]string
	dataColumnIndexes map[TransactionDataTableColumn]int
	rowParser         TransactionDataRowParser
	addedColumns      map[TransactionDataTableColumn]bool
}

// basicDataTableToTransactionDataTableWrapperRow defines the data row structure of basic data table to transaction data table wrapper
type basicDataTableToTransactionDataTableWrapperRow struct {
	transactionDataTable *basicDataTableToTransactionDataTableWrapper
	rowData              map[TransactionDataTableColumn]string
	rowDataValid         bool
}

// basicDataTableToTransactionDataTableWrapperRowIterator defines the data row iterator structure of basic data table to transaction data table wrapper
type basicDataTableToTransactionDataTableWrapperRowIterator struct {
	transactionDataTable *basicDataTableToTransactionDataTableWrapper
	innerIterator        BasicDataTableRowIterator
}

// HasColumn returns whether the data table has specified column
func (t *basicDataTableToTransactionDataTableWrapper) HasColumn(column TransactionDataTableColumn) bool {
	index, exists := t.dataColumnIndexes[column]

	if exists && index >= 0 {
		return exists
	}

	if t.addedColumns != nil {
		_, exists = t.addedColumns[column]

		if exists {
			return exists
		}
	}

	return false
}

// TransactionRowCount returns the total count of transaction data row
func (t *basicDataTableToTransactionDataTableWrapper) TransactionRowCount() int {
	return t.innerDataTable.DataRowCount()
}

// TransactionRowIterator returns the iterator of transaction data row
func (t *basicDataTableToTransactionDataTableWrapper) TransactionRowIterator() TransactionDataRowIterator {
	return &basicDataTableToTransactionDataTableWrapperRowIterator{
		transactionDataTable: t,
		innerIterator:        t.innerDataTable.DataRowIterator(),
	}
}

// IsValid returns whether this row is valid data for importing
func (r *basicDataTableToTransactionDataTableWrapperRow) IsValid() bool {
	return r.rowDataValid
}

// GetData returns the data in the specified column type
func (r *basicDataTableToTransactionDataTableWrapperRow) GetData(column TransactionDataTableColumn) string {
	if !r.rowDataValid {
		return ""
	}

	_, exists := r.transactionDataTable.dataColumnIndexes[column]

	if exists {
		return r.rowData[column]
	}

	if r.transactionDataTable.addedColumns != nil {
		_, exists = r.transactionDataTable.addedColumns[column]

		if exists {
			return r.rowData[column]
		}
	}

	return ""
}

// HasNext returns whether the iterator does not reach the end
func (t *basicDataTableToTransactionDataTableWrapperRowIterator) HasNext() bool {
	return t.innerIterator.HasNext()
}

// Next returns the next transaction data row
func (t *basicDataTableToTransactionDataTableWrapperRowIterator) Next(ctx core.Context, user *models.User) (daraRow TransactionDataRow, err error) {
	basicDataRow := t.innerIterator.Next()

	if basicDataRow == nil {
		return nil, nil
	}

	if basicDataRow.ColumnCount() == 1 && basicDataRow.GetData(0) == "" {
		return &basicDataTableToTransactionDataTableWrapperRow{
			transactionDataTable: t.transactionDataTable,
			rowData:              nil,
			rowDataValid:         false,
		}, nil
	}

	if basicDataRow.ColumnCount() < len(t.transactionDataTable.dataColumnIndexes) {
		log.Errorf(ctx, "[basic_data_table_to_transaction_data_table_wrapper.Next] cannot parse data row, because may missing some columns (column count %d in data row is less than header column count %d)", basicDataRow.ColumnCount(), len(t.transactionDataTable.dataColumnIndexes))
		return nil, errs.ErrFewerFieldsInDataRowThanInHeaderRow
	}

	rowData := make(map[TransactionDataTableColumn]string, len(t.transactionDataTable.dataColumnIndexes))
	rowDataValid := true

	for column, columnIndex := range t.transactionDataTable.dataColumnIndexes {
		if columnIndex < 0 || columnIndex >= basicDataRow.ColumnCount() {
			continue
		}

		value := basicDataRow.GetData(columnIndex)
		rowData[column] = value
	}

	if t.transactionDataTable.rowParser != nil {
		rowData, rowDataValid, err = t.transactionDataTable.rowParser.Parse(rowData)

		if err != nil {
			log.Errorf(ctx, "[basic_data_table_to_transaction_data_table_wrapper.Next] cannot parse data row, because %s", err.Error())
			return nil, err
		}
	}

	return &basicDataTableToTransactionDataTableWrapperRow{
		transactionDataTable: t.transactionDataTable,
		rowData:              rowData,
		rowDataValid:         rowDataValid,
	}, nil
}

// CreateNewTransactionDataTableFromBasicDataTable returns transaction data table from basic data table
func CreateNewTransactionDataTableFromBasicDataTable(dataTable BasicDataTable, dataColumnMapping map[TransactionDataTableColumn]string) TransactionDataTable {
	return CreateNewTransactionDataTableFromBasicDataTableWithRowParser(dataTable, dataColumnMapping, nil)
}

// CreateNewTransactionDataTableFromBasicDataTableWithRowParser returns transaction data table from basic data table
func CreateNewTransactionDataTableFromBasicDataTableWithRowParser(dataTable BasicDataTable, dataColumnMapping map[TransactionDataTableColumn]string, rowParser TransactionDataRowParser) TransactionDataTable {
	headerLineItems := dataTable.HeaderColumnNames()
	headerItemMap := make(map[string]int, len(headerLineItems))

	for i := 0; i < len(headerLineItems); i++ {
		headerItemMap[headerLineItems[i]] = i
	}

	dataColumnIndexes := make(map[TransactionDataTableColumn]int, len(headerLineItems))

	for column, columnName := range dataColumnMapping {
		columnIndex, exists := headerItemMap[columnName]

		if exists {
			dataColumnIndexes[column] = columnIndex
		}
	}

	var addedColumns map[TransactionDataTableColumn]bool

	if rowParser != nil {
		addedColumnsByParser := rowParser.GetAddedColumns()
		addedColumns = make(map[TransactionDataTableColumn]bool, len(addedColumnsByParser))

		for i := 0; i < len(addedColumnsByParser); i++ {
			addedColumns[addedColumnsByParser[i]] = true
		}
	}

	return &basicDataTableToTransactionDataTableWrapper{
		innerDataTable:    dataTable,
		dataColumnMapping: dataColumnMapping,
		dataColumnIndexes: dataColumnIndexes,
		rowParser:         rowParser,
		addedColumns:      addedColumns,
	}
}
