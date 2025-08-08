package datatable

// SubBasicDataTable defines the structure of sub basic data table
type SubBasicDataTable struct {
	baseTable BasicDataTable
	fromIndex int
	toIndex   int
}

// SubBasicDataTableRowIterator defines the structure of sub basic data table row iterator
type SubBasicDataTableRowIterator struct {
	dataTable     *SubBasicDataTable
	innerIterator BasicDataTableRowIterator
	currentIndex  int
}

// DataRowCount returns the total count of data row
func (t *SubBasicDataTable) DataRowCount() int {
	return t.toIndex - t.fromIndex
}

// HeaderColumnNames returns the header column name list
func (t *SubBasicDataTable) HeaderColumnNames() []string {
	return t.baseTable.HeaderColumnNames()
}

// DataRowIterator returns the iterator of data row
func (t *SubBasicDataTable) DataRowIterator() BasicDataTableRowIterator {
	innerIterator := t.baseTable.DataRowIterator()
	currentIndex := -1

	// skip rows until reaching the fromIndex
	for currentIndex = -1; currentIndex < t.fromIndex-1 && innerIterator.HasNext(); currentIndex++ {
		innerIterator.Next()
	}

	return &SubBasicDataTableRowIterator{
		dataTable:     t,
		innerIterator: innerIterator,
		currentIndex:  currentIndex,
	}
}

// HasNext returns whether the iterator does not reach the end
func (t *SubBasicDataTableRowIterator) HasNext() bool {
	return t.currentIndex+1 < t.dataTable.toIndex && t.innerIterator.HasNext()
}

// CurrentRowId returns current row id
func (t *SubBasicDataTableRowIterator) CurrentRowId() string {
	return t.innerIterator.CurrentRowId()
}

// Next returns the next basic data row
func (t *SubBasicDataTableRowIterator) Next() BasicDataTableRow {
	if t.currentIndex+1 >= t.dataTable.toIndex {
		return nil
	}

	t.currentIndex++
	return t.innerIterator.Next()
}

// CreateSubBasicTable returns a sub basic data table that references a portion of the original table
func CreateSubBasicTable(dataTable BasicDataTable, fromIndex, toIndex int) *SubBasicDataTable {
	if fromIndex < 0 {
		fromIndex = 0
	}

	if fromIndex > dataTable.DataRowCount() {
		fromIndex = dataTable.DataRowCount()
	}

	if toIndex > dataTable.DataRowCount() {
		toIndex = dataTable.DataRowCount()
	}

	if toIndex < fromIndex {
		toIndex = fromIndex
	}

	return &SubBasicDataTable{
		baseTable: dataTable,
		fromIndex: fromIndex,
		toIndex:   toIndex,
	}
}
