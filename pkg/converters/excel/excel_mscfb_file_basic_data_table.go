package excel

import (
	"bytes"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"time"

	"github.com/extrame/xls"
	"github.com/xuri/excelize/v2"

	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

// ExcelMSCFBFileBasicDataTable defines the structure of excel (microsoft compound file binary) file data table
type ExcelMSCFBFileBasicDataTable struct {
	workbook              *xls.WorkBook
	sheets                []*xls.WorkSheet
	headerLineColumnNames []string
	hasTitleLine          bool
}

// ExcelMSCFBFileBasicDataTableRow defines the structure of excel (microsoft compound file binary) file data table row
type ExcelMSCFBFileBasicDataTableRow struct {
	workbook *xls.WorkBook
	sheet    *xls.WorkSheet
	rowIndex int
}

// ExcelMSCFBFileBasicDataTableRowIterator defines the structure of excel (microsoft compound file binary) file data table row iterator
type ExcelMSCFBFileBasicDataTableRowIterator struct {
	dataTable              *ExcelMSCFBFileBasicDataTable
	currentSheetIndex      int
	currentRowIndexInSheet int
}

// DataRowCount returns the total count of data row
func (t *ExcelMSCFBFileBasicDataTable) DataRowCount() int {
	totalDataRowCount := 0

	for i := 0; i < len(t.sheets); i++ {
		sheet := t.sheets[i]

		if sheet == nil {
			continue
		}

		if t.hasTitleLine {
			if sheet.MaxRow < 1 {
				continue
			}

			totalDataRowCount += int(sheet.MaxRow)
		} else {
			if sheet.MaxRow <= 0 && sheet.Row(0) == nil {
				continue
			}

			totalDataRowCount += int(sheet.MaxRow) + 1
		}
	}

	return totalDataRowCount
}

// HeaderColumnNames returns the header column name list
func (t *ExcelMSCFBFileBasicDataTable) HeaderColumnNames() []string {
	if !t.hasTitleLine {
		return nil
	}

	return t.headerLineColumnNames
}

// DataRowIterator returns the iterator of data row
func (t *ExcelMSCFBFileBasicDataTable) DataRowIterator() datatable.BasicDataTableRowIterator {
	startIndex := -1

	if t.hasTitleLine {
		startIndex = 0
	}

	return &ExcelMSCFBFileBasicDataTableRowIterator{
		dataTable:              t,
		currentSheetIndex:      0,
		currentRowIndexInSheet: startIndex,
	}
}

// ColumnCount returns the total count of column in this data row
func (r *ExcelMSCFBFileBasicDataTableRow) ColumnCount() int {
	row := r.sheet.Row(r.rowIndex)
	return row.LastCol()
}

// GetData returns the data in the specified column index
func (r *ExcelMSCFBFileBasicDataTableRow) GetData(columnIndex int) string {
	row := r.sheet.Row(r.rowIndex)
	return getExcelMSCFBFileCellData(r.workbook, row, columnIndex)
}

// HasNext returns whether the iterator does not reach the end
func (t *ExcelMSCFBFileBasicDataTableRowIterator) HasNext() bool {
	sheets := t.dataTable.sheets

	if t.currentSheetIndex >= len(sheets) {
		return false
	}

	currentSheet := sheets[t.currentSheetIndex]

	if t.currentRowIndexInSheet+1 <= int(currentSheet.MaxRow) && currentSheet.Row(t.currentRowIndexInSheet+1) != nil {
		return true
	}

	for i := t.currentSheetIndex + 1; i < len(sheets); i++ {
		sheet := sheets[i]

		if t.dataTable.hasTitleLine {
			if sheet.MaxRow < 1 {
				continue
			}
		} else {
			if sheet.MaxRow <= 0 && sheet.Row(0) == nil {
				continue
			}
		}

		return true
	}

	return false
}

// CurrentRowId returns current index
func (t *ExcelMSCFBFileBasicDataTableRowIterator) CurrentRowId() string {
	if t.currentRowIndexInSheet < 0 || t.currentSheetIndex >= len(t.dataTable.sheets) {
		return ""
	}

	if len(t.dataTable.sheets) > 1 {
		return fmt.Sprintf("sheet:%d-row:%d", t.currentSheetIndex+1, t.currentRowIndexInSheet+1)
	} else {
		return fmt.Sprintf("row:%d", t.currentRowIndexInSheet+1)
	}
}

// Next returns the next basic data row
func (t *ExcelMSCFBFileBasicDataTableRowIterator) Next() datatable.BasicDataTableRow {
	sheets := t.dataTable.sheets

	for i := t.currentSheetIndex; i < len(sheets); i++ {
		sheet := sheets[i]

		if t.currentRowIndexInSheet+1 <= int(sheet.MaxRow) && sheet.Row(t.currentRowIndexInSheet+1) != nil {
			t.currentRowIndexInSheet++
			break
		}

		t.currentSheetIndex++

		if t.dataTable.hasTitleLine {
			t.currentRowIndexInSheet = 0
		} else {
			t.currentRowIndexInSheet = -1
		}
	}

	if t.currentSheetIndex >= len(sheets) {
		return nil
	}

	currentSheet := sheets[t.currentSheetIndex]

	if t.currentRowIndexInSheet > int(currentSheet.MaxRow) || currentSheet.Row(t.currentRowIndexInSheet) == nil {
		return nil
	}

	return &ExcelMSCFBFileBasicDataTableRow{
		workbook: t.dataTable.workbook,
		sheet:    currentSheet,
		rowIndex: int(t.currentRowIndexInSheet),
	}
}

// CreateNewExcelMSCFBFileBasicDataTable returns excel (microsoft compound file binary) data table by file binary data
func CreateNewExcelMSCFBFileBasicDataTable(data []byte, hasTitleLine bool) (datatable.BasicDataTable, error) {
	reader := bytes.NewReader(data)
	workbook, err := xls.OpenReader(reader, "")

	if err != nil {
		return nil, err
	}

	var firstRowItems []string
	var sheets []*xls.WorkSheet

	for i := 0; i < workbook.NumSheets(); i++ {
		sheet := workbook.GetSheet(i)

		if sheet == nil {
			continue
		}

		if sheet.MaxRow <= 0 && sheet.Row(0) == nil {
			continue
		}

		row := sheet.Row(0)

		if row == nil {
			continue
		}

		if hasTitleLine {
			if i == 0 {
				// row.LastCol() returns "colMac" in the "Row" struct, that is an unsigned integer that specifies the one-based index of the last column.
				// But row.FirstCol() returns "colMic" in the "Row" struct, that is an unsigned integer that specifies the zero-based index of the first column.
				// Reference: https://learn.microsoft.com/en-us/openspecs/office_file_formats/ms-xls/4aab09eb-49ed-4d01-a3b1-1d726247d3c2
				for j := 0; j < row.LastCol(); j++ {
					headerItem := getExcelMSCFBFileCellData(workbook, row, j)

					if headerItem == "" {
						break
					}

					firstRowItems = append(firstRowItems, headerItem)
				}
			} else {
				for j := 0; j < min(row.LastCol(), len(firstRowItems)); j++ {
					headerItem := getExcelMSCFBFileCellData(workbook, row, j)

					if headerItem != firstRowItems[j] {
						return nil, errs.ErrFieldsInMultiTableAreDifferent
					}
				}
			}
		}

		sheets = append(sheets, sheet)
	}

	var headerLineColumnNames []string = nil

	if hasTitleLine {
		headerLineColumnNames = firstRowItems
	}

	return &ExcelMSCFBFileBasicDataTable{
		workbook:              workbook,
		sheets:                sheets,
		headerLineColumnNames: headerLineColumnNames,
		hasTitleLine:          hasTitleLine,
	}, nil
}

// CreateNewExcelMSCFBFileBasicDataTables returns excel (microsoft compound file binary) data tables by file binary data, one per worksheet
func CreateNewExcelMSCFBFileBasicDataTables(data []byte, hasTitleLine bool) ([]datatable.BasicDataTable, error) {
	reader := bytes.NewReader(data)
	workbook, err := xls.OpenReader(reader, "")

	if err != nil {
		return nil, err
	}

	var dataTables []datatable.BasicDataTable

	for i := 0; i < workbook.NumSheets(); i++ {
		sheet := workbook.GetSheet(i)

		if sheet == nil {
			continue
		}

		if sheet.MaxRow <= 0 && sheet.Row(0) == nil {
			continue
		}

		row := sheet.Row(0)

		if row == nil {
			continue
		}

		var headerLineColumnNames []string = nil

		if hasTitleLine {
			// row.LastCol() returns "colMac" in the "Row" struct, that is an unsigned integer that specifies the one-based index of the last column.
			// But row.FirstCol() returns "colMic" in the "Row" struct, that is an unsigned integer that specifies the zero-based index of the first column.
			// Reference: https://learn.microsoft.com/en-us/openspecs/office_file_formats/ms-xls/4aab09eb-49ed-4d01-a3b1-1d726247d3c2
			for j := 0; j < row.LastCol(); j++ {
				headerItem := getExcelMSCFBFileCellData(workbook, row, j)

				if headerItem == "" {
					break
				}

				headerLineColumnNames = append(headerLineColumnNames, headerItem)
			}
		}

		dataTables = append(dataTables, &ExcelMSCFBFileBasicDataTable{
			workbook:              workbook,
			sheets:                []*xls.WorkSheet{sheet},
			headerLineColumnNames: headerLineColumnNames,
			hasTitleLine:          hasTitleLine,
		})
	}

	return dataTables, nil
}

func getExcelMSCFBFileCellData(workbook *xls.WorkBook, row *xls.Row, columnIndex int) string {
	if row == nil || columnIndex < 0 || columnIndex > math.MaxUint16 {
		return ""
	}

	var cellData string
	var panicked bool

	func() {
		defer func() {
			if recover() != nil {
				panicked = true
			}
		}()
		cellData = row.Col(columnIndex)
	}()

	if !panicked {
		return cellData
	}

	// extrame/xls (*NumberCol).String dereferences wb.Formats[fNo] without checking for nil
	// when a cell uses a built-in format ID that has no explicit FORMAT record.
	return getExcelMSCFBFileCellFallbackData(workbook, row, columnIndex)
}

func getExcelMSCFBFileCellFallbackData(workbook *xls.WorkBook, row *xls.Row, columnIndex int) string {
	rowValue := reflect.ValueOf(row)
	if !rowValue.IsValid() || rowValue.Kind() != reflect.Ptr || rowValue.IsNil() {
		return ""
	}

	colsField := rowValue.Elem().FieldByName("cols")
	if !colsField.IsValid() || colsField.Kind() != reflect.Map {
		return ""
	}

	colValue := colsField.MapIndex(reflect.ValueOf(uint16(columnIndex)))
	if !colValue.IsValid() || colValue.IsNil() {
		return ""
	}

	if colValue.Kind() == reflect.Interface {
		colValue = colValue.Elem()
	}

	if colValue.Kind() != reflect.Ptr || colValue.IsNil() {
		return ""
	}

	colStruct := colValue.Elem()
	if colStruct.Type() != reflect.TypeOf(xls.NumberCol{}) {
		return ""
	}

	xfIndexField := colStruct.FieldByName("Index")
	floatField := colStruct.FieldByName("Float")
	if !xfIndexField.IsValid() || !xfIndexField.CanUint() || !floatField.IsValid() || !floatField.CanFloat() {
		return ""
	}

	return formatExcelMSCFBNumberCell(workbook, uint16(xfIndexField.Uint()), floatField.Float())
}

func formatExcelMSCFBNumberCell(workbook *xls.WorkBook, xfIndex uint16, floatVal float64) string {
	if workbook != nil && int(xfIndex) < len(workbook.Xfs) {
		formatIndex, ok := getExcelMSCBFXfFormatIndex(workbook.Xfs[xfIndex])
		if ok && isBuiltInExcelDateFormat(formatIndex) {
			isDate1904 := isExcelMSCFBWorkbookDate1904(workbook)
			if dateTime, err := excelize.ExcelDateToTime(floatVal, isDate1904); err == nil {
				return dateTime.Format(time.RFC3339)
			}
		}
	}

	return strconv.FormatFloat(floatVal, 'f', -1, 64)
}

func getExcelMSCBFXfFormatIndex(xf interface{}) (uint16, bool) {
	switch v := xf.(type) {
	case *xls.Xf8:
		if v != nil {
			return v.Format, true
		}
	case *xls.Xf5:
		if v != nil {
			return v.Format, true
		}
	}

	return 0, false
}

func isBuiltInExcelDateFormat(formatIndex uint16) bool {
	return (14 <= formatIndex && formatIndex <= 17) ||
		formatIndex == 22 ||
		(27 <= formatIndex && formatIndex <= 36) ||
		(50 <= formatIndex && formatIndex <= 58)
}

func isExcelMSCFBWorkbookDate1904(workbook *xls.WorkBook) bool {
	if workbook == nil {
		return false
	}

	wbValue := reflect.ValueOf(workbook)
	if wbValue.Kind() != reflect.Ptr || wbValue.IsNil() {
		return false
	}

	dateModeField := wbValue.Elem().FieldByName("dateMode")
	if !dateModeField.IsValid() || !dateModeField.CanUint() {
		return false
	}

	return dateModeField.Uint() == 1
}
