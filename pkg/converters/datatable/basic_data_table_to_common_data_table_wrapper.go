package datatable

// basicDataTableToCommonDataTableWrapper defines the structure of basic data table to common data table wrapper
type basicDataTableToCommonDataTableWrapper struct {
	innerDataTable    BasicDataTable
	dataColumnIndexes map[string]int
}

// basicDataTableToCommonDataTableWrapperRow defines the data row structure of basic data table to common data table wrapper
type basicDataTableToCommonDataTableWrapperRow struct {
	rowData map[string]string
}

// basicDataTableToCommonDataTableWrapperRowIterator defines the data row iterator structure of basic data table to common data table wrapper
type basicDataTableToCommonDataTableWrapperRowIterator struct {
	commonDataTable *basicDataTableToCommonDataTableWrapper
	innerIterator   BasicDataTableRowIterator
}

// HeaderColumnCount returns the total count of column in header row
func (t *basicDataTableToCommonDataTableWrapper) HeaderColumnCount() int {
	return len(t.innerDataTable.HeaderColumnNames())
}

// HasColumn returns whether the data table has specified column name
func (t *basicDataTableToCommonDataTableWrapper) HasColumn(columnName string) bool {
	index, exists := t.dataColumnIndexes[columnName]
	return exists && index >= 0
}

// DataRowCount returns the total count of common data row
func (t *basicDataTableToCommonDataTableWrapper) DataRowCount() int {
	return t.innerDataTable.DataRowCount()
}

// DataRowIterator returns the iterator of common data row
func (t *basicDataTableToCommonDataTableWrapper) DataRowIterator() CommonDataTableRowIterator {
	return &basicDataTableToCommonDataTableWrapperRowIterator{
		commonDataTable: t,
		innerIterator:   t.innerDataTable.DataRowIterator(),
	}
}

// HasData returns whether the common data row has specified column data
func (r *basicDataTableToCommonDataTableWrapperRow) HasData(columnName string) bool {
	_, exists := r.rowData[columnName]
	return exists
}

// ColumnCount returns the total count of column in this data row
func (r *basicDataTableToCommonDataTableWrapperRow) ColumnCount() int {
	return len(r.rowData)
}

// GetData returns the data in the specified column name
func (r *basicDataTableToCommonDataTableWrapperRow) GetData(columnName string) string {
	return r.rowData[columnName]
}

// HasNext returns whether the iterator does not reach the end
func (t *basicDataTableToCommonDataTableWrapperRowIterator) HasNext() bool {
	return t.innerIterator.HasNext()
}

// CurrentRowId returns current row id
func (t *basicDataTableToCommonDataTableWrapperRowIterator) CurrentRowId() string {
	return t.innerIterator.CurrentRowId()
}

// Next returns the next common data row
func (t *basicDataTableToCommonDataTableWrapperRowIterator) Next() CommonDataTableRow {
	basicDataRow := t.innerIterator.Next()

	if basicDataRow == nil {
		return nil
	}

	rowData := make(map[string]string, len(t.commonDataTable.dataColumnIndexes))

	for column, columnIndex := range t.commonDataTable.dataColumnIndexes {
		if columnIndex < 0 || columnIndex >= basicDataRow.ColumnCount() {
			continue
		}

		value := basicDataRow.GetData(columnIndex)
		rowData[column] = value
	}

	return &basicDataTableToCommonDataTableWrapperRow{
		rowData: rowData,
	}
}

// CreateNewCommonDataTableFromBasicDataTable returns common data table from basic data table
func CreateNewCommonDataTableFromBasicDataTable(dataTable BasicDataTable) CommonDataTable {
	headerLineItems := dataTable.HeaderColumnNames()
	dataColumnIndexes := make(map[string]int, len(headerLineItems))

	for i := 0; i < len(headerLineItems); i++ {
		dataColumnIndexes[headerLineItems[i]] = i
	}

	return &basicDataTableToCommonDataTableWrapper{
		innerDataTable:    dataTable,
		dataColumnIndexes: dataColumnIndexes,
	}
}
