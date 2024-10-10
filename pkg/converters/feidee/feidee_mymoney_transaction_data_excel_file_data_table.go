package feidee

import (
	"bytes"
	"time"

	"github.com/shakinm/xlsReader/xls"

	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// feideeMymoneyTransactionExcelFileDataTable defines the structure of feidee mymoney transaction plain text data table
type feideeMymoneyTransactionExcelFileDataTable struct {
	workbook              *xls.Workbook
	headerLineColumnNames []string
}

// feideeMymoneyTransactionExcelFileDataRow defines the structure of feidee mymoney transaction plain text data row
type feideeMymoneyTransactionExcelFileDataRow struct {
	sheet    *xls.Sheet
	rowIndex int
}

// feideeMymoneyTransactionExcelFileDataRowIterator defines the structure of feidee mymoney transaction plain text data row iterator
type feideeMymoneyTransactionExcelFileDataRowIterator struct {
	dataTable              *feideeMymoneyTransactionExcelFileDataTable
	currentTableIndex      int
	currentRowIndexInTable int
}

// DataRowCount returns the total count of data row
func (t *feideeMymoneyTransactionExcelFileDataTable) DataRowCount() int {
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

// HeaderLineColumnNames returns the header column name list
func (t *feideeMymoneyTransactionExcelFileDataTable) HeaderLineColumnNames() []string {
	return t.headerLineColumnNames
}

// DataRowIterator returns the iterator of data row
func (t *feideeMymoneyTransactionExcelFileDataTable) DataRowIterator() datatable.ImportedDataRowIterator {
	return &feideeMymoneyTransactionExcelFileDataRowIterator{
		dataTable:              t,
		currentTableIndex:      0,
		currentRowIndexInTable: 0,
	}
}

// IsValid returns whether this row contains valid data for importing
func (r *feideeMymoneyTransactionExcelFileDataRow) IsValid() bool {
	return true
}

// ColumnCount returns the total count of column in this data row
func (r *feideeMymoneyTransactionExcelFileDataRow) ColumnCount() int {
	row, err := r.sheet.GetRow(r.rowIndex)

	if err != nil {
		return 0
	}

	return len(row.GetCols())
}

// GetData returns the data in the specified column index
func (r *feideeMymoneyTransactionExcelFileDataRow) GetData(columnIndex int) string {
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

// GetTime returns the time in the specified column index
func (r *feideeMymoneyTransactionExcelFileDataRow) GetTime(columnIndex int, timezoneOffset int16) (time.Time, error) {
	str := r.GetData(columnIndex)

	if utils.IsValidLongDateTimeFormat(str) {
		return utils.ParseFromLongDateTime(str, timezoneOffset)
	}

	if utils.IsValidLongDateTimeWithoutSecondFormat(str) {
		return utils.ParseFromLongDateTimeWithoutSecond(str, timezoneOffset)
	}

	if utils.IsValidLongDateFormat(str) {
		return utils.ParseFromLongDateTimeWithoutSecond(str+" 00:00", timezoneOffset)
	}

	return time.Unix(0, 0), errs.ErrTransactionTimeInvalid
}

// GetTimezoneOffset returns the time zone offset in the specified column index
func (r *feideeMymoneyTransactionExcelFileDataRow) GetTimezoneOffset(columnIndex int) (*time.Location, error) {
	return nil, errs.ErrNotSupported
}

// HasNext returns whether the iterator does not reach the end
func (t *feideeMymoneyTransactionExcelFileDataRowIterator) HasNext() bool {
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
func (t *feideeMymoneyTransactionExcelFileDataRowIterator) Next(ctx core.Context, user *models.User) datatable.ImportedDataRow {
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

	return &feideeMymoneyTransactionExcelFileDataRow{
		sheet:    &currentSheet,
		rowIndex: t.currentRowIndexInTable,
	}
}

func createNewFeideeMymoneyTransactionExcelFileDataTable(data []byte) (*feideeMymoneyTransactionExcelFileDataTable, error) {
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

	return &feideeMymoneyTransactionExcelFileDataTable{
		workbook:              &workbook,
		headerLineColumnNames: headerRowItems,
	}, nil
}
