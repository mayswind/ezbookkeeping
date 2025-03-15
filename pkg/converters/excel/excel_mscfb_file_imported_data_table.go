package excel

import (
	"bytes"
	"fmt"

	"github.com/extrame/xls"

	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

// ExcelMSCFBFileImportedDataTable defines the structure of excel (microsoft compound file binary) file data table
type ExcelMSCFBFileImportedDataTable struct {
	workbook              *xls.WorkBook
	headerLineColumnNames []string
}

// ExcelMSCFBFileDataRow defines the structure of excel (microsoft compound file binary) file data table row
type ExcelMSCFBFileDataRow struct {
	sheet    *xls.WorkSheet
	rowIndex int
}

// ExcelMSCFBFileDataRowIterator defines the structure of excel (microsoft compound file binary) file data table row iterator
type ExcelMSCFBFileDataRowIterator struct {
	dataTable              *ExcelMSCFBFileImportedDataTable
	currentSheetIndex      int
	currentRowIndexInSheet uint16
}

// DataRowCount returns the total count of data row
func (t *ExcelMSCFBFileImportedDataTable) DataRowCount() int {
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
func (t *ExcelMSCFBFileImportedDataTable) HeaderColumnNames() []string {
	return t.headerLineColumnNames
}

// DataRowIterator returns the iterator of data row
func (t *ExcelMSCFBFileImportedDataTable) DataRowIterator() datatable.ImportedDataRowIterator {
	return &ExcelMSCFBFileDataRowIterator{
		dataTable:              t,
		currentSheetIndex:      0,
		currentRowIndexInSheet: 0,
	}
}

// ColumnCount returns the total count of column in this data row
func (r *ExcelMSCFBFileDataRow) ColumnCount() int {
	row := r.sheet.Row(r.rowIndex)
	return row.LastCol() + 1
}

// GetData returns the data in the specified column index
func (r *ExcelMSCFBFileDataRow) GetData(columnIndex int) string {
	row := r.sheet.Row(r.rowIndex)
	return row.Col(columnIndex)
}

// HasNext returns whether the iterator does not reach the end
func (t *ExcelMSCFBFileDataRowIterator) HasNext() bool {
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
func (t *ExcelMSCFBFileDataRowIterator) CurrentRowId() string {
	return fmt.Sprintf("table#%d-row#%d", t.currentSheetIndex, t.currentRowIndexInSheet)
}

// Next returns the next imported data row
func (t *ExcelMSCFBFileDataRowIterator) Next() datatable.ImportedDataRow {
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

	return &ExcelMSCFBFileDataRow{
		sheet:    currentSheet,
		rowIndex: int(t.currentRowIndexInSheet),
	}
}

// CreateNewExcelMSCFBFileImportedDataTable returns excel (microsoft compound file binary) data table by file binary data
func CreateNewExcelMSCFBFileImportedDataTable(data []byte) (*ExcelMSCFBFileImportedDataTable, error) {
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

		if row == nil {
			continue
		}

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

	return &ExcelMSCFBFileImportedDataTable{
		workbook:              workbook,
		headerLineColumnNames: headerRowItems,
	}, nil
}
