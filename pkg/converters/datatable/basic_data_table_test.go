package datatable

type testBasicDataTable struct {
	headerColumns []string
	rows          []*testBasicDataTableRow
}

type testBasicDataTableRow struct {
	rowId      string
	rowColumns []string
}

type testBasicDataTableRowIterator struct {
	rows         []*testBasicDataTableRow
	currentIndex int
}

func (t *testBasicDataTable) HeaderColumnNames() []string {
	return t.headerColumns
}

func (t *testBasicDataTable) DataRowCount() int {
	return len(t.rows)
}

func (t *testBasicDataTable) DataRowIterator() BasicDataTableRowIterator {
	return &testBasicDataTableRowIterator{
		rows:         t.rows,
		currentIndex: -1,
	}
}

func (r *testBasicDataTableRow) ColumnCount() int {
	return len(r.rowColumns)
}

func (r *testBasicDataTableRow) GetData(columnIndex int) string {
	if columnIndex < 0 || columnIndex >= len(r.rowColumns) {
		return ""
	}

	return r.rowColumns[columnIndex]
}

func (t *testBasicDataTableRowIterator) HasNext() bool {
	return t.currentIndex+1 < len(t.rows)
}

func (t *testBasicDataTableRowIterator) CurrentRowId() string {
	if t.currentIndex >= len(t.rows) {
		return ""
	}

	return t.rows[t.currentIndex].rowId
}

func (t *testBasicDataTableRowIterator) Next() BasicDataTableRow {
	if t.currentIndex+1 >= len(t.rows) {
		return nil
	}

	t.currentIndex++
	row := t.rows[t.currentIndex]
	return row
}
