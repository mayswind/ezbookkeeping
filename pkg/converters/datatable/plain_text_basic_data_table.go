package datatable

import (
	"fmt"
)

// PlainTextDataTableLineData defines the structure of plain text data table line data
type PlainTextDataTableLineData struct {
	LineNumber int
	Items      []string
}

// PlainTextBasicDataTable defines the structure of plain text data table
type PlainTextBasicDataTable struct {
	allLines          []*PlainTextDataTableLineData
	hasTitleLine      bool
	hasRawLineNumbers bool
}

// PlainTextBasicDataTableRow defines the structure of plain text data table row
type PlainTextBasicDataTableRow struct {
	dataTable *PlainTextBasicDataTable
	allItems  []string
}

// PlainTextBasicDataTableRowIterator defines the structure of plain text data table row iterator
type PlainTextBasicDataTableRowIterator struct {
	dataTable    *PlainTextBasicDataTable
	currentIndex int
}

// DataRowCount returns the total count of data row
func (t *PlainTextBasicDataTable) DataRowCount() int {
	if len(t.allLines) < 1 {
		return 0
	}

	if t.hasTitleLine {
		return len(t.allLines) - 1
	} else {
		return len(t.allLines)
	}
}

// HeaderColumnNames returns the header column name list
func (t *PlainTextBasicDataTable) HeaderColumnNames() []string {
	if len(t.allLines) < 1 {
		return nil
	}

	if t.hasTitleLine {
		return t.allLines[0].Items
	} else {
		return nil
	}
}

// DataRowIterator returns the iterator of data row
func (t *PlainTextBasicDataTable) DataRowIterator() BasicDataTableRowIterator {
	startIndex := -1

	if t.hasTitleLine {
		startIndex = 0
	}

	return &PlainTextBasicDataTableRowIterator{
		dataTable:    t,
		currentIndex: startIndex,
	}
}

// ColumnCount returns the total count of column in this data row
func (r *PlainTextBasicDataTableRow) ColumnCount() int {
	return len(r.allItems)
}

// GetData returns the data in the specified column index
func (r *PlainTextBasicDataTableRow) GetData(columnIndex int) string {
	if columnIndex >= len(r.allItems) {
		return ""
	}

	return r.allItems[columnIndex]
}

// HasNext returns whether the iterator does not reach the end
func (t *PlainTextBasicDataTableRowIterator) HasNext() bool {
	return t.currentIndex+1 < len(t.dataTable.allLines)
}

// CurrentRowId returns current index
func (t *PlainTextBasicDataTableRowIterator) CurrentRowId() string {
	if t.currentIndex < 0 || t.currentIndex >= len(t.dataTable.allLines) {
		return ""
	}

	if t.dataTable.hasRawLineNumbers {
		return fmt.Sprintf("line:%d", t.dataTable.allLines[t.currentIndex].LineNumber)
	}

	if t.dataTable.hasTitleLine {
		return fmt.Sprintf("record:%d", t.currentIndex)
	} else {
		return fmt.Sprintf("record:%d", t.currentIndex+1)
	}
}

// Next returns the next basic data row
func (t *PlainTextBasicDataTableRowIterator) Next() BasicDataTableRow {
	if t.currentIndex+1 >= len(t.dataTable.allLines) {
		return nil
	}

	t.currentIndex++

	rowItems := t.dataTable.allLines[t.currentIndex].Items

	return &PlainTextBasicDataTableRow{
		dataTable: t.dataTable,
		allItems:  rowItems,
	}
}

// CreateNewPlainTextDataTableWithLineNumbers returns basic data table with all lines and raw line numbers
func CreateNewPlainTextDataTableWithLineNumbers(allLines []*PlainTextDataTableLineData, hasTitleLine bool) BasicDataTable {
	return &PlainTextBasicDataTable{
		allLines:          allLines,
		hasTitleLine:      hasTitleLine,
		hasRawLineNumbers: true,
	}
}

// CreateNewPlainTextDataTableWithoutLineNumbers returns basic data table with all lines without raw line numbers
func CreateNewPlainTextDataTableWithoutLineNumbers(allOriginalLines [][]string, hasTitleLine bool) BasicDataTable {
	allLines := make([]*PlainTextDataTableLineData, len(allOriginalLines))

	for i := 0; i < len(allOriginalLines); i++ {
		allLines[i] = &PlainTextDataTableLineData{
			LineNumber: i + 1,
			Items:      allOriginalLines[i],
		}
	}

	return &PlainTextBasicDataTable{
		allLines:          allLines,
		hasTitleLine:      hasTitleLine,
		hasRawLineNumbers: false,
	}
}
