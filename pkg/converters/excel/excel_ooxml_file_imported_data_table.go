package excel

import (
	"bytes"
	"fmt"

	"github.com/xuri/excelize/v2"

	"github.com/mayswind/ezbookkeeping/pkg/converters/datatable"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

// excelOOXMLSheet defines the structure of excel (Office Open XML) file sheet
type excelOOXMLSheet struct {
	sheetName string
	allData   [][]string
}

// ExcelOOXMLFileImportedDataTable defines the structure of excel (Office Open XML) file data table
type ExcelOOXMLFileImportedDataTable struct {
	sheets                []*excelOOXMLSheet
	headerLineColumnNames []string
}

// ExcelOOXMLFileDataRow defines the structure of excel (Office Open XML) file data table row
type ExcelOOXMLFileDataRow struct {
	sheet    *excelOOXMLSheet
	rowData  []string
	rowIndex int
}

// ExcelOOXMLFileDataRowIterator defines the structure of excel (Office Open XML) file data table row iterator
type ExcelOOXMLFileDataRowIterator struct {
	dataTable              *ExcelOOXMLFileImportedDataTable
	currentSheetIndex      int
	currentRowIndexInSheet int
}

// DataRowCount returns the total count of data row
func (t *ExcelOOXMLFileImportedDataTable) DataRowCount() int {
	totalDataRowCount := 0

	for i := 0; i < len(t.sheets); i++ {
		sheet := t.sheets[i]

		if len(sheet.allData) < 1 {
			continue
		}

		totalDataRowCount += len(sheet.allData) - 1
	}

	return totalDataRowCount
}

// HeaderColumnNames returns the header column name list
func (t *ExcelOOXMLFileImportedDataTable) HeaderColumnNames() []string {
	return t.headerLineColumnNames
}

// DataRowIterator returns the iterator of data row
func (t *ExcelOOXMLFileImportedDataTable) DataRowIterator() datatable.ImportedDataRowIterator {
	return &ExcelOOXMLFileDataRowIterator{
		dataTable:              t,
		currentSheetIndex:      0,
		currentRowIndexInSheet: 0,
	}
}

// ColumnCount returns the total count of column in this data row
func (r *ExcelOOXMLFileDataRow) ColumnCount() int {
	return len(r.rowData)
}

// GetData returns the data in the specified column index
func (r *ExcelOOXMLFileDataRow) GetData(columnIndex int) string {
	if columnIndex < 0 || columnIndex >= len(r.rowData) {
		return ""
	}

	return r.rowData[columnIndex]
}

// HasNext returns whether the iterator does not reach the end
func (t *ExcelOOXMLFileDataRowIterator) HasNext() bool {
	sheets := t.dataTable.sheets

	if t.currentSheetIndex >= len(sheets) {
		return false
	}

	currentSheet := sheets[t.currentSheetIndex]

	if t.currentRowIndexInSheet+1 < len(currentSheet.allData) {
		return true
	}

	for i := t.currentSheetIndex + 1; i < len(sheets); i++ {
		sheet := sheets[i]

		if len(sheet.allData) <= 1 {
			continue
		}

		return true
	}

	return false
}

// CurrentRowId returns current index
func (t *ExcelOOXMLFileDataRowIterator) CurrentRowId() string {
	return fmt.Sprintf("table#%d-row#%d", t.currentSheetIndex, t.currentRowIndexInSheet)
}

// Next returns the next imported data row
func (t *ExcelOOXMLFileDataRowIterator) Next() datatable.ImportedDataRow {
	sheets := t.dataTable.sheets
	currentRowIndexInTable := t.currentRowIndexInSheet

	for i := t.currentSheetIndex; i < len(sheets); i++ {
		sheet := sheets[i]

		if currentRowIndexInTable+1 < len(sheet.allData) {
			t.currentRowIndexInSheet++
			currentRowIndexInTable = t.currentRowIndexInSheet
			break
		}

		t.currentSheetIndex++
		t.currentRowIndexInSheet = 0
		currentRowIndexInTable = 0
	}

	if t.currentSheetIndex >= len(sheets) {
		return nil
	}

	currentSheet := sheets[t.currentSheetIndex]

	if t.currentRowIndexInSheet >= len(currentSheet.allData) {
		return nil
	}

	return &ExcelOOXMLFileDataRow{
		sheet:    currentSheet,
		rowData:  currentSheet.allData[t.currentRowIndexInSheet],
		rowIndex: t.currentRowIndexInSheet,
	}
}

// CreateNewExcelOOXMLFileImportedDataTable returns excel (Office Open XML) data table by file binary data
func CreateNewExcelOOXMLFileImportedDataTable(data []byte) (*ExcelOOXMLFileImportedDataTable, error) {
	reader := bytes.NewReader(data)
	file, err := excelize.OpenReader(reader)

	defer file.Close()

	if err != nil {
		return nil, err
	}

	sheetNames := file.GetSheetList()
	var headerRowItems []string
	var sheets []*excelOOXMLSheet

	for i := 0; i < len(sheetNames); i++ {
		sheetName := sheetNames[i]
		allData, err := file.GetRows(sheetName)

		if err != nil {
			return nil, err
		}

		if allData == nil || len(allData) < 1 {
			continue
		}

		row := allData[0]

		if i == 0 {
			for j := 0; j < len(row); j++ {
				headerItem := row[j]

				if headerItem == "" {
					break
				}

				headerRowItems = append(headerRowItems, headerItem)
			}
		} else {
			for j := 0; j < min(len(row), len(headerRowItems)); j++ {
				headerItem := row[j]

				if headerItem != headerRowItems[j] {
					return nil, errs.ErrFieldsInMultiTableAreDifferent
				}
			}
		}

		sheets = append(sheets, &excelOOXMLSheet{
			sheetName: sheetName,
			allData:   allData,
		})
	}

	return &ExcelOOXMLFileImportedDataTable{
		sheets:                sheets,
		headerLineColumnNames: headerRowItems,
	}, nil
}
