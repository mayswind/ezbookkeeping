package datatable

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// WritableTransactionDataTable defines the structure of writable transaction data table
type WritableTransactionDataTable struct {
	allData          []map[TransactionDataTableColumn]string
	supportedColumns map[TransactionDataTableColumn]bool
	rowParser        TransactionDataRowParser
	addedColumns     map[TransactionDataTableColumn]bool
}

// WritableTransactionDataRow defines the structure of transaction data row of writable data table
type WritableTransactionDataRow struct {
	dataTable    *WritableTransactionDataTable
	rowData      map[TransactionDataTableColumn]string
	rowDataValid bool
}

// WritableTransactionDataRowIterator defines the structure of transaction data row iterator of writable data table
type WritableTransactionDataRowIterator struct {
	dataTable *WritableTransactionDataTable
	nextIndex int
}

// Add appends a new record to data table
func (t *WritableTransactionDataTable) Add(data map[TransactionDataTableColumn]string) {
	finalData := make(map[TransactionDataTableColumn]string, len(data))

	for column, value := range data {
		_, exists := t.supportedColumns[column]

		if exists {
			finalData[column] = value
		}
	}

	t.allData = append(t.allData, finalData)
}

// Get returns the record in the specified index
func (t *WritableTransactionDataTable) Get(index int) (*WritableTransactionDataRow, error) {
	if index >= len(t.allData) {
		return nil, nil
	}

	rowData := t.allData[index]
	rowDataValid := true

	if t.rowParser != nil {
		var err error
		rowData, rowDataValid, err = t.rowParser.Parse(rowData)

		if err != nil {
			return nil, err
		}
	}

	return &WritableTransactionDataRow{
		dataTable:    t,
		rowData:      rowData,
		rowDataValid: rowDataValid,
	}, nil
}

// HasColumn returns whether the data table has specified column
func (t *WritableTransactionDataTable) HasColumn(column TransactionDataTableColumn) bool {
	_, exists := t.supportedColumns[column]

	if exists {
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
func (t *WritableTransactionDataTable) TransactionRowCount() int {
	return len(t.allData)
}

// TransactionRowIterator returns the iterator of transaction data row
func (t *WritableTransactionDataTable) TransactionRowIterator() TransactionDataRowIterator {
	return &WritableTransactionDataRowIterator{
		dataTable: t,
		nextIndex: 0,
	}
}

// ColumnCount returns the total count of column in this data row
func (r *WritableTransactionDataRow) ColumnCount() int {
	if !r.rowDataValid {
		return 0
	}

	columnCount := 0

	for column := range r.rowData {
		if r.dataTable.supportedColumns[column] || r.dataTable.addedColumns[column] {
			columnCount++
		}
	}

	return columnCount
}

// IsValid returns whether this row is valid data for importing
func (r *WritableTransactionDataRow) IsValid() bool {
	return r.rowDataValid
}

// GetData returns the data in the specified column type
func (r *WritableTransactionDataRow) GetData(column TransactionDataTableColumn) string {
	if !r.rowDataValid {
		return ""
	}

	_, exists := r.dataTable.supportedColumns[column]

	if exists {
		return r.rowData[column]
	}

	if r.dataTable.addedColumns != nil {
		_, exists = r.dataTable.addedColumns[column]

		if exists {
			return r.rowData[column]
		}
	}

	return ""
}

// HasNext returns whether the iterator does not reach the end
func (t *WritableTransactionDataRowIterator) HasNext() bool {
	return t.nextIndex < len(t.dataTable.allData)
}

// Next returns the next transaction data row
func (t *WritableTransactionDataRowIterator) Next(ctx core.Context, user *models.User) (daraRow TransactionDataRow, err error) {
	if t.nextIndex >= len(t.dataTable.allData) {
		return nil, nil
	}

	rowData := t.dataTable.allData[t.nextIndex]
	rowDataValid := true

	if t.dataTable.rowParser != nil {
		rowData, rowDataValid, err = t.dataTable.rowParser.Parse(rowData)

		if err != nil {
			log.Errorf(ctx, "[writable_transaction_data_table.Next] cannot parse data row, because %s", err.Error())
			return nil, err
		}
	}

	t.nextIndex++

	return &WritableTransactionDataRow{
		dataTable:    t.dataTable,
		rowData:      rowData,
		rowDataValid: rowDataValid,
	}, nil
}

// CreateNewWritableTransactionDataTable returns a new writable transaction data table according to the specified columns
func CreateNewWritableTransactionDataTable(columns []TransactionDataTableColumn) *WritableTransactionDataTable {
	return CreateNewWritableTransactionDataTableWithRowParser(columns, nil)
}

// CreateNewWritableTransactionDataTableWithRowParser returns a new writable transaction data table according to the specified columns
func CreateNewWritableTransactionDataTableWithRowParser(columns []TransactionDataTableColumn, rowParser TransactionDataRowParser) *WritableTransactionDataTable {
	supportedColumns := make(map[TransactionDataTableColumn]bool, len(columns))

	for i := 0; i < len(columns); i++ {
		column := columns[i]
		supportedColumns[column] = true
	}

	var addedColumns map[TransactionDataTableColumn]bool

	if rowParser != nil {
		addedColumnsByParser := rowParser.GetAddedColumns()
		addedColumns = make(map[TransactionDataTableColumn]bool, len(addedColumnsByParser))

		for i := 0; i < len(addedColumnsByParser); i++ {
			addedColumns[addedColumnsByParser[i]] = true
		}
	}

	return &WritableTransactionDataTable{
		allData:          make([]map[TransactionDataTableColumn]string, 0),
		supportedColumns: supportedColumns,
		rowParser:        rowParser,
		addedColumns:     addedColumns,
	}
}
