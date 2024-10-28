package excel

import (
	"bytes"
	"fmt"

	"github.com/extrame/xls"

	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

// ExcelFileImportedDataTable defines the structure of excel file data table
type ExcelFileImportedDataTable struct {
	workbook              *xls.WorkBook
	headerLineColumnNames []string
}

// ExcelFileDataRow defines the structure of excel file data table row
type ExcelFileDataRow struct {
	sheet    *xls.WorkSheet
	rowIndex int
}

// ExcelFileDataRowIterator defines the structure of excel file data table row iterator
type ExcelFileDataRowIterator struct {
	dataTable              *ExcelFileImportedDataTable
	currentSheetIndex      int
	currentRowIndexInSheet uint16
}

// DataRowCount returns the total count of data row
func (t *ExcelFileImportedDataTable) DataRowCount() int {
	totalDataRowCount := 0

	for i := 0; i < t.workbook.NumSheets(); i++ {
		sheet := t.workbook.GetSheet(i)

		if sheet.MaxRow < 1 {
			continue
		}

		totalDataRowCount += int(sheet.MaxRow)
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
		currentSheetIndex:      0,
		currentRowIndexInSheet: 0,
	}
}

// ColumnCount returns the total count of column in this data row
func (r *ExcelFileDataRow) ColumnCount() int {
	row := r.sheet.Row(r.rowIndex)
	return row.LastCol() + 1
}

// GetData returns the data in the specified column index
func (r *ExcelFileDataRow) GetData(columnIndex int) string {
	row := r.sheet.Row(r.rowIndex)
	return row.Col(columnIndex)
}

// HasNext returns whether the iterator does not reach the end
func (t *ExcelFileDataRowIterator) HasNext() bool {
	workbook := t.dataTable.workbook

	if t.currentSheetIndex >= workbook.NumSheets() {
		return false
	}

	currentSheet := workbook.GetSheet(t.currentSheetIndex)

	if t.currentRowIndexInSheet+1 <= currentSheet.MaxRow {
		return true
	}

	for i := t.currentSheetIndex + 1; i < workbook.NumSheets(); i++ {
		sheet := workbook.GetSheet(i)

		if sheet.MaxRow < 1 {
			continue
		}

		return true
	}

	return false
}

// CurrentRowId returns current index
func (t *ExcelFileDataRowIterator) CurrentRowId() string {
	return fmt.Sprintf("table#%d-row#%d", t.currentSheetIndex, t.currentRowIndexInSheet)
}

// Next returns the next imported data row
func (t *ExcelFileDataRowIterator) Next() datatable.ImportedDataRow {
	workbook := t.dataTable.workbook
	currentRowIndexInTable := t.currentRowIndexInSheet

	for i := t.currentSheetIndex; i < workbook.NumSheets(); i++ {
		sheet := workbook.GetSheet(i)

		if currentRowIndexInTable+1 <= sheet.MaxRow {
			t.currentRowIndexInSheet++
			currentRowIndexInTable = t.currentRowIndexInSheet
			break
		}

		t.currentSheetIndex++
		t.currentRowIndexInSheet = 0
		currentRowIndexInTable = 0
	}

	if t.currentSheetIndex >= workbook.NumSheets() {
		return nil
	}

	currentSheet := workbook.GetSheet(t.currentSheetIndex)

	if t.currentRowIndexInSheet > currentSheet.MaxRow {
		return nil
	}

	return &ExcelFileDataRow{
		sheet:    currentSheet,
		rowIndex: int(t.currentRowIndexInSheet),
	}
}

// CreateNewExcelFileImportedDataTable returns excel xls data table by file binary data
func CreateNewExcelFileImportedDataTable(data []byte) (*ExcelFileImportedDataTable, error) {
	reader := bytes.NewReader(data)
	workbook, err := xls.OpenReader(reader, "")

	if err != nil {
		return nil, err
	}

	var headerRowItems []string

	for i := 0; i < workbook.NumSheets(); i++ {
		sheet := workbook.GetSheet(i)

		if sheet.MaxRow < 0 {
			continue
		}

		row := sheet.Row(0)

		if i == 0 {
			for j := 0; j <= row.LastCol(); j++ {
				headerItem := row.Col(j)

				if headerItem == "" {
					break
				}

				headerRowItems = append(headerRowItems, headerItem)
			}
		} else {
			for j := 0; j <= min(row.LastCol(), len(headerRowItems)-1); j++ {
				headerItem := row.Col(j)

				if headerItem != headerRowItems[j] {
					return nil, errs.ErrFieldsInMultiTableAreDifferent
				}
			}
		}
	}

	return &ExcelFileImportedDataTable{
		workbook:              workbook,
		headerLineColumnNames: headerRowItems,
	}, nil
}
