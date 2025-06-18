package excel

import (
	"bytes"
	"fmt"

	"github.com/extrame/xls"

	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

// ExcelMSCFBFileBasicDataTable defines the structure of excel (microsoft compound file binary) file data table
type ExcelMSCFBFileBasicDataTable struct {
	workbook              *xls.WorkBook
	headerLineColumnNames []string
}

// ExcelMSCFBFileBasicDataTableRow defines the structure of excel (microsoft compound file binary) file data table row
type ExcelMSCFBFileBasicDataTableRow struct {
	sheet    *xls.WorkSheet
	rowIndex int
}

// ExcelMSCFBFileBasicDataTableRowIterator defines the structure of excel (microsoft compound file binary) file data table row iterator
type ExcelMSCFBFileBasicDataTableRowIterator struct {
	dataTable              *ExcelMSCFBFileBasicDataTable
	currentSheetIndex      int
	currentRowIndexInSheet uint16
}

// DataRowCount returns the total count of data row
func (t *ExcelMSCFBFileBasicDataTable) DataRowCount() int {
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
func (t *ExcelMSCFBFileBasicDataTable) HeaderColumnNames() []string {
	return t.headerLineColumnNames
}

// DataRowIterator returns the iterator of data row
func (t *ExcelMSCFBFileBasicDataTable) DataRowIterator() datatable.BasicDataTableRowIterator {
	return &ExcelMSCFBFileBasicDataTableRowIterator{
		dataTable:              t,
		currentSheetIndex:      0,
		currentRowIndexInSheet: 0,
	}
}

// ColumnCount returns the total count of column in this data row
func (r *ExcelMSCFBFileBasicDataTableRow) ColumnCount() int {
	row := r.sheet.Row(r.rowIndex)
	return row.LastCol() + 1
}

// GetData returns the data in the specified column index
func (r *ExcelMSCFBFileBasicDataTableRow) GetData(columnIndex int) string {
	row := r.sheet.Row(r.rowIndex)
	return row.Col(columnIndex)
}

// HasNext returns whether the iterator does not reach the end
func (t *ExcelMSCFBFileBasicDataTableRowIterator) HasNext() bool {
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
func (t *ExcelMSCFBFileBasicDataTableRowIterator) CurrentRowId() string {
	return fmt.Sprintf("sheet#%d-row#%d", t.currentSheetIndex, t.currentRowIndexInSheet)
}

// Next returns the next basic data row
func (t *ExcelMSCFBFileBasicDataTableRowIterator) Next() datatable.BasicDataTableRow {
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

	return &ExcelMSCFBFileBasicDataTableRow{
		sheet:    currentSheet,
		rowIndex: int(t.currentRowIndexInSheet),
	}
}

// CreateNewExcelMSCFBFileBasicDataTable returns excel (microsoft compound file binary) data table by file binary data
func CreateNewExcelMSCFBFileBasicDataTable(data []byte) (datatable.BasicDataTable, error) {
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

	return &ExcelMSCFBFileBasicDataTable{
		workbook:              workbook,
		headerLineColumnNames: headerRowItems,
	}, nil
}
