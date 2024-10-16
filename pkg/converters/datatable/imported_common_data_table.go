package datatable

// ImportedCommonDataTable defines the structure of imported common data table
type ImportedCommonDataTable struct {
	innerDataTable    ImportedDataTable
	dataColumnIndexes map[string]int
}

// ImportedCommonDataRow defines the structure of imported common data row
type ImportedCommonDataRow struct {
	rowData map[string]string
}

// ImportedCommonDataRowIterator defines the structure of imported common data row iterator
type ImportedCommonDataRowIterator struct {
	commonDataTable *ImportedCommonDataTable
	innerIterator   ImportedDataRowIterator
}

// HeaderColumnCount returns the total count of column in header row
func (t *ImportedCommonDataTable) HeaderColumnCount() int {
	return len(t.innerDataTable.HeaderColumnNames())
}

// HasColumn returns whether the data table has specified column name
func (t *ImportedCommonDataTable) HasColumn(columnName string) bool {
	index, exists := t.dataColumnIndexes[columnName]
	return exists && index >= 0
}

// DataRowCount returns the total count of common data row
func (t *ImportedCommonDataTable) DataRowCount() int {
	return t.innerDataTable.DataRowCount()
}

// DataRowIterator returns the iterator of common data row
func (t *ImportedCommonDataTable) DataRowIterator() CommonDataRowIterator {
	return &ImportedCommonDataRowIterator{
		commonDataTable: t,
		innerIterator:   t.innerDataTable.DataRowIterator(),
	}
}

// HasData returns whether the common data row has specified column data
func (r *ImportedCommonDataRow) HasData(columnName string) bool {
	_, exists := r.rowData[columnName]
	return exists
}

// ColumnCount returns the total count of column in this data row
func (r *ImportedCommonDataRow) ColumnCount() int {
	return len(r.rowData)
}

// GetData returns the data in the specified column name
func (r *ImportedCommonDataRow) GetData(columnName string) string {
	return r.rowData[columnName]
}

// HasNext returns whether the iterator does not reach the end
func (t *ImportedCommonDataRowIterator) HasNext() bool {
	return t.innerIterator.HasNext()
}

// CurrentRowId returns current row id
func (t *ImportedCommonDataRowIterator) CurrentRowId() string {
	return t.innerIterator.CurrentRowId()
}

// Next returns the next common data row
func (t *ImportedCommonDataRowIterator) Next() CommonDataRow {
	importedRow := t.innerIterator.Next()

	if importedRow == nil {
		return nil
	}

	rowData := make(map[string]string, len(t.commonDataTable.dataColumnIndexes))

	for column, columnIndex := range t.commonDataTable.dataColumnIndexes {
		if columnIndex < 0 || columnIndex >= importedRow.ColumnCount() {
			continue
		}

		value := importedRow.GetData(columnIndex)
		rowData[column] = value
	}

	return &ImportedCommonDataRow{
		rowData: rowData,
	}
}

// CreateNewImportedCommonDataTable returns common data table from imported data table
func CreateNewImportedCommonDataTable(dataTable ImportedDataTable) *ImportedCommonDataTable {
	headerLineItems := dataTable.HeaderColumnNames()
	dataColumnIndexes := make(map[string]int, len(headerLineItems))

	for i := 0; i < len(headerLineItems); i++ {
		dataColumnIndexes[headerLineItems[i]] = i
	}

	return &ImportedCommonDataTable{
		innerDataTable:    dataTable,
		dataColumnIndexes: dataColumnIndexes,
	}
}
