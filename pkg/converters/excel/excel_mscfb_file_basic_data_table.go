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
	return fmt.Sprintf("sheet#%d-row#%d", t.currentSheetIndex, t.currentRowIndexInSheet)
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
