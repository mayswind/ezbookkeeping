package datatable

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// ImportedTransactionDataTable defines the structure of imported transaction data table
type ImportedTransactionDataTable struct {
	innerDataTable    ImportedDataTable
	dataColumnMapping map[TransactionDataTableColumn]string
	dataColumnIndexes map[TransactionDataTableColumn]int
	rowParser         TransactionDataRowParser
	addedColumns      map[TransactionDataTableColumn]bool
}

// ImportedTransactionDataRow defines the structure of imported transaction data row
type ImportedTransactionDataRow struct {
	transactionDataTable *ImportedTransactionDataTable
	rowData              map[TransactionDataTableColumn]string
	rowDataValid         bool
}

// ImportedTransactionDataRowIterator defines the structure of imported transaction data row iterator
type ImportedTransactionDataRowIterator struct {
	transactionDataTable *ImportedTransactionDataTable
	innerIterator        ImportedDataRowIterator
}

// HasColumn returns whether the data table has specified column
func (t *ImportedTransactionDataTable) HasColumn(column TransactionDataTableColumn) bool {
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
func (t *ImportedTransactionDataTable) TransactionRowCount() int {
	return t.innerDataTable.DataRowCount()
}

// TransactionRowIterator returns the iterator of transaction data row
func (t *ImportedTransactionDataTable) TransactionRowIterator() TransactionDataRowIterator {
	return &ImportedTransactionDataRowIterator{
		transactionDataTable: t,
		innerIterator:        t.innerDataTable.DataRowIterator(),
	}
}

// IsValid returns whether this row is valid data for importing
func (r *ImportedTransactionDataRow) IsValid() bool {
	return r.rowDataValid
}

// GetData returns the data in the specified column type
func (r *ImportedTransactionDataRow) GetData(column TransactionDataTableColumn) string {
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
func (t *ImportedTransactionDataRowIterator) HasNext() bool {
	return t.innerIterator.HasNext()
}

// Next returns the next transaction data row
func (t *ImportedTransactionDataRowIterator) Next(ctx core.Context, user *models.User) (daraRow TransactionDataRow, err error) {
	importedRow := t.innerIterator.Next()

	if importedRow == nil {
		return nil, nil
	}

	if importedRow.ColumnCount() == 1 && importedRow.GetData(0) == "" {
		return &ImportedTransactionDataRow{
			transactionDataTable: t.transactionDataTable,
			rowData:              nil,
			rowDataValid:         false,
		}, nil
	}

	if importedRow.ColumnCount() < len(t.transactionDataTable.dataColumnIndexes) {
		log.Errorf(ctx, "[imported_transaction_data_table.Next] cannot parse data row, because may missing some columns (column count %d in data row is less than header column count %d)", importedRow.ColumnCount(), len(t.transactionDataTable.dataColumnIndexes))
		return nil, errs.ErrFewerFieldsInDataRowThanInHeaderRow
	}

	rowData := make(map[TransactionDataTableColumn]string, len(t.transactionDataTable.dataColumnIndexes))
	rowDataValid := true

	for column, columnIndex := range t.transactionDataTable.dataColumnIndexes {
		if columnIndex < 0 || columnIndex >= importedRow.ColumnCount() {
			continue
		}

		value := importedRow.GetData(columnIndex)
		rowData[column] = value
	}

	if t.transactionDataTable.rowParser != nil {
		rowData, rowDataValid, err = t.transactionDataTable.rowParser.Parse(rowData)

		if err != nil {
			log.Errorf(ctx, "[imported_transaction_data_table.Next] cannot parse data row, because %s", err.Error())
			return nil, err
		}
	}

	return &ImportedTransactionDataRow{
		transactionDataTable: t.transactionDataTable,
		rowData:              rowData,
		rowDataValid:         rowDataValid,
	}, nil
}

// CreateImportedTransactionDataTable returns transaction data table from imported data table
func CreateImportedTransactionDataTable(dataTable ImportedDataTable, dataColumnMapping map[TransactionDataTableColumn]string) *ImportedTransactionDataTable {
	return CreateImportedTransactionDataTableWithRowParser(dataTable, dataColumnMapping, nil)
}

// CreateImportedTransactionDataTableWithRowParser returns transaction data table from imported data table
func CreateImportedTransactionDataTableWithRowParser(dataTable ImportedDataTable, dataColumnMapping map[TransactionDataTableColumn]string, rowParser TransactionDataRowParser) *ImportedTransactionDataTable {
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

	return &ImportedTransactionDataTable{
		innerDataTable:    dataTable,
		dataColumnMapping: dataColumnMapping,
		dataColumnIndexes: dataColumnIndexes,
		rowParser:         rowParser,
		addedColumns:      addedColumns,
	}
}
