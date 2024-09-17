package converters

import (
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// WritableDataTable defines the structure of writable data table
type WritableDataTable struct {
	allData []map[DataTableColumn]string
	columns []DataTableColumn
}

// WritableDataRow defines the structure of data row of writable data table
type WritableDataRow struct {
	dataTable *WritableDataTable
	rowData   map[DataTableColumn]string
}

// WritableDataRowIterator defines the structure of data row iterator of writable data table
type WritableDataRowIterator struct {
	dataTable *WritableDataTable
	nextIndex int
}

// Add appends a new record to data table
func (t *WritableDataTable) Add(data map[DataTableColumn]string) {
	finalData := make(map[DataTableColumn]string, len(data))

	for i := 0; i < len(t.columns); i++ {
		column := t.columns[i]

		if value, exists := data[column]; exists {
			finalData[column] = value
		}
	}

	t.allData = append(t.allData, finalData)
}

// Get returns the record in the specified index
func (t *WritableDataTable) Get(index int) ImportedDataRow {
	if index >= len(t.allData) {
		return nil
	}

	rowData := t.allData[index]

	return &WritableDataRow{
		dataTable: t,
		rowData:   rowData,
	}
}

// DataRowCount returns the total count of data row
func (t *WritableDataTable) DataRowCount() int {
	return len(t.allData)
}

// GetDataColumnMapping returns data column map for data importer
func (t *WritableDataTable) GetDataColumnMapping() map[DataTableColumn]string {
	dataColumnMapping := make(map[DataTableColumn]string, len(t.columns))

	for i := 0; i < len(t.columns); i++ {
		column := t.columns[i]
		dataColumnMapping[column] = utils.IntToString(int(column))
	}

	return dataColumnMapping
}

// HeaderLineColumnNames returns the header column name list
func (t *WritableDataTable) HeaderLineColumnNames() []string {
	columnIndexes := make([]string, len(t.columns))

	for i := 0; i < len(t.columns); i++ {
		columnIndexes[i] = utils.IntToString(int(t.columns[i]))
	}

	return columnIndexes
}

// DataRowIterator returns the iterator of data row
func (t *WritableDataTable) DataRowIterator() ImportedDataRowIterator {
	return &WritableDataRowIterator{
		dataTable: t,
		nextIndex: 0,
	}
}

// ColumnCount returns the total count of column in this data row
func (r *WritableDataRow) ColumnCount() int {
	return len(r.rowData)
}

// GetData returns the data in the specified column index
func (r *WritableDataRow) GetData(columnIndex int) string {
	if columnIndex >= len(r.dataTable.columns) {
		return ""
	}

	dataColumn := r.dataTable.columns[columnIndex]

	return r.rowData[dataColumn]
}

// GetTime returns the time in the specified column index
func (r *WritableDataRow) GetTime(columnIndex int, timezoneOffset int16) (time.Time, error) {
	return utils.ParseFromLongDateTime(r.GetData(columnIndex), timezoneOffset)
}

// GetTimezoneOffset returns the time zone offset in the specified column index
func (r *WritableDataRow) GetTimezoneOffset(columnIndex int) (*time.Location, error) {
	return utils.ParseFromTimezoneOffset(r.GetData(columnIndex))
}

// HasNext returns whether the iterator does not reach the end
func (t *WritableDataRowIterator) HasNext() bool {
	return t.nextIndex < len(t.dataTable.allData)
}

// Next returns the next imported data row
func (t *WritableDataRowIterator) Next() ImportedDataRow {
	if t.nextIndex >= len(t.dataTable.allData) {
		return nil
	}

	rowData := t.dataTable.allData[t.nextIndex]

	t.nextIndex++

	return &WritableDataRow{
		dataTable: t.dataTable,
		rowData:   rowData,
	}
}

func createNewWritableDataTable(columns []DataTableColumn) (*WritableDataTable, error) {
	return &WritableDataTable{
		allData: make([]map[DataTableColumn]string, 0),
		columns: columns,
	}, nil
}
