package excel

import (
	"bytes"
	"fmt"

	"github.com/shakinm/xlsReader/xls"

	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

// ExcelFileImportedDataTable defines the structure of excel file data table
type ExcelFileImportedDataTable struct {
	workbook              *xls.Workbook
	headerLineColumnNames []string
}

// ExcelFileDataRow defines the structure of excel file data table row
type ExcelFileDataRow struct {
	sheet    *xls.Sheet
	rowIndex int
}

// ExcelFileDataRowIterator defines the structure of excel file data table row iterator
type ExcelFileDataRowIterator struct {
	dataTable              *ExcelFileImportedDataTable
	currentTableIndex      int
	currentRowIndexInTable int
}

// DataRowCount returns the total count of data row
func (t *ExcelFileImportedDataTable) DataRowCount() int {
	allSheets := t.workbook.GetSheets()
	totalDataRowCount := 0

	for i := 0; i < len(allSheets); i++ {
		sheet := allSheets[i]

		if sheet.GetNumberRows() <= 1 {
			continue
		}

		totalDataRowCount += sheet.GetNumberRows() - 1
	}

	return totalDataRowCount
}

// HeaderColumnNames returns the header column name list
func (t *ExcelFileImportedDataTable) HeaderColumnNames() []string {
	return t.headerLineColumnNames
}

// DataRowIterator returns the iterator of data row
func (t *ExcelFileImportedDataTable) DataRowIterator() datatable.ImportedDataRowIterator {
	return &ExcelFileDataRowIterator{
		dataTable:              t,
		currentTableIndex:      0,
		currentRowIndexInTable: 0,
	}
}

// ColumnCount returns the total count of column in this data row
func (r *ExcelFileDataRow) ColumnCount() int {
	row, err := r.sheet.GetRow(r.rowIndex)

	if err != nil {
		return 0
	}

	return len(row.GetCols())
}

// GetData returns the data in the specified column index
func (r *ExcelFileDataRow) GetData(columnIndex int) string {
	row, err := r.sheet.GetRow(r.rowIndex)

	if err != nil {
		return ""
	}

	cell, err := row.GetCol(columnIndex)

	if err != nil {
		return ""
	}

	return cell.GetString()
}

// HasNext returns whether the iterator does not reach the end
func (t *ExcelFileDataRowIterator) HasNext() bool {
	allSheets := t.dataTable.workbook.GetSheets()

	if t.currentTableIndex >= len(allSheets) {
		return false
	}

	currentSheet := allSheets[t.currentTableIndex]

	if t.currentRowIndexInTable+1 < currentSheet.GetNumberRows() {
		return true
	}

	for i := t.currentTableIndex + 1; i < len(allSheets); i++ {
		sheet := allSheets[i]

		if sheet.GetNumberRows() <= 1 {
			continue
		}

		return true
	}

	return false
}

// CurrentRowId returns current index
func (t *ExcelFileDataRowIterator) CurrentRowId() string {
	return fmt.Sprintf("table#%d-row#%d", t.currentTableIndex, t.currentRowIndexInTable)
}

// Next returns the next imported data row
func (t *ExcelFileDataRowIterator) Next() datatable.ImportedDataRow {
	allSheets := t.dataTable.workbook.GetSheets()
	currentRowIndexInTable := t.currentRowIndexInTable

	for i := t.currentTableIndex; i < len(allSheets); i++ {
		sheet := allSheets[i]

		if currentRowIndexInTable+1 < sheet.GetNumberRows() {
			t.currentRowIndexInTable++
			currentRowIndexInTable = t.currentRowIndexInTable
			break
		}

		t.currentTableIndex++
		t.currentRowIndexInTable = 0
		currentRowIndexInTable = 0
	}

	if t.currentTableIndex >= len(allSheets) {
		return nil
	}

	currentSheet := allSheets[t.currentTableIndex]

	if t.currentRowIndexInTable >= currentSheet.GetNumberRows() {
		return nil
	}

	return &ExcelFileDataRow{
		sheet:    &currentSheet,
		rowIndex: t.currentRowIndexInTable,
	}
}

// CreateNewExcelFileImportedDataTable returns excel xls data table by file binary data
func CreateNewExcelFileImportedDataTable(data []byte) (*ExcelFileImportedDataTable, error) {
	reader := bytes.NewReader(data)
	workbook, err := xls.OpenReader(reader)

	if err != nil {
		return nil, err
	}

	allSheets := workbook.GetSheets()
	var headerRowItems []string

	for i := 0; i < len(allSheets); i++ {
		sheet := allSheets[i]

		if sheet.GetNumberRows() < 1 {
			continue
		}

		row, err := sheet.GetRow(0)

		if err != nil {
			return nil, err
		}

		cells := row.GetCols()

		if i == 0 {
			for j := 0; j < len(cells); j++ {
				headerItem := cells[j].GetString()

				if headerItem == "" {
					break
				}

				headerRowItems = append(headerRowItems, headerItem)
			}
		} else {
			for j := 0; j < min(len(cells), len(headerRowItems)); j++ {
				headerItem := cells[j].GetString()

				if headerItem != headerRowItems[j] {
					return nil, errs.ErrFieldsInMultiTableAreDifferent
				}
			}
		}
	}

	return &ExcelFileImportedDataTable{
		workbook:              &workbook,
		headerLineColumnNames: headerRowItems,
	}, nil
}
