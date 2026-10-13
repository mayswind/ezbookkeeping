package datatable

// SubBasicDataParseFunc defines the function type for parsing data in sub basic data table
type SubBasicDataParseFunc func(data string) string

// SubBasicDataTable defines the structure of sub basic data table
type SubBasicDataTable struct {
	baseTable             BasicDataTable
	fromIndex             int
	toIndex               int
	headerLineColumnNames []string
	dataParseFunc         SubBasicDataParseFunc
}

// SubBasicDataTableRow defines the structure of sub basic data table row
type SubBasicDataTableRow struct {
	baseRow       BasicDataTableRow
	dataParseFunc SubBasicDataParseFunc
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
	return t.headerLineColumnNames
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

// ColumnCount returns the total count of column in this data row
func (t *SubBasicDataTableRow) ColumnCount() int {
	return t.baseRow.ColumnCount()
}

// GetData returns the data in the specified column index
func (t *SubBasicDataTableRow) GetData(columnIndex int) string {
	data := t.baseRow.GetData(columnIndex)

	if t.dataParseFunc != nil {
		return t.dataParseFunc(data)
	}

	return data
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

	return &SubBasicDataTableRow{
		baseRow:       t.innerIterator.Next(),
		dataParseFunc: t.dataTable.dataParseFunc,
	}
}

// CreateSubBasicTable returns a sub basic data table that references a portion `[fromIndex, toIndex)` of the original table
func CreateSubBasicTable(dataTable BasicDataTable, fromIndex int, toIndex int, headerInFromRow bool) *SubBasicDataTable {
	return createSubBasicTable(dataTable, fromIndex, toIndex, headerInFromRow, nil)
}

// CreateSubBasicTableWithParser returns a sub basic data table that references a portion `[fromIndex, toIndex)` of the original table with a data parser
func CreateSubBasicTableWithParser(dataTable BasicDataTable, fromIndex int, toIndex int, headerInFromRow bool, dataParseFunc SubBasicDataParseFunc) *SubBasicDataTable {
	return createSubBasicTable(dataTable, fromIndex, toIndex, headerInFromRow, dataParseFunc)
}

func createSubBasicTable(dataTable BasicDataTable, fromIndex int, toIndex int, headerInFromRow bool, dataParseFunc SubBasicDataParseFunc) *SubBasicDataTable {
	if fromIndex < 0 {
		fromIndex = 0
	}

	var headerLineColumnNames []string

	if headerInFromRow {
		iterator := dataTable.DataRowIterator()
		var headerRow BasicDataTableRow

		for i := 0; i <= fromIndex && iterator.HasNext(); i++ {
			headerRow = iterator.Next()
		}

		if headerRow != nil {
			for i := 0; i < headerRow.ColumnCount(); i++ {
				headerLineColumnNames = append(headerLineColumnNames, headerRow.GetData(i))
			}
		}

		fromIndex++
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
		baseTable:             dataTable,
		fromIndex:             fromIndex,
		toIndex:               toIndex,
		headerLineColumnNames: headerLineColumnNames,
		dataParseFunc:         dataParseFunc,
	}
}
