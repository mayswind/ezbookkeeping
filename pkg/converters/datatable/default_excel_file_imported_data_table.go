package datatable

import (
	"bytes"

	"github.com/shakinm/xlsReader/xls"

	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

// DefaultExcelFileImportedDataTable defines the structure of default excel file data table
type DefaultExcelFileImportedDataTable struct {
	workbook              *xls.Workbook
	headerLineColumnNames []string
}

// DefaultExcelFileDataRow defines the structure of default excel file data table row
type DefaultExcelFileDataRow struct {
	sheet    *xls.Sheet
	rowIndex int
}

// DefaultExcelFileDataRowIterator defines the structure of default excel file data table row iterator
type DefaultExcelFileDataRowIterator struct {
	dataTable              *DefaultExcelFileImportedDataTable
	currentTableIndex      int
	currentRowIndexInTable int
}

// DataRowCount returns the total count of data row
func (t *DefaultExcelFileImportedDataTable) DataRowCount() int {
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
func (t *DefaultExcelFileImportedDataTable) HeaderColumnNames() []string {
	return t.headerLineColumnNames
}

// DataRowIterator returns the iterator of data row
func (t *DefaultExcelFileImportedDataTable) DataRowIterator() ImportedDataRowIterator {
	return &DefaultExcelFileDataRowIterator{
		dataTable:              t,
		currentTableIndex:      0,
		currentRowIndexInTable: 0,
	}
}

// ColumnCount returns the total count of column in this data row
func (r *DefaultExcelFileDataRow) ColumnCount() int {
	row, err := r.sheet.GetRow(r.rowIndex)

	if err != nil {
		return 0
	}

	return len(row.GetCols())
}

// GetData returns the data in the specified column index
func (r *DefaultExcelFileDataRow) GetData(columnIndex int) string {
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
func (t *DefaultExcelFileDataRowIterator) HasNext() bool {
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

// Next returns the next imported data row
func (t *DefaultExcelFileDataRowIterator) Next() ImportedDataRow {
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

	return &DefaultExcelFileDataRow{
		sheet:    &currentSheet,
		rowIndex: t.currentRowIndexInTable,
	}
}

// CreateNewDefaultExcelFileImportedDataTable returns default excel xls data table by file binary data
func CreateNewDefaultExcelFileImportedDataTable(data []byte) (*DefaultExcelFileImportedDataTable, error) {
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

	return &DefaultExcelFileImportedDataTable{
		workbook:              &workbook,
		headerLineColumnNames: headerRowItems,
	}, nil
}
