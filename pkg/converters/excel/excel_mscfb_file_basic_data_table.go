package excel

import (
	"bytes"
	"fmt"
	"reflect"
	"time"
	"unsafe"

	"github.com/extrame/xls"

	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

// The layouts of built-in formats are Go time layouts instead of Excel layouts, because the formatter of extrame/xls passes them through as is,
// but it formats the month, hour, minute and second of Excel layouts incorrectly.

// builtinDateTimeFormatLayout is the same output as extrame/xls produces for the date cells stored in the RK records
const builtinDateTimeFormatLayout = time.RFC3339

// builtinTimeFormatLayout is the layout for the time cells
const builtinTimeFormatLayout = "15:04:05"

// isBuiltinDateTimeFormatId returns whether the built-in number format id is a date or datetime format (the same ids as extrame/xls treats as dates)
func isBuiltinDateTimeFormatId(formatId uint16) bool {
	return 14 <= formatId && formatId <= 17 || formatId == 22 || 27 <= formatId && formatId <= 36 || 50 <= formatId && formatId <= 58
}

// isBuiltinTimeFormatId returns whether the built-in number format id is a time format
func isBuiltinTimeFormatId(formatId uint16) bool {
	return 18 <= formatId && formatId <= 21 || 45 <= formatId && formatId <= 47
}

// ExcelMSCFBFileBasicDataTable defines the structure of excel (microsoft compound file binary) file data table
type ExcelMSCFBFileBasicDataTable struct {
	sheets                []*xls.WorkSheet
	headerLineColumnNames []string
	hasTitleLine          bool
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
	return row.Col(columnIndex)
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

	fillBuiltinNumberFormats(workbook)

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
					headerItem := row.Col(j)

					if headerItem == "" {
						break
					}

					firstRowItems = append(firstRowItems, headerItem)
				}
			} else {
				for j := 0; j < min(row.LastCol(), len(firstRowItems)); j++ {
					headerItem := row.Col(j)

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

	fillBuiltinNumberFormats(workbook)

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
				headerItem := row.Col(j)

				if headerItem == "" {
					break
				}

				headerLineColumnNames = append(headerLineColumnNames, headerItem)
			}
		}

		dataTables = append(dataTables, &ExcelMSCFBFileBasicDataTable{
			sheets:                []*xls.WorkSheet{sheet},
			headerLineColumnNames: headerLineColumnNames,
			hasTitleLine:          hasTitleLine,
		})
	}

	return dataTables, nil
}

// fillBuiltinNumberFormats adds the formats which are used by cells but not stored in the workbook
// extrame/xls only reads the custom formats from the file, so a numeric cell with a built-in format id (e.g. 14) causes a nil pointer dereference
func fillBuiltinNumberFormats(workbook *xls.WorkBook) {
	for _, xf := range workbook.Xfs {
		var formatId *uint16

		switch x := xf.(type) {
		case *xls.Xf5:
			formatId = &x.Format
		case *xls.Xf8:
			formatId = &x.Format
		default:
			continue
		}

		if *formatId == 0 || *formatId >= 164 || workbook.Formats[*formatId] != nil {
			continue
		}

		var layout string

		if isBuiltinDateTimeFormatId(*formatId) {
			layout = builtinDateTimeFormatLayout
		} else if isBuiltinTimeFormatId(*formatId) {
			layout = builtinTimeFormatLayout
		} else {
			// the other built-in formats are numbers or text, so they are shown as general
			*formatId = 0
			continue
		}

		format := &xls.Format{}
		format.Head.Index = *formatId

		// the layout field is not exported
		layoutField := reflect.ValueOf(format).Elem().FieldByName("str")
		reflect.NewAt(layoutField.Type(), unsafe.Pointer(layoutField.UnsafeAddr())).Elem().SetString(layout)

		workbook.Formats[*formatId] = format
	}
}
